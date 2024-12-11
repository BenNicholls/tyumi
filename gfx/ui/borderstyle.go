package ui

import "github.com/bennicholls/tyumi/gfx"

type borderStyleFlag int

const (
	BORDER_STYLE_DEFAULT borderStyleFlag = iota //use default style
	BORDER_STYLE_INHERIT                        //use parent's style
	BORDER_STYLE_CUSTOM                         //use custom defined style
)

// Border neighbour flags
const (
	BORDER_LONELY = 0 //border cell with no nieghbours. why would this ever exist???? i don't know but it's nice to have an unusable zero value
	BORDER_L      = 1 << iota
	BORDER_R
	BORDER_U
	BORDER_D

	BORDER_LR = BORDER_L | BORDER_R
	BORDER_UD = BORDER_U | BORDER_D

	BORDER_UR = BORDER_U | BORDER_R
	BORDER_DR = BORDER_D | BORDER_R
	BORDER_UL = BORDER_U | BORDER_L
	BORDER_DL = BORDER_D | BORDER_L

	BORDER_ALL = BORDER_LR | BORDER_UD
)

// default borderstyle used by all elements
var DefaultBorderStyle BorderStyle

// some pre-defined borderstyles. current options are "Block", "Thin", and "Thick"
var BorderStyles map[string]BorderStyle

type BorderStyle struct {
	Glyphs            [BORDER_ALL]int //glyphs for border drawing, indexed by the BORDER_* constants above
	TextDecorationL   rune            //character to print on the left of titles/hints
	TextDecorationR   rune            //character to print on the right of titles/hints
	TextDecorationPad rune            //character to pad title/hint in case the decorated string isn't an even number of chars

	//scrollbar styling stuff should go here as well
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

// setup predefined border styles and set simple crummy default
func init() {

	BorderStyles = make(map[string]BorderStyle)

	var blockStyle BorderStyle
	for i := range BORDER_ALL {
		blockStyle.Glyphs[i] = gfx.GLYPH_FILL
	}
	BorderStyles["Block"] = blockStyle

	var thinStyle BorderStyle
	thinStyle.Glyphs[BORDER_LR] = gfx.GLYPH_BORDER_LR
	thinStyle.Glyphs[BORDER_UD] = gfx.GLYPH_BORDER_UD
	thinStyle.Glyphs[BORDER_UR] = gfx.GLYPH_BORDER_UR
	thinStyle.Glyphs[BORDER_DR] = gfx.GLYPH_BORDER_DR
	thinStyle.Glyphs[BORDER_UL] = gfx.GLYPH_BORDER_UL
	thinStyle.Glyphs[BORDER_DL] = gfx.GLYPH_BORDER_DL
	thinStyle.TextDecorationL = gfx.TEXT_BORDER_DECO_LEFT
	thinStyle.TextDecorationR = gfx.TEXT_BORDER_DECO_RIGHT
	thinStyle.TextDecorationPad = gfx.TEXT_BORDER_LR
	BorderStyles["Thin"] = thinStyle

	var thickStyle BorderStyle
	thickStyle.Glyphs[BORDER_LR] = gfx.GLYPH_BORDER_LLRR
	thickStyle.Glyphs[BORDER_UD] = gfx.GLYPH_BORDER_UUDD
	thickStyle.Glyphs[BORDER_UR] = gfx.GLYPH_BORDER_UURR
	thickStyle.Glyphs[BORDER_DR] = gfx.GLYPH_BORDER_DDRR
	thickStyle.Glyphs[BORDER_UL] = gfx.GLYPH_BORDER_UULL
	thickStyle.Glyphs[BORDER_DL] = gfx.GLYPH_BORDER_DDLL
	thickStyle.TextDecorationL = gfx.TEXT_BORDER_DECO_LEFT
	thickStyle.TextDecorationR = gfx.TEXT_BORDER_DECO_RIGHT
	thickStyle.TextDecorationPad = gfx.TEXT_BORDER_LR
	BorderStyles["Thick"] = thickStyle

	DefaultBorderStyle = BorderStyles["Thin"]
}

// note: changing this does NOT dynamically update the styles for ui elements already using the old default.
// TODO: maybe get that working?? when would that even be useful though, who is changing the default border
// style at runtime???
func SetDefaultBorderStyle(style BorderStyle) {
	DefaultBorderStyle = style
}
