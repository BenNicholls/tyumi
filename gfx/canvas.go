package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

const (
	COL_DEFAULT uint32 = 0x00000001 //pass this in cases where you want the canvas to use the default colours
)

// Canvas is a Z-depthed grid of Cell objects.
// All canvas drawing options are z-depth sensitive. They will never draw a lower z value cell over a higher one.
// The clear function can be used to set a region of a canvas back to -1 z level so you can redraw over it.
type Canvas struct {
	cells    []Cell
	depthmap []int
	dirty    bool //true if any cells in the Canvas are dirty and need to be drawn out. TODO: replace this with a dirty bitset

	width, height int

	foreColour uint32 //default foreground colour
	backColour uint32 //default background colour
}

func (c *Canvas) Size() vec.Dims {
	return vec.Dims{c.width, c.height}
}

func (c *Canvas) Bounds() vec.Rect {
	return vec.Rect{vec.ZERO_COORD, vec.Dims{c.width, c.height}}
}

// Initializes the canvas. Can also be used for resizing, assuming you don't mind that the contents of the canvas
// are destroyed.
func (c *Canvas) Init(w, h int) {
	c.width, c.height = util.Abs(w), util.Abs(h)
	c.cells = make([]Cell, c.Size().Area())
	c.depthmap = make([]int, c.Size().Area())
	c.foreColour = col.WHITE
	c.backColour = col.BLACK
	c.Clear()
}

func (c *Canvas) InBounds(pos vec.Coord) bool {
	if pos.X >= c.width || pos.Y >= c.height || pos.X < 0 || pos.Y < 0 {
		return false
	}
	return true
}

// Clean sets all cells in the canvas as clean (dirty = false).
func (c *Canvas) Clean() {
	for i := range c.cells {
		c.cells[i].Dirty = false
	}
}

// Sets the default colours for a canvas, then does a reset of the canvas to apply them.
func (c *Canvas) SetDefaultColours(fore uint32, back uint32) {
	c.foreColour = fore
	c.backColour = back
	c.Clear()
}

// GetCell returns the cell at pos. Returns an empty cell if pos is out of bounds.
// Note that this function just returns the value of the requested cell, not a reference,
// so you can't change the cell this way. Use the Canvas.Draw* functions for that!
func (c *Canvas) GetCell(pos vec.Coord) (cell Cell) {
	if !c.InBounds(pos) {
		return
	}
	cell = c.cells[pos.Y*c.width+pos.X]
	return
}

func (c *Canvas) getCell(pos vec.Coord) *Cell {
	return &c.cells[pos.Y*c.width+pos.X]
}

func (c *Canvas) getDepth(pos vec.Coord) int {
	return c.depthmap[pos.Y*c.width+pos.X]
}

func (c *Canvas) setDepth(pos vec.Coord, depth int) {
	c.depthmap[pos.Y*c.width+pos.X] = depth
}

func (c *Canvas) setForeColour(pos vec.Coord, depth int, col uint32) {
	if c.getDepth(pos) <= depth {
		cell := c.getCell(pos)
		if col == COL_DEFAULT {
			col = c.foreColour
		}
		cell.SetForeColour(col)
		c.dirty = cell.Dirty
		c.setDepth(pos, depth)
	}
}

func (c *Canvas) setBackColour(pos vec.Coord, depth int, col uint32) {
	if c.getDepth(pos) <= depth {
		cell := c.getCell(pos)
		if col == COL_DEFAULT {
			col = c.backColour
		}
		cell.SetBackColour(col)
		c.dirty = cell.Dirty
		c.setDepth(pos, depth)
	}
}

func (c *Canvas) setColours(pos vec.Coord, depth int, fore, back uint32) {
	c.setForeColour(pos, depth, fore)
	c.setBackColour(pos, depth, back)
}

func (c *Canvas) setGlyph(pos vec.Coord, depth, gl int) {
	if c.getDepth(pos) <= depth {
		cell := c.getCell(pos)
		cell.SetGlyph(gl)
		c.dirty = cell.Dirty
		c.setDepth(pos, depth)
	}
}

func (c *Canvas) setText(pos vec.Coord, depth int, char1, char2 rune) {
	if c.getDepth(pos) <= depth {
		cell := c.getCell(pos)
		cell.SetText(char1, char2)
		c.dirty = cell.Dirty
		c.setDepth(pos, depth)
	}
}

// Changes a single character on the canvas at position (x,y) in text mode.
// charNum: 0 = Left, 1 = Right (for ease with modulo operations). Throw whatever in here though, it gets
// modulo 2'd anyways just in case.
func (c *Canvas) setChar(pos vec.Coord, depth int, char rune, char_pos TextCellPosition) {
	if c.getDepth(pos) <= depth {
		cell := c.getCell(pos)
		cell.Mode = DRAW_TEXT
		if cell.Chars[int(char_pos)] != char {
			cell.Chars[int(char_pos)] = char
			cell.Dirty = true
			c.setDepth(pos, depth)
			c.dirty = true
		}
	}
}

// Clear resets portions of the canvas. If no areas are provided, it resets the entire canvas.
func (c *Canvas) Clear(areas ...vec.Rect) {
	if len(areas) == 0 {
		areas = append(areas, c.Bounds())
	}

	for _, area := range areas {
		for cursor := range vec.EachCoord(area) {
			if !c.InBounds(cursor) { //need to check to make sure user-provided areas are in bounds.
				continue
			}
			cell := c.getCell(cursor)
			cell.Clear()
			cell.SetBackColour(c.backColour)
			cell.SetForeColour(c.foreColour)
			c.setDepth(cursor, -1)
		}
	}
	c.dirty = true
}

//reports whether the cavas should be drawn out
func (c Canvas) Dirty() bool {
	return c.dirty
}