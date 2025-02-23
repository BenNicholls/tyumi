package sdl

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
	"github.com/veandco/go-sdl2/mix"
)

type Platform struct {
	renderer *Renderer

	mouse_position vec.Coord
}

func (p *Platform) Init() (err error) {
	p.renderer = NewRenderer()

	err = mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		log.Error("SDL: Could not set up audio mixer")
	}

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
	p.shutdownAudio()
}

// Creates a platform for use by Tyumi. Pass this into engine.SetPlatform()
func New() *Platform {
	sdl_platform := new(Platform)
	return sdl_platform
}
