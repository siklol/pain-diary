package api

import (
	"customer"
	"database/sql"
	"encoding/json"
	"github.com/satori/go.uuid"
	"net/http"
)

type CustomerController struct {
	DB *sql.DB
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

	c = pes.RebuildCustomer(customerId)
	c.CreateId(customerId)

	pes.Persist(c.Stream())
	result, _ = json.Marshal(map[string]interface{}{
		"customerId": customerId.String(),
		"error":      "",
	})

	w.Write(result)
}
