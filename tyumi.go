package tyumi

import (
	"time"

	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
)

// User controllable flags
var (
	ProfilingEnabled bool // Enables CPU profiling. Only works in debug mode.
	ShowFPS          bool
)

// ticks and fps-control vars
var (
	tick                int           //count of number of ticks since engine was initialized
	frameTargetDuration time.Duration // target duration of each frame, based on user-set framerate
	prevFrameTime       time.Time     // time we started processing the previous frame. used to calculate frame deltas.
	currentFrameTime    time.Time     // time we started processing the current frame

	overclock          bool      // if true, no framerate limiting is enforced
	fpsTicks           int       // number of ticks when fps labal was last updated
	fpsLabelUpdateTime time.Time // time that the fps label most recently updated
)

func init() {
	SetFramerate(60)
}

// Sets maximum framerate as enforced by the framerate limiter. NOTE: cannot go higher than 1000 fps.
func SetFramerate(f int) {
	if f == 0 {
		overclock = true
		return
	}
	f = util.Clamp(f, 1, 1000)
	frameTargetDuration = time.Duration(1000/float64(f)) * time.Millisecond
}

func SetFullScreen(enable bool) {
	currentPlatform.GetRenderer().SetFullscreen(enable)
}

func SetClearColour(colour col.Colour) {
	currentPlatform.GetRenderer().SetClearColour(colour)
}

// Gets the tick number for the current tick (duh)
func GetTick() int {
	return tick
}

func GetFrameDelta() time.Duration {
	return currentFrameTime.Sub(prevFrameTime)
}

func isInitialized() bool {
	if currentPlatform == nil {
		log.Error("Cannot run Tyumi: no platform set. Run SetPlatform() first.")
		return false
	}

	if !mainConsole.ready {
		log.Error("Cannot run Tyumi: console not initialized. Run InitConsole() first.")
		return false
	}

	if !renderer.Ready() {
		log.Error("Cannot run Tyumi: renderer was not set up correctly.")
		return false
	}

	if currentScene == nil {
		log.Error("Cannot run Tyumi, no initial scene set! Run SetInitialScene() first.")
		return false
	}

	return true
}
