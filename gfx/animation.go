package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
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
	area            vec.Rect
	depth           int       //depth value of the animation
	Vis             Visuals   //what to draw when the area is blinking
	originalVisuals []Visuals //base visuals drawn when area not blinking

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

// captures the state of the drawing canvas so we know what to render when ba.blinking != true
// THINK: should originalVisuals just be a canvas? then we can use the safe canvas drawing functions
// instead of this random []Visuals slice.
func (ba *BlinkAnimation) captureCanvas(c *Canvas) {
	ba.originalVisuals = make([]Visuals, 0, ba.area.Area())

	for cursor := range vec.EachCoord(ba.area) {
		if !c.InBounds(cursor) {
			//dump some crud visual data into the buffer for out of bounds positions, so we'll notice if we draw it
			ba.originalVisuals = append(ba.originalVisuals, NewGlyphVisuals(GLYPH_QUESTION_INVERSE, col.Pair{col.LIME, col.FUSCHIA}))
			continue
		}
		ba.originalVisuals = append(ba.originalVisuals, c.getCell(cursor).Visuals)
	}
}

func (ba *BlinkAnimation) Render(c *Canvas) {
	//capture original canvas state
	if c.dirty || ba.originalVisuals == nil {
		ba.captureCanvas(c)
	}

	if ba.dirty || c.dirty {
		if ba.blinking {
			for cursor := range vec.EachCoord(ba.area) {
				c.DrawVisuals(cursor, ba.depth, ba.Vis)
			}
		} else {
			for cursor := range vec.EachCoord(ba.area) {
				idx := (cursor.X - ba.area.X) + (cursor.Y-ba.area.Y)*ba.area.W
				c.DrawVisuals(cursor, ba.depth, ba.originalVisuals[idx])
			}
		}
		ba.dirty = false
	}
}

func (ba BlinkAnimation) Done() bool {
	return false
}

func (ba BlinkAnimation) Dirty() bool {
	return ba.dirty
}
