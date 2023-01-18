//go:build js && wasm

package safejs

import (
	"syscall/js"

	"github.com/hack-pad/safejs/internal/catch"
)

type Error struct {
	err js.Error
}

func (e Error) Error() string {
	errStr, err := catch.Try(e.err.Error)
	if err != nil {
		return "failed generating error message: " + err.Error()
	}
	return errStr
}
