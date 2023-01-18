package assert

import (
	"reflect"
	"testing"
)

func NoError(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
		return false
	}
	return true
}

func Equal(t *testing.T, expected, actual any) bool {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected != actual: %#v != %#v", expected, actual)
		return false
	}
	return true
}

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
