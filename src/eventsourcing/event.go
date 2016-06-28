package eventsourcing

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
)

type Event struct {
	id          uuid.UUID
	customerId  uuid.UUID
	payload     map[string]interface{}
	createdAt   time.Time
	isPersisted bool
}

func NewEvent(customerId uuid.UUID, eventId uuid.UUID, eventData map[string]interface{}) Event {
	return Event{
		id:          eventId,
		customerId:  customerId,
		createdAt:   time.Now(),
		isPersisted: false,
		payload:     eventData,
	}
}

func RebuildEvent(customerId uuid.UUID, eventId uuid.UUID, s string, createdAt time.Time, isPersisted bool) Event {
	payload := []byte(s)
	var payloadJson map[string]interface{}
	json.Unmarshal(payload, &payloadJson)

	return Event{
		id:          eventId,
		customerId:  customerId,
		createdAt:   createdAt,
		isPersisted: isPersisted,
		payload:     payloadJson,
	}
}

func (e *Event) Payload() map[string]interface{} {
	return e.payload
}

func (e *Event) Id() uuid.UUID {
	return e.id
}

func (e *Event) CustomerId() uuid.UUID {
	return e.customerId
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
