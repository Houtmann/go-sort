package pkg

import (
	"go/ast"
	"sort"

	"golang.org/x/tools/go/analysis"
)

var SliceSortAnalyzer = &analysis.Analyzer{
	Name: "slicesort",
	Doc:  "Checks if composite assignment was sorted",
	Run:  runslicesort,
}

func runslicesort(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			composite, ok := node.(*ast.CompositeLit)
			if !ok {
				return true
			}

			issorted := sort.SliceIsSorted(composite.Elts, func(i, j int) bool {
				ident1 := composite.Elts[i].(*ast.BasicLit).Value
				ident2 := composite.Elts[j].(*ast.BasicLit).Value
				return ident1 < ident2
			})
			if !issorted {
				pass.Reportf(node.Pos(), "fields of are not sorted alphabetically")
				return true
			}

			return false
		})
	}

	return nil, nil
}
