package same_package_constant_not_flagged

type Address string

func NewAddress(val string) (Address, error) { return Address(val), nil }

const EmptyAddress = Address("")
