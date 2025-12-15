package rl

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/vec"
)

var (
	EV_ENTITYBEINGDESTROYED  = event.Register("Entity being destroyed/removed from the ECS.")
	EV_ENTITYMOVED           = event.Register("Entity moved.")
	EV_ENTITYHEALTHCHANGED   = event.Register("Entity's health changed.")
	EV_ENTITYDIED            = event.Register("Entity has been killed/destroyed.")
	EV_TILECHANGEDVISIBILITY = event.Register("A Tile Changed visibility state (opaque or transparent)")
	EV_LOSTSIGHT             = event.Register("An entity has lost sight of another entity.")
	EV_GAINEDSIGHT           = event.Register("An entity has gained sight of another entity.")
)

type EntityEvent struct {
	event.EventPrototype

	Entity Entity
}

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

	Pos    vec.Coord
	Opaque bool
}

type EntitySightEvent struct {
	event.EventPrototype

	Viewer        Entity
	TrackedEntity Entity
}

type EntityHealthChangedEvent struct {
	event.EventPrototype

	Entity       Entity
	OldHP, NewHP int
}
