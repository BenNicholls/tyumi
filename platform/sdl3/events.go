package sdl3

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
	"github.com/jupiterrider/purego-sdl3/sdl"
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
	var sdlEvent sdl.Event

eventLoop:
	for ok := sdl.PollEvent(&sdlEvent); ok; ok = sdl.PollEvent(&sdlEvent) {
		switch sdlEvent.Type() {
		case sdl.EventQuit, sdl.EventWindowCloseRequested:
			event.Fire(tyumi.EV_QUIT)
			break eventLoop //don't care about other input events if we're quitting
		case sdl.EventWindowResized:
			p.renderer.onWindowResize()
		case sdl.EventKeyDown, sdl.EventKeyUp:
			var mods input.KeyModifiers
			if sdlMods&sdl.KeymodCtrl != 0 {
				mods |= input.KEYMOD_CTRL
			}
			if sdlMods&sdl.KeymodAlt != 0 {
				mods |= input.KEYMOD_ALT
			}
			if sdlMods&sdl.KeymodShift != 0 {
				mods |= input.KEYMOD_SHIFT
			}

			keyevent := sdlEvent.Key()
			if key, ok := keycodemap[keyevent.Key]; ok {
				if keyevent.Down {
					if keyevent.Repeat {
						input.FireKeyRepeatEvent(key, mods)
					} else {
						input.FireKeyPressEvent(key, mods)
					}
				} else {
					input.FireKeyReleaseEvent(key, mods)
				}
			}
		case sdl.EventMouseMotion:
			mouseEvent := sdlEvent.Motion()
			new_mouse_pos = vec.Coord{int(mouseEvent.X) / p.renderer.tileSize, int(mouseEvent.Y) / p.renderer.tileSize}
		case sdl.EventMouseButtonDown, sdl.EventMouseButtonUp:
			continue
		}
	}

	if new_mouse_pos != p.mouse_position {
		input.FireMouseMoveEvent(new_mouse_pos, new_mouse_pos.Subtract(p.mouse_position))
		p.mouse_position = new_mouse_pos
	}
}

