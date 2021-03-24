package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/util"
)

const FIT_TEXT int = 0

type Textbox struct {
	ElementPrototype

	text   string //text to be displayed
	center bool
	lines  []string //text after it has been word wrapped
}

//Creates a textbox. You can set the width or height to FIT_TEXT to have the textbox compute the dimensions for you. If
//width is set to FIT_TEXT, the box will ensure the entire text fits on 1 line (aka height will be 1). Setting height =
//FIT_TEXT will wrap the text at the provided width, and the textbox will have height = however many lines are required.
//Note that this is just for initialization, the textbox won't change dimensions to fit later changes in the text.
func NewTextbox(w, h, x, y, z int, text string, center bool) Textbox {
	tb := Textbox{
		text:   text,
		center: center,
	}

	//auto-fit text if required
	if w == FIT_TEXT {
		h = 1
		w = (len(text) + 1) / 2
		tb.lines = make([]string, 1)
		tb.lines[0] = text
	} else if h == FIT_TEXT {
		tb.lines = util.WrapText(text, w*2)
		h = len(tb.lines)
	}

	tb.ElementPrototype.Init(w, h, x, y, z)
	return tb
}

func (tb *Textbox) ChangeText(txt string) {
	if txt == tb.text {
		return
	}

	tb.text = txt
	w, h := tb.Dims()
	tb.lines = util.WrapText(txt, w*2, h)
	tb.dirty = true
}

func (tb *Textbox) Render() {
	if tb.visible && tb.dirty {
		tb.Clear()
		for i, line := range tb.lines {
			tb.DrawText(0, i, 0, line, gfx.COL_DEFAULT, gfx.COL_DEFAULT, 0)
		}

		tb.ElementPrototype.Render()
	}
}
