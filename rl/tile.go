package rl

import "github.com/bennicholls/tyumi/gfx"

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
	return tileDataCache.GetTileData(t.tileType).GetVisuals()
}
