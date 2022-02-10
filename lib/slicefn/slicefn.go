package slicefn

func Map[T1, T2 any](slice []T1, fn func(T1) T2) []T2 {
	res := make([]T2, len(slice))

	for i, v := range slice {
		res[i] = fn(v)
	}

	return res
}

func MapE[T1, T2 any](slice []T1, fn func(T1) (T2, error)) (res []T2, err error) {
	res = make([]T2, len(slice))

	for i, v := range slice {
		res[i], err = fn(v)
		if err != nil {
			return nil, err
		}
	}

	return
}

func Filter[T any](slice []T, fn func(T) bool) []T {
	res := make([]T, 0, len(slice))

	for _, v := range slice {
		if fn(v) {
			res = append(res, v)
		}
	}

	return res
}

func FilterE[T any](slice []T, fn func(T) (bool, error)) (res []T, err error) {
	res = make([]T, 0, len(slice))

	for _, v := range slice {
		if ok, err := fn(v); err != nil {
			return nil, err
		} else if ok {
			res = append(res, v)
		}
	}

	return
}

func Reduce[T1, T2 any](slice []T1, fn func(acc T2, item T1) T2, initial T2) T2 {
	acc := initial

	for _, v := range slice {
		acc = fn(acc, v)
	}

	return acc
}

func ReduceE[T1, T2 any](slice []T1, fn func(acc T2, item T1) (T2, error), initial T2) (acc T2, err error) {
	acc = initial

	for _, v := range slice {
		acc, err = fn(acc, v)
		if err != nil {
			return
		}
	}

	return
}

func Some[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if fn(v) {
			return true
		}
	}

	return false
}

func SomeE[T any](slice []T, fn func(T) (bool, error)) (bool, error) {
	for _, v := range slice {
		if ok, err := fn(v); err != nil {
			return false, err
		} else if ok {
			return true, nil
		}
	}

	return false, nil
}

func Every[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if !fn(v) {
			return false
		}
	}

	return true
}

func EveryE[T any](slice []T, fn func(T) (bool, error)) (bool, error) {
	for _, v := range slice {
		if ok, err := fn(v); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		}
	}

	return true, nil
}

func Find[T any](slice []T, fn func(T) bool) (res T, ok bool) {
	for _, v := range slice {
		if fn(v) {
			return v, true
		}
	}

	return
}

func FindE[T any](slice []T, fn func(T) (bool, error)) (res T, ok bool, err error) {
	for _, v := range slice {
		if ok, err := fn(v); err != nil {
			return zeroValue[T](), false, err
		} else if ok {
			return v, true, nil
		}
	}

	return
}
