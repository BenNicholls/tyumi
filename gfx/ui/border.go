package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// This is the depth used by the UI system for drawing borders while rendering an element. Any children that need to be
// linked to this element's borders needs to be at this depth. If you want to draw over the border during rendering,
// you must draw at a higher depth than this. Border decorations, like the title/hint and scrollbars, are drawn one
// level higher than this.
var BorderDepth int = 50

type Border struct {
	enabled bool
	title   string
	hint    string
	colours col.Pair

	styleFlag    borderStyleFlag
	custom_style *BorderStyle

	//SCROLLBAR STUFF. for now, only vertical scrollbar for lists and the like.
	scrollbar                 bool //whether the scrollbar is enabled. scrollbar will be drawn whenever content doesn't fit
	scrollbarContentHeight    int  //total height of scrolling content
	scrollbarViewportPosition int  //position of the viewed content

	dirty bool
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

// Enable the border. If no border has been setup via SetupBorder(), a default one will be created. Style defaults to
// ui.DefaultBorderStyle but you can use SetBorderStyle to use something else.
func (e *ElementPrototype) EnableBorder() {
	e.setBorder(true)
}

// Disables the border. Doesn't delete any setup border, so you can make the old border reappear by enabling it again.
func (e *ElementPrototype) DisableBorder() {
	e.setBorder(false)
}

func (e *ElementPrototype) setBorder(bordered bool) {
	if e.border.enabled == bordered {
		return
	}

	if bordered {
		e.Canvas.Resize(e.size.W+2, e.size.H+2)
		e.SetOrigin(vec.Coord{1, 1})
		e.border.dirty = true
	} else {
		e.Canvas.Resize(e.size.W, e.size.H)
		e.SetOrigin(vec.ZERO_COORD)
	}

	e.border.enabled = bordered
	e.forceRedraw = true
	e.forceParentRedraw()
}

func (e *ElementPrototype) IsBordered() bool {
	return e.border.enabled
}

// Creates and enables a border for the element. Title will be shown in the top left, and hint will be shown in the
// bottom right.
// TODO: centered titles? setting borderstyle at the same time?
func (e *ElementPrototype) SetupBorder(title, hint string) {
	e.border.title = title
	e.border.hint = hint
	e.EnableBorder()
}

// Sets the border style flag. Options are:
// - BORDER_STYLE_DEFAULT: uses the package-wide default at ui.DefaultBorderStyle
// - BORDER_STYLE_INHERIT: uses the borderstyle of its parent element
// - BORDER_STYLE_CUSTOM: uses the borderstyle provided in the 2nd argument
func (e *ElementPrototype) SetBorderStyle(styleFlag borderStyleFlag, borderStyle ...BorderStyle) {
	if styleFlag == BORDER_STYLE_CUSTOM {
		if borderStyle == nil {
			log.Error("Custom border style application failed: no borderstyle provided.")
			return
		}

		e.border.custom_style = &borderStyle[0]
	} else {
		e.border.custom_style = nil
	}

	e.border.styleFlag = styleFlag
}

// Sets the border colours. Colours set this way will override the colours in the border's set style. 
// Use gfx.COL_DEFAULT to default to indicate you want to use one of both of the element's default colours.
func (e *ElementPrototype) SetBorderColours(col col.Pair) {
	if e.border.colours == col {
		return
	}

	e.border.colours = col
	e.border.dirty = true
}

func (e *ElementPrototype) getBorderStyle() (style BorderStyle) {
	switch e.border.styleFlag {
	case BORDER_STYLE_INHERIT:
		if parent := e.GetParent(); parent != nil {
			style = parent.getBorderStyle()
		}
	case BORDER_STYLE_CUSTOM:
		if e.border.custom_style != nil {
			style = *e.border.custom_style
		}
	case BORDER_STYLE_DEFAULT:
		style = defaultBorderStyle
	}

	//find some colour to use, prioritizing current border, then the style, then falling back to the ui default
	colours := e.border.colours
	if colours.Fore == col.NONE {
		if style.Colours.Fore != col.NONE {
			colours.Fore = style.Colours.Fore
		} else {
			colours.Fore = defaultBorderStyle.Colours.Fore
		}
	}
	if colours.Back == col.NONE {
		if style.Colours.Back != col.NONE {
			colours.Back = style.Colours.Back
		} else {
			colours.Back = defaultBorderStyle.Colours.Back
		}
	}

	//if any colours are gfx.COL_DEFAULT, replace with canvas colours
	if colours.Fore == gfx.COL_DEFAULT {
		colours.Fore = e.DefaultColours().Fore
	}
	if colours.Back == gfx.COL_DEFAULT {
		colours.Back = e.DefaultColours().Back
	}

	style.Colours = colours

	return
}

func (e *ElementPrototype) DrawBorder() {
	rect := e.Canvas.Bounds()
	style := e.getBorderStyle()

	// draw box
	e.DrawBox(rect, BorderDepth, style.lineType, style.Colours)

	//decorate and draw title
	if e.border.title != "" {
		decoratedTitle := style.DecorateText(e.border.title)
		if len([]rune(decoratedTitle))%2 == 1 {
			decoratedTitle += string(style.TextDecorationPad)
		}
		e.DrawText(rect.Coord.Add(vec.Coord{1, 0}), BorderDepth+1, decoratedTitle, style.Colours, gfx.DRAW_TEXT_LEFT)
	}

	//decorate and draw hint
	if e.border.hint != "" {
		decoratedHint := style.DecorateText(e.border.hint)
		if len([]rune(decoratedHint))%2 == 1 {
			decoratedHint = string(style.TextDecorationPad) + decoratedHint
		}
		hintOffset := rect.W - len([]rune(decoratedHint))/2 - 1
		e.DrawText(rect.Coord.Add(vec.Coord{hintOffset, e.Bounds().H - 1}), BorderDepth+1, decoratedHint, style.Colours, 0)
	}

	//draw scrollbar if necessary
	if e.border.scrollbar && e.border.scrollbarContentHeight > e.size.H {
		right := rect.X + rect.W - 1                      // x coord of the right side of the border
		top := vec.Coord{right, rect.Y + 2}               //top of scrollbar area
		bottom := vec.Coord{right, rect.Y + e.size.H - 1} //bottom of scrollbar area
		e.DrawGlyph(top.Step(vec.DIR_UP), BorderDepth+1, gfx.GLYPH_TRIANGLE_UP)
		e.DrawGlyph(bottom.Step(vec.DIR_DOWN), BorderDepth+1, gfx.GLYPH_TRIANGLE_DOWN)
		e.DrawLine(vec.Line{top, bottom}, BorderDepth+1, gfx.NewGlyphVisuals(gfx.GLYPH_FILL_SPARSE, style.Colours))

		h := e.size.H - 2 //scrollbar area height (not including arrows)
		barSize := util.Clamp(util.RoundFloatToInt(float64(e.size.H)/float64(e.border.scrollbarContentHeight)*float64(h)), 1, h-1)

		barPos := top                                                                       // default to barposition at top ie. no scrolling
		if e.border.scrollbarViewportPosition == e.border.scrollbarContentHeight-e.size.H { // scrolling content is at bottom
			barPos.Y += h - barSize
		} else if e.border.scrollbarViewportPosition != 0 { //scrolling content is somewhere in the middle. must ensure bar isn't at top or bottom.
			barSize = util.Clamp(barSize, 1, h-2) // ensure bar cannot touch sides, so it shows that we scroll in both directions.
			barPos.Y += util.Clamp(util.RoundFloatToInt(float64(e.border.scrollbarViewportPosition)/float64(e.border.scrollbarContentHeight)*float64(h-barSize)), 1, h-barSize-1)
		}

		for i := range barSize {
			e.DrawGlyph(barPos.StepN(vec.DIR_DOWN, i), BorderDepth+1, gfx.GLYPH_FILL_DENSE)
		}
	}
}

func (e *ElementPrototype) linkBorder() {
	for cursor := range vec.EachCoordInPerimeter(e.Canvas) {
		cell := e.GetCell(cursor)
		switch cell.Mode {
		case gfx.DRAW_GLYPH:
			e.DrawLinkedGlyph(cursor, BorderDepth, cell.Glyph)
		case gfx.DRAW_TEXT:
			// these are some corner cases, needed because titles and hints are drawn on a higher layer than
			// normal borders. so for borders to link to special title decorations and padding characters we have
			// to manually check for them.
			if cell.Chars[0] == gfx.TEXT_BORDER_DECO_LEFT || cell.Chars[0] == gfx.TEXT_BORDER_LR {
				left := cursor.Step(vec.DIR_LEFT)
				if e.InBounds(left) {
					left_cell := e.GetCell(left)
					if left_cell.Mode == gfx.DRAW_GLYPH {
						linkedGlyph := e.CalcLinkedGlyph(left_cell.Glyph, left, BorderDepth+1)
						e.DrawGlyph(left, BorderDepth, linkedGlyph)
					}
				}
			} else if cell.Chars[1] == gfx.TEXT_BORDER_DECO_RIGHT || cell.Chars[1] == gfx.TEXT_BORDER_LR {
				right := cursor.Step(vec.DIR_RIGHT)
				if e.InBounds(right) {
					right_cell := e.GetCell(right)
					if right_cell.Mode == gfx.DRAW_GLYPH {
						linkedGlyph := e.CalcLinkedGlyph(right_cell.Glyph, right, BorderDepth+1)
						e.DrawGlyph(right, BorderDepth, linkedGlyph)
					}
				}
			}
		}
	}
}
