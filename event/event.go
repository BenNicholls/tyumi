// Event is Tyumi's generic event system. It handles registration, creation, storage, filtering,
// and eventual doling-out of events.
package event

import "github.com/bennicholls/tyumi/log"

// Definition for event objects. Compose custom events around the EventPrototype to satisfy
// this interface cleanly.
type Event interface {
	ID() EventID
	String() string
	Handled() bool
	setHandled()
}

type EventID uint32

func (id EventID) valid() bool {
	return int(id) < len(registeredEvents)
}

// Compose custom events around this
type EventPrototype struct {
	id      EventID
	handled bool
}

func New(ID EventID) EventPrototype {
	if !ID.valid() {
		log.Warning("Attempted to create event with unregistered ID: ", ID)
		return EventPrototype{id: 0}
	}

	return EventPrototype{id: ID}
}

func (e EventPrototype) ID() EventID {
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
	if !e.ID().valid() {
		log.Error("Attempted to fire unregistered event with ID ", e.ID())
		return
	}

	for s := range registeredEvents[e.ID()].listeners.EachElement() {
		s.add(e)
	}
}

// fire a simple event into the void. Produces and error if the event was not registered as a simple event.
func FireSimple(ID EventID) {
	if registeredEvents[ID].eType != SIMPLE {
		log.Error("Attempted to fire complex event with FireSimple(), id: ", ID)
		return
	}

	simpleEvent := EventPrototype{id: ID}
	Fire(&simpleEvent)
}
