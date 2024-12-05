package gfx

import (
	"github.com/bennicholls/tyumi/vec"
)

type Effect func(*Cell)

//Applies the provided effect to each cell within the provided areas. If no areas are provided,
//applies it to the whole canvas. If multiple areas overlap, the effect will be applied multiple
//times!
func (c *Canvas) DrawEffect(effect Effect, areas ...vec.Rect) {
	if len(areas) == 0 {
		areas = append(areas, c.Bounds())
	}

	for _, area := range areas {
		offset := vec.ZERO_COORD
		for i := 0; i < area.Area(); i++ {
			offset = vec.Coord{i%area.W, i/area.W}
			if cell := c.GetCell(area.Coord.Add(offset)); cell != nil {
				effect(cell)
			}
		}
	}
}

func InvertEffect(cell *Cell) {
	f := cell.ForeColour
	cell.SetForeColour(cell.BackColour)
	cell.SetBackColour(f)
}
