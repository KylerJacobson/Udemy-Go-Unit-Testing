package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func Test_app_enableCORS(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	var tests = []struct {
		name           string
		method         string
		expectedHeader bool
	}{
		{"Preflight", http.MethodOptions, true},
		{"GET", http.MethodGet, false},
	}

	for _, e := range tests {
		handlerToTest := app.enableCORS(nextHandler)
		req := httptest.NewRequest(e.method, "/v1/test", nil)
		rr := httptest.NewRecorder()
		handlerToTest.ServeHTTP(rr, req)
		if e.expectedHeader && rr.Header().Get("Access-Control-Allow-Credentials") == "" {
			t.Errorf("%s: expected Access-Control-Allow-Credentials header to be set", e.name)
		}

		if !e.expectedHeader && rr.Header().Get("Access-Control-Allow-Credentials") != "" {
			t.Errorf("%s: expected Access-Control-Allow-Credentials header to not be set", e.name)
		}
	}
}

func Test_app_authRequired(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	testUser := data.User{ID: 1, FirstName: "Admin", LastName: "User", Email: "admin@example.com"}

	tokens, _ := app.generateTokenPair(&testUser)

	var tests = []struct {
		name             string
		token            string
		expectAuthorized bool
		setHeader        bool
	}{
		{"valid token", fmt.Sprintf("Bearer %s", tokens.Token), true, true},
		{"invalid token", fmt.Sprintf("Bearer %s", tokens.Token+"dsfadsf"), false, true},
		{"no token", "", false, false},
		{"expired token", fmt.Sprintf("Bearer %s", expiredToken), false, true},
	}

	for _, e := range tests {
		handlerToTest := app.authRequired(nextHandler)
		req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
		if e.setHeader {
			req.Header.Set("Authorization", e.token)
		}
		rr := httptest.NewRecorder()
		handlerToTest.ServeHTTP(rr, req)
		if rr.Code == http.StatusUnauthorized && e.expectAuthorized {
			t.Errorf("%s: expected authorized, got unauthorized", e.name)
		}

		if rr.Code != http.StatusUnauthorized && !e.expectAuthorized {
			t.Errorf("%s: expected unauthorized, got authorized", e.name)
		}
	}
}
