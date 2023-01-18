//go:build !js

package safejs

import (
	"bufio"
	"go/build/constraint"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllWasmTags(t *testing.T) {
	const rootDir = "."
	walkErr := filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		switch {
		case err != nil:
			return err
		case path != rootDir && info.IsDir():
			return filepath.SkipDir // only walk JS packages
		case info.IsDir(),
			filepath.Ext(path) != ".go":
			return nil
		case path == "wasm_tags_test.go":
			// ignore this file, since it must run with file system support enabled
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "//") {
				// hit non-comment line, so no build tags exist (see https://golang.org/cmd/go/#hdr-Build_constraints)
				t.Errorf("File %q does not contain a js,wasm build tag", path)
				break
			}

			expr, err := constraint.Parse(line)
			if err != nil {
				t.Logf("Build constraint failed to parse line in file %q: %q; %v", path, line, err)
				continue
			}
			if isJSWasm(expr) {
				break
			}
		}
		return scanner.Err()
	})
	if walkErr != nil {
		t.Error("Walk failed:", walkErr)
	}
}

func isJSWasm(expr constraint.Expr) bool {
	switch expr := expr.(type) {
	case *constraint.AndExpr:
		x, y := expr.X.String(), expr.Y.String()
		return (x == "js" && y == "wasm") || (x == "wasm" && y == "js")
	default:
		return false
	}
}
