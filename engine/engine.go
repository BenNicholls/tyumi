package engine

import (
	"errors"
	"runtime"
	"time"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/log"
	"github.com/veandco/go-sdl2/sdl"
)

var renderer gfx.Renderer
var console Console
var mainState State

var tick int //count of number of ticks since engine was initialized

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

	if mainState == nil {
		log.Error("No game state for Tyumi to run! Ensure that engine.InitMainState() is run before the gameloop.")
		return
	}

	for running := true; running; {
		if processInput() == 1 { //gather inputs and handles/distributes them
			running = false
		}
		update()       //step forward the gamestate
		updateUI()     //update changed UI elements
		render()       //composite frame together, post process, and render to screen
		time.Sleep(15) //TODO: implement an actual framerate limiter
	}

	renderer.Cleanup()
}

//gather input events from sdl and handle/distribute accordingly
func processInput() int {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return 1 //TODO: make this return an actual event
			//break
		case *sdl.KeyboardEvent:
			//keyboard handling
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
