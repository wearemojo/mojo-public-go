package miter

import (
	"iter"
)

func CollectErr[T any](seq iter.Seq2[T, error]) ([]T, error) {
	res := []T{}
	for item, err := range seq {
		if err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	return res, nil
}

func CollectErrUnptr[T any](seq iter.Seq2[*T, error]) ([]T, error) {
	res := []T{}
	for item, err := range seq {
		if err != nil {
			return nil, err
		}

		res = append(res, *item)
	}

	return res, nil
}
