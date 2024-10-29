package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"testing"
)

func Test_routes_exists(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	routes := []string{"/authenticate"}
}

func routExists(t *testing.T, routes chi.Router, router string) {
	found := false
	_ = chi.Walk(routes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {

	})
}
