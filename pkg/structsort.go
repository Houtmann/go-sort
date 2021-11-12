package pkg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"reflect"
	"sort"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ast/inspector"
)

var StructDeclFieldsSortAnalyzer = &analysis.Analyzer{
	Name:     "struct",
	Doc:      "Checks if struct fields declaration sorted by alphabetical order",
	Run:      structSorter,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type Field struct {
	Name string
	Type reflect.Type
}

func structSorter(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		var s *ast.StructType
		var ok bool
		if s, ok = node.(*ast.StructType); !ok {
			return
		}
		if tv, ok := pass.TypesInfo.Types[s]; ok {
			_, indexes := optimalOrder(tv.Type.(*types.Struct), s.Fields.List)

			var flat []*ast.Field
			for _, f := range s.Fields.List {
				f.Comment = nil
				f.Doc = nil
				if len(f.Names) <= 1 {
					flat = append(flat, f)
					continue
				}
				for _, name := range f.Names {
					flat = append(flat, &ast.Field{
						Names: []*ast.Ident{name},
						Type:  f.Type,
					})
				}
			}

			var reordered []*ast.Field
			for _, index := range indexes {
				reordered = append(reordered, flat[index])
			}

			newStr := &ast.StructType{
				Fields: &ast.FieldList{
					List: reordered,
				},
			}
			// Write the newly aligned struct node to get the content for suggested fixes.
			var buf bytes.Buffer
			if err := format.Node(&buf, token.NewFileSet(), newStr); err != nil {
				return
			}

			pass.Report(analysis.Diagnostic{
				Pos:     s.Pos(),
				End:     s.Pos() + token.Pos(len("struct")),
				Message: "toto",
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: "Rearrange fields",
					TextEdits: []analysis.TextEdit{{
						Pos:     s.Pos(),
						End:     s.End(),
						NewText: buf.Bytes(),
					}},
				}},
			})
		}

	})
	return nil, nil
}

func optimalOrder(str *types.Struct, unorderedFields []*ast.Field) (*types.Struct, []int) {
	var t []
	sort.Slice(unorderedFields, func (i, j int) bool {
	return unorderedFields[i].Names[0].Name < unorderedFields[j].Names[0].Name
	})

	fields := make([]*types.Var, len(unorderedFields))
	indexes := make([]int, len(unorderedFields))
	for i := range unorderedFields {
	fields[i] = str.Field(i)
	indexes[i] = i
	}
	return types.NewStruct(fields, nil), indexes
}

func structSorter2(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			var (
				groupedfields map[token.Pos][]string
				nodeName      *ast.Ident
			)
			switch nodetype := node.(type) {
			// Struct declaration
			case *ast.TypeSpec:
				structtype, ok := nodetype.Type.(*ast.StructType)
				if !ok {
					return true
				}

				nodeName = nodetype.Name
				_ = structDeclSorter(structtype)

			case *ast.AssignStmt:
				if composite, ok := nodetype.Rhs[0].(*ast.CompositeLit); ok {
					nodeName = nodetype.Lhs[0].(*ast.Ident)
					groupedfields = structAssignStmtSorter(composite)
				}
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
			if nodeName == nil {
				return true
			}

			// newnode := render(pass.Fset, node)
			// pass.Report(analysis.Diagnostic{
			// 	Pos:     node.Pos(),
			// 	Message: "toto",
			// 	SuggestedFixes: []analysis.SuggestedFix{
			// 		{
			// 			Message: fmt.Sprintf("should replace with `%s`", newnode),
			// 			TextEdits: []analysis.TextEdit{
			// 				{
			// 					Pos:     node.Pos(),
			// 					End:     node.End(),
			// 					NewText: []byte(newnode),
			// 				},
			// 			},
			// 		},
			// 	},
			// })

			// printer.Fprint(os.Stdout, pass.Fset, file)
			if len(notsortedfields) > 0 {
				pass.Reportf(node.Pos(), "%s fields of are not sorted alphabetically", nodeName)
				return false
			}

			return true
		})

	}
	return nil, nil
}

func structDeclSorter(node *ast.StructType) bool {
	// var (
	// 	groupedfields = make(map[token.Pos][]*ast.Field)
	// 	lastPos       token.Pos
	// 	groupPos      = node.Fields.List[0].Pos()
	// )
	//
	// for i := range node.Fields.List {
	// 	field := node.Fields.List[i]
	// 	if diff := field.Pos() - lastPos; diff == 2 || lastPos < 0 {
	// 		groupedfields[groupPos] = append(groupedfields[groupPos], field)
	// 	} else if diff > 2 {
	// 		groupPos = field.End()
	// 		groupedfields[groupPos] = append(groupedfields[groupPos], field)
	// 	}
	// 	lastPos = field.End()
	// }
	for _ = range node.Fields.List {
		astutil.Apply(node, nil, func(cursor *astutil.Cursor) bool {
			fmt.Println(cursor)

			return true
		})
	}
	// sort.SliceStable(node.Fields.List, func(i, j int) bool {
	// 	// fmt.Println(node.Fields.List[j].Names[0].Name, node.Fields.List[j].Pos(), node.Fields.List[j].End(), node.Fields.List[i].Names[0].Name,
	// 	// 	node.Fields.List[i].Pos(), node.Fields.List[i].End())
	// 	// if node.Fields.List[i].Pos()-node.Fields.List[j].End() <= 2 {
	// 		if node.Fields.List[i].Names[0].Name < node.Fields.List[j].Names[0].Name {
	// 			cursor.InsertAfter(node.Fields.List[j])
	// 			//cursor.InsertBefore(node.Fields.List[i])
	// 			return true
	// 		}
	// 	// }
	// 	return false
	// })

	return true
}

func structAssignStmtSorter(node *ast.CompositeLit) map[token.Pos][]string {
	var (
		groupedfields = make(map[token.Pos][]string)
		lastPos       token.Pos
		groupPos      = node.Elts[0].Pos()
	)

	for i := range node.Elts {
		keyvalueExpr := node.Elts[i].(*ast.KeyValueExpr)
		ident := keyvalueExpr.Key.(*ast.Ident)

		if diff := keyvalueExpr.Pos() - lastPos; diff == 4 || lastPos < 0 {
			groupedfields[groupPos] = append(groupedfields[groupPos], ident.Name)
		} else if diff > 4 {
			groupPos = keyvalueExpr.End()
			groupedfields[groupPos] = append(groupedfields[groupPos], ident.Name)
		}
		lastPos = keyvalueExpr.End()
	}
	return groupedfields
}

func newField(field *ast.Field) *ast.Field {
	newfield := &ast.Field{
		Doc:   newDoc(field.Doc),
		Names: newNames(field.Names),
		Type:  newType(field.Type),
		Tag:   newTag(field.Tag),
		// Comment: newComment(field.Comment),
	}

	return newfield
}

func newTag(tag *ast.BasicLit) *ast.BasicLit {
	if tag == nil {
		return tag
	}
	return &ast.BasicLit{
		Value: tag.Value,
	}
}

func newType(expr ast.Expr) ast.Expr {
	return expr
}

func newNames(names []*ast.Ident) []*ast.Ident {
	ident := &ast.Ident{
		Name: names[0].Name,
	}
	return []*ast.Ident{ident}
}

func newDoc(doc *ast.CommentGroup) *ast.CommentGroup {
	if doc == nil {
		return nil
	}
	return &ast.CommentGroup{List: doc.List}
}
