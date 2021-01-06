package engine

import (
	"github.com/bennicholls/tyumi/gfx"
)

//The Console is where the UI of the game state is composited together before being sent to the renderer. It also
//defines the size of the window you're using, and must be initialized with InitConsole() before running the 
//gameloop.
type Console struct {
	gfx.Canvas
}

func InitConsole(w, h int) {
	console = new(Console)
	console.Init(w, h)
}