package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/foolin/gocsv"
	"github.com/spoonboy-io/dujour/internal/file"
	"github.com/spoonboy-io/koan"
)

const (
	DATA_FOLDER = "data"
	TYPE_CSV    = 1
	TYPE_JSON   = 2
	DATA_ARR    = 1
	DATA_OBJ    = 2
)

var logger *koan.Logger

type datasource struct {
	fileName     string
	fileType     int
	endpointName string
	dataType     int
	ObjData      map[string]interface{}
	ArrData      []map[string]interface{}
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
	// load and validate all data from the folder
	datasources, err := LoadAndValidateDatasources(DATA_FOLDER)
	if err != nil {
		logger.FatalError("problem loading data sources", err)
	}

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

		ds := InitDatasource(fv)

		if err := LoadAndValidate(&ds); err != nil {
			logger.Error("error", err)
			continue
		}

		// add the datasource
		datasources = append(datasources, ds)
	}

	//fmt.Println(data)
	return datasources, nil
}

func InitDatasource(file string) datasource {
	ext := strings.ToLower(filepath.Ext(file))
	fileType := TYPE_JSON
	if ext == ".csv" {
		fileType = TYPE_CSV
	}

	_, filename := filepath.Split(file)
	endpointName := strings.Replace(filename, filepath.Ext(file), "", 1)
	ds := datasource{
		fileName:     file,
		fileType:     fileType,
		endpointName: endpointName,
	}

	return ds
}

func LoadAndValidate(ds *datasource) error {
	var err error
	var data []byte

	switch ds.fileType {
	case TYPE_CSV:
		mp := []map[string]interface{}{}
		mp, err = gocsv.Read(ds.fileName, true)
		if err != nil {
			return err
		}

		// array data
		ds.dataType = DATA_ARR
		ds.ArrData = mp

	case TYPE_JSON:
		// get the data
		data, err = os.ReadFile(ds.fileName)
		if err != nil {
			return err
		}

		// we need to handle array and object
		arr := []map[string]interface{}{}
		obj := map[string]interface{}{}

		if err = json.Unmarshal(data, &arr); err != nil {
			// we must have an object
			if err = json.Unmarshal(data, &obj); err != nil {
				return err
			}
			ds.dataType = DATA_OBJ
			ds.ObjData = obj
		} else {
			ds.dataType = DATA_ARR
			ds.ArrData = arr
		}
	}

	return nil
}
