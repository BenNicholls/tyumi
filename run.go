package tyumi

import (
	"fmt"
	"slices"
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
	"github.com/pkg/profile"
)

var (
	running  bool
	renderer Renderer
	events   event.Stream //the main event stream for engine-level events
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

	if ShowFPS {
		fpsLabelUpdateTime = time.Now()
	}

	currentFrameTime = time.Now() // so we get a non-nonsensical frame delta for the first frame.

	for running = true; running; {
		beginFrame()
		currentPlatform.GenerateEvents() //take inputs from platform, convert to tyumi events as appropriate, and distribute
		update()                         //step forward the gamestate
		updateUI()                       //update changed UI elements
		render()                         //composite frame together, post process, and render to screen
		events.ProcessEvents()           //processes internal events
		endFrame()                       //do any end of tick cleanup, then sleep to maintain framerate if necessary
	}

	currentPlatform.Shutdown()
	log.Info("Tyumi says goodbye. ;)")
}

func beginFrame() {
	prevFrameTime = currentFrameTime
	currentFrameTime = time.Now()
}

// This is the generic tick function. Steps forward the gamestate, and performs some engine-specific per-tick functions.
func update() {
	mainConsole.ProcessEvents()

	updateScene(currentScene)

	for _, d := range slices.Backward(dialogs) {
		if !d.IsDone() {
			updateScene(d)
		} else {
			closeDialog(d)
		}
	}
}

func updateScene(s scene) {
	if !s.IsBlocked() {
		s.InputEvents().ProcessEvents()
		s.processTimers()
		s.Update(GetFrameDelta())
	} else {
		s.flushInputs()
	}

	s.ProcessEvents() //process any gameplay events from this frame.
}

// Updates any UI elements that need updating after the most recent tick in the current active scene.
func updateUI() {
	currentScene.Window().Update(GetFrameDelta())

	for _, d := range dialogs {
		d.Window().Update(GetFrameDelta())
	}
}

// builds the frame and renders using the current platform's renderer.
func render() {
	mainConsole.Render()

	if ShowFPS {
		if time.Since(fpsLabelUpdateTime) > time.Second {
			fpsLabel := fmt.Sprintf("FPS: %4d", tick-fpsTicks)
			if !overclock {
				fpsLabel += fmt.Sprintf(" (%4.1f%%)", 100*(1-sleepTime.Seconds()))
				sleepTime = 0
			}

			mainConsole.DrawText(vec.ZERO_COORD, 10000000,
				fpsLabel, col.Pair{col.ORANGE, col.MAROON},
				gfx.DRAW_TEXT_LEFT)
			fpsTicks = tick
			fpsLabelUpdateTime = time.Now()
		}
	}

	renderer.Render()
}

func endFrame() {
	//framerate limiter, so the cpu doesn't implode
	if !overclock {
		sleep := frameTargetDuration - time.Since(currentFrameTime)
		sleepTime += sleep
		time.Sleep(sleep)
	}

	ecs.ProcessQueuedEntities()

	tick++
}

// handles events from the engine's internal event stream. runs once per tick
func handleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_QUIT: //quit event, like from clicking the close window button on the window
		closeAllDialogs()
		currentScene.Shutdown()
		currentScene.cleanup()
		running = false
		event_handled = true
	case EV_CHANGESCENE:
		changeScene(e.(*SceneChangeEvent).newScene)
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
				mainConsole.DumpUI("uidump", true)
			case input.K_F12:
				if Debug {
					if !debugger.IsOpen() {
						OpenDialog(debugger)
					}
				}
			}
		}
	}

	return
}
