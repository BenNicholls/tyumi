package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/vec"
)

type Image struct {
	ElementPrototype

	image gfx.Canvas
}

func (i *Image) Init(w, h int, pos vec.Coord, depth int, path string) {
	i.ElementPrototype.Init(w, h, pos, depth)

	i.image = gfx.ImportXPData(path)
}

func (i *Image) Render() {
	if !i.image.Ready() {
		return
	}

	i.image.Draw(&i.Canvas, vec.ZERO_COORD, 0)
}
