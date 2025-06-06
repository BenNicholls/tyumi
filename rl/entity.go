package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
)

type EntityType uint32

func (et EntityType) Data() EntityData {
	return entityDataCache.GetData(et)
}

type EntityData struct {
	Name    string
	Desc    string
	Glyph   gfx.Glyph
	Colours col.Pair
}

func (ed EntityData) GetVisuals() gfx.Visuals {
	return gfx.NewGlyphVisuals(ed.Glyph, ed.Colours)
}

var entityDataCache dataCache[EntityData, EntityType]

func RegisterEntityType(entity_data EntityData) EntityType {
	return entityDataCache.RegisterDataType(entity_data)
}

// Entity represents a tilemap object. Each tile can hold one entity (at most). Examples of entities would be things
// like the player, enemies, furniture and other decorations, etc. DO NOT confuse these with ECS Entities, which can
// literally be anything. I'm sorry for the name clash but I genuinely can't think of a better name for these right now.
type Entity ecs.Entity

var INVALID_ENTITY = Entity(ecs.INVALID_ID)

func CreateEntity(entity_type EntityType) (entity Entity) {
	entity = Entity(ecs.CreateEntity())

	ecs.AddComponent(entity, EntityComponent{EntityType: entity_type})
	ecs.AddComponent(entity, PositionComponent{Coord: NOT_IN_TILEMAP})

	return
}

func (e Entity) GetVisuals() gfx.Visuals {
	return ecs.GetComponent[EntityComponent](e).EntityType.Data().GetVisuals()
}

func (e Entity) GetName() string {
	return ecs.GetComponent[EntityComponent](e).EntityType.Data().Name
}

func (e Entity) Position() vec.Coord {
	return ecs.GetComponent[PositionComponent](e).Coord
}

func (e Entity) MoveTo(pos vec.Coord) {
	position := ecs.GetComponent[PositionComponent](e)
	if position.Static && (pos != NOT_IN_TILEMAP) {
		return
	}

	position.Coord = pos
	if fov := ecs.GetComponent[FOVComponent](e); fov != nil {
		fov.Dirty = true
	}
}

func (e Entity) IsInTilemap() bool {
	return ecs.GetComponent[PositionComponent](e).Coord != NOT_IN_TILEMAP
}
