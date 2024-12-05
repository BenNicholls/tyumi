package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
)

//A single tile in a canvas.
type Cell struct {
	Visuals
	Dirty bool //this will be true if the cell has been changed since the last copy/render operation
}

//Sets the properties of a cell all at once for Glyph Mode.
func (c *Cell) SetGlyphCell(gl int, fore, back uint32) {
	c.SetGlyph(gl)
	c.SetForeColour(fore)
	c.SetBackColour(back)
}

//Sets the properties of a cell all at once for Text Mode.
func (c *Cell) SetTextCell(char1, char2 rune, fore, back uint32) {
	c.SetText(char1, char2)
	c.SetForeColour(fore)
	c.SetBackColour(back)
}

func (c *Cell) SetForeColour(col uint32) {
	if col != c.ForeColour {
		c.ForeColour = col
		c.Dirty = true
	}
}

func (c *Cell) SetBackColour(col uint32) {
	if col != c.BackColour {
		c.BackColour = col
		c.Dirty = true
	}
}

func (c *Cell) SetGlyph(gl int) {
	if gl != c.Glyph || c.Mode != DRAW_GLYPH {
		c.Mode = DRAW_GLYPH
		c.Glyph = gl
		c.Dirty = true
	}
}

func (c *Cell) SetText(char1, char2 rune) {
	if char1 != c.Chars[0] || char2 != c.Chars[1] || c.Mode != DRAW_TEXT {
		c.Mode = DRAW_TEXT
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
		c.SetGlyphCell(GLYPH_NONE, col.WHITE, col.BLACK)
	} else {
		c.SetTextCell(TEXT_NONE, TEXT_NONE, col.WHITE, col.BLACK)
	}
}
