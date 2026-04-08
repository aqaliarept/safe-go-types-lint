Feature: untyped-literal diagnostic for untyped constant literals used as custom types

  Background:
    Given the safe-go-types-lint binary is built
    And a Go package with a valid custom type and constructor:
      """
      type Address string

      func NewAddress(val string) (Address, error) {
          return Address(val), nil
      }
      """

  Scenario: untyped string literal assigned to custom type variable is flagged
    Given the package also contains:
      """
      var a Address = "foo"
      """
    When the linter runs on the package
    Then an "untyped-literal" diagnostic is reported on the assignment

  Scenario: untyped literal passed as custom type argument is flagged
    Given the package also contains:
      """
      func consume(a Address) {}

      func example() {
          consume("foo")
      }
      """
    When the linter runs on the package
    Then an "untyped-literal" diagnostic is reported on the "consume" call

  Scenario: same-package constant is not flagged
    Given the package also contains:
      """
      const EmptyAddress = Address("")
      """
    When the linter runs on the package
    Then no diagnostic is reported for "EmptyAddress"

  Scenario: constant from same package used as value is not flagged
    Given the package also contains:
      """
      const EmptyAddress = Address("")

      func example() {
          a := EmptyAddress
          _ = a
      }
      """
    When the linter runs on the package
    Then no diagnostic is reported for the use of "EmptyAddress"
