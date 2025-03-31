package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

type TextCellPosition uint8

const (
	DRAW_TEXT_LEFT  TextCellPosition = 0
	DRAW_TEXT_RIGHT TextCellPosition = 1
)

type TextMode uint8

const (
	TEXTMODE_DEFAULT TextMode = iota
	TEXTMODE_FULL
	TEXTMODE_HALF
)

var DefaultTextMode TextMode = TEXTMODE_FULL

// Draws text to the canvas, starting at pos and respecting depth. If drawing in half-width mode, start_pos determines
// which side of the cell the text begins in. Optionally takes a Textmode; if this is omitted, uses the defined default
// mode.
func (c *Canvas) DrawText(pos vec.Coord, depth int, text string, colours col.Pair, start_pos TextCellPosition, text_mode ...TextMode) {
	var mode TextMode
	if len(text_mode) == 0 {
		mode = DefaultTextMode
	} else {
		mode = text_mode[0]
		if mode == TEXTMODE_DEFAULT {
			mode = DefaultTextMode
		}
	}

	switch mode {
	case TEXTMODE_FULL:
		c.DrawFullWidthText(pos, depth, text, colours)
	case TEXTMODE_HALF:
		c.DrawHalfWidthText(pos, depth, text, colours, start_pos)
	default:
		log.Error("bad textmode????")
	}
}

// DrawText draws the provided string to the canvas using the half-width text drawing mode, beginning at pos and
// respecting depth. start_pos specifies which side of the cell we begin drawing in.
func (c *Canvas) DrawHalfWidthText(pos vec.Coord, depth int, text string, colours col.Pair, start_pos TextCellPosition) {
	//build []rune version of txt string
	var textRunes []rune = make([]rune, 0, len(text))
	if start_pos == DRAW_TEXT_RIGHT { //pad start with a space if we're starting on the right
		textRunes = append(textRunes, rune(TEXT_NONE))
	}
	textRunes = append(textRunes, []rune(text)...)
	if len(textRunes)%2 != 0 { //pad end if we're ending on the left
		textRunes = append(textRunes, rune(TEXT_NONE))
	}

	//iterate by pairs of runes, drawing 1 cell per loop
	for i := 0; i < len(textRunes); i += 2 {
		cursor := vec.Coord{pos.X + i/2, pos.Y}
		if !c.InBounds(cursor) { //make sure we're drawing in the canvas.
			continue
		}

		c.setCell(cursor, depth, NewTextVisuals(uint8(textRunes[i]), uint8(textRunes[i+1]), colours))
	}
}

// DrawFullWidthText draws the provided string to the canvas using the full-width glyph drawing mode, beginning at pos and
// respecting depth.
func (c *Canvas) DrawFullWidthText(pos vec.Coord, depth int, text string, colours col.Pair) {
	for i, textRune := range []rune(text) {
		c.setCell(pos.StepN(vec.DIR_RIGHT, i), depth, NewGlyphVisuals(Glyph(textRune), colours))
	}
}
