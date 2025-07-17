package rl

import (
	"time"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type drawableTileMap interface {
	vec.Bounded
	gfx.DrawableReporter
	CalcTileVisuals(tile_pos vec.Coord) gfx.Visuals
	getMap() *TileMap
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
	tmv.SetCameraOffset(vec.ZERO_COORD)
}

func (tmv *TileMapView) SetTileMap(tilemap drawableTileMap) {
	tmv.tilemap = tilemap
	tmv.Updated = true
	tmv.tilemap.getMap().currentCameraBounds = vec.Rect{tmv.cameraOffset, tmv.Size()}
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
	if tmv.tilemap != nil {
		tmv.tilemap.getMap().currentCameraBounds = vec.Rect{offset, tmv.Size()}
	}

	tmv.ForceRedraw()
}

func (tmv TileMapView) GetCameraOffset() vec.Coord {
	return tmv.cameraOffset
}

// DrawTilemapObject draws an object to a position defined in tilemap-space.
// TODO: look at this more closely. i think this is old code that doesn't make sense any more.
func (tmv *TileMapView) DrawTilemapObject(object gfx.Drawable, tilemap_position vec.Coord, depth int) {
	object.Draw(&tmv.Canvas, tilemap_position.Add(tmv.cameraOffset), depth)
	tmv.Updated = true
}

func (tmv *TileMapView) Update(delta time.Duration) {
	if tmv.tilemap == nil {
		return
	}

	if tmv.tilemap.Dirty() {
		tmv.Updated = true
	}
}

func (tmv *TileMapView) Render() {
	if tmv.tilemap == nil {
		return
	}

	tilemap := tmv.tilemap.getMap()
	for cursor := range vec.EachCoordInIntersection(tmv, tmv.tilemap.Bounds().Translated(tmv.cameraOffset.Scale(-1))) {
		tileCursor := cursor.Add(tmv.cameraOffset)
		if tmv.IsRedrawing() || tilemap.IsDirtyAt(tileCursor) {
			tileVisuals := tmv.tilemap.CalcTileVisuals(tileCursor)
			if tileVisuals.Mode == gfx.DRAW_NONE {
				tmv.DrawVisuals(cursor, 0, tmv.DefaultVisuals())
			} else {
				tmv.DrawVisuals(cursor, 0, tileVisuals)
			}
		}
	}

	tmv.tilemap.Clean()
}
