package g

import (
	"fmt"
	"io"
)

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf[T ~string](w io.Writer, format T, a ...any) (int, error) {
	return fmt.Fprintf(w, string(format), a...)
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf[T ~string](format T, a ...any) (int, error) { return fmt.Printf(string(format), a...) }

// Sprintf formats according to a format specifier and returns the resulting String.
func Sprintf[T ~string](str T, a ...any) String { return NewString(fmt.Sprintf(string(str), a...)) }

// Fprint writes the output to w using the default formats for its operands.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, a ...any) (int, error) { return fmt.Fprint(w, a...) }

// Print writes the output to standard output using the default formats for its operands.
// It returns the number of bytes written and any write error encountered.
func Print(a ...any) (int, error) { return fmt.Print(a...) }

// Sprint formats using the default formats for its operands and returns the resulting String.
// Spaces are added between operands when neither is a string.
func Sprint(a ...any) String { return NewString(fmt.Sprint(a...)) }

// Fprintln writes the output to w followed by a newline using the default formats for its operands.
// It returns the number of bytes written and any write error encountered.
func Fprintln(w io.Writer, a ...any) (int, error) { return fmt.Fprintln(w, a...) }

// Println writes the output to standard output followed by a newline using the default formats for its operands.
// It returns the number of bytes written and any write error encountered.
func Println(a ...any) (int, error) { return fmt.Println(a...) }

// Sprintln formats using the default formats for its operands and returns the resulting String.
// Spaces are added between operands when neither is a string. A newline is appended.
func Sprintln(a ...any) String { return NewString(fmt.Sprintln(a...)) }

// Format formats a string (str) by replacing placeholders with values from a map (args)
// and returns the result as a String. Placeholders in the format string should be enclosed
// in curly braces, e.g., "{name}". The values for placeholders are retrieved from the
// provided map using Sprint for formatting individual values.
//
// Parameters:
//   - str: A format specifier as a template for the formatting.
//   - args: A map containing values to replace placeholders in the format specifier.
//
// Returns:
//
//	A String containing the formatted result.
//
// Example:
//
//	values := map[string]any{
//	    "name":  "John",
//	    "age":   30,
//	    "city":  "New York",
//	}
//	format := "Hello, my name is {name}. I am {age} years old and live in {city}."
//	formatted := g.Format(format, values)
//	formatted.Print()
//
// Output:
//
//	Hello, my name is John. I am 30 years old and live in New York.
func Format[T, U ~string](str T, args Map[U, any]) String {
	result := String(str)

	args.Iter().
		ForEach(
			func(k U, v any) {
				result = result.ReplaceAll("{"+String(k)+"}", Sprint(v))
			})

	return result
}
