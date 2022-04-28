package file_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spoonboy-io/dujour/internal"

	"github.com/spoonboy-io/dujour/internal/file"
	"github.com/spoonboy-io/koan"
)

func TestFindFiles(t *testing.T) {

	logger := &koan.Logger{}
	testCases := []struct {
		name       string
		dataFolder string
		testFiles  []string
		wantFiles  []string
	}{
		{
			"files are all good but mixed extension case",
			"data",
			[]string{"file1.csv", "file2.CSV", "file3.json", "file4.JSON"},
			[]string{"data/file1.csv", "data/file2.CSV", "data/file3.json", "data/file4.JSON"},
		},

		{
			"files contain files we should be ignoring",
			"data",
			[]string{"file1.csv", "text.txt", "file3.json", "excel,xls"},
			[]string{"data/file1.csv", "data/file3.json"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create the folder
			if err := makeTestFolder(tc.dataFolder); err != nil {
				t.Fatalf("TestFindfiles could not create the test folder: %v", err)
			}

			// add the files
			if err := createTestFiles(tc.testFiles, tc.dataFolder); err != nil {
				t.Fatalf("TestFindfiles could not create the test files: %v", err)
			}

			gotFiles, err := file.FindFiles(tc.dataFolder, logger)
			if err != nil {
				t.Fatalf("Findfiles unexpected error: %v", err)
			}

			// compare slices
			if !reflect.DeepEqual(gotFiles, tc.wantFiles) {
				t.Errorf("failed got %v wanted %v", gotFiles, tc.wantFiles)
			}

			// tear down
			if err := removeTestFolder(tc.dataFolder); err != nil {
				t.Fatalf("TestFindfiles remove test folder %v", err)
			}
		})
	}

}

func TestInitDatasource(t *testing.T) {
	testCases := []struct {
		name           string
		filename       string
		wantDatasource internal.Datasource
	}{
		{
			"good file mixed case",
			"data/MyDataFile.CSV",
			internal.Datasource{
				FileName:     "data/MyDataFile.CSV",
				FileType:     internal.TYPE_CSV,
				EndpointName: "mydatafile",
			},
		},

		{
			"good file mixed case some underscore chars to replace",
			"data/My_Data_File.json",
			internal.Datasource{
				FileName:     "data/My_Data_File.json",
				FileType:     internal.TYPE_JSON,
				EndpointName: "my-data-file",
			},
		},

		{
			"good file mixed case some space chars to replace",
			"data/My Data File.json",
			internal.Datasource{
				FileName:     "data/My Data File.json",
				FileType:     internal.TYPE_JSON,
				EndpointName: "my-data-file",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotDatasource := file.InitDatasource(tc.filename)

			if tc.wantDatasource.FileName != gotDatasource.FileName {
				t.Errorf("failed on filename got %v wanted %v", gotDatasource.FileName, tc.wantDatasource.FileName)
			}

			if tc.wantDatasource.FileType != gotDatasource.FileType {
				t.Errorf("failed on filetype got %v wanted %v", gotDatasource.FileType, tc.wantDatasource.FileType)
			}

			if tc.wantDatasource.EndpointName != gotDatasource.EndpointName {
				t.Errorf("failed on endpointname got %v wanted %v", gotDatasource.EndpointName, tc.wantDatasource.EndpointName)
			}
		})
	}
}

func makeTestFolder(folder string) error {
	dataPath := filepath.Join(".", folder)
	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func createTestFiles(files []string, folder string) error {
	for _, v := range files {
		dataPath := filepath.Join(".", folder, "/", v)
		if err := os.WriteFile(dataPath, []byte("sample data"), 0644); err != nil {
			return err
		}
	}
	return nil
}

func removeTestFolder(folder string) error {
	dataPath := filepath.Join(".", folder)
	if err := os.RemoveAll(dataPath); err != nil {
		return err
	}
	return nil
}

/*

package main

import (
    "fmt"
    "reflect"
)

func main() {
    a := []int {4,5,6}
    b := []int {4,5,6}
    c := []int {4,5,6,7}

    fmt.Println(reflect.DeepEqual(a, b))
    fmt.Println(reflect.DeepEqual(a, c))

}
*/
