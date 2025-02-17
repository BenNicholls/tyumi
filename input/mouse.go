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

func newMouseMoveEvent(pos, delta vec.Coord) (mme *MouseMoveEvent) {
	mme = new(MouseMoveEvent)
	mme.EventPrototype = *event.New(EV_MOUSEMOVE)
	mme.Position = pos
	mme.Delta = delta
	return
}

func (mme MouseMoveEvent) String() string {
	return "Mouse Move Event: pos " + mme.Position.String() + ", delta " + mme.Delta.String()
}

func FireMouseMoveEvent(pos, delta vec.Coord) {
	if !EnableMouse {
		return
	}
	
	event.Fire(newMouseMoveEvent(pos, delta))
}
