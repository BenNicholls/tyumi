package ui

import (
	"slices"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var ACTION_LIST_NEXT = input.RegisterAction("Select Next List Element")
var ACTION_LIST_PREV = input.RegisterAction("Select Previous List Element")
var ACTION_LIST_SCROLLUP = input.RegisterAction("Scroll List Up 1 Row")
var ACTION_LIST_SCROLLDOWN = input.RegisterAction("Select List Down 1 Row")

func init() {
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_LIST_NEXT, input.K_DOWN)
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_LIST_PREV, input.K_UP)
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_LIST_SCROLLUP, input.K_PAGEUP)
	input.DefaultActionMap.AddSimpleKeyAction(ACTION_LIST_SCROLLDOWN, input.K_PAGEDOWN)
}

// A List is a container that renders it's children elements from top to bottom, like you would expect a list to do. If
// the size of the content is too large a scrollbar is activated, like magic.
type List struct {
	Element

	ReverseOrder      bool   // if true, inserted elements are displayed from most recent to least recent
	OnChangeSelection func() // callback triggered any time selection is changed, i.e. by scrolling

	selectionEnabled bool // toggle to allow user to have an item selected. selected items are always kept visible.
	selectionIndex   int  // index of element that is currently selected. selected item will be ensured to be visible
	highlight        bool // toggle for showing the currently selected element (TODO: different highlight modes)

	items      []element // items in the list.
	capacity   int       // max capacity for items in the List. if 0, no limit is enforced. If > 0, new items replace the oldest items.
	emptyLabel *Textbox  // text shown when list is empty

	padding       int // amount of padding added between list items
	contentHeight int // total height of all list contents. used for viewport and scrollbar purposes
	scrollOffset  int // number of rows (NOT elements) to scroll the list contents to keep selected item visible

	recalibrate bool // flag that indicates list elements need to be recalibrated before rendering
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
	l.selectionIndex = -1
}

// Insert adds elements to the list. Inserted elements will be added to the end of the list and automatically
// positioned.
func (l *List) Insert(items ...element) {
	if len(items) == 0 {
		return
	}

	if l.items == nil {
		l.items = make([]element, 0)
	}

	oldCount := l.Count()
	itemAdded := false
	for _, item := range items {
		if slices.Contains(l.items, item) {
			continue
		}

		if l.capacity > 0 && len(l.items) == l.capacity {
			l.RemoveAt(0)
		}

		l.items = append(l.items, item)
		l.AddChild(item)
		itemAdded = true
	}

	if !itemAdded {
		return
	}

	if oldCount == 0 && l.Count() > 0 {
		if l.selectionEnabled {
			l.Select(0)
		}

		if l.emptyLabel != nil {
			l.emptyLabel.Hide()
		}
	}

	l.recalibrate = true
	l.Updated = true
}

// InsertText adds simple textboxes to the list, one for each string passed.
func (l *List) InsertText(align Alignment, items ...string) {
	textBoxes := make([]element, 0)
	for _, item := range items {
		textBoxes = append(textBoxes, NewTextbox(vec.Dims{l.size.W, FIT_TEXT}, vec.ZERO_COORD, 0, item, align))
	}

	l.Insert(textBoxes...)
}

// Remove removes ui elements from the list, if present.
func (l *List) Remove(items ...element) {
	for _, item := range items {
		itemIndex := slices.Index(l.items, item)
		if itemIndex == -1 {
			continue
		}

		l.RemoveAt(itemIndex)
	}
}

// RemoveAt removes the item at the provided index. If the index is out of range, does nothing.
func (l *List) RemoveAt(index int) {
	if index >= l.Count() || index < 0 {
		return
	}

	l.RemoveChild(l.items[index])
	l.items = slices.Delete(l.items, index, index+1)
	if l.selectionEnabled && index <= l.selectionIndex {
		l.Select(l.selectionIndex - 1)
	}

	if l.emptyLabel != nil && l.Count() == 0 {
		l.emptyLabel.Show()
	}

	l.recalibrate = true
	l.Updated = true
}

// RemoveAll removes all list items.
func (l *List) RemoveAll() {
	if l.Count() == 0 {
		return
	}

	for _, item := range l.items {
		l.RemoveChild(item)
	}

	l.items = nil
	l.Select(-1)
	if l.emptyLabel != nil {
		l.emptyLabel.Show()
	}

	l.contentHeight = 0
	l.scrollOffset = 0
	l.Border.UpdateScrollbar(l.contentHeight, l.scrollOffset)
	l.Updated = true
}

// Count returns the number of items in the list.
func (l List) Count() int {
	return len(l.items)
}

// SetCapacity sets the maximum number of items that a list can hold. Defaults to 0, indicating no maximum limit. If
// an item is inserted while the list is at capacity, the oldest item is removed.
func (l *List) SetCapacity(cap int) {
	l.capacity = max(cap, 0)

	if l.capacity > 0 && len(l.items) > l.capacity {
		toTrim := len(l.items) - l.capacity
		for i := range toTrim {
			l.RemoveAt(toTrim - 1 - i)
		}
	}
}

// Sets the amount of padding between list items.
func (l *List) SetPadding(padding int) {
	if l.padding == padding {
		return
	}

	l.padding = padding
	l.recalibrate = true
}

// SetEmptyText creates a label that is shown in the list when the list contains no items.
func (l *List) SetEmptyText(text string) {
	if l.emptyLabel == nil {
		if text == "" {
			return
		}

		l.emptyLabel = NewTextbox(vec.Dims{l.size.W - 2, FIT_TEXT}, vec.ZERO_COORD, 1, text, ALIGN_CENTER)
		l.AddChild(l.emptyLabel)
		l.emptyLabel.Center()
	} else {
		if text != "" {
			l.emptyLabel.ChangeText(text)
			l.emptyLabel.Center()
		} else {
			l.RemoveChild(l.emptyLabel)
			l.emptyLabel = nil
		}
	}
}

// positions all the children elements so they are top to bottom, and the selected item is visible
func (l *List) calibrate() {
	l.contentHeight = 0
	for i := range l.Count() {
		var item element
		if l.ReverseOrder {
			item = l.items[l.Count()-1-i]
		} else {
			item = l.items[i]
		}
		item.MoveTo(vec.Coord{0, l.contentHeight - l.scrollOffset})
		if item.IsBordered() {
			item.Move(0, 1)
		}

		if item.Bounds().Intersects(l.DrawableArea()) {
			item.Show()
		} else {
			item.Hide()
		}
		l.contentHeight += item.Bounds().H + l.padding
	}

	l.contentHeight -= l.padding // remove the padding below the last item
	if l.selectionEnabled {
		l.updateScrollPosition()
	} else {
		l.Border.UpdateScrollbar(l.contentHeight, l.scrollOffset)
	}
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
				if child.Bounds().Intersects(l.DrawableArea()) {
					child.Show()
				} else {
					child.Hide()
				}
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
	if !l.selectionEnabled || l.highlight == highlight {
		return
	}

	l.highlight = highlight
	l.Updated = true
}

func (l *List) Select(selection int) {
	if !l.selectionEnabled {
		return
	}

	if l.Count() == 0 {
		if l.selectionIndex != -1 {
			l.selectionIndex = -1
			fireCallbacks(l.OnChangeSelection)
			l.Updated = true
		}
		return
	}

	new_selection := util.Clamp(selection, 0, l.Count()-1)
	if l.selectionIndex == new_selection {
		return
	}

	l.selectionIndex = new_selection
	l.updateScrollPosition()
	fireCallbacks(l.OnChangeSelection)
	l.Updated = true
}

func (l List) GetSelectionIndex() int {
	return l.selectionIndex
}

func (l *List) getSelected() element {
	if !l.selectionEnabled || l.Count() == 0 {
		return nil
	}

	return l.items[l.selectionIndex]
}

// Selects the next item in the list, wrapping back around to the top if necessary.
func (l *List) SelectNext() {
	if !l.selectionEnabled || l.Count() <= 1 {
		return
	}

	var nextIndex int
	if !l.ReverseOrder {
		nextIndex = util.CycleClamp(l.selectionIndex+1, 0, l.Count()-1)
	} else {
		nextIndex = util.CycleClamp(l.selectionIndex-1, 0, l.Count()-1)
	}
	l.Select(nextIndex)
}

// Selects the previous item in the list, wrapping back around to the bottom if necessary.
func (l *List) SelectPrev() {
	if !l.selectionEnabled || l.Count() <= 1 {
		return
	}

	var prevIndex int
	if !l.ReverseOrder {
		prevIndex = util.CycleClamp(l.selectionIndex-1, 0, l.Count()-1)
	} else {
		prevIndex = util.CycleClamp(l.selectionIndex+1, 0, l.Count()-1)
	}
	l.Select(prevIndex)
}

func (l *List) SelectTop() {
	if l.ReverseOrder {
		l.Select(l.Count() - 1)
	} else {
		l.Select(0)
	}
}

func (l *List) SelectBottom() {
	if !l.ReverseOrder {
		l.Select(l.Count() - 1)
	} else {
		l.Select(0)
	}
}

func (l *List) scrollTo(offset int) {
	if offset == l.scrollOffset {
		return
	}

	if l.contentHeight < l.size.H {
		return
	}

	l.scrollOffset = offset
	l.recalibrate = true
}

func (l *List) ScrollToTop() {
	l.scrollTo(0)
}

func (l *List) ScrollToBottom() {
	l.scrollTo(l.contentHeight - l.size.H + 1)
}

func (l *List) ScrollUp() {
	l.scrollTo(util.Clamp(l.scrollOffset-1, 0, l.contentHeight-l.size.H))
}

func (l *List) ScrollDown() {
	l.scrollTo(util.Clamp(l.scrollOffset+1, 0, l.contentHeight-l.size.H))
}

func (l *List) prepareRender() {
	if l.Updated {
		l.forceRedraw = true
	}

	if l.recalibrate {
		l.calibrate()
		l.recalibrate = false
	}

	l.Element.prepareRender()
}

func (l *List) Render() {
	//render highlight for selected item.
	//TODO: different options for how the selected item is highlighted. currently just inverts the colours
	if l.selectionEnabled && l.highlight && l.Count() > 0 {
		selected_area := l.getSelected().Bounds()
		highlight_area := vec.FindIntersectionRect(selected_area, l.DrawableArea())
		l.Canvas.DrawEffect(gfx.InvertEffect, highlight_area)
	}
}

func (l *List) HandleAction(action input.ActionID) (action_handled bool) {
	switch action {
	case ACTION_LIST_NEXT:
		if l.selectionEnabled {
			l.SelectNext()
		}
	case ACTION_LIST_PREV:
		if l.selectionEnabled {
			l.SelectPrev()
		}
	case ACTION_LIST_SCROLLUP:
		l.ScrollUp()
	case ACTION_LIST_SCROLLDOWN:
		l.ScrollDown()
	default:
		return false
	}

	return true
}
