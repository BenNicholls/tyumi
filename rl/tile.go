package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
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

type Tile struct {
	tileType TileType
	entity   TileMapEntity
}

func (t Tile) GetTileType() TileType {
	return t.tileType
}

func (t *Tile) SetTileType(tileType TileType) {
	if t.tileType == tileType {
		return
	}

	t.tileType = tileType
}

func (t Tile) GetVisuals() gfx.Visuals {
	vis := tileDataCache.GetData(t.tileType).GetVisuals()
	if t.entity != nil {
		entityVisuals := t.entity.GetVisuals()
		vis.Glyph = entityVisuals.Glyph
		vis.Colours.Fore = entityVisuals.Colours.Fore
		if entityVisuals.Colours.Back != col.NONE {
			vis.Colours.Back = entityVisuals.Colours.Back
		}
	}

	return vis
}

func (t Tile) IsPassable() bool {
	return t.entity == nil && tileDataCache.GetData(t.tileType).Passable
}
