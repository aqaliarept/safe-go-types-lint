package no_zero_value_struct_field

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

type Order struct{ Destination Address }

func NewOrder(d Address) (Order, error) { return Order{Destination: d}, nil }

var o Order // want `no-zero-value`
