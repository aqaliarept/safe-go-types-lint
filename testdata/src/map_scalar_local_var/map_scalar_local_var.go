package map_scalar_local_var

func example() {
	counts := map[string]int{} // want `no-scalar`
	_ = counts
}
