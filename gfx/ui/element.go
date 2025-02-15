package ui

import (
	"slices"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
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
	drawChildren()
	ForceRedraw() //Force the element to clear and redraw itself and all children from scratch
	isRedrawing() bool

	HandleKeypress(*input.KeyboardEvent) (event_handled bool)

	MoveTo(vec.Coord)
	Move(int, int)

	IsVisible() bool
	IsUpdated() bool
	IsBordered() bool
	Size() vec.Dims
	getCanvas() *gfx.Canvas
	getWindow() *Window
	getBorderStyle() BorderStyle
	getDepth() int
	getPosition() vec.Coord
}

type ElementPrototype struct {
	gfx.Canvas
	util.TreeNode[Element]
	Updated bool   //indicates this object's state has changed and needs to be re-rendered.
	Border  Border //the element's border data. use EnableBorder() to turn on

	position    vec.Coord
	size        vec.Dims
	depth       int //depth for the UI system, relative to the element's parent.
	visible     bool
	forceRedraw bool           //indicates this object needs to clear and render everything from zero
	label       string         // an optional identifier for the element
	animations  []gfx.Animator //animations on this element. these are updated once per frame.
}

func (e *ElementPrototype) Init(w, h int, pos vec.Coord, depth int) {
	e.Canvas.Init(w, h)
	e.SetDefaultVisuals(defaultCanvasVisuals)
	e.position = pos
	e.size = vec.Dims{w, h}
	e.depth = depth
	e.visible = true
	e.Updated = true
	e.TreeNode.Init(e)
}

func (e *ElementPrototype) Size() vec.Dims {
	return e.size
}

// Resizes the element. This clears the internal canvas and forces redraws of everything.
func (e *ElementPrototype) Resize(size vec.Dims) {
	if size == e.size {
		return
	}

	if e.Border.enabled {
		e.Canvas.Resize(size.W+2, size.H+2)
	} else {
		e.Canvas.Resize(size.W, size.H)
	}

	e.size = size
	e.Updated = true
	e.forceRedraw = true
	e.forceParentRedraw()
}

func (e *ElementPrototype) SetDefaultColours(colours col.Pair) {
	e.Canvas.SetDefaultColours(colours)
	e.Updated = true
}

// Returns the bounding box of the element wrt its parent.
// Use Canvas.Bounds() to get the bounds of the underlying canvas for drawing to
func (e *ElementPrototype) Bounds() vec.Rect {
	if e.Border.enabled {
		return e.Canvas.Bounds().Translated(e.position)
	}
	return vec.Rect{e.position, e.size}
}

func (e *ElementPrototype) DrawableArea() vec.Rect {
	return vec.Rect{vec.ZERO_COORD, e.size}
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

// OVERRIDABLE FUNCTIONS!
// -----------------

// Update() can be overriden to update the state of the UI Element. Note that the element's animations are updated
// separately and do not need to be managed here.
func (e *ElementPrototype) Update() {
	return
}

// Renders any changes in the element to the internal canvas. Override this to implement custom rendering behaviour.
// Note that this is called *after* any subelements are drawn to the canvas, and *before* any running animations
// are rendered.
func (e *ElementPrototype) Render() {
	return
}

// Handles keypresses. Override this to implement key input handling.
func (e *ElementPrototype) HandleKeypress(event *input.KeyboardEvent) (event_handled bool) {
	return
}

// -------------------

func (e *ElementPrototype) IsUpdated() bool {
	return e.Updated
}

func (e *ElementPrototype) updateAnimations() {
	for _, a := range e.animations {
		if a.IsPlaying() {
			a.Update()
		}
	}

	// remove finished one-shot animations
	e.animations = slices.DeleteFunc(e.animations, func(a gfx.Animator) bool {
		return a.IsOneShot() && a.IsDone()
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

// performs some pre-render operations. done for the whole tree before any rendering is done.
func (e *ElementPrototype) prepareRender() {
	if e.forceRedraw {
		e.Clear()
	}

	if e.Border.enabled && (e.Border.dirty || e.forceRedraw) {
		e.DrawBorder()
	}
}

// performs some after-render cleanups. TODO: could also put some profiling code in here once that's a thing?
func (e *ElementPrototype) finalizeRender() {
	if e.Border.enabled && (e.Border.dirty || e.forceRedraw) {
		e.linkBorder()
	}

	e.Updated = false
	e.forceRedraw = false
	e.Border.dirty = false
}

func (e *ElementPrototype) renderAnimations() {
	for _, animation := range e.animations {
		if animation.IsPlaying() && vec.Intersects(e.getCanvas(), animation) {
			animation.Render(&e.Canvas)
		}
	}
}

func (e *ElementPrototype) drawChildren() {
	for i, child := range e.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		if child.getCanvas().Dirty() || e.forceRedraw {
			child.getCanvas().Draw(&e.Canvas, child.getPosition(), child.getDepth())
			if child.IsBordered() {
				// attempt to link to siblings' borders
				for sib_i, sibling := range e.GetChildren() {
					// if we're doing a forced redraw of all children then we only need to link to siblings that have
					// already been drawn. subsequently drawn elements will then link to this one. so we can break
					// when we get to this element
					if e.forceRedraw && sib_i == i {
						break
					}

					//figure out if we should be linking to this sibling whatsoever. lots of things to consider
					if child == sibling || !sibling.IsBordered() || sibling.getDepth() != child.getDepth() || !sibling.IsVisible() {
						continue
					}

					intersection := vec.FindIntersectionRect(child, sibling)
					if intersection.Area() == 0 {
						continue
					}

					//THERE ARE LIKE 20 WAYS RECTANGLES CAN OVERLAP. LET'S CHECK THEM ALL!
					// Dear future Ben: i know what you're thinking. there must be a pattern here that we can use to
					// simplify this monstrosity. trust me, you looked and couldn't see one that covered all 20+
					// cases cleanly. maybe one exists, hell it probably does, but this appears to work and is fast.
					// it just looks awful. so leave it alone and go make a game or something.
					// - forever yours, Past Ben
					switch {
					case intersection.Area() == 1:
						e.LinkCell(intersection.Coord)
					case intersection.W == 1 || intersection.H == 1:
						corners := intersection.Corners()
						e.LinkCell(corners[0])
						e.LinkCell(corners[2])
					default:
						corners := intersection.Corners()
						c := child.Bounds()
						s := sibling.Bounds()
						switch {
						case intersection.W == s.W || intersection.H == s.H:
							for _, corner := range corners {
								e.LinkCell(corner)
							}
						case c.X < s.X && c.Y < s.Y:
							if c.X+c.W > s.X+s.W {
								e.LinkCell(corners[2])
								e.LinkCell(corners[3])
							} else if c.Y+c.H > s.Y+s.H {
								e.LinkCell(corners[1])
								e.LinkCell(corners[2])
							} else {
								e.LinkCell(corners[1])
								e.LinkCell(corners[3])
							}
						case c.X < s.X && c.Y > s.Y:
							e.LinkCell(corners[0])
							if c.X+c.W > s.X+s.W {
								e.LinkCell(corners[1])
							} else if c.Y+c.H >= s.Y+s.H {
								e.LinkCell(corners[2])
							} else {
								e.LinkCell(corners[3])
							}
						case c.X > s.X && c.Y <= s.Y:
							e.LinkCell(corners[0])
							if c.Y+c.H > s.Y+s.H {
								e.LinkCell(corners[3])
							} else if c.X+c.W > s.X+s.W {
								e.LinkCell(corners[2])
							} else {
								e.LinkCell(corners[1])
							}
						case c.X > s.X && c.Y > s.Y:
							if c.X+c.W < s.X+s.W {
								e.LinkCell(corners[2])
								e.LinkCell(corners[3])
							} else if c.Y+c.H < s.Y+s.H {
								e.LinkCell(corners[1])
								e.LinkCell(corners[2])
							} else {
								e.LinkCell(corners[1])
								e.LinkCell(corners[3])
							}
						}
					}
				}
			}
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
	if !e.visible {
		return false
	}

	if parent := e.GetParent(); parent != nil {
		if !vec.Intersects(e.Bounds(), parent.getCanvas()) {
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

func (e *ElementPrototype) getPosition() vec.Coord {
	return e.position
}

func (e *ElementPrototype) getDepth() int {
	return e.depth
}

func (e *ElementPrototype) getCanvas() *gfx.Canvas {
	return &e.Canvas
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
