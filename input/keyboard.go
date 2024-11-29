package input

import (
	"strconv"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/util"
	"github.com/veandco/go-sdl2/sdl"
)

type KeyboardEvent struct {
	event.EventPrototype

	Key Keycode
}

func NewKeyboardEvent(key Keycode) (kbe KeyboardEvent) {
	kbe.EventPrototype = event.New(EV_KEYBOARD)
	kbe.Key = key
	return
}

// If the keyboard event represents a direction. returns the x and y deltas for the direction.
func (kb KeyboardEvent) Direction() (dx, dy int) {
	switch kb.Key {
	case K_UP:
		return 0, -1
	case K_DOWN:
		return 0, 1
	case K_LEFT:
		return -1, 0
	case K_RIGHT:
		return 1, 0
	}

	return
}

// If the keyboard event represents text in some way (letter, number, anything traditionally
// representable by a string) this returns the string form. If not, it return an empty string.
// TODO: support for capital letters once modifier support is added.
func (kb KeyboardEvent) Text() (s string) {
	key := rune(kb.Key)
	if util.ValidText(key) {
		s = strconv.QuoteRuneToASCII(key)
		s, _ = strconv.Unquote(s)
	} else if kb.Key == K_SPACE {
		s = " "
	}

	return
}

// Keycodes. Once there are multiple input handlers instead of just sdl, we'll need to have a
// layer that translates key events into these keycodes. For now, the keycodes are just the same
// as the sdl keycodes.
type Keycode int

// keycodes. these names and definitions and names are ripped right from
// github.com/veandco/go-sdl2/sdl/keycode.go
const (
	K_UNKNOWN      = Keycode(sdl.K_UNKNOWN)      // "" (no name, empty string)
	K_RETURN       = Keycode(sdl.K_RETURN)       // "Return" (the Enter key (main keyboard))
	K_ESCAPE       = Keycode(sdl.K_ESCAPE)       // "Escape" (the Esc key)
	K_BACKSPACE    = Keycode(sdl.K_BACKSPACE)    // "Backspace"
	K_TAB          = Keycode(sdl.K_TAB)          // "Tab" (the Tab key)
	K_SPACE        = Keycode(sdl.K_SPACE)        // "Space" (the Space Bar key(s))
	K_EXCLAIM      = Keycode(sdl.K_EXCLAIM)      // "!"
	K_QUOTEDBL     = Keycode(sdl.K_QUOTEDBL)     // """
	K_HASH         = Keycode(sdl.K_HASH)         // "#"
	K_PERCENT      = Keycode(sdl.K_PERCENT)      // "%"
	K_DOLLAR       = Keycode(sdl.K_DOLLAR)       // "$"
	K_AMPERSAND    = Keycode(sdl.K_AMPERSAND)    // "&"
	K_QUOTE        = Keycode(sdl.K_QUOTE)        // "'"
	K_LEFTPAREN    = Keycode(sdl.K_LEFTPAREN)    // "("
	K_RIGHTPAREN   = Keycode(sdl.K_RIGHTPAREN)   // ")"
	K_ASTERISK     = Keycode(sdl.K_ASTERISK)     // "*"
	K_PLUS         = Keycode(sdl.K_PLUS)         // "+"
	K_COMMA        = Keycode(sdl.K_COMMA)        // ","
	K_MINUS        = Keycode(sdl.K_MINUS)        // "-"
	K_PERIOD       = Keycode(sdl.K_PERIOD)       // "."
	K_SLASH        = Keycode(sdl.K_SLASH)        // "/"
	K_0            = Keycode(sdl.K_0)            // "0"
	K_1            = Keycode(sdl.K_1)            // "1"
	K_2            = Keycode(sdl.K_2)            // "2"
	K_3            = Keycode(sdl.K_3)            // "3"
	K_4            = Keycode(sdl.K_4)            // "4"
	K_5            = Keycode(sdl.K_5)            // "5"
	K_6            = Keycode(sdl.K_6)            // "6"
	K_7            = Keycode(sdl.K_7)            // "7"
	K_8            = Keycode(sdl.K_8)            // "8"
	K_9            = Keycode(sdl.K_9)            // "9"
	K_COLON        = Keycode(sdl.K_COLON)        // ":"
	K_SEMICOLON    = Keycode(sdl.K_SEMICOLON)    // ";"
	K_LESS         = Keycode(sdl.K_LESS)         // "<"
	K_EQUALS       = Keycode(sdl.K_EQUALS)       // "="
	K_GREATER      = Keycode(sdl.K_GREATER)      // ">"
	K_QUESTION     = Keycode(sdl.K_QUESTION)     // "?"
	K_AT           = Keycode(sdl.K_AT)           // "@"
	K_LEFTBRACKET  = Keycode(sdl.K_LEFTBRACKET)  // "["
	K_BACKSLASH    = Keycode(sdl.K_BACKSLASH)    // "\"
	K_RIGHTBRACKET = Keycode(sdl.K_RIGHTBRACKET) // "]"
	K_CARET        = Keycode(sdl.K_CARET)        // "^"
	K_UNDERSCORE   = Keycode(sdl.K_UNDERSCORE)   // "_"
	K_BACKQUOTE    = Keycode(sdl.K_BACKQUOTE)    // "`"
	K_a            = Keycode(sdl.K_a)            // "A"
	K_b            = Keycode(sdl.K_b)            // "B"
	K_c            = Keycode(sdl.K_c)            // "C"
	K_d            = Keycode(sdl.K_d)            // "D"
	K_e            = Keycode(sdl.K_e)            // "E"
	K_f            = Keycode(sdl.K_f)            // "F"
	K_g            = Keycode(sdl.K_g)            // "G"
	K_h            = Keycode(sdl.K_h)            // "H"
	K_i            = Keycode(sdl.K_i)            // "I"
	K_j            = Keycode(sdl.K_j)            // "J"
	K_k            = Keycode(sdl.K_k)            // "K"
	K_l            = Keycode(sdl.K_l)            // "L"
	K_m            = Keycode(sdl.K_m)            // "M"
	K_n            = Keycode(sdl.K_n)            // "N"
	K_o            = Keycode(sdl.K_o)            // "O"
	K_p            = Keycode(sdl.K_p)            // "P"
	K_q            = Keycode(sdl.K_q)            // "Q"
	K_r            = Keycode(sdl.K_r)            // "R"
	K_s            = Keycode(sdl.K_s)            // "S"
	K_t            = Keycode(sdl.K_t)            // "T"
	K_u            = Keycode(sdl.K_u)            // "U"
	K_v            = Keycode(sdl.K_v)            // "V"
	K_w            = Keycode(sdl.K_w)            // "W"
	K_x            = Keycode(sdl.K_x)            // "X"
	K_y            = Keycode(sdl.K_y)            // "Y"
	K_z            = Keycode(sdl.K_z)            // "Z"
	K_CAPSLOCK     = Keycode(sdl.K_CAPSLOCK)     // "CapsLock"
	K_F1           = Keycode(sdl.K_F1)           // "F1"
	K_F2           = Keycode(sdl.K_F2)           // "F2"
	K_F3           = Keycode(sdl.K_F3)           // "F3"
	K_F4           = Keycode(sdl.K_F4)           // "F4"
	K_F5           = Keycode(sdl.K_F5)           // "F5"
	K_F6           = Keycode(sdl.K_F6)           // "F6"
	K_F7           = Keycode(sdl.K_F7)           // "F7"
	K_F8           = Keycode(sdl.K_F8)           // "F8"
	K_F9           = Keycode(sdl.K_F9)           // "F9"
	K_F10          = Keycode(sdl.K_F10)          // "F10"
	K_F11          = Keycode(sdl.K_F11)          // "F11"
	K_F12          = Keycode(sdl.K_F12)          // "F12"
	K_PRINTSCREEN  = Keycode(sdl.K_PRINTSCREEN)  // "PrintScreen"
	K_SCROLLLOCK   = Keycode(sdl.K_SCROLLLOCK)   // "ScrollLock"
	K_PAUSE        = Keycode(sdl.K_PAUSE)        // "Pause" (the Pause / Break key)
	K_INSERT       = Keycode(sdl.K_INSERT)       // "Insert" (insert on PC, help on some Mac keyboards (but does send code 73, not 117))
	K_HOME         = Keycode(sdl.K_HOME)         // "Home"
	K_PAGEUP       = Keycode(sdl.K_PAGEUP)       // "PageUp"
	K_DELETE       = Keycode(sdl.K_DELETE)       // "Delete"
	K_END          = Keycode(sdl.K_END)          // "End"
	K_PAGEDOWN     = Keycode(sdl.K_PAGEDOWN)     // "PageDown"
	K_RIGHT        = Keycode(sdl.K_RIGHT)        // "Right" (the Right arrow key (navigation keypad))
	K_LEFT         = Keycode(sdl.K_LEFT)         // "Left" (the Left arrow key (navigation keypad))
	K_DOWN         = Keycode(sdl.K_DOWN)         // "Down" (the Down arrow key (navigation keypad))
	K_UP           = Keycode(sdl.K_UP)           // "Up" (the Up arrow key (navigation keypad))
	K_NUMLOCKCLEAR = Keycode(sdl.K_NUMLOCKCLEAR) // "Numlock" (the Num Lock key (PC) / the Clear key (Mac))
	K_KP_DIVIDE    = Keycode(sdl.K_KP_DIVIDE)    // "Keypad /" (the / key (numeric keypad))
	K_KP_MULTIPLY  = Keycode(sdl.K_KP_MULTIPLY)  // "Keypad *" (the * key (numeric keypad))
	K_KP_MINUS     = Keycode(sdl.K_KP_MINUS)     // "Keypad -" (the - key (numeric keypad))
	K_KP_PLUS      = Keycode(sdl.K_KP_PLUS)      // "Keypad +" (the + key (numeric keypad))
	K_KP_ENTER     = Keycode(sdl.K_KP_ENTER)     // "Keypad Enter" (the Enter key (numeric keypad))
	K_KP_1         = Keycode(sdl.K_KP_1)         // "Keypad 1" (the 1 key (numeric keypad))
	K_KP_2         = Keycode(sdl.K_KP_2)         // "Keypad 2" (the 2 key (numeric keypad))
	K_KP_3         = Keycode(sdl.K_KP_3)         // "Keypad 3" (the 3 key (numeric keypad))
	K_KP_4         = Keycode(sdl.K_KP_4)         // "Keypad 4" (the 4 key (numeric keypad))
	K_KP_5         = Keycode(sdl.K_KP_5)         // "Keypad 5" (the 5 key (numeric keypad))
	K_KP_6         = Keycode(sdl.K_KP_6)         // "Keypad 6" (the 6 key (numeric keypad))
	K_KP_7         = Keycode(sdl.K_KP_7)         // "Keypad 7" (the 7 key (numeric keypad))
	K_KP_8         = Keycode(sdl.K_KP_8)         // "Keypad 8" (the 8 key (numeric keypad))
	K_KP_9         = Keycode(sdl.K_KP_9)         // "Keypad 9" (the 9 key (numeric keypad))
	K_KP_0         = Keycode(sdl.K_KP_0)         // "Keypad 0" (the 0 key (numeric keypad))
	K_KP_PERIOD    = Keycode(sdl.K_KP_PERIOD)    // "Keypad ." (the . key (numeric keypad))
)
