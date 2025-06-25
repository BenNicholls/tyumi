package rl

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
)

type EntityType uint32

func (et EntityType) Data() EntityData {
	return entityDataCache.GetData(et)
}

type EntityData struct {
	Name           string // Generic name of the entity.
	Desc           string // Generic description for the entity
	Visuals        gfx.Visuals
	SightRange     uint8
	TracksEntities bool
	HasMemory      bool
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

	ecs.Add(entity, EntityComponent{EntityType: entity_type})
	ecs.Add(entity, PositionComponent{Coord: NOT_IN_TILEMAP})

	if sight := entity_type.Data().SightRange; sight > 0 {
		ecs.Add(entity, FOVComponent{SightRange: sight,
			TrackEntities: entity_type.Data().TracksEntities})
	}

	if entity_type.Data().HasMemory {
		ecs.Add[MemoryComponent](entity)
	}

	return
}

// Destroy removes the entity from the ECS. Before doing so, it emits EV_ENTITYBEINGDESTROYED, an event fired in
// immediate mode. Systems that need to do some cleanup when an entity is destroyed can listen for this event and
// respond accordingly. Before being removed, all components on the entity will have their Cleanup() function run, if
// present.
func (e Entity) Destroy() {
	if !ecs.Alive(e) {
		log.Debug("Trying to destroy an entity that is already dead!!")
	}

	log.Debug("Now we do the immediate fire thing.")
	event.FireImmediate(EV_ENTITYBEINGDESTROYED, &EntityEvent{Entity: e})
	log.Debug("Now we destroy.")
	ecs.RemoveEntity(e)
}

func (e Entity) GetVisuals() gfx.Visuals {
	return ecs.Get[EntityComponent](e).EntityType.Data().Visuals
}

func (e Entity) GetName() string {
	return ecs.Get[EntityComponent](e).EntityType.Data().Name
}

func (e Entity) Position() vec.Coord {
	return ecs.Get[PositionComponent](e).Coord
}

func (e Entity) IsPlayer() bool {
	return ecs.Has[PlayerComponent](e)
}

func (e Entity) MoveTo(pos vec.Coord) {
	position := ecs.Get[PositionComponent](e)
	if position.Static && pos != NOT_IN_TILEMAP {
		return
	}

	if position.Coord != NOT_IN_TILEMAP {
		event.Fire(EV_ENTITYMOVED, &EntityMovedEvent{Entity: e, From: position.Coord, To: pos})
	}
	position.Coord = pos
}

func (e Entity) IsInTilemap() bool {
	return ecs.Get[PositionComponent](e).Coord != NOT_IN_TILEMAP
}
