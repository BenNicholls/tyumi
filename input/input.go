package input

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/veandco/go-sdl2/sdl"
)

var EV_KEYBOARD = event.Register()
var EV_QUIT = event.Register()

//gather input events from sdl and fire relevant tyumi events
//TODO: abstract this so different input handlers can be used instead of just SDL
func Process() {
	for sdlevent := sdl.PollEvent(); sdlevent != nil; sdlevent = sdl.PollEvent() {
		switch e := sdlevent.(type) {
		case *sdl.QuitEvent:
			event.Fire(event.New(EV_QUIT))
			break //don't care about other input events if we're quitting
		case *sdl.KeyboardEvent:
			//only want keydown events for now.
			if e.Type == sdl.KEYDOWN {
				event.Fire(NewKeyboardEvent(Keycode(e.Keysym.Sym)))
			}
		}
	}

	return
}
