package no_zero_value_bare_var

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

var a Address // want `no-zero-value`
