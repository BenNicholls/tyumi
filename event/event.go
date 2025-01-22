// Event is Tyumi's generic event system. It handles registration, creation, storage, filtering,
// and eventual doling-out of events.
package event

import "github.com/bennicholls/tyumi/log"

//this array is populated with the registered listeners for each registered event type.
var registeredEvents [][]*Stream
var eventNames []string

func init() {
	registeredEvents = make([][]*Stream, 0)
	eventNames = make([]string, 0)
}

//Definition for event objects. Compose custom events around the EventPrototype to satisfy
//this interface cleanly.
type Event interface {
	ID() int
	String() string
	Handled() bool
	SetHandled()
}

//Compose custom events around this
type EventPrototype struct {
	id      int
	handled bool
}

func (e EventPrototype) ID() int {
	return e.id
}

func (e EventPrototype) Handled() bool {
	return e.handled
}

func (e EventPrototype) String() (s string) {
	s = eventNames[e.id]
	if e.handled {
		s += " (handled)"
	}

	return
}

// Marks the event as handled. This doesn't prevent propogation/processing of the event on its own, but can be checked
// to dip out of event handling early if desired.
func (e *EventPrototype) SetHandled() {
	e.handled = true
}

func New(ID int) *EventPrototype {
	return &EventPrototype{id: ID}
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

// Adds an event handler function to the stream for event processing. If this is not set, the stream will not receive
// events.
func (s *Stream) AddHandler(h func(Event)) {
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

//Processes all events in the stream with the provided event handler function (if there is one).
func (s *Stream) Process() {
	if s.handler == nil {
		return
	}

	for e := s.Next(); e != nil; e = s.Next() {
		s.handler(e)
	}
}

//Registers an event type with the event system and returns the assigned ID. Also creates the list of
//listeners
func Register(name string) int {
	registeredEvents = append(registeredEvents, make([]*Stream, 0))
	eventNames = append(eventNames, name)

	return len(registeredEvents) - 1
}

//fire the event into the void. the event will be sent to all listening event streams
func Fire(e Event) {
	for _, s := range registeredEvents[e.ID()] {
		s.Add(e)
	}
}
