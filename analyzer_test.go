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

func TestNoConstructorDerived(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_derived")
}

func TestNoZeroValueBareVar(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_zero_value_bare_var")
}

func TestNoZeroValueConstructorNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_zero_value_constructor_not_flagged")
}

func TestNoZeroValueStructField(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_zero_value_struct_field")
}

func TestNoCastOutsideConstructor(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_cast_outside_constructor")
}

func TestNoCastInsideConstructorNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_cast_inside_constructor_not_flagged")
}

func TestNoCastReverseConversionNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_cast_reverse_conversion_not_flagged")
}

func TestUntypedLiteralVarAssign(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "untyped_literal_var_assign")
}

func TestUntypedLiteralFuncArg(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "untyped_literal_func_arg")
}

func TestSamePackageConstantNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "same_package_constant_not_flagged")
}

func TestSamePackageConstantUsedNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "same_package_constant_used_not_flagged")
}

func TestVarScalarExplicit(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "var_scalar_explicit")
}

func TestVarScalarShortAssign(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "var_scalar_short_assign")
}

func TestVarScalarFuncCallNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "var_scalar_func_call")
}

func TestSliceScalarStructField(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "slice_scalar_struct_field")
}

func TestSliceScalarLocalVar(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "slice_scalar_local_var")
}

func TestSliceCustomTypeNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "slice_custom_type_not_flagged")
}

func TestMapScalarStructField(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "map_scalar_struct_field")
}

func TestMapScalarLocalVar(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "map_scalar_local_var")
}

func TestPtrScalarStructField(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "ptr_scalar_struct_field")
}

func TestPtrScalarLocalVar(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "ptr_scalar_local_var")
}

func TestChanScalarStructField(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "chan_scalar_struct_field")
}

func TestChanScalarLocalVar(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "chan_scalar_local_var")
}

func TestChanVarDeclNotDoubleReported(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "chan_var_decl_no_double")
}

func TestNestedCompositeScalar(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "nested_composite_scalar")
}

func TestSliceCustomTypeLocalVarNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "slice_custom_type_local_var")
}
