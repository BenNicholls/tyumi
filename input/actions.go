package input

import (
	"slices"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/util"
)

var EV_ACTION int = event.Register("Action Event", event.COMPLEX)

// The default action map. Tyumi packages that rely on actions will register action triggers here.
var DefaultActionMap ActionMap

func init() {
	DefaultActionMap.keyTriggers = make(map[Keycode]*util.Set[ActionKeyTrigger])
}

var actionIDcounter int
var actionNames []string

type ActionID int

func (a ActionID) String() string {
	return actionNames[int(a)]
}

func RegisterAction(name string) ActionID {
	if index := slices.Index(actionNames, name); index == -1 {
		actionNames = append(actionNames, name)
		actionIDcounter++
		return ActionID(actionIDcounter - 1)
	} else {
		return ActionID(index)
	}
}

// ActionHandler is any function that takes an ActionID and handles the action. If the action is handled, it returns
// true
type ActionHandler func(action ActionID) (action_handled bool)

type ActionEvent struct {
	event.EventPrototype

	Action ActionID
}

func fireActionEvent(action ActionID) {
	if DefaultActionMap.disabled.Contains(action) {
		return
	}

	event.Fire(&ActionEvent{
		EventPrototype: event.New(EV_ACTION),
		Action:         action,
	})
}

type ActionKeyTrigger struct {
	Action ActionID // ID of the action to trigger
	Key    Keycode  // key to trigger on

	// presstype to trigger on. you can use KEYPRESS_EITHER to have the action trigger on both PRESS and RELEASE.
	// defaults to KEY_PRESSED
	PressType KeyPressType
}

// TriggeredBy returns true is the provided keyboard event successfully triggers the trigger. Trigger.
func (akt ActionKeyTrigger) TriggeredBy(key_event KeyboardEvent) bool {
	if akt.PressType != key_event.PressType && akt.PressType != KEYPRESS_EITHER {
		return false
	}

	return true
}

type ActionMap struct {
	keyTriggers map[Keycode]*util.Set[ActionKeyTrigger]

	disabled util.Set[ActionID] // actions that do not produce action events
}

// Adds triggers for the provided action for each provided key. These triggers default to firing on KEY_PRESSED with
// no regard to which modifier keys are held down.
func (am *ActionMap) AddSimpleKeyAction(action ActionID, keys ...Keycode) {
	for _, key := range keys {
		trigger := ActionKeyTrigger{Action: action, Key: key}
		am.AddKeyAction(trigger)
	}
}

// Adds a new key trigger for an action to the action map. The trigger defines the conditions under which the action
// fires.
func (am *ActionMap) AddKeyAction(trigger ActionKeyTrigger) {
	if _, ok := am.keyTriggers[trigger.Key]; !ok {
		am.keyTriggers[trigger.Key] = new(util.Set[ActionKeyTrigger])
	}

	am.keyTriggers[trigger.Key].Add(trigger)
}

// Disable actions from firing, even when their triggers are detected. The triggers are not removed from the actionmap,
// and can be re-enabled using EnableAction
func (am *ActionMap) DisableAction(actions ...ActionID) {
	for _, action := range actions {
		am.disabled.Add(action)
	}
}

// Enable actions to fire when their triggers are detected. Since actions default to being enabled, this function should
// only be used to re-enable actions that you've disabled temporarily for whatever reason.
func (am *ActionMap) EnableAction(actions ...ActionID) {
	for _, action := range actions {
		am.disabled.Remove(action)
	}
}

// EnableAll clears the list of disabled actions.
func (am *ActionMap) EnableAll() {
	am.disabled.RemoveAll()
}
