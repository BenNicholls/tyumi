//go:build audio

package sdl

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/veandco/go-sdl2/mix"
)

func (p *Platform) GetAudioSystem() tyumi.AudioSystem {
	// do some initialization
	err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		log.Error("SDL: Could not set up audio mixer")
		return nil
	}

	p.audio = new(AudioSystem)

	return p.audio
}

type AudioSystem struct {
	soundCache []*mix.Chunk
	musicCache []*mix.Music
}

func (as *AudioSystem) LoadSound(path string) (res_id int, err error) {
	chunk, err := mix.LoadWAV(path)
	if err != nil {
		return -1, err
	}

	as.soundCache = append(as.soundCache, chunk)

	return len(as.soundCache) - 1, nil
}

func (as *AudioSystem) UnloadSound(id int) {
	if id >= len(as.soundCache) {
		log.Debug("BAD! Too much platform id to handle!!!")
		return
	}

	chunk := as.soundCache[id]
	if chunk == nil {
		return
	}

	chunk.Free()
	as.soundCache[id] = nil
}

func (as *AudioSystem) PlaySound(id, channel, volume_pct int) {
	if id >= len(as.soundCache) {
		log.Debug("BAD! Too much platform id to handle!!!")
		return
	}

	chunk := as.soundCache[id]
	if chunk == nil {
		log.Debug("Don't play an unloaded sound!!!!")
		return
	}

	volume := util.Clamp(int(1.28*float64(volume_pct)), 0, 128)
	chunk.Volume(volume)
	chunk.Play(channel, 0)
}

func (as *AudioSystem) Shutdown() {
	for _, chunk := range as.soundCache {
		if chunk != nil {
			chunk.Free()
		}
	}

	for _, music := range as.musicCache {
		if music != nil {
			music.Free()
		}
	}

	mix.CloseAudio()
	log.Info("SDL Audio System shut down!")
}

func (as *AudioSystem) LoadMusic(path string) (music_id int, err error) {
	music, err := mix.LoadMUS(path)
	if err != nil {
		return -1, err
	}

	as.musicCache = append(as.musicCache, music)

	return len(as.musicCache) - 1, nil
}

func (as *AudioSystem) UnloadMusic(music_id int) {
	if music_id >= len(as.musicCache) {
		log.Debug("BAD! Too much platform id to handle!!!")
		return
	}

	music := as.musicCache[music_id]
	if music == nil {
		return
	}

	music.Free()
	as.musicCache[music_id] = nil
}

func (as *AudioSystem) PlayMusic(music_id int, looping bool) {
	if music_id >= len(as.musicCache) {
		log.Debug("BAD! Too much id for platform to handle!")
		return
	}

	music := as.musicCache[music_id]
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

func (as *AudioSystem) SetMusicVolume(volume_pct int) {
	volume := util.Clamp(int(1.28*float64(volume_pct)), 0, 128)
	mix.VolumeMusic(volume)
}

func (as *AudioSystem) PauseMusic() {
	if mix.PlayingMusic() {
		mix.PauseMusic()
	}
}

func (as *AudioSystem) ResumeMusic() {
	if mix.PausedMusic() && mix.PlayingMusic() { // need the 2nd check because music can be halted after being paused
		mix.ResumeMusic()
	}
}

func (as *AudioSystem) StopMusic() {
	if mix.PlayingMusic() {
		mix.HaltMusic()
	}
}
