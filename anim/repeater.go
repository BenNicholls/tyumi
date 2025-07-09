package anim

import "time"

// Repeater is an animation that triggers a provided function when complete. When set to repeat, the effect is the
// provided function repeating at a set rate. Hooray.
type Repeater struct {
	Animation

	repeatFunc func()
}

func (r *Repeater) Update(delta time.Duration) {
	r.Animation.Update(delta)

	if r.justLooped {
		r.repeatFunc()
	}
}

// Repeaters always repeat so they cannot be OneShot. This is a no-op.
func (r *Repeater) SetOneShot(oneShot bool) {
	return
}

// Creates a repeater animation. While it is playing, repeater_function will be called repeatedly after each duration
// seconds.
func NewRepeaterAnimation(duration time.Duration, repeater_function func()) Repeater {
	return Repeater{
		Animation: Animation{
			Duration: duration,
			Repeat:   true,
		},

		repeatFunc: repeater_function,
	}
}
