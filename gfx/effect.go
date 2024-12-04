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
		areas = append(areas, vec.Rect{vec.ZERO_COORD, c.Size()})
	}

	for _, area := range areas {
		for i := 0; i < area.W*area.H; i++ {
			if cell := c.GetCell(i%area.W+area.X, i/area.W+area.Y); cell != nil {
				effect(cell)
			}
		}
	}
}

func InvertEffect(cell *Cell) {
	f := cell.ForeColour
	cell.SetForeColour(cell.Z, cell.BackColour)
	cell.SetBackColour(cell.Z, f)
}
