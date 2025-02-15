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
