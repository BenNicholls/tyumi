package gfx

import (
	"fmt"

	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

// Canvas is a Z-depthed grid of Cell objects.
// All canvas drawing options are z-depth sensitive. They will never draw a lower z value cell over a higher one.
// The clear function can be used to set a region of a canvas back to -1 z level so you can redraw over it.
type Canvas struct {
	cells            []Cell
	depthmap         []int
	dirty            bool //true if any cells in the Canvas are dirty and need to be drawn out. TODO: replace this with a dirty bitset
	transparentCells int  //number of cells with some sort of transparency. if >0, the whole canvas is reported as transparent

	size   vec.Dims
	offset vec.Coord //coordinate of the top-left corner. generally (0,0)

	defaultVisuals Visuals // Visuals drawn when the canvas is cleared.
}

// Initializes the canvas, setting all cells to a nice black and white default drawing mode.
func (c *Canvas) Init(size vec.Dims) {
	if size.Area() == 0 {
		log.Error("Canvas cannot be initialized: size has zero area.")
		return
	}

	c.defaultVisuals = Visuals{
		Mode:    DRAW_GLYPH,
		Colours: col.Pair{col.WHITE, col.BLACK},
	}

	c.Resize(size)
}

// Returns true if the canvas has been initialized and is non-zero in size
func (c Canvas) Ready() bool {
	return c.cells != nil && len(c.cells) != 0
}

func (c Canvas) String() string {
	return fmt.Sprintf("Canvas with size %v, offset %v and default visuals %v", c.size, c.offset, c.defaultVisuals)
}

// Sets the default colours for a canvas, then does a reset of the canvas to apply them.
func (c *Canvas) SetDefaultColours(colours col.Pair) {
	if c.defaultVisuals.Colours == colours {
		return
	}

	c.defaultVisuals.Colours = colours
	c.Clear()
}

// Sets the default visuals for a canvas, then does a reset of the canvas to apply them.
func (c *Canvas) SetDefaultVisuals(visuals Visuals) {
	if c.defaultVisuals == visuals {
		return
	}

	c.defaultVisuals = visuals
	c.Clear()
}

func (c *Canvas) Size() vec.Dims {
	return c.size
}

// Resizes the canvas. This also clears the canvas!
func (c *Canvas) Resize(new_size vec.Dims) {
	if c.size == new_size {
		return
	}

	c.size = new_size
	c.cells = make([]Cell, c.size.Area())
	c.depthmap = make([]int, c.size.Area())
	c.transparentCells = c.size.Area()
	c.Clear()
}

func (c Canvas) Bounds() vec.Rect {
	return vec.Rect{c.offset, c.size}
}

func (c *Canvas) InBounds(pos vec.Coord) bool {
	return pos.IsInside(c)
}

// SetOrigin sets the origin coord for the canvas. draw operations will be done relative to this point.
// must be a point in the canvas, so {0 <= x < W, 0 <= y < H}
func (c *Canvas) SetOrigin(pos vec.Coord) {
	if !pos.IsInside(c.size) {
		return
	}

	c.offset.X, c.offset.Y = -pos.X, -pos.Y
}

// Clean sets all cells in the canvas as clean (dirty = false).
func (c *Canvas) Clean() {
	for i := range c.cells {
		c.cells[i].Dirty = false
	}
	c.dirty = false
}

func (c *Canvas) cellIndex(pos vec.Coord) int {
	if c.offset == vec.ZERO_COORD {
		return pos.ToIndex(c.size.W)
	}

	return pos.Subtract(c.offset).ToIndex(c.size.W)
}

// GetCell returns the cell at pos. Returns an empty cell if pos is out of bounds.
// Note that this function just returns the value of the requested cell, not a reference,
// so you can't change the cell this way. Use the Canvas.Draw* functions for that!
func (c *Canvas) GetCell(pos vec.Coord) (cell Cell) {
	if !c.InBounds(pos) {
		return
	}

	cell = c.cells[c.cellIndex(pos)]
	return
}

func (c *Canvas) getCell(pos vec.Coord) *Cell {
	return &c.cells[c.cellIndex(pos)]
}

func (c *Canvas) GetDepth(pos vec.Coord) int {
	if !c.InBounds(pos) {
		panic("bad depth get! do a bounds check first!!!")
	}

	return c.getDepth(pos)
}

func (c *Canvas) getDepth(pos vec.Coord) int {
	return c.depthmap[c.cellIndex(pos)]
}

func (c *Canvas) setDepth(pos vec.Coord, depth int) {
	c.depthmap[c.cellIndex(pos)] = depth
}

// sets the visuals for the cell, respecting depth. if depth = -1, previous depth value is ignored.
func (c *Canvas) setCell(pos vec.Coord, depth int, vis Visuals) {
	if c.getDepth(pos) > depth && depth != -1 {
		return
	}

	c.setDepth(pos, depth)

	cell := c.getCell(pos)
	if vis.Colours.Fore == COL_DEFAULT {
		vis.Colours.Fore = c.defaultVisuals.Colours.Fore
	} else if vis.Colours.Fore == col.NONE {
		vis.Colours.Fore = cell.Colours.Fore
	}
	if vis.Colours.Back == COL_DEFAULT {
		vis.Colours.Back = c.defaultVisuals.Colours.Back
	} else if vis.Colours.Back == col.NONE {
		vis.Colours.Back = cell.Colours.Back
	}
	if vis.Mode == DRAW_TEXT {
		if vis.Chars[0] == TEXT_DEFAULT {
			vis.Chars[0] = cell.Chars[0]
		}
		if vis.Chars[1] == TEXT_DEFAULT {
			vis.Chars[1] = cell.Chars[1]
		}
	}
	if cell.Visuals == vis {
		return
	}

	transparent_before := cell.IsTransparent()
	cell.Visuals = vis
	cell.Dirty = true
	if transparent_after := cell.IsTransparent(); transparent_after != transparent_before {
		if transparent_after {
			c.transparentCells += 1
		} else {
			c.transparentCells -= 1
		}
	}
	c.dirty = true
}

func (c *Canvas) setForeColour(pos vec.Coord, depth int, colour col.Colour) {
	v := c.getCell(pos).Visuals
	v.Colours.Fore = colour
	c.setCell(pos, depth, v)
}

func (c *Canvas) setBackColour(pos vec.Coord, depth int, colour col.Colour) {
	v := c.getCell(pos).Visuals
	v.Colours.Back = colour
	c.setCell(pos, depth, v)
}

func (c *Canvas) setColours(pos vec.Coord, depth int, colours col.Pair) {
	v := c.getCell(pos).Visuals
	v.Colours = colours
	if v.Mode == DRAW_NONE && colours.Back != col.NONE {
		v.Mode = DRAW_GLYPH
	}
	c.setCell(pos, depth, v)
}

func (c *Canvas) setGlyph(pos vec.Coord, depth int, glyph Glyph) {
	v := c.getCell(pos).Visuals
	v.Mode = DRAW_GLYPH
	v.Glyph = glyph
	c.setCell(pos, depth, v)
}

func (c *Canvas) setText(pos vec.Coord, depth int, char1, char2 uint8) {
	v := c.getCell(pos).Visuals
	v.Mode = DRAW_TEXT
	v.Chars[0], v.Chars[1] = char1, char2
	c.setCell(pos, depth, v)
}

// Changes a single character on the canvas at position (x,y) in text mode.
func (c *Canvas) setChar(pos vec.Coord, depth int, char uint8, char_pos TextCellPosition) {
	v := c.getCell(pos).Visuals
	v.Mode = DRAW_TEXT
	switch char_pos {
	case DRAW_TEXT_LEFT:
		v.Chars[0] = char
	case DRAW_TEXT_RIGHT:
		v.Chars[1] = char
	}
	c.setCell(pos, depth, v)
}

// sets a cell at pos to DRAW_NONE
func (c *Canvas) setBlank(pos vec.Coord) {
	c.setCell(pos, -1, Visuals{Mode: DRAW_NONE})
}

// Clear resets portions of the canvas. If no areas are provided, it resets the entire canvas. The appearance
// of the reset cells can set using canvas.SetDefaultVisuals()
func (c *Canvas) Clear(areas ...vec.Rect) {
	c.ClearAtDepth(-1, areas...)
}

// Clears all cells in the canvas at or below a certain depth. If depth < 0, clears everything
func (c *Canvas) ClearAtDepth(depth int, areas ...vec.Rect) {
	if len(areas) == 0 {
		areas = append(areas, c.Bounds())
	}

	for _, area := range areas {
		for cursor := range vec.EachCoordInIntersection(c, area) {
			if depth < 0 || c.getDepth(cursor) <= depth {
				c.setCell(cursor, -1, c.defaultVisuals)
			}
		}
	}

	c.dirty = true
}

// reports whether the canvas should be drawn out
func (c Canvas) Dirty() bool {
	return c.dirty
}

func (c Canvas) IsTransparent() bool {
	return c.transparentCells > 0
}

func (c Canvas) DefaultColours() col.Pair {
	return c.defaultVisuals.Colours
}

func (c Canvas) DefaultVisuals() Visuals {
	return c.defaultVisuals
}

// Returns a copy of a region of the canvas. If the area is not in the canvas, copy will be empty.
func (c Canvas) CopyArea(area vec.Rect) (copy Canvas) {
	copy.Init(area.Dims)

	if !vec.Intersects(c, area) {
		return
	}

	copy.SetDefaultVisuals(c.defaultVisuals)
	for cursor := range vec.EachCoordInIntersection(c, area) {
		copy.setCell(cursor.Subtract(area.Coord), c.getDepth(cursor), c.getCell(cursor).Visuals)
	}

	return
}
