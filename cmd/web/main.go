package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RakhmanovTimur/bookings/internal/config"
	"github.com/RakhmanovTimur/bookings/internal/driver"
	"github.com/RakhmanovTimur/bookings/internal/handlers"
	"github.com/RakhmanovTimur/bookings/internal/helpers"
	"github.com/RakhmanovTimur/bookings/internal/models"
	"github.com/RakhmanovTimur/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	// Close mail channel when application stops
	defer close(app.MailChan)
	listenForMail()
	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	cache := flag.Bool("cache", true, "Use cache templates")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpassword", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()
	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}
	// Set up the mail channel
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *cache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Connect to database
	log.Println("Connecting to database...")
	connectionSettings := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		*dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionSettings)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	return db, err
}
