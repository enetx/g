package f

import (
	"cmp"
	"reflect"
	"regexp"
	"strings"

	"github.com/enetx/g/constraints"
)

// IsComparable reports whether the value v is comparable.
func IsComparable[T any](t T) bool { return reflect.ValueOf(t).Comparable() }

// IsZero is a generic function designed to check if a value is considered zero.
func IsZero[T cmp.Ordered](v T) bool { return v == *new(T) }

// IsEven is a generic function that checks if the provided integer is even.
func IsEven[T constraints.Integer](i T) bool { return i%2 == 0 }

// IsOdd is a generic function that checks if the provided integer is odd.
func IsOdd[T constraints.Integer](i T) bool { return i%2 != 0 }

// Match returns a function that checks whether a string or []byte matches a given regular expression.
func Match[T ~string | ~[]byte](t *regexp.Regexp) func(T) bool {
	return func(s T) bool {
		return t.MatchString(string(s))
	}
}

// Contains returns a function that checks whether a string or []byte contains a given substring.
func Contains[T ~string | ~[]byte](t T) func(T) bool {
	return func(s T) bool {
		return strings.Contains(string(s), string(t))
	}
}

// ContainsAnyChars returns a function that checks whether a string contains any of the characters from a given set.
func ContainsAnyChars[T ~string | ~[]byte](t T) func(T) bool {
	return func(s T) bool {
		return strings.ContainsAny(string(s), string(t))
	}
}

// StartsWith returns a function that checks whether a string starts with a given prefix.
func StartsWith[T ~string | ~[]byte](t T) func(T) bool {
	return func(s T) bool {
		return strings.HasPrefix(string(s), string(t))
	}
}

// EndsWith returns a function that checks whether a string ends with a given suffix.
func EndsWith[T ~string | ~[]byte](t T) func(T) bool {
	return func(s T) bool {
		return strings.HasSuffix(string(s), string(t))
	}
}

// Eq returns a comparison function that evaluates to true when a value is equal to the provided threshold.
func Eq[T comparable](t T) func(T) bool {
	return func(s T) bool {
		return s == t
	}
}

// Ne returns a comparison function that evaluates to true when a value is not equal to the provided threshold.
func Ne[T comparable](t T) func(T) bool {
	return func(s T) bool {
		return s != t
	}
}

// Eqd returns a comparison function that evaluates to true when a value is deeply equal to the provided threshold.
func Eqd[T any](t T) func(T) bool {
	return func(s T) bool {
		return reflect.DeepEqual(t, s)
	}
}

// Ned returns a comparison function that evaluates to true when a value is not deeply equal to the provided threshold.
func Ned[T any](t T) func(T) bool {
	return func(s T) bool {
		return !reflect.DeepEqual(t, s)
	}
}

// Eqi returns a fast equality predicate for non-comparable types.
func Eqi[T any](t T) func(T) bool {
	vt := reflect.ValueOf(t)
	if !vt.IsValid() {
		return func(s T) bool { return !reflect.ValueOf(s).IsValid() }
	}

	typ := vt.Type()

	if typ.Comparable() {
		tv := vt.Interface()
		return func(s T) bool {
			vs := reflect.ValueOf(s)
			return vs.IsValid() &&
				vs.Type() == typ &&
				tv == vs.Interface()
		}
	}

	switch vt.Kind() {
	case reflect.Func:
		if vt.IsNil() {
			return func(s T) bool {
				vs := reflect.ValueOf(s)
				return vs.IsValid() &&
					vs.Kind() == reflect.Func &&
					vs.Type() == typ &&
					vs.IsNil()
			}
		}

		fp := vt.Pointer()

		return func(s T) bool {
			vs := reflect.ValueOf(s)
			return vs.IsValid() &&
				vs.Kind() == reflect.Func &&
				vs.Type() == typ &&
				!vs.IsNil() &&
				vs.Pointer() == fp
		}
	case reflect.Slice:
		if vt.IsNil() {
			return func(s T) bool {
				vs := reflect.ValueOf(s)
				return vs.IsValid() &&
					vs.Kind() == reflect.Slice &&
					vs.Type() == typ &&
					vs.IsNil()
			}
		}

		base := vt.UnsafePointer()
		l := vt.Len()

		return func(s T) bool {
			vs := reflect.ValueOf(s)
			return vs.IsValid() &&
				vs.Kind() == reflect.Slice &&
				vs.Type() == typ &&
				!vs.IsNil() &&
				vs.UnsafePointer() == base &&
				vs.Len() == l
		}
	case reflect.Map:
		if vt.IsNil() {
			return func(s T) bool {
				vs := reflect.ValueOf(s)
				return vs.IsValid() &&
					vs.Kind() == reflect.Map &&
					vs.Type() == typ &&
					vs.IsNil()
			}
		}

		mp := vt.Pointer()

		return func(s T) bool {
			vs := reflect.ValueOf(s)
			return vs.IsValid() &&
				vs.Kind() == reflect.Map &&
				vs.Type() == typ &&
				!vs.IsNil() &&
				vs.Pointer() == mp
		}
	}

	return func(s T) bool { return reflect.DeepEqual(t, s) }
}

// Nei returns a comparison function that evaluates to true when a value is not identity equal to the provided threshold.
func Nei[T any](t T) func(T) bool {
	eq := Eqi(t)
	return func(s T) bool {
		return !eq(s)
	}
}

// Gt returns a comparison function that evaluates to true when a value is greater than the threshold.
func Gt[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s > t
	}
}

// Gte returns a comparison function that evaluates to true when a value is greater than or equal to the threshold.
func Gte[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s >= t
	}
}

// Lt returns a comparison function that evaluates to true when a value is less than the threshold.
func Lt[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s < t
	}
}

// Lte returns a comparison function that evaluates to true when a value is less than or equal to the threshold.
func Lte[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s <= t
	}
}
