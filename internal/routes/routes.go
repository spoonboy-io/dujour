package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spoonboy-io/dujour/internal"
	"github.com/spoonboy-io/koan"
)

type App struct {
	Logger      *koan.Logger
	Datasources map[string]internal.Datasource
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

	_, _ = fmt.Fprint(w, res)
}

func (a *App) ListDatasources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// this is the information we will output for list
	type listDS struct {
		Endpoint string `json:"endpoint"`
		Source   string `json:"source"`
		Elements int    `json:"elements"`
	}

	list := []listDS{}

	// iterate the datasources
	for _, v := range a.Datasources {
		ds := listDS{
			v.EndpointName,
			v.FileName,
			0, // need to add this
		}

		list = append(list, ds)
	}

	res, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(res))
}

func (a *App) DatasourceGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("GetAll"))
}

func (a *App) DatasourceGetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("Get By ID"))
}

func setHeaders() {}

/*
func DeleteGeofence(w http.ResponseWriter, r *http.Request) {
	// Secure this route, must be logged in
	if session.Values["loggedIn"] != true {
		http.Redirect(w, r, "https://"+DOMAIN+PORTTLS, http.StatusFound)
	}

	vars := mux.Vars(r)
	fenceID := vars["id"]
	userKey := session.Values["key"].(string)
	err := db.DeleteGeofence(fenceID, userKey)

	if err != nil {
		session.AddFlash("There was a problem, the geofence could not be deleted", "error")
	} else {
		session.AddFlash("Geofence was successfully deleted", "confirm")
	}

	// Save session
	session.Save(r, w)

	// Perform redirect
	http.Redirect(w, r, "https://"+DOMAIN+PORTTLS+"/account", http.StatusFound)
}
*/
