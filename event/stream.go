package event

import "github.com/bennicholls/tyumi/log"

// A Handler is called when processing events in an eventstream. It takes 1 argument e: the event being processed.
// It is expected that, if the Handler successfully handles the event, it returns true
type Handler func(e Event) (handled bool)

// Stream is a queue of events.
type Stream struct {
	stream  chan Event
	handler Handler //event handler called by Process()
}

func NewStream(size int, handler Handler) (es Stream) {
	es.stream = make(chan Event, size)
	es.handler = handler

	return
}

// Adds an event handler function to the stream for event processing. If this is not set, the stream will not receive
// events.
func (s *Stream) AddHandler(h Handler) {
	s.handler = h
}

// Adds an event to the stream, unless the stream is full. If stream does not have an event handler, we assume
// that it can't handle events so we don't add anything.
func (s *Stream) Add(e Event) {
	if s.handler == nil {
		return
	}

	if len(s.stream) == cap(s.stream) {
		log.Warning("Event stream reached cap! No event added. Maybe you should make this bigger?!?!?")
		return
	}

	s.stream <- e
}

// Flushes all events from the stream instead of processing them.
func (s *Stream) Flush() {
	for range len(s.stream) {
		<-s.stream
	}
}

// pops the next event and returns it. if there are no events, this will return nil
func (s *Stream) Next() Event {
	if len(s.stream) == 0 {
		return nil
	}

	return <-s.stream
}

// Begins listening for the specified event.
func (s *Stream) Listen(ids ...int) {
	for _, id := range ids {
		if id < 0 || id >= len(registeredEvents) {
			log.Warning("Attempted to listen for unregistered event ID: ", id)
			continue
		}
		registeredEvents[id].addListener(s)
	}
}

// Processes all events in the stream with the provided event handler function (if there is one).
func (s *Stream) Process() {
	if s.handler == nil {
		return
	}

	for event := s.Next(); event != nil; event = s.Next() {
		handled := s.handler(event)
		if handled {
			event.setHandled()
		}
	}
}
