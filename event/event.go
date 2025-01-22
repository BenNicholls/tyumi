// Event is Tyumi's generic event system. It handles registration, creation, storage, filtering,
// and eventual doling-out of events.
package event

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

func New(ID int) *EventPrototype {
	return &EventPrototype{id: ID}
}

func (e EventPrototype) ID() int {
	return e.id
}

func (e EventPrototype) Handled() bool {
	return e.handled
}

func (e EventPrototype) String() (s string) {
	s = registeredEvents[e.id].name
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

//fire the event into the void. the event will be sent to all listening event streams
func Fire(e Event) {
	for _, s := range registeredEvents[e.ID()].listeners {
		s.Add(e)
	}
}
