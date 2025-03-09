package tyumi

import (
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

var running bool

// This is the gameloop
func Run() {
	defer log.WriteToDisk()

	if !isInitialized() {
		log.Error("Tyumi must shut down now. Bye Bye.")
		return
	}

	events = event.NewStream(250, handleEvent)
	events.Listen(EV_QUIT, EV_CHANGESTATE)

	for running = true; running; {
		beginFrame()
		eventGenerator() //take inputs from platform, convert to tyumi events as appropriate, and distribute
		update()         //step forward the gamestate
		updateUI()       //update changed UI elements
		render()         //composite frame together, post process, and render to screen
		events.Process() //processes internal events
		endFrame()       //do any end of tick cleanup, then sleep to maintain framerate if necessary
	}

	currentPlatform.Shutdown()
	log.Info("Tyumi says goodbye. ;)")
}

func beginFrame() {
	frameTime = time.Now()
}

// This is the generic tick function. Steps forward the gamestate, and performs some engine-specific per-tick functions.
func update() {
	mainConsole.events.Process()

	if !currentState.IsBlocked() {
		currentState.InputEvents().Process()
	} else {
		currentState.InputEvents().Flush()
	}

	if !currentState.IsBlocked() {
		currentState.Update()
	}

	currentState.Events().Process() //process any gameplay events from this frame.
}

// Updates any UI elements that need updating after the most recent tick in the current active state.
func updateUI() {
	currentState.UpdateUI()
	currentState.Window().Update()
}

// builds the frame and renders using the current platform's renderer.
func render() {
	currentState.Window().Render()
	currentState.Window().Draw(&mainConsole.Canvas, vec.ZERO_COORD, 0)
	if mainConsole.Dirty() {
		renderer.Render()
	}
}

func endFrame() {
	//framerate limiter, so the cpu doesn't implode
	time.Sleep(frameTargetDuration - time.Since(frameTime))
	tick++
}

// handles events from the engine's internal event stream. runs once per tick
func handleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_QUIT: //quit event, like from clicking the close window button on the window
		running = false
		currentState.Shutdown()
		event_handled = true
	case EV_CHANGESTATE:
		currentState.Shutdown()
		changeEvent := e.(*StateChangeEvent)
		currentState = changeEvent.newState
		event_handled = true
	}

	return
}
