package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/rl/ecs"
)

type TileType uint32

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

func CreateTile(tile_type TileType) (te Tile) {
	te = Tile(ecs.CreateEntity())
	ecs.AddComponent(te, TerrainComponent{TileType: tile_type})

	if tileDataCache.GetData(tile_type).Passable {
		ecs.AddComponent[EntityContainerComponent](te)
	}

	return
}

func (te Tile) GetTileType() TileType {
	return ecs.GetComponent[TerrainComponent](te).TileType
}

func (te Tile) SetTileType(tile_type TileType) {
	ecs.GetComponent[TerrainComponent](te).TileType = tile_type
	if passable := tileDataCache.GetData(tile_type).Passable; passable != ecs.HasComponent[EntityContainerComponent](te) {
		if passable {
			ecs.AddComponent[EntityContainerComponent](te)
		} else {
			ecs.RemoveComponent[EntityContainerComponent](te)
		}
	}
}

func (te Tile) IsPassable() bool {
	return tileDataCache.GetData(te.GetTileType()).Passable && te.GetEntity() == nil
}

func (te Tile) IsTransparent() bool {
	return !tileDataCache.GetData(te.GetTileType()).Opaque
}

func (te Tile) GetVisuals() gfx.Visuals {
	vis := tileDataCache.GetData(te.GetTileType()).Visuals
	if entity := te.GetEntity(); entity != nil {
		entityVisuals := entity.GetVisuals()
		vis.Glyph = entityVisuals.Glyph
		vis.Colours.Fore = entityVisuals.Colours.Fore
		if entityVisuals.Colours.Back != col.NONE {
			vis.Colours.Back = entityVisuals.Colours.Back
		}
	}

	return vis
}

func (te Tile) GetEntity() TileMapEntity {
	if container := ecs.GetComponent[EntityContainerComponent](te); container != nil {
		return container.TileMapEntity
	} else {
		return nil
	}
}

func (te Tile) RemoveEntity() {
	if container := ecs.GetComponent[EntityContainerComponent](te); container != nil {
		container.TileMapEntity = nil
	}
}
