//go:build js && wasm

package web

import (
	"fmt"
	"image"
	"syscall/js"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

type pxRect struct {
	x, y, w, h int
}

// Renderer paints the console into a canvas with the page's 2D context.
type Renderer struct {
	canvas    js.Value
	ctx       js.Value
	buffer    js.Value
	bufferCtx js.Value

	glyphs *image.NRGBA
	font   *image.NRGBA

	glyphCache map[col.Colour]js.Value
	fontCache  map[col.Colour]js.Value

	tileSize int

	forceRedraw bool
	showChanges bool
	frames      int

	clearColour col.Colour
	debugColour col.Colour

	console *gfx.Canvas
	ready   bool
}

func (r *Renderer) Setup(console *gfx.Canvas, glyphPath, fontPath, title string) error {
	if err := r.attachCanvas(); err != nil {
		return err
	}
	r.console = console
	r.clearColour = col.BLACK
	r.applyClearColour()
	if title != "" {
		js.Global().Get("document").Set("title", title)
	}

	if err := r.ChangeFonts(glyphPath, fontPath); err != nil {
		return err
	}
	r.ready = true
	r.canvas.Call("focus")
	return nil
}

func (r *Renderer) Ready() bool {
	return r.ready
}

func (r *Renderer) Cleanup() {
	if !r.ready {
		return
	}
	r.ready = false
	r.glyphCache = nil
	r.fontCache = nil
	log.Info("Web renderer shut down.")
}

func (r *Renderer) ChangeFonts(glyphPath, fontPath string) error {
	glyphs, err := loadAtlas(glyphPath)
	if err != nil {
		log.Error("WEB RENDERER: Could not load glyphs at ", glyphPath, ": ", err)
		return err
	}
	font, err := loadAtlas(fontPath)
	if err != nil {
		log.Error("WEB RENDERER: Could not load font at ", fontPath, ": ", err)
		return err
	}

	tileSize := glyphs.Bounds().Dx() / 16
	if tileSize <= 0 {
		return fmt.Errorf("glyph image %s is too small to contain a 16x16 atlas", glyphPath)
	}

	r.glyphs = glyphs
	r.font = font
	r.glyphCache = make(map[col.Colour]js.Value)
	r.fontCache = make(map[col.Colour]js.Value)

	if tileSize != r.tileSize || !r.buffer.Truthy() {
		r.tileSize = tileSize
		r.createBuffer()
		r.forceRedraw = true
	}
	log.Info("WEB RENDERER: Loaded fonts. Glyph: ", glyphPath, ", Text: ", fontPath)
	return nil
}

func loadAtlas(path string) (*image.NRGBA, error) {
	data, err := ReadAsset(path)
	if err != nil {
		return nil, err
	}
	atlas, err := AtlasFromBMP(data)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return atlas, nil
}

func (r *Renderer) attachCanvas() error {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return fmt.Errorf("document is unavailable")
	}
	canvas := doc.Call("getElementById", CanvasID)
	if !canvas.Truthy() {
		canvas = doc.Call("createElement", "canvas")
		canvas.Set("id", CanvasID)
		style := canvas.Get("style")
		style.Set("width", "100vw")
		style.Set("height", "100vh")
		style.Set("objectFit", "contain")
		style.Set("imageRendering", "pixelated")
		style.Set("background", "#000")
		body := doc.Get("body")
		if body.Truthy() {
			body.Get("style").Set("margin", "0")
			body.Get("style").Set("background", "#000")
			body.Call("appendChild", canvas)
		}
	}
	canvas.Set("tabIndex", 0)
	if canvas.Get("style").Get("imageRendering").String() == "" {
		canvas.Get("style").Set("imageRendering", "pixelated")
	}
	r.canvas = canvas
	return nil
}

func (r *Renderer) createBuffer() {
	if r.console == nil || r.tileSize <= 0 {
		return
	}
	size := r.console.Size()
	w := size.W * r.tileSize
	h := size.H * r.tileSize

	r.buffer = js.Global().Get("document").Call("createElement", "canvas")
	r.buffer.Set("width", w)
	r.buffer.Set("height", h)
	r.bufferCtx = r.buffer.Call("getContext", "2d")
	configure2D(r.bufferCtx)

	r.canvas.Set("width", w)
	r.canvas.Set("height", h)
	r.ctx = r.canvas.Call("getContext", "2d")
	configure2D(r.ctx)
}

func (r *Renderer) SetFullscreen(enable bool) {
	doc := js.Global().Get("document")
	if enable {
		el := doc.Get("documentElement")
		if !el.Get("requestFullscreen").Truthy() {
			return
		}
		_, err := await(el.Call("requestFullscreen"))
		if err != nil {
			log.Error("WEB RENDERER: Fullscreen request failed: ", err)
		}
		return
	}
	if doc.Get("fullscreenElement").Truthy() && doc.Get("exitFullscreen").Truthy() {
		_, err := await(doc.Call("exitFullscreen"))
		if err != nil {
			log.Error("WEB RENDERER: Leaving fullscreen failed: ", err)
		}
	}
}

func (r *Renderer) ToggleFullscreen() {
	r.SetFullscreen(!js.Global().Get("document").Get("fullscreenElement").Truthy())
}

func (r *Renderer) SetClearColour(colour col.Colour) {
	r.clearColour = colour
	r.applyClearColour()
	r.forceRedraw = true
}

func (r *Renderer) applyClearColour() {
	body := js.Global().Get("document").Get("body")
	if body.Truthy() {
		body.Get("style").Set("background", cssColour(r.clearColour))
	}
	if r.canvas.Truthy() {
		r.canvas.Get("style").Set("background", cssColour(r.clearColour))
	}
}

func (r *Renderer) ForceRedraw() {
	r.forceRedraw = true
}

func (r *Renderer) ToggleDebugMode(m string) {
	switch m {
	case "changes":
		r.showChanges = !r.showChanges
		log.Info("WEB RENDERER: Toggled cell change display.")
	default:
		log.Error("WEB RENDERER: no debug mode called ", m)
	}
}

func (r *Renderer) Render() {
	if !r.ready || r.console == nil || !r.bufferCtx.Truthy() {
		return
	}
	if !r.console.Dirty() && !r.forceRedraw {
		return
	}

	if r.showChanges {
		r.debugColour = col.MakeOpaque(
			uint8((r.frames*10)%255),
			uint8(((r.frames+100)*10)%255),
			uint8(((r.frames+200)*10)%255),
		)
	}

	bg := make(map[col.Colour][]pxRect)
	fgGlyphs := make(map[col.Colour]map[gfx.Glyph][]vec.Coord)
	fgText := make(map[col.Colour]map[uint8][]pxRect)

	for cell, cursor := range r.console.EachCell() {
		if cell.Mode == gfx.DRAW_NONE || (!r.forceRedraw && !r.console.IsDirtyAt(cursor)) {
			continue
		}

		bgColour := cell.Colours.Back
		if r.showChanges {
			bgColour = r.debugColour
		}
		bg[bgColour] = append(bg[bgColour], pxRect{
			cursor.X * r.tileSize,
			cursor.Y * r.tileSize,
			r.tileSize,
			r.tileSize,
		})

		if !cell.HasForegroundContent() {
			continue
		}
		fg := cell.Colours.Fore
		switch cell.Mode {
		case gfx.DRAW_GLYPH:
			if fgGlyphs[fg] == nil {
				fgGlyphs[fg] = make(map[gfx.Glyph][]vec.Coord)
			}
			fgGlyphs[fg][cell.Glyph] = append(fgGlyphs[fg][cell.Glyph], cursor)
		case gfx.DRAW_TEXT:
			if fgText[fg] == nil {
				fgText[fg] = make(map[uint8][]pxRect)
			}
			for i, char := range cell.Chars {
				fgText[fg][char] = append(fgText[fg][char], pxRect{
					cursor.X*r.tileSize + i*r.tileSize/2,
					cursor.Y * r.tileSize,
					r.tileSize / 2,
					r.tileSize,
				})
			}
		}
	}

	for colour, rects := range bg {
		r.bufferCtx.Set("fillStyle", cssColour(colour))
		for _, rect := range rects {
			r.bufferCtx.Call("clearRect", rect.x, rect.y, rect.w, rect.h)
			if colour.A() != 0 {
				r.bufferCtx.Call("fillRect", rect.x, rect.y, rect.w, rect.h)
			}
		}
	}

	for colour, glyphs := range fgGlyphs {
		atlas := r.tinted(r.glyphs, r.glyphCache, colour)
		for glyph, coords := range glyphs {
			sx := int(glyph%16) * r.tileSize
			sy := int(glyph/16) * r.tileSize
			for _, pos := range coords {
				r.bufferCtx.Call("drawImage", atlas, sx, sy, r.tileSize, r.tileSize, pos.X*r.tileSize, pos.Y*r.tileSize, r.tileSize, r.tileSize)
			}
		}
	}

	half := r.tileSize / 2
	for colour, chars := range fgText {
		atlas := r.tinted(r.font, r.fontCache, colour)
		for char, rects := range chars {
			sx := int(char%32) * half
			sy := int(char/32) * r.tileSize
			for _, rect := range rects {
				r.bufferCtx.Call("drawImage", atlas, sx, sy, half, r.tileSize, rect.x, rect.y, rect.w, rect.h)
			}
		}
	}

	r.console.Clean()
	r.present()
	r.forceRedraw = false
	r.frames++
}

func (r *Renderer) present() {
	w := r.canvas.Get("width").Int()
	h := r.canvas.Get("height").Int()
	r.ctx.Set("fillStyle", cssColour(r.clearColour))
	r.ctx.Call("fillRect", 0, 0, w, h)
	r.ctx.Call("drawImage", r.buffer, 0, 0)
}

func (r *Renderer) tinted(src *image.NRGBA, cache map[col.Colour]js.Value, c col.Colour) js.Value {
	if canvas, ok := cache[c]; ok {
		return canvas
	}
	canvas := imageToCanvas(TintAtlas(src, c))
	cache[c] = canvas
	return canvas
}

func imageToCanvas(img *image.NRGBA) js.Value {
	if img == nil {
		return js.Undefined()
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	canvas := js.Global().Get("document").Call("createElement", "canvas")
	canvas.Set("width", w)
	canvas.Set("height", h)
	ctx := canvas.Call("getContext", "2d")
	data := ctx.Call("createImageData", w, h)
	js.CopyBytesToJS(data.Get("data"), img.Pix)
	ctx.Call("putImageData", data, 0, 0)
	return canvas
}

func cssColour(c col.Colour) string {
	return fmt.Sprintf("rgba(%d,%d,%d,%.3f)", c.R(), c.G(), c.B(), float64(c.A())/255)
}
