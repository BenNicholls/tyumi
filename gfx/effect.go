package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

type Effect func(*Cell)

// Applies the provided effect to each cell within the provided areas. If no areas are provided,
// applies it to the whole canvas. If multiple areas overlap, the effect will be applied multiple
// times!
func (c *Canvas) DrawEffect(effect Effect, areas ...vec.Rect) {
	if len(areas) == 0 {
		areas = append(areas, c.Bounds())
	}

	for _, area := range areas {
		for cursor := range vec.EachCoordInArea(area) {
			if !c.InBounds(cursor) {
				continue
			}
			cell := c.getCell(cursor)
			effect(cell)
		}
	}
}

func InvertEffect(cell *Cell) {
	cell.SetColours(col.Pair{cell.Colours.Back, cell.Colours.Fore})
}
