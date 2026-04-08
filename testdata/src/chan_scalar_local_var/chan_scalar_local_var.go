package chan_scalar_local_var

func example() {
	ch := make(chan string) // want `no-scalar`
	_ = ch
}
