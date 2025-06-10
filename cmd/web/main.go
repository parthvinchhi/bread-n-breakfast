package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const portNumber = ":8088"

// var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	// Close GORM's underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get raw DB:", err)
	}
	defer sqlDB.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Println(fmt.Sprintf("Starting Application on port %s\n", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routesChi(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*gorm.DB, error) {
	//What am I going to store in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

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

	app.Session = session

	log.Println("Connecting to GORM database...")
	db, err := connect()
	if err != nil {
		log.Fatal("Cannot connect to GORM DB:", err)
	}
	log.Println("Connected to GORM database")

	db.AutoMigrate(&models.Reservation{}, &models.User{}, &models.Room{}, &models.Restriction{})

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}

func connect() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=123 dbname=bread-n-breakfast port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
