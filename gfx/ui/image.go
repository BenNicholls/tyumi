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

func NewImage(pos vec.Coord, depth int, path string) (i *Image) {
	i = new(Image)
	i.Init(pos, depth, path)

	return
}

// loads the image at path. the image element will have the same size as the image. If no image is provided, or if the
// provided image fails to load for whatever reason, size will be 1x1.
func (i *Image) Init(pos vec.Coord, depth int, path string) {
	i.LoadImage(path)
	if i.image.Ready() {
		i.Element.Init(i.image.Size(), pos, depth)
	} else {
		i.Element.Init(vec.Dims{1, 1}, pos, depth)
	}
	i.TreeNode.Init(i)
}

// Loads an xp image from disk to be displayed in the image element.
func (i *Image) LoadImage(path string) {
	if path == "" {
		return
	}

	i.image = gfx.ImportXPData(path)
	if i.image.Ready() {
		if i.image.Size() != i.Size() {
			i.Resize(i.image.Size())
		}
		i.Updated = true
	}
}

func (i *Image) Render() {
	if !i.image.Ready() {
		return
	}

	i.image.Draw(&i.Canvas, vec.ZERO_COORD, 0)
}
