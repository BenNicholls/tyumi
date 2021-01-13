package main

import (
	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/gfx/sdlrenderer"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/log"
)

func main() {
	engine.InitConsole(40, 20)
	engine.InitRenderer(new(sdlrenderer.SDLRenderer), "res/glyphs.bmp", "res/font.bmp", "TEST WINDOW")

	state := TestState{}
	state.Init()
	engine.InitMainState(&state)

	engine.Run()
}

type TestState struct {
	engine.StatePrototype

	text ui.Textbox

	tick int
}

func (ts *TestState) Init() {
	ts.InitWindow(engine.FIT_CONSOLE, engine.FIT_CONSOLE)
	ts.text = ui.NewTextbox(ui.FIT_TEXT, ui.FIT_TEXT, 1, 1, 0, "TEST STRING DO NOT UPVOTE", true)
	ts.Window().AddElement(&ts.text)
}

func (ts *TestState) Update() {
	log.Info("TICK ", ts.tick)
	ts.tick++
	ts.text.Pos().MoveTo(ts.tick%10, 1)
	log.Info(ts.text.Pos())
}

func (ts *TestState) UpdateUI() {
	log.Info("TOCK")
}
