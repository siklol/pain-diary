package eventsourcing

import "github.com/satori/go.uuid"

type EventStore interface {
	Persist(es *EventStream)
	Stream(uuid uuid.UUID)
}
