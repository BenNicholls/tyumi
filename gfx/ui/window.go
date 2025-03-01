package ui

import (
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// Window acts as a root node for the UI system.
type Window struct {
	Element

	labels             map[string]element
	blockingAnimations int // number of running animations blocking updates
}

func NewWindow(size vec.Dims, pos vec.Coord, depth int) (wnd *Window) {
	wnd = new(Window)
	wnd.Init(size, pos, depth)
	wnd.TreeNode.Init(wnd)
	wnd.labels = make(map[string]element)
	return
}

// Updates all visible subelements in the window, as well as all visible animations.
func (wnd *Window) Update() {
	// see how many animations (if any) are blocking updates
	wnd.blockingAnimations = 0
	util.WalkTree[element](wnd, func(element element) {
		for _, a := range element.getAnimations() {
			if a.IsBlocking() && a.IsPlaying() {
				wnd.blockingAnimations += 1
			}
		}
	}, ifVisible)

	// update all visible subelements (unless window is blocked) and their animations
	util.WalkSubTrees[element](wnd, func(element element) {
		if !wnd.IsBlocked() {
			element.Update()
		}
		element.updateAnimations()
	}, ifVisible)

	wnd.updateAnimations()
}

func (wnd *Window) Render() {
	//prepare_render
	util.WalkTree[element](wnd, func(element element) {
		element.prepareRender()
	}, ifVisible)

	//render all visible subnodes
	util.WalkSubTrees[element](wnd, func(element element) {
		element.drawChildren()
		if element.IsUpdated() || element.isRedrawing() {
			element.Render()
		}
		element.renderAnimations()
		element.finalizeRender() // does this need to go in a seperate walk??
	}, ifVisible)

	wnd.drawChildren()
	wnd.renderAnimations()
	wnd.finalizeRender()
}

func (wnd *Window) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	util.WalkSubTrees[element](wnd, func(element element) {
		if !event_handled {
			event_handled = element.HandleKeypress(key_event)
		}
	}, ifVisible)

	return
}

func (wnd *Window) IsBlocked() bool {
	return wnd.blockingAnimations > 0
}

// returns this window so subelements can find this. how a window would find a parent window remains a topic
// of active and fruitless discussion. thankfully i haven't thought about nesting windows yet so it doesn't keep
// me up at night
// THINK: i think dialogs are going to be subwindows? they might be substates though, so maybe nesting windows won't
// be a thing
func (wnd *Window) getWindow() *Window {
	return wnd
}

func (wnd *Window) addLabel(label string, e element) {
	if _, ok := wnd.labels[label]; ok {
		log.Warning("Duplicate label: ", label)
		return
	}

	wnd.labels[label] = e
}

func (wnd *Window) removeLabel(label string) {
	delete(wnd.labels, label)
}

func (wnd *Window) onSubNodeAdded(subNode element) {
	//find labelled subnodes of the new element and add them to the label map
	util.WalkTree(subNode, func(e element) {
		if e.IsLabelled() {
			wnd.addLabel(e.GetLabel(), e)
		}
	})
}

func (wnd *Window) onSubNodeRemoved(subNode element) {
	//find labelled subnodes of the removed element and remove them from the label map
	util.WalkTree(subNode, func(e element) {
		if e.IsLabelled() {
			wnd.removeLabel(e.GetLabel())
		}
	})
}

func (wnd *Window) onBlockingAnimationAdded() {
	wnd.blockingAnimations += 1
}

// Labelled elements can be retrieved via their label string from the window they are in. Also the labels can be
// used for other things that I have not thought up yet.
type Labelled interface {
	SetLabel(string)
	GetLabel() string
	IsLabelled() bool
}

// predicate for window's ui-tree-walking functions. elements that are not visible do not need to be updated
// or rendered, and neither do their children, so we use this to break early
func ifVisible(e element) bool {
	return e.IsVisible()
}
