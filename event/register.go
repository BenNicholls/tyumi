package event

var registeredEvents []eventInfo

func init() {
	registeredEvents = make([]eventInfo, 1)
	Register("ERROR event (if you see this, someone did a boo boo.)")
}

// eventInfo describes each registered event type.
type eventInfo struct {
	id        int
	name      string
	listeners []*Stream
}

func (e *eventInfo) addListener(stream *Stream) {
	for _, listener := range e.listeners {
		if listener == stream {
			return
		}
	}

	e.listeners = append(e.listeners, stream)
}

// Registers an event type with the event system and returns the assigned ID.
// NOTE: name is NOT a unique identifier, only the returned ID is. Name is just for human readability.
func Register(name string) int {
	info := eventInfo{
		id:        len(registeredEvents),
		name:      name,
		listeners: make([]*Stream, 0),
	}
	registeredEvents = append(registeredEvents, info)

	return info.id
}
