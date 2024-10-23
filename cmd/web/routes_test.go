package main

import (
	"fmt"
	"testing"

	"github.com/parthvinchhi/bread-n-breakfast/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routesChi(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing
	default:
		t.Error(fmt.Sprintf("This is not *chi.Mux, type is %T", v))
	}
}
