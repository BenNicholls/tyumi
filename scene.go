package tyumi

import (
	"slices"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var EV_CHANGESCENE = event.Register("Scene Change Event")
var currentScene scene
var activeScene scene // the scene being updated, checked each frame

const (
	FIT_CONSOLE int = 0 //window size flag
)

// A gameobject to be handled by Tyumi's main game loop.
type scene interface {
	event.Listener
	InputEvents() *event.Stream

	Window() *ui.Window
	Ready() bool
	Shutdown()

	Update()
	UpdateUI()
	processTimers()
	IsBlocked() bool

	OpenDialog(dialog)

	getActiveSubScene() scene
	flushInputs()
	cleanup()
}

// Scene is the base implementation for Tyumi's game state object. Scenes contain a window, where the programs UI is
// built, as well as the machinery for handling game and input events. Custom scenes can be made by embedding this and
// overriding the virtual functions defined here. Most important is the Update() function, which runs once per-tick
// and should contain your main game code.
type Scene struct {
	util.StateMachine
	event.Stream // event stream for game events. these are collected during update() and then processed at the end of the tick (before ui updating/rendering)

	window *ui.Window

	subScene dialog
	timers   []Timer

	inputEvents          event.Stream        //for input events. processed at the start of each tick
	inputHandler         event.Handler       //user-provided input handling function. runs AFTER the UI has had a chance to process input.
	actionHandler        input.ActionHandler //user-provided action handling function. runs AFTER the UI has had a chance to process input.
	keypressInputHandler func(key_event *input.KeyboardEvent) bool

	ready bool // indicates the scene has been successfully initialized
}

// Init prepares the scene, defaulting to a window the full size of the console.
// NOTE: If you want a border drawn around the window, use InitBordered() instead since Tyumi draws borders *around*
// objects and if the window is the size of the console you wouldn't be able to see it.
func (s *Scene) Init() {
	if !mainConsole.ready {
		log.Error("Cannot fit scene window to console: console not initialized.")
		return
	}

	s.init(mainConsole.Size(), vec.ZERO_COORD, false)
}

// InitBordered prepares the scene, defaulting to a window the full size of the console with a border drawn around
// the outside.
func (s *Scene) InitBordered() {
	if !mainConsole.ready {
		log.Error("Cannot fit scene window to console: console not initialized.")
		return
	}

	s.init(mainConsole.Size().Shrink(2, 2), vec.Coord{1, 1}, true)
}

// InitCentered prepares the scene, centering the scene's window inside the console.
func (s *Scene) InitCentered(size vec.Dims) {
	if !mainConsole.ready {
		log.Error("Cannot fit scene window to console: console not initialized.")
		return
	}

	pos := vec.Coord{(mainConsole.Size().W - size.W) / 2, (mainConsole.Size().H - size.H) / 2}
	s.init(size, pos, false)
}

// InitCustom prepares the scene, creating a window using the given size and position. The position will be relative
// to the console. If the console has been initialized, you can use FIT_CONSOLE as the width and/or height to have the
// scene size itself automatically.
func (s *Scene) InitCustom(size vec.Dims, pos vec.Coord) {
	s.init(size, pos, false)
}

func (s *Scene) init(size vec.Dims, pos vec.Coord, bordered bool) {
	if s.ready {
		log.Error("Trying to initialize a scene more than once. Don't do that.")
		return
	}

	if size.W == FIT_CONSOLE || size.H == FIT_CONSOLE {
		if !mainConsole.ready {
			log.Error("Cannot fit scene window to console: console not initialized.")
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

	s.inputEvents = event.NewStream(100, s.handleInput)

	//setup automatic listening for input events.
	s.inputEvents.Listen(input.EV_ACTION, input.EV_KEYBOARD, input.EV_MOUSEBUTTON, input.EV_MOUSEMOVE)
	s.timers = make([]Timer, 0)
	s.ready = true
}

// Update is run each tick, after input has been handled and before UI is updated/rendered. Override this function
// with your primary game code!
func (s *Scene) Update() {}

// UpdateUI is called before each frame is rendered, allowing you to apply ui changes for rendering all at once if
// you prefer. Otherwise you can implement Update() functions on the individual UI elements themselves and have them
// control their own behaviour.
func (s *Scene) UpdateUI() {}

func (s *Scene) OpenDialog(subScene dialog) {
	if !subScene.Ready() {
		log.Error("Could not open dialog, dialog not initialized.")
		return
	}

	s.subScene = subScene
}

func (s *Scene) getActiveSubScene() scene {
	if s.subScene == nil {
		return nil
	}

	if s.subScene.Done() {
		s.subScene.Shutdown()
		s.subScene.cleanup()
		s.subScene = nil
		return nil
	} else {
		if subScene := s.subScene.getActiveSubScene(); subScene != nil {
			return subScene
		} else {
			return s.subScene
		}
	}
}

func (s *Scene) cleanup() {
	s.DisableListening()
	s.inputEvents.DisableListening()

	if s.subScene != nil {
		s.subScene.Shutdown()
		s.subScene.cleanup()
		s.subScene = nil
	}

	// in theory scenes should be freed from memory after being shutdown so this is pointless, but on the off chance a
	// reference is hanging around maybe this will help catch a bug.
	s.ready = false
}

// Shutdown is called when the scene is no longer needed and should cleanly pack itself away (for example, when
// switching to another scene or closing the program). Override this function and use it to free resources, save things
// to disk, whatever you need to do.
func (s *Scene) Shutdown() {}

func (s *Scene) Window() *ui.Window {
	return s.window
}

func (s *Scene) InputEvents() *event.Stream {
	return &s.inputEvents
}

// Sets the function for handling inputs to the scene object. Inputs are collected, distributed and then
// processed at the beginning of each tick(). This handler is called after the UI and any more specific input handlers
// have had a chance to handle the input. If another handler handles the event then event.Handled() will be true. You
// can still choose to ignore that and handle the event again if you like though.
func (s *Scene) SetInputHandler(handler event.Handler) {
	s.inputHandler = handler
}

// Sets the function for handling keypresses. Inputs are collected, distributed and then processed at the beginning of
// each tick(). This handler is called only for key press events (not key releases) after the UI has had a chance to
// handle the input. If the UI handles the event then event.Handled() will be true. You can still choose to ignore that
// and handle the event again if you like though.
func (s *Scene) SetKeypressHandler(keypress_handler func(keyboard_event *input.KeyboardEvent) bool) {
	s.keypressInputHandler = keypress_handler
}

// Sets the function for handling action events. Inputs are collected, distributed and then processed at the beginning of
// each tick(). This handler is called only for events that trigger actions. It runs after the UI has had a chance to
// handle the action. If the UI handles the action then event.Handled() will be true. You can still choose to ignore that
// and handle the action again if you like though.
func (s *Scene) SetActionHandler(action_handler input.ActionHandler) {
	s.actionHandler = action_handler
}

func (s *Scene) handleInput(event event.Event) (event_handled bool) {
	switch event.ID() {
	case input.EV_ACTION:
		action_event := event.(*input.ActionEvent)
		event_handled = s.window.HandleAction(action_event.Action)
		if s.actionHandler != nil {
			event_handled = s.actionHandler(action_event.Action) || event_handled
		}
	case input.EV_KEYBOARD:
		key_event := event.(*input.KeyboardEvent)
		event_handled = s.window.HandleKeypress(key_event)
		if s.keypressInputHandler != nil && key_event.PressType == input.KEY_PRESSED {
			event_handled = s.keypressInputHandler(key_event) || event_handled
		}
	}

	if s.inputHandler != nil {
		event_handled = s.inputHandler(event) || event_handled
	}

	return
}

func (s *Scene) flushInputs() {
	s.inputEvents.FlushEvents()

	if s.subScene != nil {
		s.subScene.flushInputs()
	}
}

func (s *Scene) Ready() bool {
	return s.ready
}

// Returns true if updating has been blocked. Currently this only happens from blocking animations, but later might
// also indicate that the game is paused perhaps.
func (s *Scene) IsBlocked() bool {
	return s.window.IsBlocked()
}

// SetInitialScene sets a scene to be run by Tyumi at the beginning of execution.
// This function DOES NOTHING if a scene has already been initialized.
func SetInitialScene(s scene) {
	if currentScene != nil {
		return
	}

	if !mainConsole.ready {
		log.Error("Cannot set initial scene: console not initialized. Run InitConsole() first.")
		return
	}

	if s == nil || !s.Ready() {
		log.Error("Cannot set initial scene: scene not initialized or ready.")
		return
	}

	currentScene = s
}

// CreateTimer creates a timer. After duration ticks, the function f is run and the timer is destroyed.
func (s *Scene) CreateTimer(duration int, f func()) {
	if f == nil || duration <= 0 {
		return
	}

	s.timers = append(s.timers, Timer{TimerFunction: f, Ticks: duration})
}

func (s *Scene) processTimers() {
	if len(s.timers) == 0 {
		return
	}

	for i := range s.timers {
		s.timers[i].Process()
	}

	s.timers = slices.DeleteFunc(s.timers, func(timer Timer) bool {
		return timer.Done()
	})
}

type SceneChangeEvent struct {
	event.EventPrototype

	newScene scene
}

// ChangeScene changes the current scene being run in Tyumi's gameloop. The change is done at the end of the current
// engine tick. The old scene's Shutdown() method is called before we swap in the new one. Be sure to initialize the
// new scene before calling ChangeScene(), otherwise no change will happen and the old scene will remain.
func ChangeScene(new_scene scene) {
	if new_scene == nil || !new_scene.Ready() {
		log.Error("Could not change scene: scene invalid or not initialized.")
		return
	}

	//if user tries to use this to setup the initial scene, just forgive them their sin and do it. no need to
	//harass them with "the correct way".
	if currentScene == nil {
		SetInitialScene(new_scene)
		return
	}

	event.Fire(EV_CHANGESCENE, &SceneChangeEvent{newScene: new_scene})
}
