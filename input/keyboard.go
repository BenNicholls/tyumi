package input

import (
	"strconv"
	"strings"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// KEY_PRESSED for key downs, KEY_RELEASED for key ups.
type KeyPressType uint8

const (
	KEY_PRESSED KeyPressType = iota
	KEY_RELEASED

	KEYPRESS_EITHER // only used for defining allowable presstypes for actions.
)

// Flags representing modifier keys. These will be OR'd together in the case that multiple are pressed concurrently.
type KeyModifiers uint8

const (
	KEYMOD_NONE  KeyModifiers = 0
	KEYMOD_CTRL  KeyModifiers = 0b001
	KEYMOD_ALT   KeyModifiers = 0b010
	KEYMOD_SHIFT KeyModifiers = 0b100

	KEYMOD_CTRLSHIFT = KEYMOD_CTRL | KEYMOD_SHIFT
	KEYMOD_CTRLALT   = KEYMOD_CTRL | KEYMOD_ALT
	KEYMOD_ALTSHIFT  = KEYMOD_ALT | KEYMOD_SHIFT
)

type KeyboardEvent struct {
	event.EventPrototype

	Key       Keycode
	PressType KeyPressType
	Mods      KeyModifiers // Modifier keys for this key event, ex. CTRL, ALT, SHIFT.
	Repeat    bool         //will be true if this is the key is being held down
}

func fireKeyboardEvent(key_event KeyboardEvent) {
	event.Fire(EV_KEYBOARD, &key_event)
	key_event.fireActions()
}

func (kbe KeyboardEvent) fireActions() {
	var triggeredActions util.Set[ActionID]
	if triggerSet, ok := DefaultActionMap.keyTriggers[kbe.Key]; ok {
		for trigger := range triggerSet.EachElement() {
			if trigger.TriggeredBy(kbe) {
				triggeredActions.Add(trigger.Action)
			}
		}
	}

	for action := range triggeredActions.EachElement() {
		fireActionEvent(action)
	}
}

// Emits keypress event.
func FireKeyPressEvent(key Keycode, mods ...KeyModifiers) {
	fireKeyboardEvent(KeyboardEvent{
		Key:  key,
		Mods: util.OrAll(mods),
	})
}

// Emits keyrelease event.
func FireKeyReleaseEvent(key Keycode, mods ...KeyModifiers) {
	if SuppressKeyUpEvents {
		return
	}

	fireKeyboardEvent(KeyboardEvent{
		Key:       key,
		Mods:      util.OrAll(mods),
		PressType: KEY_RELEASED,
	})
}

// Emits key repeated event. The KeyPressType of repeat events is always KEY_PRESSED.
func FireKeyRepeatEvent(key Keycode, mods ...KeyModifiers) {
	if !AllowKeyRepeats {
		return
	}

	fireKeyboardEvent(KeyboardEvent{
		Key:    key,
		Mods:   util.OrAll(mods),
		Repeat: true,
	})
}

// Returns a vec.Direction if the keyboard event represents a direction (or vec.DIR_NONE if not).
func (kb KeyboardEvent) Direction() vec.Direction {
	switch kb.Key {
	case K_UP, K_KP_8:
		return vec.DIR_UP
	case K_DOWN, K_KP_2:
		return vec.DIR_DOWN
	case K_LEFT, K_KP_4:
		return vec.DIR_LEFT
	case K_RIGHT, K_KP_6:
		return vec.DIR_RIGHT
	case K_KP_7:
		return vec.DIR_UPLEFT
	case K_KP_9:
		return vec.DIR_UPRIGHT
	case K_KP_1:
		return vec.DIR_DOWNLEFT
	case K_KP_3:
		return vec.DIR_DOWNRIGHT
	default:
		return vec.DIR_NONE
	}
}

// If the keyboard event represents text in some way (letter, number, anything traditionally
// representable by a string) this returns the string form. If not, it return an empty string.
func (kb KeyboardEvent) Text() (s string) {
	if s, ok := keytextmap[kb.Key]; ok {
		if kb.Mods == KEYMOD_SHIFT {
			s = strings.ToUpper(s)
		}
		return s
	} else {
		return ""
	}
}

func (kb KeyboardEvent) String() (s string) {
	s = "KEY "
	if kb.PressType == KEY_PRESSED {
		s += "PRESSED: "
	} else {
		s += "RELEASED: "
	}
	if txt := kb.Text(); txt != "" {
		s += txt
	} else {
		s += "some non text char (" + strconv.Itoa(int(kb.Key)) + ")"
	}

	if kb.Mods != 0 {
		s += " +["
		if kb.Mods&KEYMOD_CTRL != 0 {
			s += "CTRL"
		}
		if kb.Mods&KEYMOD_ALT != 0 {
			s += "ALT"
		}
		if kb.Mods&KEYMOD_SHIFT != 0 {
			s += "SHIFT"
		}

		s += "]"
	}

	return
}

// keycodes. these names are ripped right from github.com/veandco/go-sdl2/sdl/keycode.go
// for now, but we can add/change things here freely as long as we update the corresponding
// mapping from platform-specific keycodes to these in the platform folder.
type Keycode uint8

const (
	K_UNKNOWN      Keycode = iota // "" (no name, empty string)
	K_RETURN                      // "Return" (the Enter key (main keyboard))
	K_ESCAPE                      // "Escape" (the Esc key)
	K_BACKSPACE                   // "Backspace"
	K_TAB                         // "Tab" (the Tab key)
	K_SPACE                       // "Space" (the Space Bar key(s))
	K_EXCLAIM                     // "!"
	K_QUOTEDBL                    // """
	K_HASH                        // "#"
	K_PERCENT                     // "%"
	K_DOLLAR                      // "$"
	K_AMPERSAND                   // "&"
	K_QUOTE                       // "'"
	K_LEFTPAREN                   // "("
	K_RIGHTPAREN                  // ")"
	K_ASTERISK                    // "*"
	K_PLUS                        // "+"
	K_COMMA                       // ","
	K_MINUS                       // "-"
	K_PERIOD                      // "."
	K_SLASH                       // "/"
	K_0                           // "0"
	K_1                           // "1"
	K_2                           // "2"
	K_3                           // "3"
	K_4                           // "4"
	K_5                           // "5"
	K_6                           // "6"
	K_7                           // "7"
	K_8                           // "8"
	K_9                           // "9"
	K_COLON                       // ":"
	K_SEMICOLON                   // ";"
	K_LESS                        // "<"
	K_EQUALS                      // "="
	K_GREATER                     // ">"
	K_QUESTION                    // "?"
	K_AT                          // "@"
	K_LEFTBRACKET                 // "["
	K_BACKSLASH                   // "\"
	K_RIGHTBRACKET                // "]"
	K_CARET                       // "^"
	K_UNDERSCORE                  // "_"
	K_BACKQUOTE                   // "`"
	K_a                           // "A"
	K_b                           // "B"
	K_c                           // "C"
	K_d                           // "D"
	K_e                           // "E"
	K_f                           // "F"
	K_g                           // "G"
	K_h                           // "H"
	K_i                           // "I"
	K_j                           // "J"
	K_k                           // "K"
	K_l                           // "L"
	K_m                           // "M"
	K_n                           // "N"
	K_o                           // "O"
	K_p                           // "P"
	K_q                           // "Q"
	K_r                           // "R"
	K_s                           // "S"
	K_t                           // "T"
	K_u                           // "U"
	K_v                           // "V"
	K_w                           // "W"
	K_x                           // "X"
	K_y                           // "Y"
	K_z                           // "Z"
	K_CAPSLOCK                    // "CapsLock"
	K_F1                          // "F1"
	K_F2                          // "F2"
	K_F3                          // "F3"
	K_F4                          // "F4"
	K_F5                          // "F5"
	K_F6                          // "F6"
	K_F7                          // "F7"
	K_F8                          // "F8"
	K_F9                          // "F9"
	K_F10                         // "F10"
	K_F11                         // "F11"
	K_F12                         // "F12"
	K_PRINTSCREEN                 // "PrintScreen"
	K_SCROLLLOCK                  // "ScrollLock"
	K_PAUSE                       // "Pause" (the Pause / Break key)
	K_INSERT                      // "Insert" (insert on PC, help on some Mac keyboards (but does send code 73, not 117))
	K_HOME                        // "Home"
	K_PAGEUP                      // "PageUp"
	K_DELETE                      // "Delete"
	K_END                         // "End"
	K_PAGEDOWN                    // "PageDown"
	K_RIGHT                       // "Right" (the Right arrow key (navigation keypad))
	K_LEFT                        // "Left" (the Left arrow key (navigation keypad))
	K_DOWN                        // "Down" (the Down arrow key (navigation keypad))
	K_UP                          // "Up" (the Up arrow key (navigation keypad))
	K_NUMLOCKCLEAR                // "Numlock" (the Num Lock key (PC) / the Clear key (Mac))
	K_KP_DIVIDE                   // "Keypad /" (the / key (numeric keypad))
	K_KP_MULTIPLY                 // "Keypad *" (the * key (numeric keypad))
	K_KP_MINUS                    // "Keypad -" (the - key (numeric keypad))
	K_KP_PLUS                     // "Keypad +" (the + key (numeric keypad))
	K_KP_ENTER                    // "Keypad Enter" (the Enter key (numeric keypad))
	K_KP_1                        // "Keypad 1" (the 1 key (numeric keypad))
	K_KP_2                        // "Keypad 2" (the 2 key (numeric keypad))
	K_KP_3                        // "Keypad 3" (the 3 key (numeric keypad))
	K_KP_4                        // "Keypad 4" (the 4 key (numeric keypad))
	K_KP_5                        // "Keypad 5" (the 5 key (numeric keypad))
	K_KP_6                        // "Keypad 6" (the 6 key (numeric keypad))
	K_KP_7                        // "Keypad 7" (the 7 key (numeric keypad))
	K_KP_8                        // "Keypad 8" (the 8 key (numeric keypad))
	K_KP_9                        // "Keypad 9" (the 9 key (numeric keypad))
	K_KP_0                        // "Keypad 0" (the 0 key (numeric keypad))
	K_KP_PERIOD                   // "Keypad ." (the . key (numeric keypad))
)

var keytextmap = map[Keycode]string{
	K_SPACE:        " ",
	K_EXCLAIM:      "!",
	K_QUOTEDBL:     "\"",
	K_HASH:         "#",
	K_PERCENT:      "%",
	K_DOLLAR:       "$",
	K_AMPERSAND:    "&",
	K_QUOTE:        "'",
	K_LEFTPAREN:    "(",
	K_RIGHTPAREN:   ")",
	K_ASTERISK:     "*",
	K_PLUS:         "+",
	K_COMMA:        ",",
	K_MINUS:        "-",
	K_PERIOD:       ".",
	K_SLASH:        "/",
	K_0:            "0",
	K_1:            "1",
	K_2:            "2",
	K_3:            "3",
	K_4:            "4",
	K_5:            "5",
	K_6:            "6",
	K_7:            "7",
	K_8:            "8",
	K_9:            "9",
	K_COLON:        ":",
	K_SEMICOLON:    ";",
	K_LESS:         "<",
	K_EQUALS:       "=",
	K_GREATER:      ">",
	K_QUESTION:     "?",
	K_AT:           "@",
	K_LEFTBRACKET:  "[",
	K_BACKSLASH:    "\\",
	K_RIGHTBRACKET: "]",
	K_CARET:        "^",
	K_UNDERSCORE:   "_",
	K_BACKQUOTE:    "`",
	K_a:            "a",
	K_b:            "b",
	K_c:            "c",
	K_d:            "d",
	K_e:            "e",
	K_f:            "f",
	K_g:            "g",
	K_h:            "h",
	K_i:            "i",
	K_j:            "j",
	K_k:            "k",
	K_l:            "l",
	K_m:            "m",
	K_n:            "n",
	K_o:            "o",
	K_p:            "p",
	K_q:            "q",
	K_r:            "r",
	K_s:            "s",
	K_t:            "t",
	K_u:            "u",
	K_v:            "v",
	K_w:            "w",
	K_x:            "x",
	K_y:            "y",
	K_z:            "z",
	K_KP_1:         "1",
	K_KP_2:         "2",
	K_KP_3:         "3",
	K_KP_4:         "4",
	K_KP_5:         "5",
	K_KP_6:         "6",
	K_KP_7:         "7",
	K_KP_8:         "8",
	K_KP_9:         "9",
	K_KP_0:         "0",
}
