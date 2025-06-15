package gfx

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/vec"
)

var EV_ANIMATION_COMPLETE = event.Register("Animation Complete")

type AnimationEvent struct {
	event.EventPrototype

	Label string // label of the animation that produced the event
}

func fireAnimationCompleteEvent(label string) {
	event.Fire(EV_ANIMATION_COMPLETE, &AnimationEvent{Label: label})
}

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
	IsBlocking() bool
	IsUpdated() bool

	JustStopped() bool
	Finish()

	GetDuration() int
}

// Base struct for animations. Embed this to satisfy Animator interface above.
type Animation struct {
	OneShot       bool //indicates animation should play once and then be deleted
	Repeat        bool //animation repeats when finished
	Backwards     bool //play the animation backwards. NOTE: not all animations implement this (sometimes it doesn't make sense)
	Blocking      bool //whether this animation should block updates until completed. NOTE: if this is true, Repeat will be set to false to prevent infinite blocking
	AlwaysUpdates bool //if true, indicates this animation updates every frame
	Updated       bool //indicates to whatever is drawing the animation that it's going to render this frame
	Depth         int  //depth value of the animation
	Duration      int  //duration of animation in ticks
	Label         string

	OnDone func() // Callback run when animation finishes.

	area        vec.Rect
	ticks       int  //incremented each update
	enabled     bool //animation is playing
	reset       bool //indicates animation should reset and start over.
	justStopped bool //indicates animation has stopped recently. use Finish() to clear this flag.
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
		if a.Label != "" {
			fireAnimationCompleteEvent(a.Label)
		}

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

func (a *Animation) SetArea(area vec.Rect) {
	a.MoveTo(area.Coord)
	a.Resize(area.Dims)
}

func (a *Animation) Resize(size vec.Dims) {
	if a.area.Dims == size {
		return
	}

	a.area.Dims = size
	a.Updated = true
}

func (a *Animation) MoveTo(pos vec.Coord) {
	if a.area.Coord == pos {
		return
	}

	a.area.Coord = pos
	a.Updated = true
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

func (a Animation) Bounds() vec.Rect {
	return a.area
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

// Returns true if the animation has stopped recently.
func (a Animation) JustStopped() bool {
	return a.justStopped
}

// Finish clears the justStopped flag.
func (a *Animation) Finish() {
	a.justStopped = false
}

// AnimationChain is a container for multiple animations. Playing the chain will play all of the
// contained animations one after the other until all sub-animations have completed.
type AnimationChain struct {
	Animation

	animations []Animator
	current    int
}

// Adds animations to the chain. These animations will be played in the order provided.
func (ac *AnimationChain) Add(animations ...Animator) {
	for _, animation := range animations {
		ac.animations = append(ac.animations, animation)
		ac.Duration += animation.GetDuration()
	}
}

func (ac *AnimationChain) Update() {
	ac.Animation.Update()

	if ac.IsDone() {
		return
	}

	//ensure all animations in the chain are reset when the chain is reset
	if ac.ticks == 0 {
		ac.current = 0
		for _, anim := range ac.animations {
			anim.Start()
		}
	} else if ac.animations[ac.current].IsDone() {
		ac.current += 1
	}

	ac.animations[ac.current].Update()
	ac.Updated = ac.animations[ac.current].IsUpdated()
}

func (ac *AnimationChain) Render(canvas *Canvas) {
	ac.animations[ac.current].Render(canvas)
	ac.Updated = false
}
