package routes

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

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
					{"id": 1, "name": "Test", "age": 100},
					{"id": 2, "name": "Test2", "age": 25},
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
