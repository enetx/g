package main

import (
	"fmt"

	. "github.com/enetx/g"
)

type ErrorValidation struct {
	Field  string
	Reason string
}

func (e *ErrorValidation) Error() string {
	return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Reason)
}

func main() {
	// Create a Result containing our custom error.
	err := &ErrorValidation{Field: "email", Reason: "must not be empty"}
	result := Err[string](err)

	processResultData(result)
	fmt.Println("---")
	processResultData(Ok("Hello"))
}

func processResultData(res Result[string]) {
	defer func() {
		// Use recover to catch the panic.
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic.\n")

			// Check if the recovered value is our specific error type.
			// This is possible because Unwrap panics with the original error object.
			if validationErr, ok := r.(*ErrorValidation); ok {
				fmt.Printf("Specific error type found: %T\n", validationErr)
				fmt.Printf("Invalid field: %s\n", validationErr.Field)
			} else {
				// Handle other types of errors.
				fmt.Printf("An unexpected error occurred: %v\n", r)
			}
		}
	}()

	fmt.Println("Processing data...")

	// This will panic because the Result is an Err.
	value := res.Unwrap()
	fmt.Println("This line will not be reached. Value:", value)
}
