package g

import "fmt"

// Sprintf formats according to a format specifier and returns the resulting String.
func Sprintf[T ~string](str T, a ...any) String { return NewString(fmt.Sprintf(string(str), a...)) }

// Sprint formats using the default formats for its operands and returns the resulting String.
// Spaces are added between operands when neither is a string.
func Sprint(a ...any) String { return NewString(fmt.Sprint(a...)) }

// GSprintf formats a string (str) by replacing placeholders with values from a map (args)
// and returns the result as a String. Placeholders in the format string should be enclosed
// in curly braces, e.g., "{Name}". The values for placeholders are retrieved from the
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
//	    "Name":  "John",
//	    "Age":   30,
//	    "City":  "New York",
//	}
//	formatString := "Hello, my name is {Name}. I am {Age} years old and live in {City}."
//	formattedString := GSprintf(formatString, values)
//	formattedString.Print()
//
// Output:
//
//	Hello, my name is John. I am 30 years old and live in New York.
func GSprintf[T, U ~string](str T, args Map[U, any]) String {
	result := String(str)

	args.Iter().
		ForEach(
			func(k U, v any) {
				result = result.ReplaceAll("{"+String(k)+"}", Sprint(v))
			})

	return result
}
