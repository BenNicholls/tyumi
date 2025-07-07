package anim

import (
	"github.com/bennicholls/tyumi/log"
)

// Anything that can do animations on a Canvas
type Animator interface {
	Start()
	Play()
	Pause()
	Stop()

	Update()

	IsPlaying() bool
	IsDone() bool
	IsOneShot() bool
	IsBlocking() bool
	IsUpdated() bool

	JustStopped() bool
	Finish()

	GetDuration() int

	SetOneShot(bool)
}

// Base struct for animations. Embed this to satisfy Animator interface above.
type Animation struct {
	OneShot       bool //indicates animation should play once and then be deleted
	Repeat        bool //animation repeats when finished
	Backwards     bool //play the animation backwards. NOTE: not all animations implement this (sometimes it doesn't make sense)
	Blocking      bool //whether this animation should block updates until completed. NOTE: if this is true, Repeat will be set to false to prevent infinite blocking
	Updated       bool //indicates to whatever is drawing the animation that it's going to render this frame
	AlwaysUpdates bool //if true, indicates this animation updates every frame
	Duration      int  //duration of animation in ticks

	OnDone func() // Callback run when animation finishes.

	enabled     bool //animation is playing
	reset       bool //indicates animation should reset and start over.
	justStopped bool //indicates animation has stopped recently. use Finish() to clear this flag.
	ticks       int  //incremented each update
}

func (a *Animation) Update() {
	if a.reset {
		a.ticks = 0
		a.Updated = true
		a.reset = false
	} else {
		if a.Repeat && a.Blocking { // make sure we don't get in an infinite blocking loop
			a.Repeat = false
		}

		if a.Repeat {
			a.ticks = (a.ticks + 1) % a.Duration
		} else {
			a.ticks += 1
		}
	}

	if a.IsDone() {
		if a.OnDone != nil {
			a.OnDone()
		}

		a.enabled = false
		a.justStopped = true
		a.reset = true
	}
}

func (a Animation) IsDone() bool {
	return !a.Repeat && a.ticks >= a.Duration
}

func (a Animation) IsOneShot() bool {
	return a.OneShot
}

func (a Animation) IsPlaying() bool {
	return a.enabled
}

func (a Animation) IsBlocking() bool {
	return a.Blocking
}

func (a Animation) IsUpdated() bool {
	return a.AlwaysUpdates || a.Updated
}

// Sets the animation to OneShot. OneShot animations play once and then can be removed/deleted/garbaged/whatever.
// If the animation is set to Repeat, repeat is removed since you can't do both.
func (a *Animation) SetOneShot(oneshot bool) {
	if a.OneShot == oneshot {
		return
	}

	a.OneShot = oneshot
	if a.OneShot && a.Repeat {
		log.Warning("Repeating animations cannot be oneshot! Removing repeat flag.")
		a.Repeat = false
	}
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
	a.justStopped = true
}

// Stops an animation and resets it.
func (a *Animation) Stop() {
	if !a.enabled {
		return
	}

	a.enabled = false
	a.justStopped = true
	a.reset = true
}

func (a Animation) GetDuration() int {
	return a.Duration
}

// gets the tick number. if the animation is being played backwards, this will count down instead of up!
func (a Animation) GetTicks() int {
	if a.Backwards {
		return a.Duration - a.ticks - 1
	}

	return a.ticks
}

// GetProgress returns a value from [0,1] indicating the progress of the animation
func (a Animation) GetProgress() float64 {
	return float64(a.ticks) / float64(a.Duration)
}

func (a Animation) IsResetting() bool {
	return a.reset
}

// Returns true if the animation has stopped recently.
func (a Animation) JustStopped() bool {
	return a.justStopped
}

// Finish clears the justStopped flag.
func (a *Animation) Finish() {
	a.justStopped = false
}
