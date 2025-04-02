package ui

import (
	"math/rand"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var EV_CHOICE_CHANGED = event.Register("Choice Changed", event.SIMPLE)

var ACTION_CHOICE_NEXT = input.RegisterAction("Select Next Choice")
var ACTION_CHOICE_PREV = input.RegisterAction("Select Previous Choice")

func init() {
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_CHOICE_NEXT, input.K_RIGHT)
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_CHOICE_PREV, input.K_LEFT)
}

// Choicebox displays one element from a list, and allows the user to cycle through the options
type ChoiceBox struct {
	Textbox

	choices            []string
	currentChoiceIndex int            //will be -1 if no choices present
	arrowVisuals       [2]gfx.Visuals //LEFT and RIGHT
	arrowAnimations    [2]gfx.FlashAnimation
}

func (cb *ChoiceBox) Init(size vec.Dims, pos vec.Coord, depth int, choices ...string) {
	cb.Textbox.Init(size, pos, depth, "No Choice", JUSTIFY_CENTER)
	cb.TreeNode.Init(cb)

	cb.currentChoiceIndex = -1
	cb.choices = choices
	if len(cb.choices) > 0 {
		cb.selectChoice(0)
	}

	cb.arrowVisuals[0] = gfx.NewGlyphVisuals(gfx.GLYPH_TRIANGLE_LEFT, cb.DefaultColours())
	cb.arrowVisuals[1] = gfx.NewGlyphVisuals(gfx.GLYPH_TRIANGLE_RIGHT, cb.DefaultColours())

	cb.arrowAnimations[0] = gfx.NewFlashAnimation(vec.Rect{vec.Coord{0, 0}, vec.Dims{1, 1}}, 1, col.Pair{col.RED, gfx.COL_DEFAULT}, 15)
	cb.arrowAnimations[1] = gfx.NewFlashAnimation(vec.Rect{vec.Coord{cb.Size().W - 1, 0}, vec.Dims{1, 1}}, 1, col.Pair{col.RED, gfx.COL_DEFAULT}, 15)
	cb.AddAnimation(&cb.arrowAnimations[0])
	cb.AddAnimation(&cb.arrowAnimations[1])
}

func NewChoiceBox(size vec.Dims, pos vec.Coord, depth int, choices ...string) (cb *ChoiceBox) {
	cb = new(ChoiceBox)
	cb.Init(size, pos, depth, choices...)

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
	event.FireSimple(EV_CHOICE_CHANGED)
}

func (cb *ChoiceBox) Prev() {
	if len(cb.choices) < 2 {
		return
	}

	cb.selectChoice(util.CycleClamp(cb.currentChoiceIndex-1, 0, len(cb.choices)-1))
	cb.arrowAnimations[0].Start()
}

func (cb *ChoiceBox) Next() {
	if len(cb.choices) < 2 {
		return
	}

	cb.selectChoice(util.CycleClamp(cb.currentChoiceIndex+1, 0, len(cb.choices)-1))
	cb.arrowAnimations[1].Start()
}

// Selects a random choice from the choices available.
func (cb *ChoiceBox) RandomizeChoice() {
	cb.selectChoice(rand.Intn(len(cb.choices)))
}

func (cb *ChoiceBox) GetChoiceIndex() int {
	return cb.currentChoiceIndex
}

func (cb *ChoiceBox) Render() {
	cb.Textbox.Render()

	cb.DrawVisuals(vec.Coord{0, 0}, 1, cb.arrowVisuals[0])
	cb.DrawVisuals(vec.Coord{cb.Size().W - 1, 0}, 1, cb.arrowVisuals[1])
}

func (cb *ChoiceBox) HandleAction(action input.ActionID) (action_handled bool) {
	switch action {
	case ACTION_CHOICE_NEXT:
		cb.Next()
	case ACTION_CHOICE_PREV:
		cb.Prev()
	default:
		return false
	}

	return true
}
