package platform

import (
	"errors"

	"github.com/bennicholls/tyumi/event"
)

var EV_QUIT = event.Register("Quit Event")

type Platform interface {
	Init() error
	GetRenderer() Renderer
	GetEventGenerator() EventGenerator
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
func GetNewRenderer() (renderer Renderer, err error) {
	if current == nil {
		err = errors.New("Platform not set up! Use platform.Set() with your chosen platform first.")
		return
	}

	renderer = current.GetRenderer()
	return
}

func GetEventGenerator() (generator EventGenerator, err error) {
	if current == nil {
		err = errors.New("Platform not set up! Use platform.Set() with your chosen platform first.")
		return
	}

	generator = current.GetEventGenerator()
	return
}
