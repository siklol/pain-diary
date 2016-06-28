package main

import (
	"database/sql"
	"net/http"
	"fmt"
	"log"
	"api"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DATABASE_NAME     = "rbl"
	DATABASE_USER     = "rbl"
	DATABASE_PASSWORD = "rbl"
	SERVER_PORT       = ":8080"
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", DATABASE_NAME, DATABASE_USER, DATABASE_PASSWORD, "disable"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	customerController := api.CustomerController{db}

	router := mux.NewRouter()
	router.HandleFunc("/customer/create", customerController.Create).Methods("POST")

	log.Fatal(http.ListenAndServe(SERVER_PORT, router))
}
