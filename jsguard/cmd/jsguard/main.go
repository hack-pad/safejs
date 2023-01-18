//go:build !js

// Command jsguard reports unsafe calls to Go's syscall/js package
package main

import (
	"github.com/hack-pad/safejs/jsguard"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(jsguard.Analyzer)
}
