package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var tests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/fern", http.StatusNotFound},
	}

	var app application
	routes := app.routes()

	// create a test server
	s := httptest.NewTLSServer(routes)
	defer s.Close()

	pathToTemplates = "./../../templates/"
	for _, e := range tests {
		resp, err := s.Client().Get(s.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s: expected %d, but go %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}

}
