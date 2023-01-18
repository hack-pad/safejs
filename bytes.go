//go:build js && wasm

package safejs

import (
	"syscall/js"

	"github.com/hack-pad/safejs/internal/catch"
)

func CopyBytesToGo(dst []byte, src Value) (int, error) {
	return catch.Try(func() int {
		return js.CopyBytesToGo(dst, src.jsValue)
	})
}

func CopyBytesToJS(dst Value, src []byte) (int, error) {
	return catch.Try(func() int {
		return js.CopyBytesToJS(dst.jsValue, src)
	})
}
