package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

// A handy text button that performs an action when pressed. Also plays an animation (defaults to a quick pulse).
type Button struct {
	Textbox

	OnPressCallback  func()
	OnPressAnimation gfx.Animator
}

func (b *Button) Init(size vec.Dims, pos vec.Coord, depth int, text string, on_press func()) {
	b.Textbox.Init(size, pos, depth, text, JUSTIFY_CENTER)
	b.OnPressCallback = on_press
	pressPulse := gfx.NewPulseAnimation(b.DrawableArea(), 0, 20, col.Pair{col.WHITE, col.WHITE})
	pressPulse.OneShot = true
	b.OnPressAnimation = &pressPulse
}

func NewButton(size vec.Dims, pos vec.Coord, depth int, text string, on_press func()) (b *Button) {
	b = new(Button)
	b.Init(size, pos, depth, text, on_press)

	return
}

func (b *Button) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.Handled() || key_event.PressType != input.KEY_PRESSED {
		return
	}

	switch key_event.Key {
	case input.K_RETURN:
		if b.OnPressCallback != nil {
			b.OnPressCallback()
			event_handled = true
		}

		if b.OnPressAnimation != nil {
			if !b.OnPressAnimation.IsPlaying() {
				b.AddAnimation(b.OnPressAnimation)
			}
			b.OnPressAnimation.Start()
		}
	}

	return
}
