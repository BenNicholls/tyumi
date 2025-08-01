package ui

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// ProgressBar is a textbox that can display a coloured bar in the background to indicate some progress. Defaults to
// 100% progress.
type ProgressBar struct {
	Textbox

	progressColour col.Colour
	progress       int
}

func (pb *ProgressBar) Init(size vec.Dims, pos vec.Coord, depth int, progress_colour col.Colour, text string) {
	pb.Textbox.Init(size, pos, depth, text, ALIGN_CENTER)
	pb.TreeNode.Init(pb)

	pb.progressColour = progress_colour
	pb.progress = 100
}

func (pb ProgressBar) GetProgress() int {
	return pb.progress
}

// Returns a value [0,1] representing the progress value of the progressbar.
func (pb ProgressBar) GetProgressNormalized() float64 {
	switch pb.progress {
	case 100:
		return 1
	case 0:
		return 0
	default:
		return float64(pb.progress) / 100
	}
}

// SetProgress determines the length of the progress bar. It takes a percentage (which it clamps to [0, 100] just in
// case you forget).
func (pb *ProgressBar) SetProgress(progress_pct int) {
	progress := util.Clamp(progress_pct, 0, 100)
	if pb.progress == progress {
		return
	}

	pb.progress = progress
	pb.Updated = true
}

// SetProgressColour does precisely what you expect it to.
func (pb *ProgressBar) SetProgressColour(colour col.Colour) {
	if pb.progressColour == colour {
		return
	}

	pb.progressColour = colour
	pb.Updated = true
}

func (pb *ProgressBar) Render() {
	pb.Textbox.Render()
	if pb.progress == 0 {
		return
	}

	barLength := util.Lerp(0, pb.size.W, pb.progress, 100)

	for x := range barLength {
		for y := range pb.DrawableArea().H {
			pb.DrawColours(vec.Coord{x, y}, 0, col.Pair{col.NONE, pb.progressColour})
		}
	}
}
