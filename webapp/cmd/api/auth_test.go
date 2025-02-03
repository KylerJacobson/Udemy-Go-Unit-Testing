package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func Test_app_getTokenFromHeaderAndVerify(t *testing.T) {
	testUser := data.User{ID: 1, FirstName: "admin", LastName: "User", Email: "admin@example.com"}

	tokens, _ := app.generateTokenPair(&testUser)

	test := []struct {
		name          string
		token         string
		errorExpected bool
		setHeader     bool
		issuer        string
	}{
		{"valid token", fmt.Sprintf("Bearer %s", tokens.Token), false, true, app.Domain},
		{"valid but expire", fmt.Sprintf("Bearer %s", expiredToken), true, true, app.Domain},
		{"no header", "", true, false, app.Domain},
		{"no token", fmt.Sprintf("Bearer %s", ""), true, true, app.Domain},
		{"bad token", fmt.Sprintf("Bearer %s", tokens.Token+"dsfadsf"), true, true, app.Domain},
		{"bad header", fmt.Sprintf("Bear %s", tokens.Token), true, true, app.Domain},
		// Make sure the next test is the last one to run
		{"wrong issuer", fmt.Sprintf("Bearer %s", tokens.Token), true, true, "baddomain.com"},
	}
	for _, e := range test {
		if e.issuer != app.Domain {
			app.Domain = e.issuer
			tokens, _ = app.generateTokenPair(&testUser)
		}
		req, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
		if e.setHeader {
			req.Header.Set("Authorization", e.token)
		}
		rr := httptest.NewRecorder()
		_, _, err := app.getTokenFromHeaderAndVerify(rr, req)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: expected no error, got %v", e.name, err)
		}

		if err == nil && e.errorExpected {
			t.Errorf("%s: expected error, got none", e.name)
		}
		app.Domain = "example.com"
	}
}
