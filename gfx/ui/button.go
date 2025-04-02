package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

var ACTION_BUTTON_PRESS = input.RegisterAction("Button Press")

func init() {
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_BUTTON_PRESS, input.K_RETURN)
}

// A handy text button that performs an action when pressed. Also plays an animation (defaults to a quick pulse).
type Button struct {
	Textbox

	OnPressCallback  func()
	OnPressAnimation gfx.Animator
}

func (b *Button) Init(size vec.Dims, pos vec.Coord, depth int, text string, on_press func()) {
	b.Textbox.Init(size, pos, depth, text, JUSTIFY_CENTER)
	b.TreeNode.Init(b)

	b.OnPressCallback = on_press
	pressPulse := gfx.NewPulseAnimation(b.DrawableArea(), 0, 20, col.Pair{col.WHITE, col.WHITE})
	pressPulse.OneShot = true
	pressPulse.Label = "Button Pressed"
	b.OnPressAnimation = &pressPulse
}

func NewButton(size vec.Dims, pos vec.Coord, depth int, text string, on_press func()) (b *Button) {
	b = new(Button)
	b.Init(size, pos, depth, text, on_press)

	return
}

func (b *Button) Press() {
	if b.OnPressCallback != nil {
		b.OnPressCallback()
	}

	if b.OnPressAnimation != nil {
		if !b.OnPressAnimation.IsPlaying() {
			b.AddAnimation(b.OnPressAnimation)
		}
		b.OnPressAnimation.Start()
	}
}

func (b *Button) HandleAction(action input.ActionID) (action_handled bool) {
	switch action {
	case ACTION_BUTTON_PRESS:
		b.Press()
	default:
		return false
	}

	return true
}
