package nolint_specific

type User struct {
	Name string //nolint:safe-go-types/no-scalar
	Age  int    // want `no-scalar`
}
