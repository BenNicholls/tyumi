package tyumi

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

// A dialog is a state that can report when it is done and can be shutdown.
type dialog interface {
	state

	Done() bool
}

// MessageDialog is a dialog that displays a simple message and an okay button.
type MessageDialog struct {
	State

	okayButton ui.Button
	done       bool
}

func NewMessageDialog(title, message string) (md *MessageDialog) {
	md = new(MessageDialog)
	md.Init(title, message)

	return
}

func (md *MessageDialog) Init(title, message string) {
	md.State.InitCentered(vec.Dims{mainConsole.Size().W / 2, 12})
	md.Window().EnableBorder()

	messageText := ui.NewTextbox(vec.Dims{md.Window().Size().W, ui.FIT_TEXT}, vec.Coord{0, 1}, 0, message, ui.JUSTIFY_CENTER)
	md.Window().AddChild(messageText)
	messageText.MoveTo(vec.Coord{0, (9 - messageText.Size().H) / 2})
	messageText.CenterHorizontal()

	md.okayButton.Init(vec.Dims{6, 1}, vec.Coord{0, 10}, 1, "Okay", nil)
	md.okayButton.EnableBorder()
	md.okayButton.Focus()
	md.Window().AddChild(&md.okayButton)
	md.okayButton.CenterHorizontal()

	md.Events().Listen(gfx.EV_ANIMATION_COMPLETE)
	md.SetEventHandler(md.HandleEvent)
}

func (md *MessageDialog) HandleEvent(game_event event.Event) (event_handled bool) {
	if game_event.ID() == gfx.EV_ANIMATION_COMPLETE {
		md.done = true
		return true
	}

	return
}

func (md MessageDialog) Done() bool {
	return md.done
}
