Feature: no-zero-value, no-cast, and untyped-literal diagnostics

  Background:
    Given the safe-go-types-lint binary is built
    And a Go package with a valid custom type and constructor:
      """
      type Address string

      func NewAddress(val string) (Address, error) {
          return Address(val), nil
      }
      """

  # no-zero-value

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

  # no-cast

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

  # untyped-literal

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
