package input

import (
	"github.com/bennicholls/tyumi/event"
)

var EV_KEYBOARD = event.Register("Key Event")
var EV_MOUSEMOVE = event.Register("Mouse Move Event")
var EV_MOUSEBUTTON = event.Register("Mouse Button Event")
