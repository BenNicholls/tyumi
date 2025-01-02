package input

import (
	"github.com/bennicholls/tyumi/event"
)

var EV_KEYBOARD = event.Register("Key Event")
var EV_QUIT = event.Register("Quit Event")

// definition of whatever system is grabbing events from the system
type Processor func()

