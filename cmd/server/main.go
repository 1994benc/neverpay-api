package main

import (
	"1994benc/neverpay-user-service/internal/database"
	transportHTTP "1994benc/neverpay-user-service/internal/transport/http"
	"1994benc/neverpay-user-service/internal/user"
	"net/http"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type App struct {
	Name    string
	Version string
}

func main() {
	app := App{
		Name:    "Neverpay",
		Version: "1.0.0",
	}
	err := app.RunServer()
	if err != nil {
		log.Fatalf("Error starting the server %s", err)
	}
}

func (app *App) RunServer() error {
	app.setUpLogger()
	db, err := app.setUpDatabase()
	if err != nil {
		log.Fatalln("Error setting up database")
	}
	t := &user.TokenValidator{}
	userService := user.NewUserService(db, t)
	handler := transportHTTP.NewHandler(userService)
	handler.SetupRoutes()
	err = http.ListenAndServe(":8080", handler.Router)
	return err
}

func (app *App) setUpLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName":    app.Name,
			"AppVersion": app.Version,
		},
	).Info("Setting up app info")
}

func (app *App) setUpDatabase() (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	db, err = database.New()
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}
	err = database.MigrateDB(db)
	if err != nil {
		log.Fatalf("Error migrating DB: %s", err)
	}
	return db, err
}
