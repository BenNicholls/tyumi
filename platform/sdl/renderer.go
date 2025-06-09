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

	//store render colours so we don't have to set them for every renderer.Copy()
	clearColour         col.Colour
	backDrawColour      col.Colour
	foreDrawColourText  col.Colour
	foreDrawColourGlyph col.Colour

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
		r.backDrawColour = col.MakeOpaque(
			uint8((r.frames*10)%255),
			uint8(((r.frames+100)*10)%255),
			uint8(((r.frames+200)*10)%255))
	}

	for cursor := range vec.EachCoordInArea(r.console) {
		cell := r.console.GetCell(cursor)
		if cell.Dirty || r.forceRedraw {
			if cell.Mode != gfx.DRAW_NONE {
				r.copyToRenderer(cell.Visuals, cursor)
			}
		}
	}

	r.console.Clean()

	r.renderer.SetRenderTarget(t) //point renderer at window again
	r.renderer.Copy(r.canvasBuffer, nil, nil)
	r.renderer.Present()
	r.renderer.SetDrawColor(r.clearColour.RGBA())
	r.renderer.Clear()
	r.renderer.SetDrawColor(r.backDrawColour.RGBA())

	r.forceRedraw = false
	r.frames++
}

// Copies a rect of pixeldata from a source texture to a rect on the renderer's target.
func (r *Renderer) copyToRenderer(vis gfx.Visuals, pos vec.Coord) {
	//change backcolour if it is different from previous draw
	if !r.showChanges && vis.Colours.Back != r.backDrawColour {
		r.backDrawColour = vis.Colours.Back
		r.renderer.SetDrawColor(r.backDrawColour.RGBA())
	}

	dst := makeRect(pos.X*r.tileSize, pos.Y*r.tileSize, r.tileSize, r.tileSize)
	r.renderer.FillRect(&dst)

	//if we're drawing a nothing character (space, whatever), skip next part.
	if !vis.HasForegroundContent() {
		return
	}

	//change texture color mod if it is different from previous draw, then draw glyph/text
	switch vis.Mode {
	case gfx.DRAW_GLYPH:
		if vis.Colours.Fore != r.foreDrawColourGlyph {
			r.setTextureColour(r.glyphs, vis.Colours.Fore, r.foreDrawColourGlyph.A() != vis.Colours.Fore.A())
			r.foreDrawColourGlyph = vis.Colours.Fore
		}
		src := makeRect(int(vis.Glyph%16)*r.tileSize, int(vis.Glyph/16)*r.tileSize, r.tileSize, r.tileSize)
		r.renderer.Copy(r.glyphs, &src, &dst)
	case gfx.DRAW_TEXT:
		if vis.Colours.Fore != r.foreDrawColourText {
			r.setTextureColour(r.font, vis.Colours.Fore, r.foreDrawColourText.A() != vis.Colours.Fore.A())
			r.foreDrawColourText = vis.Colours.Fore
		}
		dst.W = dst.W / 2
		for i, char := range vis.Chars {
			if char == 0 || char == gfx.TEXT_NONE {
				continue
			}
			src := makeRect(int(char%32)*r.tileSize/2, int(char/32)*r.tileSize, r.tileSize/2, r.tileSize)
			dst.X = dst.X + int32(i*(r.tileSize/2))
			r.renderer.Copy(r.font, &src, &dst)
		}
	}
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
