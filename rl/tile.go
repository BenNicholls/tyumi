package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
)

type TileType uint32

type TileData struct {
	Name     string
	Desc     string
	Glyph    gfx.Glyph
	Colours  col.Pair
	Passable bool
	Opaque   bool
}

func (td TileData) GetVisuals() gfx.Visuals {
	return gfx.NewGlyphVisuals(td.Glyph, td.Colours)
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
	return tileDataCache.GetData(t.tileType).GetVisuals()
}
