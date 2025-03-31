package col

type BlendMode int

const (
	BLEND_MULTIPLY BlendMode = iota
	BLEND_SCREEN
)

// Blends 2 colours c1 and c2. c1 is the active colour (it's on "top").
func Blend(c1, c2 Colour, mode BlendMode) Colour {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	var r, g, b, a int // these need to be ints because the blend calculations can overflow a uint8 during internediate steps

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

	return Make(uint8(a), uint8(r), uint8(g), uint8(b))
}
