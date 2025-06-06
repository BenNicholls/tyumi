package rl

import (
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
)

func init() {
	ecs.RegisterComponent[PositionComponent]()
	ecs.RegisterComponent[TerrainComponent]()
	ecs.RegisterComponent[EntityComponent]()
	ecs.RegisterComponent[EntityContainerComponent]()
}

type PositionComponent struct {
	ecs.Component
	vec.Coord

	Static bool
}

type TerrainComponent struct {
	ecs.Component
	TileType
}

type EntityComponent struct {
	ecs.Component
	EntityType
}

type EntityContainerComponent struct {
	ecs.Component

	Entity
}

func (ecc EntityContainerComponent) Empty() bool {
	return ecc.Entity == Entity(ecs.INVALID_ID)
}

func (ecc *EntityContainerComponent) Add(entity Entity) {
	if ecc.Entity != Entity(ecs.INVALID_ID) {
		log.Debug("Overwriting entity!!")
	}
	ecc.Entity = entity
}

func (ecc *EntityContainerComponent) Remove() {
	ecc.Entity = Entity(ecs.INVALID_ID)
}
