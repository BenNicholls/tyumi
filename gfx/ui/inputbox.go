package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

// Inputbox is a textbox that can accept and display keyboard input
type InputBox struct {
	Textbox

	cursor *InputCursorAnimation
}

func NewInputbox(w, h int, pos vec.Coord, depth int) (ib *InputBox) {
	ib = new(InputBox)
	ib.Textbox.Init(w, h, pos, depth, "", false)
	ib.cursor = NewInputCursorAnimation(vec.Coord{0, 0}, 0, 30)

	ib.AddAnimation(ib.cursor)

	return
}

func (ib *InputBox) HandleKeypress(event input.KeyboardEvent) {
	if text := event.Text(); text != "" {
		ib.Insert(text)
	} else if event.Key == input.K_BACKSPACE {
		ib.Delete()
	}
}

// Appends the provided string to the contents of the inputbox.
func (ib *InputBox) Insert(text string) {
	if w := ib.Size().W; len(ib.text) == w*2 {
		return
	}

	ib.ChangeText(ib.text + text)
	ib.cursor.MoveTo(len(ib.text)/2, 0, len(ib.text)%2)
}

// Deletes the final character of the contents of the Inputbox
func (ib *InputBox) Delete() {
	if len(ib.text) == 0 {
		return
	}

	ib.ChangeText(ib.text[:len(ib.text)-1])
	ib.cursor.MoveTo(len(ib.text)/2, 0, len(ib.text)%2)
}

type InputCursorAnimation struct {
	gfx.BlinkAnimation
}

func NewInputCursorAnimation(pos vec.Coord, depth, rate int) (cursor *InputCursorAnimation) {
	vis := gfx.NewTextVisuals(gfx.TEXT_BORDER_UD, gfx.TEXT_DEFAULT, col.Pair{col.WHITE, col.BLACK})
	cursor = &InputCursorAnimation{
		BlinkAnimation: *gfx.NewBlinkAnimation(pos, vec.Dims{1, 1}, depth, vis, rate),
	}
	cursor.Start()

	return
}

// Moves the cursor to (x, y), and blinks the indicated character (0 for left side, 1 for right side)
func (cursor *InputCursorAnimation) MoveTo(x, y, charNum int) {
	cursor.BlinkAnimation.MoveTo(vec.Coord{x, y})
	if charNum%2 == 0 {
		cursor.Vis.ChangeChars(gfx.TEXT_BORDER_UD, gfx.TEXT_DEFAULT)
	} else {
		cursor.Vis.ChangeChars(gfx.TEXT_DEFAULT, gfx.TEXT_BORDER_UD)
	}
}
