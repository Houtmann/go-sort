package pkg

import (
	"go/ast"
	"sort"
)

type Field struct {
	ast.Field
}

func (f *Field) Identifier() string {
	if len(f.Names) > 0 {
		return f.Names[0].Name
	}
	return ""
}

type BasicLit struct {
	ast.BasicLit
}

func (b *BasicLit) Identifier() string {
	return b.Value
}

type Ident struct {
	ast.Ident
}

func (i *Ident) Identifier() string {
	return i.Name
}

type Elem interface {
	Identifier() string
}

func IsSorted(sortable []Elem) bool {
	issorted := sort.SliceIsSorted(sortable, func(i, j int) bool {
		return sortable[i].Identifier() < sortable[j].Identifier()
	})
	return issorted
}

func sortElement(sortable []Elem) []string {
	var sortedNames []string
	for i := range sortable {
		sortedNames = append(sortedNames, sortable[i].Identifier())
	}
	sort.Strings(sortedNames)
	return sortedNames
}
