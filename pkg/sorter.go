package pkg

import (
	"go/ast"
	"sort"
)

type Field struct {
	ast.Field
}

func (f *Field) Name() string {
	if len(f.Names) > 0 {
		return f.Names[0].Name
	}
	return ""
}

type BasicLit struct {
	ast.BasicLit
}

func (b *BasicLit) Name() string {
	return b.Value
}

type Ident struct {
	ast.Ident
}

func (i *Ident) Name() string {
	return i.Name()
}

type Sortable interface {
	Name() string
}

func isSorted[T Sortable](sortable []T) bool {
	return sort.SliceIsSorted(sortable, func(i, j int) bool {
		return sortable[i].Name() < sortable[j].Name()
	})
}
