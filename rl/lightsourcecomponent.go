package rl

import (
	"math/rand"
	"time"

	"github.com/bennicholls/tyumi/anim"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

func init() {
	ecs.Register[LightSourceComponent]()
}

type LightSourceComponent struct {
	ecs.Component

	Disabled     bool
	Dirty        bool  // if true, light needs to be reapplied to its area
	AreaDirty    bool  // if true, light needs to recompute which tiles it is affecting
	Power        uint8 // light level applied at the source
	FalloffRate  uint8 // amount light level diminishes every 1 tile away from source
	MaxRange     uint8 // maximum range of the light. If 0 (default), the light's distance is computed from Power and FalloffRate
	FlickerSpeed uint8 // How many ticks between flickers (at 60 ticks per second). If 0, flickering is disabled. NOTE: changing this after init does not update the flicker yet.
	//Colour      col.Colour // Colour of light

	photons   []photon // a photon is an amount of light being applied to a specific location on the tilemap
	litbounds vec.Rect // rough bounding box containing all lit positions (and more, of course)

	flickerAnimation  *anim.Repeater
	basePower         uint8 // used when computing flickers
	baseFalloff       uint8 // used when computing flickers
	lastComputedRange uint8 // used to prevent unnecessary area computations while flickering
}

type photon struct {
	amount uint8
	pos    vec.Coord
}

func (lsc *LightSourceComponent) Init() {
	if !lsc.Disabled {
		lsc.AreaDirty = true
	}

	if lsc.FlickerSpeed > 0 {
		flickerTime := (time.Second / 60) * time.Duration(lsc.FlickerSpeed)
		flicker := anim.NewRepeaterAnimation(flickerTime, func() {
			if light := ecs.Get[LightSourceComponent](lsc.GetEntity()); light != nil {
				light.flicker()
			}
		})

		if !lsc.Disabled {
			flicker.Start()
		}

		AddAnimation(lsc.GetEntity(), &flicker, false)
		lsc.flickerAnimation = &flicker
	}
}

func (lsc *LightSourceComponent) Cleanup() {
	if lsc.flickerAnimation != nil {
		entity := lsc.GetEntity()
		if !ecs.Alive(entity) {
			return
		}

		if animComp := ecs.Get[AnimationComponent](entity); animComp != nil {
			animComp.RemoveAnimation(lsc.flickerAnimation)
		}
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
	if lsc.Disabled {
		if lsc.flickerAnimation != nil {
			lsc.flickerAnimation.Stop()
		}
	} else {
		if lsc.flickerAnimation != nil {
			lsc.flickerAnimation.Start()
		}
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

	if lsc.MaxRange == 0 && lsc.GetMaxRange() > lsc.lastComputedRange {
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

	if lsc.MaxRange == 0 && lsc.GetMaxRange() > lsc.lastComputedRange {
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

func (lsc *LightSourceComponent) GetMaxRange() (lightRange uint8) {
	if lsc.MaxRange != 0 {
		return lsc.MaxRange
	}

	if lsc.FalloffRate == 0 {
		return lsc.Power // 0 falloffrate makes the lightrange infinite, so let's avoid that yeah?
	}

	return uint8(float32(lsc.Power) / float32(lsc.FalloffRate))
}

func (lsc *LightSourceComponent) flicker() {
	if lsc.Disabled || lsc.FlickerSpeed == 0 {
		return
	}

	if lsc.basePower == 0 {
		lsc.basePower = lsc.Power
		lsc.baseFalloff = max(lsc.FalloffRate, 1)
	}

	f := float32(lsc.baseFalloff) - rand.Float32()*(float32(lsc.baseFalloff)/10)
	p := float32(lsc.basePower) - rand.Float32()*(float32(lsc.basePower)/10)

	lsc.SetFalloff(uint8(util.Clamp(f, 1, 255)))
	lsc.SetPower(uint8(min(p, 255)))
}

type LightSystem struct {
	System

	tileMap               *TileMap
	changedVisbilityTiles util.Set[vec.Coord]
	globalLight           uint8 // amount of light automatically applied to every tile. If 255, disables system.

	lightmap []uint16 // the light applied to each tile!
}

func (ls *LightSystem) Init(tm *TileMap) {
	ls.tileMap = tm
	area := ls.tileMap.Bounds().Area()
	ls.lightmap = make([]uint16, area, area)
	ls.Listen(EV_ENTITYMOVED, EV_TILECHANGEDVISIBILITY)
	ls.SetImmediateEventHandler(ls.immediateHandleEvent)
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

func (ls LightSystem) GetLightLevel(pos vec.Coord) uint8 {
	if !ls.Enabled {
		return 255
	}

	return uint8(min(255, uint16(ls.globalLight)+ls.lightmap[pos.ToIndex(ls.tileMap.size.W)]))
}

func (ls *LightSystem) immediateHandleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_ENTITYMOVED:
		moveEvent := e.(*EntityMovedEvent)
		if light := ecs.Get[LightSourceComponent](moveEvent.Entity); light != nil {
			light.AreaDirty = true
		}
	case EV_TILECHANGEDVISIBILITY:
		visEvent := e.(*TileChangedVisibilityEvent)
		ls.changedVisbilityTiles.Add(visEvent.Pos)
	default:
		return
	}

	return true
}

func (ls *LightSystem) Update(delta time.Duration) {
	if !ls.Enabled || ls.tileMap == nil {
		return
	}

	ls.System.Update(delta)

	for light := range ecs.EachComponent[LightSourceComponent]() {
		if light.Disabled {
			continue
		}

		// check if light should trigger an area update due to nearby tiles changing visibility
		if !light.AreaDirty {
			for pos := range ls.changedVisbilityTiles.EachElement() {
				if pos.IsInside(light.litbounds) {
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

	// light application has to go in a separate pass to prevent certain accumulation errors arising from weird
	// situations where tiles are replaced.
	for light := range ecs.EachComponent[LightSourceComponent]() {
		if light.Dirty && light.litbounds.Intersects(ls.tileMap.currentCameraBounds) {
			ls.applyLight(light)
		}
	}
}

// LightTileVisuals applies the light level at the position to the computed tile visuals.
func (ls *LightSystem) LightTileVisuals(vis gfx.Visuals, pos vec.Coord) (lit_vis gfx.Visuals) {
	if !ls.Enabled {
		return vis
	}

	// TODO: this lighting function will act pretty weird if the backcolour is a light colour (like if something
	// inverts the tile colours) should probably do this better somehow....
	lit_vis = vis
	lit_vis.Colours.Fore = vis.Colours.Back.Lerp(vis.Colours.Fore, int(ls.GetLightLevel(pos)), 255)

	return
}

func (ls *LightSystem) removeAppliedLight(light *LightSourceComponent) {
	for idx, photon := range light.photons {
		pos, amount := photon.pos, photon.amount
		if amount == 0 {
			continue
		}

		ls.removeLight(pos, amount)
		light.photons[idx].amount = 0
		ls.tileMap.SetDirty(pos)
	}
}

// computeLitArea runs the shadowcaster to determine the area lit by the light. if the light was already lighting
// some terrain, those terrains are de-lighted first. how delightful!
func (ls *LightSystem) computeLightArea(light *LightSourceComponent) {
	light.AreaDirty = false

	// clear photon array, we're recomputing from scratch!
	if light.photons != nil {
		ls.removeAppliedLight(light)
		light.photons = light.photons[0:0]
		light.litbounds = vec.Rect{}
	} else {
		light.photons = make([]photon, 0)
	}

	if light.Disabled || light.Power == 0 {
		return
	}

	source := ecs.Get[PositionComponent](light.GetEntity()).Coord
	if source == NOT_IN_TILEMAP {
		return
	}

	lightRange := light.GetMaxRange()
	ls.tileMap.ShadowCast(source, int(lightRange), func(tm *TileMap, pos vec.Coord, d, r int) {
		light.photons = append(light.photons, photon{0, pos})
		light.litbounds = light.litbounds.CalcExtendedRect(pos)
	})

	light.Dirty = true
	light.lastComputedRange = lightRange
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

	for idx, photon := range light.photons {
		pos := photon.pos
		newLight := max(0, int(light.Power)-int(source.DistanceTo(pos)*float64(light.FalloffRate)))
		if delta := newLight - int(photon.amount); delta != 0 {
			ls.modLight(pos, delta)
			ls.tileMap.SetDirty(pos)
		}

		light.photons[idx].amount = uint8(newLight)
	}
}

// Applies a delta to the light level on a tile.
func (ls *LightSystem) modLight(pos vec.Coord, delta int) {
	if delta > 0 {
		ls.addLight(pos, uint8(delta))
	} else if delta < 0 {
		ls.removeLight(pos, uint8(-delta))
	}
}

func (ls *LightSystem) addLight(pos vec.Coord, light uint8) {
	ls.lightmap[pos.ToIndex(ls.tileMap.size.W)] += uint16(light)
}

func (ls *LightSystem) removeLight(pos vec.Coord, light uint8) {
	level := ls.GetLightLevel(pos)
	if level < light {
		ls.lightmap[pos.ToIndex(ls.tileMap.size.W)] = 0
	} else {
		ls.lightmap[pos.ToIndex(ls.tileMap.size.W)] -= uint16(light)
	}
}
