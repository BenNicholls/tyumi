package engine

import (
	"runtime"
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

var renderer Renderer
var event_generator EventGenerator
var main_console console
var main_state State

var tick int                       //count of number of ticks since engine was initialized
var frame_target_dur time.Duration // target duration of each frame, based on user-set framerate
var frame_time time.Time           // time current frame began

var running bool

var events event.Stream //the main event stream for engine-level events

func init() {
	SetFramerate(60)
}

// Sets maximum framerate as enforced by the framerate limiter. NOTE: cannot go higher than 1000 fps.
func SetFramerate(f int) {
	f = min(f, 1000)
	frame_target_dur = time.Duration(1000/float64(f+1)) * time.Millisecond
}

// Gets the tick number for the current tick (duh)
func GetTick() int {
	return tick
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
		event_generator() //take inputs from platform, convert to tyumi events as appropriate, and distribute
		update()          //step forward the gamestate
		updateUI()        //update changed UI elements
		render()          //composite frame together, post process, and render to screen
		events.Process()  //processes internal events
		endFrame()
	}

	current_platform.Shutdown()
	log.Info("Tyumi says goodbye. ;)")
}

func beginFrame() {
	frame_time = time.Now()
}

// This is the generic tick function. Steps forward the gamestate, and performs some engine-specific per-tick functions.
func update() {
	main_console.events.Process()

	if !main_state.IsBlocked() {
		main_state.InputEvents().Process()
	} else {
		main_state.InputEvents().Flush()
	}

	if !main_state.IsBlocked() {
		main_state.Update()
	}

	main_state.Events().Process() //process any gameplay events from this frame.
}

// Updates any UI elements that need updating after the most recent tick in the current active state.
func updateUI() {
	main_state.UpdateUI()
	main_state.Window().Update()
}

// builds the frame and renders using the current platform's renderer.
func render() {
	main_state.Window().Render()
	main_state.Window().Draw(&main_console.Canvas, vec.ZERO_COORD, 0)
	if main_console.Dirty() {
		renderer.Render()
	}
}

func endFrame() {
	//framerate limiter, so the cpu doesn't implode
	time.Sleep(frame_target_dur - time.Since(frame_time))
	tick++
}

// handles events from the engine's internal event stream. runs once per tick
func handleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_QUIT: //quit event, like from clicking the close window button on the window
		running = false
		main_state.Shutdown()
		event_handled = true
	}

	return
}

func engineIsInitialized() bool {
	if !main_console.ready {
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

	if event_generator == nil {
		log.Error("Cannot run Tyumi: platform did not provide an event generator.")
		return false
	}

	if main_state == nil {
		log.Error("Cannot run Tyumi, no MainState set! Run SetInitialMainState() first.")
		return false
	}

	return true
}
