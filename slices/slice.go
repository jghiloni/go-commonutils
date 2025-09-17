package slices

import (
	"cmp"
	"crypto/rand"
	"fmt"
	"iter"
	"math"
	"math/big"
	"slices"
)

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

// Shuffle will randomize the elements in a slice and return an iterator to the next item
func Shuffle[S ~[]E, E any](list S) S {
	used := make(map[int64]bool, len(list))
	ub := big.NewInt(int64(len(list)))

	idxAttempts := max(min(math.MaxInt, len(list)*5), 5000)
	next := func() int {

		for range idxAttempts {
			i, _ := rand.Int(rand.Reader, ub)
			if _, found := used[i.Int64()]; found {
				continue
			}

			used[i.Int64()] = true
			return int(i.Int64())
		}

		return -1
	}

	retVal := make(S, len(list))
	for i := range list {
		idx := next()
		if idx == -1 {
			panic(fmt.Errorf("could not find an unused index after %d attempts", idxAttempts))
		}

		retVal[i] = list[idx]
	}

	return retVal
}

// EVERYTHING BELOW THIS IS A DIRECT CALL TO THE STDLIB SLICES PACKAGE SO THAT OUR SLICES PACKAGE CAN BE A DROP IN
// REPLACEMENT

func All[Slice ~[]E, E any](s Slice) iter.Seq2[int, E] {
	return slices.All(s)
}
func AppendSeq[Slice ~[]E, E any](s Slice, seq iter.Seq[E]) Slice {
	return slices.AppendSeq(s, seq)
}
func Backward[Slice ~[]E, E any](s Slice) iter.Seq2[int, E] {
	return slices.Backward(s)
}
func BinarySearch[S ~[]E, E cmp.Ordered](x S, target E) (int, bool) {
	return slices.BinarySearch(x, target)
}
func BinarySearchFunc[S ~[]E, E, T any](x S, target T, cmp func(E, T) int) (int, bool) {
	return slices.BinarySearchFunc(x, target, cmp)
}
func Chunk[Slice ~[]E, E any](s Slice, n int) iter.Seq[Slice] {
	return slices.Chunk(s, n)
}
func Clip[S ~[]E, E any](s S) S {
	return slices.Clip(s)
}
func Clone[S ~[]E, E any](s S) S {
	return slices.Clone(s)
}
func Collect[E any](seq iter.Seq[E]) []E {
	return slices.Collect(seq)
}
func Compact[S ~[]E, E comparable](s S) S {
	return slices.Compact(s)
}
func CompactFunc[S ~[]E, E any](s S, eq func(E, E) bool) S {
	return slices.CompactFunc(s, eq)
}
func Compare[S ~[]E, E cmp.Ordered](s1, s2 S) int {
	return slices.Compare(s1, s2)
}
func CompareFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](s1 S1, s2 S2, cmp func(E1, E2) int) int {
	return slices.CompareFunc(s1, s2, cmp)
}
func Concat[S ~[]E, E any](s ...S) S {
	return slices.Concat(s...)
}
func Contains[S ~[]E, E comparable](s S, v E) bool {
	return slices.Contains(s, v)
}
func ContainsFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	return slices.ContainsFunc(s, f)
}
func Delete[S ~[]E, E any](s S, i, j int) S {
	return slices.Delete(s, i, j)
}
func DeleteFunc[S ~[]E, E any](s S, del func(E) bool) S {
	return slices.DeleteFunc(s, del)
}
func Equal[S ~[]E, E comparable](s1, s2 S) bool {
	return slices.Equal(s1, s2)
}
func EqualFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](s1 S1, s2 S2, eq func(E1, E2) bool) bool {
	return slices.EqualFunc(s1, s2, eq)
}
func Grow[S ~[]E, E any](s S, n int) S {
	return slices.Grow(s, n)
}
func Index[S ~[]E, E comparable](s S, v E) int {
	return slices.Index(s, v)
}
func IndexFunc[S ~[]E, E any](s S, f func(E) bool) int {
	return slices.IndexFunc(s, f)
}
func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	return slices.Insert(s, i, v...)
}
func IsSorted[S ~[]E, E cmp.Ordered](x S) bool {
	return slices.IsSorted(x)
}
func IsSortedFunc[S ~[]E, E any](x S, cmp func(a, b E) int) bool {
	return slices.IsSortedFunc(x, cmp)
}
func Max[S ~[]E, E cmp.Ordered](x S) E {
	return slices.Max(x)
}
func MaxFunc[S ~[]E, E any](x S, cmp func(a, b E) int) E {
	return slices.MaxFunc(x, cmp)
}
func Min[S ~[]E, E cmp.Ordered](x S) E {
	return slices.Min(x)
}
func MinFunc[S ~[]E, E any](x S, cmp func(a, b E) int) E {
	return slices.MinFunc(x, cmp)
}
func Repeat[S ~[]E, E any](x S, count int) S {
	return slices.Repeat(x, count)
}
func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	return slices.Replace(s, i, j, v...)
}
func Reverse[S ~[]E, E any](s S) {
	slices.Reverse(s)
}
func Sort[S ~[]E, E cmp.Ordered](x S) {
	slices.Sort(x)
}
func SortFunc[S ~[]E, E any](x S, cmp func(a, b E) int) {
	slices.SortFunc(x, cmp)
}
func SortStableFunc[S ~[]E, E any](x S, cmp func(a, b E) int) {
	slices.SortStableFunc(x, cmp)
}
func Sorted[E cmp.Ordered](seq iter.Seq[E]) []E {
	return slices.Sorted(seq)
}
func SortedFunc[E any](seq iter.Seq[E], cmp func(E, E) int) []E {
	return slices.SortedFunc(seq, cmp)
}
func SortedStableFunc[E any](seq iter.Seq[E], cmp func(E, E) int) []E {
	return slices.SortedStableFunc(seq, cmp)
}
func Values[Slice ~[]E, E any](s Slice) iter.Seq[E] {
	return slices.Values(s)
}
