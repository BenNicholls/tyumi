package ui

import (
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// Window acts as a root node for the UI system.
type Window struct {
	ElementPrototype

	labels map[string]Element
}

func NewWindow(w, h int, pos vec.Coord, depth int) (wnd *Window) {
	wnd = new(Window)
	wnd.Init(w, h, pos, depth)
	wnd.TreeNode.Init(wnd)
	wnd.labels = make(map[string]Element)
	return
}

// Updates all visible subelements in the window, as well as all visible animations.
func (wnd *Window) Update() {
	util.WalkSubTrees[Element](wnd, func(element Element) {
		if element.IsVisible() {
			element.Update()
			element.updateAnimations()
		}
	})

	wnd.updateAnimations()
}

func (wnd *Window) Render() {
	//prepare_render
	util.WalkTree[Element](wnd, func(element Element) {
		if element.IsVisible() {
			element.prepareRender()
		}
	})

	//render all subnodes
	util.WalkSubTrees[Element](wnd, func(element Element) {
		if element.IsVisible() {
			element.drawChildren()
			if element.IsUpdated() || element.isRedrawing() {
				element.Render()
			}
			element.renderAnimations()
			element.finalizeRender() // does this need to go in a seperate walk??
		}
	})

	wnd.renderAnimations()
	wnd.drawChildren()
	wnd.finalizeRender()
}

func (wnd *Window) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	util.WalkSubTrees[Element](wnd, func(element Element) {
		if !event_handled && element.IsVisible() {
			event_handled = element.HandleKeypress(key_event)
		}
	})

	return
}

// returns this window so subelements can find this. how a window would find a parent window remains a topic
// of active and fruitless discussion. thankfully i haven't thought about nesting windows yet so it doesn't keep
// me up at night
// THINK: i think dialogs are going to be subwindows? they might be substates though, so maybe nesting windows won't
// be a thing
func (wnd *Window) getWindow() *Window {
	return wnd
}

func (wnd *Window) addLabel(label string, e Element) {
	if _, ok := wnd.labels[label]; ok {
		log.Warning("Duplicate label: ", label)
		return
	}

	wnd.labels[label] = e
}

func (wnd *Window) removeLabel(label string) {
	delete(wnd.labels, label)
}

func (wnd *Window) onSubNodeAdded(subNode Element) {
	//find labelled subnodes of the new element and add them to the label map
	util.WalkTree[Element](subNode, func(e Element) {
		if e.IsLabelled() {
			wnd.addLabel(e.GetLabel(), e)
		}
	})
}

func (wnd *Window) onSubNodeRemoved(subNode Element) {
	//find labelled subnodes of the removed element and remove them from the label map
	util.WalkTree[Element](subNode, func(e Element) {
		if e.IsLabelled() {
			wnd.removeLabel(e.GetLabel())
		}
	})
}

// Labelled elements can be retrieved via their label string from the window they are in. Also the labels can be
// used for other things that I have not thought up yet.
type Labelled interface {
	SetLabel(string)
	GetLabel() string
	IsLabelled() bool
}
