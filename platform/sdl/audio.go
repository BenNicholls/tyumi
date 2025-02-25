package sdl

import (
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/veandco/go-sdl2/mix"
)

var audioCache []*mix.Chunk

func (p *Platform) LoadAudio(path string) (res_id int, err error) {
	chunk, err := mix.LoadWAV(path)
	if err != nil {
		return -1, err
	}

	audioCache = append(audioCache, chunk)

	return len(audioCache) - 1, nil
}

func (p *Platform) UnloadAudio(id int) {
	if id >= len(audioCache) {
		log.Debug("BAD! Too much platform id to handle!!!")
		return
	}

	chunk := audioCache[id]
	if chunk == nil {
		return
	}

	chunk.Free()
	audioCache[id] = nil
}

func (p *Platform) PlayAudio(id, channel, volume_pct int) {
	if id >= len(audioCache) {
		log.Debug("BAD! Too much platform id to handle!!!")
		return
	}

	chunk := audioCache[id]
	if chunk == nil {
		log.Debug("Don't play an unloaded sound!!!!")
		return
	}

	volume := util.Clamp(int(1.28*float64(volume_pct)), 0, 128)
	chunk.Volume(volume)
	chunk.Play(channel, 0)
}

func (p *Platform) shutdownAudio() {
	for _, chunk := range audioCache {
		if chunk != nil {
			chunk.Free()
		}
	}

	mix.CloseAudio()
}
