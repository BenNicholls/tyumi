package rl

import (
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
)

type drawableTileMap interface {
	vec.Bounded
	gfx.DrawableReporter
	CalcTileVisuals(tile_pos vec.Coord) gfx.Visuals
	getMap() *TileMap
}

type TileMapView struct {
	ui.Element
	MapLabelSystem

	FocusedEntity Entity // Entity that the tile map view remains centered on.
	ViewingEntity Entity // Entity that the tile map is drawn from the perspective of.

	tilemap      drawableTileMap
	cameraOffset vec.Coord // area we're viewing

	labelLayer ui.Element
}

func NewTileMapView(size vec.Dims, pos vec.Coord, depth int, tilemap drawableTileMap) (tmv *TileMapView) {
	tmv = new(TileMapView)
	tmv.Init(size, pos, depth, tilemap)

	return
}

func (tmv *TileMapView) Init(size vec.Dims, pos vec.Coord, depth int, tilemap drawableTileMap) {
	tmv.Element.Init(size, pos, depth)
	tmv.TreeNode.Init(tmv)

	tmv.SetEventHandler(tmv.HandleEvent)
	tmv.Listen(EV_ENTITYMOVED, EV_ENTITYBEINGDESTROYED)

	tmv.SetTileMap(tilemap)
	tmv.SetCameraOffset(vec.ZERO_COORD)

	tmv.labelLayer.Init(size, vec.ZERO_COORD, 1)
	tmv.labelLayer.OnRender = tmv.RenderLabels
	tmv.labelLayer.SetDefaultVisuals(
		gfx.Visuals{
			Mode:    gfx.DRAW_NONE,
			Colours: tmv.DefaultColours(),
		})

	tmv.MapLabelSystem.Init(&tmv.labelLayer)

	tmv.AddChild(&tmv.labelLayer)
}

func (tmv *TileMapView) HandleEvent(e event.Event) (event_handled bool) {
	// tilemapview events
	switch e.ID() {
	case EV_ENTITYMOVED:
		entity := e.(*EntityMovedEvent).Entity
		if entity == tmv.FocusedEntity {
			tmv.CenterOnTileMapCoord(entity.Position())
			event_handled = true
		}
	case EV_ENTITYBEINGDESTROYED:
		entity := e.(*EntityEvent).Entity
		if entity == tmv.FocusedEntity {
			tmv.FocusedEntity = INVALID_ENTITY
			event_handled = true
		}
	}

	return
}

func (tmv *TileMapView) SetTileMap(tilemap drawableTileMap) {
	tmv.tilemap = tilemap
	tmv.Updated = true
	tmv.tilemap.getMap().currentCameraBounds = vec.Rect{tmv.cameraOffset, tmv.Size()}
	tmv.Clear()
}

func (tmv *TileMapView) CenterTileMap() {
	offset := vec.Coord{(tmv.tilemap.Bounds().W - tmv.Bounds().W) / 2, (tmv.tilemap.Bounds().H - tmv.Bounds().H) / 2}
	tmv.SetCameraOffset(offset)
}

func (tmv *TileMapView) CenterOnTileMapCoord(tilemap_pos vec.Coord) {
	offset := tilemap_pos.Subtract(vec.Coord{tmv.Bounds().W/2 - 1, tmv.Bounds().H/2 - 1})
	tmv.SetCameraOffset(offset)
}

func (tmv TileMapView) MapCoordToViewCoord(map_coord vec.Coord) vec.Coord {
	return map_coord.Subtract(tmv.cameraOffset)
}

func (tmv *TileMapView) MoveCamera(dx, dy int) {
	tmv.SetCameraOffset(tmv.cameraOffset.Add(vec.Coord{dx, dy}))
}

func (tmv *TileMapView) SetCameraOffset(offset vec.Coord) {
	tmv.cameraOffset = offset
	if tmv.tilemap != nil {
		tmv.tilemap.getMap().currentCameraBounds = vec.Rect{offset, tmv.Size()}
	}

	tmv.ForceRedraw()
	tmv.labelLayer.ForceRedraw()
}

func (tmv TileMapView) GetCameraOffset() vec.Coord {
	return tmv.cameraOffset
}

// DrawTilemapObject draws an object to a position defined in tilemap-space.
// TODO: look at this more closely. i think this is old code that doesn't make sense any more.
func (tmv *TileMapView) DrawTilemapObject(object gfx.Drawable, tilemap_position vec.Coord, depth int) {
	object.Draw(&tmv.Canvas, tilemap_position.Add(tmv.cameraOffset), depth)
	tmv.Updated = true
}

func (tmv *TileMapView) Update(delta time.Duration) {
	tmv.MapLabelSystem.Update(delta)

	if tmv.tilemap == nil {
		return
	}

	if tmv.tilemap.Dirty() {
		tmv.Updated = true
	}
}

func (tmv *TileMapView) Render() {
	if tmv.tilemap == nil {
		return
	}

	var fovComp *FOVComponent
	var memoryComp *MemoryComponent

	if tmv.ViewingEntity != INVALID_ENTITY {
		fovComp = ecs.Get[FOVComponent](tmv.ViewingEntity)
		memoryComp = ecs.Get[MemoryComponent](tmv.ViewingEntity)
	}

	tilemap := tmv.tilemap.getMap()
	for cursor := range vec.EachCoordInIntersection(tmv, tmv.tilemap.Bounds().Translated(tmv.cameraOffset.Scale(-1))) {
		tileCursor := cursor.Add(tmv.cameraOffset)
		if tmv.IsRedrawing() || tilemap.IsDirtyAt(tileCursor) {
			var tileVisuals gfx.Visuals
			tileVisuals.Mode = gfx.DRAW_NONE

			if fovComp == nil || fovComp.InFOV(tileCursor) {
				// if there is no viewing entity, or if the viewing entity does not have an FOV component, we just
				// assume the camera is omniscient. otherwise we check to see if the tile is in the fov we found.
				tileVisuals = tmv.tilemap.CalcTileVisuals(tileCursor)
			} else if memoryComp != nil {
				// otherwise we try to pull the visuals from the memory of the viewer, if it has one.
				if memory, ok := memoryComp.GetMemory(tileCursor); ok {
					tileVisuals = memory.Visuals
					tileVisuals.Colours = memoryComp.Colours
					tileVisuals.Colours = tileVisuals.Colours.Replace(col.NONE, tmv.DefaultColours())
				}
			}

			if tileVisuals.Mode == gfx.DRAW_NONE {
				tmv.DrawVisuals(cursor, 0, tmv.DefaultVisuals())
			} else {
				tmv.DrawVisuals(cursor, 0, tileVisuals)
			}
		}
	}

	tmv.tilemap.Clean()
}

func (tmv *TileMapView) RenderLabels() {
	var fovComp *FOVComponent
	if tmv.ViewingEntity != INVALID_ENTITY {
		fovComp = ecs.Get[FOVComponent](tmv.ViewingEntity)
	}

	for label := range ecs.EachComponent[MapLabelComponent]() {
		if fovComp != nil && !label.ShowOutOfFOV {
			// cull labels that are out of the viewing entity's FOV if necessary
			// if pos is NOT_IN_TILEMAP then this is an absolute label and we can draw it regardless.
			if pos := label.EntityPosition(); pos != NOT_IN_TILEMAP && !fovComp.InFOV(pos) {
				continue
			}
		}

		label.Draw(&tmv.labelLayer.Canvas, tmv.cameraOffset, 0)
	}
}
