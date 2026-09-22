package web

import (
	"testing"

	"github.com/bennicholls/tyumi/input"
)

func TestKeycodeFromJS(t *testing.T) {
	cases := []struct {
		code string
		key  string
		want input.Keycode
		ok   bool
	}{
		{"KeyQ", "q", input.K_q, true},
		{"KeyQ", "Q", input.K_q, true},
		{"Digit1", "1", input.K_1, true},
		{"Digit1", "!", input.K_EXCLAIM, true},
		{"ArrowUp", "ArrowUp", input.K_UP, true},
		{"Numpad4", "4", input.K_KP_4, true},
		{"F5", "F5", input.K_F5, true},
		{"Enter", "Enter", input.K_RETURN, true},
		{"ShiftLeft", "Shift", input.K_UNKNOWN, false},
	}

	for _, tc := range cases {
		got, ok := KeycodeFromJS(tc.code, tc.key)
		if ok != tc.ok || got != tc.want {
			t.Errorf("KeycodeFromJS(%q, %q) = (%v, %v), want (%v, %v)", tc.code, tc.key, got, ok, tc.want, tc.ok)
		}
	}
}
