package sdl

import (
	"errors"
	"image/color"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
	"github.com/veandco/go-sdl2/sdl"
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
	bg       map[col.Colour][]sdl.Rect
	fgGlyphs map[col.Colour]map[gfx.Glyph][]vec.Coord
	fgText   map[col.Colour]map[uint8][]sdl.Rect

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
	r.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 800, 600, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
	if err != nil {
		log.Error("SDL RENDERER: Failed to create window. sdl: ", sdl.GetError())
		return errors.New("Failed to create window.")
	}

	r.renderer, err = sdl.CreateRenderer(r.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Error("SDL RENDERER: Failed to create renderer. sdl: ", sdl.GetError())
		return errors.New("Failed to create renderer.")
	}

	r.renderer.SetLogicalSize(800, 600)

	r.renderer.Clear()

	r.console = console
	err = r.ChangeFonts(glyphPath, fontPath)
	if err != nil {
		return err
	}

	r.bg = make(map[col.Colour][]sdl.Rect)
	r.fgGlyphs = make(map[col.Colour]map[gfx.Glyph][]vec.Coord)
	r.fgText = make(map[col.Colour]map[uint8][]sdl.Rect)

	r.ready = true

	return
}

func (r *Renderer) Ready() bool {
	return r.ready
}

// Deletes special graphics structures, closes files, etc. Defer this function!
func (r *Renderer) Cleanup() {
	r.glyphs.Destroy()
	r.font.Destroy()
	r.canvasBuffer.Destroy()
	r.renderer.Destroy()
	r.window.Destroy()
	log.Info("SDL Renderer shut down!")
}

// Loads new fonts to the renderer and changes the tilesize (and by extension, the window size)
func (r *Renderer) ChangeFonts(glyphPath, fontPath string) (err error) {
	if r.glyphs != nil {
		r.glyphs.Destroy()
	}
	r.glyphs, err = r.loadTexture(glyphPath)
	if err != nil {
		log.Error("SDL RENDERER: Could not load font at ", glyphPath)
		return
	}

	if r.font != nil {
		r.font.Destroy()
	}
	r.font, err = r.loadTexture(fontPath)
	if err != nil {
		log.Error("SDL RENDERER: Could not load font at ", fontPath)
		return
	}
	log.Info("SDL RENDERER: Loaded fonts! Glyph: " + glyphPath + ", Text: " + fontPath)

	//reset window size if fontsize changed
	_, _, gw, _, _ := r.glyphs.Query()
	if int(gw/16) != r.tileSize {
		r.tileSize = int(gw / 16)
		if r.console == nil {
			log.Error("SDL RENDERER: Console not initialized, cannot determine screen size.")
			err = errors.New("Console not intialized")
			return
		}
		console_size := r.console.Size()
		w, h := int32(r.tileSize*console_size.W), int32(r.tileSize*console_size.H)
		r.window.SetSize(w, h)
		r.renderer.SetLogicalSize(w, h)
		_ = r.createCanvasBuffer() //TODO: handle this error?
		r.window.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)
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
	bmpImage, err := sdl.LoadBMP(path)
	defer bmpImage.Free()

	image, err := bmpImage.ConvertFormat(sdl.PIXELFORMAT_RGBA8888, 0)
	defer image.Free()

	if err != nil {
		log.Error("SDL RENDERER: Failed to load image: ", sdl.GetError())
		return nil, errors.New("Failed to load image")
	}

	keyColour := color.NRGBA{255, 0, 255, 255}
	transparent := color.NRGBA{0, 0, 0, 0}

	image.Lock()
	for cursor := range vec.EachCoordInArea(vec.Dims{int(image.W), int(image.H)}) {
		colour := image.At(cursor.X, cursor.Y).(color.NRGBA)
		if colour.G != 0xFF {
			if colour == keyColour {
				image.Set(cursor.X, cursor.Y, transparent)
			} else {
				colour.A = colour.G
				colour.G = 0xFF
				image.Set(cursor.X, cursor.Y, colour)
			}
		}
	}
	image.Unlock()
	image.SetRLE(true)

	texture, err := r.renderer.CreateTextureFromSurface(image)
	if err != nil {
		log.Error("SDL RENDERER: Failed to create texture: ", sdl.GetError())
		return nil, errors.New("Failed to create texture")
	}

	err = texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		texture.Destroy()
		log.Error("SDL RENDERER: Failed to set blendmode: ", sdl.GetError())
		return nil, errors.New("Failed to set blendmode")
	}

	return texture, nil
}

func (r *Renderer) createCanvasBuffer() (err error) {
	if r.canvasBuffer != nil {
		r.canvasBuffer.Destroy()
	}

	console_size := r.console.Size()
	r.canvasBuffer, err = r.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int32(console_size.W*r.tileSize), int32(console_size.H*r.tileSize))
	if err != nil {
		log.Error("SDL RENDERER: Failed to create buffer texture. sdl:", sdl.GetError())
	}

	return
}

func (r *Renderer) onWindowResize() {
	r.renderer.SetDrawColor(r.clearColour.RGBA())
	r.renderer.Clear()
	r.forceRedraw = true
}

// Enables or disables fullscreen. All tyumi consoles use borderless fullscreen instead of native
// and the output is scaled to the monitor size.
func (r *Renderer) SetFullscreen(enable bool) {
	if enable {
		r.window.SetFullscreen(uint32(sdl.WINDOW_FULLSCREEN_DESKTOP))
		r.window.SetBordered(false)
		log.Info("SDL RENDERER: Fullscreen enabled.")
	} else {
		r.window.SetFullscreen(0)
		r.window.SetBordered(true)
		log.Info("SDL RENDERER: Fullscreen disabled.")
	}
}

// Toggles between fullscreen modes.
func (r *Renderer) ToggleFullscreen() {
	if (r.window.GetFlags() & sdl.WINDOW_FULLSCREEN_DESKTOP) != 0 {
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

	t := r.renderer.GetRenderTarget()          //store window texture, we'll switch back to it once we're done with the buffer.
	r.renderer.SetRenderTarget(r.canvasBuffer) //point renderer at buffer texture, we'll draw there

	if r.showChanges {
		r.debugColour = col.MakeOpaque(
			uint8((r.frames*10)%255),
			uint8(((r.frames+100)*10)%255),
			uint8(((r.frames+200)*10)%255))
	}

	//collect rects and coords, sorted by colour
	for cursor := range r.console.Bounds().EachCoord() {
		cell := r.console.GetCell(cursor)
		if !(cell.Dirty || r.forceRedraw) || cell.Mode == gfx.DRAW_NONE {
			continue
		}

		bgColour := cell.Colours.Back
		if r.showChanges {
			bgColour = r.debugColour
		}

		if _, ok := r.bg[bgColour]; !ok {
			r.bg[bgColour] = make([]sdl.Rect, 0)
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
				r.fgText[fgColour] = make(map[uint8][]sdl.Rect)
			}

			textMap := r.fgText[fgColour]
			for c_i, char := range cell.Chars {
				if _, ok := textMap[char]; !ok {
					textMap[char] = make([]sdl.Rect, 0)
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
			continue
		}
		r.renderer.SetDrawColor(colour.RGBA())
		r.renderer.FillRects(rects)
		r.bg[colour] = rects[0:0]
	}

	currentDrawColour := col.NONE

	// copy glyphs
	src := makeRect(0, 0, r.tileSize, r.tileSize)
	for colour, glyphMap := range r.fgGlyphs {
		r.setTextureColour(r.glyphs, colour, colour.A() != currentDrawColour.A())
		currentDrawColour = colour
		for glyph, coords := range glyphMap {
			src.X, src.Y = int32(int(glyph%16)*r.tileSize), int32(int(glyph/16)*r.tileSize)
			for _, pos := range coords {
				dst := makeRect(pos.X*r.tileSize, pos.Y*r.tileSize, r.tileSize, r.tileSize)
				r.renderer.Copy(r.glyphs, &src, &dst)
			}
			glyphMap[glyph] = coords[0:0]
		}
		r.fgGlyphs[colour] = glyphMap
	}

	// copy text
	src.W = src.W / 2
	for colour, textMap := range r.fgText {
		r.setTextureColour(r.font, colour, colour.A() != currentDrawColour.A())
		currentDrawColour = colour
		for char, rects := range textMap {
			src.X, src.Y = int32(int(char%32)*r.tileSize/2), int32(int(char/32)*r.tileSize)
			for _, rect := range rects {
				r.renderer.Copy(r.font, &src, &rect)
			}
			textMap[char] = rects[0:0]
		}
		r.fgText[colour] = textMap
	}

	r.console.Clean()

	r.renderer.SetRenderTarget(t) //point renderer at window again
	r.renderer.Copy(r.canvasBuffer, nil, nil)
	r.renderer.Present()

	r.forceRedraw = false
	r.frames++
}

func (r *Renderer) setTextureColour(tex *sdl.Texture, colour col.Colour, update_alpha bool) {
	tex.SetColorMod(colour.R(), colour.G(), colour.B())
	if update_alpha {
		tex.SetAlphaMod(colour.A())
	}
}

func (r *Renderer) ForceRedraw() {
	r.forceRedraw = true
}

func (r *Renderer) ToggleDebugMode(m string) {
	switch m {
	case "fps":
		//sdlr.showFPS = !sdlr.showFPS
		log.Warning("SDL RENDERER: FPS display doesn't work, largely due to laziness. Oops.")
	case "changes":
		r.showChanges = !r.showChanges
		log.Debug("SDL RENDERER: Enabled cell change display debug mode.")
	default:
		log.Error("SDL RENDERER: no debug mode called ", m)
	}
}

// int32 for rect arguments. what a world.
func makeRect(x, y, w, h int) sdl.Rect {
	return sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
}
