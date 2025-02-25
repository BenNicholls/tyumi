package tyumi

import (
	"os"
	"strings"

	"github.com/bennicholls/tyumi/log"
)

type AudioResource struct {
	Resource
}

func (ar *AudioResource) Unload() {
	if !ar.ready {
		return
	}

	currentPlatform.UnloadAudio(ar.platform_id)
	ar.Resource.Unload()
}

func LoadAudioResource(path string) (resource_id ResourceID) {
	if currentPlatform == nil {
		log.Error("Could not load audio at", path, "platform not set up yet.")
		return invalidResource
	}

	if id := getResourceIDByPath(path); id != resourceNotFound {
		log.Debug("resource already loaded!!")
		return id
	}

	platformID, err := currentPlatform.LoadAudio(path)
	if err != nil {
		log.Error("Could not load audio: ", err)
		return
	}

	res := AudioResource{
		Resource: Resource{
			platform_id: platformID,
			ready:       true,
			path:        path,
		},
	}

	return addResourceToCache(&res)
}

func PlayAudio(audio_resource_id ResourceID, channel int) {
	if audio_resource_id == invalidResource {
		log.Debug("Could not play audio, invalid resource ID.")
		return
	}

	audioResource := getResource[*AudioResource](audio_resource_id)
	if audioResource == nil {
		log.Error("Could not fetch audio resource...")
		return
	}

	if !audioResource.ready {
		log.Error("Audio resource not ready, has it been unloaded perhaps?")
		return
	}

	currentPlatform.PlayAudio(audioResource.platform_id, channel)
}

func UnloadAudio(audio_resource_id ResourceID) {
	if audio_resource_id == invalidResource {
		return
	}

	audioResource := getResource[*AudioResource](audio_resource_id)
	if audioResource == nil {
		log.Error("Could not fetch audio resource...")
		return
	}

	audioResource.Unload()
}

// LoadSoundLibrary loads all sounds in a provided directory and returns a map whose keys are the filenames
// of the sounds (minus extension) and the values are the loaded sounds' associated ResourceIDs. So for example
// the sound "beep.wav" will be named beep. If dir_path is invalid, library will be nil.
func LoadSoundLibrary(dir_path string) (library map[string]ResourceID) {
	dir, err := os.ReadDir(dir_path)
	if err != nil {
		log.Error("Could not load sound library:", err)
	}

	library = make(map[string]ResourceID)

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".wav") {
			id := LoadAudioResource(dir_path + "/" + file.Name())
			key := strings.TrimSuffix(file.Name(), ".wav")
			library[key] = id
		}
	}

	return
}
