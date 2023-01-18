//go:build js && wasm

package safejs

import (
	"syscall/js"

	"github.com/hack-pad/safejs/internal/catch"
)

type Func struct {
	fn js.Func
}

func FuncOf(fn func(this Value, args []Value) any) (Func, error) {
	jsFunc, err := toJSFunc(fn)
	return Func{
		fn: jsFunc,
	}, err
}

func toJSFunc(fn func(this Value, args []Value) any) (js.Func, error) {
	jsFunc := func(this js.Value, args []js.Value) any {
		result := fn(Safe(this), toValues(args))
		return toJSValue(result)
	}
	return catch.Try(func() js.Func {
		return js.FuncOf(jsFunc)
	})
}

func (f Func) Release() {
	f.fn.Release()
}

func (f Func) Value() Value {
	return Safe(f.fn.Value)
}
