package rl

import (
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
)

func init() {
	ecs.RegisterComponent[PositionComponent]()
	ecs.RegisterComponent[TerrainComponent]()
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

type EntityContainerComponent struct {
	ecs.Component

	TileMapEntity
}

func (ecc EntityContainerComponent) Empty() bool {
	return ecc.TileMapEntity == nil
}
