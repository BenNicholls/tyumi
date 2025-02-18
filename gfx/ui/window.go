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

	labels              map[string]Element
	blocking_animations int // number of running animations blocking updates
}

func NewWindow(size vec.Dims, pos vec.Coord, depth int) (wnd *Window) {
	wnd = new(Window)
	wnd.Init(size, pos, depth)
	wnd.TreeNode.Init(wnd)
	wnd.labels = make(map[string]Element)
	return
}

// Updates all visible subelements in the window, as well as all visible animations.
func (wnd *Window) Update() {
	// see how many animations (if any) are blocking updates
	wnd.blocking_animations = 0
	util.WalkTree[Element](wnd, func(element Element) {
		if element.IsVisible() {
			for _, a := range element.getAnimations() {
				if a.IsBlocking() && a.IsPlaying() {
					wnd.blocking_animations += 1
				}
			}
		}
	})

	util.WalkSubTrees[Element](wnd, func(element Element) {
		if element.IsVisible() {
			if !wnd.IsBlocked() {
				element.Update()
			}
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

	wnd.drawChildren()
	wnd.renderAnimations()
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

func (wnd *Window) IsBlocked() bool {
	return wnd.blocking_animations > 0
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

func (wnd *Window) onBlockingAnimationAdded() {
	wnd.blocking_animations += 1
}

// Labelled elements can be retrieved via their label string from the window they are in. Also the labels can be
// used for other things that I have not thought up yet.
type Labelled interface {
	SetLabel(string)
	GetLabel() string
	IsLabelled() bool
}
