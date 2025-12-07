package sdl3

import (
	"errors"
	"unsafe"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
	"github.com/jupiterrider/purego-sdl3/sdl"
)

type Renderer struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	glyphs       *sdl.Texture
	font         *sdl.Texture
	canvasBuffer *sdl.Texture

	tileSize int

	forceRedraw bool
	showFPS     bool
	showChanges bool

	frames int // frames rendered. NOTE: this can differ from engine.tick since the renderer may not render every tick

	clearColour col.Colour
	debugColour col.Colour // for background when show changes is on

	// caches for draw batching
	bg       map[col.Colour][]sdl.FRect
	fgGlyphs map[col.Colour]map[gfx.Glyph][]vec.Coord
	fgText   map[col.Colour]map[uint8][]sdl.FRect

	console *gfx.Canvas

	ready bool
}

// create and get a reference to an SDL Renderer. any sensible defaults can go here too, but the renderer is not
// valid until Setup() has been run on it.
func NewRenderer() *Renderer {
	sdl_renderer := new(Renderer)
	sdl_renderer.ready = false //i know false is already the default value, this is for emphasis.
	return sdl_renderer
}

func (r *Renderer) Setup(console *gfx.Canvas, glyphPath, fontPath, title string) (err error) {
	//renderer defaults to 800x600, once fonts are loaded it figures out the resolution to use and resizes accordingly
	r.window = sdl.CreateWindow(title, 800, 600, sdl.WindowVulkan|sdl.WindowResizable)
	if r.window == nil {
		log.Error("SDL RENDERER: Failed to create window. sdl: ", sdl.GetError())
		return errors.New("failed to create window")
	}

	r.renderer = sdl.CreateRenderer(r.window, "")
	if r.renderer == nil {
		log.Error("SDL RENDERER: Failed to create renderer. sdl: ", sdl.GetError())
		return errors.New("failed to create renderer")
	}

	sdl.SetRenderLogicalPresentation(r.renderer, 800, 600, sdl.LogicalPresentationLetterbox)

	sdl.RenderClear(r.renderer)

	r.console = console
	err = r.ChangeFonts(glyphPath, fontPath)
	if err != nil {
		return err
	}

	r.bg = make(map[col.Colour][]sdl.FRect)
	r.fgGlyphs = make(map[col.Colour]map[gfx.Glyph][]vec.Coord)
	r.fgText = make(map[col.Colour]map[uint8][]sdl.FRect)

	r.ready = true

	return
}

func (r *Renderer) Ready() bool {
	return r.ready
}

// Deletes special graphics structures, closes files, etc. Defer this function!
func (r *Renderer) Cleanup() {
	sdl.DestroyTexture(r.glyphs)
	sdl.DestroyTexture(r.font)
	sdl.DestroyTexture(r.canvasBuffer)
	sdl.DestroyRenderer(r.renderer)
	sdl.DestroyWindow(r.window)
	log.Info("SDL Renderer shut down!")
}

// Loads new fonts to the renderer and changes the tilesize (and by extension, the window size)
func (r *Renderer) ChangeFonts(glyphPath, fontPath string) (err error) {
	if r.glyphs != nil {
		sdl.DestroyTexture(r.glyphs)
	}
	r.glyphs, err = r.loadTexture(glyphPath)
	if err != nil {
		log.Error("SDL RENDERER: Could not load font at ", glyphPath)
		return
	}

	if r.font != nil {
		sdl.DestroyTexture(r.font)
	}
	r.font, err = r.loadTexture(fontPath)
	if err != nil {
		log.Error("SDL RENDERER: Could not load font at ", fontPath)
		return
	}
	log.Info("SDL RENDERER: Loaded fonts! Glyph: " + glyphPath + ", Text: " + fontPath)

	//reset window size if fontsize changed
	var gw, gh float32
	sdl.GetTextureSize(r.glyphs, &gw, &gh)
	if int(gw/16) != r.tileSize {
		r.tileSize = int(gw / 16)
		if r.console == nil {
			log.Error("SDL RENDERER: Console not initialized, cannot determine screen size.")
			err = errors.New("console not intialized")
			return
		}
		console_size := r.console.Size()
		w, h := int32(r.tileSize*console_size.W), int32(r.tileSize*console_size.H)
		sdl.SetWindowSize(r.window, w, h)
		sdl.SetRenderLogicalPresentation(r.renderer, w, h, sdl.LogicalPresentationLetterbox)

		_ = r.createCanvasBuffer() //TODO: handle this error?
		sdl.SetWindowPosition(r.window, sdl.WindowPosCentered, sdl.WindowPosCentered)
		r.forceRedraw = true
		log.Info("SDL RENDERER: resized window.")
	}

	return
}

// Loads a bmp font into the GPU using the current window renderer. White pixels (RGB 255, 255, 255) are modulated with
// a cell's colour, and Fuschia pixels (RGB 255, 0, 255) are transparent.  If the image contains pixels with any other
// G value it converts those pixels to partially transparent white pixels with Alpha = G.
// TODO: support more than bmps?
func (r *Renderer) loadTexture(path string) (*sdl.Texture, error) {
	bmpImage := sdl.LoadBMP(path)
	if bmpImage == nil {
		log.Error("SDL RENDERER: Failed to load image: ", sdl.GetError())
		return nil, errors.New("failed to load image")
	}
	defer sdl.DestroySurface(bmpImage)

	image := sdl.ConvertSurface(bmpImage, sdl.PixelFormatARGB8888)
	defer sdl.DestroySurface(image)

	keyColour := col.FUSCHIA
	transparent := col.NONE

	sdl.LockSurface(image)
	for cursor := range vec.EachCoordInArea(vec.Dims{int(image.W), int(image.H)}) {
		colour := getPixel(image, cursor)
		if colour.G() != 0xFF {
			if colour == keyColour {
				setPixel(image, cursor, transparent)
			} else {
				colour = col.Make(colour.G(), colour.R(), 0xFF, colour.B())
				setPixel(image, cursor, colour)
			}
		}
	}
	sdl.UnlockSurface(image)
	sdl.SetSurfaceRLE(image, true)

	texture := sdl.CreateTextureFromSurface(r.renderer, image)
	if texture == nil {
		log.Error("SDL RENDERER: Failed to create texture: ", sdl.GetError())
		return nil, errors.New("failed to create texture")
	}

	return texture, nil
}

func getPixel(surface *sdl.Surface, cursor vec.Coord) (colour col.Colour) {
	pixels := surface.Pixels
	i := cursor.Y*int(surface.Pitch) + cursor.X*4
	pixel := (*col.Colour)(unsafe.Pointer(uintptr(pixels) + uintptr(i)))
	return *pixel
}

func setPixel(surface *sdl.Surface, cursor vec.Coord, colour col.Colour) {
	pixels := surface.Pixels
	i := cursor.Y*int(surface.Pitch) + cursor.X*4
	pixel := (*col.Colour)(unsafe.Pointer(uintptr(pixels) + uintptr(i)))
	*pixel = colour
	return
}

func (r *Renderer) createCanvasBuffer() (err error) {
	if r.canvasBuffer != nil {
		sdl.DestroyTexture(r.canvasBuffer)
	}

	console_size := r.console.Size()
	r.canvasBuffer = sdl.CreateTexture(r.renderer, sdl.PixelFormatARGB8888, sdl.TextureAccessTarget, int32(console_size.W*r.tileSize), int32(console_size.H*r.tileSize))
	if r.canvasBuffer == nil {
		log.Error("SDL RENDERER: Failed to create buffer texture. sdl:", sdl.GetError())
	}

	return
}

func (r *Renderer) onWindowResize() {
	R, G, B, A := r.clearColour.RGBA()
	sdl.SetRenderDrawColor(r.renderer, R, G, B, A)
	sdl.RenderClear(r.renderer)
	r.forceRedraw = true
}

// Enables or disables fullscreen. All tyumi consoles use borderless fullscreen instead of native
// and the output is scaled to the monitor size.
func (r *Renderer) SetFullscreen(enable bool) {
	if enable {
		sdl.SetWindowFullscreen(r.window, true)
		sdl.SetWindowBordered(r.window, false)
		log.Info("SDL RENDERER: Fullscreen enabled.")
	} else {
		sdl.SetWindowFullscreen(r.window, false)
		sdl.SetWindowBordered(r.window, true)
		log.Info("SDL RENDERER: Fullscreen disabled.")
	}
}

// Toggles between fullscreen modes.
func (r *Renderer) ToggleFullscreen() {
	if sdl.GetWindowFlags(r.window)&sdl.WindowFullscreen != 0 {
		r.SetFullscreen(false)
	} else {
		r.SetFullscreen(true)
	}
}

func (r *Renderer) SetClearColour(colour col.Colour) {
	r.clearColour = colour
	r.forceRedraw = true
}

// Renders the console to the GPU and flips the buffer.
func (r *Renderer) Render() {
	if !r.console.Dirty() && !r.forceRedraw {
		return
	}

	t := sdl.GetRenderTarget(r.renderer)            //store window texture, we'll switch back to it once we're done with the buffer.
	sdl.SetRenderTarget(r.renderer, r.canvasBuffer) //point renderer at buffer texture, we'll draw there

	if r.showChanges {
		r.debugColour = col.MakeOpaque(
			uint8((r.frames*10)%255),
			uint8(((r.frames+100)*10)%255),
			uint8(((r.frames+200)*10)%255),
		)
	}

	//collect rects and coords, sorted by colour
	for cell, cursor := range r.console.EachCell() {
		if cell.Mode == gfx.DRAW_NONE || (!r.forceRedraw && !r.console.IsDirtyAt(cursor)) {
			continue
		}

		bgColour := cell.Colours.Back
		if r.showChanges {
			bgColour = r.debugColour
		}

		if _, ok := r.bg[bgColour]; !ok {
			r.bg[bgColour] = make([]sdl.FRect, 0)
		}

		rect := makeRect(cursor.X*r.tileSize, cursor.Y*r.tileSize, r.tileSize, r.tileSize)
		r.bg[bgColour] = append(r.bg[bgColour], rect)

		if !cell.HasForegroundContent() {
			continue
		}

		fgColour := cell.Colours.Fore

		switch cell.Mode {
		case gfx.DRAW_GLYPH:
			if _, ok := r.fgGlyphs[fgColour]; !ok {
				r.fgGlyphs[fgColour] = make(map[gfx.Glyph][]vec.Coord)
			}

			glyphMap, glyph := r.fgGlyphs[fgColour], cell.Glyph
			if _, ok := glyphMap[glyph]; !ok {
				glyphMap[glyph] = make([]vec.Coord, 0)
			}

			glyphMap[glyph] = append(glyphMap[glyph], cursor)
			r.fgGlyphs[fgColour] = glyphMap
		case gfx.DRAW_TEXT:
			if _, ok := r.fgText[fgColour]; !ok {
				r.fgText[fgColour] = make(map[uint8][]sdl.FRect)
			}

			textMap := r.fgText[fgColour]
			for c_i, char := range cell.Chars {
				if _, ok := textMap[char]; !ok {
					textMap[char] = make([]sdl.FRect, 0)
				}

				dst := makeRect(cursor.X*r.tileSize+c_i*r.tileSize/2, cursor.Y*r.tileSize, r.tileSize/2, r.tileSize)
				textMap[char] = append(textMap[char], dst)
			}

			r.fgText[fgColour] = textMap
		}
	}

	// apply background cell fills
	for colour, rects := range r.bg {
		if len(rects) == 0 {
			delete(r.bg, colour)
			continue
		}
		R, G, B, A := colour.RGBA()
		sdl.SetRenderDrawColor(r.renderer, R, G, B, A)
		sdl.RenderFillRects(r.renderer, rects)
		r.bg[colour] = rects[0:0]
	}

	currentDrawColour := col.NONE

	// copy glyphs
	src := makeRect(0, 0, r.tileSize, r.tileSize)
	for colour, glyphMap := range r.fgGlyphs {
		if len(glyphMap) == 0 {
			delete(r.fgGlyphs, colour)
			continue
		}

		r.setTextureColour(r.glyphs, colour, colour.A() != currentDrawColour.A())
		currentDrawColour = colour
		for glyph, coords := range glyphMap {
			if len(coords) == 0 {
				delete(glyphMap, glyph)
				continue
			}
			src.X, src.Y = float32(int(glyph%16)*r.tileSize), float32(int(glyph/16)*r.tileSize)
			for _, pos := range coords {
				dst := makeRect(pos.X*r.tileSize, pos.Y*r.tileSize, r.tileSize, r.tileSize)
				sdl.RenderTexture(r.renderer, r.glyphs, &src, &dst)
			}
			glyphMap[glyph] = coords[0:0]
		}
		r.fgGlyphs[colour] = glyphMap
	}

	// copy text
	src.W = src.W / 2
	for colour, textMap := range r.fgText {
		if len(textMap) == 0 {
			delete(r.fgText, colour)
			continue
		}

		r.setTextureColour(r.font, colour, colour.A() != currentDrawColour.A())
		currentDrawColour = colour
		for char, rects := range textMap {
			if len(rects) == 0 {
				delete(textMap, char)
				continue
			}

			src.X, src.Y = float32(int(char%32)*r.tileSize/2), float32(int(char/32)*r.tileSize)
			for _, rect := range rects {
				sdl.RenderTexture(r.renderer, r.font, &src, &rect)

			}
			textMap[char] = rects[0:0]
		}
		r.fgText[colour] = textMap
	}

	r.console.Clean()

	sdl.SetRenderTarget(r.renderer, t) //point renderer at window again
	sdl.RenderTexture(r.renderer, r.canvasBuffer, nil, nil)
	sdl.RenderPresent(r.renderer)

	r.forceRedraw = false
	r.frames++
}

func (r *Renderer) setTextureColour(tex *sdl.Texture, colour col.Colour, update_alpha bool) {
	sdl.SetTextureColorMod(tex, colour.R(), colour.G(), colour.B())
	if update_alpha {
		sdl.SetTextureAlphaMod(tex, colour.A())
	}
}

func (r *Renderer) ForceRedraw() {
	r.forceRedraw = true
}

func (r *Renderer) ToggleDebugMode(m string) {
	switch m {
	case "changes":
		r.showChanges = !r.showChanges
		log.Debug("SDL RENDERER: Enabled cell change display debug mode.")
	default:
		log.Error("SDL RENDERER: no debug mode called ", m)
	}
}

func makeRect(x, y, w, h int) sdl.FRect {
	return sdl.FRect{X: float32(x), Y: float32(y), W: float32(w), H: float32(h)}
}
