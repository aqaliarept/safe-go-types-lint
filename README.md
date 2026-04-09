# safe-go-types-lint

A Go static analysis linter that enforces the use of custom named types over raw scalar types, preventing type invariant violations at compile time.

## Why

Go allows passing any `string` where a `string` is expected — which means nothing stops you from passing an email address where a username is required, or leaving a domain object in an invalid zero state. This linter makes that impossible by requiring every domain value to be represented by a custom type and created only through a validated constructor.

## Example

```go
// ❌ flagged by safe-go-types
type User struct {
    Email    string // no-scalar: raw scalar in struct field
    Username string // no-scalar: raw scalar in struct field
}

var u User // no-zero-value: zero-value custom type

// ✅ correct
type Email string

func NewEmail(val string) (Email, error) {
    if val == "" {
        return Email(""), errors.New("email cannot be empty")
    }
    return Email(val), nil
}

type Username string

func NewUsername(val string) (Username, error) {
    if val == "" {
        return Username(""), errors.New("username cannot be empty")
    }
    return Username(val), nil
}

type User struct {
    Email    Email
    Username Username
}
```

## Installation

```sh
go install github.com/aqaliarept/safe-go-types-lint/cmd/safe-go-types-lint@latest
```

## Standalone usage

Run directly on a package or module:

```sh
safe-go-types-lint ./...
```

With path exclusions:

```sh
safe-go-types-lint -exclude-paths="**/generated/**,**/vendor/**" ./...
```

## golangci-lint integration

This linter integrates with golangci-lint v2 via the [module plugin system](https://golangci-lint.run/plugins/module-plugins/).

**1.** Create `custom-gcl.yml` in your project root:

```yaml
version: v2.8.0
name: custom-golangci-lint
destination: ./bin/

plugins:
  - module: github.com/aqaliarept/safe-go-types-lint
    import: github.com/aqaliarept/safe-go-types-lint/plugin
    version: latest
```

**2.** Add to your `.golangci.yml`:

```yaml
version: "2"

linters:
  default: none
  enable:
    - safe-go-types
  settings:
    custom:
      safe-go-types:
        type: module
        description: Enforce custom types over raw scalars
        original-url: github.com/aqaliarept/safe-go-types-lint
        settings:
          exclude-paths:
            - "**/generated/**"
            - "**/vendor/**"
```

**3.** Build the custom golangci-lint binary and run it:

```sh
golangci-lint custom
./bin/custom-golangci-lint run ./...
```

## Configuration

### Excluding paths

Use `-exclude-paths` with comma-separated glob patterns. Files matching any pattern produce no diagnostics. Supports `*`, `**`, and `?`.

```sh
safe-go-types-lint -exclude-paths="**/generated/**,internal/legacy/**" ./...
```

In golangci-lint, use the `exclude-paths` setting under `linters.settings.custom.safe-go-types.settings:` (see above).

### Suppressing individual violations

Use standard `//nolint` comments:

```go
// Suppress a specific diagnostic:
Name string //nolint:safe-go-types/no-scalar

// Suppress all safe-go-types diagnostics on this line:
Name string //nolint:safe-go-types
```

## Diagnostic codes

| Code | Description |
|------|-------------|
| `safe-go-types/no-scalar` | Raw scalar type used in a struct field, variable declaration, or composite type |
| `safe-go-types/no-constructor` | Custom type defined without a valid `New<Type>(val) (Type, error)` constructor in the same package |
| `safe-go-types/no-zero-value` | Variable of a custom type declared with its zero value, bypassing the constructor |
| `safe-go-types/no-cast` | Explicit conversion to a custom type outside its constructor body |
| `safe-go-types/untyped-literal` | Untyped constant literal assigned to or passed where a custom type is expected |

For full specification including constructor contract, exempt patterns, and diagnostic message formats, see [SPEC.md](SPEC.md).

## License

MIT
