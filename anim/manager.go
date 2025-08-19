package anim

import (
	"iter"
	"slices"
	"time"

	"github.com/bennicholls/tyumi/util"
)

type Manager interface {
	AddAnimation(Animator)
	RemoveAnimation(Animator)
	UpdateAnimations(time.Duration)
	HasBlockingAnimation() bool
	EachAnimation() iter.Seq[Animator]
	EachPlayingAnimation() iter.Seq[Animator]
}

// An animation manager is an embeddable struct that gives you the ability to add, update, and generally take care
// of animations.
type AnimationManager struct {
	animations []Animator

	// will be true if an animation has stopped in the most recent frame (the time since the most recent call to Update)
	AnimationJustStopped bool

	// will be true if an animation has updated in the most recent frame (the time since the most recent call to Update)
	AnimationJustUpdated bool
}

func (am AnimationManager) CountAnimations() int {
	return len(am.animations)
}

// Adds an animation. Note that this does NOT start the animation, you have to do that manually. Does not allow adding
// a duplicate animation.
func (am *AnimationManager) AddAnimation(animation Animator) {
	if am.animations == nil {
		am.animations = make([]Animator, 0)
	}

	if slices.Contains(am.animations, animation) {
		return
	}

	am.animations = append(am.animations, animation)
}

// Adds an animation in one shot mode. This sets the OneShot flag for you, and it starts the animation! OneShot animations
// play once, stop, and are disposed of automatically. Just fire and forget!!
func (am *AnimationManager) AddOneShotAnimation(animation Animator) {
	animation.SetOneShot(true)
	animation.Start()
	am.AddAnimation(animation)
}

// Removes an animation. Note that OneShot animations are disposed of automatically, you do not need to remove them
// yourself.
func (am *AnimationManager) RemoveAnimation(animation Animator) {
	am.animations = util.DeleteElement(am.animations, animation)
}

// Updates all playing animations. Also updates the AnimationJustStopped and AnimationJustUpdated flags; if an animation
// has stopped/updated since the previous call to UpdateAnimations these flags will be true.
//
// NOTE that this just updates the animation (ticks the internal counter forward, handles pause/play/reset/done/whatever
// states, stuff like that), it does NOT apply the effects of the animations. For example, it doesn't make a
// CanvasAnimation render and draw things on a canvas. Specific types of animations need to be applied in whatever way
// or time is appropriate for them.
func (am *AnimationManager) UpdateAnimations(delta time.Duration) {
	am.AnimationJustStopped = false
	am.AnimationJustUpdated = false
	for animation := range am.EachAnimation() {
		if animation.IsPlaying() {
			animation.Update(delta)
			if animation.IsUpdated() {
				am.AnimationJustUpdated = true
			}
		}

		if animation.stoppedSinceLastUpdate() {
			am.AnimationJustStopped = true
			animation.clearFlags()
		}
	}

	am.animations = slices.DeleteFunc(am.animations, func(a Animator) bool {
		return a.IsOneShot() && a.IsDone()
	})
}

// Checks and reports if there is an active blocking animation in the Animation Manager.
func (am *AnimationManager) HasBlockingAnimation() bool {
	for animation := range am.EachPlayingAnimation() {
		if animation.IsBlocking() {
			return true
		}
	}

	return false
}

func (am *AnimationManager) EachAnimation() iter.Seq[Animator] {
	return func(yield func(Animator) bool) {
		for _, a := range am.animations {
			if !yield(a) {
				return
			}
		}
	}
}

func (am *AnimationManager) EachPlayingAnimation() iter.Seq[Animator] {
	return func(yield func(Animator) bool) {
		for _, a := range am.animations {
			if !a.IsPlaying() {
				continue
			}

			if !yield(a) {
				return
			}
		}
	}
}
