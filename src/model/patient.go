package patient

import (
	"eventsourcing"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

type Patient struct {
	eventStream *eventsourcing.EventStream

	patientId uuid.UUID
	firstname string
	lastname  string
	createdAt time.Time
}

func (patient *Patient) change(event eventsourcing.Event) {
	patient.eventStream.Add(event)
}

func (patient *Patient) Stream() []eventsourcing.Event {
	return patient.eventStream.Stream()
}

func (patient *Patient) CreateId(patientId uuid.UUID) {
	patient.patientId = patientId
	patient.createdAt = time.Now()

	eventData := fmt.Sprintf("{\"name\": \"Event.CreateId\", \"patientId\": \"%s\", \"createdAt\": \"%s\"}", patientId, patient.createdAt.String())
	patient.change(eventsourcing.NewEvent(uuid.NewV4(), eventData, time.Now(), false))
}

func (patient *Patient) ChangeName(firstname string, lastname string) {
	patient.firstname = firstname
	patient.lastname = lastname

	eventData := fmt.Sprintf("{\"name\": \"Event.ChangeName\", \"firstname\": \"%s\", \"lastname\": \"%s\"}", patient.firstname, patient.lastname)
	patient.change(eventsourcing.NewEvent(uuid.NewV4(), eventData, time.Now(), false))
}
