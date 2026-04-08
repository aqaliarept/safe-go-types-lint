package domain2

// Not excluded — diagnostics expected.

type Product struct {
	Name  string // want `no-scalar`
	Price int    // want `no-scalar`
}
