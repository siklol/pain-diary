package main

import (
	"../src/model"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", "test", "test", "test123", "disable"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	pes := patient.CreateEventStore(db)

	v4, _ := uuid.FromString("af282579-f8e7-4fd8-877c-183573de608b")
	patient := pes.RebuildPatient(v4)

	// patient.CreateId(uuid.NewV4())
	// patient.ChangeName("Peter", "Mustermann")
	patient.ExperiencePain("9")

	pes.Persist(patient)

	log.Println(patient)
}
