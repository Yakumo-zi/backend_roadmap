package slice

import (
	"cmp"
)

func Filter[T any](arr []T, f func(T) bool) (res []T) {
	for _, v := range arr {
		if f(v) {
			res = append(res, v)
		}
	}
	return
}

func Map[T any](arr []T, f func(T) T) []T {
	res := make([]T, len(arr))
	for i, v := range arr {
		res[i] = f(v)
	}
	return res
}

func Range[T any](arr []T, f func(any, T)) {
	for i, v := range arr {
		f(i, v)
	}
}

func Reduce[T any, V any](arr []T, f func(V, T) V, init V) V {
	res := init
	for _, v := range arr {
		res = f(res, v)
	}
	return res
}

func Distinct[T cmp.Ordered](arr []T) []T {
	set := map[T]int{}
	for i, v := range arr {
		set[v] = i
	}
	res := make([]T, len(arr))
	for k, v := range set {
		res[v] = k
	}
	return res
}
