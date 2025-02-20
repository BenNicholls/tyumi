package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// A List is a container that renders it's children elements from top to bottom, like you would expect
// a list to do.
type List struct {
	ElementPrototype

	selected  int  //element that is currently selected. selected item will be ensured to be visible
	highlight bool //toggle to highlight currently selected list item

	contentHeight int //total height of all list contents. used for viewport and scrollbar purposes
	scrollOffset  int //number of rows (NOT elements) to scroll the list contents to keep selected item visible
	padding       int //amount of padding added between list items
}

func NewList(size vec.Dims, pos vec.Coord, depth int) (l *List) {
	l = new(List)
	l.ElementPrototype.Init(size, pos, depth)
	l.Border.EnableScrollbar(0, 0)
	return
}

func (l *List) AddChild(elem Element) {
	l.ElementPrototype.AddChild(elem)
	l.calibrate()
	l.Updated = true
}

func (l *List) AddChildren(elems ...Element) {
	l.ElementPrototype.AddChildren(elems...)
	l.calibrate()
	l.Updated = true
}

func (l *List) RemoveChild(e Element) {
	l.ElementPrototype.RemoveChild(e)
	l.calibrate()
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

// Toggles highlighting of currently selected item.
func (l *List) ToggleHighlight() {
	l.highlight = !l.highlight
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
	if l.selected == util.Clamp(selection, 0, l.ChildCount()-1) {
		return
	}

	l.selected = util.Clamp(selection, 0, l.ChildCount()-1)
	l.updateScrollPosition()
}

func (l List) GetSelectionIndex() int {
	return l.selected
}

func (l *List) getSelected() Element {
	return l.GetChildren()[l.selected]
}

// Selects the next item
func (l *List) Next() {
	if l.ChildCount() <= 1 {
		return
	}

	l.selected = util.CycleClamp(l.selected+1, 0, l.ChildCount()-1)
	l.updateScrollPosition()
	l.Updated = true
}

// Selects the previous item in the list
func (l *List) Prev() {
	if l.ChildCount() <= 1 {
		return
	}

	l.selected = util.CycleClamp(l.selected-1, 0, l.ChildCount()-1)
	l.updateScrollPosition()
	l.Updated = true
}

func (l *List) prepareRender() {
	if l.Updated {
		l.forceRedraw = true
	}

	l.ElementPrototype.prepareRender()
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

func (l *List) HandleKeypress(event *input.KeyboardEvent) (event_handled bool) {
	if event.PressType == input.KEY_RELEASED {
		return
	}

	switch event.Key {
	case input.K_UP, input.K_PAGEUP:
		l.Prev()
		event_handled = true
	case input.K_DOWN, input.K_PAGEDOWN:
		l.Next()
		event_handled = true
	}

	return
}
