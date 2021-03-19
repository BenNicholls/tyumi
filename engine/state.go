package engine

import (
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/event"
)

var currentState State //the current state object

const (
	FIT_CONSOLE int = 0 //window size flag
)

//A gameobject to be handled by Tyumi.
type State interface {
	Update()
	UpdateUI()
	HandleEvent(event.Event)
	Shutdown()
	Window() *ui.Container	
}

//An embeddable prototype that satisfies the State interface. Build around this
//for easier gamestate management.
type StatePrototype struct {
	window ui.Container
} 

//InitWindow prepares the window of the gamestate. If the console has been initialized, you can use FIT_CONSOLE as the
//width and/or height to have the state size itself automatically.
func (sp *StatePrototype) InitWindow(w, h int) {
	if w == FIT_CONSOLE || h == FIT_CONSOLE {
		if !console.ready {
			log.Error("Cannot fit state window to console: console not initialized.")
			return
		}
		cw, ch := console.Dims()
		if w == FIT_CONSOLE {
			w = cw
		}
		if h == FIT_CONSOLE {
			h = ch
		}
	}

	sp.window = ui.NewContainer(w, h, 0, 0, 0)
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
	return
}

func (sp *StatePrototype) HandleEvent(e event.Event) {
	return
}

func (sp *StatePrototype) Window() *ui.Container {
	return &sp.window
}

//InitMainState initializes a state to be run by Tyumi at the beginning of execution.
//This function DOES NOTHING if a state has already been initialized. 
func InitMainState(s State) {
	if mainState != nil || s == nil {
		return
	}

	mainState = s
}