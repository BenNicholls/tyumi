package tyumi

import (
	"strings"

	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

var debugger *debugDialog

type debugDialog struct {
	Dialog

	// UI things
	input ui.InputBox

	commands map[string]func([]string)
}

func (d *debugDialog) Init() {
	d.InitCentered(vec.Dims{20, 1})
	d.Window().SetupBorder("Debugger", "type a debug command and ENTER")

	d.input.Init(d.Window().Size(), vec.ZERO_COORD, 0, 0)
	d.input.Focus()
	d.Window().AddChild(&d.input)

	d.keypressInputHandler = d.handleKeyEvent
	d.Persistent = true
}

func (d *debugDialog) handleKeyEvent(ke *input.KeyboardEvent) (event_handled bool) {
	switch ke.Key {
	case input.K_F12, input.K_ESCAPE:
		d.Done = true
	case input.K_RETURN:
		d.parseCommand(d.input.InputtedText())
		d.input.DeleteAll()
	default:
		return
	}

	return true
}

func (d *debugDialog) parseCommand(commandText string) {
	tokens := strings.Split(commandText, " ")
	if commandText == "" || len(tokens) == 0 {
		return
	}

	commandName := tokens[0]
	if handler, ok := debugger.commands[commandName]; ok {
		handler(tokens[1:])
	} else {
		log.Debug("No command found for ", commandName)
	}
}

func RegisterDebugCommand(name string, handler func(args []string)) {
	if !Debug || handler == nil {
		return
	}

	if debugger.commands == nil {
		debugger.commands = make(map[string]func([]string))
	}

	if _, ok := debugger.commands[name]; ok {
		log.Debug("Overwriting debug command!")
	}

	debugger.commands[name] = handler
	log.Debug("Registered debug command ", name)
}
