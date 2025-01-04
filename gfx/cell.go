package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
)

// A single tile in a canvas.
type Cell struct {
	Visuals
	Dirty bool //this will be true if the cell has been changed since the last copy/render operation
}

// Sets the properties of a cell all at once for Glyph Mode.
func (c *Cell) SetGlyphCell(gl int, colours col.Pair) {
	c.SetGlyph(gl)
	c.SetColours(colours)
}

// Sets the properties of a cell all at once for Text Mode.
func (c *Cell) SetTextCell(char1, char2 rune, colours col.Pair) {
	c.SetText(char1, char2)
	c.SetColours(colours)
}

func (c *Cell) SetForeColour(col uint32) {
	if col == c.Colours.Fore {
		return
	}
	
	c.Colours.Fore = col
	c.Dirty = true
}

func (c *Cell) SetBackColour(col uint32) {
	if col == c.Colours.Back {
		return
	}

	c.Colours.Back = col
	c.Dirty = true
}

func (c *Cell) SetColours(colours col.Pair) {
	if colours == c.Colours {
		return
	}

	c.Colours = colours
	c.Dirty = true
}

func (c *Cell) SetGlyph(gl int) {
	if gl == c.Glyph && c.Mode == DRAW_GLYPH {
		return
	}

	c.Mode = DRAW_GLYPH
	c.Glyph = gl
	c.Dirty = true
}

func (c *Cell) SetText(char1, char2 rune) {
	if char1 == c.Chars[0] && char2 == c.Chars[1] && c.Mode == DRAW_TEXT {
		return
	}

	c.Mode = DRAW_TEXT
	if char1 != TEXT_DEFAULT {
		c.Chars[0] = char1
	}
	if char2 != TEXT_DEFAULT {
		c.Chars[1] = char2
	}
	c.Dirty = true
}

func (c *Cell) SetChar(char rune, char_pos TextCellPosition) {
	if c.Chars[int(char_pos)] == char && c.Mode == DRAW_TEXT {
		return
	}

	c.Mode = DRAW_TEXT
	if char != TEXT_DEFAULT {
		c.Chars[int(char_pos)] = char
	}
	c.Dirty = true
}

// Re-inits a cell back to default blankness.
func (c *Cell) Clear() {
	if c.Mode == DRAW_GLYPH {
		c.SetGlyphCell(GLYPH_NONE, col.Pair{col.WHITE, col.BLACK})
	} else {
		c.SetTextCell(TEXT_NONE, TEXT_NONE, col.Pair{col.WHITE, col.BLACK})
	}
}
