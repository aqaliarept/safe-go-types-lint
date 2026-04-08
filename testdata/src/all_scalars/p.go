package all_scalars

type AllScalars struct {
	F1  string     // want `no-scalar`
	F2  bool       // want `no-scalar`
	F3  int        // want `no-scalar`
	F4  int8       // want `no-scalar`
	F5  int16      // want `no-scalar`
	F6  int32      // want `no-scalar`
	F7  int64      // want `no-scalar`
	F8  uint       // want `no-scalar`
	F9  uint8      // want `no-scalar`
	F10 uint16     // want `no-scalar`
	F11 uint32     // want `no-scalar`
	F12 uint64     // want `no-scalar`
	F13 float32    // want `no-scalar`
	F14 float64    // want `no-scalar`
	F15 complex64  // want `no-scalar`
	F16 complex128 // want `no-scalar`
	F17 byte       // want `no-scalar`
	F18 rune       // want `no-scalar`
}
