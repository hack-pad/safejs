//go:build js && wasm

package safejs

import (
	"syscall/js"
	"testing"

	"github.com/hack-pad/safejs/internal/assert"
)

func TestSafeUnsafe(t *testing.T) {
	jsValue := js.ValueOf("foo")
	assert.Equal(t, jsValue, Unsafe(Safe(jsValue)))
	assert.Equal(t, Value{
		jsValue: jsValue,
	}, Safe(jsValue))
}

func TestValueBool(t *testing.T) {
	const someBool = true
	foo, err := ValueOf(someBool)
	assert.NoError(t, err)

	result, err := foo.Bool()
	assert.NoError(t, err)
	assert.Equal(t, someBool, result)
}

func TestValueFloat(t *testing.T) {
	const someFloat = 42.42
	foo, err := ValueOf(someFloat)
	assert.NoError(t, err)

	result, err := foo.Float()
	assert.NoError(t, err)
	assert.Equal(t, someFloat, result)
}

func TestValueInt(t *testing.T) {
	const someInt = 42
	foo, err := ValueOf(someInt)
	assert.NoError(t, err)

	result, err := foo.Int()
	assert.NoError(t, err)
	assert.Equal(t, someInt, result)
}

func TestValueString(t *testing.T) {
	const someString = "foo"
	foo, err := ValueOf(someString)
	assert.NoError(t, err)

	result, err := foo.String()
	assert.NoError(t, err)
	assert.Equal(t, someString, result)
}

func TestNull(t *testing.T) {
	result, err := ValueOf(nil)
	assert.NoError(t, err)
	assert.Equal(t, result, Null())
}

func TestUndefined(t *testing.T) {
	result, err := ValueOf(js.Undefined())
	assert.NoError(t, err)
	assert.Equal(t, result, Undefined())
}

func TestValueCall(t *testing.T) {
	obj, err := ValueOf(map[string]any{
		"foo": js.FuncOf(func(this js.Value, args []js.Value) any {
			return "bar"
		}),
	})
	assert.NoError(t, err)
	result, err := obj.Call("foo")
	assert.NoError(t, err)
	resultStr, err := result.String()
	assert.NoError(t, err)
	assert.Equal(t, "bar", resultStr)
}

func TestValueDelete(t *testing.T) {
	obj, err := ValueOf(map[string]any{
		"foo": 1,
		"bar": 2,
	})
	assert.NoError(t, err)

	fooValue, err := obj.Get("foo")
	assert.NoError(t, err)
	fooInt, err := fooValue.Int()
	assert.NoError(t, err)
	assert.Equal(t, 1, fooInt)

	assert.NoError(t, obj.Delete("foo"))

	fooValue, err = obj.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, Undefined(), fooValue)
}

func NaN(t *testing.T) Value {
	number, err := Global().Get("Number")
	assert.NoError(t, err)
	valueNaN, err := number.Get("NaN")
	assert.NoError(t, err)
	return valueNaN
}

func TestValueEqual(t *testing.T) {
	value1, err := ValueOf(1)
	assert.NoError(t, err)
	value2, err := ValueOf(2)
	assert.NoError(t, err)

	assert.Equal(t, true, value1.Equal(value1))
	assert.Equal(t, false, value1.Equal(value2))

	valueNaN := NaN(t)
	assert.Equal(t, false, valueNaN.Equal(valueNaN))
}

func TestValueIndex(t *testing.T) {
	arr, err := ValueOf([]any{1, 2, 3})
	assert.NoError(t, err)

	value, err := arr.Index(0)
	assert.NoError(t, err)
	valueInt, err := value.Int()
	assert.NoError(t, err)
	assert.Equal(t, 1, valueInt)
}

func TestValueInstanceOf(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		value, err := ValueOf("foo")
		assert.NoError(t, err)

		_, err = value.InstanceOf(value)
		assert.EqualError(t, err, "invalid type for instanceof: string")
	})

	t.Run("wrong constructor prototype type", func(t *testing.T) {
		fakeClass, err := FuncOf(func(this Value, args []Value) any {
			return nil
		})
		assert.NoError(t, err)
		err = fakeClass.Value().Set("prototype", 1)
		assert.NoError(t, err)
		value, err := ValueOf("foo")
		assert.NoError(t, err)

		_, err = value.InstanceOf(fakeClass.Value())
		assert.EqualError(t, err, "invalid constructor type for instanceof: number")
	})

	t.Run("non-matching type", func(t *testing.T) {
		numberType, err := Global().Get("Number")
		assert.NoError(t, err)
		stringType, err := Global().Get("String")
		assert.NoError(t, err)
		value, err := stringType.New("foo")
		assert.NoError(t, err)

		isInstance, err := value.InstanceOf(numberType)
		assert.NoError(t, err)
		assert.Equal(t, false, isInstance)
	})

	t.Run("matching type", func(t *testing.T) {
		stringType, err := Global().Get("String")
		assert.NoError(t, err)
		value, err := stringType.New("foo")
		assert.NoError(t, err)

		isInstance, err := value.InstanceOf(stringType)
		assert.NoError(t, err)
		assert.Equal(t, true, isInstance)

		isInstance = value.jsValue.InstanceOf(stringType.jsValue)
		assert.Equal(t, true, isInstance)
	})
}

/*
// TODO
func TestValueInvoke(t *testing.T) {
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
	returned, err := fn.Value().Invoke("bar")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, true, returnValue.Equal(returned))
	assert.Equal(t, true, fnThis.Equal(fn.Value()))
	if assert.Equal(t, 1, len(fnArgs)) {
		assert.Equal(t, true, argValue.Equal(fnArgs[0]))
	}
}
*/

func TestValueIsNaN(t *testing.T) {
	value, err := ValueOf("foo")
	assert.NoError(t, err)
	valueNaN := NaN(t)

	assert.Equal(t, false, value.IsNaN())
	assert.Equal(t, true, valueNaN.IsNaN())
}

func TestValueIsNull(t *testing.T) {
	value, err := ValueOf("foo")
	assert.NoError(t, err)

	assert.Equal(t, false, value.IsNull())
	assert.Equal(t, true, Null().IsNull())
}

func TestValueIsUndefined(t *testing.T) {
	value, err := ValueOf("foo")
	assert.NoError(t, err)

	assert.Equal(t, false, value.IsUndefined())
	assert.Equal(t, true, Undefined().IsUndefined())
}

func TestValueLength(t *testing.T) {
	value, err := ValueOf([]any{1, 2, 3})
	assert.NoError(t, err)

	length, err := value.Length()
	assert.NoError(t, err)
	assert.Equal(t, 3, length)
}

func TestValueSetIndex(t *testing.T) {
	value, err := ValueOf([]any{1, 2, 3})
	assert.NoError(t, err)

	const index = 1
	err = value.SetIndex(index, 4)
	assert.NoError(t, err)

	updatedValue, err := value.Index(index)
	assert.NoError(t, err)
	updatedInt, err := updatedValue.Int()
	assert.NoError(t, err)
	assert.Equal(t, 4, updatedInt)
}

func TestValueTruthy(t *testing.T) {
	valueString, err := ValueOf("foo")
	assert.NoError(t, err)
	valueTrue, err := ValueOf(true)
	assert.NoError(t, err)
	valueFalse, err := ValueOf(false)
	assert.NoError(t, err)

	isTruthy, err := valueString.Truthy()
	assert.Equal(t, true, isTruthy && err == nil)
	isTruthy, err = valueTrue.Truthy()
	assert.Equal(t, true, isTruthy && err == nil)
	isTruthy, err = valueFalse.Truthy()
	assert.Equal(t, false, isTruthy && err == nil)
}
