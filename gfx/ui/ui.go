//UI is Tyumi's UI system. A base UI element is defined and successively more complex elements are composed from it. UI
//elements can be nested in Containers for composition alongside/inside one another.
package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
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
	HandleKeypress(input.KeyboardEvent)
	MoveTo(int, int)
	Move(int, int)
	IsVisible() bool
}

type ElementPrototype struct {
	gfx.Canvas

	position vec.Coord
	depth        int //depth for the UI system, relative to the element's parent. if no parent, relative to the console
	visible  bool
	bordered bool
	dirty    bool //indicates this object needs to be re-rendered.

	border     Border
	animations []gfx.Animator //animations on this element. these are updated once per frame.

	parent *Container
}

func (e *ElementPrototype) Init(w, h, x, y, depth int) {
	e.Canvas.Init(w, h)
	e.position = vec.Coord{x, y}
	e.depth = depth
	e.visible = true
	e.dirty = true
}

func (e *ElementPrototype) SetDefaultColours(fore uint32, back uint32) {
	e.Canvas.SetDefaultColours(fore, back)
	e.dirty = true
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

func (e *ElementPrototype) Pos() vec.Coord {
	return e.position
}

func (e *ElementPrototype) MoveTo(x, y int) {
	if x == e.position.X && y == e.position.Y {
		return
	}

	e.position.MoveTo(x, y)
	if e.parent != nil {
		e.parent.Redraw()
	}
}

func (e *ElementPrototype) Move(dx, dy int) {
	e.MoveTo(e.position.X+dx, e.position.Y+dy)
}

//update() is the internal update function. handles any internal update behaviour, and calls the UpdateState function
//to allow user-defined update behaviour to occur.
func (e *ElementPrototype) update() {
	e.UpdateState()

	//tick animations
	for _, a := range e.animations {
		if a.Done() {
			//TODO: remove animation from list if it's done
			continue
		}

		a.Update()

		if a.Dirty() {
			e.dirty = true
		}
	}
}

//UpdateState() is a virtual function. Implement this to provide ui update behaviour on a thread-safe, per-frame basis,
//instead of updating the state of the element as the gamestate progresses.
func (e *ElementPrototype) UpdateState() {

}

//Renders any changes in the element to the internal canvas. If the element is not visible, we don't waste precious cpus
//rendering to nothing.
func (e *ElementPrototype) Render() {
	if !e.visible {
		return
	}

	if e.dirty {
		e.dirty = false
	}

	if e.bordered {
		e.border.Render()
	}

	for _, a := range e.animations {
		a.Render(&e.Canvas)
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

	e.DrawToCanvas(&e.parent.Canvas, e.position.X, e.position.Y, e.depth)
	if e.bordered {
		e.border.DrawToCanvas(&e.parent.Canvas, e.position.X, e.position.Y, e.depth)
	}
}

func (e *ElementPrototype) HandleKeypress(event input.KeyboardEvent) {
	return
}

//Adds an animation to the ui element.
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

//Sets the visibility of the element. If we're making it visible, we trigger a render of the element.
//We also trigger a redraw of any parent element, in either case.
func (e *ElementPrototype) SetVisible(v bool) {
	if e.visible == v {
		return
	}

	e.visible = v

	if e.visible {
		e.dirty = true
	}

	if e.parent != nil {
		e.parent.Redraw()
	}
}
