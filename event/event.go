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
	setID(id EventID)
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

func (e EventPrototype) ID() EventID {
	return e.id
}

func (e *EventPrototype) setID(id EventID) {
	e.id = id
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

// Fire an event into the void. The event will be sent to all listening event streams. Optionally lets you provide
// events to fire; use this to fire complex events that you create yourself. All provided events will have their IDs
// set to the provided ID. If no event is provided, a simple event with the provided ID will be fired.
func Fire(ID EventID, events ...Event) {
	if !ID.valid() {
		log.Error("Attempted to fire unregistered event with ID ", ID)
		return
	}

	if len(events) == 0 {
		e := EventPrototype{id: ID}
		for stream := range registeredEvents[ID].listeners.EachElement() {
			stream.add(&e)
		}
	} else { // no provided event, fire a simple event with just the id
		for _, e := range events {
			e.setID(ID)
			for stream := range registeredEvents[ID].listeners.EachElement() {
				stream.add(e)
			}
		}
	}
}
