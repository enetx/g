package g

import "fmt"

type (
	ErrFileNotExist struct{ Msg string }
	ErrFileClosed   struct{ Msg string }
)

func (e *ErrFileNotExist) Error() string {
	return fmt.Sprintf("no such file: %s", e.Msg)
}

func (e *ErrFileClosed) Error() string {
	return fmt.Sprintf("%s: file is already closed and unlocked", e.Msg)
}
