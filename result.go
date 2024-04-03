package g

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Ok returns a new Result[T] containing the given value.
func Ok[T any](value T) Result[T] { return Result[T]{value: &value, err: nil} }

// Err returns a new Result[T] containing the given error.
func Err[T any](err error) Result[T] { return Result[T]{value: nil, err: err} }

// ResultOf returns a new Result[T] based on the provided value and error.
// If err is not nil, it returns an Result containing the error.
// Otherwise, it returns an Result containing the value.
func ResultOf[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}

	return Ok(value)
}

// ResultMap applies the given function to the value inside the Result, producing a new Result with the transformed value.
// If the input Result contains a value, the provided function is applied to it.
// If the input Result contains an error, the output Result will also contain the same error.
// Parameters:
//   - r: The input Result to map over.
//   - fn: The function that returns a Result to apply to the value inside the Result.
//
// Returns:
//
//	A new Result with the transformed value, or the same error if the input Result contained an error.
func ResultMap[T, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if r.IsErr() {
		return Err[U](r.Err())
	}

	return fn(r.Ok())
}

// ResultOfMap applies the given function to the value inside the Result, producing a new Result with the transformed value.
// If the input Result contains a value, the provided function is applied to it.
// If the input Result contains an error, the output Result will also contain the same error.
// Parameters:
//   - r: The input Result to map over.
//   - fn: The function that returns a tuple (U, error) to apply to the value inside the Result.
//
// Returns:
//
//	A new Result with the transformed value, or the same error if the input Result contained an error.
func ResultOfMap[T, U any](r Result[T], fn func(T) (U, error)) Result[U] {
	if r.IsErr() {
		return Err[U](r.Err())
	}

	return ResultOf(fn(r.Ok()))
}

// Ok returns the value held in the Result.
func (r Result[T]) Ok() T { return *r.value }

// Err returns the error held in the Result.
func (r Result[T]) Err() error { return r.err }

// IsOk returns true if the Result contains a value (no error).
func (r Result[T]) IsOk() bool { return r.err == nil }

// IsErr returns true if the Result contains an error.
func (r Result[T]) IsErr() bool { return r.err != nil }

// Result returns the value held in the Result and its error.
func (r Result[T]) Result() (T, error) {
	if r.IsErr() {
		return *new(T), r.Err()
	}

	return r.Ok(), nil
}

// Unwrap returns the value held in the Result. If the Result contains an error, it panics.
func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		if pc, file, line, ok := runtime.Caller(1); ok {
			out := fmt.Sprintf("[%s:%d] [%s] %v", filepath.Base(file), line, runtime.FuncForPC(pc).Name(), r.err)
			fmt.Fprintln(os.Stderr, out)
		}

		panic(r.err)
	}

	return r.Ok()
}

// UnwrapOr returns the value held in the Result. If the Result contains an error, it returns the provided default value.
func (r Result[T]) UnwrapOr(value T) T {
	if r.IsErr() {
		return value
	}

	return r.Ok()
}

// UnwrapOrDefault returns the value held in the Result. If the Result contains an error,
// it returns the default value for type T. Otherwise, it returns the value held in the Result.
func (r Result[T]) UnwrapOrDefault() T {
	if r.IsErr() {
		return *new(T)
	}

	return r.Ok()
}

// Expect returns the value held in the Result. If the Result contains an error, it panics with the provided message.
func (r Result[T]) Expect(msg string) T {
	if r.IsErr() {
		out := fmt.Sprintf("%s: %v", msg, r.err)
		fmt.Fprintln(os.Stderr, out)
		panic(out)
	}

	return r.Ok()
}

// Then applies the function fn to the value inside the Result and returns a new Result.
// If the Result contains an error, it returns the same Result without applying fn.
func (r Result[T]) Then(fn func(T) Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}

	return fn(r.Ok())
}

// ThenOf applies the function fn to the value inside the Result, expecting
// fn to return a tuple (T, error), and returns a new Result based on the
// returned tuple. If the Result contains an error, it returns the same Result
// without applying fn.
func (r Result[T]) ThenOf(fn func(T) (T, error)) Result[T] {
	if r.IsErr() {
		return r
	}

	return ResultOf(fn(r.Ok()))
}

// Option converts a Result into an Option.
// If the Result contains an error, it returns None.
// If the Result contains a value, it returns Some with the value.
// Parameters:
//   - r: The input Result to convert into an Option.
//
// Returns:
//
//	An Option representing the value of the Result, if any.
func (r Result[T]) Option() Option[T] {
	if r.IsErr() {
		return None[T]()
	}

	return Some(r.Ok())
}

func (r Result[T]) String() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok(%v)", r.Ok())
	}

	return fmt.Sprintf("Err(%s)", r.Err().Error())
}
