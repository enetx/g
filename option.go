package g

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Some creates an Option containing a non-nil value.
func Some[T any](value T) Option[T] { return Option[T]{value: value, isSome: true} }

// None creates an Option containing a nil value.
func None[T any]() Option[T] { return Option[T]{isSome: false} }

// OptionOf creates an Option[T] based on the provided value and status flag.
// If ok is true, it returns an Option containing the value.
// Otherwise, it returns an Option representing no value.
func OptionOf[T any](value T, ok bool) Option[T] {
	if ok {
		return Some(value)
	}

	return None[T]()
}

// TransformOption applies the given function to the value inside the Option, producing a new Option with the transformed value.
// If the input Option is None, the output Option will also be None.
// Parameters:
//   - o: The input Option to map over.
//   - fn: The function to apply to the value inside the Option.
//
// Returns:
//
//	A new Option with the transformed value.
func TransformOption[T, U any](o Option[T], fn func(T) Option[U]) Option[U] {
	if o.IsNone() {
		return None[U]()
	}

	return fn(o.Some())
}

// Some returns the value held in the Option.
func (o Option[T]) Some() T { return o.value }

// IsSome returns true if the Option contains a non-nil value.
func (o Option[T]) IsSome() bool { return o.isSome }

// IsNone returns true if the Option contains a nil value.
func (o Option[T]) IsNone() bool { return !o.isSome }

// Unwrap returns the value held in the Option. If the Option contains a nil value, it panics.
func (o Option[T]) Unwrap() T {
	if o.isSome {
		return o.Some()
	}

	err := errors.New("can't unwrap none value")
	if pc, file, line, ok := runtime.Caller(1); ok {
		out := fmt.Sprintf("[%s:%d] [%s] %v", filepath.Base(file), line, runtime.FuncForPC(pc).Name(), err)
		fmt.Fprintln(os.Stderr, out)
	}

	panic(err)
}

// UnwrapOr returns the value held in the Option. If the Option contains a nil value, it returns the provided default value.
func (o Option[T]) UnwrapOr(value T) T {
	if o.isSome {
		return o.Some()
	}

	return value
}

// Expect returns the value held in the Option. If the Option contains a nil value, it panics with the provided message.
func (o Option[T]) Expect(msg string) T {
	if o.isSome {
		return o.Some()
	}

	panic(msg)
}

// Then applies the function fn to the value inside the Option and returns a new Option.
// If the Option is None, it returns the same Option without applying fn.
func (o Option[T]) Then(fn func(T) Option[T]) Option[T] {
	if o.isSome {
		return fn(o.Some())
	}

	return o
}

// String returns a string representation of the Option.
// If the Option contains a value, it returns a string in the format "Some(value)".
// Otherwise, it returns "None".
func (o Option[T]) String() string {
	if o.isSome {
		return fmt.Sprintf("Some(%v)", o.Some())
	}

	return "None"
}
