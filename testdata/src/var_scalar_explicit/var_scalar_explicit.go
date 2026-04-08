package var_scalar_explicit

func example() {
	var x string // want `no-scalar`
	_ = x
}
