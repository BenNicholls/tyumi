package platform_sdl

import "github.com/bennicholls/tyumi/gfx"
import "github.com/bennicholls/tyumi/input"

type SDLPlatform struct {
	renderer gfx.Renderer
	event_processor input.Processor
}

func (sdlp *SDLPlatform) Init() (err error) {
	sdlp.renderer = NewRenderer()
	sdlp.event_processor = processEvents
	return
}

func (sdlp *SDLPlatform) GetRenderer() gfx.Renderer {
	return sdlp.renderer
}

func (sdlp *SDLPlatform) GetEventProcessor() input.Processor {
	return sdlp.event_processor
}

func (sdlp *SDLPlatform) Shutdown() {

}

func New() *SDLPlatform {
	sdlp := new(SDLPlatform)
	return sdlp
}