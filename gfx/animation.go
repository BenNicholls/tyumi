package gfx

import (
	"github.com/bennicholls/tyumi/vec"
)

// Anything that can do animations on a Canvas
type Animator interface {
	Update()
	Render(*Canvas)
	Done() bool
	Dirty() bool
}

// Animation that makes an area blink
type BlinkAnimation struct {
	area  vec.Rect
	depth int     //depth value of the animation
	Vis   Visuals //what to draw when the area is blinking

	enabled bool //animation is playing
	dirty   bool //animation needs to be re-rendered

	blinkRate   int  //frames between state changes
	blinkFrames int  //frames since last blink
	blinking    bool //whether the area is rendering a blink or not
}

func NewBlinkAnimation(pos vec.Coord, size vec.Dims, depth int, vis Visuals, rate int) (ba *BlinkAnimation) {
	ba = &BlinkAnimation{
		area:      vec.Rect{pos, size},
		depth:     depth,
		Vis:       vis,
		enabled:   true,
		blinkRate: rate,
	}

	return
}

func (ba *BlinkAnimation) Move(dx, dy int) {
	ba.area.Move(dx, dy)
	ba.dirty = true
}

func (ba *BlinkAnimation) MoveTo(x, y int) {
	ba.area.MoveTo(x, y)
	ba.dirty = true
}

func (ba *BlinkAnimation) Update() {
	if !ba.enabled {
		return
	}

	ba.blinkFrames++
	if ba.blinkFrames == ba.blinkRate {
		ba.blinking = !ba.blinking
		ba.blinkFrames = 0
		ba.dirty = true
	}
}

func (ba *BlinkAnimation) Render(c *Canvas) {
	ba.dirty = false

	if !ba.blinking {
		return
	}

	w := ba.area.W

	for i := 0; i < ba.area.Area(); i++ {
		offset := vec.Coord{i%w, i/w}
		c.DrawVisuals(ba.area.Coord.Add(offset), ba.depth, ba.Vis)
	}
}

func (ba BlinkAnimation) Done() bool {
	return false
}

func (ba BlinkAnimation) Dirty() bool {
	return ba.dirty
}
