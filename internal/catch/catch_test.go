//go:build js && wasm

package catch

import (
	"errors"
	"syscall/js"
	"testing"

	"github.com/hack-pad/safejs/internal/assert"
)

func TestTry(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		result, err := Try(func() string {
			return "foo"
		})
		assert.Equal(t, "foo", result)
		assert.NoError(t, err)
	})

	t.Run("panic string", func(t *testing.T) {
		t.Parallel()
		result, err := Try(func() string {
			panic("some error")
		})
		assert.Equal(t, "", result)
		assert.EqualError(t, err, "some error")
	})

	t.Run("panic error", func(t *testing.T) {
		t.Parallel()
		result, err := Try(func() string {
			panic(errors.New("some error"))
		})
		assert.Equal(t, "", result)
		assert.EqualError(t, err, "some error")
	})

	t.Run("panic value", func(t *testing.T) {
		t.Parallel()
		result, err := Try(func() string {
			panic(js.ValueOf(map[string]any{
				"foo": 1,
			}))
		})
		assert.Equal(t, "", result)
		assert.EqualError(t, err, "JavaScript error: <undefined>")
	})

	t.Run("throws error", func(t *testing.T) {
		t.Parallel()
		result, err := Try(func() string {
			js.Global().Call("Array", -1)
			return "foo"
		})
		assert.Equal(t, "", result)
		assert.EqualError(t, err, "JavaScript error: Invalid array length")
	})
}

func TestTrySideEffect(t *testing.T) {
	t.Parallel()
	err := TrySideEffect(func() {
		t.Log("just a print")
	})
	assert.NoError(t, err)

	err = TrySideEffect(func() {
		panic("some error")
	})
	assert.EqualError(t, err, "some error")
}
