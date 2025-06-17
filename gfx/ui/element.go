package ui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

// element is the base api of anything handled by the UI system.
type element interface {
	vec.Bounded
	util.TreeType[element]
	Labelled
	event.Listener

	Update()
	IsUpdated() bool
	updateAnimations()

	prepareRender()
	renderIfDirty()
	Render()
	Draw(dst_canvas *gfx.Canvas, force bool)
	renderAnimations()
	finalizeRender()
	drawChildren()
	ForceRedraw() //Force the element to clear and redraw itself and all children from scratch
	IsRedrawing() bool

	acceptsInput() bool
	HandleKeypress(*input.KeyboardEvent) (event_handled bool)
	HandleAction(action input.ActionID) (action_handled bool)

	MoveTo(vec.Coord)
	Move(dx, dy int)

	Focus()
	Defocus()
	IsFocused() bool

	// some other getters
	IsVisible() bool
	IsTransparent() bool
	IsBordered() bool
	Size() vec.Dims
	ID() ElementID

	// internal getters
	getCanvas() *gfx.Canvas
	getWindow() *Window
	getBorder() *Border
	getBorderStyle() BorderStyle
	getDepth() int
	getPosition() vec.Coord
	getAnimations() []gfx.Animator

	// debug functions
	dumpUI(dir_name string, depth int)
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
	event.Stream
	Updated     bool   //indicates this object's state has changed and needs to be re-rendered.
	AcceptInput bool   // if true, the element will be sent inputs when in a window with SendEventsToUnfocused = true
	OnRender    func() // a callback, called when the element renders (unless the element has custom rendering logic)
	Border      Border //the element's border data. use EnableBorder() to turn on

	visible     bool           //visibility, controlled via Show() and Hide()
	focused     bool           //focus state. by default, only focused elements receive input
	forceRedraw bool           //indicates this object needs to clear and render everything from zero
	position    vec.Coord      //position relative to parent
	size        vec.Dims       //size of the drawable area of the element
	depth       int            //depth for the UI system, relative to the element's parent.
	id          ElementID      //a unique ID for the element
	label       string         //an optional identifier for the element
	animations  []gfx.Animator //animations on this element. these are updated once per frame while playing
}

func (e *Element) String() string {
	desc := fmt.Sprintf("[UI Element] pos: %v, size: %v, depth: %d, visible: %t\n +-- %v", e.position, e.size, e.depth, e.visible, e.Canvas)
	if e.label != "" {
		desc += "\n +-- Label: " + e.label
	}
	if e.Border.enabled {
		desc += fmt.Sprintf("\n +-- Border Enabled. Title: %s, Hint: %s", e.Border.title, e.Border.hint)
	}
	if len(e.animations) != 0 {
		playing := 0
		for _, anim := range e.animations {
			if anim.IsPlaying() {
				playing++
			}
		}
		desc += fmt.Sprintf("\n +-- Animations: %d, (%d playing)", len(e.animations), playing)
	}

	return desc
}

// Init initializes the element, setting its size, as well as its position and depth relative to its parent. This must
// be done for all elements! Don't forget!
func (e *Element) Init(size vec.Dims, pos vec.Coord, depth int) {
	e.Canvas.Init(size)
	e.SetDefaultVisuals(DefaultElementVisuals)
	e.position = pos
	e.size = size
	e.depth = depth
	e.visible = true
	e.Updated = true
	e.forceRedraw = true
	e.TreeNode.Init(e)
	e.id = generate_id()

	// if element is not in a window, ensure it and all children that may already be attached have their event streams
	// disabled
	if e.getWindow() == nil {
		util.WalkTree[element](e, func(element element) { element.DisableListening() })
	}
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
	if e.Border.enabled {
		e.Border.dirty = true
	}
}

func (e *Element) SetDefaultVisuals(vis gfx.Visuals) {
	e.Canvas.SetDefaultVisuals(vis)
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

func (e *Element) getPosition() vec.Coord {
	return e.position
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

// Centers the element within its parent. If not a child of another element, does nothing.
// NOTE: this does not keep the object centered, if the object changes shape or the parent does something, this must
// be called again.
func (e *Element) Center() {
	e.CenterHorizontal()
	e.CenterVertical()
}

// Centers the element horizontally within its parent. If not a child of another element, does nothing.
// NOTE: this does not keep the object centered, if the object changes shape or the parent does something, this must
// be called again.
func (e *Element) CenterHorizontal() {
	parent := e.GetParent()
	if parent == nil {
		return
	}

	e.MoveTo(vec.Coord{(parent.Size().W - e.size.W) / 2, e.position.Y})
}

// Centers the element vertically within its parent. If not a child of another element, does nothing.
// NOTE: this does not keep the object centered, if the object changes shape or the parent does something, this must
// be called again.
func (e *Element) CenterVertical() {
	parent := e.GetParent()
	if parent == nil {
		return
	}

	e.MoveTo(vec.Coord{e.position.X, (parent.Size().H - e.size.H) / 2})
}

// AddChild adds a child element to this one. Child elements are composited together along with their parent to
// produce the final visuals for the element.
func (e *Element) AddChild(child element) {
	if child.ID() == e.id {
		log.Error("Tried to add an element as a child of itself! Why???")
		return
	}

	if parent := child.GetParent(); parent != nil && parent.ID() == e.ID() {
		log.Warning("Tried to add an element to its own parent! Why add twice????")
		return
	}

	e.TreeNode.AddChild(child)
	if window := e.getWindow(); window != nil {
		window.onSubNodeAdded(child)
		util.WalkTree(child, func(element element) { element.EnableListening() }, ifVisible)
	}
	e.ForceRedraw()
}

func (e *Element) AddChildren(children ...element) {
	for _, child := range children {
		e.AddChild(child)
	}
}

func (e *Element) RemoveChild(child element) {
	oldChildCount := e.ChildCount()
	e.TreeNode.RemoveChild(child)
	if e.ChildCount() == oldChildCount {
		log.Debug("Child not actually a child, no remove done.")
		return
	}

	if window := e.getWindow(); window != nil {
		window.onSubNodeRemoved(child)
		util.WalkTree(child, func(element element) { element.DisableListening() })
	}
	e.ForceRedraw()
}

func (e *Element) RemoveAllChildren() {
	for _, child := range slices.Backward(e.GetChildren()) {
		e.RemoveChild(child)
	}
}

// OVERRIDABLE FUNCTIONS!
// -----------------

// Update() can be overriden to update the state of the UI Element. Update() is called on each tick. If the element's
// state is changed and need to be redrawn, you can set its Updated flag to true to trigger a render on the next frame.
// Note that the element's animations are updated separately and do not need to be managed here.
func (e *Element) Update() {}

// Renders any changes in the element to the internal canvas. Override this to implement custom rendering behaviour. If
// this method has not been overriden, it attempts to call the user-provided OnRender() callback, if any.
// Elements are rendered if their Updated flag is true. Note that an element's children are composited seperately and
// you do not have to handle that here. Render() is called *after* child elements are drawn, and *before* any playing
// animations are drawn.
func (e *Element) Render() {
	if e.OnRender != nil {
		e.OnRender()
	}
}

// Handles keypresses. Override this to implement key input handling.
func (e *Element) HandleKeypress(event *input.KeyboardEvent) (event_handled bool) { return }

// Handles Actions. Override this to implement action handling.
func (e *Element) HandleAction(action input.ActionID) (action_handled bool) { return }

// -------------------

func (e *Element) acceptsInput() bool {
	return e.AcceptInput || e.focused
}

func (e *Element) IsUpdated() bool {
	return e.Updated
}

// ForceRedraw will force an element to clear itself, redraw all of its children, and perform a Render(). This generally
// isn't necessary as the UI system will trigger these operations automatically, only when strictly needed. But in cases
// where this can't be done you can use ForceRedraw to trigger the process manually.
func (e *Element) ForceRedraw() {
	e.forceRedraw = true
}

func (e *Element) IsRedrawing() bool {
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

// this is called during rendering if you want code run in the case a child-draw occurs and dirties the canvas.
func (e *Element) renderIfDirty() {}

// performs some after-render cleanups. TODO: could also put some profiling code in here once that's a thing?
func (e *Element) finalizeRender() {
	for _, anim := range e.animations {
		anim.Finish()
	}

	e.Updated = false
	e.forceRedraw = false
	if e.Border.enabled {
		e.Border.dirty = false
		e.Border.internalLinksRecalculated = false
	}
}

// Draws the element to a provided canvas, based on the element's position and respecting depth.
func (e *Element) Draw(dst_canvas *gfx.Canvas, force bool) {
	if force {
		e.Canvas.Draw(dst_canvas, e.position, e.depth, gfx.DRAWFLAG_FORCE)
	} else {
		e.Canvas.Draw(dst_canvas, e.position, e.depth)
	}
}

func (e *Element) drawChildren() {
	if e.ChildCount() == 0 {
		return
	}

	if !e.forceRedraw {
		// some pre-draw checks.
		// NOTE TO FUTURE BEN: this has to be done here and NOT in prepareRender() because a child might become
		// transparent between the prepare render phase and this one.
		for _, child := range e.GetChildren() {
			if !child.IsVisible() {
				continue
			}

			//check for transparent dirty children. if we find one, trigger a redraw
			if child.IsTransparent() && child.getCanvas().Dirty() {
				e.forceRedraw = true
				e.Clear()
				if e.Border.enabled {
					e.drawBorder()
				}
				break
			}

			// trigger an internal border link recalculation if a child's internal links were changed this frame.
			// this happens automatically when forcing a redraw, but otherwise we have to check for it here.
			if e.Border.enabled && child.IsBordered() && child.getBorder().internalLinksRecalculated {
				e.Border.internalLinksRecalculated = true
			}
		}
	}

	// collect opaque and transparent children, sort accordingly, and recombine into a drawlist.
	// opaque children can be drawn high to low to prevent overdraw, but transparent ones must be drawn low to high
	// like Bob Ross would.
	var opaque, transparent []element

	for _, child := range e.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		if child.IsTransparent() {
			transparent = append(transparent, child)
		} else {
			if e.forceRedraw || child.getCanvas().Dirty() {
				opaque = append(opaque, child)
			}
		}
	}

	if len(opaque) == 0 && len(transparent) == 0 {
		//nothing to draw! our work here is done.
		return
	}

	slices.SortStableFunc(opaque, func(e1, e2 element) int {
		return cmp.Compare(e2.getDepth(), e1.getDepth()) // sort by descending reverse depth
	})
	slices.SortStableFunc(transparent, func(e1, e2 element) int {
		return cmp.Compare(e1.getDepth(), e2.getDepth()) // sort by ascending depth
	})

	if !e.forceRedraw && len(transparent) > 0 {
		//if we have transparent children, we need to flatten the canvas down to the depth of the highest-depth
		//opaque child. that way opaque children that have changed will redraw successfully even if they are below
		//a transparent child.
		var maxOpaqueDepth int
		for _, child := range e.GetChildren() {
			if d := child.getDepth(); !child.IsTransparent() && d > maxOpaqueDepth {
				maxOpaqueDepth = d
			}
		}

		e.Canvas.FlattenTo(maxOpaqueDepth, e.DrawableArea())
	}

	drawlist := append(opaque, transparent...)

	// precompute cells that will need to be relinked once dirty elements are drawn. this needs to be done
	// beforehand because linking must be done after drawing, but drawing sets canvases as clean. we need to know which
	// canvases are dirty to do this properly. if this element is bordered, we ensure the cache of internal border links
	// (cells in the element's border that should be linked to) is rebuilt if necessary.
	borderLinks := e.computeBorderLinks(drawlist)

	for _, child := range drawlist {
		child.Draw(&e.Canvas, e.forceRedraw || child.IsTransparent())
	}

	if !e.Dirty() {
		return
	}

	for coord := range borderLinks.EachElement() {
		if e.IsDirtyAt(coord) { // this check may be too restrictive...
			e.linkBorderCell(coord, e.GetDepth(coord))
		}
	}
}

// Adds an animation to the ui element. Note that this does NOT start the animation.
func (e *Element) AddAnimation(animation gfx.Animator) {
	if e.animations == nil {
		e.animations = make([]gfx.Animator, 0)
	}

	//check for duplicate add
	if slices.Contains(e.animations, animation) {
		return
	}

	e.animations = append(e.animations, animation)

	//if we're adding a blocking animation during an update, make sure the window knows to stop updating
	if animation.IsBlocking() && animation.IsPlaying() {
		if wnd := e.getWindow(); wnd != nil {
			wnd.onBlockingAnimationAdded()
		}
	}
}

// AddOneShotAnimation adds an animation to the element. The animation will automatically start, play once, and then
// be removed.
func (e *Element) AddOneShotAnimation(animation gfx.Animator) {
	animation.SetOneShot(true)
	animation.Start()
	e.AddAnimation(animation)
}

func (e *Element) updateAnimations() {
	for _, a := range e.animations {
		if a.IsPlaying() {
			a.Update()
		}

		if a.JustStopped() {
			// if animation has stopped, trigger a redraw to clean up anything the animation might have left on the canvas
			e.forceRedraw = true
		}
	}

	// remove finished one-shot animations
	e.animations = slices.DeleteFunc(e.animations, func(a gfx.Animator) bool {
		return a.IsOneShot() && a.IsDone()
	})
}

func (e *Element) renderAnimations() {
	for _, animation := range e.animations {
		if animation.IsPlaying() && vec.Intersects(e.getCanvas(), animation) {
			animation.Render(&e.Canvas)
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

// Hide hides the element, preventing it and its children (if any) from receiving input, or being updated/rendered. If
// the element is focused, it loses focus.
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
		util.WalkTree[element](e, func(element element) { element.EnableListening() }, ifVisible)
	} else {
		e.setFocus(false)
		util.WalkTree[element](e, func(element element) { element.DisableListening() })
	}

	e.forceParentRedraw()
}

// Focus focuses the element, obviously. This will also defocus any focused elements in the same window.
func (e *Element) Focus() {
	e.setFocus(true)
}

// Defocus removes focus from the element.
func (e *Element) Defocus() {
	e.setFocus(false)
}

func (e *Element) setFocus(focus bool) {
	if e.focused == focus {
		return
	}

	e.focused = focus
	e.Border.dirty = true

	if window := e.getWindow(); window != nil {
		if focus {
			window.onSubNodeFocused(e.GetSelf())
		} else {
			window.onSubNodeDefocused(e.GetSelf())
		}
	}
}

func (e *Element) IsFocused() bool {
	return e.focused
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

func (e *Element) dumpUI(dir_name string, depth int) {
	if !strings.HasSuffix(dir_name, "/") {
		dir_name += "/"
	}

	if parent := e.GetParent(); parent != nil {
		e.ExportToXP(fmt.Sprintf("%s[%d] - PID %d ID %d", dir_name, depth, parent.ID(), e.ID()))
	} else {
		e.ExportToXP(fmt.Sprintf("%s[%d] - %d", dir_name, depth, e.ID()))
	}

	for _, child := range e.GetChildren() {
		child.dumpUI(dir_name, depth+1)
	}
}

type ElementID int

var id_counter int //counter for element ids
func generate_id() ElementID {
	id_counter += 1
	return ElementID(id_counter)
}

// ID returns the unique id for this element. Use this for comparisons between arbitrary elements.
func (e *Element) ID() ElementID {
	return e.id
}
