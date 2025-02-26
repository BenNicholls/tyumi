package sdl

import (
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/veandco/go-sdl2/mix"
)

var soundCache []*mix.Chunk
var musicCache []*mix.Music

func (p *Platform) LoadSound(path string) (res_id int, err error) {
	chunk, err := mix.LoadWAV(path)
	if err != nil {
		return -1, err
	}

	soundCache = append(soundCache, chunk)

	return len(soundCache) - 1, nil
}

func (p *Platform) UnloadSound(id int) {
	if id >= len(soundCache) {
		log.Debug("BAD! Too much platform id to handle!!!")
		return
	}

	chunk := soundCache[id]
	if chunk == nil {
		return
	}

	chunk.Free()
	soundCache[id] = nil
}

func (p *Platform) PlaySound(id, channel, volume_pct int) {
	if id >= len(soundCache) {
		log.Debug("BAD! Too much platform id to handle!!!")
		return
	}

	chunk := soundCache[id]
	if chunk == nil {
		log.Debug("Don't play an unloaded sound!!!!")
		return
	}

	volume := util.Clamp(int(1.28*float64(volume_pct)), 0, 128)
	chunk.Volume(volume)
	chunk.Play(channel, 0)
}

func (p *Platform) shutdownAudio() {
	for _, chunk := range soundCache {
		if chunk != nil {
			chunk.Free()
		}
	}

	for _, music := range musicCache {
		if music != nil {
			music.Free()
		}
	}

	mix.CloseAudio()
}

func (p *Platform) LoadMusic(path string) (music_id int, err error) {
	music, err := mix.LoadMUS(path)
	if err != nil {
		return -1, err
	}

	musicCache = append(musicCache, music)

	return len(musicCache) - 1, nil
}

func (p *Platform) UnloadMusic(music_id int) {
	if music_id >= len(musicCache) {
		log.Debug("BAD! Too much platform id to handle!!!")
		return
	}

	music := musicCache[music_id]
	if music == nil {
		return
	}

	music.Free()
	musicCache[music_id] = nil
}

func (p *Platform) PlayMusic(music_id int, looping bool) {
	if music_id >= len(musicCache) {
		log.Debug("BAD! Too much id for platform to handle!")
		return
	}

	music := musicCache[music_id]
	if music == nil {
		log.Debug("Don't play unloaded music!!")
		return
	}

	if looping {
		music.Play(-1)
	} else {
		music.Play(1)
	}
}

func (p *Platform) SetMusicVolume(volume_pct int) {
	volume := util.Clamp(int(1.28*float64(volume_pct)), 0, 128)
	mix.VolumeMusic(volume)
}

func (p *Platform) PauseMusic() {
	if mix.PlayingMusic() {
		mix.PauseMusic()
	}
}

func (p *Platform) ResumeMusic() {
	if mix.PausedMusic() && mix.PlayingMusic() { // need the 2nd check because music can be halted after being paused
		mix.ResumeMusic()
	}
}

func (p *Platform) StopMusic() {
	if mix.PlayingMusic() {
		mix.HaltMusic()
	}
}
