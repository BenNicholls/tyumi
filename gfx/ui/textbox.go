package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

const FIT_TEXT int = 50000

// Determines how to align text
type Alignment uint8

const (
	ALIGN_LEFT Alignment = iota
	ALIGN_CENTER
	ALIGN_RIGHT
)

type Textbox struct {
	Element

	fit_width  bool
	fit_height bool
	alignment  Alignment
	textMode   gfx.TextMode
	text       string   //text to be displayed
	lines      []string //text after it has been word wrapped
}

// Creates a textbox. You can set the width or height to FIT_TEXT to have the textbox compute the dimensions for you. If
// width is set to FIT_TEXT, the box will ensure the entire text fits on 1 line (aka height will be 1). Setting height =
// FIT_TEXT will wrap the text at the provided width, and the textbox will have height = however many lines are required.
func NewTextbox(size vec.Dims, pos vec.Coord, depth int, text string, align Alignment) (tb *Textbox) {
	tb = new(Textbox)
	tb.Init(size, pos, depth, text, align)

	return
}

// Creates a textbox, defaulting to centered text and enabling the default border.
func NewTitleTextbox(size vec.Dims, pos vec.Coord, depth int, text string) (tb *Textbox) {
	tb = new(Textbox)
	tb.InitTitle(size, pos, depth, text)

	return
}

func (tb *Textbox) Init(size vec.Dims, pos vec.Coord, depth int, text string, align Alignment) {
	tb.text = text
	tb.alignment = align
	tb.lines = make([]string, 0)

	if size.W == FIT_TEXT {
		tb.fit_width = true
		tb.Element.Init(vec.Dims{1, 1}, pos, depth)
	} else if size.H == FIT_TEXT {
		tb.Element.Init(vec.Dims{size.W, 1}, pos, depth)
		tb.fit_height = true
	} else {
		tb.Element.Init(size, pos, depth)
	}

	tb.TreeNode.Init(tb)
	tb.wrapText()
}

func (tb *Textbox) InitTitle(size vec.Dims, pos vec.Coord, depth int, text string) {
	tb.Init(size, pos, depth, text, ALIGN_CENTER)
	tb.EnableBorder()
}

func (tb *Textbox) SetTextMode(text_mode gfx.TextMode) {
	if tb.textMode == text_mode {
		return
	}

	tb.textMode = text_mode
	tb.wrapText()
	tb.Updated = true
}

func (tb *Textbox) ChangeText(txt string) {
	if txt == tb.text {
		return
	}

	tb.text = txt
	tb.wrapText()
	tb.Updated = true
}

func (tb *Textbox) AppendText(txt string) {
	if txt == "" {
		return
	}

	tb.ChangeText(tb.text + txt)
}

func (tb *Textbox) wrapText() {
	size := tb.size
	switch tb.getTextMode() {
	case gfx.TEXTMODE_FULL:
		if tb.fit_width {
			size = vec.Dims{len(tb.text), 1}
			tb.lines = make([]string, 1)
			tb.lines[0] = tb.text
		} else if tb.fit_height {
			tb.lines = util.WrapText(tb.text, tb.size.W)
			size.H = len(tb.lines)
		} else {
			tb.lines = util.WrapText(tb.text, tb.size.W, tb.size.H)
		}
	case gfx.TEXTMODE_HALF:
		if tb.fit_width {
			size = vec.Dims{(len(tb.text) + 1) / 2, 1}
			tb.lines = make([]string, 1)
			tb.lines[0] = tb.text
		} else if tb.fit_height {
			tb.lines = util.WrapText(tb.text, size.W*2)
			size.H = len(tb.lines)
		} else {
			tb.lines = util.WrapText(tb.text, tb.size.W*2, tb.size.H)
		}
	}

	if size != tb.size {
		tb.Resize(size)
	}
}

func (tb Textbox) getTextMode() gfx.TextMode {
	if tb.textMode == gfx.TEXTMODE_DEFAULT {
		return gfx.DefaultTextMode
	}

	return tb.textMode
}

func (tb *Textbox) Render() {
	tb.ClearAtDepth(0)
	for i, line := range tb.lines {
		x_offset := 0
		pos := vec.Coord{0, i}
		switch tb.alignment {
		case ALIGN_CENTER:
			switch tb.getTextMode() {
			case gfx.TEXTMODE_FULL:
				x_offset = (tb.size.W - len(line)) / 2
				pos.X = x_offset
			case gfx.TEXTMODE_HALF:
				x_offset = (tb.size.W*2 - len(line)) / 2
				pos.X = x_offset / 2
			}
		case ALIGN_RIGHT:
			switch tb.getTextMode() {
			case gfx.TEXTMODE_FULL:
				x_offset = tb.size.W - len(line)
				pos.X = x_offset
			case gfx.TEXTMODE_HALF:
				x_offset = tb.size.W*2 - len(line)
				pos.X = x_offset / 2
			}
		}

		tb.DrawText(pos, 0, line, col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT}, gfx.TextCellPosition(x_offset%2), tb.textMode)
	}
}
