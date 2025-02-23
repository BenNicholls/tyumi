package sdl

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
	"github.com/veandco/go-sdl2/sdl"
)

// sdl event processor
func (p *Platform) processEvents() {
	//save mouse position so we can detect if we've moved to a new cell and fire a mouse move event
	new_mouse_pos := p.mouse_position

	for sdlevent := sdl.PollEvent(); sdlevent != nil; sdlevent = sdl.PollEvent() {
		switch e := sdlevent.(type) {
		case *sdl.QuitEvent:
			event.Fire(event.New(tyumi.EV_QUIT))
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
			new_mouse_pos = vec.Coord{int(e.X) / p.renderer.tileSize, int(e.Y) / p.renderer.tileSize}
		case *sdl.MouseButtonEvent:
			continue
		}
	}

	if new_mouse_pos != p.mouse_position {
		input.FireMouseMoveEvent(new_mouse_pos, new_mouse_pos.Subtract(p.mouse_position))
		p.mouse_position = new_mouse_pos
	}

	return
}
