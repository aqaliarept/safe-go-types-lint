Feature: configuration, path exclusions, and golangci-lint integration

  Background:
    Given the safe-go-types-lint binary is built

  # Path exclusions

  Scenario: file matching exclude-paths glob produces no diagnostics
    Given a Go package at path "internal/legacy/types.go" with scalar struct fields
    And the linter is configured with:
      """
      exclude-paths:
        - "internal/legacy/**"
      """
    When the linter runs on the package
    Then no diagnostics are reported for files under "internal/legacy/"

  Scenario: file not matching any exclude-paths glob is still checked
    Given a Go package at path "internal/domain/types.go" with scalar struct fields
    And the linter is configured with:
      """
      exclude-paths:
        - "internal/legacy/**"
      """
    When the linter runs on the package
    Then diagnostics are reported for "internal/domain/types.go"

  Scenario: multiple glob patterns are all applied
    Given Go packages at paths:
      | path                        |
      | internal/legacy/types.go    |
      | pkg/generated/models.go     |
      | internal/domain/types.go    |
    And the linter is configured with:
      """
      exclude-paths:
        - "internal/legacy/**"
        - "pkg/generated/**"
      """
    When the linter runs on all packages
    Then no diagnostics are reported for "internal/legacy/types.go"
    And no diagnostics are reported for "pkg/generated/models.go"
    And diagnostics are reported for "internal/domain/types.go"

  Scenario: wildcard glob pattern matches nested paths
    Given a Go file at path "third_party/vendor/lib/types.go" with violations
    And the linter is configured with:
      """
      exclude-paths:
        - "**/vendor/**"
      """
    When the linter runs
    Then no diagnostics are reported for the file

  # nolint comments — per diagnostic code

  Scenario: nolint comment for specific code suppresses only that diagnostic
    Given a Go package with:
      """
      type User struct {
          Name string //nolint:safe-go-types/no-scalar
          Age  int
      }
      """
    When the linter runs on the package
    Then no "no-scalar" diagnostic is reported for "Name string"
    And a "no-scalar" diagnostic is reported for "Age int"

  Scenario: blanket nolint comment suppresses all safe-go-types diagnostics on that line
    Given a Go package with:
      """
      type User struct {
          Name string //nolint:safe-go-types
      }
      """
    When the linter runs on the package
    Then no diagnostic is reported for "Name string"

  Scenario: nolint comment for a different linter does not suppress safe-go-types
    Given a Go package with:
      """
      type User struct {
          Name string //nolint:some-other-linter
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported for "Name string"

  # golangci-lint custom linter wiring

  Scenario: linter runs as a golangci-lint custom linter
    Given a repository with a ".golangci.yml" containing:
      """
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
      """
    When golangci-lint runs on the repository
    Then safe-go-types violations are reported in the golangci-lint output
    And violations in paths matching "**/generated/**" are not reported
