package g

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidBinaryLength = errors.New("binary string length must be multiple of 8")
	ErrInvalidBinaryDigit  = errors.New("binary string must contain only '0' and '1'")
)

// ErrFileNotExist represents an error for when a file does not exist.
type ErrFileNotExist struct{ Msg string }

// Error returns the error message for ErrFileNotExist.
func (e *ErrFileNotExist) Error() string { return fmt.Sprintf("no such file: %s", e.Msg) }

// ErrFileClosed represents an error for when a file is already closed.
type ErrFileClosed struct{ Msg string }

// Error returns the error message for ErrFileClosed.
func (e *ErrFileClosed) Error() string {
	return fmt.Sprintf("%s: file is already closed and unlocked", e.Msg)
}

type wrappedError struct {
	msg  string
	errs []error
}

func (e *wrappedError) Error() string   { return e.msg }
func (e *wrappedError) Unwrap() []error { return e.errs }
