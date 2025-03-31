package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

const (
	COL_DEFAULT col.Colour = 0x00000001 //pass this in cases where you want a canvas to use the default colours
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

// DrawableReporter defines drawable types that can report when they are dirty and need to be redrawn
type DrawableReporter interface {
	Drawable
	Dirty() bool
}

// Defines anything that can report a set of visuals for drawing to a single cell.
type VisualObject interface {
	GetVisuals() Visuals
}

// Draw draws the canvas c to a destination canvas, offset by some Coord at depth z. This process will mark
// any copied cells in c as clean.
// TODO: this function should take in flags to determine how the canvas is copied
func (c *Canvas) Draw(dst_canvas *Canvas, offset vec.Coord, depth int) {
	for dstCursor := range vec.EachCoordInIntersection(dst_canvas, c.Bounds().Translated(offset)) {
		srcCursor := dstCursor.Subtract(offset)
		cell := c.getCell(srcCursor)
		cell.Dirty = false

		//draw cell if depth is higher, or if the cell in the destination canvas is DRAW_NONE
		if dst_canvas.getDepth(dstCursor) <= depth || dst_canvas.getCell(dstCursor).Mode == DRAW_NONE {
			if cell.Visuals.Mode != DRAW_NONE {
				dst_canvas.setCell(dstCursor, depth, cell.Visuals)
			}
		}
	}

	c.dirty = false
}

func (c *Canvas) DrawVisuals(pos vec.Coord, depth int, visuals Visuals) {
	if !c.InBounds(pos) {
		return
	}

	c.setCell(pos, depth, visuals)
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
