package input

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/vec"
)

type MouseMoveEvent struct {
	event.EventPrototype

	Position vec.Coord //current position of the mouse in the console
	Delta    vec.Coord //movement of the mouse since last frame
}

func (mme MouseMoveEvent) String() string {
	return "Mouse Move Event: pos " + mme.Position.String() + ", delta " + mme.Delta.String()
}

func FireMouseMoveEvent(pos, delta vec.Coord) {
	if !EnableMouse {
		return
	}

	event.Fire(EV_MOUSEMOVE, &MouseMoveEvent{Position: pos, Delta: delta})
}
