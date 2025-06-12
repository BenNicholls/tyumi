package tyumi

import (
	"fmt"
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
	"github.com/pkg/profile"
)

var (
	running        bool
	renderer       Renderer
	eventGenerator EventGenerator
	events         event.Stream //the main event stream for engine-level events
)

// This is the gameloop
func Run() {
	if ProfilingEnabled {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	if !isInitialized() {
		log.Error("Tyumi must shut down now. Bye Bye.")
		return
	}

	events = event.NewStream(250, handleEvent)
	events.Listen(EV_QUIT, EV_CHANGESCENE)
	if Debug {
		events.Listen(input.EV_KEYBOARD)
	}

	fpsTime = time.Now()

	for running = true; running; {
		beginFrame()
		eventGenerator()       //take inputs from platform, convert to tyumi events as appropriate, and distribute
		update()               //step forward the gamestate
		updateUI()             //update changed UI elements
		render()               //composite frame together, post process, and render to screen
		events.ProcessEvents() //processes internal events
		endFrame()             //do any end of tick cleanup, then sleep to maintain framerate if necessary
	}

	currentPlatform.Shutdown()
	log.Info("Tyumi says goodbye. ;)")
}

func beginFrame() {
	frameTime = time.Now()
	activeScene = currentScene
	if activeSubScene := currentScene.getActiveSubScene(); activeSubScene != nil {
		activeScene = activeSubScene
	}
}

// This is the generic tick function. Steps forward the gamestate, and performs some engine-specific per-tick functions.
func update() {
	mainConsole.ProcessEvents()

	if !activeScene.IsBlocked() {
		activeScene.InputEvents().ProcessEvents()
	}

	// make sure we don't accumulate a bunch of inputs in scenes that aren't being updated for whatever reason
	currentScene.flushInputs()

	if !activeScene.IsBlocked() {
		activeScene.processTimers()
		activeScene.Update()
	}

	activeScene.ProcessEvents() //process any gameplay events from this frame.
}

// Updates any UI elements that need updating after the most recent tick in the current active scene.
func updateUI() {
	activeScene.UpdateUI()
	activeScene.Window().Update()
}

// builds the frame and renders using the current platform's renderer.
func render() {
	activeScene.Window().Render()
	activeScene.Window().Draw(&mainConsole.Canvas)

	if ShowFPS {
		if time.Since(fpsTime) > time.Second {
			mainConsole.DrawText(vec.ZERO_COORD, 10000000,
				fmt.Sprintf("FPS: %4d", tick-fpsTicks), col.Pair{col.ORANGE, col.MAROON},
				gfx.DRAW_TEXT_LEFT)
			fpsTicks = tick
			fpsTime = time.Now()
		}
	}

	renderer.Render()
}

func endFrame() {
	//framerate limiter, so the cpu doesn't implode
	if !overclock {
		time.Sleep(frameTargetDuration - time.Since(frameTime))
	}

	tick++
}

// handles events from the engine's internal event stream. runs once per tick
func handleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_QUIT: //quit event, like from clicking the close window button on the window
		currentScene.Shutdown()
		currentScene.cleanup()
		running = false
		event_handled = true
	case EV_CHANGESCENE:
		currentScene.Shutdown()
		currentScene.cleanup()
		currentScene = e.(*SceneChangeEvent).newScene
		event_handled = true
	}

	if Debug {
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
				log.Info("Dumping UI of current scene! Saving files to directory 'uidump'")
				currentScene.Window().DumpUI("uidump", true)
				if activeScene != currentScene {
					activeScene.Window().DumpUI("uidump", false)
				}
			}
		}
	}

	return
}
