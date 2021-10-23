package pkg

import (
	"go/ast"
	"go/token"
	"sort"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "gosort",
	Doc:  "Checks if struct fields sorted by alphabetical order",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		structtype, ok := node.(*ast.StructType)
		if !ok {
			return true
		}

		var (
			groupedfields = make(map[token.Pos][]string)
			lastPos token.Pos
			groupPos = structtype.Fields.Pos()
		)


		for i := range structtype.Fields.List{
			if diff := structtype.Fields.List[i].Pos()-lastPos; diff == 2 || lastPos == 0{
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
							return fieldnames[i] < fieldnames[j]})
			if !issorted{
				notsortedfields = append(notsortedfields, fieldnames...)
			}
		}
		if len(notsortedfields) > 0 {
			pass.Reportf(node.Pos(), "fields %v of are not sorted alphabetically", notsortedfields)
			return true
		}
		return false
	}
	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
