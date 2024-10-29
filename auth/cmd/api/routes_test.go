package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"testing"
)

func Test_routes_exists(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router) // actual routes

	routes := []string{"/authenticate"} // expected routes

	for _, route := range routes {
		routExists(t, chiRoutes, route)
	}
}

func routExists(t *testing.T, routes chi.Router, route string) {
	found := false
	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("did not find %v in registered routes", route)
	}
}
