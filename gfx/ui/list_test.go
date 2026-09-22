package ui

import (
	"testing"

	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

func newTestList() *List {
	return NewList(vec.Dims{20, 10}, vec.Coord{0, 0}, 0)
}

func TestListInsertAndCount(t *testing.T) {
	l := newTestList()

	if l.Count() != 0 {
		t.Fatalf("Count() = %d, want 0 for new list", l.Count())
	}

	l.InsertText(ALIGN_LEFT, "one", "two", "three")

	if got := l.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}
}

func TestListInsertNoOpForEmptyArgs(t *testing.T) {
	l := newTestList()
	l.Insert() // no items, should be a no-op and not panic

	if l.Count() != 0 {
		t.Errorf("Count() = %d, want 0", l.Count())
	}
}

func TestListInsertSkipsDuplicates(t *testing.T) {
	l := newTestList()
	item := NewTextbox(vec.Dims{10, FIT_TEXT}, vec.ZERO_COORD, 0, "hello", ALIGN_LEFT)

	l.Insert(item)
	l.Insert(item) // duplicate, should be skipped

	if got := l.Count(); got != 1 {
		t.Errorf("Count() = %d, want 1 (duplicate insert should be skipped)", got)
	}
}

func TestListRemove(t *testing.T) {
	l := newTestList()
	item1 := NewTextbox(vec.Dims{10, FIT_TEXT}, vec.ZERO_COORD, 0, "one", ALIGN_LEFT)
	item2 := NewTextbox(vec.Dims{10, FIT_TEXT}, vec.ZERO_COORD, 0, "two", ALIGN_LEFT)
	l.Insert(item1, item2)

	l.Remove(item1)

	if got := l.Count(); got != 1 {
		t.Errorf("Count() = %d, want 1", got)
	}
}

func TestListRemoveAt(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three")

	l.RemoveAt(1)

	if got := l.Count(); got != 2 {
		t.Errorf("Count() = %d, want 2", got)
	}
}

func TestListRemoveAtOutOfRangeIsNoOp(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one")

	l.RemoveAt(5)  // out of range
	l.RemoveAt(-1) // negative

	if got := l.Count(); got != 1 {
		t.Errorf("Count() = %d, want 1 (out-of-range RemoveAt should be a no-op)", got)
	}
}

func TestListRemoveAll(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three")

	l.RemoveAll()

	if got := l.Count(); got != 0 {
		t.Errorf("Count() = %d, want 0", got)
	}
}

func TestListSetCapacityTrimsExcessItems(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three", "four")

	l.SetCapacity(2)

	if got := l.Count(); got != 2 {
		t.Errorf("Count() after SetCapacity(2) = %d, want 2", got)
	}
}

func TestListInsertRespectsCapacity(t *testing.T) {
	l := newTestList()
	l.SetCapacity(2)

	l.InsertText(ALIGN_LEFT, "one", "two", "three")

	if got := l.Count(); got != 2 {
		t.Errorf("Count() = %d, want 2 (capacity should evict oldest items)", got)
	}
}

// Selection is gated behind the unexported selectionEnabled flag, which currently has no exported setter on this
// branch (a List.EnableSelection() method exists on a separate in-progress branch/PR, but not here) - so these
// tests enable it directly via the unexported field, which is fine from inside the package.
func TestListSelectNextPrev(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three")
	l.selectionEnabled = true
	l.Select(0)

	l.SelectNext()
	if got := l.GetSelectionIndex(); got != 1 {
		t.Errorf("GetSelectionIndex() after SelectNext = %d, want 1", got)
	}

	l.SelectNext()
	if got := l.GetSelectionIndex(); got != 2 {
		t.Errorf("GetSelectionIndex() after SelectNext = %d, want 2", got)
	}

	l.SelectNext() // wraps back to 0
	if got := l.GetSelectionIndex(); got != 0 {
		t.Errorf("GetSelectionIndex() after wrapping SelectNext = %d, want 0", got)
	}

	l.SelectPrev() // wraps back to the end
	if got := l.GetSelectionIndex(); got != 2 {
		t.Errorf("GetSelectionIndex() after wrapping SelectPrev = %d, want 2", got)
	}
}

func TestListSelectTopBottom(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three")
	l.selectionEnabled = true

	l.SelectBottom()
	if got := l.GetSelectionIndex(); got != 2 {
		t.Errorf("GetSelectionIndex() after SelectBottom = %d, want 2", got)
	}

	l.SelectTop()
	if got := l.GetSelectionIndex(); got != 0 {
		t.Errorf("GetSelectionIndex() after SelectTop = %d, want 0", got)
	}
}

func TestListSelectClampsToValidRange(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three")
	l.selectionEnabled = true

	l.Select(100)
	if got := l.GetSelectionIndex(); got != 2 {
		t.Errorf("GetSelectionIndex() after Select(100) = %d, want 2 (clamped)", got)
	}

	l.Select(-100)
	if got := l.GetSelectionIndex(); got != 0 {
		t.Errorf("GetSelectionIndex() after Select(-100) = %d, want 0 (clamped)", got)
	}
}

func TestListSelectDisabledIsNoOp(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three")
	// selectionEnabled left false

	l.Select(1)
	if got := l.GetSelectionIndex(); got != -1 {
		t.Errorf("GetSelectionIndex() with selection disabled = %d, want -1 (unchanged default)", got)
	}
}

func TestListSelectOnEmptyList(t *testing.T) {
	l := newTestList()
	l.selectionEnabled = true

	l.Select(0) // no items, should not panic and should leave selection at -1

	if got := l.GetSelectionIndex(); got != -1 {
		t.Errorf("GetSelectionIndex() on empty list = %d, want -1", got)
	}
}

func TestListOnChangeSelectionCallback(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two")
	l.selectionEnabled = true

	calls := 0
	l.OnChangeSelection = func() { calls++ }

	l.Select(1)
	if calls != 1 {
		t.Errorf("OnChangeSelection called %d times after Select, want 1", calls)
	}

	l.Select(1) // same selection, should not trigger callback again
	if calls != 1 {
		t.Errorf("OnChangeSelection called %d times after re-selecting same index, want 1", calls)
	}
}

func TestListOnItemInsertedCallback(t *testing.T) {
	l := newTestList()

	calls := 0
	l.OnItemInserted = func() { calls++ }

	l.InsertText(ALIGN_LEFT, "one", "two", "three")

	if calls != 3 {
		t.Errorf("OnItemInserted called %d times, want 3", calls)
	}
}

func TestListRemoveAtAdjustsSelection(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three")
	l.selectionEnabled = true
	l.Select(2)

	l.RemoveAt(0) // removing an item before/at the selection should shift selection down

	if got := l.GetSelectionIndex(); got != 1 {
		t.Errorf("GetSelectionIndex() after removing item before selection = %d, want 1", got)
	}
}

func TestListHandleActionScrolling(t *testing.T) {
	l := newTestList()
	l.InsertText(ALIGN_LEFT, "one", "two", "three")

	if !l.HandleAction(ACTION_LIST_SCROLLDOWN) {
		t.Error("expected HandleAction(ACTION_LIST_SCROLLDOWN) to report handled")
	}
	if !l.HandleAction(ACTION_LIST_SCROLLUP) {
		t.Error("expected HandleAction(ACTION_LIST_SCROLLUP) to report handled")
	}
}

func TestListHandleActionUnknownReturnsFalse(t *testing.T) {
	l := newTestList()

	if l.HandleAction(input.ActionID(99999)) {
		t.Error("expected HandleAction with an unrecognized action to report not handled")
	}
}
