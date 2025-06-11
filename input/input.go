package input

import (
	"github.com/bennicholls/tyumi/event"
)

var (
	EV_KEYBOARD    = event.Register("Key Event")
	EV_MOUSEMOVE   = event.Register("Mouse Move Event")
	EV_MOUSEBUTTON = event.Register("Mouse Button Event")
)

// Set this to true to have Tyumi emit key-repeat events when keys are held down
var AllowKeyRepeats bool

// Set this to true to keep Tyumi from emitting key events when a key is released.
var SuppressKeyUpEvents bool

// Enables mouse input, allowing mouse events to be fired.
var EnableMouse bool
