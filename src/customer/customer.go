package customer

import (
	"eventsourcing"
	"log"
	"time"

	"github.com/satori/go.uuid"
	"strconv"
)

type Customer struct {
	eventStream *eventsourcing.EventStream

	customerId uuid.UUID
	firstname  string
	lastname   string
	createdAt  time.Time
	pain       []int64
}

func (customer *Customer) Replay(eventStream *eventsourcing.EventStream) {
	customer.eventStream = eventStream

	customer.mutate()
}

func (customer *Customer) apply(event eventsourcing.Event) {
	customer.eventStream.Add(event)

	customer.mutate()
}

func (customer *Customer) mutate() {
	stream := customer.eventStream.Stream()
	var err error

	for _, e := range stream {
		switch true {
		case e.Name() == "Event.CreateId":
			customer.customerId, _ = uuid.FromString(e.Payload()["customerId"].(string))
			customer.createdAt, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", e.Payload()["createdAt"].(string))
			if err != nil {
				log.Fatal(err)
			}
			break
		case e.Name() == "Event.ChangeName":
			customer.firstname = e.Payload()["firstname"].(string)
			customer.lastname = e.Payload()["lastname"].(string)
			break
		case e.Name() == "Event.ExperiencePain":
			i, _ := strconv.ParseInt(e.Payload()["pain"].(string), 10, 64)
			customer.pain = append(customer.pain, i)
			break
		}
	}
}

func (customer *Customer) Stream() *eventsourcing.EventStream {
	return customer.eventStream
}

func (customer *Customer) CreateId(customerId uuid.UUID) {
	customer.apply(eventsourcing.NewEvent(customerId, uuid.NewV4(), map[string]interface{}{
		"name":       "Event.CreateId",
		"customerId": customerId.String(),
		"createdAt":  time.Now().String(),
	}))
}

func (customer *Customer) ChangeName(firstname string, lastname string) {
	customer.apply(eventsourcing.NewEvent(customer.customerId, uuid.NewV4(), map[string]interface{}{
		"name":      "Event.ChangeName",
		"firstname": firstname,
		"lastname":  lastname,
	}))
}

func (customer *Customer) ExperiencePain(pain string) {
	customer.apply(eventsourcing.NewEvent(customer.customerId, uuid.NewV4(), map[string]interface{}{
		"name": "Event.ExperiencePain",
		"pain": pain,
	}))
}
