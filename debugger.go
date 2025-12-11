package tyumi

import (
	"fmt"
	"strings"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

var EV_LOGMESSAGE = event.Register("Message Logged")

type LogEvent struct {
	event.EventPrototype

	Entry log.Entry
}

var debugger *debugDialog

type debugDialog struct {
	Dialog

	container ui.PageContainer

	commandPage    *ui.Page
	commandInput   ui.InputBox
	commandDisplay ui.List

	logPage    *ui.Page
	logDisplay ui.List

	commands map[string]debugCommand
}

func (d *debugDialog) Init() {
	d.InitCentered(vec.Dims{30, 15})
	d.Window().SetupBorder("Debugger", "")
	d.Window().SendEventsToUnfocused = true

	d.container.Init(d.window.Size(), vec.ZERO_COORD, ui.BorderDepth)
	d.container.EnableBorder()
	d.container.AcceptInput = true
	d.Window().AddChild(&d.container)

	d.commandPage = d.container.CreatePage("CMD")
	d.commandInput.Init(vec.Dims{d.commandPage.Size().W-2, 1}, vec.Coord{2, d.commandPage.Size().H - 1}, 0, 0)
	d.commandInput.AcceptInput = true
	d.commandDisplay.Init(d.commandPage.Size().Shrink(0, 2), vec.ZERO_COORD, ui.BorderDepth)
	d.commandDisplay.EnableBorder()
	d.commandDisplay.SetCapacity(100)
	d.commandDisplay.SetPadding(1)
	d.commandDisplay.OnItemInserted = func() { d.commandDisplay.ScrollToBottom() }
	d.commandDisplay.AcceptInput = true
	d.commandPage.AddChildren(&d.commandInput, &d.commandDisplay)
	d.commandPage.AddChild(ui.NewTextbox(vec.Dims{2,1}, vec.Coord{0, d.commandPage.Size().H - 1}, 0, ">>>", ui.ALIGN_LEFT))

	d.logPage = d.container.CreatePage("LOG")
	d.logDisplay.Init(d.logPage.Size().Shrink(0, 2), vec.ZERO_COORD, ui.BorderDepth)
	d.logDisplay.SetCapacity(200)
	d.logDisplay.OnItemInserted = func() { d.logDisplay.ScrollToBottom() }
	d.logDisplay.AcceptInput = true
	d.logPage.AddChild(&d.logDisplay)

	// add already existing log messages
	for _, entry := range log.GetLogs() {
		d.logDisplay.InsertText(ui.ALIGN_LEFT, entry.SimpleString())
	}

	d.keypressInputHandler = d.handleKeyEvent
	d.Persistent = true
}

func (d *debugDialog) handleKeyEvent(ke *input.KeyboardEvent) (event_handled bool) {
	switch ke.Key {
	case input.K_F12, input.K_ESCAPE:
		d.Done = true
		return true
	}

	switch d.container.GetPageIndex() {
	case 0: //CMD
		switch ke.Key {
		case input.K_RETURN:
			d.parseCommand(d.commandInput.InputtedText())
			d.commandInput.DeleteAll()
		default:
			return
		}
	}

	return true
}

func (d *debugDialog) parseCommand(commandText string) {
	tokens := strings.Split(commandText, " ")
	if commandText == "" || len(tokens) == 0 {
		return
	}

	commandName := strings.TrimSpace(tokens[0])

	switch commandName {
	case "help":
		d.addToCommandDisplay(commandText, "This is the debugger! Type 'commands' to see available commands. Use PGUP/PGDOWN to scroll through messages. Press TAB to cycle to the next page.")
	case "commands":
		commandList := ""
		for name, command := range d.commands {
			commandList += fmt.Sprintf("|||%s: %s/n", name, command.desc)
		}

		d.addToCommandDisplay(commandText, "Available commands are: /n/n"+commandList)
	default:
		if command, ok := debugger.commands[commandName]; ok {
			output := command.handler(tokens[1:])
			d.addToCommandDisplay(commandText, output)
		} else {
			d.addToCommandDisplay(commandText, "No command found for "+commandName)
		}
	}
}

func (d *debugDialog) addToCommandDisplay(command, text string) {
	d.commandDisplay.InsertText(ui.ALIGN_LEFT, fmt.Sprintf(">>> %s/n   %s", command, text))
}

func (d *debugDialog) addToLogDisplay(entry log.Entry) {
	d.logDisplay.InsertText(ui.ALIGN_LEFT, entry.SimpleString())
}

func RegisterDebugCommand(name, desc string, handler func(args []string) string) {
	if !Debug || handler == nil {
		return
	}

	if debugger.commands == nil {
		debugger.commands = make(map[string]debugCommand)
	}

	if _, ok := debugger.commands[name]; ok {
		log.Debug("Overwriting debug command!")
	}

	debugger.commands[name] = debugCommand{
		name:    name,
		desc:    desc,
		handler: handler,
	}
	log.Debug("Registered debug command ", name)
}

type debugCommand struct {
	name    string
	desc    string
	handler func(args []string) string
}
