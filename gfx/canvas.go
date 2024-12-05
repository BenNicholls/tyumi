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

// GetCell returns a reference to the cell at pos. Returns nil if pos is out of bounds.
func (c *Canvas) GetCell(pos vec.Coord) *Cell {
	if !c.InBounds(pos) {
		return nil
	}
	return &c.cells[pos.Y*c.width+pos.X]
}

// Gets the depth value at pos. If pos is not in the canvas returns -1
func (c *Canvas) GetDepth(pos vec.Coord) int {
	if !c.InBounds(pos) {
		return -1
	}
	return c.depthmap[pos.Y*c.width+pos.X]
}

// Sets the depth value at pos. If pos is not in canvas, does nothing.
func (c *Canvas) SetDepth(pos vec.Coord, depth int) {
	if c.InBounds(pos) {
		c.depthmap[pos.Y*c.width+pos.X] = depth
	}
}

func (c *Canvas) InBounds(pos vec.Coord) bool {
	if pos.X >= c.width || pos.Y >= c.height || pos.X < 0 || pos.Y < 0 {
		return false
	}
	return true
}

// Sets the default colours for a canvas, then does a reset of the canvas to apply them.
func (c *Canvas) SetDefaultColours(fore uint32, back uint32) {
	c.foreColour = fore
	c.backColour = back
	c.Clear()
}

func (c *Canvas) SetForeColour(pos vec.Coord, depth int, col uint32) {
	if cell := c.GetCell(pos); cell != nil && c.GetDepth(pos) <= depth {
		if col == COL_DEFAULT {
			col = c.foreColour
		}
		cell.SetForeColour(col)
		c.SetDepth(pos, depth)
	}
}

func (c *Canvas) SetBackColour(pos vec.Coord, depth int, col uint32) {
	if cell := c.GetCell(pos); cell != nil && c.GetDepth(pos) <= depth {
		if col == COL_DEFAULT {
			col = c.backColour
		}
		cell.SetBackColour(col)
		c.SetDepth(pos, depth)
	}
}

func (c *Canvas) SetColours(pos vec.Coord, depth int, fore, back uint32) {
	c.SetForeColour(pos, depth, fore)
	c.SetBackColour(pos, depth, back)
}

func (c *Canvas) SetGlyph(pos vec.Coord, depth, gl int) {
	if cell := c.GetCell(pos); cell != nil && c.GetDepth(pos) <= depth {
		cell.SetGlyph(gl)
		c.SetDepth(pos, depth)
	}
}

func (c *Canvas) SetText(pos vec.Coord, depth int, char1, char2 rune) {
	if cell := c.GetCell(pos); cell != nil && c.GetDepth(pos) <= depth {
		cell.SetText(char1, char2)
		c.SetDepth(pos, depth)
	}
}

// Changes a single character on the canvas at position (x,y) in text mode.
// charNum: 0 = Left, 1 = Right (for ease with modulo operations). Throw whatever in here though, it gets
// modulo 2'd anyways just in case.
func (c *Canvas) SetChar(pos vec.Coord, depth int, char rune, char_pos TextCellPosition) {
	if cell := c.GetCell(pos); cell != nil && c.GetDepth(pos) <= depth {
		cell.Mode = DRAW_TEXT
		if cell.Chars[int(char_pos)] != char {
			cell.Chars[int(char_pos)] = char
			c.SetDepth(pos, depth)
			cell.Dirty = true
		}
	}
}

// Clear resets portions of the canvas. If no areas are provided, it resets the entire canvas.
func (c *Canvas) Clear(areas ...vec.Rect) {
	if len(areas) == 0 {
		areas = append(areas, c.Bounds())
	}

	for _, r := range areas {
		for i := 0; i < r.Area(); i++ {
			pos := vec.Coord{r.X + i%r.W, r.Y + i/r.W}
			if cell := c.GetCell(pos); cell != nil {
				cell.Clear()
				cell.SetBackColour(c.backColour)
				cell.SetForeColour(c.foreColour)
				c.SetDepth(pos, -1)
			}
		}
	}
}
