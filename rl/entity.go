package rl

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

type EntityType uint32

type EntityData struct {
	Name    string
	Desc    string
	Glyph   gfx.Glyph
	Colours col.Pair
}

func (ed EntityData) GetVisuals() gfx.Visuals {
	return gfx.NewGlyphVisuals(ed.Glyph, ed.Colours)
}

var entityDataCache dataCache[EntityData, EntityType]

func RegisterEntityType(entity_data EntityData) EntityType {
	return entityDataCache.RegisterDataType(entity_data)
}

func init() {
	entityDataCache.Init()
}

type Entity struct {
	Name string
	Desc string

	entityType EntityType
	position   vec.Coord
}

func (e *Entity) Init(entity_type EntityType) {
	e.entityType = entity_type
}

func (e Entity) GetVisuals() gfx.Visuals {
	return entityDataCache.GetData(e.entityType).GetVisuals()
}
