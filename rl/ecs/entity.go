package ecs

import (
	"github.com/bennicholls/tyumi/log"
)

var entities []Entity
var freeIDs chan uint32
var generations []uint8

const maxFreeIDs int = 32

const INVALID_ID Entity = 0xffffffff
const indexMask uint32 = 0x00ffffff
const generationMask uint32 = 0xff000000

func init() {
	entities = make([]Entity, 0)
	freeIDs = make(chan uint32, 256)
	generations = make([]uint8, 0)
}

type Entity uint32

func (e Entity) index() uint32 {
	return uint32(e) & indexMask
}

func (e Entity) Valid() bool {
	return e != INVALID_ID
}

// Reports whether the entity is still active (has not been removed/deleted)
func (e Entity) Alive() bool {
	return e.Valid() && e == entities[e.index()]
}

// Creates an Entity. Entities are just a number
func CreateEntity() (e Entity) {
	if len(freeIDs) < maxFreeIDs { //append to entities list, return ID with generation 0
		e = Entity(len(entities))
		entities = append(entities, e)
		generations = append(generations, 0)
	} else { // take first free ID, retrieve generation for that slot, increment, compile ID, store new ID and gen, return
		idx := <-freeIDs
		gen := uint32(generations[idx]) + 1
		if gen == 255 {
			log.Warning("GENERATION LIMIT REACHED!!! (If you see this, it's not a *big* deal but there's a very small chance a bug could occur going forward.)")
		}
		generations[idx] = uint8(gen)
		e = Entity(idx | (gen << 24))
		entities[idx] = e
	}

	return
}

func RemoveEntity(e Entity) {
	if !e.Alive() {
		log.Debug("Double removing entity??")
		return
	}

	entities[e.index()] = INVALID_ID
	addFreeID(e.index())

	for _, cache := range componentCaches {
		cache.removeComponent(e)
	}
}

func addFreeID(idx uint32) {
	if len(freeIDs) == cap(freeIDs) {
		newChannel := make(chan uint32, int(float32(cap(freeIDs))*1.5))
		for range len(freeIDs) {
			newChannel <- <-freeIDs
		}
		freeIDs = newChannel
		log.Debug("FreeID channel resized!!")
	}

	freeIDs <- idx
}
