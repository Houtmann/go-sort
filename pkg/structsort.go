package pkg

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var StructFieldsSortAnalyzer = &analysis.Analyzer{
	Name: "sort",
	Doc:  "Checks if struct fields sorted by alphabetical order",
	Run:  run,
}

type StructSorter struct {
	fset    *token.FileSet
	structs []*ast.StructType
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
				fieldNames    = make(map[string]string)
			)

		loop:
			for i := range structtype.Fields.List {
				if len(structtype.Fields.List[i].Names) == 0 {
					continue
				}

				for j := range ignoreRules {
					if ignoreRules[j](structtype.Fields.List[i]) {
						continue loop
					}
				}

				if diff := structtype.Fields.List[i].Pos() - lastPos; diff == 2 || lastPos == 0 {
					fieldname := strings.ToLower(structtype.Fields.List[i].Names[0].Obj.Name)
					fieldNames[fieldname] = structtype.Fields.List[i].Names[0].Obj.Name
					groupedfields[groupPos] = append(groupedfields[groupPos], fieldname)
					lastPos = structtype.Fields.List[i].End()
				} else if diff >= 3 {
					fieldname := strings.ToLower(structtype.Fields.List[i].Names[0].Obj.Name)
					fieldNames[fieldname] = structtype.Fields.List[i].Names[0].Obj.Name

					groupPos = structtype.Fields.List[i].Pos()
					groupedfields[groupPos] = append(groupedfields[groupPos], fieldname)
					lastPos = structtype.Fields.List[i].End()
				}

			}

			var notsortedfields [][]string
			for _, fieldnames := range groupedfields {
				issorted := sort.SliceIsSorted(fieldnames, func(i, j int) bool {
					return fieldnames[i] < fieldnames[j]
				})
				if !issorted {
					notsortedfields = append(notsortedfields, fieldnames)
				}
			}
			if len(notsortedfields) > 0 {
				for i := range notsortedfields {
					sort.Strings(notsortedfields[i])
					for j := range notsortedfields[i] {
						fieldname, _ := fieldNames[notsortedfields[i][j]]
						notsortedfields[i][j] = fieldname
					}
				}
				pass.Reportf(node.Pos(), "fields of are not sorted alphabetically and should be %v", notsortedfields)
				return true
			}
			return false
		})
	}

	return nil, nil
}

func ignoreMutexField(node ast.Node) bool {
	field, ok := node.(*ast.Field)
	if !ok {
		return false
	}
	typ, ok := field.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if typ.Sel.Name == "Mutex" {
		return true
	}
	return false
}

type ignoreRule func(node ast.Node) bool

var ignoreRules = []ignoreRule{ignoreMutexField}
