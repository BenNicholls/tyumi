package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/vec"
)

// An image. Right now only supports xp images loaded on element initialization, but in the future we'll have pre-
// loaded resources and other formats and all of that good stuff. For now just maybe don't initialize these every frame
// or whatever.
type Image struct {
	Element

	image gfx.Canvas
}

// loads the image at path. the image element will have the same size as the image.
func (i *Image) Init(pos vec.Coord, depth int, path string) {
	i.image = gfx.ImportXPData(path)
	i.Element.Init(i.image.Size(), pos, depth)
	i.TreeNode.Init(i)
}

func (i *Image) Render() {
	if !i.image.Ready() {
		return
	}

	i.image.Draw(&i.Canvas, vec.ZERO_COORD, 0)
}
