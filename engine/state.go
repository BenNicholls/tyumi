package engine

var currentState State //the current state object

//The currently running game object.
type State interface {

}

//An embeddable prototype that satisfies the State interface. Build around this
//for easier gamestate management.
type StatePrototype struct {

} 