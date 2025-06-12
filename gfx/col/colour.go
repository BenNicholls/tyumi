// package for dealing with, forming, deforming, and manipulating colours!
// colours are uint32s in ARGB format.
package col

import (
	"fmt"
	"math/rand"

	"github.com/bennicholls/tyumi/util"
)

// Colour is an ARGB8888 encoded colour.
type Colour uint32

func (c Colour) String() string {
	if name, ok := ColourNames[c]; ok {
		return name
	}

	return fmt.Sprintf("(%v, %v, %v, %v)", c.A(), c.R(), c.G(), c.B())
}

// Make returns an ARGB8888 colour formed from provided uint8 components.
func Make(a, r, g, b uint8) (colour Colour) {
	colour = Colour(a) << 24
	colour |= Colour(r) << 16
	colour |= Colour(g) << 8
	colour |= Colour(b)

	return
}

// Make returns an ARGB8888 colour formed from provided uint8 components, with alpha set to 255 (0xFF)
func MakeOpaque(r, g, b uint8) Colour {
	return Make(255, r, g, b)
}

// Returns the Alpha component of a colour.
func (c Colour) A() uint8 {
	return uint8(c >> 24)
}

// Returns the Red component of a colour.
func (c Colour) R() uint8 {
	return uint8((c >> 16) & 0xFF)
}

// Returns the Green component of a colour.
func (c Colour) G() uint8 {
	return uint8((c >> 8) & 0xFF)
}

// Returns the Blue component of a colour.
func (c Colour) B() uint8 {
	return uint8(c & 0xFF)
}

// RGBA returns the RGBA components of an ARGB8888 formatted colour.
func (c Colour) RGBA() (r, g, b, a uint8) {
	return c.R(), c.G(), c.B(), c.A()
}

// RGB returns the RGB components of an ARGB8888 formatted colour.
func (c Colour) RGB() (r, g, b uint8) {
	return c.R(), c.G(), c.B()
}

func (c Colour) IsTransparent() bool {
	return c.A() != 0xFF
}

// Lineraly interpolates the colour between c and c2 over (steps) number of steps, returning the (val)th value.
// NOTE: this completely disregards transparent colours, except for NONE. If lerping to NONE, it just doesn't do it and
// returns the other colour.
func (c Colour) Lerp(c2 Colour, val, steps int) Colour {
	if c == NONE || val >= steps {
		return c2
	} else if c2 == NONE || val <= 0 {
		return c
	}

	return MakeOpaque(
		util.Lerp(c.R(), c2.R(), val, steps),
		util.Lerp(c.G(), c2.G(), val, steps),
		util.Lerp(c.B(), c2.B(), val, steps),
	)
}

// Random returns a random opaque colour.
func Random() Colour {
	return MakeOpaque(uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32()))
}

// A Pair of colours, fore and back
type Pair struct {
	Fore, Back Colour
}

func (p Pair) String() string {
	return fmt.Sprintf("{%s, %s}", p.Fore, p.Back)
}

func (p Pair) Inverted() Pair {
	return Pair{p.Back, p.Fore}
}

// Linearly interpolates between p and p2 over (steps) number of steps, returning the (val)th value.
func (p Pair) Lerp(p2 Pair, val, steps int) Pair {
	return Pair{p.Fore.Lerp(p2.Fore, val, steps), p.Back.Lerp(p2.Back, val, steps)}
}

type Gradient []Colour

// Generate a gredient with num items, passing from colour c1 to c2. The colours are
// lineraly interpolated evenly from one to the next. Gradient is NOT circular.
// TODO: Circular gradient function?
func GenerateGradient(num int, c1, c2 Colour) (p Gradient) {
	p = make(Gradient, num)
	for i := range num {
		p[i] = c1.Lerp(c2, i, num-1)
	}

	return
}
