package main

import (
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/gfx/sdlrenderer"
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

	tick int
}

func (ts *TestState) Init() {
	
}

func (ts *TestState) Update() {
	log.Info("TICK", ts.tick)
	ts.tick++
}

func (ts *TestState) UpdateUI() {
	log.Info("TOCK")
}

