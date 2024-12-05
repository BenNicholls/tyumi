package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

//A List is a container that renders it's children elements from top to bottom, like you would expect
//a list to do.
type List struct {
	Container

	selected  int  //element that is currently selected. selected item will be ensured to be visible
	highlight bool //toggle to highlight currently selected list item

	contentHeight int //total height of all list contents. used for viewport and scrollbar purposes
	scrollOffset  int //number of rows (NOT elements) to scroll the list contents to keep selected item visible
}

func NewList(w, h int, pos vec.Coord, depth int) (l List) {
	l = List{
		Container: NewContainer(w, h, pos, depth),
	}

	return
}

func (l *List) AddElement(elems ...Element) {
	l.Container.AddElement(elems...)
	l.calibrate()
}

func (l *List) RemoveElement(e Element) {
	l.Container.RemoveElement(e)
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

//positions all the children elements so they are top to bottom, and the selected item is visible
func (l *List) calibrate() {
	l.contentHeight = 0
	for _, child := range l.children {
		child.MoveTo(vec.Coord{0, l.contentHeight-l.scrollOffset})
		l.contentHeight += child.Bounds().H
	}

	//if there is more list content than can be displayed at once, ensure selected item is shown via scrolling
	if l.contentHeight > l.Bounds().H {
		intersect := vec.FindIntersectionRect(l.children[l.selected], vec.Rect{vec.ZERO_COORD, l.Size()})
		if sh := l.children[l.selected].Bounds().H; intersect.H != sh {
			scrollDelta := 0

			sy := l.children[l.selected].Bounds().Y
			if sy < 0 {
				scrollDelta = sy
			} else if sy >= l.Bounds().H {
				scrollDelta = sy - l.Bounds().H + sh
			} else { // element is in list, but not fully visible
				scrollDelta = sh - intersect.H
			}

			for _, child := range l.children {
				child.Move(0, -scrollDelta)
			}

			l.scrollOffset += scrollDelta

		}
	} else { //if content fits in the list, no need to remember some old scroll offset
		l.scrollOffset = 0
	}

	l.border.UpdateScrollbar(l.contentHeight, l.scrollOffset)
}

//Toggles highlighting of currently selected item.
func (l *List) ToggleHighlight() {
	l.highlight = !l.highlight
}

func (l *List) Select(selection int) {
	if l.selected == util.Clamp(selection, 0, len(l.children)-1) {
		return
	}

	l.selected = util.Clamp(selection, 0, len(l.children)-1)
	l.calibrate()
}

//Selects the next item
func (l *List) Next() {
	if len(l.children) <= 1 {
		return
	}

	l.selected, _ = util.ModularClamp(l.selected+1, 0, len(l.children)-1)
	l.calibrate()
}

//Selects the previous item in the list
func (l *List) Prev() {
	if len(l.children) <= 1 {
		return
	}

	l.selected, _ = util.ModularClamp(l.selected-1, 0, len(l.children)-1)
	l.calibrate()
}

func (l *List) Render() {
	l.Container.Render()

	//render highlight for selected item.
	//TODO: different options for how the selected item is highlighted. currently just inverts the colours
	if l.highlight {
		area := l.children[l.selected].Bounds()
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
