// Package certificate provides code to create and write a self-signed TLS certificate
package certificate

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/spoonboy-io/dujour/internal"

	"github.com/spoonboy-io/koan"
)

// Make generates a self-signed X.509 certificate for a TLS server, code based on example code
// from the crypto/tls package found here https://go.dev/src/crypto/tls/generate_cert.go
func Make(logger *koan.Logger) error {
	// make private key
	var priv interface{}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("Failed to generate private key : %v", err)
	}

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature
	if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	validFrom := time.Now()
	validTo := validFrom.Add(internal.TLS_VALID_FOR)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{internal.TLS_ORG},
		},
		NotBefore: validFrom,
		NotAfter:  validTo,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// get the preferred outbound IP to register on the certificate
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return fmt.Errorf("Failed making UDP connection: %v", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	logger.Info(fmt.Sprintf("Using IP address '%v' for certifcate", localAddr.IP))

	template.IPAddresses = append(template.IPAddresses, localAddr.IP)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("Failed to create certificate: %v", err)

	}

	// file targets
	certDest := fmt.Sprintf("%s/cert.pem", internal.TLS_FOLDER)
	keyDest := fmt.Sprintf("%s/key.pem", internal.TLS_FOLDER)

	// write the certificate
	certOut, err := os.Create(certDest)
	if err != nil {
		return fmt.Errorf("Failed to open cert.pem for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("Failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("Error closing cert.pem: %v", err)
	}

	logger.Info("Created 'cert.pem' file")

	keyOut, err := os.OpenFile(keyDest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open key.pem for writing: %v", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("Unable to marshal private key: %v", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("Failed to write data to key.pem: %v", err)
	}

	if err := keyOut.Close(); err != nil {
		return fmt.Errorf("Error closing key.pem: %v", err)
	}

	logger.Info("Created 'key.pem' file")

	return nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}
