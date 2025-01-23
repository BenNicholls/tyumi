package gfx

type Renderer interface {
	Setup(window *Canvas, glyphPath, fontPath, title string) error
	Ready() bool
	Cleanup()
	ChangeFonts(glyphPath, fontPath string) error
	SetFullscreen(bool)
	ToggleFullscreen()
	Render()
	ForceRedraw()
	ToggleDebugMode(string)
}
