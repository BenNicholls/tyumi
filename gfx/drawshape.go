package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

// Draws the provided visuals to the entire rectangular area. Optionally takes a 2nd Visuals that will be drawn
// on the interior of the rect, so the normal visuals will be the border.
func (c *Canvas) DrawFilledRect(area vec.Rect, depth int, brush Visuals, inner_brush ...Visuals) {
	if area.Area() == 0 {
		return
	}

	if len(inner_brush) == 0 || area.Area() == 1 {
		for cursor := range vec.EachCoordInArea(area) {
			c.DrawVisuals(cursor, depth, brush)
		}
	} else {
		for cursor := range vec.EachCoordInArea(area) {
			if cursor.IsInPerimeter(area) {
				c.DrawVisuals(cursor, depth, brush)
			} else {
				c.DrawVisuals(cursor, depth, inner_brush[0])
			}
		}
	}
}

// Draws the provided visuals to the border of the rectangular area.
func (c *Canvas) DrawRect(area vec.Rect, depth int, brush Visuals) {
	if area.Area() == 0 {
		return
	}

	//single-cell sized rect.
	if area.Area() == 1 {
		c.DrawVisuals(area.Coord, depth, brush)
		return
	}

	var sides [4]vec.Line
	corners := area.Corners()
	sides[0] = vec.Line{area.Coord, corners[1].Step(vec.DIR_LEFT)}  // top
	sides[1] = vec.Line{corners[1], corners[2].Step(vec.DIR_UP)}    // right
	sides[2] = vec.Line{corners[2], corners[3].Step(vec.DIR_RIGHT)} // bottom
	sides[3] = vec.Line{corners[3], area.Coord.Step(vec.DIR_DOWN)}  // left
	for _, line := range sides {
		c.DrawLine(line, depth, brush)
	}
}

// draws a circle of radius r centered at (px, py), copying the visuals from v, with option to fill the circle with same
// visuals
func (c *Canvas) DrawCircle(center vec.Coord, depth, radius int, visuals Visuals, fill bool) {
	vec.Circle(center, radius, func(pos vec.Coord) {
		c.DrawVisuals(pos, depth, visuals)
	})

	if fill {
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

	//draw corners
	corners := box.Corners()
	c.DrawVisuals(corners[0], depth, NewGlyphVisuals(style.Glyphs[LINK_DR], colours))
	c.DrawVisuals(corners[1], depth, NewGlyphVisuals(style.Glyphs[LINK_DL], colours))
	c.DrawVisuals(corners[2], depth, NewGlyphVisuals(style.Glyphs[LINK_UL], colours))
	c.DrawVisuals(corners[3], depth, NewGlyphVisuals(style.Glyphs[LINK_UR], colours))

	//draw sides
	brush := NewGlyphVisuals(GLYPH_NONE, colours)
	if box.W > 2 {
		brush.Glyph = style.Glyphs[LINK_LR]
		top := vec.Line{corners[0].Step(vec.DIR_RIGHT), corners[1].Step(vec.DIR_LEFT)}
		bottom := vec.Line{corners[2].Step(vec.DIR_LEFT), corners[3].Step(vec.DIR_RIGHT)}
		c.DrawLine(top, depth, brush)
		c.DrawLine(bottom, depth, brush)
	}

	if box.H > 2 {
		brush.Glyph = style.Glyphs[LINK_UD]
		right := vec.Line{corners[1].Step(vec.DIR_DOWN), corners[2].Step(vec.DIR_UP)}
		left := vec.Line{corners[3].Step(vec.DIR_UP), box.Coord.Step(vec.DIR_DOWN)}
		c.DrawLine(left, depth, brush)
		c.DrawLine(right, depth, brush)
	}
}
