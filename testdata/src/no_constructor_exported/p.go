package no_constructor_exported

type Address string

func NewAddress(val string) (Address, error) {
	return Address(val), nil
}
