package gfx

import "github.com/bennicholls/tyumi/vec"

// Draws the provided visuals to the rectangular area
func (c *Canvas) DrawRect(area vec.Rect, depth int, v Visuals) {
	for cursor := range vec.EachCoord(area) {
		c.DrawVisuals(cursor, depth, v)
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
