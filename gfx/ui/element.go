package ui

import (
	"slices"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// An element is the base structure of anything handled by the UI system.
type Element interface {
	vec.Bounded
	util.TreeType[Element]
	Labelled

	Update()
	updateAnimations()

	prepareRender()
	Render()
	renderAnimations()
	finalizeRender()
	drawToParent()
	drawChildren()
	ForceRedraw() //Force the element to clear and redraw itself and all children from scratch
	isRedrawing() bool

	HandleKeypress(*input.KeyboardEvent) (event_handled bool)

	MoveTo(vec.Coord)
	Move(int, int)

	IsVisible() bool
	IsUpdated() bool
	IsBordered() bool
	getBorder() *Border
	getCanvas() *gfx.Canvas
	getWindow() *Window
}

type ElementPrototype struct {
	gfx.Canvas
	util.TreeNode[Element]
	Updated bool //indicates this object's state has changed and needs to be re-rendered.

	position    vec.Coord
	depth       int //depth for the UI system, relative to the element's parent. if no parent, relative to the console
	visible     bool
	bordered    bool
	forceRedraw bool   //indicates this object needs to clear and render everything from zero
	label       string // an optional identifier for the element
	border      *Border
	animations  []gfx.Animator //animations on this element. these are updated once per frame.
}

func (e *ElementPrototype) Init(w, h int, pos vec.Coord, depth int) {
	e.Canvas.Init(w, h)
	e.SetDefaultVisuals(defaultCanvasVisuals)
	e.position = pos
	e.depth = depth
	e.visible = true
	e.Updated = true
	e.TreeNode.Init(e)
}

// Resizes the element. This clears the internal canvas and forces redraws of everything.
func (e *ElementPrototype) Resize(size vec.Dims) {
	if size == e.Size() {
		return
	}

	e.Canvas.Resize(size.W, size.H)

	if e.border != nil {
		e.border.resize(size)
	}

	e.Updated = true
	e.forceRedraw = true
	e.forceParentRedraw()
}

func (e *ElementPrototype) SetDefaultColours(colours col.Pair) {
	e.Canvas.SetDefaultColours(colours)
	if e.border != nil {
		e.border.setColours(colours)
	}
	e.Updated = true
}

// Enable the border. Defaults to ui.DefaultBorderStyle. Use SetBorderStyle to use something else. If no border
// has been setup via SetupBorder(), a default one will be created.
func (e *ElementPrototype) EnableBorder() {
	e.setBorder(true)
}

func (e *ElementPrototype) DisableBorder() {
	e.setBorder(false)
}

func (e *ElementPrototype) setBorder(bordered bool) {
	if bordered == e.bordered {
		return
	}

	e.bordered = bordered
	if e.bordered == true && e.border == nil {
		//setup default border if none has been setup
		e.SetupBorder("", "")
	}

	e.forceParentRedraw()
}

func (e *ElementPrototype) IsBordered() bool {
	return e.bordered
}

// Creates and enables a border for the element. Title will be shown in the top left, and hint will be shown in the
// bottom right.
// TODO: centered titles? setting borderstyle at the same time?
func (e *ElementPrototype) SetupBorder(title, hint string) {
	e.border = NewBorder(e.Size())
	e.border.title = title
	e.border.hint = hint
	e.border.setColours(e.DefaultColours())
	e.border.style = &DefaultBorderStyle
	e.EnableBorder()
}

// Sets the border style flag and, if possible, updates the used style. Sometimes you can't though...
// for example, setting the flag to BORDER_INHERIT while the element does not have a parent.
func (e *ElementPrototype) SetBorderStyle(styleFlag borderStyleFlag, borderStyle ...BorderStyle) {
	if e.border == nil {
		log.Error("UI: Could not apply borderstyle, element has no border to style.")
		return
	}

	switch styleFlag {
	case BORDER_STYLE_DEFAULT:
		e.border.style = &DefaultBorderStyle
	case BORDER_STYLE_INHERIT:
		if parent := e.GetParent(); parent != nil {
			if parent_border := parent.getBorder(); parent_border != nil {
				e.border.style = parent_border.style
			}
		}
	case BORDER_STYLE_CUSTOM:
		if borderStyle == nil {
			log.Error("Custom border style application failed: no borderstyle provided.")
			return
		}

		e.border.style = &borderStyle[0]
	}

	e.border.styleFlag = styleFlag
}

func (e *ElementPrototype) getBorder() *Border {
	return e.border
}

func (e *ElementPrototype) getWindow() *Window {
	parent := e.GetParent()
	if parent == nil {
		//if this element *is* a window, we can find out by grabbing the self pointer from the internal
		//treenode and try casting it.
		if wnd, ok := e.GetSelf().(*Window); ok {
			return wnd
		}
		return nil
	}

	return parent.getWindow()
}

func (e *ElementPrototype) Bounds() vec.Rect {
	return vec.Rect{e.position, e.Size()}
}

func (e *ElementPrototype) MoveTo(pos vec.Coord) {
	if e.position == pos {
		return
	}

	e.position = pos
	e.forceParentRedraw()
}

// THINK: should this take a coord too? or a Vec2i?
func (e *ElementPrototype) Move(dx, dy int) {
	e.MoveTo(vec.Coord{e.position.X + dx, e.position.Y + dy})
}

func (e *ElementPrototype) AddChild(child Element) {
	if child_border := child.getBorder(); child_border != nil && child_border.styleFlag == BORDER_STYLE_INHERIT {
		child_border.style = e.border.style
	}

	e.TreeNode.AddChild(child)
	if window := e.getWindow(); window != nil {
		window.onSubNodeAdded(child)
	}
	e.ForceRedraw()
}

func (e *ElementPrototype) AddChildren(children ...Element) {
	for _, child := range children {
		e.AddChild(child)
	}
}

func (e *ElementPrototype) RemoveChild(child Element) {
	e.TreeNode.RemoveChild(child)
	if window := e.getWindow(); window != nil {
		window.onSubNodeRemoved(child)
	}
	e.ForceRedraw()
}

// Update() can be overriden to update the state of the UI Element. Note that the element's animations are updated
// separately and do not need to be managed here.
func (e *ElementPrototype) Update() {
	return
}

func (e *ElementPrototype) IsUpdated() bool {
	return e.Updated
}

func (e *ElementPrototype) updateAnimations() {
	for _, a := range e.animations {
		if a.Playing() {
			a.Update()
		}
	}

	// remove finished one-shot animations
	e.animations = slices.DeleteFunc[[]gfx.Animator](e.animations, func(a gfx.Animator) bool {
		return a.IsOneShot() && a.Done()
	})
}

func (e *ElementPrototype) ForceRedraw() {
	e.forceRedraw = true
}

func (e *ElementPrototype) isRedrawing() bool {
	return e.forceRedraw
}

func (e *ElementPrototype) forceParentRedraw() {
	parent := e.GetParent()
	if parent != nil {
		parent.ForceRedraw()
	}
}

// performs some pre-render checks. done for the whole tree before any rendering is done.
func (e *ElementPrototype) prepareRender() {
	if !e.visible {
		return
	}

	if e.forceRedraw {
		if e.bordered {
			e.border.dirty = true
		}

		for _, child := range e.GetChildren() { //make sure siblings recompute border links
			if child.IsBordered() {
				child.getBorder().dirty = true
			}
		}

		e.Clear()
	}
}

// Renders any changes in the element to the internal canvas. Override this to implement custom rendering behaviour.
// Note that this is called *after* any subelements are drawn to the canvas, and *before* any running animations
// are rendered.
func (e *ElementPrototype) Render() {
	return
}

// performs some after-render cleanups. TODO: could also put some profiling code in here once that's a thing?
func (e *ElementPrototype) finalizeRender() {
	e.Updated = false
	e.forceRedraw = false
}

func (e *ElementPrototype) renderAnimations() {
	for _, animation := range e.animations {
		if animation.Playing() && vec.Intersects(e.getCanvas(), animation) {
			animation.Render(&e.Canvas)
		}
	}
}

func (e *ElementPrototype) drawToParent() {
	parent := e.GetParent()
	if parent == nil {
		return
	}

	e.Draw(parent.getCanvas(), e.position, e.depth)
	if e.bordered {
		e.border.DrawToCanvas(parent.getCanvas(), e.position, e.depth)
	}
}

func (e *ElementPrototype) drawChildren() {
	for _, child := range e.GetChildren() {
		if child.IsVisible() {
			if child.getCanvas().Dirty() || e.forceRedraw {
				child.drawToParent()
			}
		}
	}
}

func (e *ElementPrototype) HandleKeypress(event *input.KeyboardEvent) (event_handled bool) {
	return
}

// Adds an animation to the ui element.
func (e *ElementPrototype) AddAnimation(a gfx.Animator) {
	if e.animations == nil {
		e.animations = make([]gfx.Animator, 0)
	}

	//check for duplicate add
	for _, anim := range e.animations {
		if a == anim {
			return
		}
	}

	e.animations = append(e.animations, a)
}

func (e *ElementPrototype) IsVisible() bool {
	if !e.visible {
		return false
	}

	if parent := e.GetParent(); parent != nil {
		bounds := e.Bounds()
		if e.bordered {
			bounds.X, bounds.Y = bounds.X-1, bounds.Y-1
			bounds.W, bounds.H = bounds.W+2, bounds.H+2
		}
		if !vec.Intersects(bounds, parent.getCanvas()){
			return false
		}
	}

	return true
}

func (e *ElementPrototype) ToggleVisible() {
	e.SetVisible(!e.visible)
}

// Sets the visibility of the element. If we're making it visible, we trigger a render of the element.
// We also trigger a redraw of the parent element, in either case.
func (e *ElementPrototype) SetVisible(v bool) {
	if e.visible == v {
		return
	}

	e.visible = v
	if e.visible {
		e.Updated = true
	}

	e.forceParentRedraw()
}

func (e *ElementPrototype) getCanvas() *gfx.Canvas {
	return &e.Canvas
}

func (e *ElementPrototype) SetLabel(label string) {
	if e.label != "" {
		//changing labels. if we're in a window, remove the old label from the map
		if window := e.getWindow(); window != nil {
			window.removeLabel(e.label)
		}
	}

	e.label = label

	//get window, if it exists, and update the label map
	if window := e.getWindow(); window != nil {
		window.addLabel(e.label, e)
	}
}

func (e *ElementPrototype) GetLabel() string {
	return e.label
}

func (e *ElementPrototype) IsLabelled() bool {
	return e.label != ""
}
