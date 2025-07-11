package ui

import (
	"time"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

var ACTION_INPUT_DELETE = input.RegisterAction("Delete Text")

func init() {
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_INPUT_DELETE, input.K_BACKSPACE)
}

// Inputbox is a textbox that can accept and display keyboard input
type InputBox struct {
	Textbox

	OnTextChanged  func() //callback triggered when input changes
	OnTextInputted func() //callback triggered when input is added
	OnTextDeleted  func() //callback triggered when input is deleted

	cursor         InputCursorAnimation
	inputLengthMax int //limit for input length. defaults to the width of the box
}

func NewInputbox(size vec.Dims, pos vec.Coord, depth, input_length int) (ib *InputBox) {
	ib = new(InputBox)
	ib.Init(size, pos, depth, input_length)

	return
}

// Initializes the inputbox. input_length limits the number of characters that can be written. if <= 0,
// input will instead be limited to the width of the inputbox
func (ib *InputBox) Init(size vec.Dims, pos vec.Coord, depth, input_length int) {
	ib.Textbox.Init(size, pos, depth, "", ALIGN_LEFT)
	ib.TreeNode.Init(ib)

	if input_length > 0 {
		ib.inputLengthMax = input_length
	} else {
		ib.inputLengthMax = size.W * 2
	}
	ib.cursor = NewInputCursorAnimation(vec.ZERO_COORD, 0, time.Second/2)
	ib.AddAnimation(&ib.cursor)
}

func (ib *InputBox) ChangeText(text string) {
	if ib.text == text {
		return
	}

	ib.Textbox.ChangeText(text)
	ib.cursor.MoveTo(len(ib.text)/2, 0, len(ib.text)%2)
	fireCallbacks(ib.OnTextChanged)
}

func (ib *InputBox) Focus() {
	if ib.focused {
		return
	}

	ib.setFocus(true)
	ib.cursor.Play()
}

func (ib *InputBox) Defocus() {
	if !ib.focused {
		return
	}

	ib.setFocus(false)
	ib.cursor.Stop()
	ib.Updated = true
}

func (ib *InputBox) HandleKeypress(event *input.KeyboardEvent) (event_handled bool) {
	if event.PressType == input.KEY_RELEASED {
		return
	}

	if text := event.Text(); text != "" {
		ib.Insert(text)
		event_handled = true
	}

	return
}

func (ib *InputBox) HandleAction(action input.ActionID) (action_handled bool) {
	switch action {
	case ACTION_INPUT_DELETE:
		ib.Delete()
	default:
		return
	}

	return true
}

// Appends the provided string to the contents of the inputbox.
func (ib *InputBox) Insert(input string) {
	new_text := ib.text + input
	if len(new_text) > ib.inputLengthMax {
		return
	}

	ib.ChangeText(new_text)
	fireCallbacks(ib.OnTextInputted)
}

// Deletes the final character of the contents of the Inputbox
func (ib *InputBox) Delete() {
	if len(ib.text) == 0 {
		return
	}

	ib.ChangeText(ib.text[:len(ib.text)-1])
	fireCallbacks(ib.OnTextDeleted)
}

func (ib InputBox) InputtedText() string {
	return ib.text
}

type InputCursorAnimation struct {
	gfx.BlinkAnimation
}

func NewInputCursorAnimation(pos vec.Coord, depth int, rate time.Duration) (cursor InputCursorAnimation) {
	vis := gfx.NewTextVisuals(gfx.TEXT_BORDER_UD, gfx.TEXT_DEFAULT, col.Pair{gfx.COL_DEFAULT, gfx.COL_DEFAULT})
	cursor = InputCursorAnimation{
		BlinkAnimation: gfx.NewBlinkAnimation(vec.Rect{pos, vec.Dims{1, 1}}, depth, vis, rate),
	}
	cursor.Start()

	return
}

// Moves the cursor to (x, y), and blinks the indicated character (0 for left side, 1 for right side)
func (cursor *InputCursorAnimation) MoveTo(x, y, charNum int) {
	cursor.BlinkAnimation.MoveTo(vec.Coord{x, y})
	if charNum%2 == 0 {
		cursor.Vis.SetText(gfx.TEXT_BORDER_UD, gfx.TEXT_DEFAULT)
	} else {
		cursor.Vis.SetText(gfx.TEXT_DEFAULT, gfx.TEXT_BORDER_UD)
	}
}
