package no_cast_reverse_conversion_not_flagged

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

func example(a Address) string {
	return string(a)
}
