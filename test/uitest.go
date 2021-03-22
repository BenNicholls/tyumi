package main

import (
	"strconv"

	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/sdlrenderer"
	"github.com/bennicholls/tyumi/gfx/ui"
)

func main() {
	engine.InitConsole(40, 20)
	engine.InitRenderer(new(sdlrenderer.SDLRenderer), "res/glyphs12x24.bmp", "res/font12x24.bmp", "TEST WINDOW")

	state := TestState{}
	state.Setup()
	engine.InitMainState(&state)

	engine.Run()
}

type TestState struct {
	engine.StatePrototype

	text ui.Textbox

	tick int
}

func (ts *TestState) Setup() {
	ts.Init(engine.FIT_CONSOLE, engine.FIT_CONSOLE)
	ts.text = ui.NewTextbox(ui.FIT_TEXT, ui.FIT_TEXT, 1, 1, 0, "TEST STRING DO NOT UPVOTE", true)
	ts.Window().AddElement(&ts.text)
	ts.AddInputHandler(ts.HandleInputs)
}

func (ts *TestState) Update() {
	ts.tick++
	//ts.text.MoveTo(ts.tick%10, 1)
}

func (ts *TestState) UpdateUI() {
	return
}

func (ts *TestState) HandleInputs(e event.Event) {
	switch e.ID() {
	case engine.EV_KEYBOARD:
		ev := e.(engine.KeyboardEvent)
		ts.text.ChangeText(strconv.Itoa(ev.Key))
	}
}
