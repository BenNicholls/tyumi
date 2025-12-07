//go:build !audio

package sdl3

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/log"
)

func (p *Platform) GetAudioSystem() tyumi.AudioSystem {
	log.Error("Could not get audio system: audio system not enabled during build. (Add build tag 'audio' to compile sdl platform with audio support.)")

	return nil
}
