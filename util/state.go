package util

import "github.com/bennicholls/tyumi/log"

type StateID uint32

const STATE_NONE StateID = 0

// A State is an element of a state machine that defines the behaviour that occurs when states are changed. Both OnEnter
// and OnLeave callbacks are optional.
type State struct {
	OnEnter func(previous_state StateID) // Callback run when state is entered.
	OnLeave func(next_state StateID)     // Callback run when state is left.
}

// A StateMachine allows you to register different states, one of which will be the CurrentState. You can change to
// another state, and the StateMachine will run any callbacks as appropriate. See the ChangeState() function for specifics.
// One of the basic fundamentals of a state machine is that all states are mutually exclusive; only one state will be
// active at a time. Use this in cases where something can be one of multiple states and you need to define behaviour
// for swapping between them.
// StateMachines do not need to be initialized and default to being in the state STATE_NONE, which has no callbacks. You
// can use this as a sort of default state if you like, or define your own and change away on setup, never to return
// again...
type StateMachine struct {
	current  StateID
	states   []State
	changing bool // will be true if we're in the middle of a change. used to detect callbacks that try to change the state again

	// OnStateChange is a callback that is run whenever the state machine changes state.
	OnStateChange func(previous, next StateID)
}

// RegisterState adds a state to the state machine, and returns an ID that you can use to refer to the state. Note that
// once set up States can't be modified so don't bother retaining raw State objects, just use the IDs.
// THINK: maybe mutable state objects would be useful? not sure how but... maybe?
func (sm *StateMachine) RegisterState(new_state State) StateID {
	sm.states = append(sm.states, new_state)
	return StateID(len(sm.states))
}

// ChangeState is the heart of the StateMachine. If next is a valid StateID, and is not the same as the current state, a
// change will be initiated. Callbacks are called in this order:
// 1) The Current State's OnLeave()
// 2) The StateMachine's OnStateChange()
// 3) The Next State's OnEnter()
// Of course all of these are optional, they do not need to be explicitly set.
// NOTE: callbacks CANNOT contain more state changes inside! this can cause an infinite loop if you're not careful which
// will make Tyumi cry. If this is detected it will throw an error and all but the initial call to ChangeState will fail.
func (sm *StateMachine) ChangeState(next StateID) {
	if next == sm.current {
		return
	}

	if sm.changing {
		log.Error("Cannot change state while already changing state! Did you put a ChangeState in a state callback?? Don't do that!")
		return
	}

	if !sm.isIDValid(next) {
		log.Error("Cannot change state! Invalid ID: ", next)
		return
	}

	sm.changing = true

	if sm.current != STATE_NONE {
		if currentState := sm.getStateByID(sm.current); currentState.OnLeave != nil {
			currentState.OnLeave(next)
		}
	}

	if sm.OnStateChange != nil {
		sm.OnStateChange(sm.current, next)
	}

	if next != STATE_NONE {
		if nextState := sm.getStateByID(next); nextState.OnEnter != nil {
			nextState.OnEnter(sm.current)
		}
	}

	sm.current = next
	sm.changing = false
}

// CurrentState retrieves the ID for the currently active state.
func (sm StateMachine) CurrentState() StateID {
	return sm.current
}

func (sm *StateMachine) getStateByID(id StateID) State {
	return sm.states[int(id)-1]
}

func (sm StateMachine) isIDValid(id StateID) bool {
	return int(id) <= len(sm.states)
}
