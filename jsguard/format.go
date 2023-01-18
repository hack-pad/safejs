package jsguard

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"path/filepath"
)

func ignorePanic() {
	recover()
}

func formatNodeLog(fset *token.FileSet, node ast.Node) string {
	defer ignorePanic()
	position := fset.PositionFor(node.Pos(), true)
	position.Filename = filepath.Base(position.Filename)
	return fmt.Sprint(position, " ", formatNode(fset, node))
}

func formatNode(fset *token.FileSet, x interface{}) string {
	defer ignorePanic()
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, x)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
