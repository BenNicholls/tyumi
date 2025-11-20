//go:build debug

package tyumi

import "github.com/bennicholls/tyumi/log"

const Debug bool = true

func init() {
	log.Info("Tyumi Startup Up in Debug Mode!")
	ShowFPS = true

	// create the debugger. the debugger is centered in the console, so we init it in InitConsole().
	debugger = new(debugDialog)
}
