//package for dealing with, forming, deforming, and manipulating colours!
//colours are uint32s in ARGB format.

package col

import "github.com/bennicholls/tyumi/util"

//Make returns a uint32 colour in ARGB formed from provided int components
func Make(a, r, g, b int) (colour uint32) {
	colour = uint32((a % 256) << 24)
	colour |= uint32(r%256) << 16
	colour |= uint32(g%256) << 8
	colour |= uint32(b % 256)

	return
}

//Takes r,g,b ints and creates a colour with alpha 255 in ARGB format.
func MakeOpaque(r, g, b int) uint32 {
	return Make(255, r, g, b)
}

//RGBA returns the RGBA components of an ARGB8888 formatted uint32 colour.
func RGBA(colour uint32) (r, g, b, a uint8) {
	b = uint8(colour & 0x000000FF)
	g = uint8((colour >> 8) & 0x000000FF)
	r = uint8((colour >> 16) & 0x000000FF)
	a = uint8(colour >> 24)

	return
}

//RGB returns the RGB components of an ARGB8888 formatted uint32 colour.
func RGB(colour uint32) (r, g, b uint8) {
	r, g, b, _ = RGBA(colour)
	return
}

//Blends 2 colours c1 and c2. c1 is the active colour (it's on "top").
func Blend(c1, c2 uint32, mode BlendMode) uint32 {
	r1, g1, b1, a1 := RGBA(c1)
	r2, g2, b2, a2 := RGBA(c2)

	var r, g, b, a int

	switch mode {
	case BLEND_MULTIPLY:
		r = int(r1) * int(r2) / 255
		g = int(g1) * int(g2) / 255
		b = int(b1) * int(b2) / 255
		a = int(a1) * int(a2) / 255
	case BLEND_SCREEN:
		r = 255 - int(255-r1)*int(255-r2)/255
		g = 255 - int(255-g1)*int(255-g2)/255
		b = 255 - int(255-b1)*int(255-b2)/255
		a = 255 - int(255-a1)*int(255-a2)/255
	}

	return Make(a, r, g, b)
}

type BlendMode int

const (
	BLEND_MULTIPLY BlendMode = iota
	BLEND_SCREEN
)

//colours! hardcoded for your pleasure.
const (
	NONE      uint32 = 0x00000000
	WHITE     uint32 = 0xFFFFFFFF
	BLACK     uint32 = 0xFF000000
	RED       uint32 = 0xFFFF0000
	BLUE      uint32 = 0xFF0000FF
	LIME      uint32 = 0xFF00FF00
	LIGHTGREY uint32 = 0xFF444444
	GREY      uint32 = 0xFF888888
	DARKGREY  uint32 = 0xFFCCCCCC
	YELLOW    uint32 = 0xFFFFFF00
	FUSCHIA   uint32 = 0xFFFF00FF
	CYAN      uint32 = 0xFF00FFFF
	MAROON    uint32 = 0xFF800000
	OLIVE     uint32 = 0xFF808000
	GREEN     uint32 = 0xFF008000
	TEAL      uint32 = 0xFF008080
	NAVY      uint32 = 0xFF000080
	PURPLE    uint32 = 0xFF800080
)

const KEY = FUSCHIA //key color. renderers should use this to support spritesheet transparency for BMP

type Palette []uint32

//Adds the palette p2 to the end of p.
func (p *Palette) Add(p2 Palette) {
	if (*p)[len(*p)-1] == p2[0] {
		*p = append(*p, p2[1:]...)
	} else {
		*p = append(*p, p2...)
	}
}

//Generate a palette with num items, passing from colour c1 to c2. The colours are
//lineraly interpolated evenly from one to the next. Gradient is NOT circular.
//TODO: Circular palette function?
func GenerateGradient(num int, c1, c2 uint32) (p Palette) {
	p = make(Palette, num)

	r1, g1, b1 := RGB(c1)
	r2, g2, b2 := RGB(c2)

	for i := range p {
		p[i] = MakeOpaque(util.Lerp(int(r1), int(r2), i, len(p)), util.Lerp(int(g1), int(g2), i, len(p)), util.Lerp(int(b1), int(b2), i, len(p)))
	}

	p[num-1] = c2 //fix end of palette rounding lerp stuff.

	return
}

// Lineraly interpolates the colour between colour1 and colour2 over (steps) number of steps, returning the (val)th value.
// NOTE: this completely disregards tranparent colours.
func Lerp(colour1, colour2 uint32, val, steps int) uint32 {
	r1, g1, b1 := RGB(colour1)
	r2, g2, b2 := RGB(colour2)
	return MakeOpaque(util.Lerp(int(r1), int(r2), val, steps), util.Lerp(int(g1), int(g2), val, steps), util.Lerp(int(b1), int(b2), val, steps))
}

// A Pair of colours, fore and back
type Pair struct {
	Fore uint32
	Back uint32
}

// Linearly interpolates between p and p2 over (steps) number of steps, returning the (val)th value.
func (p Pair) Lerp(p2 Pair, val, steps int) Pair {
	return Pair{Lerp(p.Fore, p2.Fore, val, steps), Lerp(p.Back, p2.Back, val, steps)}
}
