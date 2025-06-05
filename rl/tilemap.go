package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
)

type TileMap struct {
	size vec.Dims

	tiles []Tile
	dirty bool // if the tilemap has been changed and needs to be redrawn
}

// Initialize the TileMap. All tiles in the map will be set to defaultTile
func (tm *TileMap) Init(size vec.Dims, defaultTile TileType) {
	tm.size = size

	tm.tiles = make([]Tile, 0, size.Area())
	for range tm.Bounds().Area() {
		tm.tiles = append(tm.tiles, CreateTile(defaultTile))
	}
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
func (tm TileMap) SetTile(pos vec.Coord, tile Tile) {
	if !ecs.Alive(tile) {
		return
	}

	if !pos.IsInside(tm) {
		ecs.RemoveEntity(tile)
		return
	}

	ecs.RemoveEntity(tm.tiles[pos.ToIndex(tm.size.W)])
	tm.tiles[pos.ToIndex(tm.size.W)] = tile
	tm.dirty = true
}

func (tm *TileMap) SetTileType(pos vec.Coord, tileType TileType) {
	if !pos.IsInside(tm) {
		return
	}

	tm.tiles[pos.ToIndex(tm.size.W)].SetTileType(tileType)
	tm.dirty = true
}

func (tm *TileMap) AddEntity(entity TileMapEntity, pos vec.Coord) {
	if !pos.IsInside(tm) {
		return
	}

	tile := tm.GetTile(pos)
	if !ecs.Alive(tile) || !tile.IsPassable() {
		return
	}

	if container := ecs.GetComponent[EntityContainerComponent](tile); container != nil && container.Empty() {
		entity.MoveTo(pos)
		container.TileMapEntity = entity
		tm.dirty = true
	}
}

func (tm *TileMap) RemoveEntity(entity TileMapEntity) {
	tm.RemoveEntityAt(entity.Position())
}

func (tm *TileMap) RemoveEntityAt(pos vec.Coord) {
	if !pos.IsInside(tm) {
		return
	}

	tile := tm.GetTile(pos)
	if container := ecs.GetComponent[EntityContainerComponent](tile); container != nil && !container.Empty() {
		container.TileMapEntity.MoveTo(vec.Coord{-1, -1})
		container.TileMapEntity = nil
		tm.dirty = true
	}
}

func (tm *TileMap) GetEntityAt(pos vec.Coord) TileMapEntity {
	tile := tm.GetTile(pos)
	if !ecs.Valid(tile) {
		return nil
	}

	return tile.GetEntity()
}

func (tm *TileMap) MoveEntity(entity TileMapEntity, to vec.Coord) {
	from := entity.Position()
	if !from.IsInside(tm) || !to.IsInside(tm) {
		return
	}

	fromTile, toTile := tm.GetTile(from), tm.GetTile(to)
	if fromTile.GetEntity() != entity || !toTile.IsPassable() {
		return
	}

	ecs.GetComponent[EntityContainerComponent](toTile).TileMapEntity = entity
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
