package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

// Animation that makes an area blink. The entire provided area will be filled with visuals Vis while blinking,
// otherwise will draw what what is underneath.
type BlinkAnimation struct {
	Animation
	Vis Visuals //what to draw when the area is blinking

	blinking bool //whether the area is rendering a blink or not
}

func NewBlinkAnimation(pos vec.Coord, size vec.Dims, depth int, vis Visuals, rate int) (ba BlinkAnimation) {
	ba = BlinkAnimation{
		Animation: Animation{
			area:     vec.Rect{pos, size},
			Depth:    depth,
			Repeat:   true,
			Duration: rate,
			reset:    true,
		},
		Vis: vis,
	}

	return
}

func (ba *BlinkAnimation) MoveTo(pos vec.Coord) {
	ba.Animation.MoveTo(pos)
}

func (ba *BlinkAnimation) Update() {
	ba.Animation.Update()

	if ba.ticks == 0 {
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
	Animation

	ToColours   col.Pair
	FromColours col.Pair
}

// Sets up a Fade Animation. Optionally takes a col.Pair for the fade to start from. Omit this to just fade from
// whatever the canvas colours are, which is generally what you want.
func NewFadeAnimation(area vec.Rect, depth int, duration_frames int, fade_colours col.Pair, start_colours ...col.Pair) (fa FadeAnimation) {
	fa = FadeAnimation{
		Animation: Animation{
			area:          area,
			Depth:         depth,
			Duration:      duration_frames,
			AlwaysUpdates: true,
		},
		ToColours: fade_colours,
	}

	if len(start_colours) > 0 {
		fa.FromColours = start_colours[0]
	}

	return
}

func (fa *FadeAnimation) Render(c *Canvas) {
	for cursor := range vec.EachCoordInIntersection(c, fa) {
		dst_cell := c.getCell(cursor)
		if dst_cell.Mode == DRAW_NONE {
			continue
		}

		colours := fa.FromColours
		if colours.Fore == col.NONE {
			colours.Fore = dst_cell.Colours.Fore
		}
		if colours.Back == col.NONE {
			colours.Back = dst_cell.Colours.Back
		}

		c.DrawColours(cursor, fa.Depth, colours.Lerp(fa.ToColours, fa.GetTicks(), fa.Duration-1))
	}

	fa.Updated = false
}

// FlashAnimation makes an area flash once.
type FlashAnimation struct {
	FadeAnimation
}

func NewFlashAnimation(area vec.Rect, depth int, flash_colours col.Pair, duration_frames int) (fa FlashAnimation) {
	fa.FadeAnimation = NewFadeAnimation(area, depth, duration_frames, flash_colours)
	fa.Backwards = true

	return
}

// PulseAnimation makes an area pulse, fading to a set of colours and then fading back
type PulseAnimation struct {
	Animation

	fade FadeAnimation
}

// Creates a pulse animation. duration_frames is the duration of the entire cycle: start -> fade to pulse colour -> fade back
func NewPulseAnimation(area vec.Rect, depth int, duration_frames int, pulse_colours col.Pair) (pa PulseAnimation) {
	pa.Animation = Animation{
		area:          area,
		Depth:         depth,
		Duration:      duration_frames,
		AlwaysUpdates: true,
	}

	pa.fade = NewFadeAnimation(area, depth, duration_frames/2, pulse_colours)
	pa.fade.Start()
	return
}

func (pa *PulseAnimation) Update() {
	if pa.reset {
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
	pa.Animation.SetArea(area)
	pa.fade.SetArea(area)
}

func (pa *PulseAnimation) MoveTo(pos vec.Coord) {
	pa.Animation.MoveTo(pos)
	pa.fade.MoveTo(pos)
}

func (pa *PulseAnimation) Render(canvas *Canvas) {
	pa.fade.Render(canvas)

	pa.Updated = false // don't actually think this is necessary... more of a guard for if the user does something weird.
}
