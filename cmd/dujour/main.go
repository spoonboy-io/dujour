package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spoonboy-io/dujour/internal"

	"github.com/spoonboy-io/dujour/internal/file"
	"github.com/spoonboy-io/koan"
)

var logger *koan.Logger

func init() {
	logger = &koan.Logger{}

	// check/create data folder
	dataPath := filepath.Join(".", internal.DATA_FOLDER)
	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		logger.FatalError("problem checking/creating data folder", err)
	}

	// check/create certificates folder
	tlsPath := filepath.Join(".", internal.TLS_FOLDER)
	if err := os.MkdirAll(tlsPath, os.ModePerm); err != nil {
		logger.FatalError("problem checking/creating 'certificates' folder", err)
	}

	// add self-signed certificate if folder empty
	// TODO
}

func main() {
	datasources, err := file.LoadAndValidateDatasources(internal.DATA_FOLDER, logger)
	if err != nil {
		logger.FatalError("problem loading data sources", err)
	}

	// add a watch to the folder for hot reload

	// create a server running as service

	// TODO remove debug
	for k, v := range datasources {
		fmt.Println("Datasource", k)
		fmt.Println("-------------------")
		fmt.Printf("\n%+v\n\n", v.Data)
		fmt.Println("------------------------------------------------")
		fmt.Println("")
	}
}

// TODO
// LoadAndValidate value return rather than operate on pointer
// generate the self-signed cert
// implement logging messages for the validate/load operations
// get tests in place on the work done so far
// create the server which offers HTTPS using the cert
// implement routing to server an API for the files
// implement a watcher to reload/create new data when files are added to the data folder
// revisit the camelcaseing for the map keys (will need to be recursive)
