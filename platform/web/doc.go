// Package web is Tyumi's browser platform. It turns the engine's draw, input, and
// audio calls into canvas, DOM, and Web Audio calls.
//
// Build the game with GOOS=js GOARCH=wasm. A native build and a web build are two
// compiles of the same game; each main picks a platform. The HTML page needs a
// canvas whose id is CanvasID, wasm_exec.js from the Go toolchain, and the wasm
// binary. Serve them over HTTP. Browsers refuse to instantiate wasm from a file URL,
// and the wasm response should use Content-Type application/wasm. example/web is a
// working page.
//
// The page has to load memfs.js before wasm_exec.js. Go's browser port sends
// os.Open, os.WriteFile, and directory listings to that script. Without it, those
// calls fail. Mount fetches a URL and writes it at the same path, so later
// os.Open, reximage, and audio loads succeed. Font and audio loads also fetch a
// path on their own when it is missing. LoadSoundLibrary and any other directory
// listing still needs each file mounted first, because the page cannot list a web
// directory.
//
// Closing the tab ends the page without delivering EV_QUIT. Fullscreen uses the
// browser Fullscreen API and has to follow a key or a click. Audio stays silent
// until the first key or click. SetFramerate(0) never sleeps, so the tab stops
// painting and stops receiving input. Screenshots and UI dumps stay in that
// in-memory filesystem. The program is single-threaded, and the binary is several
// megabytes.
package web

// CanvasID is the id of the canvas element the renderer draws into.
const CanvasID = "tyumi-canvas"
