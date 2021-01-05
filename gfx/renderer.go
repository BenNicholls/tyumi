package gfx

type Renderer interface {
	Setup(string, string, string) error
	Ready() bool
	Cleanup()
	ChangeFonts(string, string) error
	SetFullscreen(bool)
	ToggleFullscreen()
	SetFramerate(int)
	Render()
	ForceRedraw()
	ToggleDebugMode(string)
}