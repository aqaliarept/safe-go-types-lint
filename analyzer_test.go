package safegotypes_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	safegotypes "safe-go-types-lint"
)

func testdata() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

func TestStructFieldRawString(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "struct_string_field")
}

func TestStructFieldRawInt(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "struct_int_field")
}

func TestStructFieldCustomTypeNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "struct_custom_type_field")
}

func TestTypeDefinitionNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "type_definition")
}

func TestFuncParamNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "func_param")
}

func TestFuncReturnNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "func_return")
}

func TestAllScalarTypesInStructFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "all_scalars")
}

func TestNoConstructorBasic(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_basic")
}

func TestNoConstructorExported(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_exported")
}

func TestNoConstructorUnexported(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_unexported")
}

func TestNoConstructorMissingError(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_missing_error")
}

func TestNoConstructorExtraReturns(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_extra_returns")
}

func TestNoConstructorWrongName(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_wrong_name")
}

func TestNoConstructorDifferentPackage(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "domain")
}
