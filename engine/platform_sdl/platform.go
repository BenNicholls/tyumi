package platform_sdl

import (
	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/vec"
)

type SDLPlatform struct {
	renderer *SDLRenderer

	mouse_position vec.Coord
}

func (sdlp *SDLPlatform) Init() (err error) {
	sdlp.renderer = NewRenderer()
	return
}

func (sdlp *SDLPlatform) GetRenderer() engine.Renderer {
	return sdlp.renderer
}

func (sdlp *SDLPlatform) GetEventGenerator() engine.EventGenerator {
	return sdlp.processEvents
}

func (sdlp *SDLPlatform) Shutdown() {
	sdlp.renderer.Cleanup()
}

func New() *SDLPlatform {
	sdlp := new(SDLPlatform)
	return sdlp
}
