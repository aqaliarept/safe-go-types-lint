package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	safegotypes "github.com/aqaliarept/safe-go-types-lint"
)

func main() {
	singlechecker.Main(safegotypes.Analyzer)
}
