package engine

import (
	"runtime"
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var renderer Renderer
var eventGenerator EventGenerator
var console Console
var mainState State

var tick int                     //count of number of ticks since engine was initialized
var frameTargetDur time.Duration // target duration of each frame, based on user-set framerate
var frameTime time.Time          // time current frame began

var running bool

var events event.Stream //the main event stream for engine-level events

func init() {
	SetFramerate(60)
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

	if !engineIsInitialized() {
		log.Error("Tyumi must shut down now. Bye Bye.")
		return
	}

	events = event.NewStream(250, handleEvent)
	events.Listen(EV_QUIT)

	for running = true; running; {
		beginFrame()
		eventGenerator() //take inputs from platform, convert to tyumi events as appropriate, and distribute
		update()         //step forward the gamestate
		updateUI()       //update changed UI elements
		render()         //composite frame together, post process, and render to screen
		events.Process() //processes internal events
		endFrame()
	}

	current_platform.Shutdown()
	log.Info("Tyumi says goodbye. ;)")
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

// builds the frame and renders using the current platform's renderer.
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

// handles events from the engine's internal event stream. runs once per tick
func handleEvent(e event.Event) {
	switch e.ID() {
	case EV_QUIT: //quit event, like from clicking the close window button on the window
		running = false
		mainState.Shutdown()
		e.SetHandled()
	}
}

func engineIsInitialized() bool {
	if !console.ready {
		log.Error("Cannot run Tyumi: console not initialized. Run InitConsole() first.")
		return false
	}

	if current_platform == nil {
		log.Error("Cannot run Tyumi: no platform set. Run SetPlatform() first.")
		return false
	}

	if !renderer.Ready() {
		log.Error("Cannot run Tyumi: renderer not set up. Run SetupRenderer() first.")
		return false
	}

	if eventGenerator == nil {
		log.Error("Cannot run Tyumi: platform did not provide an event generator.")
		return false
	}
	
	if mainState == nil {
		log.Error("Cannot run Tyumi, no MainState set! Run SetInitialMainState() first.")
		return false
	}

	return true
}