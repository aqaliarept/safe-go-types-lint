package nolint_other

type User struct {
	Name string //nolint:some-other-linter // want `no-scalar`
}
