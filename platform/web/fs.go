//go:build js && wasm

package web

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall/js"
)

var (
	assetMu     sync.Mutex
	remoteCache = map[string][]byte{}
)

// Mount fetches path and writes those bytes into the in-memory filesystem at the
// same path. Later os.Open, directory listings, and Tyumi file loads can use it.
// path is a URL relative to the page, such as "fonts/curses.bmp".
func Mount(path string) error {
	data, err := fetchBytes(path)
	if err != nil {
		return err
	}
	if isRemote(path) {
		assetMu.Lock()
		remoteCache[path] = data
		assetMu.Unlock()
		return nil
	}
	return writeMemFile(path, data)
}

// MountAll mounts each path. Every path is attempted. The returned error joins
// the failures.
func MountAll(paths ...string) error {
	var errs []error
	for _, path := range paths {
		if err := Mount(path); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// ReadAsset returns the bytes of path. A file already written into the in-memory
// filesystem wins. Otherwise the path is fetched and, when it is a relative path,
// stored for later os.Open calls.
func ReadAsset(path string) ([]byte, error) {
	if isRemote(path) {
		assetMu.Lock()
		data, ok := remoteCache[path]
		assetMu.Unlock()
		if ok {
			return data, nil
		}
	} else if data, err := os.ReadFile(path); err == nil {
		return data, nil
	}

	data, err := fetchBytes(path)
	if err != nil {
		return nil, err
	}
	if isRemote(path) {
		assetMu.Lock()
		remoteCache[path] = data
		assetMu.Unlock()
	} else if err := writeMemFile(path, data); err != nil {
		return data, nil
	}
	return data, nil
}

func writeMemFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0o644)
}

func isRemote(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "blob:")
}

func fetchBytes(path string) ([]byte, error) {
	opts := js.Global().Get("Object").New()
	if signal := js.Global().Get("AbortSignal"); signal.Truthy() && signal.Get("timeout").Type() == js.TypeFunction {
		opts.Set("signal", signal.Call("timeout", 30000))
	}

	resp, err := await(js.Global().Call("fetch", path, opts))
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", path, err)
	}
	if !resp.Get("ok").Bool() {
		return nil, fmt.Errorf("fetch %s: status %d", path, resp.Get("status").Int())
	}

	buf, err := await(resp.Call("arrayBuffer"))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	u8 := js.Global().Get("Uint8Array").New(buf)
	data := make([]byte, u8.Get("length").Int())
	js.CopyBytesToGo(data, u8)
	return data, nil
}
