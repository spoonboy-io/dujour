package internal

const (
	// folder locations
	DATA_FOLDER = "data"
	TLS_FOLDER  = "certs"

	// storage
	TYPE_CSV  = 1
	TYPE_JSON = 2
	DATA_ARR  = 1
	DATA_OBJ  = 2

	// tls configuration
	TLS_HOST        = ""
	TLS_VALID_FROM  = ""
	TLS_VALID_FOR   = ""
	TLS_IS_CA       = false
	TLS_RSA_BITS    = ""
	TLS_ECDSA_CURVE = ""
	TLS_ED25519_KEY = false
	TLS_CERT_NAME   = ""
	TLS_KEY_NAME    = ""

	/*
		host       = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
		validFrom  = flag.String("start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
		validFor   = flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
		isCA       = flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
		rsaBits    = flag.Int("rsa-bits", 2048, "Size of RSA key to generate. Ignored if --ecdsa-curve is set")
		ecdsaCurve = flag.String("ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521")
		ed25519Key = flag.Bool("ed25519", false, "Generate an Ed25519 key")
	*/
)

// Datasource contains both the data and metadata of a discovered and validated datasource
type Datasource struct {
	FileName     string
	FileType     int
	EndpointName string
	DataType     int
	Data         interface{}
}
