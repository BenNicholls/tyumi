package platform_sdl

import (
	"github.com/bennicholls/tyumi/platform"
	"github.com/bennicholls/tyumi/vec"
)

type SDLPlatform struct {
	renderer *SDLRenderer

}

func (sdlp *SDLPlatform) Init() (err error) {
	sdlp.renderer = NewRenderer()
	return
}

func (sdlp *SDLPlatform) GetRenderer() platform.Renderer {
	return sdlp.renderer
}

func (sdlp *SDLPlatform) GetEventGenerator() platform.EventGenerator {
	return sdlp.processEvents
}

func (sdlp *SDLPlatform) Shutdown() {

}

func New() *SDLPlatform {
	sdlp := new(SDLPlatform)
	return sdlp
}
