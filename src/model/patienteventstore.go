package patient

import (
	"database/sql"
	"encoding/json"
	"eventsourcing"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type PatientEventStore struct {
	db *sql.DB
}

const (
	TIMEFORMAT = "2006-01-02T15:04:05.000000Z"
)

func CreateEventStore(database *sql.DB) *PatientEventStore {
	return &PatientEventStore{db: database}
}

func (pes *PatientEventStore) Persist(patient *Patient) {
	for _, event := range patient.Stream() {
		if !event.IsPersisted() {
			s, _ := json.Marshal(event.Payload())
			_, err := pes.db.Exec("INSERT INTO patienteventstore (eventid, patientid, eventdata, createdat) VALUES ($1, $2, $3, $4)", event.Id().String(), patient.patientId.String(), s, event.CreatedAt().Format(TIMEFORMAT))

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (pes *PatientEventStore) RebuildPatient(patientId uuid.UUID) *Patient {
	rows, err := pes.db.Query("SELECT eventid, eventdata, createdat FROM patienteventstore WHERE patientid = $1 ORDER BY createdat", patientId.String())

	if err != nil {
		log.Fatal(err)
	}

	eventstream := eventsourcing.NewStream()

	var (
		eventId   uuid.UUID
		eventData string
		event     eventsourcing.Event
		eventCreatedAt time.Time
	)

	for rows.Next() {
		rows.Scan(&eventId, &eventData, &eventCreatedAt)
		event = eventsourcing.NewEvent(eventId, eventData, eventCreatedAt, true)
		eventstream.Add(event)
	}

	patient := Patient{}
	patient.eventStream = eventstream
	stream := eventstream.Stream()

	for _, e := range stream {
		switch true {
		case e.Name() == "Event.CreateId":
			patient.patientId, _ = uuid.FromString(e.Payload()["patientId"].(string))
			patient.createdAt, _ = time.Parse(TIMEFORMAT, e.Payload()["createdAt"].(string))
			break
		case e.Name() == "Event.ChangeName":
			patient.firstname = e.Payload()["firstname"].(string)
			patient.lastname = e.Payload()["lastname"].(string)
			break
		}
	}

	return &patient
}
