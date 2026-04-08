package domain_check

// This file is NOT in a legacy package — diagnostics should still be emitted.

type User struct {
	Name string // want `no-scalar`
	Age  int    // want `no-scalar`
}
