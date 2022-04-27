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

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/text")

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

func (a *App) ListDatasources(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
	}

	a.Logger.Info("Served GET /list request - 200 OK")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(res))
}

func (a *App) DatasourceGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	dsReq := strings.ToLower(vars["datasource"])

	// this is suboptimal, but key is used to store the filename and we use the key in
	// hotreload add/delete operations, so we have to iterate and match endpoint for now
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
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
	a.Mtx.Unlock()

	if !foundMarker {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logMsg := fmt.Sprintf("Served GET /%s request - 200 OK", dsReq)
	a.Logger.Info(logMsg)
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(res))
}

// DatasourceGetByID will process a request for a datasource and return the element that matches the ID
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
				for _, v1 := range v.Data.([]map[string]interface{}) {
					// key could be string
					if fid, ok := v1["id"].(string); ok {
						if fid == id {
							foundMarker = true
							res, err = json.MarshalIndent(v1, "", "  ")
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}
					}
					// key could be integer
					if fid, ok := v1["id"].(int); ok {
						if string(fid) == id {
							foundMarker = true
							res, err = json.MarshalIndent(v1, "", "  ")
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}
					}

				}
			case map[string]interface{}:
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
											w.WriteHeader(http.StatusInternalServerError)
											return
										}
									}
								}
								// key could be integer
								if fid, ok := v4["id"].(int); ok {
									if string(fid) == id {
										foundMarker = true
										res, err = json.MarshalIndent(v4, "", "  ")
										if err != nil {
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
				for _, v := range v.Data.([]map[string]string) {
					if fid, ok := v["id"]; ok {
						if fid == id {
							foundMarker = true
							res, err = json.MarshalIndent(v, "", "  ")
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								return
							}
						}
					}
				}
			default:
				// something has gone wrong
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	a.Mtx.Unlock()

	if !foundMarker {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logMsg := fmt.Sprintf("Served GET /%s/%s request - 200 OK", dsReq, id)
	a.Logger.Info(logMsg)
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(res))
}
