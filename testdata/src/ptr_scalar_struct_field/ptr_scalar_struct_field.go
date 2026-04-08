package ptr_scalar_struct_field

type Record struct {
	Name *string // want `no-scalar`
}
