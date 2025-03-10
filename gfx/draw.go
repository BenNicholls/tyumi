package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

const (
	COL_DEFAULT uint32 = 0x00000001 //pass this in cases where you want a canvas to use the default colours
)

type DrawMode uint8

const (
	DRAW_GLYPH DrawMode = iota // draws cell in glyph mode (square font)
	DRAW_TEXT                  // draws cell in text mode (half-width font)
	DRAW_NONE                  //do not draw this cell
)

// Defines anything with the ability to be drawn to a canvas
type Drawable interface {
	Draw(dst_canvas *Canvas, offset vec.Coord, depth int)
}

// Defines anything that can report a set of visuals for drawing to a single cell.
type VisualObject interface {
	GetVisuals() Visuals
}

// Draw draws the canvas c to a destination canvas, offset by some Coord at depth z. This process will mark
// any copied cells in c as clean.
// TODO: this function should take in flags to determine how the canvas is copied
func (c *Canvas) Draw(dst_canvas *Canvas, offset vec.Coord, depth int) {
	drawArea := vec.FindIntersectionRect(dst_canvas, c.Bounds().Translated(offset))
	for dstCursor := range vec.EachCoordInArea(drawArea) {
		if dst_canvas.getDepth(dstCursor) > depth { //skip cell if it wouldn't be drawn to the destination canvas
			continue
		}
		srcCursor := dstCursor.Subtract(offset)
		cell := c.getCell(srcCursor)
		if cell.Visuals.Mode != DRAW_NONE {
			dst_canvas.DrawVisuals(dstCursor, depth, cell.Visuals)
		}
		cell.Dirty = false
	}

	c.dirty = false
}

// THINK: this checks/sets the depth 3-4 times i think. hmmm.
func (c *Canvas) DrawVisuals(pos vec.Coord, depth int, visuals Visuals) {
	if !c.InBounds(pos) {
		return
	}

	switch visuals.Mode {
	case DRAW_GLYPH:
		c.setGlyph(pos, depth, visuals.Glyph)
	case DRAW_TEXT:
		c.setText(pos, depth, visuals.Chars[0], visuals.Chars[1])
	case DRAW_NONE:
		c.setBlank(pos)
		return // if we are not drawing this cell we can skip setting the colours below
	}

	c.setColours(pos, depth, visuals.Colours)
}

// Draws a single-celled object to the canvas.
func (c *Canvas) DrawObject(pos vec.Coord, depth int, object VisualObject) {
	c.DrawVisuals(pos, depth, object.GetVisuals())
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

// DrawGlyph draws a glyph to a cell at pos, respecting depth.
func (c *Canvas) DrawGlyph(pos vec.Coord, depth int, glyph Glyph) {
	if !c.InBounds(pos) {
		return
	}

	c.setGlyph(pos, depth, glyph)
}

// Floodfill performs a floodfill starting at x,y. it fills with visuals v, also using v as criteria for looking for
// edges. any cell with a higher z value will also count as an edge and impede the flood
// TODO: write this function
func (c *Canvas) FloodFill(pos vec.Coord, depth int, visuals Visuals) {
	//hey, write this function. it'll be fun i promise
}
