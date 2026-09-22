package util

import "testing"

func TestStateMachineDefaultState(t *testing.T) {
	var sm StateMachine

	if sm.CurrentState() != STATE_NONE {
		t.Errorf("CurrentState() = %v, want STATE_NONE", sm.CurrentState())
	}
}

func TestStateMachineRegisterAndChange(t *testing.T) {
	var sm StateMachine

	id1 := sm.RegisterState(State{})
	id2 := sm.RegisterState(State{})

	sm.ChangeState(id1)
	if sm.CurrentState() != id1 {
		t.Errorf("CurrentState() = %v, want %v", sm.CurrentState(), id1)
	}

	sm.ChangeState(id2)
	if sm.CurrentState() != id2 {
		t.Errorf("CurrentState() = %v, want %v", sm.CurrentState(), id2)
	}
}

func TestStateMachineChangeToSameStateIsNoOp(t *testing.T) {
	var sm StateMachine
	entered := 0

	id := sm.RegisterState(State{OnEnter: func(StateID) { entered++ }})
	sm.ChangeState(id)
	sm.ChangeState(id) // should be a no-op, no additional OnEnter call

	if entered != 1 {
		t.Errorf("OnEnter called %d times, want 1", entered)
	}
}

func TestStateMachineInvalidStateIgnored(t *testing.T) {
	var sm StateMachine
	id := sm.RegisterState(State{})
	sm.ChangeState(id)

	sm.ChangeState(StateID(999)) // invalid, should be ignored

	if sm.CurrentState() != id {
		t.Errorf("CurrentState() = %v, want %v (invalid change should be rejected)", sm.CurrentState(), id)
	}
}

func TestStateMachineCallbackOrder(t *testing.T) {
	var sm StateMachine
	var order []string

	idA := sm.RegisterState(State{
		OnLeave: func(StateID) { order = append(order, "A.OnLeave") },
	})
	idB := sm.RegisterState(State{
		OnEnter: func(StateID) { order = append(order, "B.OnEnter") },
	})
	sm.OnStateChange = func(previous, next StateID) { order = append(order, "OnStateChange") }

	sm.ChangeState(idA)
	order = nil // reset after initial transition into A

	sm.ChangeState(idB)

	want := []string{"A.OnLeave", "OnStateChange", "B.OnEnter"}
	if len(order) != len(want) {
		t.Fatalf("callback order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("callback order = %v, want %v", order, want)
			break
		}
	}
}

func TestStateMachineCallbackArguments(t *testing.T) {
	var sm StateMachine
	var leavePrev, enterPrev StateID
	var leaveNext StateID
	var changeFromArg, changeToArg StateID

	idA := sm.RegisterState(State{
		OnLeave: func(next StateID) { leaveNext = next },
	})
	idB := sm.RegisterState(State{
		OnEnter: func(previous StateID) { enterPrev = previous },
	})
	sm.OnStateChange = func(previous, next StateID) {
		changeFromArg, changeToArg = previous, next
	}

	sm.ChangeState(idA)
	sm.ChangeState(idB)

	_ = leavePrev
	if leaveNext != idB {
		t.Errorf("OnLeave next = %v, want %v", leaveNext, idB)
	}
	if enterPrev != idA {
		t.Errorf("OnEnter previous = %v, want %v", enterPrev, idA)
	}
	if changeFromArg != idA || changeToArg != idB {
		t.Errorf("OnStateChange(prev=%v, next=%v), want (prev=%v, next=%v)", changeFromArg, changeToArg, idA, idB)
	}
}

func TestStateMachineNestedChangeStateIgnored(t *testing.T) {
	var sm StateMachine
	nestedCallCompleted := false

	idA := sm.RegisterState(State{})
	idB := sm.RegisterState(State{
		OnEnter: func(StateID) {
			sm.ChangeState(idA) // nested change should be rejected
			nestedCallCompleted = true
		},
	})

	sm.ChangeState(idB)

	if !nestedCallCompleted {
		t.Fatal("expected the OnEnter callback to run to completion")
	}
	if sm.CurrentState() != idB {
		t.Errorf("CurrentState() = %v, want %v (nested ChangeState call should be rejected)", sm.CurrentState(), idB)
	}
}
