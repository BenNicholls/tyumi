package input

import (
	"github.com/bennicholls/tyumi/event"
)

var EV_KEYBOARD = event.Register()
var EV_QUIT = event.Register()

// definition of whatever system is grabbing events from the system
type Processor func()

