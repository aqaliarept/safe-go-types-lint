package no_constructor_missing_error

type Address string // want `no-constructor`

func NewAddress(val string) Address {
	return Address(val)
}
