package event

import "github.com/bennicholls/tyumi/log"

// A Handler is called when processing events in an eventstream. It takes 1 argument e: the event being processed.
// It is expected that, if the Handler successfully handles the event, it returns true
type Handler func(e Event) (handled bool)

// Listener defines anything that can listen for and process events.
type Listener interface {
	SetEventHandler(handler Handler)
	SetStreamSize(size int)
	ProcessEvents()
	FlushEvents()
	Listen(ids ...int)
	DeListen(ids ...int)
	StopListening()
}

// Stream is a queue of events. Use Listen() to have the stream collect events of a certain type, and then ProcessEvents()
// to call the assigned event handler on each accumulated event.
type Stream struct {
	stream  chan Event
	handler Handler //event handler called by Process()
}

// Initializes a stream. Size is the maximum number of events that can be accumulated before processing. Handler is the
// function called on the events during processing. If handler is nil, no events will be sent to this stream.
// Event streams do not need to be initialized explicitly; if not, a default size of 100 will be used.
func (s *Stream) Init(size int, handler Handler) {
	s.stream = make(chan Event, size)
	s.handler = handler
}

func NewStream(size int, handler Handler) (s Stream) {
	s.Init(size, handler)

	return
}

// Sets the event handler function to the stream for event processing. If this is not set, the stream will not receive
// events. If the stream's size defaults to 100 events. Use Init() or SetStreamSize() to set another size if needed.
func (s *Stream) SetEventHandler(handler Handler) {
	if s.stream == nil {
		s.Init(100, handler)
	} else {
		s.handler = handler
	}
}

// Sets the maximum number of events that the stream can hold before needing to be processed. If this is not called, a
// default value of 100 will be used.
func (s *Stream) SetStreamSize(size int) {
	if len(s.stream) > 0 {
		log.Warning("Setting stream size on an active event stream! All accumulated events flushed.")
	}

	s.Init(size, s.handler)
}

// Flushes all events from the stream instead of processing them.
func (s *Stream) FlushEvents() {
	for range len(s.stream) {
		<-s.stream
	}
}

// Begins listening for the specified event(s).
func (s *Stream) Listen(ids ...int) {
	for _, id := range ids {
		if !validID(id) {
			log.Warning("Attempted to listen for unregistered event ID: ", id)
			continue
		}
		registeredEvents[id].addListener(s)
	}
}

// DeListen will prevent the stream from receiving anymore of the specified events.
func (s *Stream) DeListen(ids ...int) {
	for id := range ids {
		if !validID(id) {
			continue
		}
		registeredEvents[id].removeListener(s)
	}
}

// Closes an event stream, effectively de-listening for all listened events. Also removes any assigned event handler.
func (s *Stream) StopListening() {
	for i := range registeredEvents {
		registeredEvents[i].removeListener(s)
	}

	s.handler = nil
}

// Processes all events in the stream with the provided event handler function (if there is one).
func (s *Stream) ProcessEvents() {
	if s.handler == nil {
		return
	}

	for event := s.next(); event != nil; event = s.next() {
		handled := s.handler(event)
		if handled {
			event.setHandled()
		}
	}
}

// Adds an event to the stream, unless the stream is full. If stream does not have an event handler, we assume
// that it can't handle events so we don't add anything.
func (s *Stream) add(e Event) {
	if s.handler == nil {
		return
	}

	if len(s.stream) == cap(s.stream) {
		log.Warning("Event stream full! Event not added. Either this means the stream is too small, or you've forgotten to close a stream that is no longer being processed.")
		return
	}

	s.stream <- e
}

// pops the next event and returns it. if there are no events, this will return nil
func (s *Stream) next() Event {
	if len(s.stream) == 0 {
		return nil
	}

	return <-s.stream
}
