package rl

import (
	"time"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/rl/ecs"
	"github.com/bennicholls/tyumi/util"
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
	if mlc.Text == txt {
		return
	}

	mlc.Text = txt
	event.Fire(EV_LABELUPDATED)
}

func (mlc *MapLabelComponent) SetColours(colours col.Pair) {
	if mlc.Colours == colours {
		return
	}

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

type MapLabelSystem struct {
	System

	labelLayer    *ui.Element
	movedEntities util.Set[Entity]
}

func (mls *MapLabelSystem) Init(labelLayer *ui.Element) {
	mls.labelLayer = labelLayer

	mls.SetImmediateEventHandler(mls.ImmediateHandleEvent)
	mls.Listen(EV_ENTITYMOVED, EV_ENTITYBEINGDESTROYED, EV_LABELUPDATED)
}

func (mls *MapLabelSystem) ImmediateHandleEvent(e event.Event) (event_handled bool) {
	if !mls.labelLayer.IsRedrawing() {
		switch e.ID() {
		case EV_LABELUPDATED:
			mls.labelLayer.ForceRedraw()
			event_handled = true
		case EV_ENTITYMOVED:
			mls.movedEntities.Add(e.(*EntityMovedEvent).Entity)
			event_handled = true
		}
	}

	switch e.ID() {
	case EV_ENTITYBEINGDESTROYED:
		entity := e.(*EntityEvent).Entity
		if !mls.labelLayer.IsRedrawing() {
			if ecs.Has[MapLabelComponent](entity) {
				mls.labelLayer.ForceRedraw()
			}
		}

		// if an entity is being destroyed, we need to check all labels to see if they are parented to it. if so, we
		// destroy the whole entity for the label (we're assuming the entity is entirely just for the label)
		// THINK: are there situations where this assumption is too much??
		for label, labelEntity := range ecs.EachComponent[MapLabelComponent]() {
			if label.Parent == entity {
				ecs.QueueDestroyEntity(labelEntity)
			}
		}
		event_handled = true
	}

	return
}

func (mls *MapLabelSystem) Update(delta time.Duration) {
	defer mls.movedEntities.RemoveAll()
	if mls.labelLayer.IsRedrawing() {
		return
	}

	for entity := range mls.movedEntities.EachElement() {
		if ecs.Has[MapLabelComponent](entity) {
			mls.labelLayer.ForceRedraw()
			break
		} else {
			for label := range ecs.EachComponent[MapLabelComponent]() {
				if label.Parent == entity {
					mls.labelLayer.ForceRedraw()
					break
				}
			}
		}
	}
}
