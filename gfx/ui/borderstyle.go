package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

type borderStyleFlag int

const (
	BORDER_STYLE_DEFAULT borderStyleFlag = iota //use default style
	BORDER_STYLE_INHERIT                        //use parent's style
	BORDER_STYLE_CUSTOM                         //use custom defined style
)

// Border neighbour flags
const (
	BORDER_L = 1 << iota
	BORDER_R
	BORDER_U
	BORDER_D

	BORDER_UD = BORDER_U | BORDER_D
	BORDER_LR = BORDER_L | BORDER_R

	BORDER_UL = BORDER_U | BORDER_L
	BORDER_UR = BORDER_U | BORDER_R
	BORDER_DL = BORDER_D | BORDER_L
	BORDER_DR = BORDER_D | BORDER_R

	BORDER_UDL = BORDER_UD | BORDER_L
	BORDER_UDR = BORDER_UD | BORDER_R
	BORDER_ULR = BORDER_LR | BORDER_U
	BORDER_DLR = BORDER_LR | BORDER_D

	BORDER_ALL    = BORDER_LR | BORDER_UD
	BORDER_LONELY = 0 //border cell with no neighbours. why would this ever exist???? i don't know but it's nice to have an unusable zero value
)

// default borderstyle used by all elements
var DefaultBorderStyle BorderStyle

// some pre-defined borderstyles. current options are "Block", "Thin", and "Thick"
var BorderStyles map[string]BorderStyle

type BorderStyle struct {
	Glyphs            [BORDER_ALL + 1]int //glyphs for border drawing, indexed by the BORDER_* constants above
	flagMap           map[int]int         //map of glyphs to borderflags
	TextDecorationL   rune                //character to print on the left of titles/hints
	TextDecorationR   rune                //character to print on the right of titles/hints
	TextDecorationPad rune                //character to pad title/hint in case the decorated string isn't an even number of chars

	Colours     col.Pair //colours for the border. use gfx.COL_DEFAULT to use the default colours of the ui element instead
	DisableLink bool     //toggle this to disable linking.

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

// Returns the order neighbour flags for a particular glyph. If the glyph is invalid, returns 0 (BORDER_LONELY) :(
func (bs *BorderStyle) getBorderFlags(glyph int) int {
	if bs.flagMap == nil {
		bs.buildFlagMap()
	}

	if flags, ok := bs.flagMap[glyph]; ok {
		return flags
	}

	return BORDER_LONELY
}

func getBorderFlagByDirection(dir vec.Direction) int {
	switch dir {
	case vec.DIR_UP:
		return BORDER_U
	case vec.DIR_DOWN:
		return BORDER_D
	case vec.DIR_LEFT:
		return BORDER_L
	case vec.DIR_RIGHT:
		return BORDER_R
	default:
		return BORDER_LONELY
	}
}

func (bs *BorderStyle) glyphIsLinkable(glyph int) bool {
	return bs.getBorderFlags(glyph) != BORDER_LONELY
}

func (bs *BorderStyle) buildFlagMap() {
	bs.flagMap = make(map[int]int)
	for i, glyph := range bs.Glyphs {
		if glyph != 0 {
			bs.flagMap[glyph] = i
		}
	}
}

// setup predefined border styles and set simple default
func init() {
	BorderStyles = make(map[string]BorderStyle)

	var blockStyle BorderStyle
	for i := range BORDER_ALL {
		blockStyle.Glyphs[i] = gfx.GLYPH_FILL
	}
	blockStyle.Colours = col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT}
	blockStyle.DisableLink = true
	BorderStyles["Block"] = blockStyle

	var thinStyle BorderStyle
	thinStyle.Glyphs[BORDER_LR] = gfx.GLYPH_BORDER_LR
	thinStyle.Glyphs[BORDER_UD] = gfx.GLYPH_BORDER_UD
	thinStyle.Glyphs[BORDER_UR] = gfx.GLYPH_BORDER_UR
	thinStyle.Glyphs[BORDER_DR] = gfx.GLYPH_BORDER_DR
	thinStyle.Glyphs[BORDER_UL] = gfx.GLYPH_BORDER_UL
	thinStyle.Glyphs[BORDER_DL] = gfx.GLYPH_BORDER_DL
	thinStyle.Glyphs[BORDER_UDL] = gfx.GLYPH_BORDER_UDL
	thinStyle.Glyphs[BORDER_UDR] = gfx.GLYPH_BORDER_UDR
	thinStyle.Glyphs[BORDER_ULR] = gfx.GLYPH_BORDER_ULR
	thinStyle.Glyphs[BORDER_DLR] = gfx.GLYPH_BORDER_DLR
	thinStyle.Glyphs[BORDER_ALL] = gfx.GLYPH_BORDER_UDLR
	thinStyle.TextDecorationL = gfx.TEXT_BORDER_DECO_LEFT
	thinStyle.TextDecorationR = gfx.TEXT_BORDER_DECO_RIGHT
	thinStyle.TextDecorationPad = gfx.TEXT_BORDER_LR
	thinStyle.Colours = col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT}
	BorderStyles["Thin"] = thinStyle

	var thickStyle BorderStyle
	thickStyle.Glyphs[BORDER_LR] = gfx.GLYPH_BORDER_LLRR
	thickStyle.Glyphs[BORDER_UD] = gfx.GLYPH_BORDER_UUDD
	thickStyle.Glyphs[BORDER_UR] = gfx.GLYPH_BORDER_UURR
	thickStyle.Glyphs[BORDER_DR] = gfx.GLYPH_BORDER_DDRR
	thickStyle.Glyphs[BORDER_UL] = gfx.GLYPH_BORDER_UULL
	thickStyle.Glyphs[BORDER_DL] = gfx.GLYPH_BORDER_DDLL
	thickStyle.Glyphs[BORDER_UDL] = gfx.GLYPH_BORDER_UUDDLL
	thickStyle.Glyphs[BORDER_UDR] = gfx.GLYPH_BORDER_UUDDRR
	thickStyle.Glyphs[BORDER_ULR] = gfx.GLYPH_BORDER_UULLRR
	thickStyle.Glyphs[BORDER_DLR] = gfx.GLYPH_BORDER_DDLLRR
	thickStyle.Glyphs[BORDER_ALL] = gfx.GLYPH_BORDER_UUDDLLRR
	thickStyle.TextDecorationL = gfx.TEXT_BORDER_DECO_LEFT
	thickStyle.TextDecorationR = gfx.TEXT_BORDER_DECO_RIGHT
	thickStyle.TextDecorationPad = gfx.TEXT_BORDER_LR
	thickStyle.Colours = col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT}
	BorderStyles["Thick"] = thickStyle

	DefaultBorderStyle = BorderStyles["Thin"]
}

// note: changing this does NOT dynamically update the styles for ui elements already using the old default.
// TODO: maybe get that working?? when would that even be useful though, who is changing the default border
// style at runtime???
func SetDefaultBorderStyle(style BorderStyle) {
	DefaultBorderStyle = style
}
