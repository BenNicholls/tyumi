package gfx

import (
	"testing"

	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

func TestCanvasInit(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{10, 10})

	if !c.Ready() {
		t.Error("expected canvas to be Ready() after Init")
	}
	if got := c.Size(); got != (vec.Dims{10, 10}) {
		t.Errorf("Size() = %v, want {10 10}", got)
	}
}

func TestCanvasInitZeroArea(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{0, 0})

	if c.Ready() {
		t.Error("canvas initialized with zero area should not be Ready()")
	}
}

func TestCanvasResizeIsNoOpForSameSize(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})
	c.DrawGlyph(vec.Coord{1, 1}, 0, GLYPH_STAR)
	c.Clean()

	c.Resize(vec.Dims{5, 5}) // same size, should be a no-op and not clear existing content

	if got := c.GetCell(vec.Coord{1, 1}).Glyph; got != GLYPH_STAR {
		t.Errorf("Resize to same size cleared canvas, glyph = %v, want GLYPH_STAR", got)
	}
}

func TestCanvasBoundsAndInBounds(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	if got := c.Bounds(); got != (vec.Rect{vec.Coord{0, 0}, vec.Dims{5, 5}}) {
		t.Errorf("Bounds() = %v, want {(0,0), (5,5)}", got)
	}

	if !c.InBounds(vec.Coord{2, 2}) {
		t.Error("expected (2,2) to be in bounds")
	}
	if c.InBounds(vec.Coord{5, 5}) {
		t.Error("expected (5,5) to be out of bounds")
	}
	if c.InBounds(vec.Coord{-1, 0}) {
		t.Error("expected (-1,0) to be out of bounds")
	}
}

func TestCanvasGetCellOutOfBounds(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	got := c.GetCell(vec.Coord{100, 100})
	want := Visuals{}
	if got != want {
		t.Errorf("GetCell out of bounds = %v, want zero value %v", got, want)
	}
}

func TestCanvasDrawGlyphAndGetCell(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	c.DrawGlyph(vec.Coord{2, 2}, 0, GLYPH_STAR)

	got := c.GetCell(vec.Coord{2, 2})
	if got.Glyph != GLYPH_STAR {
		t.Errorf("GetCell glyph = %v, want GLYPH_STAR", got.Glyph)
	}
	if got.Mode != DRAW_GLYPH {
		t.Errorf("GetCell mode = %v, want DRAW_GLYPH", got.Mode)
	}
}

func TestCanvasDrawGlyphOutOfBoundsIsNoOp(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	c.DrawGlyph(vec.Coord{100, 100}, 0, GLYPH_STAR) // should not panic
}

func TestCanvasDepthRespect(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	c.DrawGlyph(vec.Coord{2, 2}, 5, GLYPH_STAR)
	c.DrawGlyph(vec.Coord{2, 2}, 2, GLYPH_HEART) // lower depth, should not overwrite

	got := c.GetCell(vec.Coord{2, 2})
	if got.Glyph != GLYPH_STAR {
		t.Errorf("lower-depth draw overwrote higher-depth cell, glyph = %v, want GLYPH_STAR", got.Glyph)
	}

	c.DrawGlyph(vec.Coord{2, 2}, 10, GLYPH_HEART) // higher depth, should overwrite
	got2 := c.GetCell(vec.Coord{2, 2})
	if got2.Glyph != GLYPH_HEART {
		t.Errorf("higher-depth draw did not overwrite cell, glyph = %v, want GLYPH_HEART", got2.Glyph)
	}
}

func TestCanvasDrawColours(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	colours := col.Pair{Fore: col.RED, Back: col.BLUE}
	c.DrawColours(vec.Coord{1, 1}, 0, colours)

	got := c.GetCell(vec.Coord{1, 1})
	if got.Colours != colours {
		t.Errorf("GetCell colours = %v, want %v", got.Colours, colours)
	}
}

func TestCanvasDrawNone(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	c.DrawGlyph(vec.Coord{1, 1}, 0, GLYPH_STAR)
	c.DrawNone(vec.Coord{1, 1})

	got := c.GetCell(vec.Coord{1, 1})
	if got.Mode != DRAW_NONE {
		t.Errorf("GetCell mode after DrawNone = %v, want DRAW_NONE", got.Mode)
	}
}

func TestCanvasClear(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	c.DrawGlyph(vec.Coord{1, 1}, 0, GLYPH_STAR)
	c.Clear()

	got := c.GetCell(vec.Coord{1, 1})
	if got.Glyph == GLYPH_STAR {
		t.Error("expected Clear() to reset cell content")
	}
}

func TestCanvasClearArea(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{10, 10})

	c.DrawGlyph(vec.Coord{1, 1}, 0, GLYPH_STAR)
	c.DrawGlyph(vec.Coord{8, 8}, 0, GLYPH_STAR)

	c.Clear(vec.Rect{vec.Coord{0, 0}, vec.Dims{3, 3}})

	if got := c.GetCell(vec.Coord{1, 1}).Glyph; got == GLYPH_STAR {
		t.Error("expected cell inside cleared area to be reset")
	}
	if got := c.GetCell(vec.Coord{8, 8}).Glyph; got != GLYPH_STAR {
		t.Error("expected cell outside cleared area to be preserved")
	}
}

func TestCanvasClearAtDepth(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	c.DrawGlyph(vec.Coord{1, 1}, 5, GLYPH_STAR)
	c.ClearAtDepth(2) // only clear cells at depth <= 2, this one is at 5

	if got := c.GetCell(vec.Coord{1, 1}).Glyph; got != GLYPH_STAR {
		t.Error("expected cell above clear depth to be preserved")
	}

	c.ClearAtDepth(10) // now clear at a depth above the cell's
	if got := c.GetCell(vec.Coord{1, 1}).Glyph; got == GLYPH_STAR {
		t.Error("expected cell at or below clear depth to be reset")
	}
}

func TestCanvasFlattenTo(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	c.DrawGlyph(vec.Coord{1, 1}, 10, GLYPH_STAR)
	c.FlattenTo(3)

	if got := c.GetDepth(vec.Coord{1, 1}); got != 3 {
		t.Errorf("GetDepth after FlattenTo(3) = %d, want 3", got)
	}

	// a subsequent draw at depth 5 should now succeed, since the flattened depth (3) is lower
	c.DrawGlyph(vec.Coord{1, 1}, 5, GLYPH_HEART)
	if got := c.GetCell(vec.Coord{1, 1}).Glyph; got != GLYPH_HEART {
		t.Error("expected draw at depth 5 to succeed after flattening cell to depth 3")
	}
}

func TestCanvasDirtyTracking(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})
	c.Clean() // Init leaves everything dirty (from the initial Clear), start from a clean slate

	if c.IsDirtyAt(vec.Coord{1, 1}) {
		t.Error("expected cell to be clean after Clean()")
	}

	c.DrawGlyph(vec.Coord{1, 1}, 0, GLYPH_STAR)
	if !c.IsDirtyAt(vec.Coord{1, 1}) {
		t.Error("expected cell to be dirty after drawing to it")
	}
}

func TestCanvasIsTransparent(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})
	c.SetDefaultVisuals(Visuals{Mode: DRAW_GLYPH, Colours: col.Pair{Fore: col.WHITE, Back: col.BLACK}})

	if c.IsTransparent() {
		t.Error("canvas with fully opaque default visuals should not be transparent")
	}

	// NOTE: col.NONE can't be used here to make a cell transparent - setCell treats a Back colour of exactly
	// col.NONE as a sentinel meaning "keep the existing back colour" (see Canvas.setCell's use of
	// Pair.Replace(col.NONE, ...)), so we use an arbitrary non-opaque (but non-NONE) colour instead.
	halfAlphaBack := col.Make(0x80, 0, 0, 0)
	c.DrawColours(vec.Coord{1, 1}, 0, col.Pair{Fore: col.WHITE, Back: halfAlphaBack})
	if !c.IsTransparent() {
		t.Error("canvas with a transparent cell should be reported as transparent")
	}
}

func TestCanvasDefaultColoursAndVisuals(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	colours := col.Pair{Fore: col.RED, Back: col.GREEN}
	c.SetDefaultColours(colours)

	if got := c.DefaultColours(); got != colours {
		t.Errorf("DefaultColours() = %v, want %v", got, colours)
	}
	if got := c.DefaultVisuals().Colours; got != colours {
		t.Errorf("DefaultVisuals().Colours = %v, want %v", got, colours)
	}
}

func TestCanvasCopyArea(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{10, 10})
	c.DrawGlyph(vec.Coord{5, 5}, 0, GLYPH_STAR)

	copyArea := vec.Rect{vec.Coord{4, 4}, vec.Dims{3, 3}}
	copied := c.CopyArea(copyArea)

	if got := copied.Size(); got != (vec.Dims{3, 3}) {
		t.Errorf("CopyArea size = %v, want {3 3}", got)
	}

	// (5,5) in the original maps to (1,1) in the copy (offset by the area's top-left corner)
	if got := copied.GetCell(vec.Coord{1, 1}).Glyph; got != GLYPH_STAR {
		t.Errorf("copied cell glyph = %v, want GLYPH_STAR", got)
	}
}

func TestCanvasCopyAreaNoIntersection(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	copied := c.CopyArea(vec.Rect{vec.Coord{100, 100}, vec.Dims{3, 3}})
	// CopyArea always Init()s the result to the requested size regardless of intersection, so it's still "ready" -
	// it just has no content copied into it (every cell retains the plain default visuals).
	if !copied.Ready() {
		t.Error("CopyArea should still initialize the result canvas even with no intersection")
	}
	if got := copied.Size(); got != (vec.Dims{3, 3}) {
		t.Errorf("CopyArea size = %v, want {3 3}", got)
	}
}

func TestCanvasEachCell(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{2, 2})
	c.DrawGlyph(vec.Coord{1, 1}, 0, GLYPH_STAR)

	found := false
	count := 0
	for cell, pos := range c.EachCell() {
		count++
		if pos == (vec.Coord{1, 1}) && cell.Glyph == GLYPH_STAR {
			found = true
		}
	}

	if count != 4 {
		t.Errorf("EachCell() visited %d cells, want 4", count)
	}
	if !found {
		t.Error("EachCell() did not find the drawn cell with correct position")
	}
}

func TestCanvasEachCellWithOffset(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{3, 3})
	c.SetOrigin(vec.Coord{1, 1}) // sets offset so that (1,1) is now local (0,0)

	var positions []vec.Coord
	for _, pos := range c.EachCell() {
		positions = append(positions, pos)
	}

	// with origin at (1,1), the canvas spans from (-1,-1) to (1,1)
	found := false
	for _, p := range positions {
		if p == (vec.Coord{-1, -1}) {
			found = true
		}
	}
	if !found {
		t.Errorf("expected EachCell to report offset coordinates, got %v", positions)
	}
}

func TestCanvasSetOriginOutOfBoundsIsNoOp(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	c.SetOrigin(vec.Coord{100, 100}) // out of bounds, should be a no-op

	if got := c.Bounds().Coord; got != (vec.Coord{0, 0}) {
		t.Errorf("SetOrigin out of bounds changed offset, Bounds().Coord = %v, want {0 0}", got)
	}
}

func TestCanvasGetDepthPanicsOutOfBounds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected GetDepth to panic for an out-of-bounds coordinate")
		}
	}()

	var c Canvas
	c.Init(vec.Dims{5, 5})
	c.GetDepth(vec.Coord{100, 100})
}

func TestCanvasString(t *testing.T) {
	var c Canvas
	c.Init(vec.Dims{5, 5})

	if got := c.String(); got == "" {
		t.Error("String() should not be empty")
	}
}
