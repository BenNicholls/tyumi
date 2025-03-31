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

func (tmv *TileMapView) SetTileMap(tilemap drawableTileMap) {
	tmv.tilemap = tilemap
	tmv.Updated = true
	tmv.Clear()
}

func (tmv *TileMapView) SetCameraOffset(offset vec.Coord) {
	tmv.cameraOffset = offset
	tmv.Updated = true
}

func (tmv TileMapView) GetCameraOffset() vec.Coord {
	return tmv.cameraOffset
}

// DrawTilemapObject draws an object to a position defined in tilemap-space.
func (tmv *TileMapView) DrawTilemapObject(object gfx.Drawable, tilemap_position vec.Coord, depth int) {
	object.Draw(&tmv.Canvas, tilemap_position.Add(tmv.cameraOffset), depth)
	tmv.Updated = true
}

func (tmv *TileMapView) Update() {
	if tmv.tilemap != nil && tmv.tilemap.Dirty() {
		tmv.Updated = true
	}
}

func (tmv *TileMapView) Render() {
	if tmv.tilemap == nil {
		return
	}
	
	tmv.tilemap.Draw(&tmv.Canvas, tmv.cameraOffset, 0)
}
