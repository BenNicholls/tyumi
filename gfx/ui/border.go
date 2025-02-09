package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type Border struct {
	enabled bool
	title   string
	hint    string
	colours col.Pair

	styleFlag borderStyleFlag
	style     *BorderStyle

	//SCROLLBAR STUFF. for now, only vertical scrollbar for lists and the like.
	scrollbar                 bool //whether the scrollbar is enabled. scrollbar will be drawn whenever content doesn't fit
	scrollbarContentHeight    int  //total height of scrolling content
	scrollbarViewportPosition int  //position of the viewed content

	dirty bool
}

func (b *Border) setColours(col col.Pair) {
	if b.colours == col {
		return
	}

	b.colours = col
	b.dirty = true
}

func (b Border) getStyle() (style *BorderStyle) {
	if b.style != nil {
		style = b.style
	} else {
		style = &defaultBorderStyle
	}

	return
}

func (b *Border) Draw(canvas *gfx.Canvas, area vec.Rect) {
	style := b.getStyle()

	//determine colours
	colours := style.Colours
	if colours.Fore == gfx.COL_DEFAULT {
		if b.colours.Fore == gfx.COL_DEFAULT {
			colours.Fore = canvas.DefaultColours().Fore
		} else {
			colours.Fore = b.colours.Fore
		}
	}

	if colours.Back == gfx.COL_DEFAULT {
		if b.colours.Back == gfx.COL_DEFAULT {
			colours.Back = canvas.DefaultColours().Back
		} else {
			colours.Back = b.colours.Back
		}
	}

	// draw box
	canvas.DrawBox(area, 0, style.lineType, colours)

	//decorate and draw title
	if b.title != "" {
		decoratedTitle := style.DecorateText(b.title)
		if len([]rune(decoratedTitle))%2 == 1 {
			decoratedTitle += string(style.TextDecorationPad)
		}
		canvas.DrawText(vec.Coord{1, 0}.Add(area.Coord), 0, decoratedTitle, colours, gfx.DRAW_TEXT_LEFT)
	}

	//decorate and draw hint
	if b.hint != "" {
		decoratedHint := style.DecorateText(b.hint)
		if len([]rune(decoratedHint))%2 == 1 {
			decoratedHint = string(style.TextDecorationPad) + decoratedHint
		}
		hintOffset := area.W - len([]rune(decoratedHint))/2 - 1
		canvas.DrawText(area.Coord.Add(vec.Coord{hintOffset, area.H - 1}), 0, decoratedHint, colours, 0)
	}

	//draw scrollbar if necessary
	if b.scrollbar {
		h := area.H - 2
		right := area.X + area.W - 1
		canvas.DrawGlyph(vec.Coord{right, area.Y + 1}, 0, gfx.GLYPH_TRIANGLE_UP)
		canvas.DrawGlyph(vec.Coord{right, area.Y + h}, 0, gfx.GLYPH_TRIANGLE_DOWN)

		barSize := util.Clamp(util.RoundFloatToInt(float64(h)/float64(b.scrollbarContentHeight)*float64(h-2)), 1, h-3)

		var barPos int
		if b.scrollbarViewportPosition+h >= b.scrollbarContentHeight {
			barPos = h - 2 - barSize
		} else {
			barPos = util.Clamp(util.RoundFloatToInt(float64(b.scrollbarViewportPosition)/float64(b.scrollbarContentHeight)*float64(h-2)), 0, h-3-barSize)
		}

		for i := range barSize {
			canvas.DrawGlyph(vec.Coord{right, area.Y + i + barPos + 2}, 0, gfx.GLYPH_FILL)
		}
	}

	for cursor := range vec.EachCoordInPerimeter(area) {
		if !canvas.InBounds(cursor) {
			continue
		}

		cell := canvas.GetCell(cursor)
		if cell.Mode == gfx.DRAW_GLYPH {
			linkedGlyph := canvas.CalcLinkedGlyph(canvas.GetCell(cursor).Glyph, cursor, 0)
			canvas.DrawGlyph(cursor, 0, linkedGlyph)
			canvas.DrawColours(cursor, 0, col.Pair{col.ORANGE, col.OLIVE})
		}
	}
}

func (b *Border) EnableScrollbar(content_height, pos int) {
	if !b.scrollbar {
		b.dirty = true
	}

	b.scrollbar = true
	b.UpdateScrollbar(content_height, pos)
}

func (b *Border) DisableScrollbar() {
	if !b.scrollbar {
		return
	}

	b.scrollbar = false
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

	if b.scrollbar {
		b.dirty = true
	}
}
