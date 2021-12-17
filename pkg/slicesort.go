package pkg

import (
	"fmt"
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
				fmt.Println(composite)
				issorted = sort.SliceIsSorted(composite.Elts, func(i, j int) bool {
					if ident1, ok := composite.Elts[i].(*ast.Ident); ok {
						if ident2, ok := composite.Elts[j].(*ast.Ident); ok {
							return ident1.Name < ident2.Name
						}
					}
					if ident1, ok := composite.Elts[i].(*ast.BasicLit); ok {
						if ident2, ok := composite.Elts[j].(*ast.BasicLit); ok {
							return ident1.Value < ident2.Value
						}
					}
					return false
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
