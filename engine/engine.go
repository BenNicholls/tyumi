package engine

import (
	"errors"
	"runtime"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/platform"
	"github.com/bennicholls/tyumi/vec"
)

var renderer gfx.Renderer
var inputProcessor input.Processor
var console Console
var mainState State

var tick int //count of number of ticks since engine was initialized
var running bool

var events event.Stream //the main event stream for the engine. all events will go and be distributed from here

// Sets up the renderer. This must be done after initializaing the console, but before running the main game loop.
// logs and returns an error if this was unsuccessful.
func SetupRenderer(glyphPath, fontPath, title string) error {
	//if no renderer has been initialized, get one from the platform package.
	if renderer == nil {
		r, err := platform.GetNewRenderer()
		if err != nil {
			return err
		}
		renderer = r
	}

	if !console.ready {
		log.Error("Cannot initialize renderer: console not initialized. Run InitConsole() first.")
		return errors.New("NO CONSOLE.")
	}

	err := renderer.Setup(&console.Canvas, glyphPath, fontPath, title)
	if err != nil {
		return err
	}

	return nil
}

// Sets up a custom user-defined renderer. This must be done after initializaing the console, but before running
// the main game loop.
func SetupCustomRenderer(r gfx.Renderer, glyphPath, fontPath, title string) error {
	renderer = r
	error := SetupRenderer(glyphPath, fontPath, title)
	return error
}

// This is the gameloop
func Run() {
	runtime.LockOSThread() //most of sdl is single threaded

	defer log.WriteToDisk()

	events = event.NewStream(250)
	events.AddHandler(handleEvent)
	events.Listen(input.EV_QUIT)
	events.Listen(input.EV_KEYBOARD)

	if mainState == nil {
		log.Error("No game state for Tyumi to run! Ensure that engine.InitMainState() is run before the gameloop.")
		return
	}

	var err error
	inputProcessor, err = platform.GetInputProcessor()
	if err != nil {
		log.Error("Could not get input processor from platform: ", err.Error())
		return
	}

	for running = true; running; {
		inputProcessor() //take inputs from platform, convert to tyumi events as appropriate, and distribute
		update()         //step forward the gamestate
		updateUI()       //update changed UI elements
		render()         //composite frame together, post process, and render to screen
		events.Process() //processes internal events
	}

	renderer.Cleanup()
}

// This is the generic tick function. Steps forward the gamestate, and performs some engine-specific per-tick functions.
func update() {
	mainState.InputEvents().Process()
	mainState.Update()
	mainState.Events().Process() //process any gameplay events from this frame.
	tick++
}

// This function updates any UI elements that need updating after the most recent tick in the current active state.
func updateUI() {
	mainState.UpdateUI()
	mainState.Window().Update()
	//    - tick animations
}

// builds the frame and renders using whatever the current renderer is (sdl, web, terminal, whatever)
// this runs at speed determined by user-input FPS, defaulting to 60 FPS. this also updates any current animations in the active
// state's ui tree
func render() {
	mainState.Window().Render()
	mainState.Window().RenderAnimations()
	mainState.Window().DrawToCanvas(&console.Canvas, vec.ZERO_COORD, 0)
	mainState.Window().FinalizeRender()
	renderer.Render()
}

// handles events from the engine's internal event stream. runs once per tick()
func handleEvent(e event.Event) {
	switch e.ID() {
	case input.EV_QUIT: //quit input event, like from clicking the close window button on the window
		running = false
		mainState.Shutdown()
	case input.EV_KEYBOARD:
		//pass keyboard events to the main ui to be processed by any relevant input handlers
		mainState.Window().HandleKeypress(e.(input.KeyboardEvent))
	}
}
