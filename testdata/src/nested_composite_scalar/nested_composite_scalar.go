package nested_composite_scalar

type Index struct {
	Lookup map[string][]int // want `no-scalar`
}
