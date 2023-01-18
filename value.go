//go:build js && wasm

package safejs

import (
	"fmt"
	"syscall/js"

	"github.com/hack-pad/safejs/internal/catch"
)

type Value struct {
	jsValue js.Value
}

func Safe(value js.Value) Value {
	return Value{
		jsValue: value,
	}
}

func Unsafe(value Value) js.Value {
	return value.jsValue
}

func Null() Value {
	return Safe(js.Null())
}

func Undefined() Value {
	return Safe(js.Undefined())
}

func toJSValue(jsValue any) any {
	switch value := jsValue.(type) {
	case Value:
		return value.jsValue
	case Func:
		return value.fn
	case Error:
		return value.err
	case map[string]any:
		newValue := make(map[string]any)
		for mapKey, mapValue := range value {
			newValue[mapKey] = toJSValue(mapValue)
		}
		return newValue
	case []any:
		newValue := make([]any, len(value))
		for i, arg := range value {
			newValue[i] = toJSValue(arg)
		}
		return newValue
	default:
		return jsValue
	}
}

func toJSValues(args []any) []any {
	return toJSValue(args).([]any)
}

func toValues(args []js.Value) []Value {
	newArgs := make([]Value, len(args))
	for i, arg := range args {
		newArgs[i] = Safe(arg)
	}
	return newArgs
}

func ValueOf(value any) (Value, error) {
	jsValue, err := catch.Try(func() js.Value {
		return js.ValueOf(value)
	})
	return Safe(jsValue), err
}

func (v Value) Bool() (bool, error) {
	return catch.Try(v.jsValue.Bool)
}

func (v Value) Call(m string, args ...any) (Value, error) {
	args = toJSValues(args)
	return catch.Try(func() Value {
		return Safe(v.jsValue.Call(m, args...))
	})
}

func (v Value) Delete(p string) error {
	return catch.TrySideEffect(func() {
		v.jsValue.Delete(p)
	})
}

func (v Value) Equal(w Value) bool {
	return v.jsValue.Equal(w.jsValue)
}

func (v Value) Float() (float64, error) {
	return catch.Try(v.jsValue.Float)
}

func (v Value) Get(p string) (Value, error) {
	return catch.Try(func() Value {
		return Safe(v.jsValue.Get(p))
	})
}

func (v Value) Index(i int) (Value, error) {
	return catch.Try(func() Value {
		return Safe(v.jsValue.Index(i))
	})
}

func (v Value) InstanceOf(t Value) (bool, error) {
	// Type failures in JS throw "TypeError: Right-hand side of 'instanceof' is not an object"
	// so catch those cases here.
	//
	// A valid type is a function with a field "prototype" which is an object.
	if t.Type() != TypeFunction {
		return false, fmt.Errorf("invalid type for instanceof: %v", t.Type())
	}
	prototype, err := t.Get("prototype")
	if err != nil {
		return false, fmt.Errorf("invalid constructor type for instanceof: %v", err)
	} else if prototype.Type() != TypeObject {
		return false, fmt.Errorf("invalid constructor type for instanceof: %v", prototype.Type())
	}
	return catch.Try(func() bool {
		return v.jsValue.InstanceOf(t.jsValue)
	})
}

func (v Value) Int() (int, error) {
	return catch.Try(v.jsValue.Int)
}

func (v Value) Invoke(args ...any) (Value, error) {
	args = toJSValues(args)
	return catch.Try(func() Value {
		return Safe(v.jsValue.Invoke(args...))
	})
}

func (v Value) IsNaN() bool {
	return v.jsValue.IsNaN()
}

func (v Value) IsNull() bool {
	return v.jsValue.IsNull()
}

func (v Value) IsUndefined() bool {
	return v.jsValue.IsUndefined()
}

func (v Value) Length() (int, error) {
	return catch.Try(v.jsValue.Length)
}

func (v Value) New(args ...any) (Value, error) {
	args = toJSValues(args)
	return catch.Try(func() Value {
		return Safe(v.jsValue.New(args...))
	})
}

func (v Value) Set(p string, x any) error {
	x = toJSValue(x)
	return catch.TrySideEffect(func() {
		v.jsValue.Set(p, x)
	})
}

func (v Value) SetIndex(i int, x any) error {
	x = toJSValue(x)
	return catch.TrySideEffect(func() {
		v.jsValue.SetIndex(i, x)
	})
}

func (v Value) String() (string, error) {
	return catch.Try(v.jsValue.String)
}

func (v Value) Truthy() (bool, error) {
	return catch.Try(v.jsValue.Truthy)
}

func (v Value) Type() Type {
	return Type(v.jsValue.Type())
}
