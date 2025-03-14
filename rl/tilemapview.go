package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type drawableTileMap interface {
	vec.Bounded
	gfx.DrawableReporter
}

type TileMapView struct {
	ui.Element

	tilemap      drawableTileMap
	cameraOffset vec.Coord // area we're viewing
}

func NewTileMapView(size vec.Dims, pos vec.Coord, depth int, tilemap drawableTileMap) (tmv *TileMapView) {
	tmv = new(TileMapView)
	tmv.Init(size, pos, depth, tilemap)

	return
}

func (tmv *TileMapView) Init(size vec.Dims, pos vec.Coord, depth int, tilemap drawableTileMap) {
	tmv.Element.Init(size, pos, depth)
	tmv.TreeNode.Init(tmv)

	tmv.tilemap = tilemap
	tmv.cameraOffset = vec.ZERO_COORD
}

func (tmv *TileMapView) Update() {
	if tmv.tilemap.Dirty() {
		tmv.Updated = true
	}
}

func (tmv *TileMapView) Render() {
	tmv.tilemap.Draw(&tmv.Canvas, tmv.cameraOffset, 0)
}
