package sdl

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/vec"
)

type Platform struct {
	renderer *Renderer

	mouse_position vec.Coord
}

func (p *Platform) Init() (err error) {
	p.renderer = NewRenderer()
	return
}

func (p *Platform) GetRenderer() tyumi.Renderer {
	return p.renderer
}

func (p *Platform) GetEventGenerator() tyumi.EventGenerator {
	return p.processEvents
}

func (p *Platform) Shutdown() {
	p.renderer.Cleanup()
}

// Creates a platform for use by Tyumi. Pass this into engine.SetPlatform()
func New() *Platform {
	sdl_platform := new(Platform)
	return sdl_platform
}
