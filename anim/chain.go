package anim

import (
	"iter"
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
