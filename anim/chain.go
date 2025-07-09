package anim

import (
	"iter"
	"time"
)

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

func (ac *AnimationChain) Update(delta time.Duration) {
	ac.Animation.Update(delta)

	if ac.IsDone() {
		return
	}

	if ac.justLooped {
		// if we're a repeating chain, this detects if we've looped around and need to restart
		// TODO: this creates a little desync when chain loops! for 1 frame, the animation at
		// the start of the chain will be at zero as it resets, and then afterwards it begins adding
		// deltas, but there will be a little offset unless we set the chain back to zero since we
		// wrapped the chain elapsed duration around and the chain likely is not at exactly 0. so we
		// resynchronize that here, but it remains to be seen whether this makes repeating animations
		// stutter or something???? we'll see.
		ac.elapsed = 0
	}

	//ensure all animations in the chain are reset when the chain is reset
	if ac.elapsed == 0 {
		ac.resetChain()
	} else if ac.animations[ac.current].IsDone() {
		ac.current += 1
	}

	ac.animations[ac.current].Update(delta)
	ac.Updated = ac.animations[ac.current].IsUpdated()
}

func (ac *AnimationChain) resetChain() {
	ac.current = 0
	for _, anim := range ac.animations {
		anim.Start()
	}
}

func (ac *AnimationChain) GetCurrentAnimation() Animator {
	return ac.animations[ac.current]
}

func (ac *AnimationChain) EachAnimation() iter.Seq[Animator] {
	return func(yield func(Animator) bool) {
		for _, a := range ac.animations {
			if !yield(a) {
				return
			}
		}
	}
}
