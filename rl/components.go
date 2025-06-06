package rl

import (
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

func init() {
	ecs.RegisterComponent[PositionComponent]()
	ecs.RegisterComponent[TerrainComponent]()
	ecs.RegisterComponent[EntityComponent]()
	ecs.RegisterComponent[EntityContainerComponent]()
	ecs.RegisterComponent[FOVComponent]()
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

// FOVComponent is for anything that can see. 
type FOVComponent struct {
	ecs.Component

	blind      bool  // if true, no fov is possible
	Dirty      bool  // if true, FOV needs to recompute
	SightRange uint8 // range of FOV in tiles
	FOV        util.Set[vec.Coord]
}

func (fov FOVComponent) Blind() bool {
	return fov.blind
}

func (fov *FOVComponent) SetBlind(blind bool) {
	if fov.blind == blind {
		return
	}

	fov.blind = blind
	if blind {
		fov.FOV.RemoveAll()
		fov.Dirty = false
	} else {
		fov.Dirty = true
	}
}

func (fov *FOVComponent) SetSightRange(sight_range uint8) {
	if fov.SightRange == sight_range {
		return
	}

	fov.SightRange = sight_range
	fov.Dirty = true
}

// call this when a change that could affect FOV happens. If the change is within the current FOV, sets the dirty flag
// so FOV is recomputed before being accessed again.
func (fov *FOVComponent) OnEnvironmentChange(pos vec.Coord) {
	if !fov.Dirty && fov.FOV.Contains(pos) {
		fov.Dirty = true
	}
}

// Runs the shadowcaster for the tilemap and updates this entity's FOV set.
func (fov *FOVComponent) UpdateFOV(tileMap *TileMap) {
	if fov.blind {
		return
	}

	pos := ecs.GetComponent[PositionComponent](fov.GetEntity()).Coord
	if pos == NOT_IN_TILEMAP {
		return
	}

	tileMap.ShadowCast(pos, int(fov.SightRange), GetSpacesCast(&fov.FOV))
	fov.Dirty = false
}
