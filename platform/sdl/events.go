package sdl

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
	"github.com/veandco/go-sdl2/sdl"
)

// sdl event processor
func (p *Platform) GenerateEvents() {
	//save mouse position so we can detect if we've moved to a new cell and fire a mouse move event
	new_mouse_pos := p.mouse_position

	// fetch the state of the mods on the keyboard (because windows eats certain modifiers sometimes??)
	// we did this because windows doesn't report the SHIFT key on keypad inputs when numlock is off. by fetching the
	// keyboard mod state early, now this problem only affects the key-release event. so still not pefect but it's
	// better than nothing.
	// INVESTIGATE: this seems like it'll have other consequences that I can't see yet. if keyboard inputs get weird
	// we can look into this further.
	sdlMods := sdl.GetModState()

eventLoop:
	for sdlevent := sdl.PollEvent(); sdlevent != nil; sdlevent = sdl.PollEvent() {
		switch e := sdlevent.(type) {
		case *sdl.QuitEvent:
			event.Fire(tyumi.EV_QUIT)
			break eventLoop //don't care about other input events if we're quitting
		case *sdl.WindowEvent:
			if e.Event == sdl.WINDOWEVENT_RESIZED {
				p.renderer.onWindowResize()
			}
		case *sdl.KeyboardEvent:
			var mods input.KeyModifiers
			if sdlMods&sdl.KMOD_CTRL != 0 {
				mods |= input.KEYMOD_CTRL
			}
			if sdlMods&sdl.KMOD_ALT != 0 {
				mods |= input.KEYMOD_ALT
			}
			if sdlMods&sdl.KMOD_SHIFT != 0 {
				mods |= input.KEYMOD_SHIFT
			}

			if key, ok := keycodemap[e.Keysym.Sym]; ok {
				switch e.State {
				case sdl.PRESSED:
					if e.Repeat == 0 {
						input.FireKeyPressEvent(key, mods)
					} else {
						input.FireKeyRepeatEvent(key, mods)
					}
				case sdl.RELEASED:
					input.FireKeyReleaseEvent(key, mods)
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
}
