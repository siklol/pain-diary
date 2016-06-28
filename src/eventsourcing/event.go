package eventsourcing

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
)

type Event struct {
	id          uuid.UUID
	payload     map[string]interface{}
	createdAt   time.Time
	isPersisted bool
}

func NewEvent(eventId uuid.UUID, eventData map[string]interface{}) Event {
	e := Event{id: eventId, createdAt: time.Now(), isPersisted: false}
	e.payload = eventData

	return e
}

func RebuildEvent(eventId uuid.UUID, s string, createdAt time.Time, isPersisted bool) Event {
	payload := []byte(s)
	var f map[string]interface{}
	e := Event{id: eventId}

	json.Unmarshal(payload, &f)
	e.payload = f
	e.createdAt = createdAt
	e.isPersisted = isPersisted

	return e
}

func (e *Event) Payload() map[string]interface{} {
	return e.payload
}

func (e *Event) Id() uuid.UUID {
	return e.id
}

func (e *Event) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Event) Name() string {
	return e.payload["name"].(string)
}

func (e *Event) IsPersisted() bool {
	return e.isPersisted
}
