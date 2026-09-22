//go:build js && wasm

package main

import (
	"strings"
	"time"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/platform/web"
	"github.com/bennicholls/tyumi/vec"
)

const (
	glyphFile = "assets/glyph.bmp"
	fontFile  = "assets/font.bmp"
)

type demo struct {
	tyumi.Scene
	status *ui.Textbox
}

func main() {
	log.EnableConsoleOutput()
	if err := writeFonts(); err != nil {
		log.Error("Could not write fonts: ", err)
		return
	}
	if err := tyumi.SetPlatform(web.NewPlatform()); err != nil {
		log.Error(err)
		return
	}

	tyumi.InitConsole("Tyumi Web", vec.Dims{48, 16}, glyphFile, fontFile)
	tyumi.SetClearColour(col.BLACK)

	scene := new(demo)
	scene.Init()

	title := ui.NewTextbox(vec.Dims{44, 1}, vec.Coord{2, 2}, 1, "TYUMI ON THE WEB", ui.ALIGN_LEFT)
	title.SetDefaultColours(col.Pair{col.WHITE, col.NAVY})
	hint := ui.NewTextbox(vec.Dims{44, 1}, vec.Coord{2, 4}, 1, "PRESS A KEY", ui.ALIGN_LEFT)
	hint.SetDefaultColours(col.Pair{col.LIME, col.BLACK})
	scene.status = ui.NewTextbox(vec.Dims{44, 1}, vec.Coord{2, 6}, 1, "LAST KEY: NONE", ui.ALIGN_LEFT)
	scene.status.SetDefaultColours(col.Pair{col.YELLOW, col.BLACK})
	scene.Window().AddChildren(title, hint, scene.status)
	scene.SetKeypressHandler(scene.onKey)

	tyumi.EnableCursor()
	tyumi.SetInitialScene(scene)
	tyumi.Run()
}

func (d *demo) Update(time.Duration) {
	d.Window().DrawVisuals(vec.Coord{2, 9}, 1, gfx.NewGlyphVisuals(gfx.GLYPH_BLOCK, col.Pair{col.ORANGE, col.BLACK}))
}

func (d *demo) onKey(ev *input.KeyboardEvent) bool {
	label := strings.ToUpper(ev.Text())
	if label == "" {
		switch ev.Key {
		case input.K_LEFT:
			label = "LEFT"
		case input.K_RIGHT:
			label = "RIGHT"
		case input.K_UP:
			label = "UP"
		case input.K_DOWN:
			label = "DOWN"
		case input.K_ESCAPE:
			label = "ESCAPE"
		case input.K_RETURN:
			label = "ENTER"
		default:
			label = "KEY"
		}
	}
	d.status.ChangeText("LAST KEY: " + label)
	return true
}
