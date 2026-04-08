package slice_custom_type_local_var

type Tag string // want `no-constructor`

func example() {
	var tags []Tag // no diagnostic for no-scalar
	_ = tags
}
