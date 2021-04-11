package ui

import "github.com/bennicholls/tyumi/input"

//Container allows you to group a number of children objects together. Children can be any UI Element, including other
//containers. Go nuts why don't ya.
type Container struct {
	ElementPrototype

	children []Element
	redraw   bool //total redraw required. see Container.Redraw()
}

func NewContainer(w, h, x, y, z int) (c Container) {
	c.Init(w, h, x, y, z)
	c.children = make([]Element, 0)

	return
}

func (c *Container) AddElement(elems ...Element) {
	for _, e := range elems {
		//check for duplicate entry
		dupe := false
		for _, child := range c.children {
			if child == e {
				dupe = true
				break
			}
		}
		if dupe {
			continue
		}

		e.AddParent(c)
		c.children = append(c.children, e)
	}
}

func (c *Container) RemoveElement(e Element) {
	for i, child := range c.children {
		if child == e {
			copy(c.children[i:], c.children[i+1:])
			c.children[len(c.children)-1] = nil
			c.children = c.children[:len(c.children)-1]
		}
	}
}

//Call this to indicate to the container that it must do a complete redraw from scratch.
//Usually this is because a child element has moved, toggled visibility, or done some other
//cool move. The next call to render will perform a clear first, then reset the flag.
func (c *Container) Redraw() {
	c.redraw = true
}

func (c *Container) update() {
	c.UpdateChildren()
	c.UpdateState()
}

func (c *Container) UpdateChildren() {
	for _, e := range c.children {
		e.update()
	}
}

//Render composites all internal elements into the container's canvas.
func (c *Container) Render() {
	if c.redraw {
		c.Canvas.Clear()
		c.redraw = false
	}

	for _, child := range c.children {
		child.Render()
		child.DrawToParent()
	}

	c.ElementPrototype.Render()
}

//Containers take keyboard input events and pass them to their children, in case any of them want
//to handle the input
func (c *Container) HandleKeypress(e input.KeyboardEvent) {
	for _, child := range c.children {
		child.HandleKeypress(e)
	}
}
