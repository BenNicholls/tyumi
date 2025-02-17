package engine

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

// The console is where the UI of the game state is composited together before being sent to the renderer. It also
// defines the size of the window you're using, and must be initialized with InitConsole() before running the
// gameloop.
type console struct {
	gfx.Canvas
	ready bool

	mouseCursorVisuals gfx.Visuals
	mouseCursorPos     vec.Coord

	events event.Stream
}

func (c *console) handleEvents(e event.Event) (event_handled bool) {
	switch e.ID() {
	case input.EV_MOUSEMOVE:
			c.Clear(vec.Rect{c.mouseCursorPos, vec.Dims{1, 1}})
			c.mouseCursorPos = e.(*input.MouseMoveEvent).Position
			c.DrawVisuals(c.mouseCursorPos, 100000000, c.mouseCursorVisuals) //TODO: cursor should probably have a proper depth level just for it
	}

	return
}

func InitConsole(w, h int) {
	main_console.Init(w, h)
	main_console.ready = true

	main_console.events = event.NewStream(50, main_console.handleEvents)
	main_console.events.Listen(input.EV_MOUSEMOVE)

	main_console.mouseCursorVisuals = gfx.NewGlyphVisuals(gfx.GLYPH_BORDER_UUDDLLRR, col.Pair{col.WHITE, col.NONE})
}
