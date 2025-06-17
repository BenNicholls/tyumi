// UI is Tyumi's UI system. A base UI element is defined and successively more complex elements are composed from it. UI
// elements can be nested for composition alongside/inside one another.
package ui

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
)

// DefaultElementVisuals sets the default visuals for all UI elements. These are the visuals that an element will
// draw when cleared.
var DefaultElementVisuals gfx.Visuals

// DefaultFocusColor sets the colour for focused elements. Right now this just applies to the border of elements,
// but later on we'll use this for more advanced theming. Probably. Possibly. Oh get off my back.
var DefaultFocusColour col.Colour

// default borderstyle used by all ui elements
// NOTE: changing this does not apply retroactively to elements already created.
var DefaultBorderStyle BorderStyle

func init() {
	DefaultElementVisuals = gfx.Visuals{
		Mode:    gfx.DRAW_GLYPH,
		Colours: col.Pair{col.WHITE, col.BLACK},
	}
	DefaultFocusColour = col.PURPLE

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

// predicate for ui-tree-walking functions. we use this to break early on walks that only apply to visible sections
// of the ui tree
func ifVisible(e element) bool {
	return e != nil && e.IsVisible()
}
