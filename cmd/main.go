package main

import (
	"../src/customer"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
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

	// TODO remove debug
	db.Exec("TRUNCATE customereventstore;")

	pes := customer.CreateEventStore(db)

	v4, _ := uuid.FromString("af282579-f8e7-4fd8-877c-183573de608b")
	customer := pes.RebuildCustomer(v4)

	customer.CreateId(v4)
	customer.ChangeName("Peter", "Mustermann")
	customer.ExperiencePain("9")

	pes.Persist(customer.Stream())

	log.Println(customer)
}
