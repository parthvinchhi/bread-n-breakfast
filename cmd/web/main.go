package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/parthvinchhi/bread-n-breakfast/internal/config"
	my_driver "github.com/parthvinchhi/bread-n-breakfast/internal/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	// msg := models.MailData{
	// 	To:      "john@do.ca",
	// 	From:    "me@here.com",
	// 	Subject: "some subject",
	// 	Content: "",
	// }

	// app.MailChan <- msg

	fmt.Println(fmt.Sprintf("Starting Application on port %s\n", portNumber))
	// fmt.Println(fmt.Sprintf("Starting Application on port #{portNumber}"))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routesChi(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*my_driver.DB, error) {
	//What am I going to store in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	inProduction := flag.Bool("production", true, "Application in Production")
	useCache := flag.Bool("cache", true, "Use Template Cache")
	dbName := flag.String("dbname", "", "Database name")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database pass")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, required)")

	flag.Parse()

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//Change this to "True" when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

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

	// Connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := my_driver.ConnectSQL(connectionString)
	// db, err := my_driver.ConnectSQL("host=localhost port=5432 dbname=bread-n-breakfast user=postgres password=123 sslmode=disable")
	// db, err := my_driver.ConnectSQL("host=localhost port=5432 dbname=bread-n-breakfast user=postgres password= sslmode=disable")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...", err)
	}
	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
