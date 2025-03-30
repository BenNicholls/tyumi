//go:build debug

package tyumi

import "github.com/bennicholls/tyumi/log"

const debug bool = true

func init() {
	log.SetMinimumLogLevel(log.LVL_DEBUG)
	log.EnableConsoleOutput()
	log.Info("Tyumi Startup Up in Debug Mode!")
}