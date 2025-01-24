package platform_sdl

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/platform"
	"github.com/veandco/go-sdl2/sdl"
)

// sdl event processor
func (sdlp *SDLPlatform) processEvents() {
	for sdlevent := sdl.PollEvent(); sdlevent != nil; sdlevent = sdl.PollEvent() {
		switch e := sdlevent.(type) {

		case *sdl.QuitEvent:
			event.Fire(event.New(platform.EV_QUIT))
			break //don't care about other input events if we're quitting

		case *sdl.KeyboardEvent:
			//only want keydown events for now.
			if e.Type == sdl.KEYDOWN {
				if key, ok := keycodemap[e.Keysym.Sym]; ok {
					input.FireKeydownEvent(key)
				}
			}
		}
	}

	return
}
