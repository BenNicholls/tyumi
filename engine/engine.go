package engine

import (
	"errors"
	"runtime"
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/log"
	"github.com/veandco/go-sdl2/sdl"
)

var renderer gfx.Renderer
var console Console
var mainState State

var tick int //count of number of ticks since engine was initialized
var running bool

var events *event.Stream //the main event stream for the engine. all events will go and be distributed from here
var EV_QUIT = event.Register()

//Initializes the renderer. This must be done after initializaing the console, but before running the main game loop.
//logs and returns an error if this was unsuccessful.
func InitRenderer(r gfx.Renderer, glyphPath, fontPath, title string) error {
	if !console.ready {
		log.Error("Cannot initialize renderer: console not initialized. Run InitConsole() first.")
		return errors.New("NO CONSOLE.")
	}

	err := r.Setup(&console.Canvas, glyphPath, fontPath, title)
	if err != nil {
		return err
	}
	renderer = r

	return nil
}

//This is the gameloop
func Run() {
	runtime.LockOSThread() //most of sdl is single threaded

	defer log.WriteToDisk()

	events = event.NewStream(250)

	if mainState == nil {
		log.Error("No game state for Tyumi to run! Ensure that engine.InitMainState() is run before the gameloop.")
		return
	}

	for running = true; running; {
		processInput() //take inputs from sdl, convert to tyumi events as appropriate, and distribute
		update()       //step forward the gamestate
		updateUI()     //update changed UI elements
		render()       //composite frame together, post process, and render to screen
		handleEvents() //handle events generated this frame, both internal and external
		time.Sleep(15) //TODO: implement an actual framerate limiter
	}

	renderer.Cleanup()
}

//gather input events from sdl and handle/distribute accordingly
func processInput() int {
	for sdlevent := sdl.PollEvent(); sdlevent != nil; sdlevent = sdl.PollEvent() {
		switch e := sdlevent.(type) {
		case *sdl.QuitEvent:
			events.Add(event.New(EV_QUIT))
			break //don't care about other input events if we're quitting
		case *sdl.KeyboardEvent:
			mainState.HandleEvent(NewKeyboardEvent(int(e.Keysym.Sym)))
		}
	}

	return 0
}

//This is the generic tick function. Steps forward the gamestate, and performs some engine-specific per-tick functions.
func update() {
	mainState.Update()
	tick++
}

//This function updates any UI elements that need updating after the most recent tick in the current active state.
func updateUI() {
	mainState.UpdateUI()
	mainState.Window().UpdateChildren()
	//    - tick animations
}

//builds the frame and renders using whatever the current renderer is (sdl, web, terminal, whatever)
//this runs at speed determined by user-input FPS, defaulting to 60 FPS. this also updates any current animations in the active
//state's ui tree
func render() {
	mainState.Window().Render()
	mainState.Window().DrawToCanvas(&console.Canvas, 0, 0, 1)
	//  - render animations
	renderer.Render()
}

//handles events from the main event stream. passes them to the main gamestate first, then handles any required internal
//engine behaviour
func handleEvents() {
	for e := events.Next(); e != nil; e = events.Next() {
		mainState.HandleEvent(e)

		switch e.ID() {
		case EV_QUIT:
			running = false
			mainState.Shutdown()
		}
	}
}
