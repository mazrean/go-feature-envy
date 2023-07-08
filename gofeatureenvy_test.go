package gofeatureenvy_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	gofeatureenvy "github.com/mazrean/go-feature-envy"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, gofeatureenvy.Analyzer, "a", "b")
}
