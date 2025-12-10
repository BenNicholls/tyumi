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
	bgPositions    []sdl.FPoint
	bgColours      []sdl.FColor
	bgIndices      []int32
	glyphPositions []sdl.FPoint
	glyphColours   []sdl.FColor
	glyphUVs       []sdl.FPoint
	glyphIndices   []int32
	textPositions  []sdl.FPoint
	textColours    []sdl.FColor
	textUVs        []sdl.FPoint
	textIndices    []int32

	console *gfx.Canvas

	ready bool
}

// create and get a reference to an SDL Renderer. any sensible defaults can go here too, but the renderer is not
// valid until Setup() has been run on it.
func NewRenderer() *Renderer {
	sdl_renderer := new(Renderer)
	sdl_renderer.ready = false // i know false is already the default value, this is for emphasis.
	return sdl_renderer
}

func (r *Renderer) Setup(console *gfx.Canvas, glyphPath, fontPath, title string) (err error) {
	//renderer defaults to 800x600, once fonts are loaded it figures out the resolution to use and resizes accordingly
	r.window = sdl.CreateWindow(title, 800, 600, sdl.WindowVulkan|sdl.WindowResizable|sdl.WindowHighPixelDensity)
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

	r.bgPositions = make([]sdl.FPoint, 0)
	r.bgColours = make([]sdl.FColor, 0)
	r.bgIndices = make([]int32, 0)
	r.glyphPositions = make([]sdl.FPoint, 0)
	r.glyphUVs = make([]sdl.FPoint, 0)
	r.glyphColours = make([]sdl.FColor, 0)
	r.glyphIndices = make([]int32, 0)
	r.textPositions = make([]sdl.FPoint, 0)
	r.textUVs = make([]sdl.FPoint, 0)
	r.textColours = make([]sdl.FColor, 0)
	r.textIndices = make([]int32, 0)

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

	// process RGB bitmap into alpha-aware surface. pixels with the keycolour are fully transparent,
	// and pixels that are white are fully opaque. grey pixels have some level of transparency, defined
	// by their greyness.
	keyColour := col.FUSCHIA
	transparent := col.NONE

	sdl.LockSurface(image)
	for cursor := range vec.EachCoordInArea(vec.Dims{int(image.W), int(image.H)}) {
		colour := getPixel(image, cursor)
		if colour.G() != 0xFF {
			if colour == keyColour {
				setPixel(image, cursor, transparent)
			} else {
				// gray pixel. use green channel as alpha value. (really it doesn't matter which one you
				// use, they're all the same.)
				colour = col.Make(colour.G(), 0xFF, 0xFF, 0xFF)
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
	i := cursor.Y*int(surface.Pitch) + cursor.X*4 // tyumi colours are 4 bytes wide.
	pixel := (*col.Colour)(unsafe.Pointer(uintptr(surface.Pixels) + uintptr(i)))
	return *pixel
}

func setPixel(surface *sdl.Surface, cursor vec.Coord, colour col.Colour) {
	i := cursor.Y*int(surface.Pitch) + cursor.X*4 // tyumi colours are 4 bytes wide
	pixel := (*col.Colour)(unsafe.Pointer(uintptr(surface.Pixels) + uintptr(i)))
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

	if r.showChanges {
		r.debugColour = col.MakeOpaque(
			uint8((r.frames*10)%255),
			uint8(((r.frames+100)*10)%255),
			uint8(((r.frames+200)*10)%255),
		)
	}

	for cell, cursor := range r.console.EachCell() {
		if cell.Mode == gfx.DRAW_NONE || (!r.forceRedraw && !r.console.IsDirtyAt(cursor)) {
			continue
		}

		cursorPixel := cursor.Scale(r.tileSize) // location of cursor in pixelspace

		bgColour := cell.Colours.Back
		if r.showChanges {
			bgColour = r.debugColour
		}

		addQuadPositions(&r.bgPositions, cursorPixel, r.tileSize, r.tileSize)
		addQuadColours(&r.bgColours, bgColour)
		addQuadIndices(&r.bgIndices)

		if !cell.HasForegroundContent() {
			continue
		}

		fgColour := cell.Colours.Fore

		switch cell.Mode {
		case gfx.DRAW_GLYPH:
			w, h := float32(r.glyphs.W), float32(r.glyphs.H)
			src := vec.Coord{(int(cell.Glyph%16) * r.tileSize), (int(cell.Glyph/16) * r.tileSize)}

			addQuadPositions(&r.glyphPositions, cursorPixel, r.tileSize, r.tileSize)
			addQuadColours(&r.glyphColours, fgColour)
			addQuadIndices(&r.glyphIndices)
			addQuadUVs(&r.glyphUVs, src, w, h, r.tileSize, r.tileSize)
		case gfx.DRAW_TEXT:
			for c_i, char := range cell.Chars {
				if char == 0 || char == 32 {
					continue
				}

				textCursor := vec.Coord{cursorPixel.X + c_i*r.tileSize/2, cursorPixel.Y}
				w, h := float32(r.font.W), float32(r.font.H)
				src := vec.Coord{(int(char%32) * r.tileSize / 2), (int(char/32) * r.tileSize)}

				addQuadPositions(&r.textPositions, textCursor, r.tileSize, r.tileSize/2)
				addQuadColours(&r.textColours, fgColour)
				addQuadIndices(&r.textIndices)
				addQuadUVs(&r.textUVs, src, w, h, r.tileSize, r.tileSize/2)
			}
		}
	}

	t := sdl.GetRenderTarget(r.renderer)            //store window texture, we'll switch back to it once we're done with the buffer.
	sdl.SetRenderTarget(r.renderer, r.canvasBuffer) //point renderer at buffer texture, we'll draw there

	// render background rects
	if len(r.bgPositions) > 0 {
		sdl.RenderGeometryRaw(r.renderer, nil, r.bgPositions, r.bgColours, nil, r.bgIndices)
		r.bgPositions = r.bgPositions[0:0]
		r.bgColours = r.bgColours[0:0]
		r.bgIndices = r.bgIndices[0:0]
	}

	// render glyphs
	if len(r.glyphPositions) > 0 {
		sdl.RenderGeometryRaw(r.renderer, r.glyphs, r.glyphPositions, r.glyphColours, r.glyphUVs, r.glyphIndices)
		r.glyphPositions = r.glyphPositions[0:0]
		r.glyphColours = r.glyphColours[0:0]
		r.glyphIndices = r.glyphIndices[0:0]
		r.glyphUVs = r.glyphUVs[0:0]
	}

	// render text
	if len(r.textPositions) > 0 {
		sdl.RenderGeometryRaw(r.renderer, r.font, r.textPositions, r.textColours, r.textUVs, r.textIndices)
		r.textPositions = r.textPositions[0:0]
		r.textColours = r.textColours[0:0]
		r.textIndices = r.textIndices[0:0]
		r.textUVs = r.textUVs[0:0]
	}

	r.console.Clean()

	sdl.SetRenderTarget(r.renderer, t) //point renderer at window again
	sdl.RenderTexture(r.renderer, r.canvasBuffer, nil, nil)
	sdl.RenderPresent(r.renderer)

	r.forceRedraw = false
	r.frames++
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

// adds the 4 points to the list for draw batching. width is the width of the quad (for half width drawing),
// dst_pixel is the pixel in the top left corner
func addQuadPositions(positions *[]sdl.FPoint, dst_pixel vec.Coord, tileSize, width int) {
	*positions = append(*positions,
		sdl.FPoint{float32(dst_pixel.X), float32(dst_pixel.Y)},
		sdl.FPoint{float32(dst_pixel.X + width), float32(dst_pixel.Y)},
		sdl.FPoint{float32(dst_pixel.X), float32(dst_pixel.Y + tileSize)},
		sdl.FPoint{float32(dst_pixel.X + width), float32(dst_pixel.Y + tileSize)},
	)
}

func addQuadColours(colours *[]sdl.FColor, col col.Colour) {
	sdlColour := sdl.FColor{
		A: float32(col.A()) / 255,
		R: float32(col.R()) / 255,
		G: float32(col.G()) / 255,
		B: float32(col.B()) / 255,
	}
	*colours = append(*colours, sdlColour, sdlColour, sdlColour, sdlColour)
}

func addQuadIndices(indices *[]int32) {
	count := int32(len(*indices) / 6)
	*indices = append(*indices, []int32{
		4*count + 1, 4 * count, 4*count + 2, // triangle 1
		4*count + 1, 4*count + 2, 4*count + 3, // triangle 2
	}...)
}

func addQuadUVs(uvs *[]sdl.FPoint, src_pixel vec.Coord, w, h float32, tileSize int, width int) {
	*uvs = append(*uvs,
		sdl.FPoint{float32(src_pixel.X) / w, float32(src_pixel.Y) / h},
		sdl.FPoint{float32(src_pixel.X+(width)) / w, float32(src_pixel.Y) / h},
		sdl.FPoint{float32(src_pixel.X) / w, float32(src_pixel.Y+tileSize) / h},
		sdl.FPoint{float32(src_pixel.X+(width)) / w, float32(src_pixel.Y+tileSize) / h},
	)
}
