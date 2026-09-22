//go:build js && wasm

package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/platform/web"
)

func writeFonts() error {
	const tile = 16
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	glyphs := image.NewNRGBA(image.Rect(0, 0, 16*tile, 16*tile))
	font := image.NewNRGBA(image.Rect(0, 0, 32*(tile/2), 8*tile))
	draw.Draw(glyphs, glyphs.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	draw.Draw(font, font.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	for ch, rows := range font5x7 {
		paintChar(glyphs, font, tile, ch, parseRows(rows), white)
	}
	fillCell(glyphs, int(gfx.GLYPH_BLOCK), tile, 2, white)
	strokeCell(glyphs, int(gfx.GLYPH_BORDER_UUDDLLRR), tile, white)

	if err := writeBMP(glyphFile, glyphs); err != nil {
		return err
	}
	return writeBMP(fontFile, font)
}

func paintChar(glyphs, font *image.NRGBA, tile int, ch byte, rows [7]byte, c color.NRGBA) {
	gx := int(ch%16)*tile + 3
	gy := int(ch/16)*tile + 1
	blit(glyphs, gx, gy, 2, rows, c)

	fx := int(ch%32)*(tile/2) + 1
	fy := int(ch/32)*tile + 4
	blit(font, fx, fy, 1, rows, c)
}

func blit(img *image.NRGBA, ox, oy, scale int, rows [7]byte, c color.NRGBA) {
	for y := 0; y < 7; y++ {
		for x := 0; x < 5; x++ {
			if rows[y]&(1<<uint(4-x)) == 0 {
				continue
			}
			for sy := 0; sy < scale; sy++ {
				for sx := 0; sx < scale; sx++ {
					img.SetNRGBA(ox+x*scale+sx, oy+y*scale+sy, c)
				}
			}
		}
	}
}

func fillCell(img *image.NRGBA, index, tile, inset int, c color.NRGBA) {
	x0 := (index%16)*tile + inset
	y0 := (index/16)*tile + inset
	x1 := (index%16)*tile + tile - inset
	y1 := (index/16)*tile + tile - inset
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.SetNRGBA(x, y, c)
		}
	}
}

func strokeCell(img *image.NRGBA, index, tile int, c color.NRGBA) {
	x0 := (index % 16) * tile
	y0 := (index / 16) * tile
	x1 := x0 + tile - 1
	y1 := y0 + tile - 1
	for x := x0; x <= x1; x++ {
		img.SetNRGBA(x, y0, c)
		img.SetNRGBA(x, y1, c)
	}
	for y := y0; y <= y1; y++ {
		img.SetNRGBA(x0, y, c)
		img.SetNRGBA(x1, y, c)
	}
}

func parseRows(rows [7]string) [7]byte {
	var out [7]byte
	for y, row := range rows {
		for x, ch := range row {
			if x < 5 && ch == '#' {
				out[y] |= 1 << uint(4-x)
			}
		}
	}
	return out
}

func writeBMP(path string, img image.Image) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	data, err := web.EncodeBMP(img)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

var font5x7 = map[byte][7]string{
	'A': {" ### ", "#   #", "#   #", "#####", "#   #", "#   #", "#   #"},
	'B': {"#### ", "#   #", "#   #", "#### ", "#   #", "#   #", "#### "},
	'C': {" ####", "#    ", "#    ", "#    ", "#    ", "#    ", " ####"},
	'D': {"###  ", "#  # ", "#   #", "#   #", "#   #", "#  # ", "###  "},
	'E': {"#####", "#    ", "#    ", "#### ", "#    ", "#    ", "#####"},
	'F': {"#####", "#    ", "#    ", "#### ", "#    ", "#    ", "#    "},
	'G': {" ### ", "#    ", "#    ", "# ###", "#   #", "#   #", " ### "},
	'H': {"#   #", "#   #", "#   #", "#####", "#   #", "#   #", "#   #"},
	'I': {"#####", "  #  ", "  #  ", "  #  ", "  #  ", "  #  ", "#####"},
	'J': {"  ###", "   # ", "   # ", "   # ", "#  # ", "#  # ", " ##  "},
	'K': {"#   #", "#  # ", "# #  ", "##   ", "# #  ", "#  # ", "#   #"},
	'L': {"#    ", "#    ", "#    ", "#    ", "#    ", "#    ", "#####"},
	'M': {"#   #", "## ##", "# # #", "#   #", "#   #", "#   #", "#   #"},
	'N': {"#   #", "##  #", "# # #", "#  ##", "#   #", "#   #", "#   #"},
	'O': {" ### ", "#   #", "#   #", "#   #", "#   #", "#   #", " ### "},
	'P': {"#### ", "#   #", "#   #", "#### ", "#    ", "#    ", "#    "},
	'Q': {" ### ", "#   #", "#   #", "#   #", "# # #", "#  # ", " ## #"},
	'R': {"#### ", "#   #", "#   #", "#### ", "# #  ", "#  # ", "#   #"},
	'S': {" ####", "#    ", "#    ", " ### ", "    #", "    #", "#### "},
	'T': {"#####", "  #  ", "  #  ", "  #  ", "  #  ", "  #  ", "  #  "},
	'U': {"#   #", "#   #", "#   #", "#   #", "#   #", "#   #", " ### "},
	'V': {"#   #", "#   #", "#   #", "#   #", "#   #", " # # ", "  #  "},
	'W': {"#   #", "#   #", "#   #", "# # #", "# # #", "## ##", "#   #"},
	'X': {"#   #", "#   #", " # # ", "  #  ", " # # ", "#   #", "#   #"},
	'Y': {"#   #", "#   #", " # # ", "  #  ", "  #  ", "  #  ", "  #  "},
	'Z': {"#####", "    #", "   # ", "  #  ", " #   ", "#    ", "#####"},
	'0': {" ### ", "#   #", "#  ##", "# # #", "##  #", "#   #", " ### "},
	'1': {"  #  ", " ##  ", "  #  ", "  #  ", "  #  ", "  #  ", " ### "},
	'2': {" ### ", "#   #", "    #", "  ## ", " #   ", "#    ", "#####"},
	'3': {" ### ", "#   #", "    #", " ### ", "    #", "#   #", " ### "},
	'4': {"#   #", "#   #", "#   #", "#####", "    #", "    #", "    #"},
	'5': {"#####", "#    ", "#### ", "    #", "    #", "#   #", " ### "},
	'6': {" ### ", "#    ", "#    ", "#### ", "#   #", "#   #", " ### "},
	'7': {"#####", "    #", "   # ", "  #  ", " #   ", " #   ", " #   "},
	'8': {" ### ", "#   #", "#   #", " ### ", "#   #", "#   #", " ### "},
	'9': {" ### ", "#   #", "#   #", " ####", "    #", "    #", " ### "},
	':': {"     ", "  #  ", "  #  ", "     ", "  #  ", "  #  ", "     "},
}
