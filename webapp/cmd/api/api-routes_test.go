package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"testing"
)

func Test_application_routes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{"/v1/auth", "POST"},
		{"/v1/refresh-token", "POST"},
		{"/v1/users/", "GET"},
		{"/v1/users/{id}", "GET"},
		{"/v1/users/{id}", "DELETE"},
		{"/v1/users/{id}", "PUT"},
		{"/v1/users/", "PATCH"},
	}

	mux := app.routes()

	chiRoutes := mux.(chi.Routes)

	for _, route := range registered {
		// check to see if the route exists
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})

	return found
}
