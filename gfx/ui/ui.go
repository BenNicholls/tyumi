//UI is Tyumi's UI system. A base UI element is defined and successively more complex elements are composed from it. UI 
//elements can be nested in Containers for composition alongside one another.
package ui

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/vec"
)

//An element is the base structure of anything handled by the UI system.
type Element interface {
	vec.Bounded
	Pos() *vec.Coord
	Update()
	Render()
}

type ElementPrototype struct {
	gfx.Canvas

	position vec.Coord
	z        int //depth for the UI system, relative to the element's parent. if no parent, relative to the console
	visible  bool
	dirty    bool //indicates this object needs to be re-rendered.

	foreColour uint32 //defaults to col.WHITE
	backColour uint32 //defaults to col.BLACK

	animations []gfx.Animator //animations on this element. these are updated once per frame.
}

func (e *ElementPrototype) Init(w, h, x, y, z int) {
	e.Canvas.Init(w, h)
	e.position = vec.Coord{x, y}
	e.z = z
	e.foreColour = col.WHITE
	e.backColour = col.BLACK
	e.animations = make([]gfx.Animator, 0)
	e.dirty = true
} 

func (e *ElementPrototype) Bounds() vec.Rect {
	w, h := e.Dims()
	return vec.Rect{w, h, e.position.X, e.position.Y}
}

func (e *ElementPrototype) Pos() *vec.Coord {
	return &e.position
}

//Update is a virtual function. Override this to provide ui update behaviour on a thread-safe, per-frame basis, 
//instead of updating the element as the gamestate progresses.
func (e *ElementPrototype) Update() {

}

//Renders any changes in the element to the internal canvas.
func (e *ElementPrototype) Render() {

}
