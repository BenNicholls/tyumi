package sdl

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/vec"
)

type Platform struct {
	renderer Renderer
	audio    tyumi.AudioSystem

	mouse_position vec.Coord
}

func (p *Platform) Init() (err error) {
	gfx.DefaultTextMode = gfx.TEXTMODE_HALF

	return
}

func (p *Platform) GetRenderer() tyumi.Renderer {
	return &p.renderer
}

func (p *Platform) GetEventGenerator() tyumi.EventGenerator {
	return p.processEvents
}

func (p *Platform) Shutdown() {
	p.renderer.Cleanup()

	if p.audio != nil {
		p.audio.Shutdown()
	}
}

// Creates a platform for use by Tyumi. Pass this into engine.SetPlatform()
func NewPlatform() *Platform {
	sdl_platform := new(Platform)
	return sdl_platform
}
