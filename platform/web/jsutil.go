//go:build js && wasm

package web

import (
	"errors"
	"fmt"
	"sync"
	"syscall/js"
)

func await(promise js.Value) (result js.Value, err error) {
	if !promise.Truthy() || promise.Get("then").Type() != js.TypeFunction {
		return js.Undefined(), errors.New("missing promise")
	}

	done := make(chan struct{})
	var once sync.Once
	finish := func() { once.Do(func() { close(done) }) }
	then := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			result = args[0]
		}
		finish()
		return nil
	})
	defer then.Release()

	catch := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			err = jsError(args[0])
		} else {
			err = errors.New("promise rejected")
		}
		finish()
		return nil
	})
	defer catch.Release()

	promise.Call("then", then, catch)
	<-done
	return result, err
}

func jsError(v js.Value) error {
	if v.Type() == js.TypeString {
		return errors.New(v.String())
	}
	if v.Get("message").Truthy() {
		return errors.New(v.Get("message").String())
	}
	return fmt.Errorf("%v", v.Call("toString").String())
}

func bytesToJS(b []byte) js.Value {
	buf := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(buf, b)
	return buf
}

func configure2D(ctx js.Value) {
	if !ctx.Truthy() {
		return
	}
	ctx.Set("imageSmoothingEnabled", false)
	ctx.Set("webkitImageSmoothingEnabled", false)
}
