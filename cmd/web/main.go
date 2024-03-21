package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Pdv2323/bread-n-breakfast/pkg/config"
	"github.com/Pdv2323/bread-n-breakfast/pkg/handlers"
	"github.com/Pdv2323/bread-n-breakfast/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8000"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	//Change this to "True" when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf(fmt.Sprintf("Starting Application on port %s\n", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routesChi(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
