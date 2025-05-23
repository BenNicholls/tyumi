package gfx

import (
	"fmt"

	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

// The basic visual definition of a single-tile object that can be drawn to the screen.
// Visuals can be one of 3 Modes: Glyph drawing, or Text drawing, or disabled.
// Each mode uses a different spritesheet, and Text drawing can draw 2 letters to a cell, hence the 2 Chars.
type Visuals struct {
	Mode    DrawMode
	Glyph   Glyph
	Chars   [2]uint8
	Colours col.Pair
}

func NewGlyphVisuals(glyph Glyph, colours col.Pair) Visuals {
	return Visuals{
		Mode:    DRAW_GLYPH,
		Glyph:   glyph,
		Colours: colours,
	}
}

func NewTextVisuals(char1, char2 uint8, colours col.Pair) Visuals {
	return Visuals{
		Mode:    DRAW_TEXT,
		Colours: colours,
		Chars:   [2]uint8{char1, char2},
	}
}

func (v Visuals) String() string {
	switch v.Mode {
	case DRAW_GLYPH:
		return fmt.Sprintf("Vis{Mode: Glyph, G: %s, C: %s}", v.Glyph, v.Colours)
	case DRAW_TEXT:
		return fmt.Sprintf("Vis{Mode: Text, Txt: %v, C: %s}", v.Chars, v.Colours)
	case DRAW_NONE:
		return "Vis{Mode: None}"
	default:
		return "Vis{Unknown??}"
	}
}

// Changes the glyph. Also enables glyph drawmode.
func (v *Visuals) SetGlyph(glyph Glyph) {
	v.Glyph = glyph
	v.Mode = DRAW_GLYPH
}

// Changes the characters. Also enables text drawmode.
func (v *Visuals) SetText(char1, char2 uint8) {
	v.Chars[0] = char1
	v.Chars[1] = char2
	v.Mode = DRAW_TEXT
}

func (v Visuals) IsTransparent() bool {
	return v.Mode == DRAW_NONE || v.Colours.Back.IsTransparent()
}

func (v Visuals) Draw(dst_canvas Canvas, offset vec.Coord, depth int) {
	dst_canvas.DrawVisuals(offset, depth, v)
}

func (v Visuals) GetVisuals() Visuals {
	return v
}

// A glyph to be displayed. Tyumi fonts are based on the old codepage 437 fonts from ancient terminal applications,
// which had a whopping 256 characters to choose from.
type Glyph uint8

func (g Glyph) String() string {
	return GlyphNames[g]
}

// code page 437 glyphs
const (
	GLYPH_NONE Glyph = iota
	GLYPH_FACE1
	GLYPH_FACE2
	GLYPH_HEART
	GLYPH_DIAMOND
	GLYPH_CLUB
	GLYPH_SPADE
	GLYPH_PILL
	GLYPH_PILL_INVERSE
	GLYPH_DONUT
	GLYPH_DONUT_INVERSE //10
	GLYPH_MALE
	GLYPH_FEMALE
	GLYPH_MUSIC1
	GLYPH_MUSIC2
	GLYPH_STAR
	GLYPH_TRIANGLE_RIGHT
	GLYPH_TRIANGLE_LEFT
	GLYPH_ARROW_UPDOWN
	GLYPH_EXCLAMATION2
	GLYPH_PARAGRAPH //20
	GLYPH_SECTION
	GLYPH_LOWERCURSOR
	GLYPH_ARROW_UPDOWN_UNDERLINE
	GLYPH_ARROW_UP
	GLYPH_ARROW_DOWN
	GLYPH_ARROW_RIGHT
	GLYPH_ARROW_LEFT
	GLYPH_RETURN
	GLYPH_ARROW_LEFTRIGHT
	GLYPH_TRIANGLE_UP //30
	GLYPH_TRIANGLE_DOWN
	GLYPH_SPACE
	GLYPH_EXCLAMATION
	GLYPH_QUOTE
	GLYPH_HASH
	GLYPH_DOLLAR
	GLYPH_PERCENT
	GLYPH_AMPERSAND
	GLYPH_SINGLEQUOTE
	GLYPH_PARENTHESIS_LEFT //40
	GLYPH_PARENTHESIS_RIGHT
	GLYPH_ASTERISK
	GLYPH_PLUS
	GLYPH_COMMA
	GLYPH_MINUS
	GLYPH_PERIOD
	GLYPH_SLASH
	GLYPH_0
	GLYPH_1
	GLYPH_2 //50
	GLYPH_3
	GLYPH_4
	GLYPH_5
	GLYPH_6
	GLYPH_7
	GLYPH_8
	GLYPH_9
	GLYPH_COLON
	GLYPH_SEMICOLON
	GLYPH_ANGLEBRACE_LEFT //60
	GLYPH_EQUALS
	GLYPH_ANGLEBRACE_RIGHT
	GLYPH_QUESTION
	GLYPH_AT
	GLYPH_A
	GLYPH_B
	GLYPH_C
	GLYPH_D
	GLYPH_E
	GLYPH_F //70
	GLYPH_G
	GLYPH_H
	GLYPH_I
	GLYPH_J
	GLYPH_K
	GLYPH_L
	GLYPH_M
	GLYPH_N
	GLYPH_O
	GLYPH_P //80
	GLYPH_Q
	GLYPH_R
	GLYPH_S
	GLYPH_T
	GLYPH_U
	GLYPH_V
	GLYPH_W
	GLYPH_X
	GLYPH_Y
	GLYPH_Z //90
	GLYPH_SQUAREBRACE_LEFT
	GLYPH_BACKSLASH
	GLYPH_SQUAREBRACE_RIGHT
	GLYPH_CARET
	GLYPH_UNDERSCORE
	GLYPH_BACKTICK
	GLYPH_a
	GLYPH_b
	GLYPH_c
	GLYPH_d //100
	GLYPH_e
	GLYPH_f
	GLYPH_g
	GLYPH_h
	GLYPH_i
	GLYPH_j
	GLYPH_k
	GLYPH_l
	GLYPH_m
	GLYPH_n //110
	GLYPH_o
	GLYPH_p
	GLYPH_q
	GLYPH_r
	GLYPH_s
	GLYPH_t
	GLYPH_u
	GLYPH_v
	GLYPH_w
	GLYPH_x //120
	GLYPH_y
	GLYPH_z
	GLYPH_CURLYBRACE_LEFT
	GLYPH_BAR
	GLYPH_CURLYBRACE_RIGHT
	GLYPH_TILDE
	GLYPH_HOUSE
	GLYPH_C_CEDILLA
	GLYPH_u_UMLAUT
	GLYPH_e_ACUTE //130
	GLYPH_a_CIRCUMFLEX
	GLYPH_a_UMLAUT
	GLYPH_a_GRAVE
	GLYPH_a_RING
	GLYPH_c_CEDILLA
	GLYPH_e_CIRCUMFLEX
	GLYPH_e_UMLAUT
	GLYPH_e_GRAVE
	GLYPH_i_UMLAUT
	GLYPH_i_CIRCUMFLEX //140
	GLYPH_i_GRAVE
	GLYPH_A_UMLAUT
	GLYPH_A_RING
	GLYPH_E_ACUTE
	GLYPH_ae
	GLYPH_AE
	GLYPH_o_CIRCUMFLEX
	GLYPH_o_UMLAUT
	GLYPH_o_GRAVE
	GLYPH_u_CIRCUMFLEX //150
	GLYPH_u_GRAVE
	GLYPH_y_UMLAUT
	GLYPH_O_UMLAUT
	GLYPH_U_UMLAUT
	GLYPH_CENT
	GLYPH_POUND
	GLYPH_YEN
	GLYPH_PESETA
	GLYPH_FLORIN
	GLYPH_a_ACUTE //160
	GLYPH_i_ACUTE
	GLYPH_o_ACUTE
	GLYPH_u_ACUTE
	GLYPH_n_TILDE
	GLYPH_N_TILDE
	GLYPH_a_ORDINAL
	GLYPH_o_ORDINAL
	GLYPH_QUESTION_INVERSE
	GLYPH_NOT_REVERSE
	GLYPH_NOT //170
	GLYPH_HALF
	GLYPH_QUARTER
	GLYPH_EXCLAMATION_INVERSE
	GLYPH_DOUBLEANGLEBRACE_LEFT
	GLYPH_DOUBLEANGLEBRACE_RIGHT
	GLYPH_FILL_SPARSE
	GLYPH_FILL
	GLYPH_FILL_DENSE
	GLYPH_BORDER_UD
	GLYPH_BORDER_UDL //180
	GLYPH_BORDER_UDLL
	GLYPH_BORDER_UUDDL
	GLYPH_BORDER_DDL
	GLYPH_BORDER_DLL
	GLYPH_BORDER_UUDDLL
	GLYPH_BORDER_UUDD
	GLYPH_BORDER_DDLL
	GLYPH_BORDER_UULL
	GLYPH_BORDER_UUL
	GLYPH_BORDER_ULL //190
	GLYPH_BORDER_DL
	GLYPH_BORDER_UR
	GLYPH_BORDER_ULR
	GLYPH_BORDER_DLR
	GLYPH_BORDER_UDR
	GLYPH_BORDER_LR
	GLYPH_BORDER_UDLR
	GLYPH_BORDER_UDRR
	GLYPH_BORDER_UUDDR
	GLYPH_BORDER_UURR //200
	GLYPH_BORDER_DDRR
	GLYPH_BORDER_UULLRR
	GLYPH_BORDER_DDLLRR
	GLYPH_BORDER_UUDDRR
	GLYPH_BORDER_LLRR
	GLYPH_BORDER_UUDDLLRR
	GLYPH_BORDER_ULLRR
	GLYPH_BORDER_UULR
	GLYPH_BORDER_DLLRR
	GLYPH_BORDER_DDLR //210
	GLYPH_BORDER_UUR
	GLYPH_BORDER_URR
	GLYPH_BORDER_DRR
	GLYPH_BORDER_DDR
	GLYPH_BORDER_UUDDLR
	GLYPH_BORDER_UDLLRR
	GLYPH_BORDER_UL
	GLYPH_BORDER_DR
	GLYPH_BLOCK
	GLYPH_HALFBLOCK_DOWN //220
	GLYPH_HALFBLOCK_LEFT
	GLYPH_HALFBLOCK_RIGHT
	GLYPH_HALFBLOCK_UP
	GLYPH_ALPHA
	GLYPH_BETA
	GLYPH_GAMMA
	GLYPH_PI
	GLYPH_SIGMA
	GLYPH_SIGMA_LOWER
	GLYPH_MU //230
	GLYPH_TAU
	GLYPH_PHI
	GLYPH_THETA
	GLYPH_OMEGA
	GLYPH_DELTA
	GLYPH_INFINITY
	GLYPH_PHI_LOWER
	GLYPH_EPSILON
	GLYPH_NU
	GLYPH_IDENTICAL //240
	GLYPH_PLUSMINUS
	GLYPH_GREATEREQUAL
	GLYPH_LESSEQUAL
	GLYPH_INTEGRAL_TOP
	GLYPH_INTEGRAL_BOTTOM
	GLYPH_DIVIDE
	GLYPH_APPROX
	GLYPH_DEGREE
	GLYPH_DOT
	GLYPH_DOT_SMALL //250
	GLYPH_SQRT
	GLYPH_n_SUPERSCRIPT
	GLYPH_2_SUPERSCRIPT
	GLYPH_CURSOR
	GLYPH_BLANK
)

var GlyphNames = map[Glyph]string{
	GLYPH_NONE:                   "None",
	GLYPH_FACE1:                  "Face 1",
	GLYPH_FACE2:                  "Face 2",
	GLYPH_HEART:                  "Heart",
	GLYPH_DIAMOND:                "Diamond",
	GLYPH_CLUB:                   "Club",
	GLYPH_SPADE:                  "Spade",
	GLYPH_PILL:                   "Pill",
	GLYPH_PILL_INVERSE:           "Inverted Pill",
	GLYPH_DONUT:                  "Donut",
	GLYPH_DONUT_INVERSE:          "Inverted Donut",
	GLYPH_MALE:                   "Male",
	GLYPH_FEMALE:                 "Female",
	GLYPH_MUSIC1:                 "Music 1",
	GLYPH_MUSIC2:                 "Music 2",
	GLYPH_STAR:                   "Star",
	GLYPH_TRIANGLE_RIGHT:         "Triangle Right",
	GLYPH_TRIANGLE_LEFT:          "Triangle Left",
	GLYPH_ARROW_UPDOWN:           "Arrow Up-Down",
	GLYPH_EXCLAMATION2:           "Exclamation 2",
	GLYPH_PARAGRAPH:              "Paragraph",
	GLYPH_SECTION:                "Section",
	GLYPH_LOWERCURSOR:            "Lower Cursor",
	GLYPH_ARROW_UPDOWN_UNDERLINE: "Arrow Up-Down with Underline",
	GLYPH_ARROW_UP:               "Arrow Up",
	GLYPH_ARROW_DOWN:             "Arrow Down",
	GLYPH_ARROW_RIGHT:            "Arrow Right",
	GLYPH_ARROW_LEFT:             "Arrow Left",
	GLYPH_RETURN:                 "Return",
	GLYPH_ARROW_LEFTRIGHT:        "Arrow Left-Right",
	GLYPH_TRIANGLE_UP:            "Triangle Up",
	GLYPH_TRIANGLE_DOWN:          "Triangle Down",
	GLYPH_SPACE:                  "Space",
	GLYPH_EXCLAMATION:            "Exclamation",
	GLYPH_QUOTE:                  "Quote",
	GLYPH_HASH:                   "Hash",
	GLYPH_DOLLAR:                 "Dollar",
	GLYPH_PERCENT:                "Percent",
	GLYPH_AMPERSAND:              "Ampersand",
	GLYPH_SINGLEQUOTE:            "Single Quote",
	GLYPH_PARENTHESIS_LEFT:       "Paren Left",
	GLYPH_PARENTHESIS_RIGHT:      "Paren Right",
	GLYPH_ASTERISK:               "Asterisk",
	GLYPH_PLUS:                   "Plus",
	GLYPH_COMMA:                  "Comma",
	GLYPH_MINUS:                  "Minus",
	GLYPH_PERIOD:                 "Period",
	GLYPH_SLASH:                  "Slash",
	GLYPH_0:                      "0",
	GLYPH_1:                      "1",
	GLYPH_2:                      "2",
	GLYPH_3:                      "3",
	GLYPH_4:                      "4",
	GLYPH_5:                      "5",
	GLYPH_6:                      "6",
	GLYPH_7:                      "7",
	GLYPH_8:                      "8",
	GLYPH_9:                      "9",
	GLYPH_COLON:                  "Colon",
	GLYPH_SEMICOLON:              "Semicolon",
	GLYPH_ANGLEBRACE_LEFT:        "Angled Brace Left",
	GLYPH_EQUALS:                 "Equals",
	GLYPH_ANGLEBRACE_RIGHT:       "Angled Brace Right",
	GLYPH_QUESTION:               "Question Mark",
	GLYPH_AT:                     "At Sign",
	GLYPH_A:                      "A",
	GLYPH_B:                      "B",
	GLYPH_C:                      "C",
	GLYPH_D:                      "D",
	GLYPH_E:                      "E",
	GLYPH_F:                      "F",
	GLYPH_G:                      "G",
	GLYPH_H:                      "H",
	GLYPH_I:                      "I",
	GLYPH_J:                      "J",
	GLYPH_K:                      "K",
	GLYPH_L:                      "L",
	GLYPH_M:                      "M",
	GLYPH_N:                      "N",
	GLYPH_O:                      "O",
	GLYPH_P:                      "P",
	GLYPH_Q:                      "Q",
	GLYPH_R:                      "R",
	GLYPH_S:                      "S",
	GLYPH_T:                      "T",
	GLYPH_U:                      "U",
	GLYPH_V:                      "V",
	GLYPH_W:                      "W",
	GLYPH_X:                      "X",
	GLYPH_Y:                      "Y",
	GLYPH_Z:                      "Z",
	GLYPH_SQUAREBRACE_LEFT:       "Square Brace Left",
	GLYPH_BACKSLASH:              "Backslash",
	GLYPH_SQUAREBRACE_RIGHT:      "Square Brace Right",
	GLYPH_CARET:                  "Caret",
	GLYPH_UNDERSCORE:             "Underscore",
	GLYPH_BACKTICK:               "Backtick",
	GLYPH_a:                      "a",
	GLYPH_b:                      "b",
	GLYPH_c:                      "c",
	GLYPH_d:                      "d",
	GLYPH_e:                      "e",
	GLYPH_f:                      "f",
	GLYPH_g:                      "g",
	GLYPH_h:                      "h",
	GLYPH_i:                      "i",
	GLYPH_j:                      "j",
	GLYPH_k:                      "k",
	GLYPH_l:                      "l",
	GLYPH_m:                      "m",
	GLYPH_n:                      "n",
	GLYPH_o:                      "o",
	GLYPH_p:                      "p",
	GLYPH_q:                      "q",
	GLYPH_r:                      "r",
	GLYPH_s:                      "s",
	GLYPH_t:                      "t",
	GLYPH_u:                      "u",
	GLYPH_v:                      "v",
	GLYPH_w:                      "w",
	GLYPH_x:                      "x",
	GLYPH_y:                      "y",
	GLYPH_z:                      "z",
	GLYPH_CURLYBRACE_LEFT:        "Curly Brace Left",
	GLYPH_BAR:                    "Vertical Bar",
	GLYPH_CURLYBRACE_RIGHT:       "Curly Brace Right",
	GLYPH_TILDE:                  "Tilde",
	GLYPH_HOUSE:                  "House",
	GLYPH_C_CEDILLA:              "C with Cedilla",
	GLYPH_u_UMLAUT:               "u with Unlaut",
	GLYPH_e_ACUTE:                "e with Accent Acute",
	GLYPH_a_CIRCUMFLEX:           "a with Accent Circumflex",
	GLYPH_a_UMLAUT:               "a with Umlaut",
	GLYPH_a_GRAVE:                "a with Accent Grave",
	GLYPH_a_RING:                 "a with Ring",
	GLYPH_c_CEDILLA:              "c with Cedilla",
	GLYPH_e_CIRCUMFLEX:           "e with Accent Circumflex",
	GLYPH_e_UMLAUT:               "e with Umlaut",
	GLYPH_e_GRAVE:                "e with Accent Grave",
	GLYPH_i_UMLAUT:               "i with Umlaut",
	GLYPH_i_CIRCUMFLEX:           "i with Accent Circumflex",
	GLYPH_i_GRAVE:                "i with Accent Grave",
	GLYPH_A_UMLAUT:               "A with Umlaut",
	GLYPH_A_RING:                 "A with Ring",
	GLYPH_E_ACUTE:                "E with Accent Acute",
	GLYPH_ae:                     "ae",
	GLYPH_AE:                     "AE",
	GLYPH_o_CIRCUMFLEX:           "o with Accent Circumflex",
	GLYPH_o_UMLAUT:               "o with Umlaut",
	GLYPH_o_GRAVE:                "o with Accent Grave",
	GLYPH_u_CIRCUMFLEX:           "u with Accent Circumflex",
	GLYPH_u_GRAVE:                "u with Accent Grave",
	GLYPH_y_UMLAUT:               "y with Umlaut",
	GLYPH_O_UMLAUT:               "O with Umlaut",
	GLYPH_U_UMLAUT:               "U with Umlaut",
	GLYPH_CENT:                   "Cent",
	GLYPH_POUND:                  "Pound",
	GLYPH_YEN:                    "Yen",
	GLYPH_PESETA:                 "Peseta",
	GLYPH_FLORIN:                 "Florin",
	GLYPH_a_ACUTE:                "a with Accent Acute",
	GLYPH_i_ACUTE:                "i with Accent Acute",
	GLYPH_o_ACUTE:                "o with Accent Acute",
	GLYPH_u_ACUTE:                "u with Accent Acute",
	GLYPH_n_TILDE:                "n with Tilde",
	GLYPH_N_TILDE:                "N with Tilde",
	GLYPH_a_ORDINAL:              "a Ordinal",
	GLYPH_o_ORDINAL:              "o Ordinal",
	GLYPH_QUESTION_INVERSE:       "Inverted Question Mark",
	GLYPH_NOT_REVERSE:            "Reversed Boolean Not",
	GLYPH_NOT:                    "Boolean Not",
	GLYPH_HALF:                   "Half",
	GLYPH_QUARTER:                "Quarter",
	GLYPH_EXCLAMATION_INVERSE:    "Inverted Exclamation Mark",
	GLYPH_DOUBLEANGLEBRACE_LEFT:  "Double Angled Brace Left",
	GLYPH_DOUBLEANGLEBRACE_RIGHT: "Double Angled Brace Right",
	GLYPH_FILL_SPARSE:            "Fill (Sparse)",
	GLYPH_FILL:                   "Fill",
	GLYPH_FILL_DENSE:             "Fill (Dense)",
	GLYPH_BORDER_UD:              "Border UD",
	GLYPH_BORDER_UDL:             "Border (UDL)",
	GLYPH_BORDER_UDLL:            "Border (UDLL)",
	GLYPH_BORDER_UUDDL:           "Border (UUDDL)",
	GLYPH_BORDER_DDL:             "Border (DDL)",
	GLYPH_BORDER_DLL:             "Border (DLL)",
	GLYPH_BORDER_UUDDLL:          "Border (UUDDLL)",
	GLYPH_BORDER_UUDD:            "Border (UUDD)",
	GLYPH_BORDER_DDLL:            "Border (DDLL)",
	GLYPH_BORDER_UULL:            "Border (UULL)",
	GLYPH_BORDER_UUL:             "Border (UUL)",
	GLYPH_BORDER_ULL:             "Border (ULL)",
	GLYPH_BORDER_DL:              "Border (DL)",
	GLYPH_BORDER_UR:              "Border (UR)",
	GLYPH_BORDER_ULR:             "Border (ULR)",
	GLYPH_BORDER_DLR:             "Border (DLR)",
	GLYPH_BORDER_UDR:             "Border (UDR)",
	GLYPH_BORDER_LR:              "Border (LR)",
	GLYPH_BORDER_UDLR:            "Border (UDLR)",
	GLYPH_BORDER_UDRR:            "Border (UDRR)",
	GLYPH_BORDER_UUDDR:           "Border (UUDDR)",
	GLYPH_BORDER_UURR:            "Border (UURR)",
	GLYPH_BORDER_DDRR:            "Border (DDRR)",
	GLYPH_BORDER_UULLRR:          "Border (UULLRR)",
	GLYPH_BORDER_DDLLRR:          "Border (DDLLRR)",
	GLYPH_BORDER_UUDDRR:          "Border (UUDDRR)",
	GLYPH_BORDER_LLRR:            "Border (LLRR)",
	GLYPH_BORDER_UUDDLLRR:        "Border (UUDDLLRR)",
	GLYPH_BORDER_ULLRR:           "Border (ULLRR)",
	GLYPH_BORDER_UULR:            "Border (UULR)",
	GLYPH_BORDER_DLLRR:           "Border (DLLRR)",
	GLYPH_BORDER_DDLR:            "Border (DDLR)",
	GLYPH_BORDER_UUR:             "Border (UUR)",
	GLYPH_BORDER_URR:             "Border (URR)",
	GLYPH_BORDER_DRR:             "Border (DRR)",
	GLYPH_BORDER_DDR:             "Border (DDR)",
	GLYPH_BORDER_UUDDLR:          "Border (UUDDLR)",
	GLYPH_BORDER_UDLLRR:          "Border (UDLLRR)",
	GLYPH_BORDER_UL:              "Border (UL)",
	GLYPH_BORDER_DR:              "Border (DR)",
	GLYPH_BLOCK:                  "Block",
	GLYPH_HALFBLOCK_DOWN:         "Half Block (Down)",
	GLYPH_HALFBLOCK_LEFT:         "Half Block (Left)",
	GLYPH_HALFBLOCK_RIGHT:        "Half Block (Right)",
	GLYPH_HALFBLOCK_UP:           "Half Block (Up)",
	GLYPH_ALPHA:                  "Alpha",
	GLYPH_BETA:                   "Beta",
	GLYPH_GAMMA:                  "Gamma",
	GLYPH_PI:                     "Pi",
	GLYPH_SIGMA:                  "Sigma",
	GLYPH_SIGMA_LOWER:            "Sigma (lower case)",
	GLYPH_MU:                     "Mu",
	GLYPH_TAU:                    "Tau",
	GLYPH_PHI:                    "Phi",
	GLYPH_THETA:                  "Theta",
	GLYPH_OMEGA:                  "Omega",
	GLYPH_DELTA:                  "Delta",
	GLYPH_INFINITY:               "Infinity",
	GLYPH_PHI_LOWER:              "Phi (lower case)",
	GLYPH_EPSILON:                "Epsilon",
	GLYPH_NU:                     "Nu",
	GLYPH_IDENTICAL:              "Identical",
	GLYPH_PLUSMINUS:              "Plus/Minus",
	GLYPH_GREATEREQUAL:           "Greater Than Or Equal",
	GLYPH_LESSEQUAL:              "Less than or Equal",
	GLYPH_INTEGRAL_TOP:           "Integral (top)",
	GLYPH_INTEGRAL_BOTTOM:        "Integral (bottom)",
	GLYPH_DIVIDE:                 "Divide",
	GLYPH_APPROX:                 "Approximate",
	GLYPH_DEGREE:                 "Degree",
	GLYPH_DOT:                    "Dot",
	GLYPH_DOT_SMALL:              "Small Dot",
	GLYPH_SQRT:                   "Square Root",
	GLYPH_n_SUPERSCRIPT:          "Superscript n",
	GLYPH_2_SUPERSCRIPT:          "Superscript 2",
	GLYPH_CURSOR:                 "Cursor",
	GLYPH_BLANK:                  "Blank",
}

// Special text characters.
const (
	TEXT_BORDER_LR         uint8 = 196
	TEXT_BORDER_UD         uint8 = 179
	TEXT_BORDER_DECO_LEFT  uint8 = 180
	TEXT_BORDER_DECO_RIGHT uint8 = 195
	TEXT_DEFAULT           uint8 = 255 //indicates that text character should preserve what was there previously
	TEXT_NONE              uint8 = 32  //just a space
)
