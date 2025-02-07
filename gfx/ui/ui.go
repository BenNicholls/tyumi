// UI is Tyumi's UI system. A base UI element is defined and successively more complex elements are composed from it. UI
// elements can be nested for composition alongside/inside one another.
package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
)

var defaultCanvasVisuals gfx.Visuals = gfx.Visuals{
	Mode: gfx.DRAW_GLYPH,
	Colours: col.Pair{col.WHITE, col.BLACK},
}

// SetDefaultElementVisuals sets the default visuals for all UI elements. These are the visuals that an element will 
// draw when cleared. This will not affect elements that have already been created.
// THINK: this could emit an event that could alert windows/elements to reset their default visuals to the new one,
// though this might be better implemented as a sort of theming functionality someday
func SetDefaultElementVisuals(vis gfx.Visuals) {
	defaultCanvasVisuals = vis
}

// Retrieves a reference to the element in window with the supplied label. If the element is not found, or is not
// right type, returns nil.
func GetLabelledElement[T Element](window *Window, label string) (element T) {
	if e, ok := window.labels[label]; ok {
		if t, ok := e.(T); ok {
			return t
		}
	}

	return
}
