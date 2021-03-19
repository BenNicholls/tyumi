package engine

import (
	"github.com/bennicholls/tyumi/event"
)

var EV_KEYBOARD = event.Register()

type KeyboardEvent struct {
	event.EventPrototype

	Key int
}

func NewKeyboardEvent(key int) (kbe KeyboardEvent) {
	kbe.EventPrototype = event.New(EV_KEYBOARD)
	kbe.Key = key
	return
}