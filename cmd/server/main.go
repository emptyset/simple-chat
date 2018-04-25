package main

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/emptyset/simple-chat/internal/app"
	"github.com/emptyset/simple-chat/internal/models"
	"github.com/emptyset/simple-chat/internal/storage"
)

func main() {
	// TODO: configure from environment variables
	log.SetLevel(log.DebugLevel)

	// TODO: build connection string from environment variables
	log.Debug("opening mysql database")
	database, err := sql.Open("mysql", "root:password@tcp(database:3306)/chat?charset=utf8")
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}

	retries := 0
	for retries < 5 {
		delay := 1 << uint(retries)
		time.Sleep(time.Duration(delay) * time.Second)

		err = database.Ping()
		if err != nil {
			log.Errorf("unable to ping database: %s", err)
		} else {
			break
		}

		retries++
	}

	log.Debug("creating sql data store")
	store := storage.NewSQLDataStore(database)
	model := models.New(store)
	handler, err := app.NewHandler(model)
	if err != nil {
		log.Fatalf("unable to instantiate handler: %s", err)
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Debug("starting server")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to start server: %s", err)
	}
}
