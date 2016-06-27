package eventsourcing

import (
	"log"
)

type EventStream struct {
	stream []Event
}

func NewStream() *EventStream {
	var stream = []Event{}

	return &EventStream{stream: stream}
}

func (es *EventStream) Stream() []Event {
	return es.stream
}

func (es *EventStream) Add(e Event) {
	log.Println(e)
	es.stream = append(es.stream, e)
}

func (es *EventStream) HasEvents() bool {
	return es.Count() > 0
}

func (es *EventStream) Count() int {
	return len(es.stream)
}
