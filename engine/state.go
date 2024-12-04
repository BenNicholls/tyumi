package engine

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
)

var currentState State //the current state object

const (
	FIT_CONSOLE int = 0 //window size flag
)

//A gameobject to be handled by Tyumi's state machine.
type State interface {
	Update()
	UpdateUI()
	Shutdown()
	Window() *ui.Container
	InputEvents() *event.Stream
	Events() *event.Stream
}

//An embeddable prototype that satisfies the State interface. Build around this for easier gamestate management.
type StatePrototype struct {
	window ui.Container

	events event.Stream       //for engine events, game events, etc. processed at the end of each tick
	inputEvents event.Stream  //for input events. processed at the start of each tick
}

//Init prepares the gamestate. If the console has been initialized, you can use FIT_CONSOLE as the
//width and/or height to have the state size itself automatically.
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

	sp.window = ui.NewContainer(w, h, 0, 0, 0)

	sp.events = event.NewStream(100)
	sp.inputEvents = event.NewStream(100)

	//setup automatic listening for input events.
	sp.inputEvents.Listen(input.EV_KEYBOARD)
}

func (sp *StatePrototype) Update() {
	return
}

//UpdateUI is called before each frame, allowing the user to apply ui changes for rendering all at once if they prefer.
//Otherwise they can implement Update() routines on the individual UI elements themselves and have them control their
//own behaviour.
func (sp *StatePrototype) UpdateUI() {
	return
}

func (sp *StatePrototype) Shutdown() {
	//TODO MAYBE: de-listen for input events??
	return
}

func (sp *StatePrototype) Window() *ui.Container {
	return &sp.window
}

func (sp *StatePrototype) Events() *event.Stream {
	return &sp.events
}

func (sp *StatePrototype) InputEvents() *event.Stream {
	return &sp.inputEvents
}

//sets the function for handling generic game events. these are collected during the tick(), and then processed
//at the end of each tick() in the order they were received.
func (sp *StatePrototype) AddEventHandler(h func(event.Event)) {
	sp.events.AddHandler(h)
}

//sets the function for handling inputs to the state object. inputs are collected, distributed and then
//processed at the beginning of each tick()
func (sp *StatePrototype) AddInputHandler(h func(event.Event)) {
	sp.inputEvents.AddHandler(h)
}

//InitMainState initializes a state to be run by Tyumi at the beginning of execution.
//This function DOES NOTHING if a state has already been initialized.
func InitMainState(s State) {
	if mainState != nil || s == nil {
		return
	}

	mainState = s
}
