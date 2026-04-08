package safegotypes_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	safegotypes "safe-go-types-lint"
)

func testdata() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

// Scenario: struct field with raw string is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_StructFieldWithRawStringIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "struct_string_field")
}

// Scenario: struct field with raw int is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_StructFieldWithRawIntIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "struct_int_field")
}

// Scenario: struct field with custom type is not flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_StructFieldWithCustomTypeIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "struct_custom_type_field")
}

// Scenario: underlying type in a type definition is not flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_UnderlyingTypeInDefinitionIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "type_definition")
}

// Scenario: function parameter with raw scalar is not flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_FunctionParameterIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "func_param")
}

// Scenario: function return type with raw scalar is not flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_FunctionReturnTypeIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "func_return")
}

// Scenario: all scalar types in struct fields are flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_AllScalarTypesInStructFieldsAreFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "all_scalars")
}

// Scenario: custom type with no constructor is flagged
// Feature: plans/features/no-constructor.feature
func TestNoConstructor_CustomTypeWithNoConstructorIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_basic")
}

// Scenario: custom type with valid exported constructor is not flagged
// Feature: plans/features/no-constructor.feature
func TestNoConstructor_ValidExportedConstructorIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_exported")
}

// Scenario: custom type with valid unexported constructor is not flagged
// Feature: plans/features/no-constructor.feature
func TestNoConstructor_ValidUnexportedConstructorIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_unexported")
}

// Scenario: constructor missing error return is not recognized
// Feature: plans/features/no-constructor.feature
func TestNoConstructor_ConstructorMissingErrorReturnIsNotRecognized(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_missing_error")
}

// Scenario: constructor with extra return values is not recognized
// Feature: plans/features/no-constructor.feature
func TestNoConstructor_ConstructorWithExtraReturnValuesIsNotRecognized(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_extra_returns")
}

// Scenario: constructor with wrong name prefix is not recognized
// Feature: plans/features/no-constructor.feature
func TestNoConstructor_ConstructorWithWrongNamePrefixIsNotRecognized(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_wrong_name")
}

// Scenario: constructor in a different package is not recognized
// Feature: plans/features/no-constructor.feature
func TestNoConstructor_ConstructorInDifferentPackageIsNotRecognized(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "domain")
}

// Scenario: custom type derived from another custom type requires its own constructor
// Feature: plans/features/no-constructor.feature
func TestNoConstructor_DerivedCustomTypeRequiresItsOwnConstructor(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_constructor_derived")
}

// Scenario: bare var declaration of custom type is flagged
// Feature: plans/features/no-zero-value.feature
func TestNoZeroValue_BareVarDeclarationIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_zero_value_bare_var")
}

// Scenario: custom type obtained via constructor is not flagged
// Feature: plans/features/no-zero-value.feature
func TestNoZeroValue_CustomTypeObtainedViaConstructorIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_zero_value_constructor_not_flagged")
}

// Scenario: struct field of custom type with no initializer is flagged
// Feature: plans/features/no-zero-value.feature
func TestNoZeroValue_StructFieldWithNoInitializerIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_zero_value_struct_field")
}

// Scenario: explicit cast to custom type outside constructor is flagged
// Feature: plans/features/no-cast.feature
func TestNoCast_ExplicitCastOutsideConstructorIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_cast_outside_constructor")
}

// Scenario: cast inside constructor body is not flagged
// Feature: plans/features/no-cast.feature
func TestNoCast_CastInsideConstructorBodyIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_cast_inside_constructor_not_flagged")
}

// Scenario: reverse conversion from custom type to scalar is not flagged
// Feature: plans/features/no-cast.feature
func TestNoCast_ReverseConversionFromCustomTypeIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "no_cast_reverse_conversion_not_flagged")
}

// Scenario: untyped string literal assigned to custom type variable is flagged
// Feature: plans/features/untyped-literal.feature
func TestUntypedLiteral_UntypedStringLiteralAssignedToCustomTypeIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "untyped_literal_var_assign")
}

// Scenario: untyped literal passed as custom type argument is flagged
// Feature: plans/features/untyped-literal.feature
func TestUntypedLiteral_UntypedLiteralPassedAsCustomTypeArgumentIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "untyped_literal_func_arg")
}

// Scenario: same-package constant is not flagged
// Feature: plans/features/untyped-literal.feature
func TestUntypedLiteral_SamePackageConstantIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "same_package_constant_not_flagged")
}

// Scenario: constant from same package used as value is not flagged
// Feature: plans/features/untyped-literal.feature
func TestUntypedLiteral_ConstantFromSamePackageUsedAsValueIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "same_package_constant_used_not_flagged")
}

// Scenario: explicit scalar var declaration in function body is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_ExplicitScalarVarDeclarationInFunctionBodyIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "var_scalar_explicit")
}

// Scenario: short assignment with scalar literal is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_ShortAssignmentWithScalarLiteralIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "var_scalar_short_assign")
}

// Scenario: inferred type from function call is not flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_InferredTypeFromFunctionCallIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "var_scalar_func_call")
}

// Scenario: slice of scalar as struct field is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_SliceOfScalarAsStructFieldIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "slice_scalar_struct_field")
}

// Scenario: slice of scalar as local variable is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_SliceOfScalarAsLocalVariableIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "slice_scalar_local_var")
}

// Scenario: slice of custom type is not flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_SliceOfCustomTypeIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "slice_custom_type_not_flagged")
}

// Scenario: map with scalar key or value as struct field is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_MapWithScalarAsStructFieldIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "map_scalar_struct_field")
}

// Scenario: map with scalar value as local variable is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_MapWithScalarAsLocalVariableIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "map_scalar_local_var")
}

// Scenario: pointer to scalar as struct field is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_PointerToScalarAsStructFieldIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "ptr_scalar_struct_field")
}

// Scenario: pointer to scalar as local variable is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_PointerToScalarAsLocalVariableIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "ptr_scalar_local_var")
}

// Scenario: channel of scalar as struct field is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_ChannelOfScalarAsStructFieldIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "chan_scalar_struct_field")
}

// Scenario: channel of scalar as local variable is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_ChannelOfScalarAsLocalVariableIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "chan_scalar_local_var")
}

// Scenario: channel of scalar as local variable is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_ChannelVarDeclNotDoubleReported(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "chan_var_decl_no_double")
}

// Scenario: nested composite with scalar is flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_NestedCompositeWithScalarIsFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "nested_composite_scalar")
}

// Scenario: slice of custom type as local variable is not flagged
// Feature: plans/features/no-scalar.feature
func TestNoScalar_SliceOfCustomTypeAsLocalVariableIsNotFlagged(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "slice_custom_type_local_var")
}

// Scenario: file matching exclude-paths glob produces no diagnostics
// Feature: plans/features/configuration.feature
func TestConfiguration_FileMatchingExcludePathsProducesNoDiagnostics(t *testing.T) {
	_ = safegotypes.Analyzer.Flags.Set("exclude-paths", "legacy/**")
	t.Cleanup(func() { _ = safegotypes.Analyzer.Flags.Set("exclude-paths", "") })
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "legacy")
}

// Scenario: file not matching any exclude-paths glob is still checked
// Feature: plans/features/configuration.feature
func TestConfiguration_FileNotMatchingExcludePathsIsStillChecked(t *testing.T) {
	_ = safegotypes.Analyzer.Flags.Set("exclude-paths", "legacy/**")
	t.Cleanup(func() { _ = safegotypes.Analyzer.Flags.Set("exclude-paths", "") })
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "domain_check")
}

// Scenario: multiple glob patterns are all applied
// Feature: plans/features/configuration.feature
func TestConfiguration_MultipleGlobPatternsAreAllApplied(t *testing.T) {
	_ = safegotypes.Analyzer.Flags.Set("exclude-paths", "legacy2/**,generated2/**")
	t.Cleanup(func() { _ = safegotypes.Analyzer.Flags.Set("exclude-paths", "") })
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "legacy2")
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "generated2")
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "domain2")
}

// Scenario: wildcard glob pattern matches nested paths
// Feature: plans/features/configuration.feature
func TestConfiguration_WildcardGlobPatternMatchesNestedPaths(t *testing.T) {
	_ = safegotypes.Analyzer.Flags.Set("exclude-paths", "**/vendored/**")
	t.Cleanup(func() { _ = safegotypes.Analyzer.Flags.Set("exclude-paths", "") })
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "vendored")
}

// Scenario: nolint comment for specific code suppresses only that diagnostic
// Feature: plans/features/configuration.feature
func TestConfiguration_NolintCommentForSpecificCodeSuppressesOnlyThatDiagnostic(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "nolint_specific")
}

// Scenario: blanket nolint comment suppresses all safe-go-types diagnostics on that line
// Feature: plans/features/configuration.feature
func TestConfiguration_BlanketNolintCommentSuppressesAllSafeGoTypesDiagnostics(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "nolint_blanket")
}

// Scenario: nolint comment for a different linter does not suppress safe-go-types
// Feature: plans/features/configuration.feature
func TestConfiguration_NolintCommentForDifferentLinterDoesNotSuppressSafeGoTypes(t *testing.T) {
	analysistest.Run(t, testdata(), safegotypes.Analyzer, "nolint_other")
}

// Scenario: linter runs as a golangci-lint custom linter
// Feature: plans/features/configuration.feature
func TestConfiguration_GolangciLintYmlExists(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	yamlPath := filepath.Join(root, ".golangci.yml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf(".golangci.yml not found: %v", err)
	}
	if !strings.Contains(string(data), "safe-go-types") {
		t.Errorf(".golangci.yml does not mention safe-go-types")
	}
}
