package ui

import (
	"github.com/bennicholls/tyumi/gfx"
)

type Border struct {
	top    gfx.Canvas //top, including upper left corner
	bottom gfx.Canvas //bottom, including bottom right corner
	left   gfx.Canvas //left, including bottom left corner
	right  gfx.Canvas //right, including top right corner

	title string
	hint  string

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
}

func (b *Border) DrawToCanvas(canvas *gfx.Canvas, x, y, z int) {
	w, _ := b.top.Dims()
	_, h := b.left.Dims()
	b.top.DrawToCanvas(canvas, x-1, y-1, z)
	b.bottom.DrawToCanvas(canvas, x, y+h-1, z)
	b.left.DrawToCanvas(canvas, x-1, y, z)
	b.right.DrawToCanvas(canvas, x+w-1, y-1, z)
}
