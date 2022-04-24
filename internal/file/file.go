package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"

	"github.com/spoonboy-io/dujour/internal"
	"github.com/spoonboy-io/koan"
)

// FindFiles identifies all JSON and CSV files in the target dataFolder, files which
// are not JSON or CSV (as determined by the extension) will be skipped but logged
func FindFiles(dataFolder string, logger *koan.Logger) ([]string, error) {
	var files []string
	dataPath := filepath.Join(".", dataFolder)
	filepath.WalkDir(dataPath, func(s string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		extension := strings.ToLower(filepath.Ext(f.Name()))

		if (extension == ".csv") || (extension == ".json") {
			files = append(files, s)
		} else {
			if extension != "" {
				logger.Warn(fmt.Sprintf("Skipping file: '%s', (file extension '%s')", f.Name(), extension))
			}
		}
		return nil
	})
	return files, nil
}

// InitDatasource create a Datasource type and partially configures it with known info about the datasource
func InitDatasource(file string) internal.Datasource {
	ext := strings.ToLower(filepath.Ext(file))
	fileType := internal.TYPE_JSON
	if ext == ".csv" {
		fileType = internal.TYPE_CSV
	}

	_, filename := filepath.Split(file)
	endpointName := strings.Replace(filename, filepath.Ext(file), "", 1)
	ds := internal.Datasource{
		FileName:     file,
		FileType:     fileType,
		EndpointName: endpointName,
	}

	return ds
}

// LoadAndValidateDatasources finds, loads and validates all data at application startup
func LoadAndValidateDatasources(dataFolder string, logger *koan.Logger) (map[string]internal.Datasource, error) {
	datasources := map[string]internal.Datasource{}

	files, err := FindFiles(dataFolder, logger)
	if err != nil {
		return nil, err
	}

	for _, fv := range files {
		ds := InitDatasource(fv)
		if err := LoadAndValidate(&ds, logger); err != nil {
			logger.Error("error", err) // TODO tempory for debug
			continue
		}
		datasources[fv] = ds
	}

	return datasources, nil
}

// LoadAndValidate performs the load and validation at the individual datasource level for both JSON and CSV
// file formats, it also logs non fatal warnings and errors which may prevent proper parsing of a datasource
func LoadAndValidate(ds *internal.Datasource, logger *koan.Logger) error {

	// TODO - log or return. We also need more log information

	var err error
	var data []byte

	switch ds.FileType {
	case internal.TYPE_CSV:
		// load the CSV data
		data, err = os.ReadFile(ds.FileName)
		if err != nil {
			return err
		}

		rdr := bytes.NewReader(data)
		mp, err := gocsv.CSVToMaps(rdr)
		if err != nil {
			return err
		}

		// if good always array data when we have a CSV
		ds.DataType = internal.DATA_ARR
		ds.Data = mp

	case internal.TYPE_JSON:
		// load the JSON data
		data, err = os.ReadFile(ds.FileName)
		if err != nil {
			return err
		}

		// we potentially need to handle array and object when dealing with unknown JSON
		arr := []map[string]interface{}{}
		obj := map[string]interface{}{}

		if err = json.Unmarshal(data, &arr); err != nil {
			// we 'may' have an object or it could just be bad data
			if err = json.Unmarshal(data, &obj); err != nil {
				return err
			}
			// it was an object
			ds.DataType = internal.DATA_OBJ
			ds.Data = obj
		} else {
			// it was an array
			ds.DataType = internal.DATA_ARR
			ds.Data = arr
		}
	}

	return nil
}
