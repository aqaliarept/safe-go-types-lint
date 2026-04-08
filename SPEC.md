# safe-go-types Linter Specification

## Purpose

`safe-go-types` is a Go static analysis linter that enforces the use of custom named types over raw scalar types. Its goal is to prevent type invariant violations: a domain value (e.g. `Address`, `UserID`, `Amount`) can only be created through a validated constructor, making it impossible to accidentally pass the wrong value to the wrong function or leave a domain object in an invalid zero state.

---

## Diagnostic codes

| Code | Description |
|------|-------------|
| `safe-go-types/no-scalar` | A raw scalar type is used where a custom type is required (struct field, variable declaration, composite type) |
| `safe-go-types/no-constructor` | A custom type is defined but has no valid constructor in the same package |
| `safe-go-types/no-zero-value` | A variable of a custom type is declared with its zero value (bypassing the constructor) |
| `safe-go-types/no-cast` | An explicit conversion to a custom type occurs outside its constructor |
| `safe-go-types/untyped-literal` | An untyped constant literal is assigned to a custom type or passed where a custom type is expected |

---

## Scalar types covered

All Go built-in scalar types:

`string`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `complex64`, `complex128`, `byte`, `rune`

---

## Rules

### `no-scalar`

**Scenarios**: [`plans/features/no-scalar.feature`](plans/features/no-scalar.feature)

**Flagged:**
- Struct fields whose type is a raw scalar: `type User struct { Name string }`
- Struct fields whose type is a composite containing a scalar: `Tags []string`, `Settings map[string]int`, `Ptr *string`, `Ch chan int`
- Explicit local variable declarations with a scalar type: `var x string`
- Short assignments where the RHS is a scalar literal: `result := "pending"`
- Local variable composite types containing scalars: `var tags []string`, `counts := map[string]int{}`
- Nested composite types: `map[string][]int`

**Not flagged:**
- The underlying type in a custom type definition: `type Address string`
- Function or method parameters and return types
- Types inferred from function calls: `x := someFunc()` (even if `someFunc` returns a scalar)
- Custom types used as field or variable types: `type Name string; type User struct { Name Name }`

---

### `no-constructor`

**Scenarios**: [`plans/features/no-constructor.feature`](plans/features/no-constructor.feature)

**Flagged:**
- Any custom type `type Foo <scalar>` that has no valid constructor in the same package
- Any custom type `type Foo Bar` (derived from another custom type) that has no valid constructor

**Not flagged:**
- Custom types with a valid constructor (see [Constructor contract](#constructor-contract))

---

### `no-zero-value`

**Scenarios**: [`plans/features/no-zero-value.feature`](plans/features/no-zero-value.feature)

**Flagged:**
- Bare variable declarations of a custom type: `var a Address`
- Struct variable declarations where the struct contains custom-type fields: `var o Order`

**Not flagged:**
- Variables initialized via constructor: `a, err := NewAddress("x")`

---

### `no-cast`

**Scenarios**: [`plans/features/no-cast.feature`](plans/features/no-cast.feature)

**Flagged:**
- Explicit conversion to a custom type outside its constructor body: `a := Address("foo")`

**Not flagged:**
- Conversion inside the constructor body: `return Address(val), nil` inside `NewAddress`
- Reverse conversion from a custom type to a scalar: `string(addr)`
- Same-package constant declarations: `const Empty = Address("")`

---

### `untyped-literal`

**Scenarios**: [`plans/features/untyped-literal.feature`](plans/features/untyped-literal.feature)

**Flagged:**
- Untyped literal assigned to a custom type variable: `var a Address = "foo"`
- Untyped literal passed as argument where a custom type is expected: `consume("foo")` where `func consume(a Address)`

**Not flagged:**
- Same-package constant declarations: `const Empty = Address("")`
- Using a same-package constant as a value: `a := EmptyAddress`

---

## Constructor contract

A valid constructor for type `Foo` must satisfy all of the following:

- **Name**: `NewFoo` (exported) or `newFoo` (unexported) — exactly the `New`/`new` prefix followed by the type name
- **Return signature**: exactly `(Foo, error)` — no additional return values
- **Location**: same package as the type definition

Any function that fails any of these conditions is not recognized as a constructor. The type is then treated as having no constructor and emits `no-constructor`.

**Examples of invalid constructors:**
- `func NewAddress(val string) Address` — missing `error` return
- `func NewAddress(val string) (Address, string, error)` — extra return value
- `func CreateAddress(val string) (Address, error)` — wrong name prefix
- Constructor defined in a different package than the type

---

## Diagnostic message formats

| Diagnostic | Example message |
|------------|----------------|
| `no-scalar` | `field "Name" has raw scalar type "string"` |
| `no-constructor` | `type "Address" has no valid constructor in this package` |
| `no-zero-value` | `variable "a" of custom type "Address" declared with zero value` |
| `no-cast` | `explicit cast to custom type "Address" outside its constructor` |
| `untyped-literal` | `untyped literal used where custom type "Address" is expected` |

All messages are prefixed with the diagnostic code by the analysis framework: e.g. `safe-go-types/no-scalar: field "Name" has raw scalar type "string"`.

---

## Configuration

**Scenarios**: [`plans/features/configuration.feature`](plans/features/configuration.feature)

### Excluding paths

Pass the `-exclude-paths` flag to the analyzer with a comma-separated list of glob patterns. Files whose paths match any pattern produce no diagnostics.

```
-exclude-paths="**/generated/**,**/vendor/**"
```

Glob patterns support `*`, `**`, and `?`. `**` matches across directory separators.

In golangci-lint:
```yaml
linters-settings:
  safe-go-types:
    exclude-paths:
      - "**/generated/**"
      - "**/vendor/**"
```

### Suppressing individual diagnostics

Use standard `//nolint` comments:

```go
// Suppress a specific diagnostic on this line:
Name string //nolint:safe-go-types/no-scalar

// Suppress all safe-go-types diagnostics on this line:
Name string //nolint:safe-go-types
```

Comments for other linters (e.g. `//nolint:some-other-linter`) do not suppress `safe-go-types` diagnostics.

---

## Distribution

The linter is distributed as a standalone binary `safe-go-types-lint`. Install with:

```
go install safe-go-types-lint/cmd/safe-go-types-lint@latest
```

### golangci-lint integration

Configure as a `custom-linters` entry in `.golangci.yml`:

```yaml
linters:
  enable:
    - safe-go-types

custom-linters:
  safe-go-types:
    path: ./safe-go-types-lint
    description: Enforce custom types over raw scalars

linters-settings:
  safe-go-types:
    exclude-paths:
      - "**/generated/**"
```

---

## Implementation notes

### Analysis framework

The linter is implemented as a single `analysis.Analyzer` named `safe-go-types` using `golang.org/x/tools/go/analysis`. All five diagnostic codes are emitted by this one analyzer. This keeps the implementation cohesive and avoids the complexity of shared state across multiple analyzers.

### Pass order within a package

The analyzer runs several checks in sequence over each package:

1. **Struct field scalar check** — walks `ast.StructType` nodes, flags fields whose type contains a raw scalar.
2. **Local variable scalar check** — walks function bodies, flags explicit var declarations and short assignments with scalar types.
3. **Constructor registry** — two sub-passes: first collect all custom type definitions (scalar-backed or derived), then collect all valid `New<T>`/`new<T>` functions. Types with no matching constructor emit `no-constructor`.
4. **Zero-value check** — using the custom type registry, flags bare `var x CustomType` declarations.
5. **Cast check** — flags explicit conversions to custom types outside constructor bodies.
6. **Untyped literal check** — flags `BasicLit` nodes assigned to or passed as custom types.

The custom type registry (step 3) is computed once and shared with steps 4, 5, and 6.

### Why two passes for custom types

A package may define `type ShippingAddress Address` after `type Address string`. A single forward pass would miss `ShippingAddress` as a custom type because `Address` hasn't been registered yet. The two-pass approach (first collect all scalar-backed types, then expand transitively) handles arbitrary definition order.

### Glob matching

Path exclusion uses `github.com/bmatcuk/doublestar/v4` for `**` support, since the standard `filepath.Match` does not handle double-star patterns.

### Testing approach

Tests use `golang.org/x/tools/go/analysis/analysistest`. Each test function corresponds to one Gherkin scenario. Test data lives in `testdata/src/<package>/` with `// want` annotations marking expected diagnostics.
