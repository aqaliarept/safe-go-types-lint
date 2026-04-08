package no_constructor_unexported

type address string

func newAddress(val string) (address, error) {
	return address(val), nil
}
