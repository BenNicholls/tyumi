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
	amount [4]uint8 // UP, RIGHT, DOWN, LEFT
	pos    vec.Coord
}

type tileLight [4]uint16 // UP, RIGHT, DOWN, LEFT

func (tl tileLight) GetLevel(dir vec.Direction) uint16 {
	return tl[dir/2]
}

func (tl *tileLight) SetLevel(level uint8, dir vec.Direction) {
	tl[dir/2] = uint16(level)
}

func (tl tileLight) IsZero() bool {
	return tl == tileLight{0, 0, 0, 0}
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
	globalLight           uint8       // amount of light automatically applied to every tile. If 255, disables system.
	lightmap              []tileLight // the light applied to each tile!
}

func (ls *LightSystem) Init(tm *TileMap) {
	ls.tileMap = tm
	area := ls.tileMap.Bounds().Area()
	ls.lightmap = make([]tileLight, area, area)
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

func (ls LightSystem) GetLightLevel(pos, view_pos vec.Coord) uint8 {
	if !ls.Enabled {
		return 255
	}

	light := ls.getTileLight(pos)
	var viewedLight uint16
	if ls.tileMap.IsTileOpaque(pos) {
		if view_pos == NOT_IN_TILEMAP {
			// if viewer is NOT IN TILEMAP we assume they are omniscient, and we light the tile
			// with the largest light value among all directions
			for _, level := range light {
				viewedLight = max(viewedLight, level)
			}
		} else {
			viewDelta := view_pos.Subtract(pos)

			if viewDelta.X > 0 && !ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_RIGHT)) {
				viewedLight = light.GetLevel(vec.DIR_RIGHT)
			} else if viewDelta.X < 0 && !ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_LEFT)) {
				viewedLight = light.GetLevel(vec.DIR_LEFT)
			}

			if viewDelta.Y < 0 && !ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_UP)) {
				viewedLight = max(viewedLight, light.GetLevel(vec.DIR_UP))
			} else if viewDelta.Y > 0 && !ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_DOWN)) {
				viewedLight = max(viewedLight, light.GetLevel(vec.DIR_DOWN))
			}

			// sampling method for "technically invisible" corners that we still want to be
			// able to see because it is nicer. :)
			if viewedLight == 0 {
				if viewDelta.X < 0 && viewDelta.Y < 0 {
					if ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_UP)) && ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_LEFT)) {
						if !ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_UPLEFT)) {
							viewedLight = ls.getTileLight(pos.Step(vec.DIR_UPLEFT))[0]
						}
					}
				} else if viewDelta.X < 0 && viewDelta.Y > 0 {
					if ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_DOWN)) && ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_LEFT)) {
						if !ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_DOWNLEFT)) {
							viewedLight = ls.getTileLight(pos.Step(vec.DIR_DOWNLEFT))[0]
						}
					}
				} else if viewDelta.X > 0 && viewDelta.Y < 0 {
					if ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_UP)) && ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_RIGHT)) {
						if !ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_UPRIGHT)) {
							viewedLight = ls.getTileLight(pos.Step(vec.DIR_UPRIGHT))[0]
						}
					}
				} else if viewDelta.X > 0 && viewDelta.Y > 0 {
					if ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_DOWN)) && ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_RIGHT)) {
						if !ls.tileMap.IsTileOpaque(pos.Step(vec.DIR_DOWNRIGHT)) {
							viewedLight = ls.getTileLight(pos.Step(vec.DIR_DOWNRIGHT))[0]
						}
					}
				}
			}

		}
	} else {
		viewedLight = light.GetLevel(vec.DIR_UP)
	}

	return uint8(min(255, uint16(ls.globalLight)+viewedLight))
}

func (ls LightSystem) getTileLight(pos vec.Coord) *tileLight {
	if !ls.Enabled {
		return &tileLight{255, 255, 255, 255}
	}

	return &ls.lightmap[pos.ToIndex(ls.tileMap.size.W)]
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
func (ls *LightSystem) LightTileVisuals(vis gfx.Visuals, pos, view_pos vec.Coord) (lit_vis gfx.Visuals) {
	if !ls.Enabled {
		return vis
	}

	// TODO: this lighting function will act pretty weird if the backcolour is a light colour (like if something
	// inverts the tile colours) should probably do this better somehow....
	lit_vis = vis
	lit_vis.Colours.Fore = vis.Colours.Back.Lerp(vis.Colours.Fore, int(ls.GetLightLevel(pos, view_pos)), 255)

	return
}

var zeroLight = [4]uint8{0, 0, 0, 0}

func (ls *LightSystem) removeAppliedLight(light *LightSourceComponent) {
	for idx, photon := range light.photons {
		pos, amount := photon.pos, photon.amount
		if amount == zeroLight {
			continue
		}

		tileLight := ls.getTileLight(pos)

		for i, value := range photon.amount {
			tileLight[i] -= uint16(value)
		}
		light.photons[idx].amount = zeroLight
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
		light.photons = append(light.photons, photon{zeroLight, pos})
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

	source := Entity(light.GetEntity()).Position()
	if source == NOT_IN_TILEMAP {
		return
	}

	for idx, photon := range light.photons {
		pos := photon.pos
		amplitude := uint8(max(0, int(light.Power)-int(source.DistanceTo(pos)*float64(light.FalloffRate))))

		var newLight tileLight
		if ls.tileMap.IsTileOpaque(pos) {
			posDelta := source.Subtract(pos)

			if posDelta == vec.ZERO_COORD {
				newLight = tileLight{uint16(amplitude), uint16(amplitude), uint16(amplitude), uint16(amplitude)}
			} else {
				if posDelta.X > 0 {
					newLight.SetLevel(amplitude, vec.DIR_RIGHT)
				} else if posDelta.X < 0 {
					newLight.SetLevel(amplitude, vec.DIR_LEFT)
				} else {
					newLight.SetLevel(amplitude/2, vec.DIR_LEFT)
					newLight.SetLevel(amplitude/2, vec.DIR_RIGHT)
				}

				if posDelta.Y < 0 {
					newLight.SetLevel(amplitude, vec.DIR_UP)
				} else if posDelta.Y > 0 {
					newLight.SetLevel(amplitude, vec.DIR_DOWN)
				} else {
					newLight.SetLevel(amplitude/2, vec.DIR_DOWN)
					newLight.SetLevel(amplitude/2, vec.DIR_UP)
				}
			}
		} else {
			newLight = tileLight{uint16(amplitude), uint16(amplitude), uint16(amplitude), uint16(amplitude)}
		}

		tileLight := ls.getTileLight(pos)

		for i := range newLight {
			if photon.amount[i] == uint8(newLight[i]) {
				continue
			}

			delta := int(newLight[i]) - int(photon.amount[i])
			if delta > 0 {
				tileLight[i] += uint16(delta)
			} else if delta < 0 {
				tileLight[i] -= uint16(-delta)
			}

			light.photons[idx].amount[i] = uint8(newLight[i])

			ls.tileMap.SetDirty(pos)
		}
	}
}
