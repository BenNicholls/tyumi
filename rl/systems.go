package rl

import (
	"github.com/bennicholls/tyumi/event"
)

type System struct {
	event.Stream

	Enabled bool
}

func (s *System) setEnabled(enabled bool) {
	if s.Enabled == enabled {
		return
	}

	s.Enabled = enabled
	if s.Enabled {
		s.EnableListening()
	} else {
		s.DisableListening()
	}
}

func (s *System) Enable() {
	s.setEnabled(true)
}

func (s *System) Disable() {
	s.setEnabled(false)
}

func (s *System) Update() {
	s.Stream.ProcessEvents()
}

func (s *System) Shutdown() {
	s.DisableListening()
}
