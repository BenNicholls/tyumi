package tyumi

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

var currentState state

const (
	FIT_CONSOLE int = 0 //window size flag
)

// A gameobject to be handled by Tyumi's state machine.
type state interface {
	Update()
	UpdateUI()
	Shutdown()
	Window() *ui.Window
	InputEvents() *event.Stream
	Events() *event.Stream
	Ready() bool
	IsBlocked() bool
}

// State is the base implementation for Tyumi game state object. States contain a window, where the programs UI is built,
// as well as the machinery for handling game and input events. Custom states can be made by embedding this and
// overriding the virtual functions defined here. Most important is the Update() function, which runs once per-tick
// and should contain your main game code.
type State struct {
	window *ui.Window

	events       event.Stream  //for engine events, game events, etc. processed at the end of each tick
	inputEvents  event.Stream  //for input events. processed at the start of each tick
	inputHandler event.Handler //user-provided input handling function. runs AFTER the UI has had a chance to process input.

	ready bool // indicates the state has been successfully initialized
}

// Init prepares the gamestate. If the console has been initialized, you can use FIT_CONSOLE as the
// width and/or height to have the state size itself automatically.
func (sp *State) Init(size vec.Dims) {
	if size.W == FIT_CONSOLE || size.H == FIT_CONSOLE {
		if !mainConsole.ready {
			log.Error("Cannot fit state window to console: console not initialized.")
			return
		}

		if size.W == FIT_CONSOLE {
			size.W = mainConsole.Size().W
		}
		if size.H == FIT_CONSOLE {
			size.H = mainConsole.Size().H
		}
	}

	sp.window = ui.NewWindow(size, vec.ZERO_COORD, 0)

	sp.events = event.NewStream(100, nil)
	sp.inputEvents = event.NewStream(100, sp.handleInput)

	//setup automatic listening for input events.
	sp.inputEvents.Listen(input.EV_KEYBOARD, input.EV_MOUSEBUTTON, input.EV_MOUSEMOVE)
	sp.ready = true
}

// Update is run each tick, after input has been handled and before UI is updated/rendered. Override this function
// with your primary game code!
func (sp *State) Update() {
	return
}

// UpdateUI is called before each frame is rendered, allowing you to apply ui changes for rendering all at once if
// you prefer. Otherwise you can implement Update() functions on the individual UI elements themselves and have them
// control their own behaviour.
func (sp *State) UpdateUI() {
	return
}

func (sp *State) Shutdown() {
	//TODO MAYBE: de-listen for input events??
	return
}

func (sp *State) Window() *ui.Window {
	return sp.window
}

func (sp *State) Events() *event.Stream {
	return &sp.events
}

func (sp *State) InputEvents() *event.Stream {
	return &sp.inputEvents
}

// sets the function for handling game events. these are collected during Update(), and then processed
// at the end of each Update() in the order they were received.
func (sp *State) SetEventHandler(handler event.Handler) {
	sp.events.AddHandler(handler)
}

func (sp *State) handleInput(event event.Event) (event_handled bool) {
	switch event.ID() {
	case input.EV_KEYBOARD:
		event_handled = sp.window.HandleKeypress(event.(*input.KeyboardEvent))
	}

	if sp.inputHandler != nil {
		event_handled = event_handled || sp.inputHandler(event)
	}

	return
}

// Sets the function for handling inputs to the state object. Inputs are collected, distributed and then
// processed at the beginning of each tick(). This handler is called after the UI has had a chance to handle
// the input. If the UI handles the input, event.Handled() will be true. You can still choose to ignore that and
// handle the event again if you like though.
func (sp *State) SetInputHandler(handler event.Handler) {
	sp.inputHandler = handler
}

func (sp *State) Ready() bool {
	return sp.ready
}

// Returns true if updating has been blocked. Currently this only happens from blocking animations, but later might
// also indicate that the game is paused perhaps.
func (sp *State) IsBlocked() bool {
	return sp.window.IsBlocked()
}

// SetInitialMainState sets a state to be run by Tyumi at the beginning of execution.
// This function DOES NOTHING if a state has already been initialized.
func SetInitialMainState(s state) {
	if currentState != nil {
		return
	}

	if !mainConsole.ready {
		log.Error("Cannot set main state: console not initialized. Run InitConsole() first.")
		return
	}

	if s == nil || !s.Ready() {
		log.Error("Cannot set main state: state not initialized or ready.")
		return
	}

	currentState = s
}
