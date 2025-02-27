package tyumi

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
)

var audioSystem AudioSystem

var masterVolume float64 = 1
var sfxVolume float64 = 1
var musicVolume float64 = 1

func EnableAudio() {
	if currentPlatform == nil {
		log.Error("Cannot enable audio: platform not set.")
		return
	}

	audioSystem = currentPlatform.GetAudioSystem()
	if audioSystem != nil {
		log.Info("Audio system enabled.")
	} else {
		log.Info("Audio system not enabled: platform did not supply system.")
	}
}

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

func PlayMusic(music_resource AudioResource) {
	if !music_resource.ready {
		log.Debug("Cannot play music, music not ready.")
		return
	}

	if music_resource.audioType != AUDIO_MUSIC {
		log.Debug("Cannot play ", music_resource.name, ", not a music resource.")
		return
	}

	music_resource.Play()
}

func PauseMusic() {
	if audioSystem != nil {
		audioSystem.PauseMusic()
	}
}

func ResumeMusic() {
	if audioSystem != nil {
		audioSystem.ResumeMusic()
	}
}

func StopMusic() {
	if audioSystem != nil {
		audioSystem.StopMusic()
	}
}

type AudioType int

const (
	AUDIO_SOUND AudioType = iota
	AUDIO_MUSIC
)

// AudioResource describes a loaded sound.
type AudioResource struct {
	Looping     bool    // used for music. if true, music loops until stopped.

	audioType   AudioType
	channel     int     // channel to play on. only used for sounds, music gets its own special channel.
	volume      float64 // volume of the sound, from [0 - 1]
	name        string  // sound name, by default this is the file name (minus extension)
	platform_id int
	ready       bool
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

func (ar AudioResource) Play() {
	if audioSystem == nil {
		return
	}
	
	if !ar.ready {
		log.Error("Audio resource not ready, has it been unloaded perhaps?")
		return
	}

	switch ar.audioType {
	case AUDIO_SOUND:
		mixedVolume := masterVolume * sfxVolume * ar.volume
		audioSystem.PlaySound(ar.platform_id, ar.channel, int(mixedVolume*100))
	case AUDIO_MUSIC:
		mixedVolume := masterVolume * musicVolume * ar.volume
		audioSystem.SetMusicVolume(int(mixedVolume * 100))
		audioSystem.PlayMusic(ar.platform_id, ar.Looping)
	}
}

func (ar AudioResource) Ready() bool {
	return ar.ready
}

func (ar *AudioResource) Unload() {
	if audioSystem == nil || !ar.ready{
		return
	}

	switch ar.audioType {
	case AUDIO_SOUND:
		audioSystem.UnloadSound(ar.platform_id)
	case AUDIO_MUSIC:
		audioSystem.UnloadMusic(ar.platform_id)
	}

	ar.ready = false
}

func LoadSound(path string) (sound AudioResource) {
	log.Info("Loading sound at ", path)
	return loadAudioResource(path, AUDIO_SOUND)
}

func LoadMusic(path string) (music AudioResource) {
	log.Info("Loading music at ", path)
	return loadAudioResource(path, AUDIO_MUSIC)
}

// LoadAudioResource loads a WAV file at path and if successful returns a playable audio resource. If not successfully
// loaded audio_resource will be nil.
func loadAudioResource(path string, audio_type AudioType) (audio_resource AudioResource) {
	if audioSystem == nil {
		log.Error("Could not load audio at", path, "audio system not enabled.")
		return
	}

	var platformID int
	var err error

	switch audio_type {
	case AUDIO_SOUND:
		platformID, err = audioSystem.LoadSound(path)
	case AUDIO_MUSIC:
		platformID, err = audioSystem.LoadMusic(path)
	default:
		return
	}

	if err != nil {
		log.Error("Could not load audio: ", err)
		return
	}

	audio_resource.platform_id = platformID
	audio_resource.audioType = audio_type
	audio_resource.volume = 1
	audio_resource.name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	audio_resource.ready = true

	return
}

type SoundLibrary struct {
	names  map[string]int // map of names to index
	sounds []AudioResource
}

func (sl *SoundLibrary) AddSound(audio_resource AudioResource) {
	if !audio_resource.ready {
		log.Debug("Sound not added to library. Not ready.")
		return
	}

	sl.sounds = append(sl.sounds, audio_resource)

	if sl.names == nil {
		sl.names = make(map[string]int)
	}

	sl.names[audio_resource.name] = len(sl.sounds) - 1
}

// Returns a reference to a sound in the library. If the name is invalid, the resource will be nil.
func (sl *SoundLibrary) Get(sound_name string) *AudioResource {
	if i, ok := sl.names[sound_name]; ok {
		return &sl.sounds[i]
	} else {
		log.Debug("No sound called ", sound_name)
		return nil
	}
}

// Plays a sound! If sound_name is invalid, does nothing.
func (sl *SoundLibrary) Play(sound_name string) {
	sound := sl.Get(sound_name)
	if sound == nil || !sound.ready {
		return
	}

	sound.Play()
}

// Plays a random sound from the library. Optionally you can provide a list of sound names to randomize between.
func (sl *SoundLibrary) PlayRandom(sound_names ...string) {
	if !sl.containsReadySounds() {
		return
	}

	switch len(sound_names) {
	case 0:
		for { // need to loop here just in case we randomly land on an unloaded or otherwise invalid sound
			if sound := util.PickOne(sl.sounds); sound.ready {
				sound.Play()
				break
			}
		}
	case 1:
		sl.Play(sound_names[0])
	default:
		// go through user-provided names, extract those which are valid and point to ready sounds.
		sounds := make([]*AudioResource, 0)
		for _, name := range sound_names {
			if sound := sl.Get(name); sound != nil && sound.ready {
				sounds = append(sounds, sound)
			}
		}

		util.PickOne(sounds).Play()
	}
}

func (sl SoundLibrary) containsReadySounds() bool {
	for _, sound := range sl.sounds {
		if sound.ready {
			return true
		}
	}

	return false
}

// LoadSoundLibrary loads all sounds in a provided directory and returns a library of those sounds.
func LoadSoundLibrary(dir_path string) (library SoundLibrary) {
	dir, err := os.ReadDir(dir_path)
	if err != nil {
		log.Error("Could not load sound library:", err)
		return
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".wav") {
			res := LoadSound(filepath.Join(dir_path, file.Name()))
			if res.ready {
				library.AddSound(res)
			}
		}
	}

	return
}
