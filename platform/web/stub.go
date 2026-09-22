//go:build !js || !wasm

package web

// Platform is the browser backend. Its tyumi.Platform methods are compiled only
// for GOOS=js GOARCH=wasm, so a desktop build that passes this to SetPlatform
// fails at compile time.
type Platform struct{}

// NewPlatform panics on a desktop build. Build the game with GOOS=js GOARCH=wasm
// to get the browser implementation.
func NewPlatform() *Platform {
	panic("github.com/bennicholls/tyumi/platform/web requires GOOS=js GOARCH=wasm")
}
