package rl

import (
	"github.com/bennicholls/tyumi/gfx"
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

	tm.tiles = make([]Tile, size.Area())
	for cursor := range vec.EachCoordInArea(tm) {
		tm.SetTileType(cursor, defaultTile)
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

func (tm *TileMap) getTile(pos vec.Coord) (tile_ptr *Tile) {
	return &tm.tiles[pos.ToIndex(tm.size.W)]
}

func (tm *TileMap) SetTile(pos vec.Coord, tile Tile) {
	if !pos.IsInside(tm) {
		return
	}

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

	tile := tm.getTile(pos)
	if !tile.IsPassable() {
		return
	}

	entity.MoveTo(pos)
	tile.entity = entity
	tm.dirty = true
}

func (tm *TileMap) RemoveEntity(entity TileMapEntity) {
	pos := entity.Position()
	if !pos.IsInside(tm) {
		return
	}

	if tile := tm.getTile(pos); tile.entity == entity {
		tm.RemoveEntityAt(pos)
	}
}

func (tm *TileMap) RemoveEntityAt(pos vec.Coord) {
	if !pos.IsInside(tm) {
		return
	}

	if tile := tm.getTile(pos); tile.entity != nil {
		tile.entity.MoveTo(vec.Coord{-1, -1})
		tile.entity = nil
		tm.dirty = true
	}
}

func (tm *TileMap) MoveEntity(entity TileMapEntity, to vec.Coord) {
	from := entity.Position()
	if !from.IsInside(tm) || !to.IsInside(tm) {
		return
	}

	fromTile, toTile := tm.getTile(from), tm.getTile(to)
	if fromTile.entity != entity || !toTile.IsPassable() {
		return
	}

	toTile.entity = fromTile.entity
	fromTile.entity = nil
	toTile.entity.MoveTo(to)
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
	if tile.tileType == TILE_NONE {
		return gfx.Visuals{Mode: gfx.DRAW_NONE}
	} else {
		return tile.GetVisuals()
	}
}

func (tm TileMap) CopyToTileMap(dst_map *TileMap, offset vec.Coord) {
	for cursor := range vec.EachCoordInArea(tm) {
		dst_map.SetTile(cursor.Add(offset), tm.GetTile(cursor))
	}
}
