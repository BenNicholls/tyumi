package engine

import (
	"github.com/bennicholls/tyumi/gfx/ui"
)

var currentState State //the current state object

//A gameobject to be handled by Tyumi.
type State interface {
	Update()
	Shutdown()
}

//An embeddable prototype that satisfies the State interface. Build around this
//for easier gamestate management.
type StatePrototype struct {
	window *ui.Container
} 

func (sp *StatePrototype) Update() {
	return
}

func (sp *StatePrototype) Shutdown() {
	return
}

//InitMainState initializes a state to be run by Tyumi at the beginning of execution.
//This function DOES NOTHING if a state has already been initialized. 
func InitMainState(s State) {
	if mainState != nil || s == nil {
		return
	}

	mainState = s
}