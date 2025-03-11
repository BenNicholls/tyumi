package input

import (
	"github.com/bennicholls/tyumi/event"
)

var EV_KEYBOARD = event.Register("Key Event", event.COMPLEX)
var EV_MOUSEMOVE = event.Register("Mouse Move Event", event.COMPLEX)
var EV_MOUSEBUTTON = event.Register("Mouse Button Event", event.COMPLEX)

// Set this to true to have Tyumi emit key-repeat events when keys are held down
var AllowKeyRepeats bool

// Enables mouse input, allowing mouse events to be fired.
var EnableMouse bool