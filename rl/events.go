package rl

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/vec"
)

var EV_ENTITYMOVED = event.Register("Entity moved.")
var EV_TILECHANGEDVISIBILITY = event.Register("A Tile Changed visibility state (opaque or transparent)")

type EntityMovedEvent struct {
	event.EventPrototype

	Entity   Entity
	From, To vec.Coord
}

func fireEntityMovedEvent(entity Entity, from, to vec.Coord) {
	event.Fire(EV_ENTITYMOVED, &EntityMovedEvent{Entity: entity, From: from, To: to})
}

type TileChangedVisibilityEvent struct {
	event.EventPrototype

	Pos vec.Coord
}
