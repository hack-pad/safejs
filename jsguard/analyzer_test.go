//go:build !js

package jsguard

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hack-pad/safejs/internal/assert"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func init() {
	if err := os.Setenv("GOOS", "js"); err != nil {
		panic(err)
	}
	if err := os.Setenv("GOARCH", "wasm"); err != nil {
		panic(err)
	}
}

func makePackageDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for filePath, contents := range files {
		filePath = filepath.Join(dir, filepath.FromSlash(filePath))
		assert.NoError(t, os.MkdirAll(filepath.Dir(filePath), 0755))
		assert.NoError(t, os.WriteFile(filePath, []byte(contents), 0600))
	}
	return dir
}

func astFile(t *testing.T, pass *analysis.Pass, fileName string) *ast.File {
	t.Helper()
	fileName = strings.TrimSuffix(fileName, ".go")
	for _, file := range pass.Files {
		if file.Name.Name == fileName {
			return file
		}
		t.Log("found file:", file)
	}
	t.Fatal("AST file not found in Pass with name:", fileName)
	return nil
}

func filePos(t *testing.T, pass *analysis.Pass, fileName string, offset int) token.Pos {
	t.Helper()
	astFile := astFile(t, pass, fileName)
	tokenFile := pass.Fset.File(astFile.Pos())
	return tokenFile.Pos(offset)
}

type ignoreTestingErrorf struct{}

func (i ignoreTestingErrorf) Errorf(string, ...interface{}) {}

func TestPackageJSCall(t *testing.T) {
	t.Parallel()
	const (
		fooName = "foo.go"
		fooFile = `
//go:build js && wasm

package foo

import (
	"syscall/js"
)

func Foo() {
	value := js.Null()
	js.CopyBytesToGo(nil, value)
	js.CopyBytesToJS(value, nil)

	js.FuncOf(nil)
	js.Global()
	js.Null()
	js.Undefined()
	js.ValueOf(nil)
}
`
	)
	dir := makePackageDir(t, map[string]string{
		fooName: fooFile,
	})

	result := analysistest.Run(ignoreTestingErrorf{}, dir, Analyzer)
	if !assert.Equal(t, 1, len(result)) {
		t.FailNow()
	}
	result0 := result[0]
	expected := []analysis.Diagnostic{
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "js.CopyBytesToGo(")),
			Message: "unsafe call to syscall/js found: js.CopyBytesToGo(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "js.CopyBytesToJS(")),
			Message: "unsafe call to syscall/js found: js.CopyBytesToJS(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "js.FuncOf(")),
			Message: "unsafe call to syscall/js found: js.FuncOf(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "js.ValueOf(")),
			Message: "unsafe call to syscall/js found: js.ValueOf(...)",
		},
	}
	assert.Equal(t, len(expected), len(result0.Diagnostics))
	if len(result0.Diagnostics) < len(expected) {
		expected = expected[:len(result0.Diagnostics)]
	}
	for i := range expected {
		assert.Equal(t, expected[i], result0.Diagnostics[i])
	}
}

func TestAliasPackageCall(t *testing.T) {
	t.Parallel()
	const (
		fooName = "foo.go"
		fooFile = `
//go:build js && wasm

package foo

import (
	alias "syscall/js"
)

func Foo() {
	value := alias.ValueOf("foo")
	value.String()
}
`
	)
	dir := makePackageDir(t, map[string]string{
		fooName: fooFile,
	})

	result := analysistest.Run(ignoreTestingErrorf{}, dir, Analyzer)
	if !assert.Equal(t, 1, len(result)) {
		t.FailNow()
	}
	result0 := result[0]
	expected := []analysis.Diagnostic{
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "alias.ValueOf(")),
			Message: "unsafe call to syscall/js found: alias.ValueOf(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.String(")),
			Message: "unsafe method call on syscall/js.Value found: value.String(...)",
		},
	}
	assert.Equal(t, len(expected), len(result0.Diagnostics))
	if len(result0.Diagnostics) < len(expected) {
		expected = expected[:len(result0.Diagnostics)]
	}
	for i := range expected {
		assert.Equal(t, expected[i], result0.Diagnostics[i])
	}
}

func TestMethodCall(t *testing.T) {
	t.Parallel()
	const (
		fooName = "foo.go"
		fooFile = `
//go:build js && wasm

package foo

import (
	"syscall/js"
)

func Foo(err js.Error, value js.Value, typ js.Type) {
	err.Error()
	value.Bool()
	value.Call("")
	value.Delete("")
	value.Equal(value)
	value.Float()
	value.Get("")
	value.Index(0)
	value.InstanceOf(value)
	value.Int()
	value.Invoke()
	value.IsNaN()
	value.IsNull()
	value.IsUndefined()
	value.Length()
	value.New()
	value.Set("", nil)
	value.SetIndex(0, nil)
	value.String()
	value.Truthy()
	value.Type()
	typ.String()
}
`
	)
	dir := makePackageDir(t, map[string]string{
		fooName: fooFile,
	})

	result := analysistest.Run(ignoreTestingErrorf{}, dir, Analyzer)
	if !assert.Equal(t, 1, len(result)) {
		t.FailNow()
	}
	result0 := result[0]
	expected := []analysis.Diagnostic{
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "err.Error(")),
			Message: "unsafe method call on syscall/js.Error found: err.Error(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Bool(")),
			Message: "unsafe method call on syscall/js.Value found: value.Bool(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Call(")),
			Message: "unsafe method call on syscall/js.Value found: value.Call(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Delete(")),
			Message: "unsafe method call on syscall/js.Value found: value.Delete(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Float(")),
			Message: "unsafe method call on syscall/js.Value found: value.Float(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Get(")),
			Message: "unsafe method call on syscall/js.Value found: value.Get(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Index(")),
			Message: "unsafe method call on syscall/js.Value found: value.Index(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.InstanceOf(")),
			Message: "unsafe method call on syscall/js.Value found: value.InstanceOf(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Int(")),
			Message: "unsafe method call on syscall/js.Value found: value.Int(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Invoke(")),
			Message: "unsafe method call on syscall/js.Value found: value.Invoke(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Length(")),
			Message: "unsafe method call on syscall/js.Value found: value.Length(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.New(")),
			Message: "unsafe method call on syscall/js.Value found: value.New(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Set(")),
			Message: "unsafe method call on syscall/js.Value found: value.Set(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.SetIndex(")),
			Message: "unsafe method call on syscall/js.Value found: value.SetIndex(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.String(")),
			Message: "unsafe method call on syscall/js.Value found: value.String(...)",
		},
		{
			Pos:     filePos(t, result0.Pass, fooName, strings.Index(fooFile, "value.Truthy(")),
			Message: "unsafe method call on syscall/js.Value found: value.Truthy(...)",
		},
	}
	assert.Equal(t, len(expected), len(result0.Diagnostics))
	if len(result0.Diagnostics) < len(expected) {
		expected = expected[:len(result0.Diagnostics)]
	}
	for i := range expected {
		assert.Equal(t, expected[i], result0.Diagnostics[i])
	}
}
