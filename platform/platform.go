package platform

import (
	"errors"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/platform/platform_sdl"
)

type Platform int

// THINK: should there be a mix-and-match platform??
const (
	NONE Platform = iota
	SDL
	//CONSOLE
	//WEB
	//OTHER THINGS???
)

var current Platform = NONE

func Set(p Platform) {
	current = p
}

// initializes and returns a renderer for the selected platform. Note that the renderer is NOT SETUP YET
// and cannot be used until renderer.Setup() is called.
func GetNewRenderer() (renderer gfx.Renderer, err error) {
	switch current {
	case NONE:
		err = errors.New("No platform selected, cannot get renderer")
	case SDL:
		return platform_sdl.NewRenderer(), nil
	}
	return nil, errors.New("Weird platform???")
}
