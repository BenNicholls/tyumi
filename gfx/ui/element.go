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
	Update()           // do not override this, this is what the engine uses to tick the gamestate
	UpdateState() bool //if you want custom update code, implement this.
	Render()
	RenderAnimations()
	FinalizeRender()
	ForceRedraw() //Force the element to clear and redraw itself and all children from scratch
	HandleKeypress(*input.KeyboardEvent)
	MoveTo(vec.Coord)
	Move(int, int)
	IsVisible() bool
	DrawToParent()
	getCanvas() *gfx.Canvas
	getBorder() *Border
	getWindow() *Window
}

type ElementPrototype struct {
	gfx.Canvas
	util.TreeNode[Element]

	position    vec.Coord
	depth       int //depth for the UI system, relative to the element's parent. if no parent, relative to the console
	visible     bool
	bordered    bool
	updated     bool   //indicates this object's state has changed and needs to be re-rendered.
	forceRedraw bool   //indicates this object needs to clear and render everything from zero
	label       string // an optional identifier for the element

	border     Border
	animations []gfx.Animator //animations on this element. these are updated once per frame.
}

func (e *ElementPrototype) Init(w, h int, pos vec.Coord, depth int) {
	e.Canvas.Init(w, h)
	e.position = pos
	e.depth = depth
	e.visible = true
	e.updated = true
	e.TreeNode.Init(e)
}

func (e *ElementPrototype) SetDefaultColours(colours col.Pair) {
	e.Canvas.SetDefaultColours(colours)
	if e.bordered {
		e.border.setColours(e.DefaultColours())
	}
	e.updated = true
}

// Enable the border. Defaults to ui.DefaultBorderStyle. Use SetBorderStyle to use something else.
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
	e.forceParentRedraw()
}

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
	switch styleFlag {
	case BORDER_STYLE_DEFAULT:
		e.border.style = &DefaultBorderStyle
	case BORDER_STYLE_INHERIT:
		if parent := e.GetParent(); parent != nil {
			parent_border := parent.getBorder()
			e.border.style = parent_border.style
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
	return &e.border
}

func (e *ElementPrototype) getWindow() *Window {
	parent := e.GetParent()
	if parent == nil {
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
	if child_border := child.getBorder(); child_border.styleFlag == BORDER_STYLE_INHERIT {
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

// update() is the internal update function. handles any internal update behaviour, and calls the UpdateState function
// to allow user-defined update behaviour to occur.
func (e *ElementPrototype) Update() {
	for _, e := range e.GetChildren() {
		e.Update()
	}

	//run user-provided state update function.
	if e.UpdateState() {
		e.updated = true
	}

	//tick animations
	for _, a := range e.animations {
		if a.Playing() {
			a.Update()
		}
	}

	//remove finished one-shot animations
	e.animations = slices.DeleteFunc[[]gfx.Animator](e.animations, func(a gfx.Animator) bool {
		return a.IsOneShot() && a.Done()
	})
}

// UpdateState() is a virtual function. Implement this to provide ui update behaviour on a thread-safe, per-frame basis,
// instead of updating the state of the element as the gamestate progresses. Return true if you want to trigger a
// render of the ui element.
func (e *ElementPrototype) UpdateState() bool {
	return false
}

func (e *ElementPrototype) ForceRedraw() {
	e.forceRedraw = true
}

func (e *ElementPrototype) forceParentRedraw() {
	parent := e.GetParent()
	if parent != nil {
		parent.ForceRedraw()
	}
}

// Renders any changes in the element to the internal canvas. If the element is not visible, we don't waste precious cpus
// rendering to nothing.
func (e *ElementPrototype) Render() {
	if !e.visible {
		return
	}

	if e.bordered {
		if e.forceRedraw {
			e.border.dirty = true
		}
		e.border.Render()
	}

	if e.forceRedraw {
		for _, child := range e.GetChildren() { //make sure siblings recompute border links
			child.getBorder().dirty = true
		}
		e.Canvas.Clear()
	}

	for _, child := range e.GetChildren() {
		//BUG: visibility culling doesn't take the border of the child into account.
		//instead of fixing, might be better to redesign how borders work. Or make elements
		//dynamically adjust their bounds when borders are activated?? hmm.
		if child.IsVisible() && vec.FindIntersectionRect(e.getCanvas(), child).Area() > 0 {
			child.Render()
			child.RenderAnimations()
			if child.getCanvas().Dirty() || e.forceRedraw {
				child.DrawToParent()
			}
			child.FinalizeRender()
		}
	}
}

func (e *ElementPrototype) RenderAnimations() {
	for _, animation := range e.animations {
		if animation.Playing() && vec.FindIntersectionRect(e.getCanvas(), animation).Area() > 0 {
			animation.Render(&e.Canvas)
		}
	}
}

// performs some after-render cleanups. TODO: could also put some profiling code in here once that's a thing?
func (e *ElementPrototype) FinalizeRender() {
	e.updated = false
	e.forceRedraw = false
}

func (e *ElementPrototype) DrawToParent() {
	parent := e.GetParent()
	if parent == nil {
		return
	}

	e.DrawToCanvas(parent.getCanvas(), e.position, e.depth)
	if e.bordered {
		e.border.DrawToCanvas(parent.getCanvas(), e.position, e.depth)
	}
}

func (e *ElementPrototype) HandleKeypress(event *input.KeyboardEvent) {
	for _, child := range e.GetChildren() {
		if event.Handled() {
			break
		}

		if child.IsVisible() {
			child.HandleKeypress(event)
		}
	}
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
	return e.visible
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
		e.updated = true
	}

	e.forceParentRedraw()
}

func (e *ElementPrototype) getCanvas() *gfx.Canvas {
	return &e.Canvas
}

func (e *ElementPrototype) SetLabel(label string) {
	e.label = label
	//get window, if it exists, and update the label map
}

func (e *ElementPrototype) GetLabel() string {
	return e.label
}

func (e *ElementPrototype) IsLabelled() bool {
	return e.label != ""
}
