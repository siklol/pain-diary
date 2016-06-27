package patient

import (
	"eventsourcing"
	"github.com/satori/go.uuid"
	"time"
	"log"
)

type Patient struct {
	eventStream *eventsourcing.EventStream

	patientId uuid.UUID
	firstname string
	lastname  string
	createdAt time.Time
	pain string
}

func (patient *Patient) Replay(eventStream *eventsourcing.EventStream) {
	patient.eventStream = eventStream

	patient.mutate()
}

func (patient *Patient) apply(event eventsourcing.Event) {
	patient.eventStream.Add(event)

	patient.mutate()
}

func (patient *Patient) mutate() {
	stream := patient.eventStream.Stream()
	var err error

	for _, e := range stream {
		switch true {
		case e.Name() == "Event.CreateId":
			patient.patientId, _ = uuid.FromString(e.Payload()["patientId"].(string))
			patient.createdAt, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", e.Payload()["createdAt"].(string))
			if err != nil {
				log.Fatal(err)
			}
			break
		case e.Name() == "Event.ChangeName":
			patient.firstname = e.Payload()["firstname"].(string)
			patient.lastname = e.Payload()["lastname"].(string)
			break
		case e.Name() == "Event.ExperiencePain":
			patient.pain = e.Payload()["pain"].(string)
			break
		}
	}
}

func (patient *Patient) Stream() []eventsourcing.Event {
	return patient.eventStream.Stream()
}

func (patient *Patient) CreateId(patientId uuid.UUID) {
	patient.apply(eventsourcing.NewEvent(uuid.NewV4(), map[string]interface{} {
		"name": "Event.CreateId",
		"patientId": patientId.String(),
		"createdAt": patient.createdAt.String(),
	}))
}

func (patient *Patient) ChangeName(firstname string, lastname string) {
	patient.apply(eventsourcing.NewEvent(uuid.NewV4(), map[string]interface{} {
		"name": "Event.ChangeName",
		"firstname": patient.firstname,
		"lastname": patient.lastname,
	}))
}

func (patient *Patient) ExperiencePain(pain string) {
	patient.apply(eventsourcing.NewEvent(uuid.NewV4(), map[string]interface{} {
		"name": "Event.ExperiencePain",
		"pain": pain,
	}))
}
