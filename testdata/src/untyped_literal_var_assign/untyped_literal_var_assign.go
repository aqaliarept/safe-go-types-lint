package untyped_literal_var_assign

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

var a Address = "foo" // want `untyped-literal`
