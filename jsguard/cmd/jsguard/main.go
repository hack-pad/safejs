//go:build !js

package main

import (
	"github.com/hack-pad/safejs/jsguard"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(jsguard.Analyzer)
}
