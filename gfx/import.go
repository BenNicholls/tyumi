package gfx

import (
	"github.com/bennicholls/reximage"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

// imports image data from an XP file at the provided, drawing it at depth zero on the returned canvas
func ImportXPData(path string) (c Canvas) {
	imageData, err := reximage.Import(path)
	if err != nil {
		log.Error(err)
		return
	}

	c.Init(vec.Dims{imageData.Width, imageData.Height})
	for cursor := range vec.EachCoordInArea(c) {
		cell_data, err := imageData.GetCell(cursor.X, cursor.Y)
		if err != nil {
			log.Debug(err)
			return
		}

		if cell_data.Undrawn() {
			c.DrawNone(cursor)
		} else {
			fore, back := cell_data.ARGB()
			c.DrawVisuals(cursor, 0, NewGlyphVisuals(Glyph(cell_data.Glyph), col.Pair{fore, back}))
		}
	}

	return c
}
