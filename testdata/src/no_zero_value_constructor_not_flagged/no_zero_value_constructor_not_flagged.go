package no_zero_value_constructor_not_flagged

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

func example() {
	a, err := NewAddress("x")
	_ = a
	_ = err
}
