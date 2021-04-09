package gfx

import (
	"github.com/bennicholls/tyumi/vec"
)

func (c *Canvas) Draw(x, y, z int, d Drawable) {
	c.DrawVisuals(x, y, z, d.Visuals())
}

func (c *Canvas) DrawVisuals(x, y, z int, v Visuals) {
	if cell := c.GetCell(x, y); cell != nil && cell.Z <= z {
		if v.Mode == DRAW_GLYPH {
			cell.SetGlyphCell(v.Glyph, v.ForeColour, v.BackColour, z)
		} else {
			cell.SetTextCell(v.Chars[0], v.Chars[1], v.ForeColour, v.BackColour, z)
		}
	}
}

func (c *Canvas) DrawText(x, y, z int, txt string, fore, back uint32, charNum int) {
	i := 0 //can't use the index from the range loop since it is counting bytes, not code-points
	for _, char := range txt {
		if vec.IsInside(x+(i+charNum)/2, y, c) {
			c.SetChar(x+(i+charNum)/2, y, z, char, (i+charNum)%2)
			c.SetColours(x+(i+charNum)/2, y, z, fore, back)
			if i == len(txt)-1 && (i+charNum)%2 == 0 {
				//if final character is in the left-side of a cell, blank the right side.
				c.SetChar(x+(i+charNum)/2, y, z, 32, 1)
			}
		}
		i++
	}
}

//draws a circle of radius r centered at (px, py), copying the visuals from v, with option to fill the circle with same
//visuals
func (c *Canvas) DrawCircle(px, py, z, r int, v Visuals, fill bool) {
	drawFunc := func(x, y int) {
		c.DrawVisuals(x, y, z, v)
	}

	vec.Circle(vec.Coord{px, py}, r, drawFunc)

	if fill {
		c.FloodFill(px, py, z, v)
	}
}

//Floodfill performs a floodfill starting at x,y. it fills with visuals v, also using v as criteria for looking for
//edges. any cell with a higher z value will also count as an edge and impede the flood
func (c *Canvas) FloodFill(x, y, z int, v Visuals) {
	//hey, write this function. it'll be fun i promise
}

//DrawToCanvas draws the canvas c to a destination canvas, offset by some (x, y) at depth z. This process will mark
//any copied cells in c as clean.
//TODO: this function should take in flags to determine how the canvas is copied
//      could also pass this a rect to indicate subaras of the canvas that need to be copied
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
