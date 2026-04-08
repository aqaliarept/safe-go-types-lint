package slice_scalar_local_var

func example() {
	var tags []string // want `no-scalar`
	_ = tags
}
