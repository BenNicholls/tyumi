package rl

import (
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// Shadowcatser runs 8 times over different quadrants. rMatrix supplies rotation coefficients. Linear algebra to the
// rescue.
var rMatrix = [8][4]int{{1, 0, 0, 1}, {-1, 0, 0, 1}, {0, -1, 1, 0}, {0, -1, -1, 0}, {-1, 0, 0, -1}, {1, 0, 0, -1}, {0, 1, -1, 0}, {0, 1, 1, 0}}

// THE BIG CHEESE - The one and only Shadowcaster! For all of your FOV needs. fn is a function for the shadowcaster to
// apply to open spaces it finds.
func (m *TileMap) ShadowCast(pos vec.Coord, radius int, fn Cast) {
	if radius <= 0 {
		return
	}
	fn(m, pos, 0, radius)
	for i := range 8 {
		m.scan(pos, 1, 1.0, 0.0, radius, rMatrix[i], i%2 == 0, fn)
	}
}

// NOTE: The 'cull' bool controls the logic for ensuring the 8 passes don't overlap at the edges. It is set to true for
// the odd-numbered scans. The shadowcaster still visits these squares twice, but the function fn is not run twice.
// Trust me Ben, this was the best way you could think of and your other solutions created crazy behaviour. Leave it alone!
func (tm *TileMap) scan(pos vec.Coord, row int, slope1, slope2 float32, radius int, r [4]int, cull bool, fn Cast) {
	if slope1 < slope2 {
		return
	}
	blocked := false
	r2 := radius * radius

	//scan #radius rows
	for j := row; j <= radius && !blocked; j++ {
		//scan row
		for dx, dy, newStart := -j, -j, slope1; dx <= 0; dx++ {
			mapPos := vec.Coord{pos.X + dx*r[0] + dy*r[1], pos.Y + dx*r[2] + dy*r[3]} //map coordinates
			if !tm.Bounds().Contains(mapPos) {
				continue
			}
			lSlope, rSlope := (float32(dx)-0.5)/(float32(dy)+0.5), (float32(dx)+0.5)/(float32(dy)-0.5)

			if newStart < rSlope {
				continue
			} else if slope2 > lSlope {
				break
			}

			if d := vec.ZERO_COORD.DistanceSqTo(vec.Coord{dx, dy}); d < r2 {
				if !cull || !(dx == 0 || dy == 0 || dx == dy) {
					fn(tm, mapPos, d, radius)
				}
			}
			
			//scanning a block
			if blocked {
				if !tm.IsTileOpaque(mapPos) {
					blocked = false
					slope1 = newStart
				} else {
					newStart = rSlope
				}
			} else {
				//blocked square, commence child scan
				if j < radius && tm.IsTileOpaque(mapPos) {
					blocked = true
					tm.scan(pos, j+1, newStart, lSlope, radius, r, cull, fn)
					newStart = rSlope
				}
			}
		}
	}
}

// type specifying precisely what you can pass to the shadowcaster. parameters here are the info that the shadowcaster
// will deliver. d2 is the distance squared from the center of the cast.
type Cast func(tm *TileMap, pos vec.Coord, d2, r int)

// Collects coords of all spaces visited by the shadowcater into the provided slice. Reslices the slice down
// before running the cast.
func GetSpacesCast(spaces []vec.Coord) Cast {
	spaces = spaces[0:0]
	return func(tm *TileMap, pos vec.Coord, d, r int) {
		spaces = append(spaces, pos)
	}
}

// Collects coords of all spaces visited by the shadowcater into the provided set. Clears the set before
// running the cast, so any contents in there will be destroyed.
func GetSpacesSetCast(spaces *util.Set[vec.Coord]) Cast {
	spaces.RemoveAll()
	return func(tm *TileMap, pos vec.Coord, d, r int) {
		spaces.Add(pos)
	}
}
