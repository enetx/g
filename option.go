package g

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Some creates an Option containing a non-nil value.
func Some[T any](value T) Option[T] { return Option[T]{&value} }

// None creates an Option containing a nil value.
func None[T any]() Option[T] { return Option[T]{nil} }

// Some returns the value held in the Option.
func (o Option[T]) Some() T { return *o.value }

// IsSome returns true if the Option contains a non-nil value.
func (o Option[T]) IsSome() bool { return o.value != nil }

// IsNone returns true if the Option contains a nil value.
func (o Option[T]) IsNone() bool { return o.value == nil }

// Unwrap returns the value held in the Option. If the Option contains a nil value, it panics.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		err := errors.New("can't unwrap none value")
		if pc, file, line, ok := runtime.Caller(1); ok {
			out := fmt.Sprintf("[%s:%d] [%s] %v", filepath.Base(file), line, runtime.FuncForPC(pc).Name(), err)
			fmt.Fprintln(os.Stderr, out)
		}

		panic(err)
	}

	return o.Some()
}

// UnwrapOr returns the value held in the Option. If the Option contains a nil value, it returns the provided default value.
func (o Option[T]) UnwrapOr(value T) T {
	if o.IsNone() {
		return value
	}

	return o.Some()
}

// UnwrapOrDefault returns the value held in the Option. If the Option contains a value,
// it returns the value. If the Option is None, it returns the default value for type T.
func (o Option[T]) UnwrapOrDefault() T {
	if o.IsNone() {
		return *new(T)
	}

	return o.Some()
}

// Expect returns the value held in the Option. If the Option contains a nil value, it panics with the provided message.
func (o Option[T]) Expect(msg string) T {
	if o.IsNone() {
		panic(msg)
	}

	return o.Some()
}
