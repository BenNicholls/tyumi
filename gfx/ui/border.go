package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type Border struct {
	top    gfx.Canvas //top, including upper left corner
	bottom gfx.Canvas //bottom, including bottom right corner
	left   gfx.Canvas //left, including bottom left corner
	right  gfx.Canvas //right, including top right corner

	title string
	hint  string
	style BorderStyle

	//SCROLLBAR STUFF. for now, only vertical scrollbar for lists and the like.
	scrollbar                 bool //whether the scrollbar is enabled. scrollbar will be drawn whenever content doesn't fit
	scrollbarContentHeight    int  //total height of scrolling content
	scrollbarViewportPosition int  //position of the viewed content

	dirty bool //flag indicates when border needs to be re-rendered.
}

func NewBorder(element_size vec.Dims) Border {
	b := Border{}

	b.style = DefaultBorderStyle

	b.top.Init(element_size.W+1, 1)
	b.bottom.Init(element_size.W+1, 1)
	b.left.Init(1, element_size.H+1)
	b.right.Init(1, element_size.H+1)

	b.dirty = true

	return b
}

// renders the border to the internal canvas
func (b *Border) Render() {
	if !b.dirty {
		return
	}

	for cursor := range vec.EachCoord(b.top.Bounds()) { //top and bottom
		b.top.DrawGlyph(cursor, 0, b.style.Glyphs[BORDER_LR])
		b.bottom.DrawGlyph(cursor, 0, b.style.Glyphs[BORDER_LR])
	}

	for cursor := range vec.EachCoord(b.left.Bounds()) {
		b.left.DrawGlyph(cursor, 0, b.style.Glyphs[BORDER_UD])
		b.right.DrawGlyph(cursor, 0, b.style.Glyphs[BORDER_UD])
	}

	w, h := b.top.Size().W, b.left.Size().H
	b.top.DrawGlyph(vec.Coord{0, 0}, 0, b.style.Glyphs[BORDER_DR])        //upper left corner
	b.right.DrawGlyph(vec.Coord{0, 0}, 0, b.style.Glyphs[BORDER_DL])      //upper right corner
	b.bottom.DrawGlyph(vec.Coord{w - 1, 0}, 0, b.style.Glyphs[BORDER_UL]) //bottom right corner
	b.left.DrawGlyph(vec.Coord{0, h - 1}, 0, b.style.Glyphs[BORDER_UR])   //bottom left corner

	//decorate and draw title
	if b.title != "" {
		decoratedTitle := b.style.DecorateText(b.title)
		if len([]rune(decoratedTitle))%2 == 1 {
			decoratedTitle += string(b.style.TextDecorationPad)
		}
		b.top.DrawText(vec.Coord{1, 0}, 0, decoratedTitle, gfx.COL_DEFAULT, gfx.COL_DEFAULT, gfx.DRAW_TEXT_LEFT)
	}

	//decorate and draw hint
	if b.hint != "" {
		decoratedHint := b.style.DecorateText(b.hint)
		if len([]rune(decoratedHint))%2 == 1 {
			decoratedHint = string(b.style.TextDecorationPad) + decoratedHint
		}
		hintOffset := w - len([]rune(decoratedHint))/2 - 1
		b.bottom.DrawText(vec.Coord{hintOffset, 0}, 0, decoratedHint, gfx.COL_DEFAULT, gfx.COL_DEFAULT, 0)
	}

	//draw scrollbar if necessary
	if b.scrollbar && b.scrollbarContentHeight > h-1 {
		b.right.DrawGlyph(vec.Coord{0, 1}, 0, gfx.GLYPH_TRIANGLE_UP)
		b.right.DrawGlyph(vec.Coord{0, h - 1}, 0, gfx.GLYPH_TRIANGLE_DOWN)

		barSize := util.Clamp(util.RoundFloatToInt(float64(h-1)/float64(b.scrollbarContentHeight)*float64(h-3)), 1, h-4)

		var barPos int
		if b.scrollbarViewportPosition+h-1 >= b.scrollbarContentHeight {
			barPos = h - 3 - barSize
		} else {
			barPos = util.Clamp(util.RoundFloatToInt(float64(b.scrollbarViewportPosition)/float64(b.scrollbarContentHeight)*float64(h-3)), 0, h-4-barSize)
		}

		for i := range barSize {
			b.right.DrawGlyph(vec.Coord{0, i + barPos + 2}, 0, gfx.GLYPH_FILL)
		}
	}

	b.dirty = false
}

func (b *Border) DrawToCanvas(canvas *gfx.Canvas, x, y, depth int) {
	w, h := b.top.Size().W, b.left.Size().H
	b.top.DrawToCanvas(canvas, vec.Coord{x - 1, y - 1}, depth)
	b.bottom.DrawToCanvas(canvas, vec.Coord{x, y + h - 1}, depth)
	b.left.DrawToCanvas(canvas, vec.Coord{x - 1, y}, depth)
	b.right.DrawToCanvas(canvas, vec.Coord{x + w - 1, y - 1}, depth)
}

func (b *Border) EnableScrollbar(height, pos int) {
	b.scrollbar = true
	b.scrollbarContentHeight = height
	b.scrollbarViewportPosition = pos
	b.dirty = true
}

// Updates the position/size of the scrollbar.
// NOTE: this does NOT enable the scrollbar. you have to do that manually during setup.
func (b *Border) UpdateScrollbar(height, pos int) {
	if b.scrollbarContentHeight == height && b.scrollbarViewportPosition == pos {
		return
	}

	b.scrollbarContentHeight = height
	b.scrollbarViewportPosition = pos
	b.dirty = true
}

