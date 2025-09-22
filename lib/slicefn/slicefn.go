package slicefn

import (
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

const ErrNotFound = merr.Code("not_found")

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

	return res, err
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

func FilterE[T any](slice []T, fn func(T) (bool, error)) ([]T, error) {
	res := make([]T, 0, len(slice))

	for _, v := range slice {
		if ok, err := fn(v); err != nil {
			return nil, err
		} else if ok {
			res = append(res, v)
		}
	}

	return res, nil
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
			return acc, err
		}
	}

	return acc, err
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

func FindIndex[T any](slice []T, fn func(T) bool) int {
	for i, v := range slice {
		if fn(v) {
			return i
		}
	}

	return -1
}

func FindIndexE[T any](slice []T, fn func(T) (bool, error)) (int, error) {
	for i, v := range slice {
		if ok, err := fn(v); err != nil {
			return -1, err
		} else if ok {
			return i, nil
		}
	}

	return -1, nil
}

func FindLastIndex[T any](slice []T, fn func(T) bool) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if fn(slice[i]) {
			return i
		}
	}

	return -1
}

func FindLastIndexE[T any](slice []T, fn func(T) (bool, error)) (int, error) {
	for i := len(slice) - 1; i >= 0; i-- {
		if ok, err := fn(slice[i]); err != nil {
			return -1, err
		} else if ok {
			return i, nil
		}
	}

	return -1, nil
}

func Find[T any](slice []T, fn func(T) bool) (res T, ok bool) {
	if i := FindIndex(slice, fn); i >= 0 {
		return slice[i], true
	}

	return res, ok
}

func FindE[T any](slice []T, fn func(T) (bool, error)) (res T, ok bool, err error) {
	i, err := FindIndexE(slice, fn)
	if err != nil {
		return res, ok, err
	} else if i >= 0 {
		return slice[i], true, nil
	}

	return res, ok, err
}

func FindPtr[T any](slice []T, fn func(T) bool) *T {
	if i := FindIndex(slice, fn); i >= 0 {
		return &slice[i]
	}

	return nil
}

func FindPtrE[T any](slice []T, fn func(T) (bool, error)) (*T, error) {
	i, err := FindIndexE(slice, fn)
	if err != nil {
		return nil, err
	} else if i >= 0 {
		return &slice[i], nil
	}

	//nolint:gocritic // we don't have a context here
	return nil, ErrNotFound
}

func FindLast[T any](slice []T, fn func(T) bool) (res T, ok bool) {
	if i := FindLastIndex(slice, fn); i >= 0 {
		return slice[i], true
	}

	return res, ok
}

func FindLastE[T any](slice []T, fn func(T) (bool, error)) (res T, ok bool, err error) {
	i, err := FindLastIndexE(slice, fn)
	if err != nil {
		return res, ok, err
	} else if i >= 0 {
		return slice[i], true, nil
	}

	return res, ok, err
}

func FindLastPtr[T any](slice []T, fn func(T) bool) *T {
	if i := FindLastIndex(slice, fn); i >= 0 {
		return &slice[i]
	}

	return nil
}

func FindLastPtrE[T any](slice []T, fn func(T) (bool, error)) (*T, error) {
	i, err := FindLastIndexE(slice, fn)
	if err != nil {
		return nil, err
	} else if i >= 0 {
		return &slice[i], nil
	}

	//nolint:gocritic // we don't have a context here
	return nil, ErrNotFound
}
