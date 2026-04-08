package ptr_scalar_local_var

func example() {
	var x *string // want `no-scalar`
	_ = x
}
