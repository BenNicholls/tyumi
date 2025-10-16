package tyumi

import (
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

// A dialog is a scene that can report when it is done and can be shutdown.
type dialog interface {
	scene

	open()
	close()

	IsDone() bool
}

type Dialog struct {
	Scene

	Done bool // set this to true to have Tyumi close the dialog

	OnOpen func()
	OnDone func()
}

func (d *Dialog) open() {
	if d.OnOpen != nil {
		d.OnOpen()
	}

	d.window.Show()
}

func (d *Dialog) close() {
	if d.OnDone != nil {
		d.OnDone()
	}

	d.window.Hide()
}

// MessageDialog is a dialog that displays a simple message and an okay button.
type MessageDialog struct {
	Dialog

	okayButton ui.Button
}

func NewMessageDialog(title, message string) (md *MessageDialog) {
	md = new(MessageDialog)
	md.Init(title, message)

	return
}

func (md *MessageDialog) Init(title, message string) {
	md.Scene.InitCentered(vec.Dims{mainConsole.Size().W / 2, 12})
	md.Window().EnableBorder()

	messageText := ui.NewTextbox(vec.Dims{md.Window().Size().W, ui.FIT_TEXT}, vec.Coord{0, 1}, 0, message, ui.ALIGN_CENTER)
	md.Window().AddChild(messageText)
	messageText.MoveTo(vec.Coord{0, (9 - messageText.Size().H) / 2})
	messageText.CenterHorizontal()

	md.okayButton.Init(vec.Dims{6, 1}, vec.Coord{0, 10}, 1, "Okay", func() {
		md.CreateTimer(20, func() {
			md.Done = true
		})
	})
	md.okayButton.EnableBorder()
	md.okayButton.Focus()
	md.Window().AddChild(&md.okayButton)
	md.okayButton.CenterHorizontal()
}

func (md MessageDialog) IsDone() bool {
	return md.Done
}
