package no_cast_inside_constructor_not_flagged

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }
