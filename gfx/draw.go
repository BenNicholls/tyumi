package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

type TextCellPosition int

const (
	DRAW_TEXT_LEFT  TextCellPosition = 0
	DRAW_TEXT_RIGHT TextCellPosition = 1
)

func (c *Canvas) Draw(pos vec.Coord, depth int, d Drawable) {
	c.DrawVisuals(pos, depth, d.Visuals())
}

// THINK: this checks/sets the depth 3-4 times i think. hmmm.
func (c *Canvas) DrawVisuals(pos vec.Coord, depth int, v Visuals) {
	if !c.InBounds(pos) {
		return
	}

	switch v.Mode {
	case DRAW_GLYPH:
		c.setGlyph(pos, depth, v.Glyph)
	case DRAW_TEXT:
		c.setText(pos, depth, v.Chars[0], v.Chars[1])
	case DRAW_NONE:
		c.setBlank(pos)
		return // if we are not drawing this cell we can skip setting the colours below
	}

	c.setColours(pos, depth, v.Colours)
}

// DrawNone sets the cell at pos to mode DRAW_NONE, which prevents it from being drawn.
func (c *Canvas) DrawNone(pos vec.Coord) {
	if !c.InBounds(pos) {
		return
	}

	c.setBlank(pos)
}

// DrawColours draws a colour pair (fore/back) to a cell at pos, respecting depth.
func (c *Canvas) DrawColours(pos vec.Coord, depth int, colours col.Pair) {
	if !c.InBounds(pos) {
		return
	}

	c.setColours(pos, depth, colours)
}

func (c *Canvas) DrawGlyph(pos vec.Coord, depth int, glyph int) {
	if !c.InBounds(pos) {
		return
	}

	c.setGlyph(pos, depth, glyph)
}

func (c *Canvas) DrawText(pos vec.Coord, depth int, txt string, colours col.Pair, start_pos TextCellPosition) {
	//build []rune version of txt string
	var text_runes []rune = make([]rune, 0, len(txt))
	if start_pos == DRAW_TEXT_RIGHT { //pad start with a space if we're starting on the right
		text_runes = append(text_runes, TEXT_NONE)
	}
	text_runes = append(text_runes, []rune(txt)...)
	if len(text_runes)%2 != 0 { //pad end if we're ending on the left
		text_runes = append(text_runes, TEXT_NONE)
	}

	//iterate by pairs of runes, drawing 1 cell per loop
	for i := 0; i < len(text_runes); i += 2 {
		cursor := vec.Coord{pos.X + i/2, pos.Y}
		if !c.InBounds(cursor) { //make sure we're drawing in the canvas.
			continue
		}

		c.setText(cursor, depth, text_runes[i], text_runes[i+1])
		c.setColours(cursor, depth, colours)
	}
}

// Draws the provided visuals to the rectangular area
func (c *Canvas) DrawRect(area vec.Rect, depth int, v Visuals) {
	for cursor := range vec.EachCoord(area) {
		c.DrawVisuals(cursor, depth, v)
	}
}

// draws a circle of radius r centered at (px, py), copying the visuals from v, with option to fill the circle with same
// visuals
func (c *Canvas) DrawCircle(center vec.Coord, depth, r int, v Visuals, fill bool) {
	drawFunc := func(pos vec.Coord) {
		c.DrawVisuals(pos, depth, v)
	}

	vec.Circle(center, r, drawFunc)

	if fill {
		c.FloodFill(center, depth, v)
	}
}

// Floodfill performs a floodfill starting at x,y. it fills with visuals v, also using v as criteria for looking for
// edges. any cell with a higher z value will also count as an edge and impede the flood
func (c *Canvas) FloodFill(pos vec.Coord, depth int, v Visuals) {
	//hey, write this function. it'll be fun i promise
}

// DrawToCanvas draws the canvas c to a destination canvas, offset by some Coord at depth z. This process will mark
// any copied cells in c as clean.
// TODO: this function should take in flags to determine how the canvas is copied
func (c *Canvas) DrawToCanvas(dst *Canvas, offset vec.Coord, depth int) {
	draw_area := vec.FindIntersectionRect(dst, vec.Rect{offset, c.Size()})
	for dst_cursor := range vec.EachCoord(draw_area) {
		if dst.getDepth(dst_cursor) > depth { //skip cell if it wouldn't be drawn to the destination canvas
			continue
		}
		src_cursor := dst_cursor.Subtract(offset)
		cell := c.getCell(src_cursor)
		if cell.Visuals.Mode != DRAW_NONE {
			dst.DrawVisuals(dst_cursor, depth, cell.Visuals)
		}
		cell.Dirty = false
	}

	c.dirty = false
}
