package main

import (
	"database/sql"
	"fmt"
	"log"

	"api"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

const (
	DATABASE_NAME     = "rbl"
	DATABASE_USER     = "rbl"
	DATABASE_PASSWORD = "rbl"
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", DATABASE_NAME, DATABASE_USER, DATABASE_PASSWORD, "disable"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	customerController := api.CustomerController{db}

	router := mux.NewRouter()
	router.HandleFunc("/", customerController.Index)
	router.HandleFunc("/customer/create", customerController.Create).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
