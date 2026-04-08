Feature: no-scalar diagnostic for raw scalar types

  Background:
    Given the safe-go-types-lint binary is built

  # Struct fields

  Scenario: struct field with raw string is flagged
    Given a Go package with:
      """
      type User struct {
          Name string
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on the "Name string" field

  Scenario: struct field with raw int is flagged
    Given a Go package with:
      """
      type Order struct {
          Quantity int
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on the "Quantity int" field

  Scenario: struct field with custom type is not flagged
    Given a Go package with:
      """
      type Name string
      type User struct {
          Name Name
      }
      """
    When the linter runs on the package
    Then no "no-scalar" diagnostic is reported

  Scenario: underlying type in a type definition is not flagged
    Given a Go package with:
      """
      type Address string
      """
    When the linter runs on the package
    Then no "no-scalar" diagnostic is reported for the type definition

  Scenario: function parameter with raw scalar is not flagged
    Given a Go package with:
      """
      func Process(name string) {}
      """
    When the linter runs on the package
    Then no "no-scalar" diagnostic is reported

  Scenario: function return type with raw scalar is not flagged
    Given a Go package with:
      """
      func GetName() string { return "" }
      """
    When the linter runs on the package
    Then no "no-scalar" diagnostic is reported

  Scenario: all scalar types in struct fields are flagged
    Given a Go package with a struct containing fields of types:
      | field type |
      | string     |
      | bool       |
      | int        |
      | int8       |
      | int16      |
      | int32      |
      | int64      |
      | uint       |
      | uint8      |
      | uint16     |
      | uint32     |
      | uint64     |
      | float32    |
      | float64    |
      | complex64  |
      | complex128 |
      | byte       |
      | rune       |
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported for each field

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

  # Composite types — slices

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

  # Composite types — maps

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

  # Composite types — pointers

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

  # Composite types — channels

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

  # Composite types — nested

  Scenario: nested composite with scalar is flagged
    Given a Go package with:
      """
      type Index struct {
          Lookup map[string][]int
      }
      """
    When the linter runs on the package
    Then a "no-scalar" diagnostic is reported on "Lookup map[string][]int"

  # Exempt: type definitions, function signatures, inferred types

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
