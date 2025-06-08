package rl

import (
	"github.com/bennicholls/tyumi/event"
)

type System struct {
	event.Stream
}

func (s *System) Update() {
	s.Stream.ProcessEvents()
}

func (s *System) Shutdown() {
	s.DisableListening()
}
