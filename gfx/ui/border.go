package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/util"
)

type Border struct {
	top    gfx.Canvas //top, including upper left corner
	bottom gfx.Canvas //bottom, including bottom right corner
	left   gfx.Canvas //left, including bottom left corner
	right  gfx.Canvas //right, including top right corner

	title string
	hint  string

	//SCROLLBAR STUFF. for now, only vertical scrollbar for lists and the like.
	scrollbar                 bool //whether the scrollbar is enabled. scrollbar will be drawn whenever content doesn't fit
	scrollbarContentHeight    int  //total height of scrolling content
	scrollbarViewportPosition int  //position of the viewed content

	dirty bool //flag indicates when border needs to be re-rendered.
}

func NewBorder(w, h int) Border {
	b := Border{}

	b.top.Init(w+1, 1)
	b.bottom.Init(w+1, 1)
	b.left.Init(1, h+1)
	b.right.Init(1, h+1)

	b.dirty = true

	return b
}

//renders the border to the internal canvas
func (b *Border) Render() {
	if !b.dirty {
		return
	}

	w, _ := b.top.Dims()
	for i := 0; i < w; i++ { //top and bottom
		b.top.SetGlyph(i, 0, 0, gfx.GLYPH_BORDER_LR)
		b.bottom.SetGlyph(i, 0, 0, gfx.GLYPH_BORDER_LR)
	}
	_, h := b.left.Dims()
	for i := 0; i < h; i++ {
		b.left.SetGlyph(0, i, 0, gfx.GLYPH_BORDER_UD)
		b.right.SetGlyph(0, i, 0, gfx.GLYPH_BORDER_UD)
	}
	b.top.SetGlyph(0, 0, 0, gfx.GLYPH_BORDER_DR)      //upper left corner
	b.right.SetGlyph(0, 0, 0, gfx.GLYPH_BORDER_DL)    //upper right corner
	b.bottom.SetGlyph(w-1, 0, 0, gfx.GLYPH_BORDER_UL) //bottom right corner
	b.left.SetGlyph(0, h-1, 0, gfx.GLYPH_BORDER_UR)   //bottom left corner

	//decorate and draw title
	if b.title != "" {
		decoratedTitle := string(gfx.TEXT_BORDER_DECO_LEFT) + b.title + string(gfx.TEXT_BORDER_DECO_RIGHT)
		if len(b.title)%2 == 1 {
			decoratedTitle += string(gfx.TEXT_BORDER_LR)
		}
		b.top.DrawText(1, 0, 0, decoratedTitle, gfx.COL_DEFAULT, gfx.COL_DEFAULT, 0)
	}

	//decorate and draw hint
	if b.hint != "" {
		decoratedHint := string(gfx.TEXT_BORDER_DECO_LEFT) + b.hint + string(gfx.TEXT_BORDER_DECO_RIGHT)
		hintOffset := w - len(b.hint)/2 - 1
		if len(b.hint)%2 == 1 {
			decoratedHint = string(gfx.TEXT_BORDER_LR) + decoratedHint
			hintOffset -= 1
		}
		b.bottom.DrawText(hintOffset-1, 0, 0, decoratedHint, gfx.COL_DEFAULT, gfx.COL_DEFAULT, 0)
	}

	//draw scrollbar if necessary
	if b.scrollbar && b.scrollbarContentHeight > h-1 {
		b.right.SetGlyph(0, 1, 0, gfx.GLYPH_TRIANGLE_UP)
		b.right.SetGlyph(0, h-1, 0, gfx.GLYPH_TRIANGLE_DOWN)

		barSize := util.Clamp(util.RoundFloatToInt(float64(h-1)/float64(b.scrollbarContentHeight)*float64(h-3)), 1, h-4)
		
		var barPos int
		if b.scrollbarViewportPosition+h-1 >= b.scrollbarContentHeight {
			barPos = h - 3 - barSize
		} else {
			barPos = util.Clamp(util.RoundFloatToInt(float64(b.scrollbarViewportPosition)/float64(b.scrollbarContentHeight)*float64(h-3)), 0, h-4-barSize)
		}

		for i := 0; i < barSize; i++ {
			b.right.SetGlyph(0, i+barPos+2, 0, gfx.GLYPH_FILL)
		}
	}

	b.dirty = false
}

func (b *Border) DrawToCanvas(canvas *gfx.Canvas, x, y, z int) {
	w, _ := b.top.Dims()
	_, h := b.left.Dims()
	b.top.DrawToCanvas(canvas, x-1, y-1, z)
	b.bottom.DrawToCanvas(canvas, x, y+h-1, z)
	b.left.DrawToCanvas(canvas, x-1, y, z)
	b.right.DrawToCanvas(canvas, x+w-1, y-1, z)
}

func (b *Border) EnableScrollbar(height, pos int) {
	b.scrollbar = true
	b.scrollbarContentHeight = height
	b.scrollbarViewportPosition = pos
	b.dirty = true
}

//Updates the position/size of the scrollbar.
//NOTE: this does NOT enable the scrollbar. you have to do that manually during setup.
func (b *Border) UpdateScrollbar(height, pos int) {
	if b.scrollbarContentHeight == height && b.scrollbarViewportPosition == pos {
		return
	}

	b.scrollbarContentHeight = height
	b.scrollbarViewportPosition = pos
	b.dirty = true
}
