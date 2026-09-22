package web

import (
	"image"
	"image/color"

	"github.com/bennicholls/tyumi/gfx/col"
)

// AtlasFromBMP decodes a BMP font and applies Tyumi's font colour key.
// White pixels stay white so they can be tinted. Magenta (255, 0, 255) becomes
// transparent. Any other pixel whose green channel is below 255 becomes white
// with alpha taken from that green channel. This matches the SDL renderer.
func AtlasFromBMP(data []byte) (*image.NRGBA, error) {
	img, err := decodeBMP(data)
	if err != nil {
		return nil, err
	}
	return PrepareAtlas(img), nil
}

// PrepareAtlas applies the font colour key to src and returns an NRGBA image
// whose origin is (0, 0).
func PrepareAtlas(src image.Image) *image.NRGBA {
	b := src.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	key := color.NRGBA{R: 255, G: 0, B: 255, A: 255}
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			c := color.NRGBAModel.Convert(src.At(b.Min.X+x, b.Min.Y+y)).(color.NRGBA)
			if c.G != 0xFF {
				if c == key {
					c = color.NRGBA{}
				} else {
					c.A = c.G
					c.G = 0xFF
				}
			}
			dst.SetNRGBA(x, y, c)
		}
	}
	return dst
}

// TintAtlas multiplies src by colour, the same way SDL colour-mod and alpha-mod
// do. White is returned as src itself.
func TintAtlas(src *image.NRGBA, c col.Colour) *image.NRGBA {
	if c == col.WHITE || src == nil {
		return src
	}
	dst := image.NewNRGBA(src.Bounds())
	cr, cg, cb, ca := c.R(), c.G(), c.B(), c.A()
	for i := 0; i < len(src.Pix); i += 4 {
		dst.Pix[i+0] = uint8(uint16(src.Pix[i+0]) * uint16(cr) / 255)
		dst.Pix[i+1] = uint8(uint16(src.Pix[i+1]) * uint16(cg) / 255)
		dst.Pix[i+2] = uint8(uint16(src.Pix[i+2]) * uint16(cb) / 255)
		dst.Pix[i+3] = uint8(uint16(src.Pix[i+3]) * uint16(ca) / 255)
	}
	return dst
}
