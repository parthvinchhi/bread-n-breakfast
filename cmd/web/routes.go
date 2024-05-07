package main

import (
	"net/http"

	"github.com/Pdv2323/bread-n-breakfast/pkg/config"
	"github.com/Pdv2323/bread-n-breakfast/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routesChi(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// To use static folder
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/single-bed", handlers.Repo.SingleBed)
	mux.Get("/double-bed", handlers.Repo.DoubleBed)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Get("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/contact", handlers.Repo.ContactUs)

	mux.Get("/make-reservation", handlers.Repo.Reservation)

	return mux
}
