package tyumi

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
)

var EV_QUIT = event.Register("Quit Event", event.SIMPLE)

// Platform defines the API for the platform-specific code that Tyumi uses to interface with the system. It's split
// into a number of subsystems, some of which are optional. At the moment, everything except the AudioSystem is required.
type Platform interface {
	Init() error
	Shutdown()

	ChangeTitle(title string)

	GetRenderer() Renderer
	GetEventGenerator() EventGenerator
	GetAudioSystem() AudioSystem
}

var currentPlatform Platform = nil

// Sets the platform to be used by Tyumi for rendering, gathering of system events, and more. This must be called
// before console initialization or running the game loop. The engine will Init() the platform for you.
func SetPlatform(p Platform) (err error) {
	err = p.Init()
	if err != nil {
		log.Error("Could not initialize platform: ", err)
		return
	}

	if currentPlatform != nil {
		log.Info("Shutting down old platform.")
		currentPlatform.Shutdown()
	}

	currentPlatform = p
	renderer = p.GetRenderer()
	eventGenerator = p.GetEventGenerator()

	return
}

// definition of whatever system is grabbing events from the system
// the expectation is that this function consumes system-level events and converts them into tyumi
// events and fires them
type EventGenerator func()

// definition of whatever system is rendering to the screen
type Renderer interface {
	Setup(console *gfx.Canvas, glyphPath, fontPath, title string) error
	Ready() bool
	Cleanup()
	ChangeFonts(glyphPath, fontPath string) error
	SetFullscreen(bool)
	SetClearColour(colour col.Colour)
	ToggleFullscreen()
	Render()
	ForceRedraw()
	ToggleDebugMode(string)
}

type AudioSystem interface {
	LoadSound(path string) (platform_audio_id int, err error)
	UnloadSound(platform_audio_id int)
	PlaySound(platform_audio_id, channel, volume_pct int)

	LoadMusic(path string) (platform_music_id int, err error)
	UnloadMusic(platform_music_id int)
	PlayMusic(platform_music_id int, looping bool)
	SetMusicVolume(volume_pct int)
	PauseMusic()
	ResumeMusic()
	StopMusic()

	Shutdown()
}
