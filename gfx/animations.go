package gfx

import (
	"time"

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

// Anything that can animate a single visuals
type VisualAnimator interface {
	anim.Animator

	ApplyToVisuals(vis Visuals) Visuals
}

// Base for Canvas Animations, satisfying the CanvasAnimator interface above. Compose canvas animations around this!
// If area is zero, applies the animation to the entire canvas.
type CanvasAnimation struct {
	anim.Animation

	Depth int      // depth value of the animation
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

func (cac CanvasAnimationChain) Bounds() vec.Rect {
	current := cac.GetCurrentAnimation()
	if canvasAnim, ok := current.(CanvasAnimator); ok {
		return canvasAnim.Bounds()
	} else {
		return vec.Rect{}
	}
}

func (cac *CanvasAnimationChain) Render(canvas *Canvas) {
	if currentCanvas, ok := cac.GetCurrentAnimation().(CanvasAnimator); ok {
		currentCanvas.Render(canvas)
	}

	cac.Updated = false
}

func (cac *CanvasAnimationChain) ApplyToVisuals(vis Visuals) Visuals {
	if currentVisualAnim, ok := cac.GetCurrentAnimation().(VisualAnimator); ok {
		return currentVisualAnim.ApplyToVisuals(vis)
	}

	return vis
}

// Animation that makes an area blink. The entire provided area will be filled with visuals Vis while blinking,
// otherwise will draw what is underneath.
type BlinkAnimation struct {
	CanvasAnimation

	Vis      Visuals //what to draw when the area is blinking
	blinking bool    //whether the area is rendering a blink or not
}

func NewBlinkAnimation(area vec.Rect, depth int, vis Visuals, rate time.Duration) (ba BlinkAnimation) {
	ba.Repeat = true
	ba.Duration = rate
	ba.Depth = depth
	ba.area = area
	ba.Vis = vis

	return
}

func (ba *BlinkAnimation) Update(delta time.Duration) {
	ba.Animation.Update(delta)

	if ba.JustLooped() {
		ba.blinking = !ba.blinking
		ba.Updated = true
	}
}

func (ba *BlinkAnimation) Render(c *Canvas) {
	ba.Updated = false
	if !ba.blinking {
		return
	}

	bounds := ba.Bounds()
	if bounds.Area() == 0 { // no bounds set, apply to full canvas
		bounds = c.Bounds()
	}

	for cursor := range vec.EachCoordInIntersection(c, bounds) {
		c.DrawVisuals(cursor, ba.Depth, ba.Vis)
	}
}

func (ba *BlinkAnimation) ApplyToVisuals(vis Visuals) (result Visuals) {
	ba.Updated = false

	if !ba.blinking {
		return vis
	}

	result = ba.Vis
	result = result.ReplaceChars(TEXT_DEFAULT, vis.Chars)
	result.Colours = result.Colours.Replace(col.NONE, vis.Colours)

	return
}

// FadeAnimation makes an area fade to the provided colours (ToColours). If FromColours is non-zero, it will start the
// fade from there. Otherwise uses whatever colours are on the canvas.
type FadeAnimation struct {
	CanvasAnimation

	ToColours, FromColours col.Pair
}

// Sets up a Fade Animation. Optionally takes a col.Pair for the fade to start from. Omit this to just fade from
// whatever the canvas colours are, which is generally what you want.
func NewFadeAnimation(area vec.Rect, depth int, duration time.Duration, to_colour col.Pair, from_colour ...col.Pair) (fa FadeAnimation) {
	fa.Duration = duration
	fa.AlwaysUpdates = true
	fa.Depth = depth
	fa.area = area
	fa.ToColours = to_colour

	if len(from_colour) > 0 {
		fa.FromColours = from_colour[0]
	}

	return
}

// Sets up a Fade Out animation. Both foreground and background are faded to the specified colour.
func NewFadeOutAnimation(area vec.Rect, depth int, duration time.Duration, colour col.Colour) FadeAnimation {
	return NewFadeAnimation(area, depth, duration, col.Pair{colour, colour})
}

// Sets up a Fade In animation. Both foreground and background are faded from the specified colour.
func NewFadeInAnimation(area vec.Rect, depth int, duration time.Duration, colour col.Colour) (fa FadeAnimation) {
	fa = NewFadeAnimation(area, depth, duration, col.Pair{colour, colour})
	fa.Backwards = true

	return
}

// Sets up a flash animation. The colour immediately is set to the provided flash colours, and then the area fades
// back to the original colours over the duration.
func NewFlashAnimation(area vec.Rect, depth int, duration time.Duration, flash_colours col.Pair) (fa FadeAnimation) {
	fa = NewFadeAnimation(area, depth, duration, flash_colours)
	fa.Backwards = true

	return
}

func (fa *FadeAnimation) Render(c *Canvas) {
	bounds := fa.Bounds()
	if bounds.Area() == 0 { // no bounds set, apply to full canvas
		bounds = c.Bounds()
	}

	toColours := fa.ToColours.Replace(COL_DEFAULT, c.DefaultColours())
	for cursor := range vec.EachCoordInIntersection(c, bounds) {
		dst_cell := c.getCell(cursor)
		fromColours := fa.FromColours.Replace(COL_DEFAULT, c.DefaultColours())
		fromColours = fromColours.Replace(col.NONE, dst_cell.Colours)

		c.DrawColours(cursor, fa.Depth, fromColours.Lerp(toColours, int(fa.GetTicks()), int(fa.Duration)))
	}

	fa.Updated = false
}

func (fa *FadeAnimation) ApplyToVisuals(vis Visuals) (result Visuals) {
	result = vis

	toColours := fa.ToColours.Replace(col.NONE, vis.Colours)
	fromColours := fa.FromColours.Replace(col.NONE, vis.Colours)

	result.Colours = fromColours.Lerp(toColours, int(fa.GetTicks()), int(fa.Duration-1))
	fa.Updated = false

	return
}

// PulseAnimation makes an area pulse, fading to a set of colours and then fading back
type PulseAnimation struct {
	CanvasAnimationChain
}

// Creates a pulse animation. duration_frames is the duration of the entire cycle: start -> fade to pulse colour -> fade back
func NewPulseAnimation(area vec.Rect, depth int, duration time.Duration, pulse_colours col.Pair) (pa PulseAnimation) {
	pa.AlwaysUpdates = true

	for i := range 2 {
		fade := NewFadeAnimation(area, depth, duration/2, pulse_colours)
		fade.Depth = depth
		fade.area = area
		if i == 1 {
			fade.Backwards = true
		}
		pa.Add(&fade)
	}

	return
}

func (pa *PulseAnimation) SetArea(area vec.Rect) {
	for anim := range pa.EachAnimation() {
		if canvasAnim, ok := anim.(*CanvasAnimation); ok {
			canvasAnim.SetArea(area)
		}
	}
}

func (pa *PulseAnimation) MoveTo(pos vec.Coord) {
	for anim := range pa.EachAnimation() {
		if canvasAnim, ok := anim.(*CanvasAnimation); ok {
			canvasAnim.MoveTo(pos)
		}
	}
}
