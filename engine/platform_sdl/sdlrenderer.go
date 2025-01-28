package platform_sdl

import (
	"errors"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
	"github.com/veandco/go-sdl2/sdl"
)

type SDLRenderer struct {
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
	backDrawColour      uint32
	foreDrawColourText  uint32
	foreDrawColourGlyph uint32

	console *gfx.Canvas

	ready bool
}

// create and get a reference to a SDL Renderer. any sensible defaults can go here too, but the renderer is not
// valid until Setup() has been run on it.
func NewRenderer() *SDLRenderer {
	sdlr := new(SDLRenderer)
	sdlr.ready = false //i know false is already the default value, this is for emphasis.
	return sdlr
}

func (sdlr *SDLRenderer) Setup(console *gfx.Canvas, glyphPath, fontPath, title string) (err error) {
	//renderer defaults to 800x600, once fonts are loaded it figures out the resolution to use and resizes accordingly
	sdlr.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 800, 600, sdl.WINDOW_OPENGL)
	if err != nil {
		log.Error("SDL RENDERER: Failed to create window. sdl: ", sdl.GetError())
		return errors.New("Failed to create window.")
	}

	sdlr.renderer, err = sdl.CreateRenderer(sdlr.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Error("SDL RENDERER: Failed to create renderer. sdl: ", sdl.GetError())
		return errors.New("Failed to create renderer.")
	}
	sdlr.renderer.Clear()

	sdlr.console = console
	err = sdlr.ChangeFonts(glyphPath, fontPath)
	if err != nil {
		return err
	}

	sdlr.ready = true

	return
}

func (sdlr *SDLRenderer) Ready() bool {
	return sdlr.ready
}

// Deletes special graphics structures, closes files, etc. Defer this function!
func (sdlr *SDLRenderer) Cleanup() {
	sdlr.glyphs.Destroy()
	sdlr.font.Destroy()
	sdlr.canvasBuffer.Destroy()
	sdlr.renderer.Destroy()
	sdlr.window.Destroy()
	log.Info("SDL Renderer shut down!")
}

// Loads new fonts to the renderer and changes the tilesize (and by extension, the window size)
func (sdlr *SDLRenderer) ChangeFonts(glyphPath, fontPath string) (err error) {
	if sdlr.glyphs != nil {
		sdlr.glyphs.Destroy()
	}
	sdlr.glyphs, err = sdlr.loadTexture(glyphPath)
	if err != nil {
		log.Error("SDL RENDERER: Could not load font at ", glyphPath)
		return
	}
	if sdlr.font != nil {
		sdlr.font.Destroy()
	}
	sdlr.font, err = sdlr.loadTexture(fontPath)
	if err != nil {
		log.Error("SDL RENDERER: Could not load font at ", fontPath)
		return
	}
	log.Info("SDL RENDERER: Loaded fonts! Glyph: " + glyphPath + ", Text: " + fontPath)

	_, _, gw, _, _ := sdlr.glyphs.Query()

	//reset window size if fontsize changed
	if int(gw/16) != sdlr.tileSize {
		sdlr.tileSize = int(gw / 16)
		if sdlr.console == nil {
			log.Error("SDL RENDERER: Console not initialized, cannot determine screen size.")
			err = errors.New("Console not intialized")
			return
		}
		console_size := sdlr.console.Size()
		sdlr.window.SetSize(int32(sdlr.tileSize*console_size.W), int32(sdlr.tileSize*console_size.H))
		_ = sdlr.createCanvasBuffer() //TODO: handle this error?
		sdlr.window.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)
		log.Info("SDL RENDERER: resized window.")
	}

	return
}

// Loads a bmp font into the GPU using the current window renderer.
// TODO: support more than bmps?
func (sdlr *SDLRenderer) loadTexture(path string) (*sdl.Texture, error) {
	image, err := sdl.LoadBMP(path)
	defer image.Free()
	if err != nil {
		log.Error("SDL RENDERER: Failed to load image: ", sdl.GetError())
		return nil, errors.New("Failed to load image")
	}
	image.SetColorKey(true, col.KEY)
	texture, err := sdlr.renderer.CreateTextureFromSurface(image)
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

func (sdlr *SDLRenderer) createCanvasBuffer() (err error) {
	if sdlr.canvasBuffer != nil {
		sdlr.canvasBuffer.Destroy()
	}
	console_size := sdlr.console.Size()
	sdlr.canvasBuffer, err = sdlr.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int32(console_size.W*sdlr.tileSize), int32(console_size.H*sdlr.tileSize))
	if err != nil {
		log.Error("SDL RENDERER: Failed to create buffer texture. sdl:", sdl.GetError())
	}
	return
}

// Enables or disables fullscreen. All tyumi consoles use borderless fullscreen instead of native
// and the output is scaled to the monitor size.
func (sdlr *SDLRenderer) SetFullscreen(enable bool) {
	if enable {
		sdlr.window.SetFullscreen(uint32(sdl.WINDOW_FULLSCREEN_DESKTOP))
		sdlr.window.SetBordered(false)
		log.Info("SDL RENDERER: Fullscreen enabled.")
	} else {
		sdlr.window.SetFullscreen(0)
		sdlr.window.SetBordered(true)
		log.Info("SDL RENDERER: Fullscreen disabled.")
	}
}

// Toggles between fullscreen modes.
func (sdlr *SDLRenderer) ToggleFullscreen() {
	if (sdlr.window.GetFlags() & sdl.WINDOW_FULLSCREEN_DESKTOP) != 0 {
		sdlr.SetFullscreen(false)
	} else {
		sdlr.SetFullscreen(true)
	}
}

// Renders the console to the GPU and flips the buffer.
func (sdlr *SDLRenderer) Render() {
	var src, dst sdl.Rect
	t := sdlr.renderer.GetRenderTarget()             //store window texture, we'll switch back to it once we're done with the buffer.
	sdlr.renderer.SetRenderTarget(sdlr.canvasBuffer) //point renderer at buffer texture, we'll draw there

	for cursor := range vec.EachCoord(sdlr.console) {
		cell := sdlr.console.GetCell(cursor)
		if cell.Dirty || sdlr.forceRedraw {
			if cell.Mode == gfx.DRAW_TEXT {
				for c_i, char := range cell.Chars {
					dst = makeRect(cursor.X*sdlr.tileSize+c_i*sdlr.tileSize/2, cursor.Y*sdlr.tileSize, sdlr.tileSize/2, sdlr.tileSize)
					src = makeRect((int(char)%32)*sdlr.tileSize/2, (int(char)/32)*sdlr.tileSize, sdlr.tileSize/2, sdlr.tileSize)
					sdlr.copyToRenderer(gfx.DRAW_TEXT, src, dst, cell.Colours, int(char))
				}
			} else {
				g := cell.Glyph
				dst = makeRect(cursor.X*sdlr.tileSize, cursor.Y*sdlr.tileSize, sdlr.tileSize, sdlr.tileSize)
				src = makeRect((g%16)*sdlr.tileSize, (g/16)*sdlr.tileSize, sdlr.tileSize, sdlr.tileSize)
				sdlr.copyToRenderer(gfx.DRAW_GLYPH, src, dst, cell.Colours, g)
			}
		}
	}

	sdlr.console.Clean()

	sdlr.renderer.SetRenderTarget(t) //point renderer at window again
	sdlr.renderer.Copy(sdlr.canvasBuffer, nil, nil)
	sdlr.renderer.Present()
	sdlr.renderer.Clear()
	sdlr.forceRedraw = false

	sdlr.frames++
}

// Copies a rect of pixeldata from a source texture to a rect on the renderer's target.
func (sdlr *SDLRenderer) copyToRenderer(mode gfx.DrawMode, src, dst sdl.Rect, colours col.Pair, g int) {
	//change backcolour if it is different from previous draw
	if colours.Back != sdlr.backDrawColour {
		sdlr.backDrawColour = colours.Back
		sdlr.renderer.SetDrawColor(col.RGBA(sdlr.backDrawColour))
	}

	if sdlr.showChanges {
		sdlr.renderer.SetDrawColor(uint8((sdlr.frames*10)%255), uint8(((sdlr.frames+100)*10)%255), uint8(((sdlr.frames+200)*10)%255), 0xFF) //Test Function
	}

	sdlr.renderer.FillRect(&dst)

	//if we're drawing a nothing character (space, whatever), skip next part.
	if mode == gfx.DRAW_GLYPH && (g == gfx.GLYPH_NONE || g == gfx.GLYPH_SPACE) {
		return
	} else if mode == gfx.DRAW_TEXT && g == 32 {
		return
	} else if colours.Fore == sdlr.backDrawColour { //skip drawing foreground if it is the same as the background
		return
	}

	//change texture color mod if it is different from previous draw, then draw glyph/text
	if mode == gfx.DRAW_GLYPH {
		if colours.Fore != sdlr.foreDrawColourGlyph {
			sdlr.foreDrawColourGlyph = colours.Fore
			sdlr.setTextureColour(sdlr.glyphs, sdlr.foreDrawColourGlyph)
		}
		sdlr.renderer.Copy(sdlr.glyphs, &src, &dst)
	} else {
		if colours.Fore != sdlr.foreDrawColourText {
			sdlr.foreDrawColourText = colours.Fore
			sdlr.setTextureColour(sdlr.font, sdlr.foreDrawColourText)
		}
		sdlr.renderer.Copy(sdlr.font, &src, &dst)
	}
}

func (sdlr *SDLRenderer) setTextureColour(tex *sdl.Texture, colour uint32) {
	r, g, b, a := col.RGBA(colour)
	tex.SetColorMod(r, g, b)
	tex.SetAlphaMod(a)
}

func (sdlr *SDLRenderer) ForceRedraw() {
	sdlr.forceRedraw = true
}

func (sdlr *SDLRenderer) ToggleDebugMode(m string) {
	switch m {
	case "fps":
		//sdlr.showFPS = !sdlr.showFPS
		log.Warning("SDL RENDERER: FPS display doesn't work, largely due to laziness. Oops.")
	case "changes":
		sdlr.showChanges = !sdlr.showChanges
		log.Debug("SDL RENDERER: Enabled cell change display debug mode.")
	default:
		log.Error("SDL RENDERER: no debug mode called ", m)
	}
}

// int32 for rect arguments. what a world.
func makeRect(x, y, w, h int) sdl.Rect {
	return sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
}
