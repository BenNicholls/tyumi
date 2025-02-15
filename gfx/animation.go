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
	animEvent := AnimationEvent{
		EventPrototype: *event.New(EV_ANIMATION_COMPLETE),
		Label:          label,
	}
	event.Fire(&animEvent)
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

	GetDuration() int
}

// Base struct for animations. Embed this to satisfy Animator interface above.
type Animation struct {
	OneShot  bool //indicates animation should play once and then be deleted
	Repeat   bool //animation repeats when finished
	Area     vec.Rect
	Depth    int //depth value of the animation
	Duration int //duration of animation in ticks
	Label    string

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
		if a.Label != "" {
			fireAnimationCompleteEvent(a.Label)
		}
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

func (a Animation) IsPlaying() bool {
	return a.enabled
}

func (a *Animation) MoveTo(pos vec.Coord) {
	a.Area.MoveTo(pos.X, pos.Y)
	a.dirty = true
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

func (a Animation) GetDuration() int {
	return a.Duration
}

// AnimationChain is a container for multiple animations. Playing the chain will play all of the
// contained animations one after the other until all sub-animations have completed.
type AnimationChain struct {
	Animation

	animations []Animator
	current    int
}

// Adds animations to the chain. These animations will be played in the order provided.
func (ac *AnimationChain) Add(anims ...Animator) {
	for _, a := range anims {
		ac.animations = append(ac.animations, a)
		ac.Duration += a.GetDuration()
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
}

func (ac *AnimationChain) Render(canvas *Canvas) {
	ac.animations[ac.current].Render(canvas)
}
