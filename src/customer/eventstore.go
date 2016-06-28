package customer

import (
	"database/sql"
	"encoding/json"
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
	for _, event := range stream.Stream() {
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

func (pes *EventStore) RebuildCustomer(customerId uuid.UUID) *Customer {
	rows, err := pes.db.Query(selectQuery, customerId.String())

	if err != nil {
		log.Fatal(err)
	}

	var (
		eventId        uuid.UUID
		eventData      string
		event          eventsourcing.Event
		eventCreatedAt time.Time
	)

	eventstream := eventsourcing.NewStream()

	for rows.Next() {
		rows.Scan(&eventId, &eventData, &eventCreatedAt)
		event = eventsourcing.RebuildEvent(customerId, eventId, eventData, eventCreatedAt, true)
		eventstream.Add(event)
	}

	customer := Customer{}
	customer.Replay(eventstream)

	return &customer
}
