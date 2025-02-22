package tyumi

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

var mainConsole console

type console struct {
	gfx.Canvas
	ready bool

	mouseCursorEnabled bool
	mouseCursorVisuals gfx.Visuals
	mouseCursorPos     vec.Coord

	events event.Stream
}

func (c *console) handleEvents(e event.Event) (event_handled bool) {
	switch e.ID() {
	case input.EV_MOUSEMOVE:
		if c.mouseCursorEnabled {
			c.Clear(vec.Rect{c.mouseCursorPos, vec.Dims{1, 1}})
			c.mouseCursorPos = e.(*input.MouseMoveEvent).Position
			c.DrawVisuals(c.mouseCursorPos, 100000000, c.mouseCursorVisuals) //TODO: cursor should probably have a proper depth level just for it
		}
	}

	return
}

// Initializes the console. The console is where Tyumi composites together the frame along with the mouse cursor and
// any post-processing effects (coming soon TM) before being sent to the renderer to be displayed. The size here is in
// Cells, not pixels. The final window size will be the size here multiplied by however big your tiles are.
func InitConsole(console_size vec.Dims) {
	mainConsole.Init(console_size)
	mainConsole.ready = true

	mainConsole.events = event.NewStream(50, mainConsole.handleEvents)
	mainConsole.events.Listen(input.EV_MOUSEMOVE)

	mainConsole.mouseCursorVisuals = gfx.NewGlyphVisuals(gfx.GLYPH_BORDER_UUDDLLRR, col.Pair{col.WHITE, col.NONE})
}

// EnableCursor enables drawing of the mouse cursor. Also sets input.EnableMouse = true
func EnableCursor() {
	mainConsole.mouseCursorEnabled = true
	input.EnableMouse = true
}
