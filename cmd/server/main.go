package main

import (
	"database/sql"
	"log"
	"net/http"
	"github.com/emptyset/simple-chat/internal/app"
)

func main() {
	// TODO: configure from environment variables
	/*
	// TODO: build connection string from environment variables
	database, err := sql.Open("mysql", "")
	if err != nil {
		log.Fatal("unable to connect to database", err)
	}
	defer database.Close()
	*/

	//handler, err := app.NewHandler(database)
	handler, err := app.NewHandler(sql.DB{})

	server := http.Server{
		Addr: ":8080",
		Handler: handler,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
