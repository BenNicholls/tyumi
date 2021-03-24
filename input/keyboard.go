package input

import (
	"github.com/bennicholls/tyumi/event"
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

//If the keyboard event represents a direction. returns the x and y deltas for the direction.
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

//Keycodes. Once there are multiple input handlers instead of just sdl, we'll need to have a
//layer that translates key events into these keycodes. For now, the keycodes are just the same
//as the sdl keycodes.
type Keycode int

//keycodes. these names and definitions and names are ripped right from
//github.com/veandco/go-sdl2/sdl/keycode.go
const (
	K_UNKNOWN      = sdl.K_UNKNOWN      // "" (no name, empty string)
	K_RETURN       = sdl.K_RETURN       // "Return" (the Enter key (main keyboard))
	K_ESCAPE       = sdl.K_ESCAPE       // "Escape" (the Esc key)
	K_BACKSPACE    = sdl.K_BACKSPACE    // "Backspace"
	K_TAB          = sdl.K_TAB          // "Tab" (the Tab key)
	K_SPACE        = sdl.K_SPACE        // "Space" (the Space Bar key(s))
	K_EXCLAIM      = sdl.K_EXCLAIM      // "!"
	K_QUOTEDBL     = sdl.K_QUOTEDBL     // """
	K_HASH         = sdl.K_HASH         // "#"
	K_PERCENT      = sdl.K_PERCENT      // "%"
	K_DOLLAR       = sdl.K_DOLLAR       // "$"
	K_AMPERSAND    = sdl.K_AMPERSAND    // "&"
	K_QUOTE        = sdl.K_QUOTE        // "'"
	K_LEFTPAREN    = sdl.K_LEFTPAREN    // "("
	K_RIGHTPAREN   = sdl.K_RIGHTPAREN   // ")"
	K_ASTERISK     = sdl.K_ASTERISK     // "*"
	K_PLUS         = sdl.K_PLUS         // "+"
	K_COMMA        = sdl.K_COMMA        // ","
	K_MINUS        = sdl.K_MINUS        // "-"
	K_PERIOD       = sdl.K_PERIOD       // "."
	K_SLASH        = sdl.K_SLASH        // "/"
	K_0            = sdl.K_0            // "0"
	K_1            = sdl.K_1            // "1"
	K_2            = sdl.K_2            // "2"
	K_3            = sdl.K_3            // "3"
	K_4            = sdl.K_4            // "4"
	K_5            = sdl.K_5            // "5"
	K_6            = sdl.K_6            // "6"
	K_7            = sdl.K_7            // "7"
	K_8            = sdl.K_8            // "8"
	K_9            = sdl.K_9            // "9"
	K_COLON        = sdl.K_COLON        // ":"
	K_SEMICOLON    = sdl.K_SEMICOLON    // ";"
	K_LESS         = sdl.K_LESS         // "<"
	K_EQUALS       = sdl.K_EQUALS       // "="
	K_GREATER      = sdl.K_GREATER      // ">"
	K_QUESTION     = sdl.K_QUESTION     // "?"
	K_AT           = sdl.K_AT           // "@"
	K_LEFTBRACKET  = sdl.K_LEFTBRACKET  // "["
	K_BACKSLASH    = sdl.K_BACKSLASH    // "\"
	K_RIGHTBRACKET = sdl.K_RIGHTBRACKET // "]"
	K_CARET        = sdl.K_CARET        // "^"
	K_UNDERSCORE   = sdl.K_UNDERSCORE   // "_"
	K_BACKQUOTE    = sdl.K_BACKQUOTE    // "`"
	K_a            = sdl.K_a            // "A"
	K_b            = sdl.K_b            // "B"
	K_c            = sdl.K_c            // "C"
	K_d            = sdl.K_d            // "D"
	K_e            = sdl.K_e            // "E"
	K_f            = sdl.K_f            // "F"
	K_g            = sdl.K_g            // "G"
	K_h            = sdl.K_h            // "H"
	K_i            = sdl.K_i            // "I"
	K_j            = sdl.K_j            // "J"
	K_k            = sdl.K_k            // "K"
	K_l            = sdl.K_l            // "L"
	K_m            = sdl.K_m            // "M"
	K_n            = sdl.K_n            // "N"
	K_o            = sdl.K_o            // "O"
	K_p            = sdl.K_p            // "P"
	K_q            = sdl.K_q            // "Q"
	K_r            = sdl.K_r            // "R"
	K_s            = sdl.K_s            // "S"
	K_t            = sdl.K_t            // "T"
	K_u            = sdl.K_u            // "U"
	K_v            = sdl.K_v            // "V"
	K_w            = sdl.K_w            // "W"
	K_x            = sdl.K_x            // "X"
	K_y            = sdl.K_y            // "Y"
	K_z            = sdl.K_z            // "Z"
	K_CAPSLOCK     = sdl.K_CAPSLOCK     // "CapsLock"
	K_F1           = sdl.K_F1           // "F1"
	K_F2           = sdl.K_F2           // "F2"
	K_F3           = sdl.K_F3           // "F3"
	K_F4           = sdl.K_F4           // "F4"
	K_F5           = sdl.K_F5           // "F5"
	K_F6           = sdl.K_F6           // "F6"
	K_F7           = sdl.K_F7           // "F7"
	K_F8           = sdl.K_F8           // "F8"
	K_F9           = sdl.K_F9           // "F9"
	K_F10          = sdl.K_F10          // "F10"
	K_F11          = sdl.K_F11          // "F11"
	K_F12          = sdl.K_F12          // "F12"
	K_PRINTSCREEN  = sdl.K_PRINTSCREEN  // "PrintScreen"
	K_SCROLLLOCK   = sdl.K_SCROLLLOCK   // "ScrollLock"
	K_PAUSE        = sdl.K_PAUSE        // "Pause" (the Pause / Break key)
	K_INSERT       = sdl.K_INSERT       // "Insert" (insert on PC, help on some Mac keyboards (but does send code 73, not 117))
	K_HOME         = sdl.K_HOME         // "Home"
	K_PAGEUP       = sdl.K_PAGEUP       // "PageUp"
	K_DELETE       = sdl.K_DELETE       // "Delete"
	K_END          = sdl.K_END          // "End"
	K_PAGEDOWN     = sdl.K_PAGEDOWN     // "PageDown"
	K_RIGHT        = sdl.K_RIGHT        // "Right" (the Right arrow key (navigation keypad))
	K_LEFT         = sdl.K_LEFT         // "Left" (the Left arrow key (navigation keypad))
	K_DOWN         = sdl.K_DOWN         // "Down" (the Down arrow key (navigation keypad))
	K_UP           = sdl.K_UP           // "Up" (the Up arrow key (navigation keypad))
	K_NUMLOCKCLEAR = sdl.K_NUMLOCKCLEAR // "Numlock" (the Num Lock key (PC) / the Clear key (Mac))
	K_KP_DIVIDE    = sdl.K_KP_DIVIDE    // "Keypad /" (the / key (numeric keypad))
	K_KP_MULTIPLY  = sdl.K_KP_MULTIPLY  // "Keypad *" (the * key (numeric keypad))
	K_KP_MINUS     = sdl.K_KP_MINUS     // "Keypad -" (the - key (numeric keypad))
	K_KP_PLUS      = sdl.K_KP_PLUS      // "Keypad +" (the + key (numeric keypad))
	K_KP_ENTER     = sdl.K_KP_ENTER     // "Keypad Enter" (the Enter key (numeric keypad))
	K_KP_1         = sdl.K_KP_1         // "Keypad 1" (the 1 key (numeric keypad))
	K_KP_2         = sdl.K_KP_2         // "Keypad 2" (the 2 key (numeric keypad))
	K_KP_3         = sdl.K_KP_3         // "Keypad 3" (the 3 key (numeric keypad))
	K_KP_4         = sdl.K_KP_4         // "Keypad 4" (the 4 key (numeric keypad))
	K_KP_5         = sdl.K_KP_5         // "Keypad 5" (the 5 key (numeric keypad))
	K_KP_6         = sdl.K_KP_6         // "Keypad 6" (the 6 key (numeric keypad))
	K_KP_7         = sdl.K_KP_7         // "Keypad 7" (the 7 key (numeric keypad))
	K_KP_8         = sdl.K_KP_8         // "Keypad 8" (the 8 key (numeric keypad))
	K_KP_9         = sdl.K_KP_9         // "Keypad 9" (the 9 key (numeric keypad))
	K_KP_0         = sdl.K_KP_0         // "Keypad 0" (the 0 key (numeric keypad))
	K_KP_PERIOD    = sdl.K_KP_PERIOD    // "Keypad ." (the . key (numeric keypad))
)
