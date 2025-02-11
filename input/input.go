package input

import (
	"github.com/bennicholls/tyumi/event"
)

var EV_KEYBOARD = event.Register("Key Event")
var EV_MOUSEMOVE = event.Register("Mouse Move Event")
var EV_MOUSEBUTTON = event.Register("Mouse Button Event")

var suppress_key_repeats bool

// DisableKeyRepeats will suppress all key events from keys being held down. 
func SuppressKeyRepeats() {
	suppress_key_repeats = true
}

// EnableKeyRepeats will re-enable key events for held down keys. Note that key repeat events are sent by
// default so if you need them you do not have to call this.
func EnableKeyRepeats() {
	suppress_key_repeats = false
}