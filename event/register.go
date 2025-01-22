package event

var registeredEvents []eventInfo

func init() {
	registeredEvents = make([]eventInfo, 0)
}

// eventInfo describes each registered event type.
type eventInfo struct {
	id        int
	name      string
	listeners []*Stream
}

func (e *eventInfo) addListener(s *Stream) {
	for _, l := range e.listeners {
		if l == s {
			return
		}
	}

	e.listeners = append(e.listeners, s)
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
