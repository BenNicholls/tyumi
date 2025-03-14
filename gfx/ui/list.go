package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// A List is a container that renders it's children elements from top to bottom, like you would expect a list to do. If
// the size of the content is too large a scrollbar is activated, like magic.
type List struct {
	Element

	OnChangeSelection func() //callback triggered any time selection is changed, i.e. by scrolling

	padding   int  //amount of padding added between list items
	selected  int  //element that is currently selected. selected item will be ensured to be visible
	highlight bool //toggle to highlight currently selected list item

	contentHeight int //total height of all list contents. used for viewport and scrollbar purposes
	scrollOffset  int //number of rows (NOT elements) to scroll the list contents to keep selected item visible
}

func NewList(size vec.Dims, pos vec.Coord, depth int) (l *List) {
	l = new(List)
	l.Init(size, pos, depth)

	return
}

func (l *List) Init(size vec.Dims, pos vec.Coord, depth int) {
	l.Element.Init(size, pos, depth)
	l.TreeNode.Init(l)
	l.Border.EnableScrollbar(0, 0)
}

func (l *List) AddChild(child element) {
	l.Element.AddChild(child)
	l.calibrate()
	l.Updated = true
}

func (l *List) AddChildren(children ...element) {
	l.Element.AddChildren(children...)
	l.calibrate()
	l.Updated = true
}

func (l *List) RemoveChild(child element) {
	l.Element.RemoveChild(child)
	l.calibrate()
}

// AddTextItems adds simple textboxes to the list, one for each string passed.
func (l *List) AddTextItems(justify Justification, items ...string) {
	for _, item := range items {
		l.AddChild(NewTextbox(vec.Dims{l.size.W, FIT_TEXT}, vec.ZERO_COORD, 0, item, justify))
	}
}

// positions all the children elements so they are top to bottom, and the selected item is visible
func (l *List) calibrate() {
	l.contentHeight = 0
	for _, child := range l.GetChildren() {
		child.MoveTo(vec.Coord{0, l.contentHeight - l.scrollOffset})
		if child.IsBordered() {
			child.Move(0, 1)
		}
		l.contentHeight += child.Bounds().H + l.padding
	}

	l.contentHeight -= l.padding // remove the padding below the last item

	l.updateScrollPosition()
}

// if there is more list content than can be displayed at once, ensure selected item is shown via scrolling
func (l *List) updateScrollPosition() {
	if l.contentHeight > l.size.H {
		selected := l.getSelected()
		intersect := vec.FindIntersectionRect(selected, l.DrawableArea())
		sh := selected.Bounds().H
		if intersect.H != sh {
			scrollDelta := 0
			sy := selected.Bounds().Y

			if sy < 0 { // element above the list's draw area
				scrollDelta = sy
			} else if sy >= l.size.H { // element below the list's draw area
				scrollDelta = sy - l.size.H + sh
			} else { // element is in list, but not fully visible
				scrollDelta = sh - intersect.H
			}

			for _, child := range l.GetChildren() {
				child.Move(0, -scrollDelta)
			}

			l.scrollOffset += scrollDelta
		}
	} else { //if content fits in the list, no need to remember some old scroll offset
		l.scrollOffset = 0
	}

	l.Border.UpdateScrollbar(l.contentHeight, l.scrollOffset)
}

// Enables list element highlighting for the currently selected element.
func (l *List) EnableHighlight() {
	l.setHighlight(true)
}

// Disables list element highlighting for the currently selected element.
func (l *List) DisableHighlight() {
	l.setHighlight(false)
}

// Toggles highlighting of currently selected element.
func (l *List) ToggleHighlight() {
	l.setHighlight(!l.highlight)
}

func (l *List) setHighlight(highlight bool) {
	if l.highlight == highlight {
		return
	}

	l.highlight = highlight
	l.Updated = true
}

// Sets the amount of padding between list items.
func (l *List) SetPadding(padding int) {
	if l.padding == padding {
		return
	}

	l.padding = padding
	l.calibrate()
}

func (l *List) Select(selection int) {
	new_selection := util.Clamp(selection, 0, l.ChildCount()-1)
	if l.selected == new_selection {
		return
	}

	l.selected = new_selection
	l.updateScrollPosition()
	fireCallbacks(l.OnChangeSelection)
	l.Updated = true
}

func (l List) GetSelectionIndex() int {
	return l.selected
}

func (l *List) getSelected() element {
	return l.GetChildren()[l.selected]
}

// Selects the next item in the list, wrapping back around to the top if necessary.
func (l *List) Next() {
	if l.ChildCount() <= 1 {
		return
	}

	nextIndex := util.CycleClamp(l.selected+1, 0, l.ChildCount()-1)
	l.Select(nextIndex)
}

// Selects the previous item in the list, wrapping back around to the bottom if necessary.
func (l *List) Prev() {
	if l.ChildCount() <= 1 {
		return
	}

	prevIndex := util.CycleClamp(l.selected-1, 0, l.ChildCount()-1)
	l.Select(prevIndex)
}

func (l *List) prepareRender() {
	if l.Updated {
		l.forceRedraw = true
	}

	l.Element.prepareRender()
}

func (l *List) Render() {
	//render highlight for selected item.
	//TODO: different options for how the selected item is highlighted. currently just inverts the colours
	if l.highlight {
		selected_area := l.getSelected().Bounds()
		highlight_area := vec.FindIntersectionRect(selected_area, l.DrawableArea())
		l.Canvas.DrawEffect(gfx.InvertEffect, highlight_area)
	}
}

func (l *List) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.PressType == input.KEY_RELEASED {
		return
	}

	switch key_event.Key {
	case input.K_UP, input.K_PAGEUP:
		l.Prev()
		event_handled = true
	case input.K_DOWN, input.K_PAGEDOWN:
		l.Next()
		event_handled = true
	}

	return
}
