package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
)

func init() {
	ecs.Register[TerrainComponent]()
	ecs.Register[EntityContainerComponent]()
	ecs.Register[PositionComponent]()
	ecs.Register[EntityComponent]()
	ecs.Register[PlayerComponent]()
	ecs.Register[MemoryComponent]()
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

type PlayerComponent struct {
	ecs.Component
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

type MemoryComponent struct {
	ecs.Component

	memory map[vec.Coord]gfx.Visuals
}

func (mc *MemoryComponent) Init() {
	mc.memory = make(map[vec.Coord]gfx.Visuals)
}

func (mc MemoryComponent) GetVisuals(pos vec.Coord) (vis gfx.Visuals, ok bool) {
	vis, ok = mc.memory[pos]
	return
}

func (mc *MemoryComponent) AddVisuals(pos vec.Coord, vis gfx.Visuals) {
	if vis.Mode == gfx.DRAW_NONE {
		delete(mc.memory, pos)
		return
	}

	mc.memory[pos] = vis
}
