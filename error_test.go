//go:build js && wasm

package safejs

import (
	"syscall/js"
	"testing"

	"github.com/hack-pad/safejs/internal/assert"
)

func TestError(t *testing.T) {
	t.Parallel()
	t.Run("valid error", func(t *testing.T) {
		t.Parallel()
		jsErr := js.Error{
			Value: js.ValueOf(map[string]any{
				"message": "foo",
			}),
		}
		const expectedErr = "JavaScript error: foo"
		if !assert.EqualError(t, jsErr, expectedErr) { // quickly verify js.Error was created correctly
			return
		}

		err := Error{
			err: jsErr,
		}
		assert.EqualError(t, err, expectedErr)
	})

	t.Run("invalid error", func(t *testing.T) {
		t.Parallel()
		jsErr := js.Error{Value: js.Undefined()}
		err := Error{err: jsErr}
		assert.EqualError(t, err, "failed generating error message: syscall/js: call of Value.Get on undefined")
	})
}
