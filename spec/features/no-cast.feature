Feature: no-cast diagnostic for explicit conversions to custom types

  Background:
    Given the safe-go-types-lint binary is built
    And a Go package with a valid custom type and constructor:
      """
      type Address string

      func NewAddress(val string) (Address, error) {
          return Address(val), nil
      }
      """

  Scenario: explicit cast to custom type outside constructor is flagged
    Given the package also contains:
      """
      func example() {
          a := Address("foo")
          _ = a
      }
      """
    When the linter runs on the package
    Then a "no-cast" diagnostic is reported on "Address(\"foo\")"

  Scenario: cast inside constructor body is not flagged
    When the linter runs on the package
    Then no "no-cast" diagnostic is reported inside "NewAddress"

  Scenario: reverse conversion from custom type to scalar is not flagged
    Given the package also contains:
      """
      func example(a Address) string {
          return string(a)
      }
      """
    When the linter runs on the package
    Then no "no-cast" diagnostic is reported on "string(a)"
