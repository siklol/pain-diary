package api

import (
	"customer"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"log"
)

type CustomerController struct {
	DB *sql.DB
}

// Creates a new customer and returns a uuid
func (controller *CustomerController) Create(w http.ResponseWriter, r *http.Request) {
	var (
		customerId uuid.UUID
		pes        *customer.EventStore
		c          *customer.Customer
	)

	// TODO add query filter for customer id uniqueness

	pes = customer.CreateEventStore(controller.DB)

	if r.FormValue("customerId") != "" {
		customerId, _ = uuid.FromString(r.FormValue("customerId"))
	} else {
		customerId = uuid.NewV4()
	}

	c = pes.Create(customerId)
	c.CreateId(customerId)

	pes.Persist(c.Stream())
	Success(w, map[string]interface{}{
		"customerId": customerId.String(),
	})
}

// Updates customer profile
func (controller *CustomerController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	pes := customer.CreateEventStore(controller.DB)

	if r.FormValue("customerId") == "" {
		Error(w, "Customer not found! Customer ID not provided.")
		return
	}
	customerId, _ := uuid.FromString(r.FormValue("customerId"))

	if r.FormValue("firstname") == "" || r.FormValue("lastname") == "" {
		Error(w, "Firstname or lastname not found!")
		return
	}

	c, err := pes.Rebuild(customerId)

	if err != nil {
		log.Fatal(err)
		Error(w, "Customer rebuild failed!")
		return
	}

	c.ChangeName(r.FormValue("firstname"), r.FormValue("lastname"))

	pes.Persist(c.Stream())
	Success(w, map[string]interface{}{
		"message": "Update customer profile",
	})
}

// Adds a pain level
func (controller *CustomerController) ExperiencePain(w http.ResponseWriter, r *http.Request) {
	pes := customer.CreateEventStore(controller.DB)

	if r.FormValue("customerId") == "" {
		Error(w, "Customer not found! Customer ID not provided.")
		return
	}
	customerId, _ := uuid.FromString(r.FormValue("customerId"))

	pain := r.FormValue("pain")
	if pain == "" {
		Error(w, "Pain not provided.")
		return
	}

	c, err := pes.Rebuild(customerId)

	if err != nil {
		log.Fatal(err)
		Error(w, "Customer rebuild failed!")
		return
	}

	c.ExperiencePain(pain)

	pes.Persist(c.Stream())
	Success(w, map[string]interface{}{
		"message": fmt.Sprintf("Pain %s was logged", pain),
	})

	log.Println(c)
}
