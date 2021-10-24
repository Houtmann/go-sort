package pkg

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var SliceSortAnalyzer = &analysis.Analyzer{
	Name: "slicesort",
	Doc:  "Checks if slice declaration",
	Run:  runslicesort,
}

func runslicesort(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			sliceexpr, ok := node.(*ast.SliceExpr)
			if !ok {
				return true
			}
			fmt.Println(sliceexpr)

			return false
		})
	}

	return nil, nil
}
