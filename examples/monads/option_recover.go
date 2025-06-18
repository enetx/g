package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	processOptionData(None[string]())
	fmt.Println("---")
	processOptionData(Some("hello"))
}

func processOptionData(opt Option[string]) {
	defer func() {
		// Use recover to catch the panic.
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic.\n")

			// A type switch is the idiomatic way to handle recovered values.
			switch v := r.(type) {
			case string:
				// This is the expected path for an Option.Unwrap() panic.
				fmt.Printf("Recovered a string panic: %q\n", v)
			case error:
				// This would handle panics from Result.Unwrap() or elsewhere.
				fmt.Printf("Recovered an error object: %v\n", v)
			default:
				// Handle any other unexpected panic types.
				fmt.Printf("Recovered an unknown panic type: %T, value: %v\n", v, v)
			}
		}
	}()

	fmt.Println("Processing data...")
	// This will panic because the Option is None.
	value := opt.Unwrap()
	fmt.Println("This line will not be reached. Value:", value)
}
