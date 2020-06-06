package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/txross1993/superman-api/db"
	"github.com/txross1993/superman-api/geolocate"
	"github.com/txross1993/superman-api/models"
	"github.com/txross1993/superman-api/superman"
	"github.com/txross1993/superman-api/testdata"
)

func TestAPISmokeTest(t *testing.T) {
	dir, _ := filepath.Abs(".")
	geoDb := path.Join(dir, "../GeoLite2-City_20200602/GeoLite2-City.mmdb")
	geoSvc, err := geolocate.NewGeoService(geoDb)
	if err != nil {
		t.Fatal(err)
	}

	localDB := "test.db"
	db, err := db.InitDB(localDB)
	if err != nil {
		t.Fatal(err)
	}

	defer geoSvc.Close()
	defer db.Cleanup()
	defer db.Close()

	api := NewAPI(Config{
		Host:     "localhost",
		Port:     "9099",
		Superman: superman.NewService(geoSvc, db),
	})

	event := testdata.GenerateCurrentEvent()
	b, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	req := newRequest(t, "POST", "/v1/", bytes.NewReader(b))
	resp := makeRequest(api.router, req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var got models.Superman
	err = json.Unmarshal(bodyBytes, &got)
	if err != nil {
		t.Fatal(err)
	}

	currentGeo, _ := geoSvc.GetCoordinatesFromIP(event.IPAddress)
	want := &models.Superman{
		CurrentGeo: currentGeo,
	}

	assert.Equal(t, 201, resp.Code)
	assert.Equal(t, got, want)
}

func newRequest(t *testing.T, method string, path string, body io.Reader) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

// makeRequest will test the http request using the httptest ResponseRecorder
func makeRequest(router http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
