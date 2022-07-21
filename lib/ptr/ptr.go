package ptr

func P[T any](v T) *T {
	return &v
}

func POrNil[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	} else {
		return &v
	}
}
