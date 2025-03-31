package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
)

// A single tile in a canvas.
type Cell struct {
	Visuals

	Dirty bool //this will be true if the cell has been changed since the last time its canvas has been drawn out
}

// Sets the properties of a cell all at once for Glyph Mode.
func (c *Cell) SetGlyphCell(glyph Glyph, colours col.Pair) {
	c.SetGlyph(glyph)
	c.SetColours(colours)
}

// Sets the properties of a cell all at once for Text Mode.
func (c *Cell) SetTextCell(char1, char2 uint8, colours col.Pair) {
	c.SetText(char1, char2)
	c.SetColours(colours)
}

func (c *Cell) SetForeColour(colour col.Colour) {
	if colour == c.Colours.Fore || colour == col.NONE {
		return
	}

	c.Colours.Fore = colour
	if c.Mode != DRAW_NONE {
		c.Dirty = true
	}
}

func (c *Cell) SetBackColour(colour col.Colour) {
	if colour == c.Colours.Back || colour == col.NONE {
		return
	}

	c.Colours.Back = colour
	if c.Mode != DRAW_NONE {
		c.Dirty = true
	}
}

func (c *Cell) SetColours(colours col.Pair) {
	if colours == c.Colours {
		return
	}

	c.Colours = colours
	if c.Mode != DRAW_NONE {
		c.Dirty = true
	}
}

func (c *Cell) SetGlyph(glyph Glyph) {
	if glyph == c.Glyph && c.Mode == DRAW_GLYPH {
		return
	}

	c.Mode = DRAW_GLYPH
	c.Glyph = glyph
	c.Dirty = true
}

func (c *Cell) SetText(char1, char2 uint8) {
	if char1 == TEXT_DEFAULT {
		char1 = c.Chars[0]
	}

	if char2 == TEXT_DEFAULT {
		char2 = c.Chars[1]
	}

	if char1 == c.Chars[0] && char2 == c.Chars[1] && c.Mode == DRAW_TEXT {
		return
	}

	c.Mode = DRAW_TEXT
	c.Chars[0], c.Chars[1] = char1, char2
	c.Dirty = true
}

func (c *Cell) SetChar(char uint8, char_pos TextCellPosition) {
	if c.Chars[int(char_pos)] == char && c.Mode == DRAW_TEXT {
		return
	}

	c.Mode = DRAW_TEXT
	if char != TEXT_DEFAULT {
		c.Chars[int(char_pos)] = char
	}
	c.Dirty = true
}

func (c *Cell) SetBlank() {
	if c.Mode == DRAW_NONE {
		return
	}

	c.Mode = DRAW_NONE
	c.Glyph = GLYPH_NONE
	c.Chars = [2]uint8{0, 0}
	c.Colours = col.Pair{col.WHITE, col.BLACK}
	c.Dirty = true
}

// Sets the cell's Visuals to vis all at once
func (c *Cell) SetVisuals(visuals Visuals) {
	if c.Visuals == visuals {
		return
	}

	c.Visuals = visuals
	c.Dirty = true
}

// Re-inits a cell back to default blankness.
func (c *Cell) Clear() {
	switch c.Mode {
	case DRAW_GLYPH:
		c.SetGlyphCell(GLYPH_NONE, col.Pair{col.WHITE, col.BLACK})
	case DRAW_TEXT:
		c.SetTextCell(TEXT_NONE, TEXT_NONE, col.Pair{col.WHITE, col.BLACK})
	case DRAW_NONE: //draw_none cells are assumed to have been already cleared
		return
	}
}
