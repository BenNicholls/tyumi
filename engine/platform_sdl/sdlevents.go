package platform_sdl

import (
	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
	"github.com/veandco/go-sdl2/sdl"
)

// sdl event processor
func (sdlp *SDLPlatform) processEvents() {

	//save mouse position so we can detect if we've moved to a new cell and fire a mouse move event
	new_mouse_pos := sdlp.mouse_position

	for sdlevent := sdl.PollEvent(); sdlevent != nil; sdlevent = sdl.PollEvent() {
		switch e := sdlevent.(type) {
		case *sdl.QuitEvent:
			event.Fire(event.New(engine.EV_QUIT))
			break //don't care about other input events if we're quitting
		case *sdl.KeyboardEvent:
			if key, ok := keycodemap[e.Keysym.Sym]; ok {
				switch e.State {
				case sdl.PRESSED:
					if e.Repeat == 0 {
						input.FireKeyPressEvent(key)
					} else {
						input.FireKeyRepeatEvent(key)
					}
				case sdl.RELEASED:
					input.FireKeyReleaseEvent(key)
				}
			}
		case *sdl.MouseMotionEvent:
			new_mouse_pos = vec.Coord{int(e.X) / sdlp.renderer.tileSize, int(e.Y) / sdlp.renderer.tileSize}
		case *sdl.MouseButtonEvent:
			continue
		}
	}

	if new_mouse_pos != sdlp.mouse_position {
		input.FireMouseMoveEvent(new_mouse_pos, new_mouse_pos.Subtract(sdlp.mouse_position))
		sdlp.mouse_position = new_mouse_pos
	}

	return
}
