package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/parthvinchhi/bread-n-breakfast/internal/config"
	"github.com/parthvinchhi/bread-n-breakfast/internal/handlers"
	"github.com/parthvinchhi/bread-n-breakfast/internal/helpers"
	"github.com/parthvinchhi/bread-n-breakfast/internal/models"
	"github.com/parthvinchhi/bread-n-breakfast/internal/render"
)

const portNumber = ":8088"

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger
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

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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
	helpers.NewHelpers(&app)

	return nil
}
