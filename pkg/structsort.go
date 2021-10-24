package pkg

import (
	"go/ast"
	"go/token"
	"sort"

	"golang.org/x/tools/go/analysis"
)

var StructFieldsSortAnalyzer = &analysis.Analyzer{
	Name: "sort",
	Doc:  "Checks if struct fields sorted by alphabetical order",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			typeDecl, ok := node.(*ast.GenDecl)
			if !ok {
				return true
			}
			structdecl, ok := typeDecl.Specs[0].(*ast.TypeSpec)
			if !ok {
				return true
			}
			structtype, ok := structdecl.Type.(*ast.StructType)
			if !ok {
				return true
			}

			var (
				groupedfields = make(map[token.Pos][]string)
				lastPos       token.Pos
				groupPos      = structtype.Fields.Pos()
			)

			for i := range structtype.Fields.List {
				if len(structtype.Fields.List[i].Names) == 0 {
					continue
				}
				if diff := structtype.Fields.List[i].Pos() - lastPos; diff == 2 || lastPos == 0 {
					groupedfields[groupPos] = append(groupedfields[groupPos], structtype.Fields.List[i].Names[0].Obj.Name)
				} else if diff == 3 {
					groupPos = structtype.Fields.List[i].Pos()
					groupedfields[groupPos] = append(groupedfields[groupPos], structtype.Fields.List[i].Names[0].Obj.Name)
				}
				lastPos = structtype.Fields.List[i].End()
			}

			var notsortedfields []string
			for _, fieldnames := range groupedfields {
				issorted := sort.SliceIsSorted(fieldnames, func(i, j int) bool {
					return fieldnames[i] < fieldnames[j]
				})
				if !issorted {
					notsortedfields = append(notsortedfields, fieldnames...)
				}
			}
			if len(notsortedfields) > 0 {
				pass.Reportf(node.Pos(), "%s fields of are not sorted alphabetically", structdecl.Name)
				return true
			}
			return false
		})
	}

	return nil, nil
}
