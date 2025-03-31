package gfx

import (
	"github.com/bennicholls/reximage"
	"github.com/bennicholls/tyumi/log"
)

// ExportToXP writes the contents of the canvas to an .xp file at the given path.
func (c Canvas) ExportToXP(path string) {
	var image reximage.ImageData
	image.Init(c.size.W, c.size.H)

	for i, cell := range c.cells {
		cellData := reximage.CellData{}
		cellData.SetColoursARGB(uint32(cell.Colours.Fore), uint32(cell.Colours.Back))
		switch cell.Mode {
		case DRAW_GLYPH:
			cellData.Glyph = uint32(cell.Glyph)
		case DRAW_TEXT:
			cellData.Glyph = uint32(cell.Chars[0])
		case DRAW_NONE:
			cellData.Clear()
		}
		err := image.SetCell(i%c.size.W, i/c.size.W, cellData)
		if err != nil {
			log.Error("Error setting cell during export: ", err)
		}
	}

	err := reximage.Export(image, path)
	if err != nil {
		log.Error("Could not export image: ", err)
	}
}
