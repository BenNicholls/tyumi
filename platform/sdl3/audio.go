//go:build audio

package sdl3

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/log"
)

func (p *Platform) GetAudioSystem() tyumi.AudioSystem {
	log.Error("Could not get audio system: audio not yet supported for this platform. Sorry.")

	return nil
}
