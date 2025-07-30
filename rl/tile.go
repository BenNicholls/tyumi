package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type TileType uint32

func (tt TileType) Data() TileData {
	return tileDataCache.GetData(tt)
}

type TileData struct {
	Name     string
	Desc     string
	Visuals  gfx.Visuals
	Passable bool
	Opaque   bool
}

func (td TileData) GetVisuals() gfx.Visuals {
	return td.Visuals
}

func RegisterTileType(tileData TileData) TileType {
	return tileDataCache.RegisterDataType(tileData)
}

var tileDataCache util.DataCache[TileData, TileType]
var TILE_NONE TileType

func init() {
	TILE_NONE = RegisterTileType(TileData{Name: "No Tile", Desc: "A void in the universe."})
}

type Tile ecs.Entity

func CreateTile(tile_type TileType, pos vec.Coord) (tile Tile) {
	tile = Tile(ecs.CreateEntity())
	ecs.Add(tile, TerrainComponent{TileType: tile_type})
	ecs.Add(tile, PositionComponent{Coord: pos, Static: true})

	if tile_type.Data().Passable {
		ecs.Add[EntityContainerComponent](tile)
	}

	return
}

func (t Tile) GetTileType() TileType {
	return ecs.Get[TerrainComponent](t).TileType
}

func (t Tile) SetTileType(tile_type TileType) {
	ecs.Get[TerrainComponent](t).TileType = tile_type
	if tile_type.Data().Passable != ecs.Has[EntityContainerComponent](t) {
		ecs.Toggle[EntityContainerComponent](t)
	}
}

func (t Tile) Position() vec.Coord {
	return ecs.Get[PositionComponent](t).Coord
}

func (t Tile) IsPassable() bool {
	return t.GetTileType().Data().Passable && t.GetEntity() == Entity(ecs.INVALID_ID)
}

func (t Tile) IsOpaque() bool {
	return t.GetTileType().Data().Opaque
}

func (t Tile) GetEntity() Entity {
	if container := ecs.Get[EntityContainerComponent](t); container != nil {
		return container.Entity
	} else {
		return INVALID_ENTITY
	}
}

func (t Tile) RemoveEntity() {
	if container := ecs.Get[EntityContainerComponent](t); container != nil {
		container.Remove()
	}
}
