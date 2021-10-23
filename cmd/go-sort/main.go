package main

import (
	"github.com/houtmann/go-sort/pkg"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(pkg.Analyzer)
}
