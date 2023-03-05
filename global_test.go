//go:build js && wasm

package safejs

import (
	"testing"

	"github.com/hack-pad/safejs/internal/assert"
	"github.com/hack-pad/safejs/internal/catch"
)

func TestGlobal(t *testing.T) {
	t.Parallel()
	global := Global()
	self, err := global.Get("self")
	assert.NoError(t, err)
	assert.Equal(t, true, global.Equal(self))
}

func TestMustGetGlobal(t *testing.T) {
	t.Parallel()
	const className = "Uint8Array"
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		jsUint8Array, err := catch.Try(func() Value {
			return MustGetGlobal(className)
		})
		assert.NoError(t, err)
		jsUint8Array2, err := Global().Get(className)
		assert.NoError(t, err)
		assert.Equal(t, true, jsUint8Array.Equal(jsUint8Array2))
	})

	t.Run("undefined global", func(t *testing.T) {
		t.Parallel()
		_, err := catch.Try(func() Value {
			return MustGetGlobal(className + "-foo")
		})
		assert.EqualError(t, err, `global "Uint8Array-foo" is not defined`)
	})
}
