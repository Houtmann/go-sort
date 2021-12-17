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
			var (
				issorted bool
			)
			switch composite.Type.(type) {
			case *ast.MapType:
				issorted = sort.SliceIsSorted(composite.Elts, func(i, j int) bool {
					key1 := composite.Elts[i].(*ast.KeyValueExpr).Key
					key2 := composite.Elts[j].(*ast.KeyValueExpr).Key
					if ident1, ok := key1.(*ast.Ident); ok {
						if ident2, ok := key2.(*ast.Ident); ok {
							return ident1.Name < ident2.Name
						}
					}
					if ident1, ok := key1.(*ast.BasicLit); ok {
						if ident2, ok := key2.(*ast.BasicLit); ok {
							return ident1.Value < ident2.Value
						}
					}
					return false
				})
			case *ast.ArrayType:
				issorted = sort.SliceIsSorted(composite.Elts, func(i, j int) bool {
					ident1 := composite.Elts[i].(*ast.BasicLit).Value
					ident2 := composite.Elts[j].(*ast.BasicLit).Value
					return ident1 < ident2
				})
			}

			if !issorted {
				pass.Reportf(node.Pos(), "fields of are not sorted alphabetically")
				return true
			}

			return false
		})
	}

	return nil, nil
}
