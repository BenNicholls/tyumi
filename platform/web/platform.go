//go:build js && wasm

package web

import (
	"sync"
	"syscall/js"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

var (
	_ tyumi.Platform    = (*Platform)(nil)
	_ tyumi.Renderer    = (*Renderer)(nil)
	_ tyumi.AudioSystem = (*AudioSystem)(nil)
)

type listener struct {
	target js.Value
	event  string
	fn     js.Func
}

type queuedKey struct {
	key      input.Keycode
	mods     input.KeyModifiers
	released bool
	repeat   bool
}

// Platform draws to a canvas, reads DOM input, and plays audio with Web Audio.
type Platform struct {
	renderer Renderer
	audio    *AudioSystem

	mouse vec.Coord

	mu           sync.Mutex
	keys         []queuedKey
	held         map[input.Keycode]input.KeyModifiers
	pointerX     float64
	pointerY     float64
	pointerMoved bool
	listeners    []listener
}

func (p *Platform) Init() error {
	gfx.DefaultTextMode = gfx.TEXTMODE_HALF
	p.held = make(map[input.Keycode]input.KeyModifiers)
	p.listenWindow()
	return nil
}

func (p *Platform) ChangeTitle(title string) {
	doc := js.Global().Get("document")
	if doc.Truthy() {
		doc.Set("title", title)
	}
}

func (p *Platform) GetRenderer() tyumi.Renderer {
	return &p.renderer
}

func (p *Platform) GetAudioSystem() tyumi.AudioSystem {
	if p.audio != nil {
		return p.audio
	}
	audio, err := newAudioSystem()
	if err != nil {
		return nil
	}
	p.audio = audio
	return p.audio
}

func (p *Platform) Shutdown() {
	for _, l := range p.listeners {
		l.target.Call("removeEventListener", l.event, l.fn)
		l.fn.Release()
	}
	p.listeners = nil
	p.renderer.Cleanup()
	if p.audio != nil {
		p.audio.Shutdown()
	}
}

// NewPlatform creates a browser platform. Pass it to tyumi.SetPlatform before
// initializing the console.
func NewPlatform() *Platform {
	return new(Platform)
}

func (p *Platform) listen(target js.Value, event string, fn func(js.Value)) {
	wrapped := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			fn(args[0])
		} else {
			fn(js.Undefined())
		}
		return nil
	})
	target.Call("addEventListener", event, wrapped)
	p.listeners = append(p.listeners, listener{target: target, event: event, fn: wrapped})
}

func (p *Platform) unlockAudio() {
	if p.audio != nil {
		p.audio.Unlock()
	}
}
