package no_constructor_wrong_name

type Address string // want `no-constructor`

func CreateAddress(val string) (Address, error) {
	return Address(val), nil
}
