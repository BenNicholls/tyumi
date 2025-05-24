package event

import (
	"github.com/bennicholls/tyumi/util"
)

var registeredEvents []eventInfo

func init() {
	registeredEvents = make([]eventInfo, 1)
	Register("ERROR event (if you see this, someone did a boo boo.)", SIMPLE)
}

type eventType uint8

const (
	SIMPLE eventType = iota
	COMPLEX
)

// eventInfo describes each registered event type.
type eventInfo struct {
	id        EventID
	name      string
	eType     eventType
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
func Register(name string, event_type eventType) EventID {
	info := eventInfo{
		id:    EventID(len(registeredEvents)),
		name:  name,
		eType: event_type,
	}
	registeredEvents = append(registeredEvents, info)

	return info.id
}
