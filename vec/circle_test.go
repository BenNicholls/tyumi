package vec

import "testing"

func TestCircleEachCoordInPerimeterCount(t *testing.T) {
	c := Circle{Radius: 5, Center: Coord{0, 0}}

	count := 0
	for range c.EachCoordInPerimeter() {
		count++
	}

	if count == 0 {
		t.Fatal("expected circle perimeter to produce points")
	}
}

func TestCircleEachCoordInPerimeterDistanceFromCenter(t *testing.T) {
	center := Coord{10, 10}
	radius := 5
	c := Circle{Radius: radius, Center: center}

	for coord := range c.EachCoordInPerimeter() {
		distSq := coord.DistanceSqTo(center)
		// bresenham-style circles are approximate, so allow some tolerance around radius^2
		if distSq < (radius-1)*(radius-1) || distSq > (radius+1)*(radius+1) {
			t.Errorf("perimeter point %v is too far from radius %d (distSq=%d)", coord, radius, distSq)
		}
	}
}

func TestCircleEachCoordInPerimeterEarlyStop(t *testing.T) {
	c := Circle{Radius: 10, Center: Coord{0, 0}}

	count := 0
	for range c.EachCoordInPerimeter() {
		count++
		if count == 3 {
			break
		}
	}

	if count != 3 {
		t.Errorf("expected iteration to stop after 3 coords, got %d", count)
	}
}

func TestCircleFunc(t *testing.T) {
	center := Coord{5, 5}
	var visited []Coord

	CircleFunc(center, 3, func(pos Coord) {
		visited = append(visited, pos)
	})

	if len(visited) == 0 {
		t.Fatal("CircleFunc did not call fn for any points")
	}
}

func TestCircleZeroRadius(t *testing.T) {
	c := Circle{Radius: 0, Center: Coord{5, 5}}

	var got []Coord
	for coord := range c.EachCoordInPerimeter() {
		got = append(got, coord)
	}

	// a zero-radius circle should just be the center point (possibly repeated across octants)
	for _, coord := range got {
		if coord != c.Center {
			t.Errorf("zero-radius circle produced point %v, want only %v", coord, c.Center)
		}
	}
}
