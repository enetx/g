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

// ToResult returns a new Result[T] based on the provided value and error.
// If err is not nil, it returns an Result containing the error.
// Otherwise, it returns an Result containing the value.
func ToResult[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}

	return Ok(value)
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
