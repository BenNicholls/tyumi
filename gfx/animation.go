package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// Anything that can do animations on a Canvas
type Animator interface {
	Update()
	Render(*Canvas)
	Playing() bool
	Done() bool
	IsOneShot() bool
}

// Base struct for animations. Embed this to satisfy Animator interface above.
type Animation struct {
	OneShot bool //indicates animation should play once and then be deleted
	Repeat  bool //animation repeats when finished

	area     vec.Rect
	depth    int  //depth value of the animation
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
		a.reset = false
	} else {
		if a.Repeat {
			a.ticks = util.CycleClamp(a.ticks+1, 0, a.duration-1)
		} else {
			a.ticks += 1
		}
	}

	if a.Done() {
		a.enabled = false
		a.reset = true
	}
}

func (a Animation) Done() bool {
	if a.Repeat || a.ticks < a.duration {
		return false
	}

	return true
}

func (a Animation) IsOneShot() bool {
	return a.OneShot
}

func (a *Animation) MoveTo(pos vec.Coord) {
	a.area.MoveTo(pos.X, pos.Y)
	a.dirty = true
}

func (a Animation) Playing() bool {
	return a.enabled
}

// Starts an animation. If the animation is paused, it restarts it. If animation is playing, does nothing.
func (a *Animation) Start() {
	if a.enabled {
		return
	}
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

// PlayPause pauses playing animations, and plays paused animations.
func (a *Animation) PlayPause() {
	a.enabled = !a.enabled
}

// Animation that makes an area blink. The entire provided area will be filled with visuals Vis while blinking,
// otherwise will draw what what is underneath.
type BlinkAnimation struct {
	Animation
	Vis             Visuals //what to draw when the area is blinking
	originalVisuals Canvas  //base visuals drawn when area not blinking
	blinking        bool    //whether the area is rendering a blink or not
	recapture       bool    //whether we need to recapture the original visuals
}

func NewBlinkAnimation(pos vec.Coord, size vec.Dims, depth int, vis Visuals, rate int) (ba *BlinkAnimation) {
	ba = &BlinkAnimation{
		Animation: Animation{
			area:     vec.Rect{pos, size},
			depth:    depth,
			Repeat:   true,
			duration: rate,
			reset:    true,
		},
		Vis:       vis,
		recapture: true,
	}

	return
}

func (ba *BlinkAnimation) MoveTo(pos vec.Coord) {
	ba.Animation.MoveTo(pos)
	ba.recapture = true
}

func (ba *BlinkAnimation) Update() {
	ba.Animation.Update()

	if ba.ticks == 0 {
		ba.blinking = !ba.blinking
		ba.dirty = true
	}
}

func (ba *BlinkAnimation) Render(c *Canvas) {
	//capture original canvas state
	if c.dirty || ba.recapture {
		ba.originalVisuals = c.CopyArea(ba.area)
		ba.recapture = false
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

// FlashAnimation makes an area flash once.
type FlashAnimation struct {
	Animation
	flashColours    col.Pair
	originalColours []col.Pair
}

func NewFlashAnimation(area vec.Rect, depth int, flashColours col.Pair, duration_frames int) (fa *FlashAnimation) {
	fa = &FlashAnimation{
		Animation: Animation{
			area:     area,
			depth:    depth,
			duration: duration_frames,
		},
		flashColours: flashColours,
	}

	return
}

func (fa *FlashAnimation) Update() {
	fa.Animation.Update()
	fa.dirty = true

	if fa.reset {
		fa.originalColours = nil
	}
}

func (fa *FlashAnimation) Render(c *Canvas) {
	if !fa.dirty || !fa.enabled {
		return
	}

	if fa.originalColours == nil {
		//populate original colours to lerp to
		fa.originalColours = make([]col.Pair, fa.area.Area())
		for cursor := range vec.EachCoord(fa.area) {
			cell := c.getCell(cursor)
			col_index := cursor.Subtract(fa.area.Coord).ToIndex(fa.area.W)
			fa.originalColours[col_index] = cell.Colours
		}
	}

	for cursor := range vec.EachCoord(fa.area) {
		col_index := cursor.Subtract(fa.area.Coord).ToIndex(fa.area.W)
		c.DrawColours(cursor, fa.depth, fa.flashColours.Lerp(fa.originalColours[col_index], fa.ticks, fa.duration))
	}

	fa.dirty = false
}
