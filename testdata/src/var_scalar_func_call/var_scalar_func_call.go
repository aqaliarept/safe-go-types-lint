package var_scalar_func_call

func getValue() string { return "" }

func example() {
	x := getValue() // no diagnostic
	_ = x
}
