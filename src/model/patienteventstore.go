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

func CreateEventStore(database *sql.DB) *PatientEventStore {
	return &PatientEventStore{db: database}
}

func (pes *PatientEventStore) Persist(patient *Patient) {
	for _, event := range patient.Stream() {
		if !event.IsPersisted() {
			s, _ := json.Marshal(event.Payload())
			_, err := pes.db.Exec("INSERT INTO patienteventstore (eventid, patientid, eventdata, createdat) VALUES ($1, $2, $3, $4)", event.Id().String(), patient.patientId.String(), s, event.CreatedAt().Format("2006-01-02T15:04:05.000000Z"))

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

	var (
		eventId   uuid.UUID
		eventData string
		event     eventsourcing.Event
		eventCreatedAt time.Time
	)

	eventstream := eventsourcing.NewStream()

	for rows.Next() {
		rows.Scan(&eventId, &eventData, &eventCreatedAt)
		event = eventsourcing.RebuildEvent(eventId, eventData, eventCreatedAt, true)
		eventstream.Add(event)
	}

	patient := Patient{}
	patient.Replay(eventstream)

	return &patient
}
