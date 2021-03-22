//UI is Tyumi's UI system. A base UI element is defined and successively more complex elements are composed from it. UI
//elements can be nested in Containers for composition alongside/inside one another.
package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/vec"
)

//An element is the base structure of anything handled by the UI system.
type Element interface {
	vec.Bounded
	AddParent(*Container)
	DrawToParent()
	Render()
	update()
	UpdateState()
}

type ElementPrototype struct {
	gfx.Canvas

	position vec.Coord
	z        int //depth for the UI system, relative to the element's parent. if no parent, relative to the console
	visible  bool
	dirty    bool //indicates this object needs to be re-rendered.

	animations []gfx.Animator //animations on this element. these are updated once per frame.

	parent *Container //parent container. if nil,
}

func (e *ElementPrototype) Init(w, h, x, y, z int) {
	e.Canvas.Init(w, h)
	e.position = vec.Coord{x, y}
	e.z = z
	e.animations = make([]gfx.Animator, 0)
	e.visible = true
	e.dirty = true
}

func (e *ElementPrototype) SetDefaultColours(fore uint32, back uint32) {
	e.Canvas.SetDefaultColours(fore, back)
	e.dirty = true
}

func (e *ElementPrototype) Bounds() vec.Rect {
	w, h := e.Dims()
	return vec.Rect{w, h, e.position.X, e.position.Y}
}

func (e *ElementPrototype) Pos() vec.Coord {
	return e.position
}

func (e *ElementPrototype) MoveTo(x, y int) {
	e.position.MoveTo(x, y)
	if e.parent != nil {
		e.parent.Redraw()
	}
}

//update() is the internal update function. handles any internal update behaviour, and calls the UpdateState function
//to allow user-defined update behaviour to occur.
func (e *ElementPrototype) update() {
	e.UpdateState()
}

//UpdateState() is a virtual function. Implement this to provide ui update behaviour on a thread-safe, per-frame basis,
//instead of updating the state of the element as the gamestate progresses.
func (e *ElementPrototype) UpdateState() {

}

//Renders any changes in the element to the internal canvas.
func (e *ElementPrototype) Render() {
	if e.dirty {
		e.dirty = false
	}
}

func (e *ElementPrototype) AddParent(c *Container) {
	if e.parent != nil {
		e.parent.RemoveElement(e)
	}

	e.parent = c
}

func (e *ElementPrototype) DrawToParent() {
	if e.parent == nil {
		return
	}

	e.DrawToCanvas(&e.parent.Canvas, e.position.X, e.position.Y, e.z)
}
