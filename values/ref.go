package values

import "reflect"

// Ref returns a pointer to val. Useful for assigning pointers of literal values.
// For example, &"foo" doesn't work, but utils.Ref("foo") does
func Ref[T any](val T) *T {
	return &val
}

// NilRefIfZero will return a pointer to val, IF val is not the zero value for
// that type. If it is zero, it returns nil
func NilRefIfZero[T any](val T) *T {
	if reflect.ValueOf(val).IsZero() {
		return nil
	}

	return &val
}
