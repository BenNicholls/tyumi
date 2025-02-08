package gfx

import "github.com/bennicholls/tyumi/vec"

func (c *Canvas) DrawLine(line vec.Line, depth int, brush Visuals) {
	for cursor := range line.EachCoord() {
		c.DrawVisuals(cursor, depth, brush)
	}
}

type LineType int

const (
	LINETYPE_THIN LineType = iota
	LINETYPE_THICK

	linetype_max

	LINETYPE_NONE LineType = -1 
)

// neighbour linking flags for linked lines
const (
	LINK_L = 1 << iota
	LINK_R
	LINK_U
	LINK_D

	LINK_UD = LINK_U | LINK_D
	LINK_LR = LINK_L | LINK_R

	LINK_UL = LINK_U | LINK_L
	LINK_UR = LINK_U | LINK_R
	LINK_DL = LINK_D | LINK_L
	LINK_DR = LINK_D | LINK_R

	LINK_UDL = LINK_UD | LINK_L
	LINK_UDR = LINK_UD | LINK_R
	LINK_ULR = LINK_LR | LINK_U
	LINK_DLR = LINK_LR | LINK_D

	LINK_ALL  = LINK_LR | LINK_UD
	LINK_NONE = 0 // not sure why this would ever happen but it's nice to have a zero value
)

func GetLinkFlagByDirection(dir vec.Direction) int {
	switch dir {
	case vec.DIR_UP:
		return LINK_U
	case vec.DIR_DOWN:
		return LINK_D
	case vec.DIR_LEFT:
		return LINK_L
	case vec.DIR_RIGHT:
		return LINK_R
	default:
		return LINK_NONE
	}
}

type LineStyle struct {
	Glyphs  [LINK_ALL + 1]int //glyphs for line drawing, indexed by the LINK_* constants above
	flagMap map[int]int       //map of glyphs to linkflags
}

// Returns the link flags for a particular glyph. If the glyph is invalid, returns 0 (LINK_NONE) :(
func (ls *LineStyle) GetBorderFlags(glyph int) int {
	if ls.flagMap == nil {
		ls.buildFlagMap()
	}

	if flags, ok := ls.flagMap[glyph]; ok {
		return flags
	}

	return LINK_NONE
}

func (ls *LineStyle) glyphIsLinkable(glyph int) bool {
	return ls.GetBorderFlags(glyph) != LINK_NONE
}

func (ls *LineStyle) buildFlagMap() {
	ls.flagMap = make(map[int]int)
	for i, glyph := range ls.Glyphs {
		if glyph != 0 {
			ls.flagMap[glyph] = i
		}
	}
}

var LineStyles [linetype_max]LineStyle

// define the linking information for thin and thick lines
func init() {
	var thinStyle LineStyle
	thinStyle.Glyphs[LINK_LR] = GLYPH_BORDER_LR
	thinStyle.Glyphs[LINK_UD] = GLYPH_BORDER_UD
	thinStyle.Glyphs[LINK_UR] = GLYPH_BORDER_UR
	thinStyle.Glyphs[LINK_DR] = GLYPH_BORDER_DR
	thinStyle.Glyphs[LINK_UL] = GLYPH_BORDER_UL
	thinStyle.Glyphs[LINK_DL] = GLYPH_BORDER_DL
	thinStyle.Glyphs[LINK_UDL] = GLYPH_BORDER_UDL
	thinStyle.Glyphs[LINK_UDR] = GLYPH_BORDER_UDR
	thinStyle.Glyphs[LINK_ULR] = GLYPH_BORDER_ULR
	thinStyle.Glyphs[LINK_DLR] = GLYPH_BORDER_DLR
	thinStyle.Glyphs[LINK_ALL] = GLYPH_BORDER_UDLR
	LineStyles[LINETYPE_THIN] = thinStyle

	var thickStyle LineStyle
	thickStyle.Glyphs[LINK_LR] = GLYPH_BORDER_LLRR
	thickStyle.Glyphs[LINK_UD] = GLYPH_BORDER_UUDD
	thickStyle.Glyphs[LINK_UR] = GLYPH_BORDER_UURR
	thickStyle.Glyphs[LINK_DR] = GLYPH_BORDER_DDRR
	thickStyle.Glyphs[LINK_UL] = GLYPH_BORDER_UULL
	thickStyle.Glyphs[LINK_DL] = GLYPH_BORDER_DDLL
	thickStyle.Glyphs[LINK_UDL] = GLYPH_BORDER_UUDDLL
	thickStyle.Glyphs[LINK_UDR] = GLYPH_BORDER_UUDDRR
	thickStyle.Glyphs[LINK_ULR] = GLYPH_BORDER_UULLRR
	thickStyle.Glyphs[LINK_DLR] = GLYPH_BORDER_DDLLRR
	thickStyle.Glyphs[LINK_ALL] = GLYPH_BORDER_UUDDLLRR
	LineStyles[LINETYPE_THICK] = thickStyle
}
