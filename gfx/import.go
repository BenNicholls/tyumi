package gfx

import (
	"github.com/bennicholls/reximage"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

// imports image data from an XP file at the provided, drawing it at depth zero on the returned canvas
func ImportXPData(path string) (image Canvas) {
	imageData, err := reximage.Import(path)
	if err != nil {
		log.Error(err)
		return
	}

	image.Init(vec.Dims{imageData.Width, imageData.Height})
	for cursor := range vec.EachCoordInArea(image) {
		cellData, err := imageData.GetCell(cursor.X, cursor.Y)
		if err != nil {
			log.Debug(err)
			return Canvas{}
		}

		if cellData.Undrawn() {
			image.DrawNone(cursor)
		} else {
			fore, back := cellData.ARGB()
			image.DrawVisuals(cursor, 0, NewGlyphVisuals(Glyph(cellData.Glyph), col.Pair{col.Colour(fore), col.Colour(back)}))
		}
	}

	return
}
