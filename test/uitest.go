package main

import (
	"fmt"

	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/sdlrenderer"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/util"
)

func main() {
	engine.InitConsole(40, 20)
	engine.InitRenderer(sdlrenderer.New(), "res/curses24x24.bmp", "res/font12x24.bmp", "TEST WINDOW")

	state := TestState{}
	state.Setup()
	engine.InitMainState(&state)

	engine.Run()
}

type TestState struct {
	engine.StatePrototype

	text ui.Textbox
	list ui.List

	tick int
}

func (ts *TestState) Setup() {
	ts.Init(engine.FIT_CONSOLE, engine.FIT_CONSOLE)
	ts.Window().SetDefaultColours(col.RED, col.LIME)
	ts.text = ui.NewTextbox(ui.FIT_TEXT, ui.FIT_TEXT, 1, 1, 1, "TEST STRING DO NOT UPVOTE", true)
	ts.text.SetDefaultColours(col.CYAN, col.FUSCHIA)
	ts.text.EnableBorder("TEST TITLE", "TEST HINT")

	text2 := ui.NewTextbox(10, ui.FIT_TEXT, 10, 5, 1, util.LoremIpsum(30), true)
	text2.EnableBorder("lorem", "")
	ts.Window().AddElement(&ts.text, &text2)
	ts.AddInputHandler(ts.HandleInputs)

	// inputbox := ui.NewInputbox(10, 1, 8, 8, 10)
	// inputbox.EnableBorder("inputs!", "do the input")
	// ts.Window().AddElement(&inputbox)

	ts.list = ui.NewList(15, 10, 8, 8, 10)
	for i := 0; i < 20; i++ {
		item := ui.NewTextbox(15, i%3+1, 0, 0, 1, "List item "+fmt.Sprint(i)+"/n", false)
		ts.list.AddElement(&item)
	}

	ts.list.EnableBorder("LIST", "")
	ts.list.ToggleHighlight()
	ts.list.ToggleScrollbar()
	ts.list.SetDefaultColours(col.BLUE, col.WHITE)

	ts.Window().AddElement(&ts.list)
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
	case input.EV_KEYBOARD:
		ev := e.(input.KeyboardEvent)

		if ev.Key == input.K_a {
			item := ui.NewTextbox(15, 1, 0, 0, 1, "new item", false)
			ts.list.AddElement(&item)
		}
	}
}
