package web

import (
	"encoding/binary"
	"image"
	"image/color"
	"testing"

	"github.com/bennicholls/tyumi/gfx/col"
)

func TestPrepareAtlas(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 4, 1))
	src.SetNRGBA(0, 0, color.NRGBA{255, 255, 255, 255})
	src.SetNRGBA(1, 0, color.NRGBA{255, 0, 255, 255})
	src.SetNRGBA(2, 0, color.NRGBA{0, 0, 0, 255})
	src.SetNRGBA(3, 0, color.NRGBA{10, 128, 20, 255})

	got := PrepareAtlas(src)

	if c := got.NRGBAAt(0, 0); c != (color.NRGBA{255, 255, 255, 255}) {
		t.Errorf("white pixel = %v", c)
	}
	if c := got.NRGBAAt(1, 0); c.A != 0 {
		t.Errorf("magenta pixel alpha = %d, want 0", c.A)
	}
	if c := got.NRGBAAt(2, 0); c.A != 0 {
		t.Errorf("black pixel alpha = %d, want 0", c.A)
	}
	if c := got.NRGBAAt(3, 0); c != (color.NRGBA{10, 255, 20, 128}) {
		t.Errorf("partial pixel = %v, want (10, 255, 20, 128)", c)
	}
}

func TestTintAtlas(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	src.SetNRGBA(0, 0, color.NRGBA{255, 255, 255, 255})

	if TintAtlas(src, col.WHITE) != src {
		t.Fatal("white tint should reuse the source image")
	}

	got := TintAtlas(src, col.MakeOpaque(255, 0, 0))
	if c := got.NRGBAAt(0, 0); c != (color.NRGBA{255, 0, 0, 255}) {
		t.Fatalf("red tint = %v", c)
	}
}

func TestAtlasFromBMP(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	src.SetNRGBA(0, 0, color.NRGBA{255, 255, 255, 255})
	src.SetNRGBA(1, 0, color.NRGBA{255, 0, 255, 255})

	encoded, err := EncodeBMP(src)
	if err != nil {
		t.Fatal(err)
	}
	got, err := AtlasFromBMP(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if c := got.NRGBAAt(0, 0); c != (color.NRGBA{255, 255, 255, 255}) {
		t.Errorf("white = %v", c)
	}
	if c := got.NRGBAAt(1, 0); c.A != 0 {
		t.Errorf("magenta alpha = %d", c.A)
	}
}

func TestAtlasFromPalettedBMP(t *testing.T) {
	// 2x1, 8-bit, two palette entries: black and white.
	buf := make([]byte, 62+4)
	buf[0], buf[1] = 'B', 'M'
	binary.LittleEndian.PutUint32(buf[2:], uint32(len(buf)))
	binary.LittleEndian.PutUint32(buf[10:], 62)
	binary.LittleEndian.PutUint32(buf[14:], 40)
	binary.LittleEndian.PutUint32(buf[18:], 2)
	binary.LittleEndian.PutUint32(buf[22:], 1)
	binary.LittleEndian.PutUint16(buf[26:], 1)
	binary.LittleEndian.PutUint16(buf[28:], 8)
	binary.LittleEndian.PutUint32(buf[46:], 2)
	// palette is BGRA
	buf[54], buf[55], buf[56], buf[57] = 0, 0, 0, 0
	buf[58], buf[59], buf[60], buf[61] = 255, 255, 255, 0
	buf[62], buf[63] = 0, 1

	got, err := AtlasFromBMP(buf)
	if err != nil {
		t.Fatal(err)
	}
	if c := got.NRGBAAt(0, 0); c.A != 0 {
		t.Errorf("black alpha = %d", c.A)
	}
	if c := got.NRGBAAt(1, 0); c != (color.NRGBA{255, 255, 255, 255}) {
		t.Errorf("white = %v", c)
	}
}
