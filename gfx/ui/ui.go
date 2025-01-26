// UI is Tyumi's UI system. A base UI element is defined and successively more complex elements are composed from it. UI
// elements can be nested for composition alongside/inside one another.
package ui

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
