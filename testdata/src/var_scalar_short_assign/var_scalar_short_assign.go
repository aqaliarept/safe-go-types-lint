package var_scalar_short_assign

func example() {
	result := "pending" // want `no-scalar`
	_ = result
}
