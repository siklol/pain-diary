package api

import (
	"customer"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"
)

type CustomerController struct {
	DB *sql.DB
}

func error(w http.ResponseWriter, errorMessage string) {
	result, _ := json.Marshal(map[string]interface{}{
		"error": errorMessage,
	})

	w.Write(result)
}

// Creates a new customer and returns a uuid
func (controller *CustomerController) Create(w http.ResponseWriter, r *http.Request) {
	var (
		result     []byte
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

	c = pes.CreateCustomer(customerId)
	c.CreateId(customerId)

	pes.Persist(c.Stream())
	result, _ = json.Marshal(map[string]interface{}{
		"customerId": customerId.String(),
		"error":      "",
	})

	w.Write(result)
}

// Updates customer profile
func (controller *CustomerController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var (
		result     []byte
		customerId uuid.UUID
		pes        *customer.EventStore
		c          *customer.Customer
	)

	pes = customer.CreateEventStore(controller.DB)

	if r.FormValue("customerId") == "" {
		error(w, "Customer not found! Customer ID not provided.")
		return
	}

	c = pes.RebuildCustomer(customerId)
	c.CreateId(customerId)

	pes.Persist(c.Stream())
	result, _ = json.Marshal(map[string]interface{}{
		"customerId": customerId.String(),
		"error":      "",
	})

	w.Write(result)
}
