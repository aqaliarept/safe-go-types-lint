# Plan: safe-go-types linter

> Source PRD: design session — safe-go-types linter grill

## Architectural decisions

- **Analyzer**: single `analysis.Analyzer` named `safe-go-types`, emitting multiple diagnostic codes
- **Diagnostic codes**: `no-scalar`, `no-zero-value`, `no-cast`, `no-constructor`, `untyped-literal`
- **Distribution**: standalone binary (`safe-go-types-lint`), golangci-lint `custom` linter
- **Constructor contract**: named `New<TypeName>` (exported or unexported), returns exactly `(TypeName, error)`, same package as type
- **Config**: golangci-lint YAML, glob `exclude-paths`, per-violation `//nolint:safe-go-types/<code>` suppression
- **Scalars covered**: `string`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `complex64`, `complex128`, `byte`, `rune`

---

## Phase 1: Scaffold + struct field scalar check + test project skeleton

**Covers**: project structure, binary wiring, first working diagnostic, test project foundation

### What to build

Set up the Go module with dependencies (`golang.org/x/tools/go/analysis`), implement the `analysis.Analyzer` skeleton, and wire a standalone binary entry point. Implement the first diagnostic — `no-scalar` — scoped to struct fields only. Create a `testdata/` project with representative struct definitions: valid (fields using custom types) and invalid (fields using raw scalars). The test project grows with each phase.

### Acceptance criteria

- [ ] `go build ./...` produces a `safe-go-types-lint` binary
- [ ] Binary accepts a package path and runs the analyzer
- [ ] `no-scalar` flags raw scalar struct fields: `type User struct { Name string }`
- [ ] `no-scalar` does not flag the underlying type in a type definition: `type Name string`
- [ ] `no-scalar` does not flag function parameters or return types
- [ ] `testdata/` package exists with at least: one struct with scalar fields (expect violations), one struct with custom-type fields (expect no violations)
- [ ] `analysistest` tests pass against `testdata/`

---

## Phase 2: Constructor registry + `no-constructor`

**Covers**: constructor detection, per-package type registry, `no-constructor` diagnostic

### What to build

Implement a per-package pass that collects all custom type definitions and all functions matching the `New<TypeName>` pattern returning `(TypeName, error)`. Build a registry mapping type names to their constructors. Emit `no-constructor` for any custom type (backed by a scalar) that has no valid constructor in the same package. Extend `testdata/` with: types that have valid constructors (no violation), types missing constructors (violation), types with malformed constructors — wrong return signature, wrong name prefix (violation).

### Acceptance criteria

- [ ] `no-constructor` flags `type Address string` when no `NewAddress`/`newAddress` exists in the package
- [ ] `no-constructor` does not flag when a valid unexported constructor exists: `func newAddress(val string) (Address, error)`
- [ ] Constructor with wrong return type (missing `error`, extra values) is not recognized — type is flagged
- [ ] Constructor named `CreateAddress` is not recognized — type is flagged
- [ ] Type defined as `type ShippingAddress Address` (underlying is custom type, not scalar) — still requires its own constructor
- [ ] `testdata/` covers all constructor edge cases with `// want` annotations

---

## Phase 3: Cast, zero-value, and untyped literal checks

**Covers**: `no-cast`, `no-zero-value`, `untyped-literal` diagnostics

### What to build

Using the constructor registry from Phase 2, implement three diagnostics:
- `no-zero-value`: flag bare declarations of custom types (`var a Address`, `address Address` as struct field with no init)
- `no-cast`: flag explicit conversions to a custom type outside its constructor (`Address("foo")` in any non-constructor context)
- `untyped-literal`: flag untyped constant literals assigned to a custom type variable (`var a Address = "foo"`) or passed where a custom type is expected (`SomeFunc("literal")`)

Same-package constants (`const Empty = Address("")`) are exempt from all three. Extend `testdata/` with cases for each.

### Acceptance criteria

- [ ] `no-zero-value` flags `var a Address`
- [ ] `no-zero-value` does not flag `a, err := NewAddress("x")`
- [ ] `no-cast` flags `Address("foo")` outside the constructor body
- [ ] `no-cast` does not flag `Address(val)` inside `NewAddress`
- [ ] `no-cast` does not flag reverse conversion: `string(addr)`
- [ ] `untyped-literal` flags `var a Address = "foo"`
- [ ] `untyped-literal` flags `SomeFunc("literal")` where `SomeFunc` expects `Address`
- [ ] `const Empty = Address("")` in same package produces no violations
- [ ] `testdata/` covers all cases with `// want` annotations

---

## Phase 4: Extend `no-scalar` to variable declarations and composite types

**Covers**: scalar detection in local variables, slices, maps, pointers, channels

### What to build

Extend `no-scalar` beyond struct fields to: explicit local variable declarations (`var x string`, `result := "pending"`), and composite types containing scalars (`[]string`, `map[string]int`, `*string`, `chan int`) in both struct fields and variable declarations. Inferred types from function calls (`x := someFunc()`) are not flagged. Extend `testdata/` with all composite and local variable cases.

### Acceptance criteria

- [ ] `no-scalar` flags `var x string` in a function body
- [ ] `no-scalar` flags `result := "pending"` (explicit scalar literal assignment)
- [ ] `no-scalar` does not flag `x := someFunc()` where `someFunc` returns a scalar
- [ ] `no-scalar` flags `[]string` as a struct field type
- [ ] `no-scalar` flags `map[string]int` as a variable declaration type
- [ ] `no-scalar` flags `*string` and `chan int`
- [ ] `no-scalar` flags nested composites: `map[string][]int`
- [ ] `testdata/` covers all composite and variable declaration cases

---

## Phase 5: Config, exclusions, and golangci-lint integration

**Covers**: config loading, glob exclusions, `//nolint` wiring, golangci-lint `custom` linter setup

### What to build

Implement config loading from golangci-lint settings (passed via analyzer flags): `exclude-paths` as glob patterns. Files matching any glob are skipped entirely. Verify that golangci-lint's built-in `//nolint:safe-go-types/no-scalar` suppression works per diagnostic code. Add a `testdata/excluded/` sub-package and verify it produces no diagnostics when matched by a glob. Write a `.golangci.yml` example and usage documentation.

### Acceptance criteria

- [ ] `exclude-paths: ["testdata/excluded/**"]` suppresses all diagnostics in matched files
- [ ] Glob patterns follow standard glob syntax (`*`, `**`, `?`)
- [ ] `//nolint:safe-go-types/no-scalar` suppresses only `no-scalar` violations on that line
- [ ] `//nolint:safe-go-types` suppresses all `safe-go-types` violations on that line
- [ ] `testdata/excluded/` package with violations produces no output when excluded via config
- [ ] Example `.golangci.yml` config included in repository
- [ ] Binary usage documented (flags, golangci-lint wiring)
