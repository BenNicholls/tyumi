package platform

import (
	"errors"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
)

type Platform interface {
	Init() error
	GetRenderer() gfx.Renderer
	GetEventProcessor() input.Processor
	Shutdown()
}

var current Platform = nil

func Set(p Platform) (err error) {
	if current != nil {
		current.Shutdown()
	}
	
	err = p.Init()
	if err == nil {
		current = p
	}

	return
}

// returns a renderer for the selected platform. Note that the renderer is NOT SETUP YET
// and cannot be used until renderer.Setup() is called.
func GetNewRenderer() (renderer gfx.Renderer, err error) {
	if current == nil {
		err = errors.New("Platform not set up! Use platform.Set() with your chosen platform first.")
		return
	}

	renderer = current.GetRenderer()
	return
}

func GetInputProcessor() (processor input.Processor, err error) {
	if current == nil {
		err = errors.New("Platform not set up! Use platform.Set() with your chosen platform first.")
		return
	}

	processor = current.GetEventProcessor()
	return
}
