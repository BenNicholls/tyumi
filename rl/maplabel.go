package rl

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/vec"
)

var EV_LABELUPDATED = event.Register("Map Label Updated.")

func init() {
	ecs.Register[MapLabelComponent]()
}

// The MapLabelComponent allows you to add a text label to be drawn on the tilemap! Labels can work in 3 ways:
//
//  1. They can be added as components to other RL entities (tiles, actors, whatever). The label will be drawn relative
//     to the entity's position.
//  2. They can be added as a component to any entity, and have their Parent field set to an entity to be attached to.
//     In this case the label's position will be relative to the parent.
//  3. They can be added to a generic entity with no position. In this case the label will be drawn at a fixed position
//     in the tilemapview based on the label's Offset field.
//
// Currently the label always draws in half text mode. If not set explicitly, the colours for the label will use the
// default colours for the tilemapview when drawn.
//
// Be sure to use the proper setter functions (SetText, SetColour, etc.) when modifying or animating the label to ensure
// the label is drawn correctly when updated. Have fun!
type MapLabelComponent struct {
	ecs.Component

	Parent       Entity    // if this is set, position of the label will be relative to the parent's position (if it has one)
	Offset       vec.Coord // an offset relative to the label's attached entity.
	Text         string
	Colours      col.Pair
	ShowOutOfFOV bool // if true, the label will be drawn even if the tilemapview's viewing entity cannot see the associated entity.
}

func (mlc *MapLabelComponent) Init() {
	event.Fire(EV_LABELUPDATED)
}

func (mlc *MapLabelComponent) Cleanup() {
	event.Fire(EV_LABELUPDATED)
}

func (mlc *MapLabelComponent) SetText(txt string) {
	mlc.Text = txt
	event.Fire(EV_LABELUPDATED)
}

func (mlc *MapLabelComponent) SetColours(colours col.Pair) {
	mlc.Colours = colours
	event.Fire(EV_LABELUPDATED)
}

func (mlc *MapLabelComponent) Move(delta vec.Coord) {
	mlc.MoveTo(mlc.Offset.Add(delta))
}

func (mlc *MapLabelComponent) MoveTo(pos vec.Coord) {
	if pos == mlc.Offset {
		return
	}

	mlc.Offset = pos
	event.Fire(EV_LABELUPDATED)
}

// Returns the postion of the entity this label is attached to. If the attached entity does not have a position,
// NOT_IN_TILEMAP is returned.
func (mlc *MapLabelComponent) EntityPosition() (pos vec.Coord) {
	pos = NOT_IN_TILEMAP

	if mlc.Parent != INVALID_ENTITY {
		if posComp := ecs.Get[PositionComponent](mlc.Parent); posComp != nil {
			pos = posComp.Coord
		}
	} else {
		if posComp := ecs.Get[PositionComponent](mlc.GetEntity()); posComp != nil {
			pos = posComp.Coord
		}
	}

	return
}

// Returns the position of the Label in tilemap coordinates. If the label has a Parent set, and that parent has a
// position component, the label's position will be relative to that. Otherwise it will be relative to the label's
// entity's position (if it has one).
func (mlc *MapLabelComponent) Position() (pos vec.Coord) {
	pos = mlc.EntityPosition()

	if pos == NOT_IN_TILEMAP {
		return
	}

	return pos.Add(mlc.Offset)
}

func (mlc *MapLabelComponent) Draw(canvas *gfx.Canvas, offset vec.Coord, depth int) {
	var pos vec.Coord

	// if the label's position is in the tilemap, we translate it to view space via the offset. otherwise, we assume
	// the label is an absolute view label and its position is just the label's offset.
	if mapPos := mlc.Position(); mapPos != NOT_IN_TILEMAP {
		pos = mapPos.Subtract(offset)
	} else {
		pos = mlc.Offset
	}

	bounds := vec.Rect{pos, vec.Dims{len(mlc.Text), 1}}
	if !bounds.Intersects(canvas.Size().Bounds()) {
		return
	}

	colours := mlc.Colours.Replace(col.NONE, canvas.DefaultColours())
	canvas.DrawHalfWidthText(pos, depth, mlc.Text, colours, gfx.DRAW_TEXT_LEFT)
}
