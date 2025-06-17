package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

// Draws the provided visuals to the entire rectangular area. Optionally takes a 2nd Visuals that will be drawn
// on the perimeter of the rect, making a border.
func (c *Canvas) DrawFilledRect(area vec.Rect, depth int, brush Visuals, border_brush ...Visuals) {
	for cursor := range vec.EachCoordInArea(area) {
		if len(border_brush) > 0 && cursor.IsInPerimeter(area) {
			c.DrawVisuals(cursor, depth, border_brush[0])
		} else {
			c.DrawVisuals(cursor, depth, brush)
		}
	}
}

// Draws the provided visuals to the border of the rectangular area.
func (c *Canvas) DrawRect(area vec.Rect, depth int, brush Visuals) {
	for cursor := range area.EachCoordInPerimeter() {
		c.DrawVisuals(cursor, depth, brush)
	}
}

// draws a circle of radius r centered at (px, py), copying the visuals from v, with option to fill the circle with same
// visuals
func (c *Canvas) DrawCircle(center vec.Coord, radius, depth int, visuals Visuals, filled bool) {
	vec.CircleFunc(center, radius, func(pos vec.Coord) {
		c.DrawVisuals(pos, depth, visuals)
	})

	if filled {
		c.FloodFill(center, depth, visuals)
	}
}

// DrawBox draws a box on the edges of the provided area, respecting depth. line is the type of line to use: thin or
// thick. Does not draw anything if the dimensions of the box aren't at least 2x2.
func (c *Canvas) DrawBox(box vec.Rect, depth int, line LineType, colours col.Pair) {
	// if not doing a linking line, just call the regular old rect drawing function
	if line == LINETYPE_NONE {
		c.DrawRect(box, depth, NewGlyphVisuals(GLYPH_BLOCK, colours))
		return
	}

	//these boxes are too small to draw
	if box.Area() == 0 || box.W == 1 || box.H == 1 {
		return
	}

	style := LineStyles[line]
	sides := box.Sides()

	//draw corners
	c.DrawVisuals(sides[0].Start, depth, NewGlyphVisuals(style.Glyphs[LINK_DR], colours)) //TOPLEFT
	c.DrawVisuals(sides[1].Start, depth, NewGlyphVisuals(style.Glyphs[LINK_DL], colours)) //TOPRIGHT
	c.DrawVisuals(sides[2].Start, depth, NewGlyphVisuals(style.Glyphs[LINK_UL], colours)) //BOTTOMRIGHT
	c.DrawVisuals(sides[3].Start, depth, NewGlyphVisuals(style.Glyphs[LINK_UR], colours)) //BOTTOMLEFT

	//draw sides
	brush := NewGlyphVisuals(GLYPH_NONE, colours)
	if box.W > 2 {
		brush.Glyph = style.Glyphs[LINK_LR]
		c.DrawLine(vec.Line{sides[0].Start.Step(vec.DIR_RIGHT), sides[0].End.Step(vec.DIR_LEFT)}, depth, brush) //TOP
		c.DrawLine(vec.Line{sides[2].Start.Step(vec.DIR_LEFT), sides[2].End.Step(vec.DIR_RIGHT)}, depth, brush) //BOTTOM
	}

	if box.H > 2 {
		brush.Glyph = style.Glyphs[LINK_UD]
		c.DrawLine(vec.Line{sides[1].Start.Step(vec.DIR_DOWN), sides[1].End.Step(vec.DIR_UP)}, depth, brush) //TOP
		c.DrawLine(vec.Line{sides[3].Start.Step(vec.DIR_UP), sides[3].End.Step(vec.DIR_DOWN)}, depth, brush) //BOTTOM
	}
}
