// UI is Tyumi's UI system. A base UI element is defined and successively more complex elements are composed from it. UI
// elements can be nested in Containers for composition alongside/inside one another.
package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// An element is the base structure of anything handled by the UI system.
type Element interface {
	vec.Bounded
	util.TreeType[Element]
	Update()           // do not override this, this is what the engine uses to tick the gamestate
	UpdateState() bool //if you want custom update code, implement this.
	Render()
	ForceRedraw() //Force the element to clear and redraw itself and all children from scratch
	HandleKeypress(input.KeyboardEvent)
	MoveTo(vec.Coord)
	Move(int, int)
	IsVisible() bool
	IsUpdated() bool
	DrawToParent()
	getCanvas() *gfx.Canvas
}

type ElementPrototype struct {
	gfx.Canvas
	util.TreeNode[Element]

	position    vec.Coord
	depth       int //depth for the UI system, relative to the element's parent. if no parent, relative to the console
	visible     bool
	bordered    bool
	updated     bool //indicates this object's state has changed and needs to be re-rendered.
	forceRedraw bool //indicates this object needs to clear and render everything from zero

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

func (e *ElementPrototype) SetDefaultColours(fore uint32, back uint32) {
	e.Canvas.SetDefaultColours(fore, back)
	e.updated = true
}

func (e *ElementPrototype) EnableBorder(title, hint string) {
	e.bordered = true
	e.border = NewBorder(e.Size())
	e.border.title = title
	e.border.hint = hint
}

func (e *ElementPrototype) Bounds() vec.Rect {
	return vec.Rect{e.position, e.Size()}
}

func (e *ElementPrototype) MoveTo(pos vec.Coord) {
	if e.position == pos {
		return
	}

	e.position = pos
	if parent := e.GetParent(); parent != nil {
		parent.ForceRedraw()
	}
}

// THINK: should this take a coord too? or a Vec2i?
func (e *ElementPrototype) Move(dx, dy int) {
	e.MoveTo(vec.Coord{e.position.X + dx, e.position.Y + dy})
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
		if a.Done() {
			//TODO: remove animation from list if it's done
			continue
		}

		a.Update()

		if a.Dirty() {
			e.updated = true
		}
	}
}

// UpdateState() is a virtual function. Implement this to provide ui update behaviour on a thread-safe, per-frame basis,
// instead of updating the state of the element as the gamestate progresses. Return true if you want to trigger a 
// render of the ui element.
func (e *ElementPrototype) UpdateState() bool {
	return false
}

func (e *ElementPrototype) IsUpdated() bool {
	return e.updated
}

func (e *ElementPrototype) ForceRedraw() {
	e.forceRedraw = true
}

// Renders any changes in the element to the internal canvas. If the element is not visible, we don't waste precious cpus
// rendering to nothing.
func (e *ElementPrototype) Render() {
	if !e.visible {
		return
	}

	if e.bordered {
		e.border.Render()
	}

	if e.forceRedraw {
		e.Canvas.Clear()
	}

	for _, a := range e.animations {
		a.Render(&e.Canvas)
	}

	for _, child := range e.GetChildren() {
		if child.IsVisible() {
			if child.IsUpdated() {
				child.Render()
			}
			if child.getCanvas().Dirty() || e.forceRedraw {
				child.DrawToParent()
			}
		}
	}

	e.updated = false
	e.forceRedraw = false
}

func (e *ElementPrototype) DrawToParent() {
	var parent Element
	if parent = e.GetParent(); parent == nil {
		return
	}

	e.DrawToCanvas(parent.getCanvas(), e.position, e.depth)
	if e.bordered {
		e.border.DrawToCanvas(parent.getCanvas(), e.position.X, e.position.Y, e.depth)
	}
}

func (e *ElementPrototype) HandleKeypress(event input.KeyboardEvent) {
	for _, child := range e.GetChildren() {
		if child.IsVisible() {
			child.HandleKeypress(event)
		}
	}
}

// Adds an animation to the ui element.
func (e *ElementPrototype) AddAnimation(a gfx.Animator) {
	if e.animations == nil {
		e.animations = make([]gfx.Animator, 0)
	} else {
		for _, anim := range e.animations {
			if a == anim {
				return
			}
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

	if parent := e.GetParent(); parent != nil {
		parent.ForceRedraw()
	}
}

func (e *ElementPrototype) getCanvas() *gfx.Canvas {
	return &e.Canvas
}
