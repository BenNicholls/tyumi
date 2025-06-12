package rl

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

func init() {
	ecs.Register[LightSourceComponent]()
}

type LightSourceComponent struct {
	ecs.Component

	Disabled    bool
	Dirty       bool  // if true, light needs to be reapplied to its area
	AreaDirty   bool  // if true, light needs to recompute which tiles it is affecting
	Power       uint8 // light level applied at the source
	FalloffRate uint8 // amount light level diminishes every 1 tile away from source
	MaxRange    uint8 // maximum range of the light. If 0 (default), the light's distance is computed from Power and FalloffRate
	//Colour      col.Colour // Colour of light

	litTiles map[vec.Coord]uint8 // map of tile position to the amount of light being cast there
}

func (lsc *LightSourceComponent) Init() {
	if !lsc.Disabled {
		lsc.AreaDirty = true
	}
}

func (lsc *LightSourceComponent) Enable() {
	lsc.setDisabled(false)
}

func (lsc *LightSourceComponent) Disable() {
	lsc.setDisabled(true)
}

func (lsc *LightSourceComponent) Toggle() {
	lsc.setDisabled(!lsc.Disabled)
}

func (lsc *LightSourceComponent) setDisabled(disabled bool) {
	if lsc.Disabled == disabled {
		return
	}

	lsc.Disabled = disabled
	if !lsc.Disabled {
		lsc.AreaDirty = true // if enabling, trigger a recompute and apply
	}
}

func (lsc *LightSourceComponent) SetPower(power uint8) {
	if lsc.Power == power {
		return
	}

	lsc.Power = power
	if lsc.Disabled {
		return
	}

	if lsc.MaxRange == 0 {
		lsc.AreaDirty = true
	} else {
		lsc.Dirty = true
	}
}

func (lsc *LightSourceComponent) SetFalloff(falloff uint8) {
	if lsc.FalloffRate == falloff {
		return
	}

	lsc.FalloffRate = falloff
	if lsc.Disabled {
		return
	}

	if lsc.MaxRange == 0 {
		lsc.AreaDirty = true
	} else {
		lsc.Dirty = true
	}
}

func (lsc *LightSourceComponent) SetMaxRange(max_range uint8) {
	if lsc.MaxRange == max_range {
		return
	}

	lsc.MaxRange = max_range
	if !lsc.Disabled {
		lsc.AreaDirty = true
	}
}

type LightSystem struct {
	System

	tileMap               *TileMap
	changedVisbilityTiles util.Set[vec.Coord]
	globalLight           uint8 // amount of light automatically applied to every tile. If 255, disables system.
}

func (ls *LightSystem) Init(tm *TileMap) {
	ls.tileMap = tm
	ls.Listen(EV_ENTITYMOVED, EV_TILECHANGEDVISIBILITY)
	ls.SetEventHandler(ls.handleEvents)
	ls.Enable()
}

// SetGlobalLight sets the amount of light automatically applied to all tiles. If 255, the lighting system is
// disabled because all tiles will be completely lit at all times.
func (ls *LightSystem) SetGlobalLight(light uint8) {
	ls.globalLight = light

	if ls.globalLight == 255 {
		ls.Disable()
	} else {
		ls.Enable()
	}
}

func (ls *LightSystem) handleEvents(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_ENTITYMOVED:
		moveEvent := e.(*EntityMovedEvent)
		if light := ecs.Get[LightSourceComponent](moveEvent.Entity); light != nil {
			light.AreaDirty = true
		}
		return true
	case EV_TILECHANGEDVISIBILITY:
		visEvent := e.(*TileChangedVisibilityEvent)
		ls.changedVisbilityTiles.Add(visEvent.Pos)
		return true
	}

	return
}

func (ls *LightSystem) Update() {
	if !ls.Enabled || ls.tileMap == nil {
		return
	}

	ls.System.Update()

	for light := range ecs.EachComponent[LightSourceComponent]() {
		// check if light should trigger an area update due to nearby tiles changing visibility
		if !light.AreaDirty {
			for pos := range ls.changedVisbilityTiles.EachElement() {
				if _, ok := light.litTiles[pos]; ok {
					light.AreaDirty = true
					break
				}
			}
		}

		if light.AreaDirty {
			ls.computeLightArea(light)
		}
	}

	ls.changedVisbilityTiles.RemoveAll()

	// light application has to go in a separate pass to prevent certain accumulations errors arising from weird
	// situations where tiles are replaced.
	for light := range ecs.EachComponent[LightSourceComponent]() {
		if light.Dirty {
			ls.applyLight(light)
		}
	}
}

func (ls *LightSystem) removeAppliedLight(light *LightSourceComponent) {
	for pos, amount := range light.litTiles {
		if amount == 0 {
			continue
		}

		ls.tileMap.GetTile(pos).RemoveLight(amount)
		light.litTiles[pos] = 0
		ls.tileMap.SetDirty(pos)
	}
}

// computeLitArea runs the shadowcaster to determine the area lit by the light. if the light was already lighting
// some terrain, those terrains are de-lighted first. how delightful!
func (ls *LightSystem) computeLightArea(light *LightSourceComponent) {
	light.AreaDirty = false

	if light.litTiles != nil {
		ls.removeAppliedLight(light)
		clear(light.litTiles)
	} else {
		light.litTiles = make(map[vec.Coord]uint8)
	}

	if light.Disabled || light.Power == 0 {
		return
	}

	source := ecs.Get[PositionComponent](light.GetEntity()).Coord
	if source == NOT_IN_TILEMAP {
		return
	}

	lightRange := int(light.MaxRange)
	if lightRange == 0 {
		if light.FalloffRate == 0 {
			lightRange = 10
		} else {
			lightRange = int(float32(light.Power) / float32(light.FalloffRate))
		}
	}

	ls.tileMap.ShadowCast(source, lightRange, func(tm *TileMap, pos vec.Coord, d, r int) {
		light.litTiles[pos] = 0
	})

	light.Dirty = true
}

func (ls *LightSystem) applyLight(light *LightSourceComponent) {
	light.Dirty = false

	if light.Disabled {
		return
	}

	source := ecs.Get[PositionComponent](light.GetEntity()).Coord
	if source == NOT_IN_TILEMAP {
		return
	}

	for pos, oldLight := range light.litTiles {
		newLight := max(0, int(light.Power)-int(source.DistanceTo(pos)*float64(light.FalloffRate)))
		if delta := newLight - int(oldLight); delta != 0 {
			ls.tileMap.GetTile(pos).ModLight(delta)
			ls.tileMap.SetDirty(pos)
		}

		light.litTiles[pos] = uint8(newLight)
	}
}
