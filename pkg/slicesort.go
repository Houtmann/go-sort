package pkg

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var SliceSortAnalyzer = &analysis.Analyzer{
	Name: "slicesorter",
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
			_, ok = node.(*ast.GenDecl)
			var elems []Elem
			switch composite.Type.(type) {
			case *ast.MapType:
				for i := range composite.Elts {
					if kve, ok := composite.Elts[i].(*ast.KeyValueExpr); ok {
						elems = append(elems, Wrap(kve.Key))
					}
				}
			case *ast.ArrayType:
				for i := range composite.Elts {
					elems = append(elems, Wrap(composite.Elts[i]))
				}
			}
			if !isSorted(elems) {
				pass.Reportf(node.Pos(), "fields of are not sorted alphabetically and should be %v ", sortElement(elems))
				return true
			}

			return false
		})
	}

	return nil, nil
}

func Wrap(expr ast.Expr) Elem {
	switch elem := expr.(type) {
	case *ast.Ident:
		return &Ident{*elem}
	case *ast.BasicLit:
		return &BasicLit{*elem}
	default:
		// TODO error
	}
	return nil
}
