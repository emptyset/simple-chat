package main

import (
	"database/sql"
	"log"
	"net/http"
	"github.com/emptyset/simple-chat/internal/app"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// TODO: configure from environment variables
	// TODO: build connection string from environment variables
	database, err := sql.Open("mysql", "root:password@tcp(database:3306)/database")
	if err != nil {
		log.Fatal("unable to connect to database", err)
	}

	handler, err := app.NewHandler(database)

	server := http.Server{
		Addr: ":8080",
		Handler: handler,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
