package no_constructor_derived

type Address string

func NewAddress(val string) (Address, error) {
	return Address(val), nil
}

type ShippingAddress Address // want `no-constructor`
