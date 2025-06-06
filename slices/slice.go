package slices

// AnySlice takes a slice and returns the same slice, but with each value
// represented as an any (or interface{})
func AnySlice[T any](src []T) []any {
	if src == nil {
		return nil
	}

	dest := make([]any, len(src))
	for i := range src {
		dest[i] = src[i]
	}

	return dest
}

// FromAnySlice takes an []any and returns the same slice with each item coerced
// to a T. If any item in the slice cannot be asserted as a T, FromAnySlice panics
func FromAnySlice[T any](src []any) []T {
	if src == nil {
		return nil
	}

	dest := make([]T, len(src))
	for i := range src {
		dest[i] = src[i].(T)
	}

	return dest
}

// SubsliceUntil will return a new slice with the first n elements represented in src, where
// filter(src[n]) == true. Averages O(lg n) time.
func SubsliceUntil[T any](src []T, filter func(item T) bool) []T {
	dst := make([]T, 0, len(src))
	for _, item := range src {
		if filter(item) {
			return dst
		}

		dst = append(dst, item)
	}

	return dst
}

// Map represents the "map" in the Map-filter-reduce pattern. That is,
// given a slice src and a function mapper, return a slice dst where
// dst[i] == mapper(src[i])
func Map[R any, T any](src []T, mapper func(T) R) []R {
	rs := make([]R, 0, len(src))
	for _, t := range src {
		rs = append(rs, mapper(t))
	}

	return rs
}

// Filter represents the filter in the Map-filter-reduce pattern. Given a slice
// src, return a slice dst where each element is in the order it appeared in src
// and where filter(src[i]) == true
func Filter[T any](src []T, filterer func(T) bool) []T {
	res := make([]T, 0, len(src))
	for _, t := range src {
		if filterer(t) {
			res = append(res, t)
		}
	}
	return res
}

// Reduce represents the reduce in the Map-filter-reduce pattern, where every
// item is passed to a function and given an initial value, each item affects
// the current value and the final value is returned
func Reduce[R any, T any](src []T, reducer func(R, T) R, initialValue R) R {
	value := initialValue
	for _, t := range src {
		value = reducer(value, t)
	}

	return value
}

// Uniq takes a slice and removes duplicate elements. Unlike unix's uniq
// command, it will work if the slice is not sorted. The returned slice's
// values will have duplicates removed, but will otherwise be in the same
// order
func Uniq[T comparable](src []T) []T {
	dest := make([]T, 0, len(src))
	set := make(map[T]bool)

	for _, t := range src {
		if len(dest) == 0 {
			dest = append(dest, t)
			set[t] = true
			continue
		}

		if _, exists := set[t]; !exists {
			dest = append(dest, t)
			set[t] = true
		}
	}

	return dest
}

// UniqFunc is similar to Uniq, but can accept slices of any type, and requires a
// func "id" that satisfies the following conditions:
//
//   - for a given T t, id(t) returns a value that satisfies the comparable interface
//   - if id(t1) == id(t2), then t1 == t2
func UniqFunc[T any, R comparable](src []T, id func(t T) R) []T {
	dest := make([]T, 0, len(src))
	set := make(map[R]bool)

	for _, t := range src {
		hash := id(t)
		if len(dest) == 0 {
			dest = append(dest, t)
			set[hash] = true
			continue
		}

		if _, exists := set[hash]; !exists {
			dest = append(dest, t)
			set[hash] = true
		}
	}

	return dest
}
