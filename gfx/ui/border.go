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

	styleFlag   borderStyleFlag
	customStyle *BorderStyle

	//SCROLLBAR STUFF. for now, only vertical scrollbar for lists and the like.
	scrollbar                 bool //whether the scrollbar is enabled. scrollbar will be drawn whenever content doesn't fit
	scrollbarContentHeight    int  //total height of scrolling content
	scrollbarViewportPosition int  //position of the viewed content

	dirty bool
}

// Sets the border style flag. Options are:
// - BORDER_STYLE_DEFAULT: uses the package-wide default at ui.DefaultBorderStyle
// - BORDER_STYLE_INHERIT: uses the borderstyle of its parent element
// - BORDER_STYLE_CUSTOM: uses the borderstyle provided in the 2nd argument
func (b *Border) SetStyle(style_flag borderStyleFlag, style ...BorderStyle) {
	if style_flag == BORDER_STYLE_CUSTOM {
		if style == nil {
			log.Error("Custom border style application failed: no borderstyle provided.")
			return
		}

		b.customStyle = &style[0]
	} else {
		b.customStyle = nil
	}

	b.styleFlag = style_flag
}

// Sets the border colours. Colours set this way will override the colours in the border's set style.
// Use gfx.COL_DEFAULT to default to indicate you want to use one of both of the element's default colours.
func (b *Border) SetColours(colours col.Pair) {
	if b.colours == colours {
		return
	}

	b.colours = colours
	b.dirty = true
}

func (b *Border) EnableScrollbar(content_height, offset int) {
	if !b.scrollbar {
		b.dirty = true
	}

	b.scrollbar = true
	b.UpdateScrollbar(content_height, offset)
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
func (b *Border) UpdateScrollbar(content_height, offset int) {
	if b.scrollbarContentHeight == content_height && b.scrollbarViewportPosition == offset {
		return
	}

	b.scrollbarContentHeight = content_height
	b.scrollbarViewportPosition = offset

	if b.scrollbar {
		b.dirty = true
	}
}

// Enable the border. If no border has been setup via SetupBorder(), a default one will be created. Style defaults to
// ui.DefaultBorderStyle but you can use SetBorderStyle to use something else.
func (e *Element) EnableBorder() {
	e.setBorder(true)
}

// Disables the border. Doesn't delete any setup border, so you can make the old border reappear by enabling it again.
func (e *Element) DisableBorder() {
	e.setBorder(false)
}

func (e *Element) setBorder(bordered bool) {
	if e.Border.enabled == bordered {
		return
	}

	if bordered {
		e.Canvas.Resize(e.size.Grow(2, 2))
		e.SetOrigin(vec.Coord{1, 1})
		e.Border.dirty = true
	} else {
		e.Canvas.Resize(e.size)
		e.SetOrigin(vec.ZERO_COORD)
	}

	e.Border.enabled = bordered
	e.forceRedraw = true
	e.forceParentRedraw()
}

func (e *Element) IsBordered() bool {
	return e.Border.enabled
}

// Creates and enables a border for the element. Title will be shown in the top left, and hint will be shown in the
// bottom right.
// TODO: centered titles? setting borderstyle at the same time?
func (e *Element) SetupBorder(title, hint string) {
	e.Border.title = title
	e.Border.hint = hint
	e.EnableBorder()
}

func (e *Element) getBorderStyle() (style BorderStyle) {
	switch e.Border.styleFlag {
	case BORDER_STYLE_INHERIT:
		if parent := e.GetParent(); parent != nil {
			style = parent.getBorderStyle()
		}
	case BORDER_STYLE_CUSTOM:
		if e.Border.customStyle != nil {
			style = *e.Border.customStyle
		}
	case BORDER_STYLE_DEFAULT:
		style = DefaultBorderStyle
	}

	//find some colour to use, prioritizing current border, then the style, then falling back to the ui default
	colours := e.Border.colours
	if e.focused {
		colours.Fore = defaultFocusColour
	}

	if colours.Fore == col.NONE {
		if style.Colours.Fore != col.NONE {
			colours.Fore = style.Colours.Fore
		} else {
			colours.Fore = DefaultBorderStyle.Colours.Fore
		}
	}
	if colours.Back == col.NONE {
		if style.Colours.Back != col.NONE {
			colours.Back = style.Colours.Back
		} else {
			colours.Back = DefaultBorderStyle.Colours.Back
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

func (e *Element) drawBorder() {
	rect := e.Canvas.Bounds()
	style := e.getBorderStyle()

	// draw box
	e.DrawBox(rect, BorderDepth, style.lineType, style.Colours)

	//decorate and draw title
	if e.Border.title != "" {
		decoratedTitle := style.DecorateText(e.Border.title, style.TitleJustification)
		var offset int
		switch style.TitleJustification {
		case JUSTIFY_LEFT:
			offset = 0
		case JUSTIFY_CENTER:
			offset = (e.size.W - len(decoratedTitle)/2) / 2
		case JUSTIFY_RIGHT:
			offset = e.size.W - len([]rune(decoratedTitle))/2 - 1
		}
		e.DrawText(vec.Coord{offset, -1}, BorderDepth+1, decoratedTitle, style.Colours, gfx.DRAW_TEXT_LEFT)
	}

	//decorate and draw hint
	if e.Border.hint != "" {
		decoratedHint := style.DecorateText(e.Border.hint, style.HintJustification)
		var offset int
		switch style.HintJustification {
		case JUSTIFY_LEFT:
			offset = 0
		case JUSTIFY_CENTER:
			offset = (e.size.W - len(decoratedHint)/2) / 2
		case JUSTIFY_RIGHT:
			offset = e.size.W - len([]rune(decoratedHint))/2
		}
		e.DrawText(vec.Coord{offset, rect.H - 2}, BorderDepth+1, decoratedHint, style.Colours, 0)
	}

	//draw scrollbar if necessary
	if e.Border.scrollbar && e.Border.scrollbarContentHeight > e.size.H {
		right := rect.X + rect.W - 1                      // x coord of the right side of the border
		top := vec.Coord{right, rect.Y + 2}               //top of scrollbar area
		bottom := vec.Coord{right, rect.Y + e.size.H - 1} //bottom of scrollbar area
		e.DrawGlyph(top.Step(vec.DIR_UP), BorderDepth+1, gfx.GLYPH_TRIANGLE_UP)
		e.DrawGlyph(bottom.Step(vec.DIR_DOWN), BorderDepth+1, gfx.GLYPH_TRIANGLE_DOWN)
		e.DrawLine(vec.Line{top, bottom}, BorderDepth+1, gfx.NewGlyphVisuals(gfx.GLYPH_FILL_SPARSE, style.Colours))

		h := e.size.H - 2 //scrollbar area height (not including arrows)
		barSize := util.Clamp(util.RoundFloatToInt(float64(e.size.H)/float64(e.Border.scrollbarContentHeight)*float64(h)), 1, h-1)

		// default to barposition at top ie. no scrolling
		barPos := top
		if e.Border.scrollbarViewportPosition == e.Border.scrollbarContentHeight-e.size.H { // scrolling content is at bottom
			barPos.Y += h - barSize
		} else if e.Border.scrollbarViewportPosition != 0 { //scrolling content is somewhere in the middle. must ensure bar isn't at top or bottom.
			barSize = util.Clamp(barSize, 1, h-2) // ensure bar cannot touch sides, so it shows that we scroll in both directions.
			barPos.Y += util.Clamp(util.RoundFloatToInt(float64(e.Border.scrollbarViewportPosition)/float64(e.Border.scrollbarContentHeight)*float64(h-barSize)), 1, h-barSize-1)
		}

		for i := range barSize {
			e.DrawGlyph(barPos.StepN(vec.DIR_DOWN, i), BorderDepth+1, gfx.GLYPH_FILL_DENSE)
		}
	}
}

func (e *Element) linkChildBorderWithElement(child vec.Bounded, other vec.Bounded) {
	intersection := vec.FindIntersectionRect(child, other)
	if intersection.Area() == 0 {
		return
	}

	switch {
	case intersection.Area() == 1:
		e.linkBorderCell(intersection.Coord)
	case intersection.W == 1 || intersection.H == 1:
		corners := intersection.Corners()
		e.linkBorderCell(corners[0])
		e.linkBorderCell(corners[2])
	default:
		corners := intersection.Corners()
		for _, corner := range corners {
			if corner.IsInPerimeter(child) && corner.IsInPerimeter(other) {
				e.linkBorderCell(corner)
			}
		}
	}
}

func (e *Element) linkBorderCell(pos vec.Coord) {
	if !e.InBounds(pos) {
		return
	}

	if cell := e.GetCell(pos); cell.Mode == gfx.DRAW_GLYPH {
		e.DrawLinkedGlyph(pos, BorderDepth, cell.Glyph)

		//also need to try and link to border titles and decorations drawn at a higher level
		e.DrawLinkedGlyph(pos, BorderDepth+1, e.GetCell(pos).Glyph)
	}
}
