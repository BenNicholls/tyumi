package gfx

import (
	"github.com/bennicholls/tyumi/vec"
)

//Anything that can do animations on a Canvas
type Animator interface {
	Update()
	Render(*Canvas)
	Done() bool
	Dirty() bool
}

//Animation that makes an area blink
type BlinkAnimation struct {
	area vec.Rect
	z    int     //z value of the animation
	vis  Visuals //what to draw when the area is blinking

	enabled bool //animation is playing
	dirty   bool //animation needs to be re-rendered

	blinkRate   int  //frames between state changes
	blinkFrames int  //frames since last blink
	blinking    bool //whether the area is rendering a blink or not
}

func NewBlinkAnimation(w, h, x, y, z int, vis Visuals, rate int) (ba *BlinkAnimation) {
	ba = &BlinkAnimation{
		area:      vec.Rect{w, h, x, y},
		z:         z,
		vis:       vis,
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

	x, y := ba.area.Pos()
	w, h := ba.area.Dims()

	for i := 0; i < w*h; i++ {
		c.DrawVisuals(x+i%w, y+i/w, ba.z, ba.vis)
	}
}

func (ba BlinkAnimation) Done() bool {
	return false
}

func (ba BlinkAnimation) Dirty() bool {
	return ba.dirty
}
