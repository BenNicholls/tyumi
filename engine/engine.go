package engine

import (
	"errors"
	"runtime"
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/platform"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var renderer platform.Renderer
var eventGenerator platform.EventGenerator
var console Console
var mainState State

var tick int //count of number of ticks since engine was initialized
var frameTargetDur time.Duration // target duration of each frame, based on user-set framerate
var frameTime time.Time // time current frame began

var running bool

var events event.Stream //the main event stream for engine-level events

func init() {
	SetFramerate(60)
}

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
func SetupCustomRenderer(r platform.Renderer, glyphPath, fontPath, title string) error {
	renderer = r
	error := SetupRenderer(glyphPath, fontPath, title)
	return error
}

// Sets maximum framerate as enforced by the framerate limiter. NOTE: cannot go higher than 1000 fps.
func SetFramerate(f int) {
	f = util.Min(f, 1000)
	frameTargetDur = time.Duration(1000/float64(f+1)) * time.Millisecond
}

// This is the gameloop
func Run() {
	runtime.LockOSThread() //most of sdl is single threaded

	defer log.WriteToDisk()

	events = event.NewStream(250, handleEvent)
	events.Listen(platform.EV_QUIT)

	if mainState == nil {
		log.Error("No game state for Tyumi to run! Ensure that engine.InitMainState() is run before the gameloop.")
		return
	}

	var err error
	eventGenerator, err = platform.GetEventGenerator()
	if err != nil {
		log.Error("Could not get input processor from platform: ", err.Error())
		return
	}

	for running = true; running; {
		beginFrame()
		eventGenerator() //take inputs from platform, convert to tyumi events as appropriate, and distribute
		update()         //step forward the gamestate
		updateUI()       //update changed UI elements
		render()         //composite frame together, post process, and render to screen
		events.Process() //processes internal events
		endFrame()
	}

	renderer.Cleanup()
}

func beginFrame() {
	frameTime = time.Now()
}

// This is the generic tick function. Steps forward the gamestate, and performs some engine-specific per-tick functions.
func update() {
	mainState.InputEvents().Process()
	mainState.Update()
	mainState.Events().Process() //process any gameplay events from this frame.
}

// Updates any UI elements that need updating after the most recent tick in the current active state.
func updateUI() {
	mainState.UpdateUI()
	mainState.Window().Update()
}

// builds the frame and renders using whatever the current renderer is (sdl, web, terminal, whatever)
// this runs at speed determined by user-input FPS, defaulting to 60 FPS.
func render() {
	mainState.Window().Render()
	mainState.Window().RenderAnimations()
	mainState.Window().DrawToCanvas(&console.Canvas, vec.ZERO_COORD, 0)
	mainState.Window().FinalizeRender()
	if console.Dirty() {
		renderer.Render()
	}
}

func endFrame() {
	//framerate limiter, so the cpu doesn't implode
	time.Sleep(frameTargetDur - time.Since(frameTime))
	tick++
}

// handles events from the engine's internal event stream. runs once per tick()
func handleEvent(e event.Event) {
	switch e.ID() {
	case platform.EV_QUIT: //quit event, like from clicking the close window button on the window
		running = false
		mainState.Shutdown()
		e.SetHandled()
	}
}
