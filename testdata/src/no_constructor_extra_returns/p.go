package no_constructor_extra_returns

type Address string // want `no-constructor`

func NewAddress(val string) (Address, string, error) {
	return Address(val), "", nil
}
