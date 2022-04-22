package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spoonboy-io/dujour/internal/file"
	"github.com/spoonboy-io/koan"
)

const (
	DATA_FOLDER = "data"
	TYPE_CSV    = 1
	TYPE_JSON   = 2
)

var logger *koan.Logger

type datasource struct {
	fileName     string
	fileType     int
	endpointName string
	data         map[string]interface{}
}

func init() {
	logger = &koan.Logger{}

	// check/create data folder
	dataPath := filepath.Join(".", DATA_FOLDER)
	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		logger.FatalError("problem checking/creating data folder", err)
	}
}

func main() {
	// load all data from the folder
	datasources, err := LoadAndValidateDatasources(DATA_FOLDER)
	if err != nil {
		logger.FatalError("problem loading data sources", err)
	}

	// perform validation

	// add a watch to the folder for hot reload

	// create a server running as service

	fmt.Println("hello", datasources)
}

// LoadAndValidateDatasources finds, loads and validates all data
func LoadAndValidateDatasources(dataFolder string) ([]datasource, error) {
	datasources := []datasource{}

	files, err := file.FindFiles(dataFolder, logger)
	if err != nil {
		return nil, err
	}

	for _, fv := range files {

		// load the file

		// validate the content
		//if ok := validate.CheckFileContent(fv); !ok {
		//	continue
		//}

		ext := strings.ToLower(filepath.Ext(fv))
		fileType := TYPE_JSON
		if ext == ".csv" {
			fileType = TYPE_CSV
		}

		_, filename := filepath.Split(fv)
		endpointName := strings.Replace(filename, filepath.Ext(fv), "", 1)
		ds := datasource{
			fileName:     fv,
			fileType:     fileType,
			endpointName: endpointName,
			data:         map[string]interface{}{},
		}

		datasources = append(datasources, ds)
	}

	//fmt.Println(data)
	return datasources, nil
}
