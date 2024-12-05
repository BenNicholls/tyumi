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
	cells []Cell

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
	c.cells = make([]Cell, c.width*c.height)
	c.foreColour = col.WHITE
	c.backColour = col.BLACK
	c.Clear()
}

// GetCell returns a reference to the cell at (x, y). Returns nil if (x,y) is out of bounds.
func (c *Canvas) GetCell(x, y int) *Cell {
	if x < 0 || y < 0 || x >= c.width || y >= c.height{
		return nil
	}
	return &c.cells[y*c.width+x]
}

// Sets the default colours for a canvas, then does a reset of the canvas to apply them.
func (c *Canvas) SetDefaultColours(fore uint32, back uint32) {
	c.foreColour = fore
	c.backColour = back
	c.Clear()
}

func (c *Canvas) SetForeColour(x, y, z int, col uint32) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		if col == COL_DEFAULT {
			col = c.foreColour
		}
		cell.SetForeColour(z, col)
	}
}

func (c *Canvas) SetBackColour(x, y, z int, col uint32) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		if col == COL_DEFAULT {
			col = c.backColour
		}
		cell.SetBackColour(z, col)
	}
}

func (c *Canvas) SetColours(x, y, z int, fore, back uint32) {
	c.SetForeColour(x, y, z, fore)
	c.SetBackColour(x, y, z, back)
}

func (c *Canvas) SetGlyph(x, y, z, gl int) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		cell.SetGlyph(z, gl)
	}
}

func (c *Canvas) SetText(x, y, z int, char1, char2 rune) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		cell.SetText(z, char1, char2)
	}
}

// Changes a single character on the canvas at position (x,y) in text mode.
// charNum: 0 = Left, 1 = Right (for ease with modulo operations). Throw whatever in here though, it gets
// modulo 2'd anyways just in case.
func (c *Canvas) SetChar(x, y, z int, char rune, charNum int) {
	if cell := c.GetCell(x, y); cell != nil && charNum >= 0 && cell.Z <= z {
		cell.Mode = DRAW_TEXT
		if cell.Chars[charNum%2] != char {
			cell.Chars[charNum%2] = char
			cell.Z = z
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
		for i := 0; i < r.W*r.H; i++ {
			if cell := c.GetCell(r.X+i%r.W, r.Y+i/r.W); cell != nil {
				cell.Clear()
				cell.SetBackColour(-1, c.backColour)
				cell.SetForeColour(-1, c.foreColour)
			}
		}
	}
}
