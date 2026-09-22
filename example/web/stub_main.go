//go:build !js || !wasm

package main

import "fmt"

func main() {
	fmt.Println("Build this example with GOOS=js GOARCH=wasm. See index.html.")
}
