package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
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

var tileDataCache util.DataCache[TileData, TileType]
var TILE_NONE TileType

// The function Tyumi will use to draw tiles. This is by default set to the provided DrawTile function. The visuals
// returned here are NOT lit, that will be done later during drawing by the light system. This also does not draw
// forground content of tiles, like items or entities or actors or whatever, that is calculated by the equivalent
// DefaultTileEntityDrawFunction.
var DefaultTileDrawFunction func(tile Tile, viewer Entity) gfx.Visuals

// The function Tyumi will use to draw the entity on a given tile. It returns the computed visuals as well as the Entity
// that it drew.
var DefaultTileEntityDrawFunction func(tile Tile, viewer Entity) (gfx.Visuals, Entity)

func init() {
	TILE_NONE = RegisterTileType(TileData{Name: "No Tile", Desc: "A void in the universe."})

	DefaultTileDrawFunction = DrawTile
	DefaultTileEntityDrawFunction = DrawTileEntity
}

type Tile ecs.Entity

func CreateTile(tile_type TileType, pos vec.Coord) (tile Tile) {
	tile = Tile(ecs.CreateEntity())
	ecs.Add(tile, TerrainComponent{TileType: tile_type})
	ecs.Add(tile, PositionComponent{Coord: pos, Static: true})

	if tile_type.Data().Passable {
		ecs.Add[EntityContainerComponent](tile)
	}

	return
}

func DrawTile(tile Tile, viewer Entity) (vis gfx.Visuals) {
	return tile.GetVisuals()
}

func DrawTileEntity(tile Tile, viewer Entity) (vis gfx.Visuals, e Entity) {
	vis, e = gfx.Visuals{Mode: gfx.DRAW_NONE}, INVALID_ENTITY
	info := tile.GetTileType().Data()

	if !info.Passable {
		return
	}

	if entity := tile.GetEntity(); entity.IsValid() {
		vis, e = entity.GetVisuals(), entity
	}

	return
}

func (t Tile) GetTileType() TileType {
	return ecs.Get[TerrainComponent](t).TileType
}

func (t Tile) SetTileType(tile_type TileType) {
	ecs.Get[TerrainComponent](t).TileType = tile_type
	if tile_type.Data().Passable != ecs.Has[EntityContainerComponent](t) {
		ecs.Toggle[EntityContainerComponent](t)
	}
}

// GetVisuals returns the visuals for the tile with any animations applied.
func (t Tile) GetVisuals() (vis gfx.Visuals) {
	vis = ApplyVisualAnimations(t, t.GetRawVisuals())
	return
}

func (t Tile) GetRawVisuals() gfx.Visuals {
	return t.GetTileType().Data().Visuals
}

func (t Tile) Position() vec.Coord {
	return ecs.Get[PositionComponent](t).Coord
}

func (t Tile) IsPassable() bool {
	return t.GetTileType().Data().Passable && t.GetEntity() == Entity(ecs.INVALID_ID)
}

func (t Tile) IsOpaque() bool {
	return t.GetTileType().Data().Opaque
}

func (t Tile) GetEntity() Entity {
	if container := ecs.Get[EntityContainerComponent](t); container != nil {
		return container.Entity
	} else {
		return INVALID_ENTITY
	}
}

func (t Tile) HasEntity() bool {
	if container := ecs.Get[EntityContainerComponent](t); container != nil {
		return container.Entity.IsValid()
	} else {
		return false
	}
}

func (t Tile) RemoveEntity() {
	if container := ecs.Get[EntityContainerComponent](t); container != nil {
		container.Remove()
	}
}
