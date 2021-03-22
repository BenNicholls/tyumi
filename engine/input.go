package engine

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/veandco/go-sdl2/sdl"
)

//gather input events from sdl and handle/distribute accordingly
func processInput() {
	for sdlevent := sdl.PollEvent(); sdlevent != nil; sdlevent = sdl.PollEvent() {
		switch e := sdlevent.(type) {
		case *sdl.QuitEvent:
			event.Fire(event.New(EV_QUIT))
			break //don't care about other input events if we're quitting
		case *sdl.KeyboardEvent:
			event.Fire(NewKeyboardEvent(int(e.Keysym.Sym)))
		}
	}

	return
}

var EV_KEYBOARD = event.Register()

type KeyboardEvent struct {
	event.EventPrototype

	Key int
}

func NewKeyboardEvent(key int) (kbe KeyboardEvent) {
	kbe.EventPrototype = event.New(EV_KEYBOARD)
	kbe.Key = key
	return
}