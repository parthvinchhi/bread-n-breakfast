package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Pdv2323/bread-n-breakfast/internal/config"
	"github.com/Pdv2323/bread-n-breakfast/internal/handlers"
	"github.com/Pdv2323/bread-n-breakfast/internal/models"
	"github.com/Pdv2323/bread-n-breakfast/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8088"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(fmt.Sprintf("Starting Application on port %s\n", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routesChi(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	//What am I going to store in session
	gob.Register(models.Reservation{})

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
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	return nil
}
