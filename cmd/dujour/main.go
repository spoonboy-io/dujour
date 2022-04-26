package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/spoonboy-io/dujour/internal"
	"github.com/spoonboy-io/dujour/internal/certificate"
	"github.com/spoonboy-io/dujour/internal/file"
	"github.com/spoonboy-io/koan"
	"github.com/spoonboy-io/reprise"
)

var (
	version   = "Development build"
	goversion = "Unknown"
)

var logger *koan.Logger

func init() {
	logger = &koan.Logger{}

	// check/create data folder
	dataPath := filepath.Join(".", internal.DATA_FOLDER)
	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		logger.FatalError("Problem checking/creating data folder", err)
	}

	// check/create certificates folder
	tlsPath := filepath.Join(".", internal.TLS_FOLDER)
	if err := os.MkdirAll(tlsPath, os.ModePerm); err != nil {
		logger.FatalError("Problem checking/creating 'certificates' folder", err)
	}

	// add self-signed certificate only if folder empty, if the cert expires it
	// it can be deleted so the code here create a new cert.pem and key.pem files
	checkExist := fmt.Sprintf("%s/cert.pem", internal.TLS_FOLDER)
	if _, err := os.Stat(checkExist); errors.Is(err, os.ErrNotExist) {
		logger.Info("Creating self-signed TLS certificate for the server")
		if err := certificate.Make(logger); err != nil {
			logger.FatalError("Problem creating the certificate/key", err)
		}
	}
}

func main() {
	// write a console banner
	reprise.WriteSimple(&reprise.Banner{
		Name:         "Dujour",
		Description:  "JSON/CSV Data Server",
		Version:      version,
		GoVersion:    goversion,
		WebsiteURL:   "https://spoonboy.io",
		VcsURL:       "https://github.com/spoonboy-io/dujour",
		VcsName:      "Github",
		EmailAddress: "hello@spoonboy.io",
	})

	datasources, err := file.LoadAndValidateDatasources(internal.DATA_FOLDER, logger)
	if err != nil {
		logger.FatalError("Problem loading data sources", err)
	}

	if len(datasources) == 0 {
		logger.Warn(fmt.Sprintf("Currently there are datasources to serve, add JSON or CSV files to the '%s' folder", internal.DATA_FOLDER))
	}

	// add a watch to the folder for hot reload

	// router for testing
	mux := mux.NewRouter()
	mux.HandleFunc("/", Home).Methods("GET")

	// create a server running as service
	hostPort := fmt.Sprintf("%s:%s", internal.SRV_HOST, internal.SRV_PORT)
	srvTLS := &http.Server{
		Addr:         hostPort,
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// start HTTPS server
	logger.Info(fmt.Sprintf("Starting HTTPS server on %s", hostPort))
	if err := srvTLS.ListenAndServeTLS(fmt.Sprintf("%s/cert.pem", internal.TLS_FOLDER), fmt.Sprintf("%s/key.pem", internal.TLS_FOLDER)); err != nil {
		logger.FatalError("Failed to start HTTPS server", err)
	}

	// TODO remove debug
	_ = datasources
	/*
		for k, v := range datasources {
			fmt.Println("Datasource", k)
			fmt.Println("-------------------")
			fmt.Printf("\n%+v\n\n", v.Data)
			fmt.Println("------------------------------------------------")
			fmt.Println("")
		}
	*/

}

// Test the TLS certificate
func Home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Serving over TLS"))
}

// TODO
// implement logging messages for the validate/load operations
// get tests in place on the work done so far
// implement routing to server an API for the files
// implement a watcher to reload/create new data when files are added to the data folder
// revisit the camelcaseing for the map keys (will need to be recursive)
