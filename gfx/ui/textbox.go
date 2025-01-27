package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

const FIT_TEXT int = 0

type Textbox struct {
	ElementPrototype

	text   string //text to be displayed
	center bool
	lines  []string //text after it has been word wrapped
}

// Creates a textbox. You can set the width or height to FIT_TEXT to have the textbox compute the dimensions for you. If
// width is set to FIT_TEXT, the box will ensure the entire text fits on 1 line (aka height will be 1). Setting height =
// FIT_TEXT will wrap the text at the provided width, and the textbox will have height = however many lines are required.
// Note that this is just for initialization, the textbox won't change dimensions to fit later changes in the text.
func NewTextbox(w, h int, pos vec.Coord, depth int, text string, center bool) (tb *Textbox) {
	tb = new(Textbox)
	tb.Init(w, h, pos, depth, text, center)

	return tb
}

func (tb *Textbox) Init(w, h int, pos vec.Coord, depth int, text string, center bool) {
	tb.text = text
	tb.center = center

	//auto-fit text if required
	if w == FIT_TEXT {
		w, h = (len(text)+1)/2, 1
		tb.lines = make([]string, 1)
		tb.lines[0] = text
	} else if h == FIT_TEXT {
		tb.lines = util.WrapText(text, w*2)
		h = len(tb.lines)
	} else {
		tb.lines = util.WrapText(text, w*2)
	}

	tb.ElementPrototype.Init(w, h, pos, depth)
}

func (tb *Textbox) ChangeText(txt string) {
	if txt == tb.text {
		return
	}

	tb.text = txt
	tb.lines = util.WrapText(txt, tb.Size().W*2, tb.Size().H)
	tb.Updated = true
}

func (tb *Textbox) Render() {
	tb.ClearAtDepth(0)
	for i, line := range tb.lines {
		x_offset := 0
		if tb.center {
			x_offset = (tb.Bounds().W*2 - len(line)) / 2
		}
		pos := vec.Coord{x_offset / 2, i}
		tb.DrawText(pos, 0, line, col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT}, gfx.TextCellPosition(x_offset%2))
	}
}
