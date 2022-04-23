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
}

func main() {
	datasources, err := file.LoadAndValidateDatasources(internal.DATA_FOLDER, logger)
	if err != nil {
		logger.FatalError("problem loading data sources", err)
	}

	// add a watch to the folder for hot reload

	// create a server running as service

	// TODO remove debug
	for _, v := range datasources {
		fmt.Println("Datasource")
		fmt.Println("-------------------")
		fmt.Printf("\n%+v\n\n", v)
		fmt.Println("------------------------------------------------")
		fmt.Println("")
	}
}
