package event

import (
	"weak"

	"github.com/bennicholls/tyumi/util"
)

var registeredEvents []eventInfo

func init() {
	registeredEvents = make([]eventInfo, 0)
	Register("ERROR event (if you see this, someone did a boo boo.)")
}

// eventInfo describes each registered event type.
type eventInfo struct {
	id        EventID
	name      string
	listeners util.Set[weak.Pointer[Stream]]
}

func (e *eventInfo) addListener(stream *Stream) {
	e.listeners.Add(weak.Make(stream))
}

func (e *eventInfo) removeListener(stream *Stream) {
	e.listeners.Remove(weak.Make(stream))
}

// Registers an event type with the event system and returns the assigned ID.
// NOTE: name is NOT a unique identifier, only the returned ID is. Name is just for human readability.
func Register(name string) EventID {
	info := eventInfo{
		id:   EventID(len(registeredEvents)),
		name: name,
	}
	registeredEvents = append(registeredEvents, info)

	return info.id
}
