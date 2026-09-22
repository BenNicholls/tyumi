package web

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/bits"
)

// EncodeBMP writes img as a 24-bit uncompressed BMP.
func EncodeBMP(img image.Image) ([]byte, error) {
	if img == nil {
		return nil, errors.New("nil image")
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return nil, errors.New("image has no area")
	}
	stride := (w*3 + 3) &^ 3
	size := 54 + stride*h
	buf := make([]byte, size)
	buf[0], buf[1] = 'B', 'M'
	binary.LittleEndian.PutUint32(buf[2:], uint32(size))
	binary.LittleEndian.PutUint32(buf[10:], 54)
	binary.LittleEndian.PutUint32(buf[14:], 40)
	binary.LittleEndian.PutUint32(buf[18:], uint32(w))
	binary.LittleEndian.PutUint32(buf[22:], uint32(h))
	binary.LittleEndian.PutUint16(buf[26:], 1)
	binary.LittleEndian.PutUint16(buf[28:], 24)

	for y := 0; y < h; y++ {
		row := buf[54+(h-1-y)*stride:]
		for x := 0; x < w; x++ {
			c := color.NRGBAModel.Convert(img.At(b.Min.X+x, b.Min.Y+y)).(color.NRGBA)
			row[x*3+0] = c.B
			row[x*3+1] = c.G
			row[x*3+2] = c.R
		}
	}
	return buf, nil
}

// decodeBMP reads the BMP files Tyumi fonts use: uncompressed 1, 4, 8, 24, and
// 32-bit images, including 32-bit bitfield masks.
func decodeBMP(data []byte) (*image.NRGBA, error) {
	if len(data) < 26 || data[0] != 'B' || data[1] != 'M' {
		return nil, errors.New("not a BMP")
	}
	pixelOff := int(binary.LittleEndian.Uint32(data[10:14]))
	dib := int(binary.LittleEndian.Uint32(data[14:18]))
	if dib != 12 && dib < 40 {
		return nil, fmt.Errorf("unsupported BMP header size %d", dib)
	}
	if len(data) < 14+dib {
		return nil, errors.New("truncated BMP header")
	}

	var width, height, bpp int
	var compression uint32
	var clrUsed int
	topDown := false
	if dib == 12 {
		width = int(binary.LittleEndian.Uint16(data[18:20]))
		height = int(binary.LittleEndian.Uint16(data[20:22]))
		bpp = int(binary.LittleEndian.Uint16(data[24:26]))
	} else {
		width = int(int32(binary.LittleEndian.Uint32(data[18:22])))
		h := int32(binary.LittleEndian.Uint32(data[22:26]))
		if h < 0 {
			topDown = true
			h = -h
		}
		height = int(h)
		bpp = int(binary.LittleEndian.Uint16(data[28:30]))
		compression = binary.LittleEndian.Uint32(data[30:34])
		clrUsed = int(binary.LittleEndian.Uint32(data[46:50]))
	}
	if width <= 0 || height <= 0 {
		return nil, errors.New("BMP has no area")
	}
	if compression != 0 && !(compression == 3 && bpp == 32) {
		return nil, fmt.Errorf("unsupported BMP compression %d", compression)
	}

	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	stride := ((width*bpp + 31) / 32) * 4
	if pixelOff < 0 || pixelOff > len(data) {
		return nil, errors.New("truncated BMP pixels")
	}

	var palette []color.NRGBA
	if bpp <= 8 {
		n := clrUsed
		if n == 0 {
			n = 1 << bpp
		}
		entry := 4
		palAt := 14 + dib
		if dib == 12 {
			entry = 3
		}
		if palAt+n*entry > len(data) {
			return nil, errors.New("truncated BMP palette")
		}
		palette = make([]color.NRGBA, n)
		for i := 0; i < n; i++ {
			p := data[palAt+i*entry:]
			c := color.NRGBA{B: p[0], G: p[1], R: p[2], A: 255}
			palette[i] = c
		}
	}

	rShift, rBits, gShift, gBits, bShift, bBits, aShift, aBits := 16, 8, 8, 8, 0, 8, 24, 8
	if compression == 3 {
		masks := data[54:]
		if dib > 40 {
			masks = data[14+40:]
		}
		if len(masks) < 12 {
			return nil, errors.New("truncated BMP bitfield masks")
		}
		rMask := binary.LittleEndian.Uint32(masks[0:4])
		gMask := binary.LittleEndian.Uint32(masks[4:8])
		bMask := binary.LittleEndian.Uint32(masks[8:12])
		aMask := uint32(0)
		if len(masks) >= 16 && dib >= 56 {
			aMask = binary.LittleEndian.Uint32(masks[12:16])
		}
		rShift, rBits = maskParts(rMask)
		gShift, gBits = maskParts(gMask)
		bShift, bBits = maskParts(bMask)
		aShift, aBits = maskParts(aMask)
	}

	for y := 0; y < height; y++ {
		srcY := y
		if !topDown {
			srcY = height - 1 - y
		}
		rowAt := pixelOff + srcY*stride
		need := (width*bpp + 7) / 8
		if rowAt < 0 || rowAt+need > len(data) {
			return nil, errors.New("truncated BMP pixels")
		}
		row := data[rowAt:]
		for x := 0; x < width; x++ {
			var c color.NRGBA
			switch bpp {
			case 1:
				c = paletteAt(palette, int((row[x/8]>>(7-(x%8)))&1))
			case 4:
				v := row[x/2]
				if x%2 == 0 {
					v >>= 4
				}
				c = paletteAt(palette, int(v&0x0F))
			case 8:
				c = paletteAt(palette, int(row[x]))
			case 24:
				c = color.NRGBA{R: row[x*3+2], G: row[x*3+1], B: row[x*3], A: 255}
			case 32:
				pix := binary.LittleEndian.Uint32(row[x*4:])
				c = color.NRGBA{
					R: expandMask(pix, rShift, rBits),
					G: expandMask(pix, gShift, gBits),
					B: expandMask(pix, bShift, bBits),
					A: 255,
				}
				if compression == 3 && aBits > 0 {
					c.A = expandMask(pix, aShift, aBits)
				} else if compression == 0 {
					c.A = uint8(pix >> 24)
					if c.A == 0 && (pix&0x00FFFFFF) != 0 {
						c.A = 255
					}
				}
			default:
				return nil, fmt.Errorf("unsupported BMP bit depth %d", bpp)
			}
			dst.SetNRGBA(x, y, c)
		}
	}
	return dst, nil
}

func paletteAt(palette []color.NRGBA, i int) color.NRGBA {
	if i < 0 || i >= len(palette) {
		return color.NRGBA{A: 255}
	}
	return palette[i]
}

func maskParts(mask uint32) (shift, n int) {
	if mask == 0 {
		return 0, 0
	}
	return bits.TrailingZeros32(mask), bits.OnesCount32(mask)
}

func expandMask(v uint32, shift, n int) uint8 {
	if n <= 0 {
		return 0
	}
	val := (v >> shift) & ((1 << n) - 1)
	if n >= 8 {
		return uint8(val >> (n - 8))
	}
	return uint8(val * 255 / ((1 << n) - 1))
}
