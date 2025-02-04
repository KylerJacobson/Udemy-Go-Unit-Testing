package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"webapp/pkg/data"
)

func Test_app_authenticate(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "valid credentials",
			requestBody:    `{"email":"admin@example.com", "password":"secret"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not json",
			requestBody:    `Not json`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "empty json",
			requestBody:    `{}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "empty email",
			requestBody:    `{"email":"", "password":"secret"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "empty password",
			requestBody:    `{"email":"admin@example.com", "password":""}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "bad user",
			requestBody:    `{"email":"randomguy@example.com", "password":"secret"}`,
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, e := range tests {
		var reader io.Reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.authenticate)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: expected status %d, got %d", e.name, e.expectedStatus, rr.Code)
		}

	}
}

func Test_app_refresh(t *testing.T) {
	var tests = []struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}{
		{"valid", "", http.StatusOK, true},
		{"valid but not yet ready to expire", "", http.StatusTooEarly, false},
		{"expired token", expiredToken, http.StatusBadRequest, false},
	}
	testUser := data.User{ID: 1, FirstName: "Admin", LastName: "User", Email: "admin@example.com"}

	oldRefreshTime := jwtRefreshTokenExpiry
	for _, e := range tests {
		var tkn string
		if e.token == "" {
			if e.resetRefreshTime {
				jwtRefreshTokenExpiry = 1 * time.Second
			}
			tokens, _ := app.generateTokenPair(&testUser)
			tkn = tokens.RefreshToken
		} else {
			tkn = e.token
		}

		postedData := url.Values{
			"refresh_token": {tkn},
		}

		req, _ := http.NewRequest(http.MethodPost, "/v1/refresh-token", strings.NewReader(postedData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.refresh)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected status %d, got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		jwtRefreshTokenExpiry = oldRefreshTime
	}
}

func Test_app_userHandlers(t *testing.T) {
	var tests = []struct {
		name               string
		method             string
		json               string
		paramID            string
		handler            http.HandlerFunc
		expectedStatusCode int
	}{
		{"all users", http.MethodGet, "", "", app.allUsers, http.StatusOK},
		{"delete user", http.MethodDelete, "", "1", app.deleteUser, http.StatusNoContent},
		{"delete user invalid param", http.MethodDelete, "", "Y", app.deleteUser, http.StatusBadRequest},
		{"get user valid", http.MethodGet, "", "1", app.getUser, http.StatusOK},
		{"get user invalid", http.MethodGet, "", "2", app.getUser, http.StatusBadRequest},
		{"get user invalid param", http.MethodGet, "", "Y", app.getUser, http.StatusBadRequest},
		{
			"update valid user", http.MethodPatch, `{"id":1, "first_name": "Administrator", "last_name": "User", "email": "admin@example.com"}`, "1", app.updateUser, http.StatusNoContent,
		},
		{
			"update invalid user", http.MethodPatch, `{"id":2, "first_name": "Administrator", "last_name": "User", "email": "admin@example.com"}`, "1", app.updateUser, http.StatusBadRequest,
		},
		{
			"update user invalid json", http.MethodPatch, `{"id":1, first_name: "Administrator", "last_name": "User", "email": "admin@example.com"}`, "1", app.updateUser, http.StatusBadRequest,
		},
		{
			"insert user", http.MethodPut, `{"first_name": "Jack", "last_name": "Smith", "email": "jack@example.com"}`, "", app.insertUser, http.StatusNoContent,
		},
		{
			"insert invalid user", http.MethodPut, `{ "foo": "bar","first_name: "Jack", "last_name": "Smith", "email": "jack@example.com"}`, "", app.insertUser, http.StatusBadRequest,
		},
		{
			"insert user invalid json", http.MethodPut, `{ first_name: "Jack", "last_name": "Smith", "email": "jack@example.com"}`, "", app.insertUser, http.StatusBadRequest,
		},
	}

	for _, e := range tests {
		var req *http.Request
		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/v1/users", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/v1/users", strings.NewReader(e.json))
		}
		if e.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected status %d, got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}
