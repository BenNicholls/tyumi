package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

// Anything that can do animations on a Canvas
type Animator interface {
	vec.Bounded

	Start()
	Play()
	Pause()
	Stop()

	Update()
	Render(*Canvas)

	IsPlaying() bool
	IsDone() bool
	IsOneShot() bool
}

// Base struct for animations. Embed this to satisfy Animator interface above.
type Animation struct {
	OneShot  bool //indicates animation should play once and then be deleted
	Repeat   bool //animation repeats when finished
	Area     vec.Rect
	Depth    int //depth value of the animation
	Duration int //duration of animation in ticks

	ticks   int  //incremented each update
	enabled bool //animation is playing
	dirty   bool //animation needs to be re-rendered
	reset   bool //indicates animation should reset and start over.
}

func (a *Animation) Update() {
	if a.reset {
		a.ticks = 0
		a.dirty = true
		a.reset = false
	} else {
		if a.Repeat {
			a.ticks = (a.ticks + 1) % a.Duration
		} else {
			a.ticks += 1
		}
	}

	if a.IsDone() {
		a.enabled = false
		a.reset = true
	}
}

func (a Animation) IsDone() bool {
	return !a.Repeat && a.ticks >= a.Duration
}

func (a Animation) IsOneShot() bool {
	return a.OneShot
}

func (a *Animation) MoveTo(pos vec.Coord) {
	a.Area.MoveTo(pos.X, pos.Y)
	a.dirty = true
}

func (a Animation) IsPlaying() bool {
	return a.enabled
}

// Starts an animation. If the animation is playing or paused, restarts it.
func (a *Animation) Start() {
	a.reset = true
	a.Play()
}

// Plays an animation. If the animation is paused, continues it.
func (a *Animation) Play() {
	a.enabled = true
}

// Pauses a playing animation.
func (a *Animation) Pause() {
	a.enabled = false
}

// Stops an animation and resets it.
func (a *Animation) Stop() {
	if !a.enabled {
		return
	}

	a.enabled = false
	a.reset = true
}

func (a Animation) Bounds() vec.Rect {
	return a.Area
}

}

}
