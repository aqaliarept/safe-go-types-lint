Feature: no-zero-value diagnostic for zero-initialized custom types

  Background:
    Given the safe-go-types-lint binary is built
    And a Go package with a valid custom type and constructor:
      """
      type Address string

      func NewAddress(val string) (Address, error) {
          return Address(val), nil
      }
      """

  Scenario: bare var declaration of custom type is flagged
    Given the package also contains:
      """
      var a Address
      """
    When the linter runs on the package
    Then a "no-zero-value" diagnostic is reported on "var a Address"

  Scenario: custom type obtained via constructor is not flagged
    Given the package also contains:
      """
      func example() {
          a, err := NewAddress("x")
          _ = a
          _ = err
      }
      """
    When the linter runs on the package
    Then no "no-zero-value" diagnostic is reported

  Scenario: struct field of custom type with no initializer is flagged
    Given a Go package with:
      """
      type Address string
      func NewAddress(val string) (Address, error) { return Address(val), nil }

      type Order struct {
          Destination Address
      }

      var o Order
      """
    When the linter runs on the package
    Then a "no-zero-value" diagnostic is reported for the zero-value "Order" initialization
