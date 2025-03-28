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
	cells    []Cell
	depthmap []int
	dirty    bool //true if any cells in the Canvas are dirty and need to be drawn out. TODO: replace this with a dirty bitset

	//width, height int
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

func (c *Canvas) setForeColour(pos vec.Coord, depth int, colour uint32) {
	if c.getDepth(pos) > depth {
		return
	}

	if colour == COL_DEFAULT {
		colour = c.defaultVisuals.Colours.Fore
	}
	cell := c.getCell(pos)
	cell.SetForeColour(colour)
	c.dirty = c.dirty || cell.Dirty
	c.setDepth(pos, depth)
}

func (c *Canvas) setBackColour(pos vec.Coord, depth int, colour uint32) {
	if c.getDepth(pos) > depth {
		return
	}

	if colour == COL_DEFAULT {
		colour = c.defaultVisuals.Colours.Back
	}
	cell := c.getCell(pos)
	cell.SetBackColour(colour)
	c.dirty = c.dirty || cell.Dirty
	c.setDepth(pos, depth)
}

func (c *Canvas) setColours(pos vec.Coord, depth int, colours col.Pair) {
	c.setForeColour(pos, depth, colours.Fore)
	c.setBackColour(pos, depth, colours.Back)
}

func (c *Canvas) setGlyph(pos vec.Coord, depth int, glyph Glyph) {
	if c.getDepth(pos) > depth {
		return
	}

	cell := c.getCell(pos)
	cell.SetGlyph(glyph)
	c.dirty = c.dirty || cell.Dirty
	c.setDepth(pos, depth)
}

func (c *Canvas) setText(pos vec.Coord, depth int, char1, char2 uint8) {
	if c.getDepth(pos) > depth {
		return
	}

	cell := c.getCell(pos)
	cell.SetText(char1, char2)
	c.dirty = c.dirty || cell.Dirty
	c.setDepth(pos, depth)
}

// Changes a single character on the canvas at position (x,y) in text mode.
func (c *Canvas) setChar(pos vec.Coord, depth int, char uint8, char_pos TextCellPosition) {
	if c.getDepth(pos) > depth {
		return
	}

	cell := c.getCell(pos)
	cell.SetChar(char, char_pos)
	c.dirty = c.dirty || cell.Dirty
	c.setDepth(pos, depth)
}

// sets a cell at pos to DRAW_NONE
func (c *Canvas) setBlank(pos vec.Coord) {
	cell := c.getCell(pos)
	cell.SetBlank()
	c.dirty = c.dirty || cell.Dirty
	c.setDepth(pos, -1) // not sure if this makes sense...
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
				cell := c.getCell(cursor)
				cell.SetVisuals(c.defaultVisuals)
				c.setDepth(cursor, -1)
			}
		}
	}

	c.dirty = true
}

// reports whether the canvas should be drawn out
func (c Canvas) Dirty() bool {
	return c.dirty
}

func (c Canvas) DefaultColours() col.Pair {
	return c.defaultVisuals.Colours
}

// Returns a copy of a region of the canvas. If the area is not in the canvas, copy will be empty.
func (c Canvas) CopyArea(area vec.Rect) (copy Canvas) {
	copy.Init(area.Dims)
	copy.defaultVisuals = c.defaultVisuals

	if !vec.Intersects(c, area) {
		return
	}

	for cursor := range vec.EachCoordInArea(area) {
		if !c.InBounds(cursor) {
			continue
		}

		cell := c.getCell(cursor)
		depth := c.getDepth(cursor)
		copy_cursor := cursor.Subtract(area.Coord)
		copy.setColours(copy_cursor, depth, cell.Colours)
		switch cell.Mode {
		case DRAW_GLYPH:
			copy.setGlyph(copy_cursor, depth, cell.Glyph)
		case DRAW_TEXT:
			copy.setText(copy_cursor, depth, cell.Chars[0], cell.Chars[1])
		case DRAW_NONE:
			copy.setBlank(copy_cursor)
		}
	}

	return
}
