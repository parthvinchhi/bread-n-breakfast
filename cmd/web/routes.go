package main

import (
	"net/http"

	"github.com/Pdv2323/bread-n-breakfast/pkg/config"
	"github.com/Pdv2323/bread-n-breakfast/pkg/handlers"
	"github.com/bmizerany/pat"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routesPat(app *config.AppConfig) http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	return mux
}

func routesChi(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/single-bed", handlers.Repo.SingleBed)
	mux.Get("/double-bed", handlers.Repo.DoubleBed)
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Get("/contact", handlers.Repo.ContactUs)
	mux.Get("/make-reservation", handlers.Repo.Reservation)

	return mux
}
