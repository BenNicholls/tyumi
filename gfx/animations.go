package gfx

import (
	"github.com/bennicholls/tyumi/anim"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

// Anything that can do animations on a Canvas
type CanvasAnimator interface {
	anim.Animator
	vec.Bounded

	Render(*Canvas)
}

// Base for Canvas Animations, satisfying the CanvasAnimator interface above. Compose canvas animations around this!
type CanvasAnimation struct {
	anim.Animation

	Depth int      //depth value of the animation
	area  vec.Rect // Area this animation affects. Use MoveTo() and Resize() to change this value.
}

func (ca CanvasAnimation) Bounds() vec.Rect {
	return ca.area
}

func (ca *CanvasAnimation) SetArea(area vec.Rect) {
	ca.MoveTo(area.Coord)
	ca.Resize(area.Dims)
}

func (ca *CanvasAnimation) Resize(size vec.Dims) {
	if ca.area.Dims == size {
		return
	}

	ca.area.Dims = size
	ca.Updated = true
}

func (ca *CanvasAnimation) MoveTo(pos vec.Coord) {
	if ca.area.Coord == pos {
		return
	}

	ca.area.Coord = pos
	ca.Updated = true
}

// An animation chain that can hold and render CanvasAnimators. If an animation in the chain is a canvas animator it
// will be rendered as normal.
type CanvasAnimationChain struct {
	anim.AnimationChain
}

func (cac *CanvasAnimationChain) Render(canvas *Canvas) {
	if currentCanvas, ok := cac.GetCurrentAnimation().(CanvasAnimator); ok {
		currentCanvas.Render(canvas)
	}

	cac.Updated = false
}

// Animation that makes an area blink. The entire provided area will be filled with visuals Vis while blinking,
// otherwise will draw what what is underneath.
type BlinkAnimation struct {
	CanvasAnimation

	Vis      Visuals //what to draw when the area is blinking
	blinking bool    //whether the area is rendering a blink or not
}

func NewBlinkAnimation(pos vec.Coord, size vec.Dims, depth int, vis Visuals, rate int) BlinkAnimation {
	return BlinkAnimation{
		CanvasAnimation: CanvasAnimation{
			Animation: anim.Animation{
				Repeat:   true,
				Duration: rate,
			},
			Depth: depth,
			area:  vec.Rect{pos, size},
		},
		Vis: vis,
	}
}

func (ba *BlinkAnimation) Update() {
	ba.Animation.Update()

	if ba.GetTicks() == 0 {
		ba.blinking = !ba.blinking
		ba.Updated = true
	}
}

func (ba *BlinkAnimation) Render(c *Canvas) {
	if ba.blinking {
		for cursor := range vec.EachCoordInArea(ba) {
			c.DrawVisuals(cursor, ba.Depth, ba.Vis)
		}
	}
	ba.Updated = false
}

// FadeAnimation makes an area fade to the provided colours (ToColours). If FromColours is non-zero, it will start the
// fade from there. Otherwise uses whatever colours are on the canvas.
type FadeAnimation struct {
	CanvasAnimation

	ToColours, FromColours col.Pair
}

// Sets up a Fade Animation. Optionally takes a col.Pair for the fade to start from. Omit this to just fade from
// whatever the canvas colours are, which is generally what you want.
func NewFadeAnimation(area vec.Rect, depth int, duration_frames int, fade_colours col.Pair, start_colours ...col.Pair) (fa FadeAnimation) {
	fa = FadeAnimation{
		CanvasAnimation: CanvasAnimation{
			Animation: anim.Animation{
				Duration:      duration_frames,
				AlwaysUpdates: true,
			},
			Depth: depth,
			area:  area,
		},
		ToColours: fade_colours,
	}

	if len(start_colours) > 0 {
		fa.FromColours = start_colours[0]
	}

	return
}

// Sets up a Fade Out animation. Both foreground and background are faded to the specified colour.
func NewFadeOutAnimation(area vec.Rect, depth int, duration_frames int, colour col.Colour) (fa FadeAnimation) {
	fa = FadeAnimation{
		CanvasAnimation: CanvasAnimation{
			Animation: anim.Animation{
				Duration:      duration_frames,
				AlwaysUpdates: true,
			},
			Depth: depth,
			area:  area,
		},
		ToColours: col.Pair{colour, colour},
	}

	return
}

// Sets up a Fade In animation. Both foreground and background are faded from the specified colour.
func NewFadeInAnimation(area vec.Rect, depth int, duration_frames int, colour col.Colour) (fa FadeAnimation) {
	fa = FadeAnimation{
		CanvasAnimation: CanvasAnimation{
			Animation: anim.Animation{
				Duration:      duration_frames,
				AlwaysUpdates: true,
				Backwards:     true,
			},
			Depth: depth,
			area:  area,
		},
		ToColours: col.Pair{colour, colour},
	}

	return
}

// Sets up a flash animation. The colour immediately is set to the provided flash colours, and then the area fades
// back to the original colours over the duration.
func NewFlashAnimation(area vec.Rect, depth int, duration_frames int, flash_colours col.Pair) (fa FadeAnimation) {
	fa = NewFadeAnimation(area, depth, duration_frames, flash_colours)
	fa.Backwards = true

	return
}

func (fa *FadeAnimation) Render(c *Canvas) {
	toColours := fa.ToColours
	if toColours.Fore == COL_DEFAULT {
		toColours.Fore = c.DefaultColours().Fore
	}
	if toColours.Back == COL_DEFAULT {
		toColours.Back = c.DefaultColours().Back
	}

	for cursor := range vec.EachCoordInIntersection(c, fa) {
		dst_cell := c.getCell(cursor)

		fromColours := fa.FromColours
		if fromColours.Fore == COL_DEFAULT {
			fromColours.Fore = c.DefaultColours().Fore
		}
		if fromColours.Back == COL_DEFAULT {
			fromColours.Back = c.DefaultColours().Back
		}

		if fromColours.Fore == col.NONE {
			fromColours.Fore = dst_cell.Colours.Fore
		}
		if fromColours.Back == col.NONE {
			fromColours.Back = dst_cell.Colours.Back
		}

		c.DrawColours(cursor, fa.Depth, fromColours.Lerp(toColours, fa.GetTicks(), fa.Duration-1))
	}

	fa.Updated = false
}

// PulseAnimation makes an area pulse, fading to a set of colours and then fading back
type PulseAnimation struct {
	CanvasAnimation

	fade FadeAnimation
}

// Creates a pulse animation. duration_frames is the duration of the entire cycle: start -> fade to pulse colour -> fade back
func NewPulseAnimation(area vec.Rect, depth int, duration_frames int, pulse_colours col.Pair) (pa PulseAnimation) {
	pa.CanvasAnimation = CanvasAnimation{
		Animation: anim.Animation{
			Duration:      duration_frames,
			AlwaysUpdates: true,
		},
		Depth: depth,
		area:  area,
	}

	pa.fade = NewFadeAnimation(area, depth, duration_frames/2, pulse_colours)
	pa.fade.Start()
	return
}

func (pa *PulseAnimation) Update() {
	if pa.IsResetting() {
		pa.fade.Backwards = false
		pa.fade.Start()
	}
	pa.Animation.Update()

	if pa.fade.IsDone() {
		pa.fade.Backwards = !pa.fade.Backwards
		pa.fade.Start()
	}
	pa.fade.Update()
}

func (pa *PulseAnimation) SetArea(area vec.Rect) {
	pa.CanvasAnimation.SetArea(area)
	pa.fade.SetArea(area)
}

func (pa *PulseAnimation) MoveTo(pos vec.Coord) {
	pa.CanvasAnimation.MoveTo(pos)
	pa.fade.MoveTo(pos)
}

func (pa *PulseAnimation) Render(canvas *Canvas) {
	pa.fade.Render(canvas)
	pa.Updated = false // don't actually think this is necessary... more of a guard for if the user does something weird.
}
