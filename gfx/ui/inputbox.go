package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
)

//Inputbox is a textbox that can accept and display keyboard input
type InputBox struct {
	Textbox

	cursor *InputCursorAnimation
}

func NewInputbox(w, h, x, y, z int) (ib InputBox) {
	ib.Textbox = NewTextbox(w, h, x, y, z, "", false)
	ib.cursor = NewInputCursorAnimation(0, 0, 1, 30)

	ib.AddAnimation(ib.cursor)

	return
}

func (ib *InputBox) HandleKeypress(e input.KeyboardEvent) {
	if text := e.Text(); text != "" {
		ib.Insert(text)
	} else if e.Key == input.K_BACKSPACE {
		ib.Delete()
	}
}

//Appends the provided string to the contents of the inputbox.
func (ib *InputBox) Insert(text string) {
	if w, _ := ib.Dims(); len(ib.text) == w*2 {
		return
	}

	ib.ChangeText(ib.text + text)
	ib.cursor.MoveTo(len(ib.text)/2, 0, len(ib.text)%2)
}

//Deletes the final character of the contents of the Inputbox
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

func NewInputCursorAnimation(x, y, z, rate int) (cursor *InputCursorAnimation) {
	cursor = new(InputCursorAnimation)
	cursor.BlinkAnimation = *gfx.NewBlinkAnimation(1, 1, x, y, z, gfx.NewTextVisuals(gfx.TEXT_BORDER_UD, gfx.TEXT_DEFAULT, col.WHITE, col.BLACK), rate)

	return
}

//Moves the cursor to (x, y), and blinks the indicated character (0 for left side, 1 for right side)
func (cursor *InputCursorAnimation) MoveTo(x, y, charNum int) {
	cursor.BlinkAnimation.MoveTo(x, y)
	if charNum%2 == 0 {
		cursor.Vis.ChangeChars(gfx.TEXT_BORDER_UD, gfx.TEXT_DEFAULT)
	} else {
		cursor.Vis.ChangeChars(gfx.TEXT_DEFAULT, gfx.TEXT_BORDER_UD)
	}
}
