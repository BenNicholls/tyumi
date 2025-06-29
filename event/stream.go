package event

import (
	"slices"

	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
)

// A Handler is called when processing events in an eventstream. It takes 1 argument e: the event being processed.
// It is expected to return true if the Handler successfully handles the event.
type Handler func(e Event) (handled bool)

// Listener defines anything that can listen for and process events.
type Listener interface {
	SetEventHandler(handler Handler)
	SetStreamSize(size int)
	ProcessEvents()
	FlushEvents()
	Listen(ids ...EventID)
	DeListen(ids ...EventID)
	EnableListening()
	DisableListening()
}

// SuppressionMode defines how duplicate events are handled in an event stream.
type SuppressionMode uint8

const (
	DoNotSuppress SuppressionMode = iota // Default suppression mode. Duplicate events are not suppressed.
	KeepFirst                            // Duplicate events are not added, the first instance of the event is kept in the stream
	KeepLast                             // Duplicate events replace earlier matching events in the stream.
)

// Stream is a queue of events. Use Listen() to have the stream collect events of a certain type, and then ProcessEvents()
// to call the assigned event handler on each accumulated event.
type Stream struct {
	stream       []Event
	handler      Handler           //event handler called by Process()
	listenIDs    util.Set[EventID] // ids that are currently being listened for
	disabled     bool              // whether the stream accepts events
	suppression  SuppressionMode   // type of duplicate event suppression
	eventIndices map[EventID]int   // indices of events with specific IDs, used while suppressing duplicate events
}

// Initializes a stream. Size is the maximum number of events that can be accumulated before processing. Handler is the
// function called on the events during processing. If handler is nil, no events will be sent to this stream.
// Event streams do not need to be initialized explicitly; if not, a default size of 100 will be used.
func (s *Stream) Init(size int, handler Handler) {
	s.stream = make([]Event, 0, size)
	s.handler = handler
}

func NewStream(size int, handler Handler) (s Stream) {
	s.Init(size, handler)

	return
}

// Sets the event handler function to the stream for event processing. If this is not set, the stream will not receive
// events.
func (s *Stream) SetEventHandler(handler Handler) {
	if s.stream == nil {
		s.Init(100, handler)
	} else {
		s.handler = handler
	}
}

// Sets the maximum number of events that the stream can hold before needing to be processed. If this is not called then
// a default value of 100 will be used.
func (s *Stream) SetStreamSize(size int) {
	if size <= 0 {
		log.Error("Attempting to set stream size to 0 or less. Don't do that.")
		return
	}

	if len(s.stream) > 0 {
		log.Warning("Setting stream size on an active event stream! All accumulated events flushed.")
	}

	s.Init(size, s.handler)
}

// Sets the type of duplicate event suppression. If suppression is enabled, the stream will discard events with the same
// EventID as an event already in the stream. Possible values are:
//
//	DoNotSuppress: allows duplicate events
//	KeepFirst: suppress duplicates, keeping the first instance of the event in the stream
//	KeepLast: suppress duplicates, replacing earlier instances of an event with the latest one
func (s *Stream) SuppressDuplicateEvents(mode SuppressionMode) {
	s.suppression = mode

	if s.suppression != DoNotSuppress {
		s.eventIndices = make(map[EventID]int)
	} else {
		s.eventIndices = nil
	}
}

// Clears the event stream of all collected events.
func (s *Stream) FlushEvents() {
	s.stream = slices.Delete(s.stream, 0, len(s.stream))
	if s.suppression != DoNotSuppress {
		clear(s.eventIndices)
	}
}

// Begins listening for the specified event(s).
func (s *Stream) Listen(ids ...EventID) {
	for _, id := range ids {
		if !id.valid() {
			log.Warning("Attempted to listen for unregistered event ID: ", id)
			continue
		}

		s.listenIDs.Add(id)
		if !s.disabled {
			registeredEvents[id].addListener(s)
		}
	}
}

// DeListen will prevent the stream from receiving anymore of the specified events.
func (s *Stream) DeListen(ids ...EventID) {
	for _, id := range ids {
		if !id.valid() {
			continue
		}

		s.listenIDs.Remove(id)
		if !s.disabled {
			registeredEvents[id].removeListener(s)
		}
	}
}

// Enables listening for events. Streams default to be listening so this does NOT need to be called to activate the
// stream. Use this only to re-activate a stream that has been manually disabled with DisableListening().
func (s *Stream) EnableListening() {
	s.setDisabled(false)
}

// Disables listening for events. Disabled streams will not receive events and ProcessEvents() becomes a no-op. Use
// EnableListening() to reactivate the stream.
func (s *Stream) DisableListening() {
	s.setDisabled(true)
}

func (s *Stream) setDisabled(disabled bool) {
	if s.disabled == disabled {
		return
	}

	s.disabled = disabled
	if s.disabled {
		s.FlushEvents()
		for i := range s.listenIDs.EachElement() {
			registeredEvents[i].removeListener(s)
		}
	} else {
		for i := range s.listenIDs.EachElement() {
			registeredEvents[i].addListener(s)
		}
	}
}

// Processes all events in the stream with the provided event handler function (if there is one).
func (s *Stream) ProcessEvents() {
	if s.disabled || s.handler == nil {
		return
	}

	for _, event := range s.stream {
		handled := s.handler(event)
		if handled {
			event.setHandled()
		}
	}

	s.FlushEvents()
}

// Adds an event to the stream, unless the stream is full. If stream does not have an event handler, we assume
// that it can't handle events so we don't add anything.
func (s *Stream) add(e Event) {
	if s.handler == nil {
		return
	}

	if len(s.stream) == cap(s.stream) {
		log.Warning(`Event stream full! Event not added. Either this means the stream is too small,
						or you've forgotten to close a stream that is no longer being processed.`)
		return
	}

	if s.suppression != DoNotSuppress {
		// if duplicate ID found, either replace or just return depending on mode.
		if i, ok := s.eventIndices[e.ID()]; ok {
			if s.suppression == KeepLast {
				s.stream[i] = e
			}

			return
		} else {
			s.eventIndices[e.ID()] = len(s.stream)
		}
	}

	s.stream = append(s.stream, e)
}
