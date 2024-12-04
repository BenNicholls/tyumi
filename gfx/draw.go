package gfx

import (
	"github.com/bennicholls/tyumi/vec"
)

type TextCellPosition int

const (
	DRAW_TEXT_LEFT  TextCellPosition = 0
	DRAW_TEXT_RIGHT TextCellPosition = 1
)

func (c *Canvas) Draw(x, y, z int, d Drawable) {
	c.DrawVisuals(x, y, z, d.Visuals())
}

//TODO: fix default colour detection
func (c *Canvas) DrawVisuals(x, y, z int, v Visuals) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		if v.Mode == DRAW_GLYPH {
			cell.SetGlyphCell(v.Glyph, v.ForeColour, v.BackColour, z)
		} else {
			cell.SetTextCell(v.Chars[0], v.Chars[1], v.ForeColour, v.BackColour, z)
		}
	}
}

func (c *Canvas) DrawText(x, y, z int, txt string, fore, back uint32, start_pos TextCellPosition) {
	//build []rune version of txt string
	var text_runes []rune = make([]rune, 0, len(txt))
	if start_pos == DRAW_TEXT_RIGHT { //pad start with a space if we're starting on the right
		text_runes = append(text_runes, TEXT_NONE)
	}
	for _, r := range txt {
		text_runes = append(text_runes, r)
	}
	if len(text_runes)%2 != 0 { //pad end if we're ending on the left
		text_runes = append(text_runes, TEXT_NONE)
	}

	//iterate by pairs of runes, drawing 1 cell per loop
	for i := 0; i < len(text_runes); i += 2 {
		pos := vec.Coord{x + i/2, y}
		cell := c.GetCell(pos.X, pos.Y)
		if cell == nil { //make sure we're drawing in the canvas. TODO: some kind of easy bounds check thing??
			continue
		}

		c.SetText(pos.X, pos.Y, z, text_runes[i], text_runes[i+1])
		c.SetColours(pos.X, pos.Y, z, fore, back)
	}
}

// draws a circle of radius r centered at (px, py), copying the visuals from v, with option to fill the circle with same
// visuals
func (c *Canvas) DrawCircle(px, py, z, r int, v Visuals, fill bool) {
	drawFunc := func(x, y int) {
		c.DrawVisuals(x, y, z, v)
	}

	vec.Circle(vec.Coord{px, py}, r, drawFunc)

	if fill {
		c.FloodFill(px, py, z, v)
	}
}

// Floodfill performs a floodfill starting at x,y. it fills with visuals v, also using v as criteria for looking for
// edges. any cell with a higher z value will also count as an edge and impede the flood
func (c *Canvas) FloodFill(x, y, z int, v Visuals) {
	//hey, write this function. it'll be fun i promise
}

// DrawToCanvas draws the canvas c to a destination canvas, offset by some (x, y) at depth z. This process will mark
// any copied cells in c as clean.
// TODO: this function should take in flags to determine how the canvas is copied
//
//	could also pass this a rect to indicate subaras of the canvas that need to be copied
func (c *Canvas) DrawToCanvas(dst *Canvas, x, y, z int) {
	for i := range c.cells {
		dx, dy := x+i%c.width, y+i/c.width
		cell := c.GetCell(i%c.width, i/c.width)
		if dcell := dst.GetCell(dx, dy); dcell != nil && z >= dcell.Z {
			if cell.Mode == DRAW_GLYPH {
				dcell.SetGlyphCell(cell.Glyph, cell.ForeColour, cell.BackColour, z)
			} else {
				dcell.SetTextCell(cell.Chars[0], cell.Chars[1], cell.ForeColour, cell.BackColour, z)
			}

			cell.Dirty = false
		}
	}
}
