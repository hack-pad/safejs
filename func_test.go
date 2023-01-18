//go:build js && wasm

package safejs

import (
	"testing"

	"github.com/hack-pad/safejs/internal/assert"
)

func TestFuncOf(t *testing.T) {
	t.Parallel()
	var fnThis Value
	var fnArgs []Value
	returnValue, err := ValueOf("foo")
	assert.NoError(t, err)
	argValue, err := ValueOf("bar")
	assert.NoError(t, err)

	fn, err := FuncOf(func(this Value, args []Value) any {
		fnThis = this
		fnArgs = args
		return returnValue
	})
	assert.NoError(t, err)
	returned, err := fn.Value().Invoke(argValue)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, true, returnValue.Equal(returned))
	assert.Equal(t, true, fnThis.IsUndefined())
	if assert.Equal(t, 1, len(fnArgs)) {
		assert.Equal(t, true, argValue.Equal(fnArgs[0]))
	}

	assert.NotPanics(t, func() {
		fn.Release()
	})
}
