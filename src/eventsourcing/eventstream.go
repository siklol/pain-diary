package eventsourcing

type EventStream struct {
	stream []Event
}

func NewStream() *EventStream {
	var stream = []Event{}

	return &EventStream{stream: stream}
}

func (es *EventStream) Events() []Event {
	return es.stream
}

func (es *EventStream) Add(e Event) {
	es.stream = append(es.stream, e)
}

func (es *EventStream) HasEvents() bool {
	return es.Count() > 0
}

func (es *EventStream) Count() int {
	return len(es.stream)
}
