package gfx

import (
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// Anything that can do animations on a Canvas
type Animator interface {
	Update()
	Render(*Canvas)
	Done() bool
}

// Base struct for animations. Embed this to satisfy Animator interface above.
type Animation struct {
	area     vec.Rect
	depth    int  //depth value of the animation
	repeat   bool //animation repeats when finished
	duration int  //duration of animation in ticks
	ticks    int  //incremented each update
	enabled  bool //animation is playing
	dirty    bool //animation needs to be re-rendered
	reset    bool //indicates animation should reset and start over.
}

func (a *Animation) Update() {
	if !a.enabled {
		return
	}

	if a.reset {
		a.ticks = 0
		a.dirty = true
	} else {
		if a.repeat {
			a.ticks = util.CycleClamp(a.ticks+1, 0, a.duration-1)
		} else {
			a.ticks += 1
		}
	}
}

func (a *Animation) Render(c *Canvas) {

}

func (a Animation) Done() bool {
	if a.repeat || a.ticks < a.duration {
		return false
	}

	return true
}

func (a Animation) Dirty() bool {
	return a.dirty
}

func (a *Animation) Move(dx, dy int) {
	a.area.Move(dx, dy)
	a.dirty = true
}

func (a *Animation) MoveTo(x, y int) {
	a.area.MoveTo(x, y)
	a.dirty = true
}

func (a *Animation) Enable() {
	a.enabled = true
}

// can also be used as a pause button.
func (a *Animation) Disable() {
	a.enabled = false
}

// can also be used as a pause/play button
func (a *Animation) ToggleEnabled() {
	a.enabled = !a.enabled
}

func (a *Animation) Start() {
	if a.enabled {
		return
	}

	a.enabled = true
	a.reset = true
}

// Animation that makes an area blink. The entire provided area will be filled with visuals Vis while blinking,
// otherwise will draw what what is underneath.
type BlinkAnimation struct {
	Animation
	Vis             Visuals //what to draw when the area is blinking
	originalVisuals Canvas  //base visuals drawn when area not blinking
	blinking        bool    //whether the area is rendering a blink or not
}

func NewBlinkAnimation(pos vec.Coord, size vec.Dims, depth int, vis Visuals, rate int) (ba *BlinkAnimation) {
	ba = &BlinkAnimation{
		Animation: Animation{
			area:     vec.Rect{pos, size},
			depth:    depth,
			repeat:   true,
			duration: rate,
			reset:    true,
		},
		Vis: vis,
	}

	return
}

func (ba *BlinkAnimation) Update() {
	if !ba.enabled {
		return
	}

	ba.Animation.Update()

	if ba.ticks == 0 {
		ba.blinking = !ba.blinking
		ba.dirty = true
	}
}

func (ba *BlinkAnimation) Render(c *Canvas) {
	//capture original canvas state
	if c.dirty || ba.reset {
		ba.originalVisuals = c.CopyArea(ba.area)
		ba.reset = false
	}

	if ba.dirty || c.dirty {
		if ba.blinking {
			for cursor := range vec.EachCoord(ba.area) {
				c.DrawVisuals(cursor, ba.depth, ba.Vis)
			}
		} else {
			ba.originalVisuals.DrawToCanvas(c, ba.area.Coord, ba.depth)
		}
		ba.dirty = false
	}
}
