package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

type TextCellPosition int

const (
	DRAW_TEXT_LEFT  TextCellPosition = 0
	DRAW_TEXT_RIGHT TextCellPosition = 1
)

// DrawText draws the provided string to the canvas using the half-width text drawing mode, beginning at pos and
// respecting depth. start_pos specifies which side of the cell we begin drawing in.
func (c *Canvas) DrawText(pos vec.Coord, depth int, txt string, colours col.Pair, start_pos TextCellPosition) {
	//build []rune version of txt string
	var text_runes []rune = make([]rune, 0, len(txt))
	if start_pos == DRAW_TEXT_RIGHT { //pad start with a space if we're starting on the right
		text_runes = append(text_runes, TEXT_NONE)
	}
	text_runes = append(text_runes, []rune(txt)...)
	if len(text_runes)%2 != 0 { //pad end if we're ending on the left
		text_runes = append(text_runes, TEXT_NONE)
	}

	//iterate by pairs of runes, drawing 1 cell per loop
	for i := 0; i < len(text_runes); i += 2 {
		cursor := vec.Coord{pos.X + i/2, pos.Y}
		if !c.InBounds(cursor) { //make sure we're drawing in the canvas.
			continue
		}

		c.setText(cursor, depth, text_runes[i], text_runes[i+1])
		c.setColours(cursor, depth, colours)
	}
}

// DrawFullText draws the provided string to the canvas using the full-width glyph drawing mode, beginning at pos and 
// respecting depth.
// TODO: write this function you dweebo.
func (c *Canvas) DrawFullText(pos vec.Coord, depth int, txt string, colours col.Pair) {

}