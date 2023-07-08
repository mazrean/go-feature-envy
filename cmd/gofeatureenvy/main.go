package main

import (
	gofeatureenvy "github.com/mazrean/go-feature-envy"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(gofeatureenvy.Analyzer) }
