// UI is Tyumi's UI system. A base UI element is defined and successively more complex elements are composed from it. UI
// elements can be nested for composition alongside/inside one another.
package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
)

// default visuals used by all ui elements
var defaultCanvasVisuals gfx.Visuals

// SetDefaultElementVisuals sets the default visuals for all UI elements. These are the visuals that an element will
// draw when cleared.
// NOTE: this does not apply retroactively to elements already created.
// THINK: this could emit an event that could alert windows/elements to reset their default visuals to the new one,
// though this might be better implemented as a sort of theming functionality someday
func SetDefaultElementVisuals(vis gfx.Visuals) {
	defaultCanvasVisuals = vis
}

// default colour for focused element borders
var defaultFocusColour uint32

// SetDefaultFocusColor sets the colour for focused elements. Right now this just applies to the border of elements,
// but later on we'll use this for more advanced theming. Probably. Possibly. Oh get off my back.
func SetDefaultFocusColour(colour uint32) {
	defaultFocusColour = colour
}

// default borderstyle used by all ui elements
// NOTE: changing this does not apply retroactively to elements already created.
var DefaultBorderStyle BorderStyle

func init() {
	defaultCanvasVisuals = gfx.Visuals{
		Mode:    gfx.DRAW_GLYPH,
		Colours: col.Pair{col.WHITE, col.BLACK},
	}
	defaultFocusColour = col.PURPLE

	createBorderStyles()
	DefaultBorderStyle = BorderStyles["Thin"]
}

// Retrieves a reference to the element in window with the supplied label. If the element is not found, or is not
// right type, returns nil.
func GetLabelled[T element](window *Window, label string) (element T) {
	if e, ok := window.labels[label]; ok {
		if t, ok := e.(T); ok {
			return t
		}
	}

	return
}

func fireCallbacks(callbacks ...func()) {
	for _, callback := range callbacks {
		if callback != nil {
			callback()
		}
	}
}
