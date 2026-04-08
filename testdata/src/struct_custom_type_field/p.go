package struct_custom_type_field

type Name string // want `no-constructor`

type User struct {
	Name Name
}
