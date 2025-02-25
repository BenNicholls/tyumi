package tyumi

import (
	"os"
	"strings"

	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
)

var masterVolume float64 = 1
var sfxVolume float64 = 1
var musicVolume float64 = 1

// Sets the master volume for all sounds and music. volume is a percentage [0 - 100]
func SetVolume(volume int) {
	masterVolume = util.Clamp(float64(volume)/100.0, 0, 1)
}

// Sets the volume for all sounds. volume is a percentage [0 - 100]
func SetSFXVolume(volume int) {
	sfxVolume = util.Clamp(float64(volume)/100.0, 0, 1)
}

// Sets the volume for all music. volume is a percentage [0 - 100]
func SetMusicVolume(volume int) {
	musicVolume = util.Clamp(float64(volume)/100.0, 0, 1)
}

// AudioResource describes a loaded sound.
type AudioResource struct {
	channel     int     // channel to play on
	volume      float64 // volume of the sound, from [0 - 1]
	ready       bool
	platform_id int
}

// Sets the volume for the sound. This is a percentage value between [0 - 100].
func (ar *AudioResource) SetVolume(volume int) {
	ar.volume = util.Clamp(float64(volume)/100.0, 0, 1)
}

// Sets which channel this sound should play on. Sounds played on the same channel will cut eachother off.
func (ar *AudioResource) SetChannel(channel int) {
	ar.channel = util.Clamp(channel, 0, 7)
}

// Tells the sound to play on the next available channel.
func (ar *AudioResource) SetChannelAny() {
	ar.channel = -1
}

func (ar *AudioResource) Play() {
	if !ar.ready {
		log.Error("Audio resource not ready, has it been unloaded perhaps?")
		return
	}

	mixedVolume := masterVolume * sfxVolume * ar.volume
	currentPlatform.PlayAudio(ar.platform_id, ar.channel, int(mixedVolume*100))
}

func (ar AudioResource) Ready() bool {
	return ar.ready
}

func (ar *AudioResource) Unload() {
	if !ar.ready {
		return
	}

	currentPlatform.UnloadAudio(ar.platform_id)
	ar.ready = false
}

// LoadAudioResource loads a WAV file at path and if successful returns a playable audio resource. If not successfully
// loaded audio_resource will be nil.
func LoadAudioResource(path string) (audio_resource *AudioResource) {
	log.Info("Loading audio at ", path)
	if currentPlatform == nil {
		log.Error("Could not load audio at", path, "platform not set up yet.")
		return
	}

	platformID, err := currentPlatform.LoadAudio(path)
	if err != nil {
		log.Error("Could not load audio: ", err)
		return
	}

	audio_resource = new(AudioResource)
	audio_resource.platform_id = platformID
	audio_resource.ready = true
	audio_resource.volume = 1

	return
}

// LoadSoundLibrary loads all sounds in a provided directory and returns a map whose keys are the filenames
// of the sounds (minus extension) and the values are the associated AudioResources. So for example
// the sound "beep.wav" will be named beep. If dir_path is invalid, library will be nil.
func LoadSoundLibrary(dir_path string) (library map[string]*AudioResource) {
	dir, err := os.ReadDir(dir_path)
	if err != nil {
		log.Error("Could not load sound library:", err)
	}

	library = make(map[string]*AudioResource)

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".wav") {
			res := LoadAudioResource(dir_path + "/" + file.Name())
			if res != nil {
				key := strings.TrimSuffix(file.Name(), ".wav")
				library[key] = res
			}
		}
	}

	return
}
