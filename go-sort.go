package main

import (
	"github.com/houtmann/go-sort/pkg"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(pkg.SliceSortAnalyzer)
}
