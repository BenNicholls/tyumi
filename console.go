package tyumi

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

var mainConsole console

type console struct {
	ui.Window

	ready bool
	title string // title of the program

	mouseCursorEnabled bool
	mouseCursorVisuals gfx.Visuals
	mouseCursorPos     vec.Coord
}

func (c *console) changeTitle(new_title string) {
	if c.title == new_title {
		return
	}

	c.title = new_title
	currentPlatform.ChangeTitle(c.title)
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
// This must be done AFTER setting the platform.
func InitConsole(title string, console_size vec.Dims, glyph_path, font_path string) {
	if currentPlatform == nil {
		log.Error("Could not initialize console: platform not set. Run SetPlatform() first!")
		return
	}

	mainConsole.Init(console_size, vec.ZERO_COORD, 0)
	mainConsole.title = title

	// now that the console is set up, we can initialize the renderer (hopefully)
	err := renderer.Setup(&mainConsole.Canvas, glyph_path, font_path, title)
	if err != nil {
		log.Error("Renderer setup failed! Console no good.")
		return
	}

	mainConsole.SetEventHandler(mainConsole.handleEvents)
	mainConsole.Listen(input.EV_MOUSEMOVE)

	mainConsole.mouseCursorVisuals = gfx.NewGlyphVisuals(gfx.GLYPH_BORDER_UUDDLLRR, col.Pair{col.WHITE, col.NONE})

	mainConsole.ready = true

	if Debug {
		debugger.Init()
	}
}

// ChangeTitle changes the title of the running program. i.e. the string shown in the title bar of the program's
// window for a windows/mac/linux program, or the string in the tab of a running web app.
func ChangeTitle(title string) {
	mainConsole.changeTitle(title)
}

// EnableCursor enables drawing of the mouse cursor. Also sets input.EnableMouse = true so mouse events start being
// fired.
func EnableCursor() {
	mainConsole.mouseCursorEnabled = true
	input.EnableMouse = true
}
