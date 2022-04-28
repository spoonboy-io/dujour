// Package routes contains all the handlers used in the application
package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/spoonboy-io/dujour/internal"
	"github.com/spoonboy-io/koan"
)

type App struct {
	Logger      *koan.Logger
	Datasources map[string]internal.Datasource
	Mtx         *sync.Mutex
}

// Home provides basic instruction on how to poll the datasources hosted by the application as text format.
func (a *App) Home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	res := "Dujour - JSON/CSV Data Server\n"
	res += "=============================\n\n"
	res += "Usage\n=====\n"
	res += "GET / \t\t\t- Text format help page\n"
	res += "GET /list \t\t- JSON array of all loaded datasources\n"
	res += "GET /{datasource} \t- JSON representing all elements/rows for requested {datasource} or 404\n"
	res += "GET /{datasource}/{id} \t- JSON representing element/row matching {id} from requested {datasource} or 404\n"

	a.Logger.Info("Served GET / reques - 200 OK")
	_, _ = fmt.Fprint(w, res)
}

// ListDatasources provides a summary of datasources hosted by the application in JSON format
func (a *App) ListDatasources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// this is the information we will output for list
	type listDS struct {
		Endpoint string `json:"endpoint"`
		Source   string `json:"source"`
	}

	list := []listDS{}

	// iterate the datasources
	a.Mtx.Lock()
	for _, v := range a.Datasources {
		ds := listDS{
			v.EndpointName,
			v.FileName,
		}

		list = append(list, ds)
	}
	a.Mtx.Unlock()

	res, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		a.Logger.Error("Marshalling ListDatasources:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.Logger.Info("Served GET /list request - 200 OK")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(res))
}

// DatasourceGetAll will retrieve all data for a datasource in JSON format
func (a *App) DatasourceGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	dsReq := strings.ToLower(vars["datasource"])
	foundMarker := false
	var res []byte
	var err error

	// we are reading from the map, should not need mutex but adding
	a.Mtx.Lock()
	for _, v := range a.Datasources {
		if v.EndpointName == dsReq {
			foundMarker = true
			res, err = json.MarshalIndent(v.Data, "", "  ")
			if err != nil {
				a.Logger.Error("Marshaling DatasourceGetAll:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	a.Mtx.Unlock()

	if !foundMarker {
		logMsg := fmt.Sprintf("Served GET /%s request - 404 Not Found", dsReq)
		a.Logger.Info(logMsg)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logMsg := fmt.Sprintf("Served GET /%s request - 200 OK", dsReq)
	a.Logger.Info(logMsg)
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(res))
}

// DatasourceGetByID will process a request for a datasource and return the element that matches the ID in JSON format
func (a *App) DatasourceGetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	dsReq := vars["datasource"]
	id := vars["id"]

	foundMarker := false
	var res []byte
	var err error

	a.Mtx.Lock()
	for _, v := range a.Datasources {
		if v.EndpointName == dsReq {
			// we have the datasource but what about the id?
			// need a type assertion to discover what we have need to parse
			switch v.Data.(type) {
			case []map[string]interface{}:
				// 1
				for _, v1 := range v.Data.([]map[string]interface{}) {
					// key could be string
					if fid, ok := v1["id"].(string); ok {
						if fid == id {
							foundMarker = true
							res, err = json.MarshalIndent(v1, "", "  ")
							if err != nil {
								a.Logger.Error("Marshaling DatasourceGetByID (case 1, string):", err)
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}
					}
					// key could be integer
					if fid, ok := v1["id"].(int); ok {
						if fmt.Sprint(fid) == id {
							foundMarker = true
							res, err = json.MarshalIndent(v1, "", "  ")
							if err != nil {
								a.Logger.Error("Marshaling DatasourceGetByID (case 1, int):", err)
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}
					}

				}
			case map[string]interface{}:
				// 2
				for _, v1 := range v.Data.(map[string]interface{}) {
					// the value stored should a slice otherwise we don't have a list of data, only an object
					if v2, ok := v1.([]interface{}); ok {
						for _, v3 := range v2 {
							if v4, ok := v3.(map[string]interface{}); ok {
								// we can now inspect v4 for data
								// key could be string
								if fid, ok := v4["id"].(string); ok {
									if fid == id {
										foundMarker = true
										res, err = json.MarshalIndent(v4, "", "  ")
										if err != nil {
											a.Logger.Error("Marshaling DatasourceGetByID (case 2, sttring):", err)
											w.WriteHeader(http.StatusInternalServerError)
											return
										}
									}
								}
								// key could be integer
								if fid, ok := v4["id"].(int); ok {
									if fmt.Sprint(fid) == id {
										foundMarker = true
										res, err = json.MarshalIndent(v4, "", "  ")
										if err != nil {
											a.Logger.Error("Marshaling DatasourceGetByID (case 2, int):", err)
											w.WriteHeader(http.StatusInternalServerError)
											return
										}
									}
								}
							}
						}
					}
				}
			case []map[string]string:
				// 3
				for _, v := range v.Data.([]map[string]string) {
					if fid, ok := v["id"]; ok {
						if fid == id {
							foundMarker = true
							res, err = json.MarshalIndent(v, "", "  ")
							if err != nil {
								a.Logger.Error("Marshaling DatasourceGetByID (case 3):", err)
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}
					}
				}
			default:
				// something has gone wrong
				a.Logger.Warn("DatasourceGetByID unexpected Type. Unhandled")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	a.Mtx.Unlock()

	if !foundMarker {
		logMsg := fmt.Sprintf("Served GET /%s/%s request - 404 Not Found", dsReq, id)
		a.Logger.Info(logMsg)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logMsg := fmt.Sprintf("Served GET /%s/%s request - 200 OK", dsReq, id)
	a.Logger.Info(logMsg)
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(res))
}
