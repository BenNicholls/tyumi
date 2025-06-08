package rl

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

func init() {
	ecs.RegisterComponent[FOVComponent]()
}

// FOVComponent is for anything that can see.
type FOVComponent struct {
	ecs.Component

	Blind      bool  // if true, no fov is possible
	Omniscient bool  // if true, all tiles are reported as within the FOV. good for testing. overrides blind.
	Dirty      bool  // if true, FOV needs to recompute
	SightRange uint8 // range of FOV in tiles
	field      util.Set[vec.Coord]
}

func (fov *FOVComponent) Init() {
	fov.Dirty = true
}

func (fov *FOVComponent) SetBlind(blind bool) {
	if fov.Blind == blind {
		return
	}

	fov.Blind = blind
	if blind {
		fov.field.RemoveAll()
		fov.Dirty = false
	} else {
		fov.Dirty = true
	}
}

func (fov *FOVComponent) SetOmniscience(omniscient bool) {
	if fov.Omniscient == omniscient {
		return
	}

	fov.Omniscient = omniscient
	if fov.Omniscient {
		fov.field.RemoveAll()
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
	if fov.Dirty || fov.Blind || fov.Omniscient {
		return
	}

	if fov.field.Contains(pos) {
		fov.Dirty = true
	}
}

func (fov *FOVComponent) InFOV(pos vec.Coord) bool {
	if fov.Omniscient {
		return true
	} else if fov.Blind {
		return false
	}

	return fov.field.Contains(pos)
}

// Runs the shadowcaster for the tilemap and updates this entity's FOV set.
func (fov *FOVComponent) UpdateFOV(tileMap *TileMap) {
	fov.Dirty = false

	if fov.Blind || fov.Omniscient {
		return
	}

	pos := ecs.GetComponent[PositionComponent](fov.GetEntity()).Coord
	if pos == NOT_IN_TILEMAP {
		return
	}

	tileMap.ShadowCast(pos, int(fov.SightRange), GetSpacesSetCast(&fov.field))
}

type FOVSystem struct {
	System

	tileMap               *TileMap
	changedVisbilityTiles util.Set[vec.Coord]
}

func (fs *FOVSystem) Init(tm *TileMap) {
	fs.tileMap = tm
	fs.Listen(EV_ENTITYMOVED, EV_TILECHANGEDVISIBILITY)
	fs.SetEventHandler(fs.handleEvents)
}

func (fs *FOVSystem) handleEvents(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_ENTITYMOVED:
		moveEvent := e.(*EntityMovedEvent)
		if fov := ecs.GetComponent[FOVComponent](moveEvent.Entity); fov != nil {
			fov.Dirty = true
		}
		return true
	case EV_TILECHANGEDVISIBILITY:
		visEvent := e.(*TileChangedVisibilityEvent)
		fs.changedVisbilityTiles.Add(visEvent.Pos)
		return true
	}

	return
}

func (fs *FOVSystem) Update() {
	fs.System.Update()

	// FOV updates
	for fov := range ecs.EachComponent[FOVComponent]() {
		if fov.Blind || fov.Omniscient {
			continue
		}

		if !fov.Dirty {
			for pos := range fs.changedVisbilityTiles.EachElement() {
				if fov.InFOV(pos) {
					fov.Dirty = true
					break
				}
			}
		}

		if fov.Dirty {
			fov.UpdateFOV(fs.tileMap)
		}
	}

	fs.changedVisbilityTiles.RemoveAll()
}
