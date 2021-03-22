//Event is Tyumi's generic event system. It handles registration, creation, storage, filtering,
//and eventual doling-out of events.
package event


//this array is populated with the registered listeners for each registered event type.
var registeredEvents [][]*Stream

func init() {
	registeredEvents = make([][]*Stream, 0)
}

//Definition for event objects. Compose custom events around the EventPrototype to satisfy
//this interface cleanly.
type Event interface {
	ID() int
}

//Compose custom events around this
type EventPrototype struct {
	id int
}

func (e EventPrototype) ID() int {
	return e.id
}

func New(ID int) EventPrototype {
	return EventPrototype{id: ID}
}

//Stream is a queue of events.
type Stream struct {
	stream chan Event

	handler func(Event) //event handler called by Process()
}

func NewStream(size int) (es Stream) {
	es.stream = make(chan Event, size)
	return
}

func (s *Stream) AddHandler(h func(Event)) {
	s.handler = h
}

//Adds an event to the stream, unless the stream is full.
//TODO: handle full event stream behaviour somehow.
func (s *Stream) Add(e Event) {
	if len(s.stream) == cap(s.stream) {
		return
	}
	s.stream <- e
}

//pops the next event and returns it. if there are no events, this will return nil
func (s *Stream) Next() Event {
	if len(s.stream) == 0 {
		return nil
	}

	return <-s.stream
}

//Begins listening for the specified event.
//TODO: check for duplicate listens
func (s *Stream) Listen(id int) {
	registeredEvents[id] = append(registeredEvents[id], s)
}

//Processes all events in the stream with the provided event handler function. If no event handler
//exists, does nothing.
//THINK: should this clear the events from the stream if no handler is present?
func (s *Stream) Process() {
	if s.handler == nil {
		return
	}

	for e := s.Next(); e != nil; e = s.Next() {
		s.handler(e)
	}
}

//Registers an event with the event system and returns the assigned ID. Also creates the list of
//listeners
func Register() int {
	registeredEvents = append(registeredEvents, make([]*Stream, 0))
	return len(registeredEvents) - 1
}

//fire the event into the void. the event will be sent to all listening event streams
func Fire(e Event) {
	for _, s := range registeredEvents[e.ID()] {
		s.Add(e)
	}
}
