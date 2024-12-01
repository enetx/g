package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestErrFileNotExist_Error(t *testing.T) {
	// Create an instance of ErrFileNotExist
	err := &ErrFileNotExist{Msg: "example.txt"}

	// Check if the error message is formatted correctly
	expectedErrorMsg := "no such file: example.txt"
	if errMsg := err.Error(); errMsg != expectedErrorMsg {
		t.Errorf("TestErrFileNotExist_Error: Expected error message '%s', got '%s'", expectedErrorMsg, errMsg)
	}
}

func TestErrFileClosed_Error(t *testing.T) {
	// Create an instance of ErrFileClosed
	err := &ErrFileClosed{Msg: "example.txt"}

	// Check if the error message is formatted correctly
	expectedErrorMsg := "example.txt: file is already closed and unlocked"
	if errMsg := err.Error(); errMsg != expectedErrorMsg {
		t.Errorf("TestErrFileClosed_Error: Expected error message '%s', got '%s'", expectedErrorMsg, errMsg)
	}
}
