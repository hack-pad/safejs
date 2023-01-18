// Package assert contains small assertion test functions to assist in writing clean tests.
package assert

import (
	"reflect"
	"testing"
)

// NoError asserts err is nil
func NoError(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
		return false
	}
	return true
}

// Equal asserts actual is equal to expected
func Equal(t *testing.T, expected, actual any) bool {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected != actual: %#v != %#v", expected, actual)
		return false
	}
	return true
}

// EqualError asserts err.Error() is equal to expected
func EqualError(t *testing.T, err error, expected string) bool {
	t.Helper()
	if err == nil {
		t.Error("Expected error, got nil")
		return false
	}
	var actual string
	NotPanics(t, func() {
		actual = err.Error()
	})
	if expected != actual {
		t.Errorf("Expected != actual: Type=%T  %#v != %#v", err, expected, actual)
		return false
	}
	return true
}

// NotPanics asserts fn() does not panic
func NotPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		val := recover()
		if val != nil {
			t.Errorf("Unexpected panic. Value: %#v", val)
		}
	}()
	fn()
}
