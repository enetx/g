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

// Format formats a string by replacing placeholders with corresponding values from a map
// and returns the resulting formatted string. Placeholders in the input string should
// be enclosed in curly braces, e.g., "{key}". Optional modifiers can be used to transform
// the values (e.g., "$upper" for uppercase, "$lower" for lowercase, etc.).
//
// Parameters:
//   - str: A template string containing placeholders enclosed in curly braces.
//   - args: A map containing keys and their associated values for replacing placeholders.
//
// Placeholder Syntax:
//   - Simple Placeholder: "{key}" - Replaced with the value from the map corresponding to "key".
//   - Placeholder with Fallback: "{key?fallback}" - Uses the "fallback" key if "key" is not found.
//   - Placeholder with Modifiers: "{$modifiers:key}" - Applies transformations to the value using the specified modifiers.
//   - Placeholder with Modifiers and Fallback: "{$modifiers:key?fallback}" - Applies modifiers and uses fallback logic.
//
// Supported Modifiers:
//   - "$upper": Converts the value to uppercase.
//   - "$lower": Converts the value to lowercase.
//   - "$title": Converts the value to title case.
//   - "$trim": Trims leading and trailing whitespace from the value.
//   - "$len": Returns the length of the value as a string.
//   - "$round": Rounds a floating-point value to the nearest integer.
//   - "$abs": Returns the absolute value for numeric inputs.
//   - "$bool": Converts a boolean value to "true" or "false".
//   - "$reverse": Reverses the string value.
//   - "$hex": Converts a numeric value to hexadecimal representation.
//   - "$oct": Converts a numeric value to octal representation.
//   - "$bin": Converts a numeric value to binary representation.
//   - "$url": Encodes the string as a URL-safe string.
//   - "$html": Encodes the string with HTML entities.
//   - "$base64": Encodes the string in Base64 format.
//   - "$rot13": Applies ROT13 encoding to the string.
//
// Returns:
//   - A formatted string with placeholders replaced by their corresponding values from the map.
//
// Example Usage:
//
//	values := map[string]any{
//	    "name":  "John",
//	    "age":   30,
//	    "city":  "New York",
//	}
//	format := "Hello, my name is {name}. I am {age} years old and live in {city}."
//	formatted := g.Format(format, values)
//	formatted.Println()
//
// Output:
//
//	Hello, my name is John. I am 30 years old and live in New York.
//
// Example with Modifiers:
//
//	values := map[string]any{
//	    "name":  "John",
//	    "work": " developer ",
//	}
//	format := "Name: {$upper:name}, Title: {$trim.$title:work}"
//	formatted := g.Format(format, values)
//	formatted.Print()
//
// Output:
//
//	Name: JOHN, Title: Developer
func Format[T, U ~string](str T, args Map[U, any]) String {
	result := String(str)

	applyRegex := func(regex String, handler func(String, Option[Slice[String]]) String) {
		re := regex.Regexp().Compile().Ok()
		result = result.Regexp().ReplaceBy(re, func(match String) String {
			return handler(match, match.Regexp().FindSubmatch(re))
		})
	}

	modifierHandlers := Map[String, func(any) any]{
		"$upper": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Upper()
			case string:
				return String(s).Upper()
			default:
				return v
			}
		},
		"$lower": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Lower()
			case string:
				return String(s).Lower()
			default:
				return v
			}
		},
		"$title": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Title()
			case string:
				return String(s).Title()
			default:
				return v
			}
		},
		"$trim": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Trim()
			case string:
				return String(s).Trim()
			default:
				return v
			}
		},
		"$len": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Len().String()
			case string:
				return String(s).Len().String()
			default:
				return v
			}
		},
		"$round": func(v any) any {
			switch f := v.(type) {
			case Float:
				return f.Round().String()
			case float64:
				return Float(f).Round().String()
			default:
				return v
			}
		},
		"$abs": func(v any) any {
			switch n := v.(type) {
			case Int:
				return n.Abs().String()
			case int:
				return Int(n).Abs().String()
			case Float:
				return n.Abs().String()
			case float64:
				return Float(n).Abs().String()
			default:
				return v
			}
		},
		"$bool": func(v any) any {
			if b, ok := v.(bool); ok {
				if b {
					return "true"
				}
				return "false"
			}
			return v
		},
		"$reverse": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Reverse()
			case string:
				return String(s).Reverse()
			default:
				return v
			}
		},
		"$hex": func(v any) any {
			switch n := v.(type) {
			case Int:
				return n.Hex()
			case int:
				return Int(n).Hex()
			case Float:
				return n.Int().Hex()
			case float64:
				return Int(n).Hex()
			case String:
				return n.Encode().Hex()
			case string:
				return String(n).Encode().Hex()
			default:
				return v
			}
		},
		"$oct": func(v any) any {
			switch n := v.(type) {
			case Int:
				return n.Octal()
			case int:
				return Int(n).Octal()
			case Float:
				return n.Int().Octal()
			case float64:
				return Int(n).Octal()
			case String:
				return n.Encode().Octal()
			case string:
				return String(n).Encode().Octal()
			default:
				return v
			}
		},
		"$bin": func(v any) any {
			switch n := v.(type) {
			case Int:
				return n.Binary()
			case int:
				return Int(n).Binary()
			case Float:
				return n.Int().Binary()
			case float64:
				return Int(n).Binary()
			case String:
				return n.Encode().Binary()
			case string:
				return String(n).Encode().Binary()
			default:
				return v
			}
		},
		"$url": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Encode().URL()
			case string:
				return String(s).Encode().URL()
			default:
				return v
			}
		},
		"$html": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Encode().HTML()
			case string:
				return String(s).Encode().HTML()
			default:
				return v
			}
		},
		"$base64": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Encode().Base64()
			case string:
				return String(s).Encode().Base64()
			default:
				return v
			}
		},
		"$rot13": func(v any) any {
			switch s := v.(type) {
			case String:
				return s.Encode().Rot13()
			case string:
				return String(s).Encode().Rot13()
			default:
				return v
			}
		},
	}
	processModifiers := func(modifiers String, v any) String {
		for modifier := range modifiers.Split(".") {
			if handler := modifierHandlers.Get(modifier); handler.IsSome() {
				v = handler.Some()(v)
			}
		}

		return Sprint(v)
	}

	applyRegex(String(`\{([\$\w\.]+):(\w+)\?(\w+)\}`), func(match String, matches Option[Slice[String]]) String {
		if matches.IsNone() || matches.Some().Len() != 4 {
			return match
		}

		var modifiers, key, fallbackKey String
		matches.Some()[1:].Unpack(&modifiers, &key, &fallbackKey)

		value := args.Get(U(key))
		if value.IsNone() {
			value = args.Get(U(fallbackKey))
		}

		if value.IsNone() {
			return match
		}

		return processModifiers(modifiers, value.Some())
	})

	applyRegex(String(`\{([\$\w\.]+):(\w+)\}`), func(match String, matches Option[Slice[String]]) String {
		if matches.IsNone() || matches.Some().Len() != 3 {
			return match
		}

		var modifiers, key String
		matches.Some()[1:].Unpack(&modifiers, &key)

		value := args.Get(U(key))
		if value.IsNone() {
			return match
		}

		return processModifiers(modifiers, value.Some())
	})

	applyRegex(String(`\{(\w+)\?(\w+)\}`), func(match String, matches Option[Slice[String]]) String {
		if matches.IsNone() || matches.Some().Len() != 3 {
			return match
		}

		var key, fallbackKey String
		matches.Some()[1:].Unpack(&key, &fallbackKey)

		value := args.Get(U(key))
		if value.IsNone() {
			value = args.Get(U(fallbackKey))
		}

		if value.IsNone() {
			return match
		}

		return Sprint(value.Some())
	})

	applyRegex(String(`\{(\w+)\}`), func(match String, matches Option[Slice[String]]) String {
		if matches.IsNone() || matches.Some().Len() != 2 {
			return match
		}

		if value := args.Get(U(matches.Some().Get(1))); value.IsSome() {
			return Sprint(value.Some())
		}

		return match
	})

	return result
}
