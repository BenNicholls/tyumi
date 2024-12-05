package vec

//returns a generator that computes successive coordinates representing 1/8th of a circle. rotate the arc to draw
//circles. gives back the ZERO_COORD when it is done. 
func ArcGenerator(radius int) func() Coord {
	x, y := 0, radius
	f := 1 - radius
	ddf_x, ddf_y := 1, -2*radius

	return func() Coord {
		if x <= y {
			c := Coord{x, y}
			if f >= 0 {
				y--
				ddf_y += 2
				f += ddf_y
			}

			x++
			ddf_x += 2
			f += ddf_x

			return c
		}

		return ZERO_COORD
	}
}

//Computes a circle, calling fn on each point of the circle. fn can be a drawing function or whatever.
func Circle(center Coord, radius int, fn func(pos Coord)) {
	c := ArcGenerator(radius)
	for p := c(); p != ZERO_COORD; p = c() {
		fn(Coord{center.X+p.X, center.Y+p.Y})
		fn(Coord{center.X+p.Y, center.Y+p.X})
		fn(Coord{center.X-p.Y, center.Y+p.X})
		fn(Coord{center.X-p.X, center.Y+p.Y})
		fn(Coord{center.X-p.X, center.Y-p.Y})
		fn(Coord{center.X-p.Y, center.Y-p.X})
		fn(Coord{center.X+p.Y, center.Y-p.X})
		fn(Coord{center.X+p.X, center.Y-p.Y})
	}
}