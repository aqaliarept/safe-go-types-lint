package no_cast_outside_constructor

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

func example() {
	a := Address("foo") // want `no-cast`
	_ = a
}
