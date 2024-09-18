package util

// NewSlice creates a new slice from the given slices.
func NewSlice[V any](sn ...[]V) []V {
	var res []V

	if len(sn) == 0 {
		return res
	}

	for _, s := range sn {
		res = append(res, s...)
	}
	return res
}

func IfValue[V comparable](fb V, sn ...V) V {
	if len(sn) == 0 {
		return fb
	}

	var zero V
	for _, s := range sn {
		if s != zero {
			return s
		}
	}

	return fb
}
