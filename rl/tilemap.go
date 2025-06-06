package rl

import (
	"runtime"
	"slices"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var NOT_IN_TILEMAP = vec.Coord{-1, -1}

type TileMap struct {
	size vec.Dims

	tiles                  []Tile
	visibilityChangedTiles util.Set[vec.Coord]

	entities []Entity
	dirty    bool // if the tilemap has been changed and needs to be redrawn
}

// Initialize the TileMap. All tiles in the map will be set to defaultTile
func (tm *TileMap) Init(size vec.Dims, defaultTile TileType) {
	tm.size = size

	tm.tiles = make([]Tile, 0, size.Area())
	for range tm.Bounds().Area() {
		tm.tiles = append(tm.tiles, CreateTile(defaultTile))
	}

	runtime.AddCleanup(tm, func(tiles []Tile) {
		for _, tile := range tiles {
			ecs.RemoveEntity(tile)
		}
	}, tm.tiles)
}

// update tilemap-controlled systems
func (tm *TileMap) Update() {
	// FOV updates
	for fov := range ecs.EachComponent[FOVComponent]() {
		if !fov.Dirty {
			for pos := range tm.visibilityChangedTiles.EachElement() {
				fov.OnEnvironmentChange(pos)
			}
		}

		if fov.Dirty {
			fov.UpdateFOV(tm)
		}
	}

	tm.visibilityChangedTiles.RemoveAll()
}

func (tm TileMap) Size() vec.Dims {
	return tm.size
}

func (tm TileMap) Bounds() vec.Rect {
	return tm.size.Bounds()
}

func (tm *TileMap) Clean() {
	tm.dirty = false
}

func (tm TileMap) GetTile(pos vec.Coord) (tile Tile) {
	if !pos.IsInside(tm) {
		return
	}

	return tm.tiles[pos.ToIndex(tm.size.W)]
}

// Sets the tile at the provided position pos. If the set fails for whatever reason (pos out of bounds, etc.), the
// provided tile entity is destroyed.
func (tm *TileMap) SetTile(pos vec.Coord, tile Tile) {
	if !ecs.Alive(tile) {
		return
	}

	if !pos.IsInside(tm) {
		ecs.RemoveEntity(tile)
		return
	}

	oldTile := tm.tiles[pos.ToIndex(tm.size.W)]

	if oldTile.IsOpaque() != tile.IsOpaque() {
		tm.visibilityChangedTiles.Add(pos)
	}

	ecs.RemoveEntity(oldTile)
	tm.tiles[pos.ToIndex(tm.size.W)] = tile
	tm.dirty = true
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

	if tile.IsOpaque() != tileType.Data().Opaque {
		tm.visibilityChangedTiles.Add(pos)
	}

	tm.tiles[pos.ToIndex(tm.size.W)].SetTileType(tileType)
	tm.dirty = true
}

func (tm *TileMap) AddEntity(entity Entity, pos vec.Coord) {
	if !pos.IsInside(tm) {
		return
	}

	tile := tm.GetTile(pos)
	if !ecs.Alive(tile) || !tile.IsPassable() {
		return
	}

	if container := ecs.GetComponent[EntityContainerComponent](tile); container != nil && container.Empty() {
		entity.MoveTo(pos)
		container.Add(entity)
		tm.entities = append(tm.entities, entity)
		tm.dirty = true

		if fov := ecs.GetComponent[FOVComponent](entity); fov != nil {
			fov.Dirty = true
		}
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
	if container := ecs.GetComponent[EntityContainerComponent](tile); container != nil && !container.Empty() {
		container.Entity.MoveTo(NOT_IN_TILEMAP)
		tm.entities = slices.DeleteFunc(tm.entities, func(e Entity) bool {
			return e == container.Entity
		})
		container.Remove()
		tm.dirty = true
	}
}

func (tm *TileMap) GetEntityAt(pos vec.Coord) Entity {
	tile := tm.GetTile(pos)
	if !ecs.Valid(tile) {
		return Entity(ecs.INVALID_ID)
	}

	return tile.GetEntity()
}

func (tm *TileMap) MoveEntity(entity Entity, to vec.Coord) {
	from := entity.Position()
	if !from.IsInside(tm) || !to.IsInside(tm) {
		return
	}

	fromTile, toTile := tm.GetTile(from), tm.GetTile(to)
	if fromTile.GetEntity() != entity || !toTile.IsPassable() {
		return
	}

	ecs.GetComponent[EntityContainerComponent](toTile).Entity = entity
	fromTile.RemoveEntity()
	entity.MoveTo(to)
	tm.dirty = true
}

func (tm TileMap) Dirty() bool {
	return tm.dirty
}

func (tm TileMap) Draw(dst_canvas *gfx.Canvas, offset vec.Coord, depth int) {
	for cursor := range vec.EachCoordInIntersection(dst_canvas, tm.Bounds().Translated(offset)) {
		dst_canvas.DrawVisuals(cursor, depth, tm.CalcTileVisuals(cursor.Subtract(offset)))
	}

	tm.dirty = false
}

func (tm TileMap) CalcTileVisuals(pos vec.Coord) gfx.Visuals {
	tile := tm.GetTile(pos)
	if tile.GetTileType() == TILE_NONE {
		return gfx.Visuals{Mode: gfx.DRAW_NONE}
	} else {
		return tile.GetVisuals()
	}
}

func (tm TileMap) CopyToTileMap(dst_map *TileMap, offset vec.Coord) {
	for cursor := range vec.EachCoordInArea(tm) {
		dst_map.SetTile(cursor.Add(offset), Tile(ecs.CopyEntity(tm.GetTile(cursor))))
	}
}
