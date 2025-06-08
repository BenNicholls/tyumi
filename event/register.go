package event

import (
	"github.com/bennicholls/tyumi/util"
)

var registeredEvents []eventInfo

func init() {
	registeredEvents = make([]eventInfo, 1)
	Register("ERROR event (if you see this, someone did a boo boo.)")
}

// eventInfo describes each registered event type.
type eventInfo struct {
	id        EventID
	name      string
	listeners util.Set[*Stream]
}

func (e *eventInfo) addListener(stream *Stream) {
	e.listeners.Add(stream)
}

func (e *eventInfo) removeListener(stream *Stream) {
	e.listeners.Remove(stream)
}

// Registers an event type with the event system and returns the assigned ID.
// event_type is SIMPLE or COMPLEX. simple events just contain the ID, complex events are composed around a larger
// event struct with additional data.
// NOTE: name is NOT a unique identifier, only the returned ID is. Name is just for human readability.
func Register(name string) EventID {
	info := eventInfo{
		id:   EventID(len(registeredEvents)),
		name: name,
	}
	registeredEvents = append(registeredEvents, info)

	return info.id
}
