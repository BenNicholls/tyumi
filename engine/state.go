package engine

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

var currentState State //the current state object

const (
	FIT_CONSOLE int = 0 //window size flag
)

// A gameobject to be handled by Tyumi's state machine.
type State interface {
	Update()
	UpdateUI()
	Shutdown()
	Window() *ui.Window
	InputEvents() *event.Stream
	Events() *event.Stream
	Ready() bool
	IsBlocked() bool
}

// An embeddable prototype that satisfies the State interface. Build around this for easier gamestate management.
type StatePrototype struct {
	window *ui.Window

	events       event.Stream  //for engine events, game events, etc. processed at the end of each tick
	inputEvents  event.Stream  //for input events. processed at the start of each tick
	inputHandler event.Handler //user-provided input handling function. runs AFTER the UI has had a chance to process input.

	ready bool // indicates the state has been successfully initialized
}

// Init prepares the gamestate. If the console has been initialized, you can use FIT_CONSOLE as the
// width and/or height to have the state size itself automatically.
func (sp *StatePrototype) Init(w, h int) {
	if w == FIT_CONSOLE || h == FIT_CONSOLE {
		if !console.ready {
			log.Error("Cannot fit state window to console: console not initialized.")
			return
		}

		if w == FIT_CONSOLE {
			w = console.Size().W
		}
		if h == FIT_CONSOLE {
			h = console.Size().H
		}
	}

	sp.window = ui.NewWindow(w, h, vec.ZERO_COORD, 0)

	sp.events = event.NewStream(100, nil)
	sp.inputEvents = event.NewStream(100, sp.handleInput)

	//setup automatic listening for input events.
	sp.inputEvents.Listen(input.EV_KEYBOARD, input.EV_MOUSEBUTTON, input.EV_MOUSEMOVE)
	sp.ready = true
}

func (sp *StatePrototype) Update() {
	return
}

// UpdateUI is called before each frame is rendered, allowing the user to apply ui changes for rendering all at once if
// they prefer. Otherwise they can implement Update() routines on the individual UI elements themselves and have them
// control their own behaviour.
func (sp *StatePrototype) UpdateUI() {
	return
}

func (sp *StatePrototype) Shutdown() {
	//TODO MAYBE: de-listen for input events??
	return
}

func (sp *StatePrototype) Window() *ui.Window {
	return sp.window
}

func (sp *StatePrototype) Events() *event.Stream {
	return &sp.events
}

func (sp *StatePrototype) InputEvents() *event.Stream {
	return &sp.inputEvents
}

// sets the function for handling generic game events. these are collected during the tick(), and then processed
// at the end of each tick() in the order they were received.
func (sp *StatePrototype) SetEventHandler(handler event.Handler) {
	sp.events.AddHandler(handler)
}

func (sp *StatePrototype) handleInput(event event.Event) (event_handled bool) {
	switch event.ID() {
	case input.EV_KEYBOARD:
		event_handled = sp.window.HandleKeypress(event.(*input.KeyboardEvent))
	}

	if sp.inputHandler != nil {
		event_handled = event_handled || sp.inputHandler(event)
	}

	return
}

// sets the function for handling inputs to the state object. inputs are collected, distributed and then
// processed at the beginning of each tick(). This handler is called after the UI has had a chance to handle
// the input. If the UI handles the input, event.Handled() will be true. You can still choose to ignore that and
// handle the event again if you like though.
func (sp *StatePrototype) SetInputHandler(handler event.Handler) {
	sp.inputHandler = handler
}

func (sp *StatePrototype) Ready() bool {
	return sp.ready
}

func (sp *StatePrototype) IsBlocked() bool {
	return sp.window.IsBlocked()
}

// SetInitialMainState sets a state to be run by Tyumi at the beginning of execution.
// This function DOES NOTHING if a state has already been initialized.
func SetInitialMainState(s State) {
	if mainState != nil {
		return
	}

	if !console.ready {
		log.Error("Cannot set main state: console not initialized. Run InitConsole() first.")
		return
	}

	if s == nil || !s.Ready() {
		log.Error("Cannot set main state: state not initialized or ready.")
		return
	}

	mainState = s
}
