package tyumi

import "github.com/bennicholls/tyumi/log"

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

func PlayAudio(audio_resource_id ResourceID) {
	audioResource := getResource[*AudioResource](audio_resource_id)
	if audioResource == nil {
		log.Error("Could not fetch audio resource... oops.")
		return
	}

	if !audioResource.ready {
		log.Error("Audio resource not ready, has it been unloaded perhaps?")
		return
	}

	currentPlatform.PlayAudio(audioResource.platform_id)
}

func UnloadAudio(audio_resource_id ResourceID) {
	audioResource := getResource[*AudioResource](audio_resource_id)
	if audioResource == nil {
		log.Error("Could not fetch audio resource... oops.")
		return
	}

	audioResource.Unload()
}
