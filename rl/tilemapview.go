package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type drawableTileMap interface {
	vec.Bounded
	gfx.DrawableReporter
	CalcTileVisuals(tile_pos vec.Coord) gfx.Visuals
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

func (tmv *TileMapView) CenterTileMap() {
	offset := vec.Coord{(tmv.tilemap.Bounds().W - tmv.Bounds().W) / 2, (tmv.tilemap.Bounds().H - tmv.Bounds().H) / 2}
	tmv.SetCameraOffset(offset)
}

func (tmv *TileMapView) CenterOnTileMapCoord(tilemap_pos vec.Coord) {
	offset := tilemap_pos.Subtract(vec.Coord{tmv.Bounds().W/2 - 1, tmv.Bounds().H/2 - 1})
	tmv.SetCameraOffset(offset)
}

func (tmv TileMapView) MapCoordToViewCoord(map_coord vec.Coord) vec.Coord {
	return map_coord.Subtract(tmv.cameraOffset)
}

func (tmv *TileMapView) MoveCamera(dx, dy int) {
	tmv.SetCameraOffset(tmv.cameraOffset.Add(vec.Coord{dx, dy}))
}

func (tmv *TileMapView) SetCameraOffset(offset vec.Coord) {
	tmv.cameraOffset = offset
	tmv.ForceRedraw()
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

	tmv.tilemap.Draw(&tmv.Canvas, tmv.cameraOffset.Scale(-1), 0)
	offset := tmv.cameraOffset.Scale(-1)
	for cursor := range vec.EachCoordInIntersection(tmv, tmv.tilemap.Bounds().Translated(offset)) {
		tileVisuals := tmv.tilemap.CalcTileVisuals(cursor.Subtract(offset))
		if tileVisuals.Mode == gfx.DRAW_NONE {
			tmv.DrawVisuals(cursor, 0, tmv.DefaultVisuals())
		} else {
			tmv.DrawVisuals(cursor, 0, tileVisuals)
		}
	}

	tmv.tilemap.Clean()
}
