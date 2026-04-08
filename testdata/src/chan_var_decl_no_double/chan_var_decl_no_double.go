package chan_var_decl_no_double

// Checks that "var ch chan int" produces exactly one no-scalar diagnostic,
// not two (one from the GenDecl handler and one from the ChanType walker).

func example() {
	var ch chan int // want `no-scalar`
	_ = ch
}
