package main

import (
	"github.com/Houtmann/go-sort/pkg"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(pkg.StructFieldsSortAnalyzer, pkg.SliceSortAnalyzer)
}
