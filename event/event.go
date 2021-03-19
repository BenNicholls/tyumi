//Event is Tyumi's generic event system. It handles registration, creation, storage, filtering,
//and eventual doling-out of events. 
package event

//this is a tracker for how many events have been registered. events are given ids in the order
//they are registered. eventually this will be an array that also holds per-event data, like
//filter status and other flags.
var registeredEvents int 

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
}

func NewStream(size int) (es *Stream) {
	es = new(Stream)
	es.stream = make(chan Event, size)
	return es
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

	return <- s.stream
}

//Registers an event with the event system and returns the assigned ID.
func Register() int {
	registeredEvents++
	return registeredEvents - 1
}