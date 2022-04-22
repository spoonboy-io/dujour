package file

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

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
