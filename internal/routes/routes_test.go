package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/mux"

	"github.com/spoonboy-io/dujour/internal"
	"github.com/spoonboy-io/koan"
)

func createTestAppContext() *App {
	testMtx := &sync.Mutex{}
	testLogger := &koan.Logger{}
	testDatasources := map[string]internal.Datasource{
		"data/people.csv": internal.Datasource{
			FileName:     "data/people.csv",
			FileType:     internal.TYPE_CSV,
			EndpointName: "people",
			Data: []map[string]string{
				{"id": "1", "name": "Test", "age": "100"},
				{"id": "2", "name": "Test2", "age": "25"},
			},
		},
		"data/people2.json": internal.Datasource{
			FileName:     "data/people2.json",
			FileType:     internal.TYPE_JSON,
			EndpointName: "people2",
			Data: []map[string]interface{}{
				{"id": 1, "name": "Test", "age": 100},
				{"id": 2, "name": "Test2", "age": 25},
			},
		},
		"data/people3.json": internal.Datasource{
			FileName:     "data/people3.json",
			FileType:     internal.TYPE_JSON,
			EndpointName: "people3",
			Data: map[string]interface{}{
				"result": []map[string]interface{}{
					{"id": "abc", "name": "Test", "age": 100},
					{"id": "DEF", "name": "Test2", "age": 25},
				},
			},
		},
	}

	testApp := &App{
		Logger:      testLogger,
		Datasources: testDatasources,
		Mtx:         testMtx,
	}

	return testApp
}

func TestHomeHandler(t *testing.T) {
	app := createTestAppContext()

	expected := "Dujour - JSON/CSV Data Server\n"
	expected += "=============================\n\n"
	expected += "Usage\n=====\n"
	expected += "GET / \t\t\t- Text format help page\n"
	expected += "GET /list \t\t- JSON array of all loaded datasources\n"
	expected += "GET /{datasource} \t- JSON representing all elements/rows for requested {datasource} or 404\n"
	expected += "GET /{datasource}/{id} \t- JSON representing element/row matching {id} from requested {datasource} or 404\n"

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Home)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListDatasources(t *testing.T) {
	testCases := []struct {
		name          string
		requestMethod string
		requestURI    string
		wantStatus    int
		wantListDS    []listDS
	}{
		{
			"simple request for the /list endpoint should be 200 OK",
			"GET",
			"/list",
			http.StatusOK,
			[]listDS{
				{
					Endpoint: "people",
					Source:   "data/people.csv",
				},
				{
					Endpoint: "people2",
					Source:   "data/people2.json",
				},
				{
					Endpoint: "people3",
					Source:   "data/people3.json",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			app := createTestAppContext()

			req, err := http.NewRequest(tc.requestMethod, tc.requestURI, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.ListDatasources)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.wantStatus)
			}

			gotListDS := []listDS{}
			err = json.Unmarshal(rr.Body.Bytes(), &gotListDS)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(gotListDS, tc.wantListDS) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					gotListDS, tc.wantListDS)
			}

		})
	}
}

func TestDatasourceGetAll(t *testing.T) {
	testCases := []struct {
		name          string
		requestMethod string
		requestURI    string
		wantStatus    int
		wantBody      string
	}{
		{
			"request for /people endpoint should be 200 OK",
			"GET",
			"/people",
			http.StatusOK,
			"[{\"age\":\"100\",\"id\":\"1\",\"name\":\"Test\"},{\"age\":\"25\",\"id\":\"2\",\"name\":\"Test2\"}]",
		},
		{
			"request for /servers endpoint should be 404 Not Found",
			"GET",
			"/servers",
			http.StatusNotFound,
			"404pagenotfound",
		},
		{
			"request for /People endpoint should be 404 Not Found (we support lowercase routes by design)",
			"GET",
			"/People",
			http.StatusNotFound,
			"404pagenotfound",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			app := createTestAppContext()

			req, err := http.NewRequest(tc.requestMethod, tc.requestURI, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			testMux := mux.NewRouter()
			testMux.HandleFunc("/{datasource:[a-z0-9=\\-\\/]+}", app.DatasourceGetAll).Methods("GET")
			testMux.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.wantStatus)
			}

			// we do this to make comparison of the output simpler
			gotBody := strings.ReplaceAll(rr.Body.String(), "\n", "")
			gotBody = strings.ReplaceAll(gotBody, " ", "")
			if gotBody != tc.wantBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					gotBody, tc.wantBody)
			}

		})
	}
}

func TestDatasourceGetByID(t *testing.T) {
	testCases := []struct {
		name          string
		requestMethod string
		requestURI    string
		wantStatus    int
		wantBody      string
	}{
		{
			"request for /people/1 endpoint should be 200 OK",
			"GET",
			"/people/1",
			http.StatusOK,
			"{\"age\":\"100\",\"id\":\"1\",\"name\":\"Test\"}",
		},
		{
			"request for /people/10 endpoint should be 404 Not Found",
			"GET",
			"/people/10",
			http.StatusNotFound,
			"404pagenotfound",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			app := createTestAppContext()

			req, err := http.NewRequest(tc.requestMethod, tc.requestURI, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			testMux := mux.NewRouter()
			testMux.HandleFunc("/{datasource:[a-z0-9=\\-\\/]+}/{id:[a-zA-Z0-9=\\-\\/]+}", app.DatasourceGetByID).Methods("GET")
			testMux.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.wantStatus)
			}

			// we do this to make comparison of the output simpler
			gotBody := strings.ReplaceAll(rr.Body.String(), "\n", "")
			gotBody = strings.ReplaceAll(gotBody, " ", "")
			if gotBody != tc.wantBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					gotBody, tc.wantBody)
			}

		})
	}
}
