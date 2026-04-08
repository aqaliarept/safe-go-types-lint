Feature: no-constructor diagnostic for custom types without valid constructors

  Background:
    Given the safe-go-types-lint binary is built

  Scenario: custom type with no constructor is flagged
    Given a Go package with:
      """
      type Address string
      """
    When the linter runs on the package
    Then a "no-constructor" diagnostic is reported for "Address"

  Scenario: custom type with valid exported constructor is not flagged
    Given a Go package with:
      """
      type Address string

      func NewAddress(val string) (Address, error) {
          return Address(val), nil
      }
      """
    When the linter runs on the package
    Then no "no-constructor" diagnostic is reported for "Address"

  Scenario: custom type with valid unexported constructor is not flagged
    Given a Go package with:
      """
      type address string

      func newAddress(val string) (address, error) {
          return address(val), nil
      }
      """
    When the linter runs on the package
    Then no "no-constructor" diagnostic is reported for "address"

  Scenario: constructor missing error return is not recognized
    Given a Go package with:
      """
      type Address string

      func NewAddress(val string) Address {
          return Address(val)
      }
      """
    When the linter runs on the package
    Then a "no-constructor" diagnostic is reported for "Address"

  Scenario: constructor with extra return values is not recognized
    Given a Go package with:
      """
      type Address string

      func NewAddress(val string) (Address, string, error) {
          return Address(val), "", nil
      }
      """
    When the linter runs on the package
    Then a "no-constructor" diagnostic is reported for "Address"

  Scenario: constructor with wrong name prefix is not recognized
    Given a Go package with:
      """
      type Address string

      func CreateAddress(val string) (Address, error) {
          return Address(val), nil
      }
      """
    When the linter runs on the package
    Then a "no-constructor" diagnostic is reported for "Address"

  Scenario: constructor in a different package is not recognized
    Given package "domain" with:
      """
      type Address string
      """
    And package "factory" with:
      """
      func NewAddress(val string) (domain.Address, error) {
          return domain.Address(val), nil
      }
      """
    When the linter runs on package "domain"
    Then a "no-constructor" diagnostic is reported for "Address"

  Scenario: custom type derived from another custom type requires its own constructor
    Given a Go package with:
      """
      type Address string

      func NewAddress(val string) (Address, error) {
          return Address(val), nil
      }

      type ShippingAddress Address
      """
    When the linter runs on the package
    Then a "no-constructor" diagnostic is reported for "ShippingAddress"
    And no "no-constructor" diagnostic is reported for "Address"
