package gfx

import "github.com/bennicholls/tyumi/vec"

// Draws the provided visuals to the entire rectangular area. Optionally takes a 2nd Visuals that will be drawn
// on the interior of the rect, making the normal visuals will be the border.
func (c *Canvas) DrawFilledRect(area vec.Rect, depth int, brush Visuals, inner_brush ...Visuals) {
	if area.Area() == 0 {
		return
	}

	if len(inner_brush) == 0 || area.Area() == 1 {
		for cursor := range vec.EachCoordInArea(area) {
			c.DrawVisuals(cursor, depth, brush)
		}
	} else {
		for cursor := range vec.EachCoordInArea(area) {
			if cursor.X == area.X || cursor.Y == area.Y || cursor.X == area.X+area.W-1 || cursor.Y == area.Y+area.H-1 {
				c.DrawVisuals(cursor, depth, brush)
			} else {
				c.DrawVisuals(cursor, depth, inner_brush[0])
			}
		}
	}
}

// Draws the provided visuals to the border of the rectangular area.
func (c *Canvas) DrawRect(area vec.Rect, depth int, brush Visuals) {
	if area.Area() == 0 {
		return
	}

	//single-cell sized rect.
	if area.Area() == 1 {
		c.DrawVisuals(area.Coord, depth, brush)
		return
	}

	var sides [4]vec.Line
	corners := area.Corners()
	sides[0] = vec.Line{area.Coord, corners[1].Step(vec.DIR_LEFT)}  // top
	sides[1] = vec.Line{corners[1], corners[2].Step(vec.DIR_UP)}    // right
	sides[2] = vec.Line{corners[2], corners[3].Step(vec.DIR_RIGHT)} // bottom
	sides[3] = vec.Line{corners[3], area.Coord.Step(vec.DIR_DOWN)}  // left
	for _, line := range sides {
		c.DrawLine(line, depth, brush)
	}
}

// draws a circle of radius r centered at (px, py), copying the visuals from v, with option to fill the circle with same
// visuals
func (c *Canvas) DrawCircle(center vec.Coord, depth, r int, v Visuals, fill bool) {
	drawFunc := func(pos vec.Coord) {
		c.DrawVisuals(pos, depth, v)
	}

	vec.Circle(center, r, drawFunc)

	if fill {
		c.FloodFill(center, depth, v)
	}
}

func (c *Canvas) DrawLine(line vec.Line, depth int, vis Visuals) {
	for cursor := range line.EachCoord() {
		c.DrawVisuals(cursor, depth, vis)
	}
}
