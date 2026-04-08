Feature: no-scalar diagnostic for local variables and composite types

  Background:
    Given the safe-go-types-lint binary is built

  # Local variable declarations

  Scenario: explicit scalar var declaration in function body is flagged
    Given a Go package with:
      """
      func example() {
          var x string
          _ = x
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "var x string"

  Scenario: short assignment with scalar literal is flagged
    Given a Go package with:
      """
      func example() {
          result := "pending"
          _ = result
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "result := \"pending\""

  Scenario: inferred type from function call is not flagged
    Given a Go package with:
      """
      func getValue() string { return "" }

      func example() {
          x := getValue()
          _ = x
      }
      """
    When the linter runs on the package
    Then no "no-scalar" diagnostic is reported on "x := getValue()"

  # Slice composite types

  Scenario: slice of scalar as struct field is flagged
    Given a Go package with:
      """
      type Order struct {
          Tags []string
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "Tags []string"

  Scenario: slice of scalar as local variable is flagged
    Given a Go package with:
      """
      func example() {
          var tags []string
          _ = tags
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "var tags []string"

  Scenario: slice of custom type is not flagged
    Given a Go package with:
      """
      type Tag string

      type Order struct {
          Tags []Tag
      }
      """
    When the linter runs on the package
    Then no "no-scalar" diagnostic is reported

  # Map composite types

  Scenario: map with scalar key or value as struct field is flagged
    Given a Go package with:
      """
      type Config struct {
          Settings map[string]string
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "Settings map[string]string"

  Scenario: map with scalar value as local variable is flagged
    Given a Go package with:
      """
      func example() {
          counts := map[string]int{}
          _ = counts
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "counts := map[string]int{}"

  # Pointer types

  Scenario: pointer to scalar as struct field is flagged
    Given a Go package with:
      """
      type Record struct {
          Name *string
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "Name *string"

  Scenario: pointer to scalar as local variable is flagged
    Given a Go package with:
      """
      func example() {
          var x *string
          _ = x
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "var x *string"

  # Channel types

  Scenario: channel of scalar as struct field is flagged
    Given a Go package with:
      """
      type Worker struct {
          Jobs chan int
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "Jobs chan int"

  Scenario: channel of scalar as local variable is flagged
    Given a Go package with:
      """
      func example() {
          ch := make(chan string)
          _ = ch
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported

  # Nested composites

  Scenario: nested composite with scalar is flagged
    Given a Go package with:
      """
      type Index struct {
          Lookup map[string][]int
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "Lookup map[string][]int"

  # Non-scalar composites are not flagged

  Scenario: slice of custom type as local variable is not flagged
    Given a Go package with:
      """
      type Tag string

      func example() {
          var tags []Tag
          _ = tags
      }
      """
    When the linter runs on the package
    Then no "no-scalar" diagnostic is reported
