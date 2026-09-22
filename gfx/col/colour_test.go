package col

import "testing"

func TestMake(t *testing.T) {
	c := Make(0x11, 0x22, 0x33, 0x44)
	if a, r, g, b := c.A(), c.R(), c.G(), c.B(); a != 0x11 || r != 0x22 || g != 0x33 || b != 0x44 {
		t.Errorf("Make(0x11, 0x22, 0x33, 0x44) components = (%x, %x, %x, %x), want (11, 22, 33, 44)", a, r, g, b)
	}
}

func TestMakeOpaque(t *testing.T) {
	c := MakeOpaque(0x22, 0x33, 0x44)
	if c.A() != 0xFF {
		t.Errorf("MakeOpaque alpha = %x, want FF", c.A())
	}
	if r, g, b := c.R(), c.G(), c.B(); r != 0x22 || g != 0x33 || b != 0x44 {
		t.Errorf("MakeOpaque components = (%x, %x, %x), want (22, 33, 44)", r, g, b)
	}
}

func TestColourComponents(t *testing.T) {
	c := RED
	r, g, b, a := c.RGBA()
	if r != 0xFF || g != 0 || b != 0 || a != 0xFF {
		t.Errorf("RGBA() = (%x, %x, %x, %x), want (FF, 0, 0, FF)", r, g, b, a)
	}

	r2, g2, b2 := c.RGB()
	if r2 != 0xFF || g2 != 0 || b2 != 0 {
		t.Errorf("RGB() = (%x, %x, %x), want (FF, 0, 0)", r2, g2, b2)
	}
}

func TestColourIsTransparent(t *testing.T) {
	if !NONE.IsTransparent() {
		t.Error("NONE should be transparent")
	}
	if WHITE.IsTransparent() {
		t.Error("WHITE should not be transparent")
	}

	halfAlpha := Make(0x80, 0xFF, 0xFF, 0xFF)
	if !halfAlpha.IsTransparent() {
		t.Error("half-alpha colour should be considered transparent (any alpha != 0xFF)")
	}
}

func TestColourString(t *testing.T) {
	if got := WHITE.String(); got != "White" {
		t.Errorf("String() for named colour = %q, want %q", got, "White")
	}

	unnamed := Make(0xFF, 0x01, 0x02, 0x03)
	got := unnamed.String()
	want := "(255, 1, 2, 3)"
	if got != want {
		t.Errorf("String() for unnamed colour = %q, want %q", got, want)
	}
}

func TestColourLerp(t *testing.T) {
	c1 := MakeOpaque(0, 0, 0)
	c2 := MakeOpaque(100, 100, 100)

	if got := c1.Lerp(c2, 0, 10); got != c1 {
		t.Errorf("Lerp at val=0 = %v, want %v", got, c1)
	}
	if got := c1.Lerp(c2, 10, 10); got != c2 {
		t.Errorf("Lerp at val=steps = %v, want %v", got, c2)
	}

	mid := c1.Lerp(c2, 5, 10)
	if mid.R() != 50 || mid.G() != 50 || mid.B() != 50 {
		t.Errorf("Lerp midpoint = %v, want R/G/B == 50", mid)
	}
}

func TestColourLerpNone(t *testing.T) {
	c := RED

	if got := NONE.Lerp(c, 5, 10); got != c {
		t.Errorf("Lerp from NONE = %v, want %v (NONE is disregarded)", got, c)
	}
	if got := c.Lerp(NONE, 5, 10); got != c {
		t.Errorf("Lerp to NONE = %v, want %v (NONE is disregarded)", got, c)
	}
}

func TestColourReplace(t *testing.T) {
	if got := RED.Replace(RED, BLUE); got != BLUE {
		t.Errorf("Replace matching colour = %v, want BLUE", got)
	}
	if got := RED.Replace(GREEN, BLUE); got != RED {
		t.Errorf("Replace non-matching colour = %v, want RED unchanged", got)
	}
}

func TestRandom(t *testing.T) {
	c := Random()
	if c.A() != 0xFF {
		t.Errorf("Random() alpha = %x, want FF (opaque)", c.A())
	}
	if c.IsTransparent() {
		t.Error("Random() colour should not be transparent")
	}
}

func TestPairInverted(t *testing.T) {
	p := Pair{Fore: RED, Back: BLUE}
	got := p.Inverted()
	want := Pair{Fore: BLUE, Back: RED}
	if got != want {
		t.Errorf("Inverted() = %v, want %v", got, want)
	}
}

func TestPairLerp(t *testing.T) {
	p1 := Pair{Fore: MakeOpaque(0, 0, 0), Back: MakeOpaque(0, 0, 0)}
	p2 := Pair{Fore: MakeOpaque(100, 100, 100), Back: MakeOpaque(200, 200, 200)}

	got := p1.Lerp(p2, 0, 10)
	if got != p1 {
		t.Errorf("Lerp at val=0 = %v, want %v", got, p1)
	}

	got2 := p1.Lerp(p2, 10, 10)
	if got2 != p2 {
		t.Errorf("Lerp at val=steps = %v, want %v", got2, p2)
	}
}

func TestPairReplace(t *testing.T) {
	p := Pair{Fore: RED, Back: NONE}
	defaults := Pair{Fore: WHITE, Back: BLACK}

	got := p.Replace(NONE, defaults)
	want := Pair{Fore: RED, Back: BLACK}
	if got != want {
		t.Errorf("Replace() = %v, want %v", got, want)
	}
}

func TestPairString(t *testing.T) {
	p := Pair{Fore: WHITE, Back: BLACK}
	got := p.String()
	want := "{White, Black}"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestGenerateGradient(t *testing.T) {
	c1 := MakeOpaque(0, 0, 0)
	c2 := MakeOpaque(100, 100, 100)

	gradient := GenerateGradient(5, c1, c2)

	if len(gradient) != 5 {
		t.Fatalf("GenerateGradient(5, ...) len = %d, want 5", len(gradient))
	}
	if gradient[0] != c1 {
		t.Errorf("gradient[0] = %v, want %v", gradient[0], c1)
	}
	if gradient[4] != c2 {
		t.Errorf("gradient[4] = %v, want %v", gradient[4], c2)
	}
}
