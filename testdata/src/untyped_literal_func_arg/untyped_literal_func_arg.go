package untyped_literal_func_arg

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

func consume(a Address) {}

func example() {
	consume("foo") // want `untyped-literal`
}
