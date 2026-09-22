//go:build js && wasm

package web

import (
	"fmt"
	"math"
	"syscall/js"
	"time"

	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
)

// AudioSystem plays sound and music through the Web Audio API.
// The context stays suspended until the player presses a key or clicks, which
// is when browsers allow audio to start.
type AudioSystem struct {
	ctx         js.Value
	musicGain   js.Value
	sounds      []js.Value
	music       []js.Value
	channels    map[int]js.Value
	musicNode   js.Value
	musicID     int
	musicLoop   bool
	musicOn     bool
	musicPaused bool
	musicStart  time.Time
	musicOff    float64
	pauseAt     float64
}

func newAudioSystem() (*AudioSystem, error) {
	ctor := js.Global().Get("AudioContext")
	if ctor.Type() != js.TypeFunction {
		ctor = js.Global().Get("webkitAudioContext")
	}
	if ctor.Type() != js.TypeFunction {
		err := fmt.Errorf("this browser has no Web Audio API")
		log.Error("WEB AUDIO: ", err)
		return nil, err
	}

	var ctx js.Value
	err := func() (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				err = fmt.Errorf("%v", rec)
			}
		}()
		ctx = ctor.New()
		return nil
	}()
	if err != nil {
		log.Error("WEB AUDIO: Could not create audio context: ", err)
		return nil, err
	}

	as := &AudioSystem{
		ctx:      ctx,
		channels: make(map[int]js.Value),
	}
	as.musicGain = ctx.Call("createGain")
	as.musicGain.Call("connect", ctx.Get("destination"))
	log.Info("Web audio system enabled.")
	return as, nil
}

func (as *AudioSystem) Unlock() {
	if as == nil || !as.ctx.Truthy() {
		return
	}
	if as.ctx.Get("state").String() == "suspended" {
		as.ctx.Call("resume")
	}
}

func (as *AudioSystem) LoadSound(path string) (int, error) {
	buf, err := as.decode(path)
	if err != nil {
		return -1, err
	}
	as.sounds = append(as.sounds, buf)
	return len(as.sounds) - 1, nil
}

func (as *AudioSystem) UnloadSound(id int) {
	if id < 0 || id >= len(as.sounds) {
		return
	}
	as.sounds[id] = js.Null()
}

func (as *AudioSystem) PlaySound(id, channel, volumePct int) {
	if id < 0 || id >= len(as.sounds) || !as.sounds[id].Truthy() {
		return
	}
	as.Unlock()
	src := as.ctx.Call("createBufferSource")
	src.Set("buffer", as.sounds[id])
	gain := as.ctx.Call("createGain")
	gain.Get("gain").Set("value", float64(util.Clamp(volumePct, 0, 100))/100)
	src.Call("connect", gain)
	gain.Call("connect", as.ctx.Get("destination"))
	if channel >= 0 {
		stopAudioNode(as.channels[channel])
		as.channels[channel] = src
	}
	src.Call("start")
}

func (as *AudioSystem) LoadMusic(path string) (int, error) {
	buf, err := as.decode(path)
	if err != nil {
		return -1, err
	}
	as.music = append(as.music, buf)
	return len(as.music) - 1, nil
}

func (as *AudioSystem) UnloadMusic(id int) {
	if id < 0 || id >= len(as.music) {
		return
	}
	if as.musicOn && as.musicID == id {
		as.StopMusic()
	}
	as.music[id] = js.Null()
}

func (as *AudioSystem) PlayMusic(id int, looping bool) {
	if id < 0 || id >= len(as.music) || !as.music[id].Truthy() {
		return
	}
	as.Unlock()
	as.stopMusicNode()
	as.startMusic(id, looping, 0)
}

func (as *AudioSystem) SetMusicVolume(volumePct int) {
	if !as.musicGain.Truthy() {
		return
	}
	as.musicGain.Get("gain").Set("value", float64(util.Clamp(volumePct, 0, 100))/100)
}

func (as *AudioSystem) PauseMusic() {
	if !as.musicOn || as.musicPaused {
		return
	}
	elapsed := time.Since(as.musicStart).Seconds() + as.musicOff
	if dur := as.duration(as.musicID); dur > 0 {
		if as.musicLoop {
			elapsed = math.Mod(elapsed, dur)
		} else if elapsed >= dur {
			as.StopMusic()
			return
		}
	}
	as.pauseAt = elapsed
	as.stopMusicNode()
	as.musicPaused = true
	as.musicOn = true
}

func (as *AudioSystem) ResumeMusic() {
	if as.musicPaused && as.musicOn {
		as.startMusic(as.musicID, as.musicLoop, as.pauseAt)
	}
}

func (as *AudioSystem) StopMusic() {
	as.stopMusicNode()
	as.musicOn = false
	as.musicPaused = false
}

func (as *AudioSystem) Shutdown() {
	as.StopMusic()
	if as.ctx.Truthy() && as.ctx.Get("close").Type() == js.TypeFunction {
		as.ctx.Call("close")
	}
	log.Info("Web audio system shut down.")
}

func (as *AudioSystem) decode(path string) (js.Value, error) {
	data, err := ReadAsset(path)
	if err != nil {
		return js.Undefined(), err
	}
	buf := js.Global().Get("ArrayBuffer").New(len(data))
	js.CopyBytesToJS(js.Global().Get("Uint8Array").New(buf), data)
	decoded, err := await(as.ctx.Call("decodeAudioData", buf))
	if err != nil {
		return js.Undefined(), fmt.Errorf("decode %s: %w", path, err)
	}
	return decoded, nil
}

func (as *AudioSystem) startMusic(id int, looping bool, offset float64) {
	src := as.ctx.Call("createBufferSource")
	src.Set("buffer", as.music[id])
	src.Set("loop", looping)
	src.Call("connect", as.musicGain)
	if offset > 0 {
		src.Call("start", 0, offset)
	} else {
		src.Call("start")
	}
	as.musicNode = src
	as.musicID = id
	as.musicLoop = looping
	as.musicOn = true
	as.musicPaused = false
	as.musicStart = time.Now()
	as.musicOff = offset
}

func (as *AudioSystem) stopMusicNode() {
	stopAudioNode(as.musicNode)
	as.musicNode = js.Null()
}

func (as *AudioSystem) duration(id int) float64 {
	if id < 0 || id >= len(as.music) || !as.music[id].Truthy() {
		return 0
	}
	return as.music[id].Get("duration").Float()
}

func stopAudioNode(node js.Value) {
	if !node.Truthy() {
		return
	}
	defer func() { recover() }()
	node.Call("stop")
}
