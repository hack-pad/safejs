//go:build !js

package jsguard

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "jsguard",
	Doc:  "protect against unsafe calls to syscall/js",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if os.Getenv("GOOS") != "js" {
		log.Println("GOOS is not set to 'js', results are likely empty.")
	}

	for _, file := range pass.Files {
		syscallJSImportName, ok := jsImportName(file.Imports)
		if ok {
			ast.Inspect(file, func(node ast.Node) bool {
				return inspectNode(pass, node, syscallJSImportName)
			})
		}
	}
	return nil, nil
}

func isSafeCall(typeName, functionName string) bool {
	name := fmt.Sprintf("%s.%s", typeName, functionName)
	const jsPrefix = "syscall/js."
	if !strings.HasPrefix(name, jsPrefix) {
		return true
	}
	name = strings.TrimPrefix(name, jsPrefix)
	switch name {
	case
		// Safe type+func calls in syscall/js package:
		"Func.Release",
		"Global",
		"Null",
		"Type.String",
		"Undefined",
		"Value.Equal",
		"Value.IsNaN",
		"Value.IsNull",
		"Value.IsUndefined",
		"Value.Type":
		return true
	default:
		return false
	}
}

func jsImportName(imports []*ast.ImportSpec) (string, bool) {
	for _, imprt := range imports {
		if imprt.Path.Value == `"syscall/js"` {
			if imprt.Name != nil {
				return imprt.Name.Name, true
			}
			return "js", true
		}
	}
	return "", false
}

func inspectNode(pass *analysis.Pass, node ast.Node, syscallJSImportName string) bool {
	recurse := false
	recurse = inspectMethodCall(pass, node) || recurse
	recurse = inspectPackageCall(pass, node, syscallJSImportName) || recurse
	return recurse
}

func inspectMethodCall(pass *analysis.Pass, node ast.Node) bool {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return true
	}
	selector, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}
	selectorIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return true
	}
	typeName := pass.TypesInfo.TypeOf(selectorIdent).String()
	if !isSafeCall(typeName, selector.Sel.Name) {
		pass.Reportf(callExpr.Pos(), "unsafe method call on %s found: %s(...)", typeName, formatNode(pass.Fset, selector))
	}
	return true
}

func inspectPackageCall(pass *analysis.Pass, node ast.Node, syscallJSImportName string) bool {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return true
	}
	selector, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}
	selectorIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return true
	}
	if selectorIdent.Name != syscallJSImportName {
		return true
	}
	if !isSafeCall("syscall/js", selector.Sel.Name) {
		pass.Reportf(callExpr.Pos(), "unsafe call to syscall/js found: %s(...)",
			formatNode(pass.Fset, selector))
	}
	return true
}
