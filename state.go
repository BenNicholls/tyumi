package tyumi

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

var EV_CHANGESTATE = event.Register("State Change Event", event.COMPLEX)
var currentState state
var activeState state // the state being updated, checked each frame

const (
	FIT_CONSOLE int = 0 //window size flag
)

// A gameobject to be handled by Tyumi's state machine.
type state interface {
	Window() *ui.Window
	Ready() bool

	Update()
	UpdateUI()
	IsBlocked() bool

	InputEvents() *event.Stream
	Events() *event.Stream

	OpenDialog(dialog)

	getActiveState() state
	flushInputs()
	shutdown()
}

// State is the base implementation for Tyumi's game state object. States contain a window, where the programs UI is
// built, as well as the machinery for handling game and input events. Custom states can be made by embedding this and
// overriding the virtual functions defined here. Most important is the Update() function, which runs once per-tick
// and should contain your main game code.
type State struct {
	window *ui.Window

	subState dialog

	events       event.Stream  //for engine events, game events, etc. processed at the end of each tick
	inputEvents  event.Stream  //for input events. processed at the start of each tick
	inputHandler event.Handler //user-provided input handling function. runs AFTER the UI has had a chance to process input.

	ready bool // indicates the state has been successfully initialized
}

// Init prepares the gamestate, defaulting to a window the full size of the console.
// NOTE: If you want a border drawn around the window, use InitBordered() instead since Tyumi draws borders *around*
// objects and if the window is the size of the console you wouldn't be able to see it.
func (s *State) Init() {
	if !mainConsole.ready {
		log.Error("Cannot fit state window to console: console not initialized.")
		return
	}

	s.init(mainConsole.Size(), vec.ZERO_COORD, false)
}

// InitBordered prepares the gamestate, defaulting to a window the full size of the console with a border drawn around
// the outside.
func (s *State) InitBordered() {
	if !mainConsole.ready {
		log.Error("Cannot fit state window to console: console not initialized.")
		return
	}

	s.init(mainConsole.Size().Shrink(2, 2), vec.Coord{1, 1}, true)
}

// InitCentered prepares the gamestate, centering the state's window inside the console.
func (s *State) InitCentered(size vec.Dims) {
	if !mainConsole.ready {
		log.Error("Cannot fit state window to console: console not initialized.")
		return
	}

	pos := vec.Coord{(mainConsole.Size().W - size.W) / 2, (mainConsole.Size().H - size.H) / 2}
	s.init(size, pos, false)
}

// InitCustom prepares the gamestate, creating a window using the given size and position. The position will be relative
// to the console. If the console has been initialized, you can use FIT_CONSOLE as the width and/or height to have the
// state size itself automatically.
func (s *State) InitCustom(size vec.Dims, pos vec.Coord) {
	s.init(size, pos, false)
}

func (s *State) init(size vec.Dims, pos vec.Coord, bordered bool) {
	if s.ready {
		log.Error("Trying to initialize a state more than once. Don't do that.")
		return
	}

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

	s.window = ui.NewWindow(size, pos, 0)
	if bordered {
		s.window.EnableBorder()
	}

	s.events = event.NewStream(100, nil)
	s.inputEvents = event.NewStream(100, s.handleInput)

	//setup automatic listening for input events.
	s.inputEvents.Listen(input.EV_KEYBOARD, input.EV_MOUSEBUTTON, input.EV_MOUSEMOVE)
	s.ready = true
}

// Update is run each tick, after input has been handled and before UI is updated/rendered. Override this function
// with your primary game code!
func (s *State) Update() {
	return
}

// UpdateUI is called before each frame is rendered, allowing you to apply ui changes for rendering all at once if
// you prefer. Otherwise you can implement Update() functions on the individual UI elements themselves and have them
// control their own behaviour.
func (s *State) UpdateUI() {
	return
}

func (s *State) OpenDialog(subState dialog) {
	if !subState.Ready() {
		log.Error("Could not open dialog, dialog not initialized.")
		return
	}

	s.subState = subState
}

func (s *State) getActiveState() state {
	if s.subState == nil {
		return s
	}

	if s.subState.Done() {
		log.Debug("shutting down subtate")
		s.subState.shutdown()
		s.subState = nil
		return s
	} else {
		return s.subState.getActiveState()
	}
}

func (s *State) shutdown() {
	s.Shutdown()

	s.events.Close()
	s.inputEvents.Close()

	// in theory states should be freed from memory after being shutdown so this is pointless, but on the off chance a
	// reference is hanging around maybe this will help catch a bug.
	s.ready = false
}

// Shutdown is called when the state is no longer needed and should cleanly pack itself away (for example, when
// switching to another state or closing the program). Override this function and use it to free resources, save things
// to disk, whatever you need to do.
func (s *State) Shutdown() {
	return
}

func (s *State) Window() *ui.Window {
	return s.window
}

func (s *State) Events() *event.Stream {
	return &s.events
}

func (s *State) InputEvents() *event.Stream {
	return &s.inputEvents
}

// sets the function for handling game events. these are collected during Update(), and then processed
// at the end of each Update() in the order they were received.
func (s *State) SetEventHandler(handler event.Handler) {
	s.events.AddHandler(handler)
}

func (s *State) handleInput(event event.Event) (event_handled bool) {
	switch event.ID() {
	case input.EV_KEYBOARD:
		event_handled = s.window.HandleKeypress(event.(*input.KeyboardEvent))
	}

	if s.inputHandler != nil {
		event_handled = event_handled || s.inputHandler(event)
	}

	return
}

// Sets the function for handling inputs to the state object. Inputs are collected, distributed and then
// processed at the beginning of each tick(). This handler is called after the UI has had a chance to handle
// the input. If the UI handles the input, event.Handled() will be true. You can still choose to ignore that and
// handle the event again if you like though.
func (s *State) SetInputHandler(handler event.Handler) {
	s.inputHandler = handler
}

func (s *State) flushInputs() {
	s.inputEvents.Flush()

	if s.subState != nil {
		s.subState.flushInputs()
	}
}

func (s *State) Ready() bool {
	return s.ready
}

// Returns true if updating has been blocked. Currently this only happens from blocking animations, but later might
// also indicate that the game is paused perhaps.
func (s *State) IsBlocked() bool {
	return s.window.IsBlocked()
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

type StateChangeEvent struct {
	event.EventPrototype

	newState state
}

// ChangeState changes the current state being run in Tyumi's gameloop. The change is done at the end of the current
// engine tick. The old state's Shutdown() method is called before we swap in the new one. Be sure to initialize the
// new state before calling ChangeState(), otherwise no change will happen and the old state will remain.
func ChangeState(new_state state) {
	if new_state == nil || !new_state.Ready() {
		log.Error("Could not change state: state invalid or not initialized.")
		return
	}

	//if user tries to use this to setup the initial main state, just forgive them their sin and do it. no need to
	//harass them with "the correct way".
	if currentState == nil {
		SetInitialMainState(new_state)
		return
	}

	event.Fire(&StateChangeEvent{
		EventPrototype: event.New(EV_CHANGESTATE),
		newState:       new_state,
	})
}
