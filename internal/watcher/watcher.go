// Package watcher monitors a data folder and performs automatic reloading of data
// including updating the datasource cache in memory for deleted and edited data files
package watcher

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spoonboy-io/dujour/internal/file"

	"github.com/spoonboy-io/dujour/internal"

	"github.com/fsnotify/fsnotify"
	"github.com/spoonboy-io/koan"
)

// Monitor creates a file watcher for the data directory
func Monitor(datasources map[string]internal.Datasource, logger *koan.Logger, mtx *sync.Mutex) error {
	watchPath := filepath.Join(".", internal.DATA_FOLDER)

	logger.Info(fmt.Sprintf("Creating Watcher for '%s' folder", watchPath))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Could not create watcher; %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				switch event.Op {
				case fsnotify.Create:
					logger.Info(fmt.Sprintf("File added '%s'", event.Name))

					// is it a file we want to process
					extension := strings.ToLower(filepath.Ext(event.Name))
					if (extension != ".csv") && (extension != ".json") {
						if extension != "" {
							logger.Warn(fmt.Sprintf("Hotloader skipping file '%s', file extension is '%s'", event.Name, extension))
							continue
						}
					}

					// init & validate the file
					var hlds internal.Datasource
					var err error
					if hlds, err = file.LoadAndValidate(file.InitDatasource(event.Name), logger); err != nil {
						logger.Error(fmt.Sprintf("Could not hotload datasource '%s'", event.Name), err)
					}

					// add the datasource
					mtx.Lock()
					datasources[event.Name] = hlds
					mtx.Unlock()
				case fsnotify.Rename, fsnotify.Remove:
					logger.Info(fmt.Sprintf("Hotloader file removed '%s'", event.Name))
					// remove
					mtx.Lock()
					delete(datasources, event.Name)
					mtx.Unlock()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("Watcher unexpected error:", err)
			}
		}
	}()

	if err := watcher.Add(watchPath); err != nil {
		return fmt.Errorf("Adding file failed; %v", err)
	}

	<-done

	return nil
}
