//go:build js && wasm

package safejs

import (
	"testing"

	"github.com/hack-pad/safejs/internal/assert"
)

func TestCopyBytesError(t *testing.T) {
	t.Parallel()
	_, err := CopyBytesToGo(nil, Undefined())
	assert.EqualError(t, err, "syscall/js: CopyBytesToGo: expected src to be an Uint8Array or Uint8ClampedArray")

	_, err = CopyBytesToJS(Undefined(), nil)
	assert.EqualError(t, err, "syscall/js: CopyBytesToJS: expected dst to be an Uint8Array or Uint8ClampedArray")
}
