package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

// An effect is a function that takes a set of visuals and returns a transformed set of visuals. An effect can be
// passed to Canvas.DrawEffect() to apply the effect to areas of the canvas
type Effect func(Visuals) Visuals

// Applies the provided effect to each cell within the provided areas. If no areas are provided, applies it to the whole
// canvas. If multiple areas overlap, the effect will be applied multiple times!
func (c *Canvas) DrawEffect(effect Effect, areas ...vec.Rect) {
	if len(areas) == 0 {
		areas = append(areas, c.Bounds())
	}

	for _, area := range areas {
		for cursor := range vec.EachCoordInIntersection(c, area) {
			v := effect(c.getCell(cursor))
			c.setCell(cursor, c.getDepth(cursor), v)
		}
	}
}

func InvertEffect(visuals Visuals) Visuals {
	visuals.Colours = col.Pair{visuals.Colours.Back, visuals.Colours.Fore}
	return visuals
}
