package ui

import (
	"slices"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var EV_FOCUS_CHANGED = event.Register("Focus Changed", event.SIMPLE)

// Window acts as a root node for the UI system.
type Window struct {
	Element

	labels             map[string]element
	blockingAnimations int //number of running animations blocking updates

	SendKeyEventsToUnfocused bool //if true, unhandled key events will be sent to all elements, not just the focused one.
	focusedElement           element
	tabbingOrder             []element
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
	if key_event.Key == input.K_TAB && key_event.PressType == input.KEY_PRESSED && len(wnd.tabbingOrder) > 0 {
		wnd.TabForward()
		return
	}

	if wnd.SendKeyEventsToUnfocused {
		util.WalkSubTrees[element](wnd, func(element element) {
			if !event_handled {
				event_handled = element.HandleKeypress(key_event)
			}
		}, ifVisible)
	} else {
		if wnd.focusedElement != nil && wnd.focusedElement.IsVisible() {
			event_handled = wnd.focusedElement.HandleKeypress(key_event)
		}
	}

	return
}

// SetTabbingOrder sets the order for tabbing between elements. Any previously set tabbing order is not retained.
func (wnd *Window) SetTabbingOrder(tabbed_elements ...element) {
	wnd.tabbingOrder = nil

	for _, e := range tabbed_elements {
		if slices.Contains(wnd.tabbingOrder, e) {
			continue
		}

		if window := e.getWindow(); window == wnd {
			wnd.tabbingOrder = append(wnd.tabbingOrder, e)
		}
	}
}

// TabForward switches the window's currently focused element to the next element in the tabbing order (if there is
// one). If no element is focused, or if the currently focused element is not in the tabbing order, the first element
// in the tabbing order is focused instead.
func (wnd *Window) TabForward() {
	if len(wnd.tabbingOrder) == 0 {
		return
	}

	if wnd.focusedElement == nil {
		wnd.tabbingOrder[0].Focus()
	} else if !slices.Contains(wnd.tabbingOrder, wnd.focusedElement) {
		wnd.tabbingOrder[0].Focus()
	} else {
		i := slices.Index(wnd.tabbingOrder, wnd.focusedElement)
		wnd.tabbingOrder[util.CycleClamp(i+1, 0, len(wnd.tabbingOrder)-1)].Focus()
	}
}

// GetFocusedElementID retrieves the ID of the focused element. If no element is focused, returns 0
func (wnd *Window) GetFocusedElementID() ElementID {
	if wnd.focusedElement == nil {
		return ElementID(0)
	}

	return wnd.focusedElement.ID()
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

		if e.IsFocused() {
			if wnd.focusedElement == nil {
				wnd.focusedElement = e
			} else {
				log.Warning("Focused element added to window that already has another focused element.")
				e.Defocus()
			}
		}
	})
}

func (wnd *Window) onSubNodeRemoved(subNode element) {
	util.WalkTree(subNode, func(e element) {
		//find labelled subnodes of the removed element and remove them from the label map
		if e.IsLabelled() {
			wnd.removeLabel(e.GetLabel())
		}

		//remove subnode from tabbing order if necessary
		if len(wnd.tabbingOrder) > 0 {
			if i := slices.Index(wnd.tabbingOrder, e); i != -1 {
				wnd.tabbingOrder = slices.Delete(wnd.tabbingOrder, i, i+1)
			}
		}

		if wnd.focusedElement == e {
			wnd.focusedElement = nil
		}
	})
}

func (wnd *Window) onSubNodeFocused(subnode element) {
	if wnd.focusedElement != nil {
		wnd.focusedElement.Defocus()
	}

	wnd.focusedElement = subnode
	event.FireSimple(EV_FOCUS_CHANGED)
}

func (wnd *Window) onSubNodeDefocused(subnode element) {
	if wnd.focusedElement == subnode {
		wnd.focusedElement = nil
	}
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
