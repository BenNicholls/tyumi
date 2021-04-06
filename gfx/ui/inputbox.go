package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
)

//Inputbox is a textbox that can accept and display keyboard input
type InputBox struct {
	Textbox

	cursor *gfx.BlinkAnimation
}

func NewInputbox(w, h, x, y, z int) (ib InputBox) {
	ib.Textbox = NewTextbox(w, h, x, y, z, "", false)
	ib.cursor = gfx.NewBlinkAnimation(1, 1, 0, 0, 0, gfx.Visuals{gfx.GLYPH_BLOCK, col.WHITE, col.BLACK}, 30)

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
	ib.cursor.MoveTo(len(ib.text)/2, 0)
}

//Deletes the final character of the contents of the Inputbox
func (ib *InputBox) Delete() {
	if len(ib.text) == 0 {
		return
	}

	ib.ChangeText(ib.text[:len(ib.text)-1])
	ib.cursor.MoveTo(len(ib.text)/2, 0)
}
