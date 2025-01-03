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
}

func NewList(w, h int, pos vec.Coord, depth int) (l *List) {
	l = new(List)
	l.ElementPrototype.Init(w, h, pos, depth)
	return
}

func (l *List) AddChild(elem Element) {
	l.ElementPrototype.AddChild(elem)
	l.calibrate()
	l.updated = true
}

func (l *List) AddChildren(elems ...Element) {
	l.ElementPrototype.AddChildren(elems...)
	l.calibrate()
	l.updated = true
}

func (l *List) RemoveChild(e Element) {
	l.ElementPrototype.RemoveChild(e)
	l.calibrate()
}

func (l *List) ToggleScrollbar() {
	if l.border.scrollbar {
		l.border.scrollbar = false
		l.border.dirty = true
	} else {
		l.border.EnableScrollbar(l.contentHeight, l.scrollOffset)
	}
}

// positions all the children elements so they are top to bottom, and the selected item is visible
func (l *List) calibrate() {
	l.contentHeight = 0
	for _, child := range l.GetChildren() {
		child.MoveTo(vec.Coord{0, l.contentHeight - l.scrollOffset})
		l.contentHeight += child.Bounds().H
	}

	//if there is more list content than can be displayed at once, ensure selected item is shown via scrolling
	if l.contentHeight > l.Bounds().H {

		intersect := vec.FindIntersectionRect(l.GetChildren()[l.selected], vec.Rect{vec.ZERO_COORD, l.Size()})
		if sh := l.GetChildren()[l.selected].Bounds().H; intersect.H != sh {
			scrollDelta := 0

			sy := l.GetChildren()[l.selected].Bounds().Y
			if sy < 0 {
				scrollDelta = sy
			} else if sy >= l.Bounds().H {
				scrollDelta = sy - l.Bounds().H + sh
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

	l.border.UpdateScrollbar(l.contentHeight, l.scrollOffset)
}

// Toggles highlighting of currently selected item.
func (l *List) ToggleHighlight() {
	l.highlight = !l.highlight
}

func (l *List) Select(selection int) {
	if l.selected == util.Clamp(selection, 0, len(l.GetChildren())-1) {
		return
	}

	l.selected = util.Clamp(selection, 0, l.ChildCount()-1)
	l.calibrate()
}

// Selects the next item
func (l *List) Next() {
	if l.ChildCount() <= 1 {
		return
	}

	l.selected = util.CycleClamp(l.selected+1, 0, l.ChildCount()-1)
	l.calibrate()
	l.updated = true
}

// Selects the previous item in the list
func (l *List) Prev() {
	if l.ChildCount() <= 1 {
		return
	}

	l.selected = util.CycleClamp(l.selected-1, 0, l.ChildCount()-1)
	l.calibrate()
	l.updated = true
}

func (l *List) Render() {
	if l.updated {
		l.forceRedraw = true
		l.updated = false
	}

	l.ElementPrototype.Render()

	//render highlight for selected item.
	//TODO: different options for how the selected item is highlighted. currently just inverts the colours
	if l.highlight {
		area := l.GetChildren()[l.selected].Bounds()
		l.Canvas.DrawEffect(gfx.InvertEffect, area)
	}
}

func (l *List) HandleKeypress(e input.KeyboardEvent) {
	switch e.Key {
	case input.K_UP, input.K_PAGEUP:
		l.Prev()
	case input.K_DOWN, input.K_PAGEDOWN:
		l.Next()
	}
}
