package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
)

//A single tile in a canvas. It tracks its Z value for use in a canvas, but does not enforce any z-depthing.
type Cell struct {
	Visuals
	Z     int
	Dirty bool //this will be true if the cell has been changed since the last copy/render operation
}

//Sets the properties of a cell all at once for Glyph Mode.
func (c *Cell) SetGlyphCell(gl int, fore, back uint32, z int) {
	c.SetGlyph(z, gl)
	c.SetForeColour(z, fore)
	c.SetBackColour(z, back)
}

//Sets the properties of a cell all at once for Text Mode.
func (c *Cell) SetTextCell(char1, char2 rune, fore, back uint32, z int) {
	c.SetText(z, char1, char2)
	c.SetForeColour(z, fore)
	c.SetBackColour(z, back)
}

func (c *Cell) SetForeColour(z int, col uint32) {
	if col != c.ForeColour {
		c.Z = z
		c.ForeColour = col
		c.Dirty = true
	}
}

func (c *Cell) SetBackColour(z int, col uint32) {
	if col != c.BackColour {
		c.Z = z
		c.BackColour = col
		c.Dirty = true
	}
}

func (c *Cell) SetGlyph(z, gl int) {
	if gl != c.Glyph || c.Mode != DRAW_GLYPH {
		c.Mode = DRAW_GLYPH
		c.Z = z
		c.Glyph = gl
		c.Dirty = true
	}
}

func (c *Cell) SetText(z int, char1, char2 rune) {
	if char1 != c.Chars[0] || char2 != c.Chars[1] || c.Mode != DRAW_TEXT {
		c.Mode = DRAW_TEXT
		c.Z = z
		if char1 != TEXT_DEFAULT {
			c.Chars[0] = char1
		}
		if char2 != TEXT_DEFAULT {
			c.Chars[1] = char2
		}
		c.Dirty = true
	}
}

//Re-inits a cell back to default blankness.
func (c *Cell) Clear() {
	if c.Mode == DRAW_GLYPH {
		c.SetGlyphCell(GLYPH_NONE, col.WHITE, col.BLACK, -1)
	} else {
		c.SetTextCell(TEXT_NONE, TEXT_NONE, col.WHITE, col.BLACK, -1)
	}
}
