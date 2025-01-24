package platform

import "github.com/bennicholls/tyumi/gfx"

// definition of whatever system is grabbing events from the system
// the expectation is that this function consumes system-level events and converts them into tyumi
// events and fires them
type EventGenerator func()

// definition of whatever system is rendering to the screen
type Renderer interface {
	Setup(console *gfx.Canvas, glyphPath, fontPath, title string) error
	Ready() bool
	Cleanup()
	ChangeFonts(glyphPath, fontPath string) error
	SetFullscreen(bool)
	ToggleFullscreen()
	Render()
	ForceRedraw()
	ToggleDebugMode(string)
}
