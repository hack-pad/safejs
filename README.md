# SafeJS
A safer, drop-in replacement for Go's `syscall/js` JavaScript package.

## What makes it safer?

Today, `syscall/js` panics when the JavaScript runtime throws errors.
While this is sensible behavior in a JavaScript runtime, it is in stark contrast with Go's pattern of returned errors.

SafeJS provides a near-identical API to `syscall/js`, but returns errors instead of panicking.

Although returned errors aren't pretty, they make it much easier to integrate with existing Go tools and code patterns.

**Please note:** This package uses the same backward compatibility guarantees as `syscall/js`. In an effort to align with the Go standard library's API, some breaking changes may become necessary and will receive their own minor version bumps.

## Quick start

1. Import `safejs`:
```go
import "github.com/hack-pad/safejs"
```
2. Replace typical uses of `syscall/js` with the `safejs` alternative. 

Before:
```go
//go:build js && wasm

package buttons

import "syscall/js"

// InsertButton creates a new button, adds it to 'container', and returns it. Usually.
func InsertButton(container js.Value) js.Value {
    // *whisper:* There's a good chance it could panic! Eh, probably don't need to document it, right?
    dom, err := js.Global().Get("document")
    if err != nil {
        return err
    }
    button, err := dom.Call("createElement", "button")
    if err != nil {
        return err
    }
    _, err = container.Call("appendChild", button)
    if err != nil {
        return err
    }
    return button, nil
}
```

After:
```go
//go:build js && wasm

package buttons

import "github.com/hack-pad/safejs"

// InsertButton creates a new button, adds it to 'container', and returns it or an error.
func InsertButton(container safejs.Value) (safejs.Value, error) {
    dom, err := safejs.Global().Get("document")
    if err != nil {
        return err
    }
    button, err := dom.Call("createElement", "button")
    if err != nil {
        return err
    }
    _, err = container.Call("appendChild", button)
    if err != nil {
        return err
    }
    return button, nil
}
```

## Even safer

If you would like additional safety when working with JavaScript, then use the `jsguard` linter as well.

`jsguard` reports the locations of unsafe JavaScript calls, which should be replaced with SafeJS.

```bash
go install github.com/hack-pad/safejs/jsguard/cmd/jsguard
export GOOS=js GOARCH=wasm
jsguard ./...
```
