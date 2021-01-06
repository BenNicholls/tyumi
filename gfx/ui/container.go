package ui

//Container allows you to group a number of children objects together. Children can be any UI Element, including other
//containers. Go nuts why don't ya.
type Container struct {
	ElementPrototype

	children []Element
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

		c.children = append(c.children, e)
	}
}