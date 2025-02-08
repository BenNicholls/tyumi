package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
)

type borderStyleFlag int

const (
	BORDER_STYLE_DEFAULT borderStyleFlag = iota //use default style
	BORDER_STYLE_INHERIT                        //use parent's style
	BORDER_STYLE_CUSTOM                         //use custom defined style
)

// some pre-defined borderstyles. current options are "Block", "Thin", and "Thick"
var BorderStyles map[string]BorderStyle

type BorderStyle struct {
	lineType     gfx.LineType
	DefaultGlyph int // default glyph to draw when not using linked line glyphs

	TextDecorationL   rune //character to print on the left of titles/hints
	TextDecorationR   rune //character to print on the right of titles/hints
	TextDecorationPad rune //character to pad title/hint in case the decorated string isn't an even number of chars

	Colours     col.Pair //colours for the border. use gfx.COL_DEFAULT to use the default colours of the ui element instead
	DisableLink bool     //toggle this to disable linking.

	//scrollbar styling stuff should go here as well
}

// GetGlyph returns the appropriate border glyph to link to neighbours as described in the naighbour_flags.
func (bs *BorderStyle) GetGlyph(neighbour_flags int) int {
	if bs.lineType == gfx.LINETYPE_NONE {
		return bs.DefaultGlyph
	}

	return gfx.LineStyles[bs.lineType].Glyphs[neighbour_flags]
}

func (bs BorderStyle) DecorateText(text string) (decoratedText string) {
	if bs.TextDecorationL != rune(0) {
		decoratedText += string(bs.TextDecorationL)
	}
	decoratedText += text
	if bs.TextDecorationR != rune(0) {
		decoratedText += string(bs.TextDecorationR)
	}

	return
}

// Returns the border neighbour flags for a particular glyph. If the glyph does not link with this border, 
// returns 0 (LINK_NONE) :(
func (bs *BorderStyle) getBorderFlags(glyph int) int {
	if bs.lineType == gfx.LINETYPE_NONE {
		return gfx.LINK_NONE
	}

	return gfx.LineStyles[bs.lineType].GetBorderFlags(glyph)
}

func (bs *BorderStyle) glyphIsLinkable(glyph int) bool {
	return bs.getBorderFlags(glyph) != gfx.LINK_NONE
}

// setup predefined border styles and set simple default
func init() {
	BorderStyles = make(map[string]BorderStyle)

	var blockStyle BorderStyle
	blockStyle.lineType = gfx.LINETYPE_NONE
	blockStyle.DefaultGlyph = gfx.GLYPH_FILL
	blockStyle.Colours = col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT}
	blockStyle.DisableLink = true
	BorderStyles["Block"] = blockStyle

	var thinStyle BorderStyle
	thinStyle.lineType = gfx.LINETYPE_THIN
	thinStyle.TextDecorationL = gfx.TEXT_BORDER_DECO_LEFT
	thinStyle.TextDecorationR = gfx.TEXT_BORDER_DECO_RIGHT
	thinStyle.TextDecorationPad = gfx.TEXT_BORDER_LR
	thinStyle.Colours = col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT}
	BorderStyles["Thin"] = thinStyle

	var thickStyle BorderStyle
	thickStyle.lineType = gfx.LINETYPE_THICK
	thickStyle.TextDecorationL = gfx.TEXT_BORDER_DECO_LEFT
	thickStyle.TextDecorationR = gfx.TEXT_BORDER_DECO_RIGHT
	thickStyle.TextDecorationPad = gfx.TEXT_BORDER_LR
	thickStyle.Colours = col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT}
	BorderStyles["Thick"] = thickStyle

	defaultBorderStyle = BorderStyles["Thin"]
}
