//go:build js && wasm

package safejs

import "syscall/js"

func Global() Value {
	return Safe(js.Global())
}
