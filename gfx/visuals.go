package gfx

//Defines anything with the ability to be drawn.
type Drawable interface {
	Visuals() Visuals
}

//The basic visual definition of a single-tile object that can be drawn to the screen.
type Visuals struct {
	Glyph      int
	ForeColour uint32
	BackColour uint32
}

func (v *Visuals) ChangeGlyph(g int) {
	v.Glyph = g
}

func (v *Visuals) ChangeForeColour(f uint32) {
	v.ForeColour = f
}

func (v *Visuals) ChangeBackColour(b uint32) {
	v.BackColour = b
}

func (v Visuals) GetVisuals() Visuals {
	return v
}

//code page 437 glyphs
const (
	GLYPH_NONE int = iota
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
	MAXGLYPHS
)

//Special text characters. Add these to text strings, it'll just work!
const (
	TEXT_BORDER_LR         string = string(rune(196))
	TEXT_BORDER_UD         string = string(rune(179))
	TEXT_BORDER_DECO_LEFT  string = string(rune(180))
	TEXT_BORDER_DECO_RIGHT string = string(rune(195))
)
