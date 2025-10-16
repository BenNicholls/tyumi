package tyumi

import (
	"slices"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
)

var currentScene scene
var dialogs []dialog

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
	mainConsole.AddChild(s.Window())
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

// Opens a dialog in the current scene.
func OpenDialog(d dialog) {
	if !d.Ready() {
		log.Error("Could not open dialog: dialog not initialized.")
		return
	}

	if slices.Contains(dialogs, d) {
		log.Error("Could not open dialog: dialog already open!")
		return
	}

	d.open()

	// disable input events for the active scene/dialog
	if len(dialogs) == 0 {
		currentScene.InputEvents().DisableListening()
	} else {
		dialogs[len(dialogs)-1].InputEvents().DisableListening()
	}

	d.Window().SetDepth(len(dialogs) + 1)
	mainConsole.AddChild(d.Window())

	dialogs = append(dialogs, d)
}

func closeDialog(d dialog) {
	idx := slices.Index(dialogs, d)
	if idx == -1 || !d.Ready() {
		log.Error("Cannot close dialog, dialog not open.")
		return
	}

	// if the dialog being closed is the top-most dialog, re-enable inputs for the next dialog/scene down in the hierarchy.
	if idx == len(dialogs)-1 {
		if idx == 0 {
			currentScene.InputEvents().EnableListening()
		} else {
			dialogs[idx-1].InputEvents().EnableListening()
		}
	}

	d.close()
	d.Shutdown()
	d.cleanup()

	dialogs = util.DeleteElement(dialogs, d)
}

func closeAllDialogs() {
	for _, d := range slices.Backward(dialogs) {
		closeDialog(d)
	}
}
