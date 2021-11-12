package pkg

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var SliceSortAnalyzer = &analysis.Analyzer{
	Name: "slicesort",
	Doc:  "check function variadic string param order",
	Run:  runslicesort,
}

func runslicesort(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			funccall, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			fun, ok := funccall.Fun.(*ast.Ident)
			if !ok {
				return true
			}
			decl, ok := fun.Obj.Decl.(*ast.FuncDecl)
			if !ok {
				return true
			}
			if len(funccall.Args) > len(decl.Type.Params.List) && len(decl.Type.Params.List) == 1 {
				fmt.Println(fun.Name)
			}

			return false
		})
	}

	return nil, nil
}
