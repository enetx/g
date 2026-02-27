package g

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Option is a generic struct for representing an optional value.
type Option[T any] struct {
	v      T    // Value.
	isSome bool // Indicator of value presence.
}

// Unit represents an empty value.
// Used in contexts where a function needs to return "something" but
// the actual value doesn't matter, only success/failure status.
type Unit struct{}

// Some creates an Option containing a value.
func Some[T any](value T) Option[T] { return Option[T]{v: value, isSome: true} }

// None creates an Option representing no value.
func None[T any]() Option[T] { return Option[T]{isSome: false} }

// OptionOf creates an Option[T] based on the provided value and a boolean flag.
// If ok is true, it returns Some(value).
// Otherwise, it returns None.
func OptionOf[T any](value T, ok bool) Option[T] {
	if ok {
		return Some(value)
	}

	return None[T]()
}

// OptionFromPtr converts a pointer into an Option.
// Returns None if ptr is nil.
func OptionFromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return None[T]()
	}

	return Some(*ptr)
}

// TransformOption applies the given function to the value inside the Option, producing a new Option with the transformed value.
// If the input Option is None, the output Option will also be None.
// Parameters:
//   - o: The input Option to map over.
//   - fn: The function that returns an Option to apply to the value inside the Option.
//
// Returns:
//
//	A new Option with the transformed value, or None if the input was None.
func TransformOption[T, U any](o Option[T], fn func(T) Option[U]) Option[U] {
	if o.isSome {
		return fn(o.v)
	}

	return None[U]()
}

// Some returns the contained value of the Option.
//
// WARNING: If the Option is None, this method will return the zero value
// for type T. Always check IsSome() before calling this method, or use safer alternatives
// like Unwrap(), or UnwrapOr().
func (o Option[T]) Some() T { return o.v }

// IsSome returns true if the Option contains a value.
func (o Option[T]) IsSome() bool { return o.isSome }

// IsNone returns true if the Option represents no value.
func (o Option[T]) IsNone() bool { return !o.isSome }

// Unwrap returns the value held in the Option. If the Option is None, it panics.
func (o Option[T]) Unwrap() T {
	if o.isSome {
		return o.v
	}

	const panicMsg = "called Option.Unwrap() on a None value"

	if pc, file, line, ok := runtime.Caller(1); ok {
		out := fmt.Sprintf("[%s:%d] [%s] %s", filepath.Base(file), line, runtime.FuncForPC(pc).Name(), panicMsg)
		fmt.Fprintln(os.Stderr, out)
	}

	panic(panicMsg)
}

// UnwrapOr returns the value held in the Option. If the Option is None, it returns the provided default value.
func (o Option[T]) UnwrapOr(value T) T {
	if o.isSome {
		return o.v
	}

	return value
}

// UnwrapOrDefault returns the contained value if Some; otherwise returns the zero value for T.
func (o Option[T]) UnwrapOrDefault() T {
	if o.isSome {
		return o.v
	}

	var zero T
	return zero
}

// Expect returns the value held in the Option. If the Option is None, it panics with the provided message.
func (o Option[T]) Expect(msg string) T {
	if o.isSome {
		return o.v
	}

	out := fmt.Sprintf("Expect() failed: %s", msg)
	fmt.Fprintln(os.Stderr, out)
	panic(out)
}

// Then applies the function fn to the value inside the Option and returns a new Option.
// If the Option is None, it returns the same Option without applying fn.
func (o Option[T]) Then(fn func(T) Option[T]) Option[T] {
	if o.isSome {
		return fn(o.v)
	}

	return o
}

// Filter returns Some(value) if the Option is Some and the predicate returns true.
// Otherwise, it returns None.
func (o Option[T]) Filter(pred func(T) bool) Option[T] {
	if o.isSome && pred(o.v) {
		return o
	}

	return None[T]()
}

// Or returns the Option if it contains a value.
// Otherwise, it returns the provided alternative Option.
func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.isSome {
		return o
	}

	return other
}

// OrElse returns the Option if it contains a value.
// Otherwise, it calls fn and returns its result.
func (o Option[T]) OrElse(fn func() Option[T]) Option[T] {
	if o.isSome {
		return o
	}

	return fn()
}

// IsSomeAnd returns true if the Option is Some
// and the predicate returns true for the contained value.
func (o Option[T]) IsSomeAnd(pred func(T) bool) bool {
	return o.isSome && pred(o.v)
}

// Insert inserts the given value into the Option,
// replacing any existing value, and returns a pointer
// to the inserted value.
func (o *Option[T]) Insert(value T) *T {
	o.v = value
	o.isSome = true

	return &o.v
}

// GetOrInsert inserts the given value if the Option is None,
// and returns a pointer to the contained value.
// If the Option already contains a value, it is left unchanged.
func (o *Option[T]) GetOrInsert(value T) *T {
	if !o.isSome {
		o.v = value
		o.isSome = true
	}
	return &o.v
}

// GetOrInsertWith inserts a value computed by fn if the Option is None,
// and returns a pointer to the contained value.
// The function fn is evaluated lazily.
func (o *Option[T]) GetOrInsertWith(fn func() T) *T {
	if !o.isSome {
		o.v = fn()
		o.isSome = true
	}

	return &o.v
}

// Take takes the value out of the Option, leaving None in its place.
// It returns Some(value) if the Option was Some,
// otherwise returns None.
func (o *Option[T]) Take() Option[T] {
	if !o.isSome {
		return None[T]()
	}

	val := o.v
	var zero T

	o.v = zero
	o.isSome = false

	return Some(val)
}

// Replace replaces the contained value with the given value,
// returning the old value as an Option.
// If the Option was None, it inserts the value and returns None.
func (o *Option[T]) Replace(value T) Option[T] {
	old := o.Take()
	o.v = value
	o.isSome = true

	return old
}

// OkOr converts the Option into a Result.
// Returns Ok(value) if Some, otherwise returns Err(err).
func (o Option[T]) OkOr(err error) Result[T] {
	if o.isSome {
		return Ok(o.v)
	}

	return Err[T](err)
}

// OkOrElse converts the Option into a Result.
// Returns Ok(value) if Some.
// Otherwise, it calls fn and returns Err(fn()).
func (o Option[T]) OkOrElse(fn func() error) Result[T] {
	if o.isSome {
		return Ok(o.v)
	}

	return Err[T](fn())
}

// Ptr returns a pointer to the contained value if Some.
// Otherwise, it returns nil.
func (o Option[T]) Ptr() *T {
	if o.isSome {
		return &o.v
	}

	return nil
}

// Result converts an Option into a Result.
// If the Option is Some, it returns an Ok Result with the value.
// If the Option is None, it returns an Err Result with the provided error.
func (o Option[T]) Result(err error) Result[T] {
	if o.isSome {
		return Ok(o.v)
	}

	return Err[T](err)
}

func (o Option[T]) Option() (T, bool) {
	if o.IsSome() {
		return o.Some(), true
	}

	var zero T
	return zero, false
}

// String returns a string representation of the Option.
// If the Option contains a value, it returns a string in the format "Some(value)".
// Otherwise, it returns "None".
func (o Option[T]) String() string {
	if o.isSome {
		return fmt.Sprintf("Some(%v)", o.v)
	}

	return "None"
}
