//go:build js && wasm

package web

import (
	"syscall/js"

	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

func (p *Platform) listenWindow() {
	win := js.Global().Get("window")
	p.listen(win, "keydown", func(ev js.Value) { p.onKey(ev, true) })
	p.listen(win, "keyup", func(ev js.Value) { p.onKey(ev, false) })
	p.listen(win, "mousemove", p.onMouseMove)
	p.listen(win, "mousedown", func(ev js.Value) {
		p.unlockAudio()
		p.focusCanvas()
		p.onMouseMove(ev)
	})
	p.listen(win, "blur", func(js.Value) { p.releaseAllKeys() })
	p.listen(win, "resize", func(js.Value) { p.renderer.forceRedraw = true })
	p.listen(win, "contextmenu", func(ev js.Value) {
		if p.renderer.canvas.Truthy() && ev.Get("target").Equal(p.renderer.canvas) {
			ev.Call("preventDefault")
		}
	})
}

func (p *Platform) focusCanvas() {
	if p.renderer.canvas.Truthy() {
		p.renderer.canvas.Call("focus")
	}
}

func (p *Platform) onKey(ev js.Value, down bool) {
	code := ev.Get("code").String()
	keyName := ev.Get("key").String()
	if !browserShortcut(ev, code, keyName) {
		ev.Call("preventDefault")
	}

	key, ok := KeycodeFromJS(code, keyName)
	if !ok {
		return
	}

	var mods input.KeyModifiers
	if ev.Get("ctrlKey").Bool() || ev.Get("metaKey").Bool() {
		mods |= input.KEYMOD_CTRL
	}
	if ev.Get("altKey").Bool() {
		mods |= input.KEYMOD_ALT
	}
	if ev.Get("shiftKey").Bool() {
		mods |= input.KEYMOD_SHIFT
	}

	p.unlockAudio()

	p.mu.Lock()
	defer p.mu.Unlock()
	if down {
		repeat := ev.Get("repeat").Bool()
		p.keys = append(p.keys, queuedKey{key: key, mods: mods, repeat: repeat})
		if !repeat {
			p.held[key] = mods
		}
		return
	}
	delete(p.held, key)
	p.keys = append(p.keys, queuedKey{key: key, mods: mods, released: true})
}

func browserShortcut(ev js.Value, code, key string) bool {
	if code == "F5" || code == "F11" || code == "F12" || ev.Get("metaKey").Bool() {
		return true
	}
	if ev.Get("ctrlKey").Bool() && (key == "r" || key == "R" || key == "l" || key == "L") {
		return true
	}
	return false
}

func (p *Platform) onMouseMove(ev js.Value) {
	if !ev.Truthy() || ev.Get("clientX").IsUndefined() {
		return
	}
	p.mu.Lock()
	p.pointerX = ev.Get("clientX").Float()
	p.pointerY = ev.Get("clientY").Float()
	p.pointerMoved = true
	p.mu.Unlock()
}

func (p *Platform) releaseAllKeys() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for key, mods := range p.held {
		p.keys = append(p.keys, queuedKey{key: key, mods: mods, released: true})
	}
	clear(p.held)
}

func (p *Platform) GenerateEvents() {
	p.mu.Lock()
	keys := p.keys
	p.keys = nil
	x, y, moved := p.pointerX, p.pointerY, p.pointerMoved
	p.pointerMoved = false
	p.mu.Unlock()

	for _, k := range keys {
		switch {
		case k.released:
			input.FireKeyReleaseEvent(k.key, k.mods)
		case k.repeat:
			input.FireKeyRepeatEvent(k.key, k.mods)
		default:
			input.FireKeyPressEvent(k.key, k.mods)
		}
	}

	if !moved {
		return
	}
	cell, ok := p.renderer.cellAtClient(x, y)
	if !ok || cell == p.mouse {
		return
	}
	input.FireMouseMoveEvent(cell, cell.Subtract(p.mouse))
	p.mouse = cell
}

func (r *Renderer) cellAtClient(clientX, clientY float64) (vec.Coord, bool) {
	if !r.canvas.Truthy() || r.tileSize <= 0 {
		return vec.Coord{}, false
	}
	rect := r.canvas.Call("getBoundingClientRect")
	width := rect.Get("width").Float()
	height := rect.Get("height").Float()
	canvasW := r.canvas.Get("width").Float()
	canvasH := r.canvas.Get("height").Float()
	if width <= 0 || height <= 0 || canvasW <= 0 || canvasH <= 0 {
		return vec.Coord{}, false
	}
	x := (clientX - rect.Get("left").Float()) * canvasW / width
	y := (clientY - rect.Get("top").Float()) * canvasH / height
	if x < 0 || y < 0 || x >= canvasW || y >= canvasH {
		return vec.Coord{}, false
	}
	return vec.Coord{int(x) / r.tileSize, int(y) / r.tileSize}, true
}
