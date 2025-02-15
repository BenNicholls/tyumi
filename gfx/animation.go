package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

// Anything that can do animations on a Canvas
type Animator interface {
	vec.Bounded
	Update()
	Render(*Canvas)
	Playing() bool
	Done() bool
	IsOneShot() bool
}

// Base struct for animations. Embed this to satisfy Animator interface above.
type Animation struct {
	OneShot  bool //indicates animation should play once and then be deleted
	Repeat   bool //animation repeats when finished
	Area     vec.Rect
	Depth    int //depth value of the animation
	Duration int //duration of animation in ticks

	ticks   int  //incremented each update
	enabled bool //animation is playing
	dirty   bool //animation needs to be re-rendered
	reset   bool //indicates animation should reset and start over.
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
			a.ticks = (a.ticks + 1) % a.Duration
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
	return !a.Repeat && a.ticks >= a.Duration
}

func (a Animation) IsOneShot() bool {
	return a.OneShot
}

func (a *Animation) MoveTo(pos vec.Coord) {
	a.Area.MoveTo(pos.X, pos.Y)
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

// Plays an animation. If the animation is paused, continues it. If it's already playing, restarts it.
func (a *Animation) Play() {
	if a.enabled {
		a.reset = true
	} else {
		a.enabled = true
	}
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

func (a Animation) Bounds() vec.Rect {
	return a.Area
}

// Animation that makes an area blink. The entire provided area will be filled with visuals Vis while blinking,
// otherwise will draw what what is underneath.
type BlinkAnimation struct {
	Animation
	Vis Visuals //what to draw when the area is blinking

	originalVisuals Canvas //base visuals drawn when area not blinking
	blinking        bool   //whether the area is rendering a blink or not
	recapture       bool   //whether we need to recapture the original visuals
}

func NewBlinkAnimation(pos vec.Coord, size vec.Dims, depth int, vis Visuals, rate int) (ba BlinkAnimation) {
	ba = BlinkAnimation{
		Animation: Animation{
			Area:     vec.Rect{pos, size},
			Depth:    depth,
			Repeat:   true,
			Duration: rate,
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
		ba.originalVisuals = c.CopyArea(ba.Area)
		ba.recapture = false
	}

	if ba.dirty || c.dirty {
		if ba.blinking {
			for cursor := range vec.EachCoordInArea(ba.Area) {
				c.DrawVisuals(cursor, ba.Depth, ba.Vis)
			}
		} else {
			ba.originalVisuals.Draw(c, ba.Area.Coord, ba.Depth)
		}
		ba.dirty = false
	}
}

// FlashAnimation makes an area flash once.
type FlashAnimation struct {
	Animation
	Colours col.Pair

	originalColours []col.Pair
}

func NewFlashAnimation(area vec.Rect, depth int, flash_colours col.Pair, duration_frames int) (fa FlashAnimation) {
	fa = FlashAnimation{
		Animation: Animation{
			Area:     area,
			Depth:    depth,
			Duration: duration_frames,
		},
		Colours: flash_colours,
	}

	return
}

func (fa *FlashAnimation) Update() {
	if fa.reset {
		fa.originalColours = nil
	}

	fa.Animation.Update()
	fa.dirty = true
}

func (fa *FlashAnimation) Render(c *Canvas) {
	if !fa.dirty || !fa.enabled {
		return
	}

	if fa.originalColours == nil {
		//populate original colours to lerp to
		fa.originalColours = make([]col.Pair, fa.Area.Area())
		for cursor := range vec.EachCoordInArea(fa.Area) {
			cell := c.getCell(cursor)
			col_index := cursor.Subtract(fa.Area.Coord).ToIndex(fa.Area.W)
			fa.originalColours[col_index] = cell.Colours
		}
	}

	for cursor := range vec.EachCoordInArea(fa.Area) {
		col_index := cursor.Subtract(fa.Area.Coord).ToIndex(fa.Area.W)
		c.DrawColours(cursor, fa.Depth, fa.Colours.Lerp(fa.originalColours[col_index], fa.ticks, fa.Duration-1))
	}

	fa.dirty = false
}
