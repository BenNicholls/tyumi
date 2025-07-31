package gfx

import (
	"cmp"
	"iter"
	"slices"

	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// DirtyTracker is an embeddable type that implements tracking of dirty coordinates in a grid.
type DirtyTracker struct {
	dirty    util.Bitset
	allDirty bool
	stride   int
}

func (dt *DirtyTracker) Init(size vec.Dims) {
	dt.dirty.Init(size.Area())
	dt.stride = size.W
}

func (dt DirtyTracker) IsDirtyAt(pos vec.Coord) bool {
	return dt.isDirtyAtIndex(pos.ToIndex(dt.stride))
}

func (dt DirtyTracker) isDirtyAtIndex(idx int) bool {
	return dt.allDirty || dt.dirty.Get(idx)
}

func (dt DirtyTracker) Dirty() bool {
	return dt.allDirty || !dt.dirty.IsEmpty()
}

func (dt *DirtyTracker) SetDirty(pos vec.Coord) {
	if dt.allDirty {
		return
	}

	dt.dirty.Set(pos.ToIndex(dt.stride))
}

func (dt *DirtyTracker) setDirtyAtIndex(idx int) {
	if dt.allDirty {
		return
	}

	dt.dirty.Set(idx)
}

func (dt *DirtyTracker) SetAllDirty() {
	dt.allDirty = true
}

func (dt DirtyTracker) CountDirty() int {
	return dt.dirty.Count()
}

func (dt *DirtyTracker) Clean() {
	dt.dirty.Clear()
	dt.allDirty = false
}

// SpecialDirtyTracker was an attempt to make a dirty tracker with minimal memory overhead. It accepts a certain number
// of Coords; after this limit it merges the Coords in Rects, with calculations to ensure only the smallest rects are
// created during the merging process. It also tries expanding existing areas, in case adding to a current area
// is better than making a new one.
//
// The system works great, but is significantly slower than the bitset-based dirty tracker. It may be worth coming
// back to this later on when more complex levels and behaviour are in. I suspect this method may show its
// quality on large maps with things happening in areas that are far away from eachother. A 200x200 tilemap
// would need a 5MB bitset, almost all of which would be zeros at any time, but the SpecialDirtyTracker but be able to
// adequately represent the same info with 20 Coords and 10 Rects, especially if a lot of the level is static.
//
// Apart from having a different Init() function signature, this is a direct drop-in replacement for the regular
// DirtyTracker.
type SpecialDirtyTracker struct {
	singles util.OrderedSet[vec.Coord]
	areas   []vec.Rect

	MaxSingles int
	MaxAreas   int
}

func (dt *SpecialDirtyTracker) Init(max_singles, max_areas int) {
	dt.MaxSingles = max_singles
	dt.MaxAreas = max_areas
	dt.areas = make([]vec.Rect, 0, max_areas)
}

func (dt *SpecialDirtyTracker) SetDirty(pos vec.Coord) {
	if dt.IsDirtyAt(pos) {
		return
	}

	dt.singles.Add(pos)

	if dt.singles.Count() <= dt.MaxSingles {
		return
	}

	var bestNew, bestExpanded vec.Rect
	var expandedDelta, area1 int

	// find best rect that combines 2 singles
	if len(dt.areas) < dt.MaxAreas {
		bestNew = dt.calcBestNewRect()
	}

	expandedDelta = 1000000
	// find smallest expansion of an area and a single
	for i, area := range dt.areas {
		smallestDim := min(area.H, area.W)
		if smallestDim >= expandedDelta {
			continue
		}
		for _, single := range dt.singles.EachElement() {
			rect := area.CalcExtendedRect(single)
			if delta := rect.Area() - area.Area(); delta < expandedDelta {
				bestExpanded = rect
				expandedDelta = delta
				area1 = i
				if delta == smallestDim {
					break
				}
			}
		}
	}

	if bestNew.Area() < expandedDelta && bestNew.Area() != 0 {
		dt.singles.RemoveFunc(func(single vec.Coord) bool {
			return bestNew.Contains(single)
		})

		dt.areas = slices.DeleteFunc(dt.areas, func(area vec.Rect) bool {
			return area.IsInside(bestNew)
		})

		dt.areas = append(dt.areas, bestNew)
	} else if bestExpanded.Area() != 0 {
		dt.areas[area1] = bestExpanded
		//remove singles that are inside the new expanded area
		dt.singles.RemoveFunc(func(single vec.Coord) bool {
			return bestExpanded.Contains(single)
		})

		//remove any areas consumed by the new expanded area
		dt.areas = slices.DeleteFunc(dt.areas, func(area vec.Rect) bool {
			return area.IsInside(bestExpanded)
		})
	}

	slices.SortFunc(dt.areas, func(r1, r2 vec.Rect) int {
		return cmp.Compare(r1.Area(), r2.Area())
	})
}

func (dt *SpecialDirtyTracker) calcBestNewRect() (best vec.Rect) {
	minArea := 1000000
	for i := range dt.singles.Count() {
		for j := i + 1; j < dt.singles.Count(); j++ {
			rect := vec.CalcRectContainingCoords(dt.singles.At(i), dt.singles.At(j))
			if area := rect.Area(); area < minArea {
				best = rect
				if area == 2 {
					return
				}
				minArea = area
			}
		}
	}

	return
}

func (dt *SpecialDirtyTracker) IsDirtyAt(pos vec.Coord) bool {
	if dt.singles.Contains(pos) {
		return true
	}

	// areas are sorted by ascending size, so searching the list backwards means searching big ones first.
	for _, area := range slices.Backward(dt.areas) {
		if area.Contains(pos) {
			return true
		}
	}

	return false
}

func (dt *SpecialDirtyTracker) CountDirty() (total int) {
	for _, rect := range dt.areas {
		total += rect.Area()
	}

	total += dt.singles.Count()

	return
}

func (dt *SpecialDirtyTracker) Clean() {
	dt.singles.RemoveAll()
	dt.areas = dt.areas[0:0]
}

func (dt *SpecialDirtyTracker) EachDirtyCoord() iter.Seq[vec.Coord] {
	return func(yield func(vec.Coord) bool) {
		for _, coord := range dt.singles.EachElement() {
			if !yield(coord) {
				return
			}
		}

		for _, area := range dt.areas {
			for coord := range area.EachCoord() {
				if !yield(coord) {
					return
				}
			}
		}
	}
}

func (dt *SpecialDirtyTracker) Dirty() bool {
	return dt.singles.Count() > 0 || len(dt.areas) > 0
}
