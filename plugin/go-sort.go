package main

import (
	"github.com/houtmann/go-sort/pkg"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		pkg.StructFieldsSortAnalyzer,
	}
}

// This must be defined and named 'AnalyzerPlugin' for golint-ci plugin
var AnalyzerPlugin analyzerPlugin
