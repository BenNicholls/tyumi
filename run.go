package tyumi

import (
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
)

var running bool

// This is the gameloop
func Run() {
	if !isInitialized() {
		log.Error("Tyumi must shut down now. Bye Bye.")
		return
	}

	events = event.NewStream(250, handleEvent)
	events.Listen(EV_QUIT, EV_CHANGESTATE)
	if debug {
		events.Listen(input.EV_KEYBOARD)
	}

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
	activeState = currentState
	if activeSubState := currentState.getActiveSubState(); activeSubState != nil {
		activeState = activeSubState
	}
}

// This is the generic tick function. Steps forward the gamestate, and performs some engine-specific per-tick functions.
func update() {
	mainConsole.events.Process()

	if !activeState.IsBlocked() {
		activeState.InputEvents().Process()
	}

	// make sure we don't accumulate a bunch of inputs in states that aren't being updated for whatever reason
	currentState.flushInputs()

	if !activeState.IsBlocked() {
		activeState.processTimers()
		activeState.Update()
	}

	activeState.Events().Process() //process any gameplay events from this frame.
}

// Updates any UI elements that need updating after the most recent tick in the current active state.
func updateUI() {
	activeState.UpdateUI()
	activeState.Window().Update()
}

// builds the frame and renders using the current platform's renderer.
func render() {
	activeState.Window().Render()
	activeState.Window().Draw(&mainConsole.Canvas)
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
		currentState.Shutdown()
		currentState.cleanup()
		running = false
		event_handled = true
	case EV_CHANGESTATE:
		currentState.Shutdown()
		currentState.cleanup()
		currentState = e.(*StateChangeEvent).newState
		event_handled = true
	}

	if debug {
		if e.ID() == input.EV_KEYBOARD {
			key_event := e.(*input.KeyboardEvent)
			if key_event.PressType == input.KEY_RELEASED {
				return
			}
			switch key_event.Key {
			case input.K_F9:
				log.Info("Taking Screenshot! Saving to 'screenshot.xp'.")
				mainConsole.ExportToXP("screenshot.xp")
			case input.K_F10:
				log.Info("Dumping UI of current state! Saving files to directory 'uidump'")
				currentState.Window().DumpUI("uidump", true)
				if activeState != currentState {
					activeState.Window().DumpUI("uidump", false)
				}
			}
		}
	}

	return
}
