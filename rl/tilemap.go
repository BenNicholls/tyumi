package rl

import (
	"runtime"
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var NOT_IN_TILEMAP = vec.Coord{-1, -1}

type TileMap struct {
	gfx.DirtyTracker

	LightSystem
	FOVSystem
	AnimationSystem

	Ready bool // set this to true once level generation is complete! suppresses events while false.

	events     event.Stream
	size       vec.Dims
	tiles      []Tile
	opacityMap util.Bitset

	currentCameraBounds vec.Rect
}

func (tm *TileMap) getMap() *TileMap {
	return tm
}

// Initialize the TileMap. All tiles in the map will be set to defaultTile
func (tm *TileMap) Init(size vec.Dims, defaultTile TileType) {
	tm.DirtyTracker.Init(size)
	tm.events.Listen(EV_ENTITYBEINGDESTROYED, EV_TILECHANGEDVISIBILITY)
	tm.events.SetEventHandler(tm.handleEvent)
	tm.size = size
	tm.opacityMap.Init(size.Area())

	tm.LightSystem.Init(tm)
	tm.FOVSystem.Init(tm)
	tm.AnimationSystem.Init(tm)

	tm.tiles = make([]Tile, 0, size.Area())
	for cursor := range vec.EachCoordInArea(tm.Bounds()) {
		tm.tiles = append(tm.tiles, CreateTile(defaultTile, cursor))
	}

	if defaultTile.Data().Opaque {
		tm.opacityMap.SetAll()
	}

	runtime.AddCleanup(tm, func(tiles []Tile) {
		for _, tile := range tiles {
			ecs.DestroyEntity(tile)
		}

		tm.LightSystem.Shutdown()
		tm.FOVSystem.Shutdown()
		tm.events.DisableListening()
	}, tm.tiles)
}

// update tilemap-controlled systems
func (tm *TileMap) Update(delta time.Duration) {
	tm.AnimationSystem.Update(delta)
	tm.LightSystem.Update(delta)

	if tm.HasBlockingAnimation {
		return
	}

	tm.FOVSystem.Update(delta)
}

func (tm *TileMap) handleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_ENTITYBEINGDESTROYED:
		entity := e.(*EntityEvent).Entity

		// ensure entity being destroyed is in the tilemap
		pos := entity.Position()
		if !tm.Bounds().Contains(pos) || tm.GetTile(pos).GetEntity() != entity {
			return
		}

		if light := ecs.Get[LightSourceComponent](entity); light != nil {
			tm.removeAppliedLight(light)
		}

		tm.RemoveEntityAt(pos)
	case EV_TILECHANGEDVISIBILITY:
		o := e.(*TileChangedVisibilityEvent)
		tm.opacityMap.SetTo(o.Pos.ToIndex(tm.size.W), o.Opaque)
	default:
		return false
	}

	return true
}

func (tm TileMap) Size() vec.Dims {
	return tm.size
}

func (tm TileMap) Bounds() vec.Rect {
	return tm.size.Bounds()
}

func (tm TileMap) GetTile(pos vec.Coord) (tile Tile) {
	if !tm.Bounds().Contains(pos) {
		return
	}

	return tm.tiles[pos.ToIndex(tm.size.W)]
}

func (tm TileMap) IsTileOpaque(pos vec.Coord) bool {
	return tm.opacityMap.Get(pos.ToIndex(tm.size.W))
}

// Sets the tile at the provided position pos. If the set fails for whatever reason (pos out of bounds, etc.), the
// provided tile entity is destroyed.
func (tm *TileMap) SetTile(pos vec.Coord, tile Tile) {
	if !ecs.Alive(tile) {
		return
	}

	if !pos.IsInside(tm) {
		ecs.DestroyEntity(tile)
		return
	}

	oldTile := tm.tiles[pos.ToIndex(tm.size.W)]

	if newTileOpacity := tile.IsOpaque(); tm.IsTileOpaque(pos) != newTileOpacity {
		if tm.Ready {
			event.Fire(EV_TILECHANGEDVISIBILITY, &TileChangedVisibilityEvent{Pos: pos, Opaque: newTileOpacity})
		} else {
			tm.opacityMap.SetTo(pos.ToIndex(tm.size.W), newTileOpacity)
		}
	}

	ecs.DestroyEntity(oldTile)
	tm.tiles[pos.ToIndex(tm.size.W)] = tile
	tm.SetDirty(pos)
}

func (tm *TileMap) SetTileType(pos vec.Coord, tileType TileType) {
	if !pos.IsInside(tm) {
		return
	}

	tile := tm.tiles[pos.ToIndex(tm.size.W)]

	if tile.GetEntity() != INVALID_ENTITY && !tileType.Data().Passable {
		// do not do the switch if there's an entity and new tiletype can't hold an entity
		return
	}

	if newTileOpacity := tileType.Data().Opaque; tm.IsTileOpaque(pos) != newTileOpacity {
		if tm.Ready {
			event.Fire(EV_TILECHANGEDVISIBILITY, &TileChangedVisibilityEvent{Pos: pos, Opaque: newTileOpacity})
		} else {
			tm.opacityMap.SetTo(pos.ToIndex(tm.size.W), newTileOpacity)
		}
	}

	tm.tiles[pos.ToIndex(tm.size.W)].SetTileType(tileType)
	tm.SetDirty(pos)
}

func (tm *TileMap) AddEntity(entity Entity, pos vec.Coord) {
	if !pos.IsInside(tm) {
		return
	}

	tile := tm.GetTile(pos)
	if !ecs.Alive(tile) || !tile.IsPassable() {
		return
	}

	if container := ecs.Get[EntityContainerComponent](tile); container != nil && container.Empty() {
		entity.MoveTo(pos)
		container.Add(entity)
		tm.SetDirty(pos)
	}
}

func (tm *TileMap) RemoveEntity(entity Entity) {
	tm.RemoveEntityAt(entity.Position())
}

func (tm *TileMap) RemoveEntityAt(pos vec.Coord) {
	if !pos.IsInside(tm) {
		return
	}

	tile := tm.GetTile(pos)
	if !ecs.Valid(tile) {
		return
	}

	container := ecs.Get[EntityContainerComponent](tile)
	if container == nil || container.Empty() {
		return
	}

	entity := container.Entity
	entity.MoveTo(NOT_IN_TILEMAP)
	container.Remove()

	tm.SetDirty(pos)
}

func (tm *TileMap) GetEntityAt(pos vec.Coord) Entity {
	tile := tm.GetTile(pos)
	if !ecs.Valid(tile) {
		return Entity(ecs.INVALID_ID)
	}

	return tile.GetEntity()
}

func (tm *TileMap) MoveEntity(entity Entity, to vec.Coord) {
	if entity == INVALID_ENTITY {
		return
	}

	from := entity.Position()
	if !from.IsInside(tm) || !to.IsInside(tm) {
		return
	}

	fromTile, toTile := tm.GetTile(from), tm.GetTile(to)
	if fromTile.GetEntity() != entity || !toTile.IsPassable() {
		return
	}

	ecs.Get[EntityContainerComponent](toTile).Entity = entity
	tm.SetDirty(to)
	ecs.Get[EntityContainerComponent](fromTile).Entity = INVALID_ENTITY
	tm.SetDirty(from)

	entity.MoveTo(to)
}

func (tm TileMap) Draw(dst_canvas *gfx.Canvas, offset vec.Coord, depth int) {
	for cursor := range vec.EachCoordInIntersection(dst_canvas, tm.Bounds().Translated(offset)) {
		dst_canvas.DrawVisuals(cursor, depth, tm.CalcTileVisuals(cursor.Subtract(offset)))
	}

	tm.Clean()
}

func (tm *TileMap) CalcTileVisuals(pos vec.Coord) (vis gfx.Visuals) {
	tile := tm.GetTile(pos)
	terrain := ecs.Get[TerrainComponent](tile)
	if terrain.TileType == TILE_NONE {
		return gfx.Visuals{Mode: gfx.DRAW_NONE}
	}

	info := terrain.Data()
	light := tm.GetLightLevel(pos)
	if light == 0 {
		return gfx.NewGlyphVisuals(gfx.GLYPH_NONE, col.Pair{col.NONE, info.Visuals.Colours.Back})
	}

	vis = info.Visuals
	if tileAnimComp := ecs.Get[AnimationComponent](tile); tileAnimComp != nil {
		vis = tileAnimComp.ApplyVisualAnimations(vis)
	}

	if info.Passable {
		if entity := tile.GetEntity(); entity != INVALID_ENTITY {
			tileVis := vis
			vis = entity.GetVisuals()
			vis.Colours.Back = vis.Colours.Back.Replace(col.NONE, tileVis.Colours.Back)
		}
	}

	//Apply lighting!
	vis = tm.LightTileVisuals(vis, pos)

	return
}

func (tm TileMap) CopyToTileMap(dst_map *TileMap, offset vec.Coord) {
	for cursor := range vec.EachCoordInArea(tm) {
		dst_map.SetTile(cursor.Add(offset), Tile(ecs.CopyEntity(tm.GetTile(cursor))))
	}
}

func (tm TileMap) OutputToXP(filename string) {
	var canvas gfx.Canvas
	canvas.Init(tm.size)

	for i, tile := range tm.tiles {
		vis := tile.GetTileType().Data().Visuals
		canvas.DrawVisuals(vec.IndexToCoord(i, tm.size.W), 1, vis)
	}

	canvas.ExportToXP(filename)
}
