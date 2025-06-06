package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/rl/ecs"
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

var tileDataCache dataCache[TileData, TileType]
var TILE_NONE TileType

func init() {
	tileDataCache.Init()
	TILE_NONE = RegisterTileType(TileData{Name: "No Tile", Desc: "A void in the universe."})
}

type Tile ecs.Entity

func CreateTile(tile_type TileType) (tile Tile) {
	tile = Tile(ecs.CreateEntity())
	ecs.AddComponent(tile, TerrainComponent{TileType: tile_type})

	if tile_type.Data().Passable {
		ecs.AddComponent[EntityContainerComponent](tile)
	}

	return
}

func (t Tile) GetTileType() TileType {
	return ecs.GetComponent[TerrainComponent](t).TileType
}

func (t Tile) SetTileType(tile_type TileType) {
	ecs.GetComponent[TerrainComponent](t).TileType = tile_type
	if tile_type.Data().Passable != ecs.HasComponent[EntityContainerComponent](t) {
		ecs.ToggleComponent[EntityContainerComponent](t)
	}
}

func (t Tile) IsPassable() bool {
	return t.GetTileType().Data().Passable && t.GetEntity() == nil
}

func (t Tile) IsTransparent() bool {
	return !t.GetTileType().Data().Opaque
}

func (t Tile) GetVisuals() gfx.Visuals {
	vis := t.GetTileType().Data().Visuals
	if entity := t.GetEntity(); entity != nil {
		entityVisuals := entity.GetVisuals()
		vis.Glyph = entityVisuals.Glyph
		vis.Colours.Fore = entityVisuals.Colours.Fore
		if entityVisuals.Colours.Back != col.NONE {
			vis.Colours.Back = entityVisuals.Colours.Back
		}
	}

	return vis
}

func (t Tile) GetEntity() TileMapEntity {
	if container := ecs.GetComponent[EntityContainerComponent](t); container != nil {
		return container.TileMapEntity
	} else {
		return nil
	}
}

func (t Tile) RemoveEntity() {
	if container := ecs.GetComponent[EntityContainerComponent](t); container != nil {
		container.TileMapEntity = nil
	}
}
