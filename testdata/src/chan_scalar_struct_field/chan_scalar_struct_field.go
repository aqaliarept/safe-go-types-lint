package chan_scalar_struct_field

type Worker struct {
	Jobs chan int // want `no-scalar`
}
