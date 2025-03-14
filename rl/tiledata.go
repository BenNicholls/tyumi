package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
)

var tileDataCache TileDataCache

type TileDataCache struct {
	cache []TileData
}

func (tdc *TileDataCache) Init() {
	tdc.cache = make([]TileData, 0)
}

func (tdc TileDataCache) ValidTileType(tile_type TileType) bool {
	return int(tile_type) < len(tdc.cache)
}

func (tdc TileDataCache) GetTileData(tile_type TileType) (td TileData) {
	if !tdc.ValidTileType(tile_type) {
		log.Error("TileType not registered.")
		return
	}

	return tdc.cache[tile_type]
}

func (tdc *TileDataCache) RegisterTileType(tile_data TileData) TileType {
	tdc.cache = append(tdc.cache, tile_data)
	return TileType(len(tdc.cache) - 1)
}

var TILE_NONE TileType

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

type TileType uint32

func RegisterTileType(tileData TileData) TileType {
	return tileDataCache.RegisterTileType(tileData)
}

func init() {
	tileDataCache.Init()
	TILE_NONE = RegisterTileType(TileData{Name: "No Tile", Desc: "A void in the universe."})
}
