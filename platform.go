package tyumi

import (
	"errors"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/log"
)

var EV_QUIT = event.Register("Quit Event")

// Platform defines the API for the platform-specific code that Tyumi uses to interface with the system. It's split
// into a number of subsystems, all of which need to be handled by the platform at the moment.
// THINK: we could support having platforms that omit certain subsystems by having them be able to report exactly
// what features they support. This is a problem for later.
type Platform interface {
	Init() error
	GetRenderer() Renderer
	GetEventGenerator() EventGenerator
	Shutdown()
}

var currentPlatform Platform = nil

// Sets the platform to be used by Tyumi for rendering, gathering of system events, and more. This must be called
// after console initialization and before running the game loop. The engine will Init() the platform for you.
func SetPlatform(p Platform) (err error) {
	err = p.Init()
	if err != nil {
		log.Error("Could not initialize platform: ", err)
		return
	}

	if currentPlatform != nil {
		log.Info("Shutting down old platform.")
		currentPlatform.Shutdown()
	}

	currentPlatform = p
	renderer = p.GetRenderer()
	eventGenerator = p.GetEventGenerator()

	return
}

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

// Sets up the renderer. This must be done after initializaing the console and setting the platform, but before
// running the main game loop.
func SetupRenderer(glyph_path, font_path, title string) error {
	if !mainConsole.ready {
		log.Error("Cannot initialize renderer: console not initialized. Run InitConsole() first.")
		return errors.New("NO CONSOLE.")
	}

	if currentPlatform == nil {
		log.Error("Cannot initialize renderer: no platform set. Run SetPlatform() first.")
		return errors.New("NO PLATFORM.")
	}

	err := renderer.Setup(&mainConsole.Canvas, glyph_path, font_path, title)
	if err != nil {
		return err
	}

	return nil
}
