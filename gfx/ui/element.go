package ui

import (
	"slices"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// element is the base api of anything handled by the UI system.
type element interface {
	vec.Bounded
	util.TreeType[element]
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
	getAnimations() []gfx.Animator
}

// Element is the base implementation for any UI Element handled by Tyumi's UI system. More complex UI elements can
// be created by embedding this and overriding the methods. Of highest importance are the Update() method, which
// describes how the element evolves with each tick, and the Render() method, which draws the element to the internal
// canvas when necessary. To trigger a render, set element.Updated = true when changing the element's state.
//
// Elements can be organized into a tree -- an element's children will be composited together with the element's
// own visuals when required. Note that children of an element are clipped by it, portions of the child not inside
// the element's bounds will not be drawn.
type Element struct {
	gfx.Canvas
	util.TreeNode[element]
	Updated bool   //indicates this object's state has changed and needs to be re-rendered.
	Border  Border //the element's border data. use EnableBorder() to turn on

	position    vec.Coord      //position relative to parent
	size        vec.Dims       //size of the drawable area of the element
	depth       int            //depth for the UI system, relative to the element's parent.
	visible     bool           //visibility, controlled via Show() and Hide()
	forceRedraw bool           //indicates this object needs to clear and render everything from zero
	label       string         //an optional identifier for the element
	animations  []gfx.Animator //animations on this element. these are updated once per frame while playing
}

// Init initializes the element, setting its size, as well as its position and depth relative to its parent. This must
// be done for all elements! Don't forget!
func (e *Element) Init(size vec.Dims, pos vec.Coord, depth int) {
	e.Canvas.Init(size)
	e.SetDefaultVisuals(defaultCanvasVisuals)
	e.position = pos
	e.size = size
	e.depth = depth
	e.visible = true
	e.Updated = true
	e.TreeNode.Init(e)
}

// Size returns the size of the drawable area of the element.
func (e *Element) Size() vec.Dims {
	return e.size
}

// Resizes the element. This clears the internal canvas and forces redraws of everything.
func (e *Element) Resize(size vec.Dims) {
	if size == e.size {
		return
	}

	if e.Border.enabled {
		e.Canvas.Resize(size.Grow(2, 2))
	} else {
		e.Canvas.Resize(size)
	}

	e.size = size
	e.Updated = true
	e.forceRedraw = true
	e.forceParentRedraw()
}

// Sets the default colours for draw operations on this element.
func (e *Element) SetDefaultColours(colours col.Pair) {
	e.Canvas.SetDefaultColours(colours)
	e.Updated = true
}

// Returns the bounding box of the element wrt its parent.
// Use Canvas.Bounds() to get the bounds of the underlying canvas for drawing to
func (e *Element) Bounds() vec.Rect {
	if e.Border.enabled {
		return e.Canvas.Bounds().Translated(e.position)
	}
	return vec.Rect{e.position, e.size}
}

func (e *Element) DrawableArea() vec.Rect {
	return e.size.Bounds()
}

func (e *Element) MoveTo(pos vec.Coord) {
	if e.position == pos {
		return
	}

	e.position = pos
	e.forceParentRedraw()
}

// THINK: should this take a coord too? or a Vec2i?
func (e *Element) Move(dx, dy int) {
	e.MoveTo(vec.Coord{e.position.X + dx, e.position.Y + dy})
}

// AddChild add a child element to this one. Child elements are composited together along with their parent to
// produce the final visuals for the element.
func (e *Element) AddChild(child element) {
	e.TreeNode.AddChild(child)
	if window := e.getWindow(); window != nil {
		window.onSubNodeAdded(child)
	}
	e.ForceRedraw()
}

func (e *Element) AddChildren(children ...element) {
	for _, child := range children {
		e.AddChild(child)
	}
}

func (e *Element) RemoveChild(child element) {
	e.TreeNode.RemoveChild(child)
	if window := e.getWindow(); window != nil {
		window.onSubNodeRemoved(child)
	}
	e.ForceRedraw()
}

// OVERRIDABLE FUNCTIONS!
// -----------------

// Update() can be overriden to update the state of the UI Element. Update() is called on each tick. If the element's
// state is changed and need to be redrawn, you can set its Updated flag to true to trigger a render on the next frame.
// Note that the element's animations are updated separately and do not need to be managed here.
func (e *Element) Update() {
	return
}

// Renders any changes in the element to the internal canvas. Override this to implement custom rendering behaviour.
// Elements are rendered if their Updated flag is true. Note that an element's children are composited seperately and
// you do not have to handle that here. Render() is called *after* child elements are drawn, and *before* any playing
// animations are drawn.
func (e *Element) Render() {
	return
}

// Handles keypresses. Override this to implement key input handling.
func (e *Element) HandleKeypress(event *input.KeyboardEvent) (event_handled bool) {
	return
}

// -------------------

func (e *Element) IsUpdated() bool {
	return e.Updated
}

func (e *Element) updateAnimations() {
	for _, a := range e.animations {
		if a.IsPlaying() {
			a.Update()
		}

		if !a.IsPlaying() { //animation just finished
			e.forceRedraw = true // make sure we reset just in case the animation left some garbage on the canvas
		}
	}

	// remove finished one-shot animations
	e.animations = slices.DeleteFunc(e.animations, func(a gfx.Animator) bool {
		return a.IsOneShot() && a.IsDone()
	})
}

// ForceRedraw will force an element to clear itself redraw all of its children, and perform a Render(). This generally
// isn't necessary as the UI system will trigger these operations automatically, only when strictly needed. But in cases
// where this can't be done you can use ForceRedraw to trigger the process manually.
func (e *Element) ForceRedraw() {
	e.forceRedraw = true
}

func (e *Element) isRedrawing() bool {
	return e.forceRedraw
}

func (e *Element) forceParentRedraw() {
	if parent := e.GetParent(); parent != nil {
		parent.ForceRedraw()
	}
}

// performs some pre-render operations. done for the whole tree before any rendering is done.
func (e *Element) prepareRender() {
	//if any animations are rendering this frame, trigger a redraw
	for _, a := range e.animations {
		if a.IsPlaying() && a.IsUpdated() {
			e.forceRedraw = true
			break
		}
	}

	if e.forceRedraw {
		e.Clear()
	}

	if e.Border.enabled && (e.Border.dirty || e.forceRedraw) {
		e.drawBorder()
	}
}

// performs some after-render cleanups. TODO: could also put some profiling code in here once that's a thing?
func (e *Element) finalizeRender() {
	if e.Border.enabled && (e.Border.dirty || e.forceRedraw) {
		e.linkBorder()
	}

	e.Updated = false
	e.forceRedraw = false
	e.Border.dirty = false
}

func (e *Element) renderAnimations() {
	for _, animation := range e.animations {
		if animation.IsPlaying() && vec.Intersects(e.getCanvas(), animation) {
			animation.Render(&e.Canvas)
		}
	}
}

func (e *Element) drawChildren() {
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

// Adds an animation to the ui element. Note that this does NOT start the animation.
func (e *Element) AddAnimation(animation gfx.Animator) {
	if e.animations == nil {
		e.animations = make([]gfx.Animator, 0)
	}

	//check for duplicate add
	for _, anim := range e.animations {
		if animation == anim {
			return
		}
	}

	e.animations = append(e.animations, animation)

	//if we're added a blocking animation during an update, make sure the window knows to stop updating
	if animation.IsBlocking() && animation.IsPlaying() {
		if wnd := e.getWindow(); wnd != nil {
			wnd.onBlockingAnimationAdded()
		}
	}
}

func (e *Element) getAnimations() []gfx.Animator {
	return e.animations
}

// IsVisible returns true if the element's visibility is true AND at least part of it is within the bounds of its
// parent.
// NOTE: this does NOT check if the parent is visible, so for a visible element whose parent is hidden this will still
// return true. Perhaps this behaviour should be changed... that would require a recursive check up the tree right to
// the window though, which could be expensive... hmm.
func (e *Element) IsVisible() bool {
	if !e.visible {
		return false
	}

	if parent := e.GetParent(); parent != nil {
		if !vec.Intersects(e, parent.getCanvas()) {
			return false
		}
	}

	return true
}

// Show makes the element visible.
func (e *Element) Show() {
	e.setVisible(true)
}

// Hide hides the element, preventing it and its children (if any) from receiving input, or being updated/rendered.
func (e *Element) Hide() {
	e.setVisible(false)
}

// ToggleVisible toggles element visibility.
func (e *Element) ToggleVisible() {
	e.setVisible(!e.visible)
}

func (e *Element) setVisible(v bool) {
	if e.visible == v {
		return
	}

	e.visible = v
	if e.visible {
		e.Updated = true
	}

	e.forceParentRedraw()
}

// SetLabel labels the element. References to labelled elements are retrievable from their parent window using
// GetLabelled().
func (e *Element) SetLabel(label string) {
	if window := e.getWindow(); window != nil {
		if e.label != "" {
			//changing labels. if we're in a window, remove the old label from the map
			window.removeLabel(e.label)
		}
		window.addLabel(label, e)
	}

	e.label = label
}

func (e *Element) GetLabel() string {
	return e.label
}

func (e *Element) IsLabelled() bool {
	return e.label != ""
}

func (e *Element) getPosition() vec.Coord {
	return e.position
}

func (e *Element) getDepth() int {
	return e.depth
}

func (e *Element) getCanvas() *gfx.Canvas {
	return &e.Canvas
}

func (e *Element) getWindow() *Window {
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
