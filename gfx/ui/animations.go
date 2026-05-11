package ui

import (
	"time"

	"github.com/bennicholls/tyumi/anim"
	"github.com/bennicholls/tyumi/vec"
)

type ElementAnimator interface {
	anim.Animator

	ApplyToElement(e ElementInterface)
}

type ElementMoveAnimation struct {
	anim.Animation

	from, to vec.Coord
}

func (ema *ElementMoveAnimation) ApplyToElement(e ElementInterface) {
	pos := ema.from.Lerp(ema.to, int(ema.GetTicks()), int(ema.GetDuration()))
	e.MoveTo(pos)
}

func NewElementMoveAnimation(from, to vec.Coord, duration time.Duration) ElementMoveAnimation {
	return ElementMoveAnimation{
		Animation: anim.Animation{
			Duration: duration,
		},
		from: from,
		to:   to,
	}
}
