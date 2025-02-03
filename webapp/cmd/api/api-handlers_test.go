package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
