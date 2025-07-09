package rl

import (
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

func init() {
	ecs.Register[FOVComponent]()
}

// FOVComponent is for anything that can see.
type FOVComponent struct {
	ecs.Component

	Blind         bool  // if true, no fov is possible
	Omniscient    bool  // if true, all tiles are reported as within the FOV. good for testing. overrides blind.
	Dirty         bool  // if true, FOV needs to recompute
	SightRange    uint8 // range of FOV in tiles
	TrackEntities bool

	field    util.Set[vec.Coord]
	entities util.Set[Entity]
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

func (fov *FOVComponent) InFOV(pos vec.Coord) bool {
	if fov.Omniscient {
		return true
	} else if fov.Blind {
		return false
	}

	return fov.field.Contains(pos)
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
		for fov := range ecs.EachComponent[FOVComponent]() {
			if Entity(fov.GetEntity()) == moveEvent.Entity {
				fov.Dirty = true
				continue
			}

			if !fov.TrackEntities {
				continue
			}

			if fov.InFOV(moveEvent.From) && !fov.InFOV(moveEvent.To) { // entity moved away
				fov.entities.Remove(moveEvent.Entity)
				event.Fire(EV_LOSTSIGHT, &EntitySightEvent{
					Viewer:        Entity(fov.GetEntity()),
					TrackedEntity: moveEvent.Entity},
				)
			} else if fov.InFOV(moveEvent.To) && !fov.InFOV(moveEvent.From) { //entity moved into the fov
				fov.entities.Add(moveEvent.Entity)
				event.Fire(EV_GAINEDSIGHT, &EntitySightEvent{
					Viewer:        Entity(fov.GetEntity()),
					TrackedEntity: moveEvent.Entity},
				)
			}
		}

		return true
	case EV_TILECHANGEDVISIBILITY:
		visEvent := e.(*TileChangedVisibilityEvent)
		fs.changedVisbilityTiles.Add(visEvent.Pos)
		return true
	}

	return
}

func (fs *FOVSystem) Update(delta time.Duration) {
	fs.System.Update(delta)

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
			fs.computeFOV(fov)
			if Entity(fov.GetEntity()).IsPlayer() {
				fs.tileMap.SetAllDirty()
			}

			if fov.TrackEntities {
				var newEntities util.Set[Entity]

				for entity := range ecs.EachEntityWith[EntityComponent]() {
					if entity == fov.GetEntity() { // don't track self
						continue
					}

					if pos := ecs.Get[PositionComponent](entity).Coord; fov.InFOV(pos) {
						newEntities.Add(Entity(entity))
					}
				}

				if !fov.entities.Equals(newEntities) {
					lostSight := fov.entities.Difference(newEntities)
					for entity := range lostSight.EachElement() {
						event.Fire(EV_LOSTSIGHT, &EntitySightEvent{
							Viewer:        Entity(fov.GetEntity()),
							TrackedEntity: entity})
					}

					gainedSight := newEntities.Difference(fov.entities)
					for entity := range gainedSight.EachElement() {
						event.Fire(EV_GAINEDSIGHT, &EntitySightEvent{
							Viewer:        Entity(fov.GetEntity()),
							TrackedEntity: entity})
					}

					fov.entities = newEntities
				}
			}
		}
	}

	fs.changedVisbilityTiles.RemoveAll()
}

func (fs *FOVSystem) computeFOV(fov *FOVComponent) {
	fov.Dirty = false

	if fov.Blind || fov.Omniscient {
		return
	}

	pos := ecs.Get[PositionComponent](fov.GetEntity()).Coord
	if pos == NOT_IN_TILEMAP {
		return
	}

	fs.tileMap.ShadowCast(pos, int(fov.SightRange), GetSpacesSetCast(&fov.field))
}
