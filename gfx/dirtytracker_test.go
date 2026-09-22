package gfx

import (
	"testing"

	"github.com/bennicholls/tyumi/vec"
)

func TestDirtyTrackerInitStartsClean(t *testing.T) {
	var dt DirtyTracker
	dt.Init(vec.Dims{5, 5})

	if dt.Dirty() {
		t.Error("freshly initialized DirtyTracker should not be dirty")
	}
	if dt.IsDirtyAt(vec.Coord{2, 2}) {
		t.Error("freshly initialized cell should not be dirty")
	}
}

func TestDirtyTrackerSetDirty(t *testing.T) {
	var dt DirtyTracker
	dt.Init(vec.Dims{5, 5})

	dt.SetDirty(vec.Coord{2, 2})

	if !dt.Dirty() {
		t.Error("expected tracker to report dirty after SetDirty")
	}
	if !dt.IsDirtyAt(vec.Coord{2, 2}) {
		t.Error("expected (2,2) to be dirty")
	}
	if dt.IsDirtyAt(vec.Coord{0, 0}) {
		t.Error("expected (0,0) to remain clean")
	}
	if got := dt.CountDirty(); got != 1 {
		t.Errorf("CountDirty() = %d, want 1", got)
	}
}

func TestDirtyTrackerSetAllDirty(t *testing.T) {
	var dt DirtyTracker
	dt.Init(vec.Dims{5, 5})

	dt.SetAllDirty()

	if !dt.IsDirtyAt(vec.Coord{4, 4}) {
		t.Error("expected all cells to be dirty after SetAllDirty")
	}
}

func TestDirtyTrackerClean(t *testing.T) {
	var dt DirtyTracker
	dt.Init(vec.Dims{5, 5})

	dt.SetDirty(vec.Coord{2, 2})
	dt.Clean()

	if dt.Dirty() {
		t.Error("expected tracker to be clean after Clean()")
	}
	if dt.IsDirtyAt(vec.Coord{2, 2}) {
		t.Error("expected (2,2) to be clean after Clean()")
	}

	dt.SetAllDirty()
	dt.Clean()
	if dt.Dirty() {
		t.Error("expected Clean() to also clear the allDirty flag")
	}
}

func TestSpecialDirtyTrackerBasics(t *testing.T) {
	var dt SpecialDirtyTracker
	dt.Init(10, 5)

	if dt.Dirty() {
		t.Error("freshly initialized SpecialDirtyTracker should not be dirty")
	}

	dt.SetDirty(vec.Coord{1, 1})
	dt.SetDirty(vec.Coord{5, 5})

	if !dt.Dirty() {
		t.Error("expected tracker to be dirty after SetDirty")
	}
	if !dt.IsDirtyAt(vec.Coord{1, 1}) || !dt.IsDirtyAt(vec.Coord{5, 5}) {
		t.Error("expected both dirtied coords to report dirty")
	}
	if dt.IsDirtyAt(vec.Coord{9, 9}) {
		t.Error("expected untouched coord to report clean")
	}
	if got := dt.CountDirty(); got != 2 {
		t.Errorf("CountDirty() = %d, want 2", got)
	}
}

func TestSpecialDirtyTrackerClean(t *testing.T) {
	var dt SpecialDirtyTracker
	dt.Init(10, 5)

	dt.SetDirty(vec.Coord{1, 1})
	dt.Clean()

	if dt.Dirty() {
		t.Error("expected tracker to be clean after Clean()")
	}
	if dt.IsDirtyAt(vec.Coord{1, 1}) {
		t.Error("expected coord to be clean after Clean()")
	}
}

func TestSpecialDirtyTrackerEachDirtyCoord(t *testing.T) {
	var dt SpecialDirtyTracker
	dt.Init(10, 5)

	dt.SetDirty(vec.Coord{1, 1})
	dt.SetDirty(vec.Coord{2, 2})

	found := make(map[vec.Coord]bool)
	for c := range dt.EachDirtyCoord() {
		found[c] = true
	}

	if !found[vec.Coord{1, 1}] || !found[vec.Coord{2, 2}] {
		t.Errorf("EachDirtyCoord() did not report both dirtied coords, got %v", found)
	}
}

func TestSpecialDirtyTrackerMergesIntoAreas(t *testing.T) {
	// with a low MaxSingles, adding enough dirty coords should force it to merge some into an area
	var dt SpecialDirtyTracker
	dt.Init(2, 5)

	dt.SetDirty(vec.Coord{0, 0})
	dt.SetDirty(vec.Coord{1, 0})
	dt.SetDirty(vec.Coord{2, 0}) // exceeds MaxSingles, should trigger a merge

	// regardless of internal representation (single vs merged area), all 3 coords should still report dirty
	for _, c := range []vec.Coord{{0, 0}, {1, 0}, {2, 0}} {
		if !dt.IsDirtyAt(c) {
			t.Errorf("expected %v to be dirty after merge", c)
		}
	}
}
