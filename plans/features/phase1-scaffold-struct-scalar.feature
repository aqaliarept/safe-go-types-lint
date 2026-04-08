Feature: no-scalar diagnostic for struct fields

  Background:
    Given the safe-go-types-lint binary is built

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
