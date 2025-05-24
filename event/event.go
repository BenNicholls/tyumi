// Event is Tyumi's generic event system. It handles registration, creation, storage, filtering,
// and eventual doling-out of events.
package event

import "github.com/bennicholls/tyumi/log"

// Definition for event objects. Compose custom events around the EventPrototype to satisfy
// this interface cleanly.
type Event interface {
	ID() int
	String() string
	Handled() bool
	setHandled()
}

// Compose custom events around this
type EventPrototype struct {
	id      int
	handled bool
}

func New(ID int) EventPrototype {
	if !validID(ID) {
		log.Warning("Attempted to create event with unregistered ID: ", ID)
		return EventPrototype{id: 0}
	}

	return EventPrototype{id: ID}
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
func (e *EventPrototype) setHandled() {
	e.handled = true
}

// fire the event into the void. the event will be sent to all listening event streams
func Fire(e Event) {
	if !validID(e.ID()) {
		log.Error("Attempted to fire unregistered event with ID ", e.ID())
		return
	}

	for s := range registeredEvents[e.ID()].listeners.EachElement() {
		s.add(e)
	}
}

// fire a simple event into the void. Produces and error if the event was not registered as a simple event.
func FireSimple(ID int) {
	if registeredEvents[ID].eType != SIMPLE {
		log.Error("Attempted to fire complex event with FireSimple(), id: ", ID)
		return
	}

	simpleEvent := EventPrototype{id: ID}
	Fire(&simpleEvent)
}

func validID(ID int) bool {
	return ID < len(registeredEvents) && ID > 0
}
