//go:build !debug

package tyumi

import "github.com/bennicholls/tyumi/log"

const Debug bool = false

func init() {
	log.Info("Tyumi Starting Up!")
}
