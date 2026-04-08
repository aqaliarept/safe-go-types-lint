package slice_custom_type_not_flagged

type Tag string // want `no-constructor`

type Order struct {
	Tags []Tag // no diagnostic for no-scalar
}
