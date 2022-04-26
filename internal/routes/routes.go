package routes

import (
	"net/http"

	"github.com/spoonboy-io/dujour/internal"
	"github.com/spoonboy-io/koan"
)

type App struct {
	Logger      *koan.Logger
	Datasources map[string]internal.Datasource
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Home"))
}

func (a *App) ListDatasources(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("List"))
}

func (a *App) DatasourceGetAll(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("GetAll"))
}

func (a *App) DatasourceGetByID(w http.ResponseWriter, r *http.Request) {
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
