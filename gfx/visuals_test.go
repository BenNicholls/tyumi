package gfx

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"testing"
)

func TestNewGlyphVisuals(t *testing.T) {
	v := NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.BLACK})
	if v.Mode != DRAW_GLYPH {
		t.Errorf("Mode = %v, want DRAW_GLYPH", v.Mode)
	}
	if v.Glyph != GLYPH_STAR {
		t.Errorf("Glyph = %v, want GLYPH_STAR", v.Glyph)
	}
}

func TestNewTextVisuals(t *testing.T) {
	v := NewTextVisuals('a', 'b', col.Pair{Fore: col.WHITE, Back: col.BLACK})
	if v.Mode != DRAW_TEXT {
		t.Errorf("Mode = %v, want DRAW_TEXT", v.Mode)
	}
	if v.Chars[0] != 'a' || v.Chars[1] != 'b' {
		t.Errorf("Chars = %v, want ['a' 'b']", v.Chars)
	}
}

func TestVisualsSetGlyph(t *testing.T) {
	var v Visuals
	v.Mode = DRAW_TEXT
	v.SetGlyph(GLYPH_HEART)

	if v.Mode != DRAW_GLYPH {
		t.Errorf("Mode after SetGlyph = %v, want DRAW_GLYPH", v.Mode)
	}
	if v.Glyph != GLYPH_HEART {
		t.Errorf("Glyph = %v, want GLYPH_HEART", v.Glyph)
	}
}

func TestVisualsSetText(t *testing.T) {
	var v Visuals
	v.SetText('x', 'y')

	if v.Mode != DRAW_TEXT {
		t.Errorf("Mode after SetText = %v, want DRAW_TEXT", v.Mode)
	}
	if v.Chars[0] != 'x' || v.Chars[1] != 'y' {
		t.Errorf("Chars = %v, want ['x' 'y']", v.Chars)
	}
}

func TestVisualsReplaceGlyph(t *testing.T) {
	v := NewGlyphVisuals(GLYPH_STAR, col.Pair{})
	got := v.ReplaceGlyph(GLYPH_STAR, GLYPH_HEART)
	if got.Glyph != GLYPH_HEART {
		t.Errorf("ReplaceGlyph matching = %v, want GLYPH_HEART", got.Glyph)
	}

	got2 := v.ReplaceGlyph(GLYPH_HEART, GLYPH_DIAMOND) // GLYPH_HEART != v.Glyph (GLYPH_STAR), so should not match
	if got2.Glyph != GLYPH_STAR {
		t.Errorf("ReplaceGlyph non-matching should leave glyph unchanged, got %v", got2.Glyph)
	}

	// text-mode visuals should be unaffected by ReplaceGlyph
	textVis := NewTextVisuals('a', 'b', col.Pair{})
	got3 := textVis.ReplaceGlyph(GLYPH_NONE, GLYPH_HEART)
	if got3 != textVis {
		t.Errorf("ReplaceGlyph on text-mode visuals should be a no-op, got %v", got3)
	}
}

func TestVisualsReplaceChars(t *testing.T) {
	v := NewTextVisuals(255, 'b', col.Pair{}) // 255 == TEXT_DEFAULT
	got := v.ReplaceChars(TEXT_DEFAULT, [2]uint8{'x', 'y'})

	if got.Chars[0] != 'x' {
		t.Errorf("ReplaceChars left char = %c, want 'x'", got.Chars[0])
	}
	if got.Chars[1] != 'b' {
		t.Errorf("ReplaceChars right char should be unaffected (didn't match) = %c, want 'b'", got.Chars[1])
	}

	// glyph-mode visuals should be unaffected by ReplaceChars
	glyphVis := NewGlyphVisuals(GLYPH_STAR, col.Pair{})
	got2 := glyphVis.ReplaceChars(TEXT_DEFAULT, [2]uint8{'x', 'y'})
	if got2 != glyphVis {
		t.Errorf("ReplaceChars on glyph-mode visuals should be a no-op, got %v", got2)
	}
}

func TestVisualsIsTransparent(t *testing.T) {
	tests := []struct {
		name string
		v    Visuals
		want bool
	}{
		{"draw none", Visuals{Mode: DRAW_NONE}, true},
		{"transparent back", NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.NONE}), true},
		{"opaque", NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.BLACK}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsTransparent(); got != tt.want {
				t.Errorf("IsTransparent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualsHasForegroundContent(t *testing.T) {
	tests := []struct {
		name string
		v    Visuals
		want bool
	}{
		{"empty glyph", NewGlyphVisuals(GLYPH_NONE, col.Pair{Fore: col.WHITE, Back: col.BLACK}), false},
		{"space glyph", NewGlyphVisuals(GLYPH_SPACE, col.Pair{Fore: col.WHITE, Back: col.BLACK}), false},
		{"real glyph", NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.BLACK}), true},
		{"draw none", Visuals{Mode: DRAW_NONE}, false},
		{"fore == back", NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.WHITE}), false},
		{"transparent fore", NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.NONE, Back: col.BLACK}), false},
		{"empty text", NewTextVisuals(0, 0, col.Pair{Fore: col.WHITE, Back: col.BLACK}), false},
		{"real text", NewTextVisuals('a', 0, col.Pair{Fore: col.WHITE, Back: col.BLACK}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.HasForegroundContent(); got != tt.want {
				t.Errorf("HasForegroundContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualsHasBackgroundContent(t *testing.T) {
	if (Visuals{Mode: DRAW_NONE}).HasBackgroundContent() {
		t.Error("DRAW_NONE visuals should not have background content")
	}

	opaque := NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.BLACK})
	if !opaque.HasBackgroundContent() {
		t.Error("opaque back colour should count as background content")
	}

	transparent := NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.NONE})
	if transparent.HasBackgroundContent() {
		t.Error("fully transparent back colour should not count as background content")
	}
}

func TestVisualsGetVisuals(t *testing.T) {
	v := NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.BLACK})
	if got := v.GetVisuals(); got != v {
		t.Errorf("GetVisuals() = %v, want %v", got, v)
	}
}

func TestVisualsString(t *testing.T) {
	glyphVis := NewGlyphVisuals(GLYPH_STAR, col.Pair{Fore: col.WHITE, Back: col.BLACK})
	if got := glyphVis.String(); got == "" {
		t.Error("String() for glyph visuals should not be empty")
	}

	noneVis := Visuals{Mode: DRAW_NONE}
	if got := noneVis.String(); got != "Vis{Mode: None}" {
		t.Errorf("String() for none visuals = %q, want %q", got, "Vis{Mode: None}")
	}
}

func TestGlyphString(t *testing.T) {
	if got := GLYPH_STAR.String(); got != "Star" {
		t.Errorf("GLYPH_STAR.String() = %q, want %q", got, "Star")
	}
}
