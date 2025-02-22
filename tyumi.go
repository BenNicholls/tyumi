package tyumi

import (
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/log"
)

var renderer Renderer
var eventGenerator EventGenerator
var events event.Stream //the main event stream for engine-level events

var tick int                          //count of number of ticks since engine was initialized
var frameTargetDuration time.Duration // target duration of each frame, based on user-set framerate
var frameTime time.Time

func init() {
	SetFramerate(60)
}

// Sets maximum framerate as enforced by the framerate limiter. NOTE: cannot go higher than 1000 fps.
func SetFramerate(f int) {
	f = min(f, 1000)
	frameTargetDuration = time.Duration(1000/float64(f+1)) * time.Millisecond
}

// Gets the tick number for the current tick (duh)
func GetTick() int {
	return tick
}

func isInitialized() bool {
	if !mainConsole.ready {
		log.Error("Cannot run Tyumi: console not initialized. Run InitConsole() first.")
		return false
	}

	if currentPlatform == nil {
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

	if currentState == nil {
		log.Error("Cannot run Tyumi, no MainState set! Run SetInitialMainState() first.")
		return false
	}

	return true
}
