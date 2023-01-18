//go:build js && wasm

package catch

import (
	"fmt"
	"syscall/js"

	"github.com/hack-pad/safejs/internal/stackerr"
)

func Try[Result any](fn func() Result) (result Result, err error) {
	defer recoverErr(&err)
	result = fn()
	return
}

func TrySideEffect(fn func()) (err error) {
	defer recoverErr(&err)
	fn()
	return
}

func recoverErr(err *error) {
	value := recover()
	valueErr := recoverValueToError(value)
	if valueErr != nil {
		*err = stackerr.WithStack(valueErr)
	}
}

func recoverValueToError(value any) error {
	if value == nil {
		return nil
	}
	switch value := value.(type) {
	case error:
		return value
	case js.Value:
		return js.Error{Value: value}
	default:
		return fmt.Errorf("%+v", value)
	}
}
