package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// Choicebox displays one element from a list, and allows the user to cycle through the options
type ChoiceBox struct {
	Textbox

	choices            []string
	currentChoiceIndex int            //will be -1 if no choices present
	arrowVisuals       [2]gfx.Visuals //LEFT and RIGHT
	arrowAnimations    [2]*gfx.FlashAnimation
}

func NewChoiceBox(w, h int, pos vec.Coord, depth int, choices ...string) (cb *ChoiceBox) {
	cb = new(ChoiceBox)
	cb.Textbox.Init(w, h, pos, depth, "No Choice", true)
	cb.currentChoiceIndex = -1

	cb.choices = choices
	if len(cb.choices) > 0 {
		cb.selectChoice(0)
	}

	cb.arrowVisuals[0] = gfx.NewGlyphVisuals(gfx.GLYPH_TRIANGLE_LEFT, cb.DefaultColours())
	cb.arrowVisuals[1] = gfx.NewGlyphVisuals(gfx.GLYPH_TRIANGLE_RIGHT, cb.DefaultColours())

	cb.arrowAnimations[0] = gfx.NewFlashAnimation(vec.Rect{vec.Coord{0, 0}, vec.Dims{1, 1}}, 1, col.Pair{col.RED, gfx.COL_DEFAULT}, 15)
	cb.arrowAnimations[1] = gfx.NewFlashAnimation(vec.Rect{vec.Coord{cb.Bounds().W-1, 0}, vec.Dims{1, 1}}, 1, col.Pair{col.RED, gfx.COL_DEFAULT}, 15)
	cb.AddAnimation(cb.arrowAnimations[0])
	cb.AddAnimation(cb.arrowAnimations[1])

	return
}

func (cb *ChoiceBox) selectChoice(index int) {
	if index == cb.currentChoiceIndex {
		return
	}

	if index >= len(cb.choices) {
		log.Error("Bad choice select!")
		return
	}

	cb.currentChoiceIndex = index
	cb.ChangeText(cb.choices[cb.currentChoiceIndex])
}

func (cb *ChoiceBox) Prev() {
	if len(cb.choices) < 2 {
		return
	}

	cb.selectChoice(util.CycleClamp(cb.currentChoiceIndex-1, 0, len(cb.choices)-1))
	cb.arrowAnimations[0].Play()
}

func (cb *ChoiceBox) Next() {
	if len(cb.choices) < 2 {
		return
	}

	cb.selectChoice(util.CycleClamp(cb.currentChoiceIndex+1, 0, len(cb.choices)-1))
	cb.arrowAnimations[1].Play()
}

func (cb *ChoiceBox) Render() {
	cb.Textbox.Render()
	
	if cb.updated {
		cb.DrawVisuals(vec.Coord{0, 0}, 1, cb.arrowVisuals[0])
		cb.DrawVisuals(vec.Coord{cb.Bounds().W - 1, 0}, 1, cb.arrowVisuals[1])
	}
}

func (cb *ChoiceBox) HandleKeypress(event *input.KeyboardEvent) {
	switch event.Key {
	case input.K_RIGHT:
		cb.Next()
		event.SetHandled()
	case input.K_LEFT:
		cb.Prev()
		event.SetHandled()
	}
}
