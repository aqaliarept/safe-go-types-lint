package same_package_constant_used_not_flagged

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

const EmptyAddress = Address("")

func example() {
	a := EmptyAddress
	_ = a
}
