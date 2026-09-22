package web

import "github.com/bennicholls/tyumi/input"

// codeToKey maps a DOM KeyboardEvent.code to a Tyumi key. These are physical keys
// whose identity should not follow the shifted character (arrows, function keys,
// the numpad).
var codeToKey = map[string]input.Keycode{
	"Enter":          input.K_RETURN,
	"NumpadEnter":    input.K_KP_ENTER,
	"Escape":         input.K_ESCAPE,
	"Backspace":      input.K_BACKSPACE,
	"Tab":            input.K_TAB,
	"Space":          input.K_SPACE,
	"CapsLock":       input.K_CAPSLOCK,
	"F1":             input.K_F1,
	"F2":             input.K_F2,
	"F3":             input.K_F3,
	"F4":             input.K_F4,
	"F5":             input.K_F5,
	"F6":             input.K_F6,
	"F7":             input.K_F7,
	"F8":             input.K_F8,
	"F9":             input.K_F9,
	"F10":            input.K_F10,
	"F11":            input.K_F11,
	"F12":            input.K_F12,
	"PrintScreen":    input.K_PRINTSCREEN,
	"ScrollLock":     input.K_SCROLLLOCK,
	"Pause":          input.K_PAUSE,
	"Insert":         input.K_INSERT,
	"Home":           input.K_HOME,
	"PageUp":         input.K_PAGEUP,
	"Delete":         input.K_DELETE,
	"End":            input.K_END,
	"PageDown":       input.K_PAGEDOWN,
	"ArrowRight":     input.K_RIGHT,
	"ArrowLeft":      input.K_LEFT,
	"ArrowDown":      input.K_DOWN,
	"ArrowUp":        input.K_UP,
	"NumLock":        input.K_NUMLOCKCLEAR,
	"NumpadDivide":   input.K_KP_DIVIDE,
	"NumpadMultiply": input.K_KP_MULTIPLY,
	"NumpadSubtract": input.K_KP_MINUS,
	"NumpadAdd":      input.K_KP_PLUS,
	"Numpad1":        input.K_KP_1,
	"Numpad2":        input.K_KP_2,
	"Numpad3":        input.K_KP_3,
	"Numpad4":        input.K_KP_4,
	"Numpad5":        input.K_KP_5,
	"Numpad6":        input.K_KP_6,
	"Numpad7":        input.K_KP_7,
	"Numpad8":        input.K_KP_8,
	"Numpad9":        input.K_KP_9,
	"Numpad0":        input.K_KP_0,
	"NumpadDecimal":  input.K_KP_PERIOD,
}

// charToKey maps the character a key produced. Shifted punctuation stays on its
// own Tyumi key, matching SDL keysyms. Letters are handled separately so "A" and
// "a" both land on K_a, with shift carried as a modifier.
var charToKey = map[string]input.Keycode{
	" ":  input.K_SPACE,
	"!":  input.K_EXCLAIM,
	"\"": input.K_QUOTEDBL,
	"#":  input.K_HASH,
	"$":  input.K_DOLLAR,
	"%":  input.K_PERCENT,
	"&":  input.K_AMPERSAND,
	"'":  input.K_QUOTE,
	"(":  input.K_LEFTPAREN,
	")":  input.K_RIGHTPAREN,
	"*":  input.K_ASTERISK,
	"+":  input.K_PLUS,
	",":  input.K_COMMA,
	"-":  input.K_MINUS,
	".":  input.K_PERIOD,
	"/":  input.K_SLASH,
	"0":  input.K_0,
	"1":  input.K_1,
	"2":  input.K_2,
	"3":  input.K_3,
	"4":  input.K_4,
	"5":  input.K_5,
	"6":  input.K_6,
	"7":  input.K_7,
	"8":  input.K_8,
	"9":  input.K_9,
	":":  input.K_COLON,
	";":  input.K_SEMICOLON,
	"<":  input.K_LESS,
	"=":  input.K_EQUALS,
	">":  input.K_GREATER,
	"?":  input.K_QUESTION,
	"@":  input.K_AT,
	"[":  input.K_LEFTBRACKET,
	"\\": input.K_BACKSLASH,
	"]":  input.K_RIGHTBRACKET,
	"^":  input.K_CARET,
	"_":  input.K_UNDERSCORE,
	"`":  input.K_BACKQUOTE,
}

// KeycodeFromJS converts a DOM keyboard event into a Tyumi keycode.
// code is KeyboardEvent.code and key is KeyboardEvent.key.
// The boolean is false for keys Tyumi does not track, such as Shift itself.
func KeycodeFromJS(code, key string) (input.Keycode, bool) {
	if k, ok := codeToKey[code]; ok {
		return k, true
	}
	if k, ok := charToKey[key]; ok {
		return k, true
	}
	if len(key) == 1 {
		r := key[0]
		switch {
		case r >= 'a' && r <= 'z':
			return input.K_a + input.Keycode(r-'a'), true
		case r >= 'A' && r <= 'Z':
			return input.K_a + input.Keycode(r-'A'), true
		}
	}
	return input.K_UNKNOWN, false
}
