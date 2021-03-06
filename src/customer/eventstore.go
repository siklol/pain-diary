package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"eventsourcing"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

var (
	insertQuery = "INSERT INTO customereventstore (eventid, customerid, eventdata, createdat) VALUES ($1, $2, $3, $4)"
	selectQuery = "SELECT eventid, eventdata, createdat FROM customereventstore WHERE customerid = $1 ORDER BY createdat"
)

type EventStore struct {
	db *sql.DB
}

func CreateEventStore(database *sql.DB) *EventStore {
	return &EventStore{db: database}
}

func (pes *EventStore) Persist(stream *eventsourcing.EventStream) {
	for _, event := range stream.Events() {
		if !event.IsPersisted() {
			s, _ := json.Marshal(event.Payload())
			eventId := event.Id().String()
			customerId := event.CustomerId().String()
			createdAt := event.CreatedAt().Format("2006-01-02T15:04:05.000000Z")
			_, err := pes.db.Exec(insertQuery, eventId, customerId, s, createdAt)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (pes *EventStore) Create(aggregateId uuid.UUID) *Customer {
	c := Customer{}
	c.Replay(eventsourcing.NewStream())
	return &c
}

func (pes *EventStore) Rebuild(aggregate uuid.UUID) (*Customer, error) {
	rows, err := pes.db.Query(selectQuery, aggregate.String())
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	eventstream := eventsourcing.NewStream()

	var (
		eventId         uuid.UUID
		eventData       string
		eventCreatedAt  time.Time
		aggregateExists bool
	)

	for rows.Next() {
		rows.Scan(&eventId, &eventData, &eventCreatedAt)
		event := eventsourcing.RebuildEvent(aggregate, eventId, eventData, eventCreatedAt, true)
		eventstream.Add(event)
		aggregateExists = true
	}

	if aggregateExists == false {
		err = errors.New("Customer does not exist")
	}

	customer := Customer{}
	customer.Replay(eventstream)

	return &customer, err
}
