package ecs

import (
	"github.com/bennicholls/tyumi/log"
)

var entities []Entity
var freeIndices chan uint32
var generations []uint8

const maxFreeIDs int = 32

const INVALID_ID Entity = 0
const indexMask uint32 = 0x00ffffff
const generationMask uint32 = 0xff000000

func init() {
	entities = make([]Entity, 0)
	freeIndices = make(chan uint32, 256)
	generations = make([]uint8, 0)
}

type Entity uint32

// Valid reports whether an entity ID is valid and properly formed.
func Valid[EntityType ~uint32](entity EntityType) (valid bool) {
	valid = entity != EntityType(INVALID_ID)
	if Debug {
		// if in debug mode, we also check to see if the id's index doesn't overflow the entity list. this should never
		// be possible with actual ids from the ecs, but we do the check just in case some user acidentally passes
		// in some other kind of uint32 from somewhere else
		 valid = valid && (index(entity) < uint32(len(entities)))
	}

	return
}

// Alive reports whether an entity is valid and has not been removed from the ECS.
func Alive[EntityType ~uint32](entity EntityType) bool {
	return Valid(entity) && Entity(entity) == entities[index(entity)]
}

func index[EntityType ~uint32](entity EntityType) uint32 {
	return (uint32(entity) & indexMask) - 1
}

// Creates an Entity. Entities are just a number
func CreateEntity() (entity Entity) {
	if len(freeIndices) < maxFreeIDs { //append to entities list, return ID with generation 0
		entity = Entity(len(entities) + 1) // REMEMBER: this is +1 because zero is the INVALID_ID
		entities = append(entities, entity)
		generations = append(generations, 0)
	} else { // take first free ID, retrieve generation for that slot, increment, compile ID, store new ID and gen, return
		idx := <-freeIndices
		gen := uint32(generations[idx]) + 1
		if Debug && gen == 255 {
			// if in debug mode, this check warns us when our generation number rolls over. this isn't necessarily an
			// issue -- it can only cause a bug in a very case where the user is holding a removed id for a slot that
			// has been reused exactly 255 times since being removed -- so it's not worth erroring out or anything, but
			// while making the game this warning could be useful just as a smoke test if something starts going wrong.
			// if this does end up causing bugs, we need to make a way to increase the number of generation bits in the
			// id.
			log.Warning("ECS: GENERATION LIMIT REACHED!!! (If you see this, it's not a *big* deal but there's a very small chance a bug could occur going forward.)")
		}
		generations[idx] = uint8(gen)
		entity = Entity((idx + 1) | (gen << 24))
		entities[idx] = entity
	}

	return
}

// CopyEntity creates a new entity that is a copy of the provided entity. All of entity e's components are copied and
// assigned to the new entity.
func CopyEntity[EntityType ~uint32](entity EntityType) (copy Entity) {
	if !Alive(entity) {
		log.Debug("ECS: Cannot copy dead entity!")
		return
	}

	copy = CreateEntity()

	for _, cache := range componentCaches {
		if cache.hasComponent(Entity(entity)) {
			cache.copyComponent(Entity(entity), copy)
		}
	}

	return
}

// RemoveEntity removes an entity from the ECS. All of its components will be removed and it will be set as dead.
func RemoveEntity[EntityType ~uint32](entity EntityType) {
	if !Alive(entity) {
		log.Debug("ECS: Removing dead/invalid entity??")
		return
	}

	entities[index(entity)] = INVALID_ID
	addFreeID(index(entity))

	for _, cache := range componentCaches {
		cache.removeComponent(Entity(entity))
	}
}

func addFreeID(idx uint32) {
	if len(freeIndices) == cap(freeIndices) {
		newChannel := make(chan uint32, int(float32(cap(freeIndices))*1.5))
		for range len(freeIndices) {
			newChannel <- <-freeIndices
		}
		freeIndices = newChannel
		log.Debug("ECS: FreeID channel resized!!")
	}

	freeIndices <- idx
}
