package g

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/enetx/g/f"
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
//   - customHandlers (optional): Additional Map of custom modifiers where the key is the
//     modifier name and the value is a function to process the transformation.
//
// Placeholder Syntax:
//   - Simple Placeholder: "{key}" - Replaced with the value from the map corresponding to "key".
//   - Placeholder with Fallback: "{key?fallback}" - Uses the "fallback" key if "key" is not found.
//   - Placeholder with Modifiers: "{key.$modifiers}" - Applies transformations to the value using the specified modifiers.
//   - Placeholder with Fallback and Modifiers: "{key?fallback.$modifiers}" - Applies modifiers and uses fallback logic.
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
//   - "$format": Formats a time.Time value using the provided format string. (e.g., {time.$format(2006-01-02)}).
//
// Custom Modifiers:
//   - You can define custom modifiers by providing a Map where the key is the modifier
//     name (String), and the value is a function with the signature: func(any, ...String) any.
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
//
//	format := "Hello, my name is {name}. I am {age} years old and live in {city}."
//	formatted := g.Format(format, values)
//	fmt.Println(formatted)
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
//	    "today":  time.Now(),
//	}
//
//	format := "Name: {name.$upper}, Work: {work.$trim.$title}, Today: {today.$format(01/02/2006)}."
//	formatted := g.Format(format, values)
//	fmt.Println(formatted)
//
// Output:
//
//	Name: JOHN, Work: Developer, Today: 01/17/2025.
//
// Example with Custom Modifiers:
//
//	handlers := Map[String, func(v any, args ...String) any]{
//		"$double": func(v any, _ ...String) any { return (v.(Int) * 2).String() },
//		"$prefix": func(v any, _ ...String) any { return "prefix_" + v.(String) },
//	}
//
//	args := map[string]any{
//		"value": Int(42),
//		"text":  String("example"),
//		"date":  time.Now(),
//	}
//
//	result := Format("{value.$double} and {text.$upper}", args, handlers)
//	result.Println()
//
// Output:
//
//	84 and EXAMPLE
func Format[T, U ~string](str T, args Map[U, any], customHandlers ...Map[String, func(any, ...String) any]) String {
	result := String(str)

	modRx := String(`(\$\w+)(?:\((.*?)\))?`).Regexp().Compile().Ok()
	fallRx := String(`\{(\w+)\?([\w]+)((?:\.\$[\w]+(?:\([^\)]*\))?)*)\}`).Regexp().Compile().Ok()
	placeRx := String(`\{(\w+)((?:\.\$[\w]+(?:\([^\)]*\))?)*)\}`).Regexp().Compile().Ok()

	handlers := Map[String, func(any, ...String) any]{
		"$format": func(v any, params ...String) any {
			if len(params) == 0 {
				return v
			}
			if date, ok := v.(time.Time); ok {
				return date.Format(params[0].Std())
			}
			return v
		},
		"$upper": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Upper()
			case string:
				return String(s).Upper()
			default:
				return v
			}
		},
		"$lower": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Lower()
			case string:
				return String(s).Lower()
			default:
				return v
			}
		},
		"$title": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Title()
			case string:
				return String(s).Title()
			default:
				return v
			}
		},
		"$trim": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Trim()
			case string:
				return String(s).Trim()
			default:
				return v
			}
		},
		"$len": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Len().String()
			case string:
				return String(s).Len().String()
			default:
				return v
			}
		},
		"$round": func(v any, _ ...String) any {
			switch f := v.(type) {
			case Float:
				return f.Round().String()
			case float64:
				return Float(f).Round().String()
			default:
				return v
			}
		},
		"$abs": func(v any, _ ...String) any {
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
		"$bool": func(v any, _ ...String) any {
			if b, ok := v.(bool); ok {
				if b {
					return "true"
				}
				return "false"
			}
			return v
		},
		"$reverse": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Reverse()
			case string:
				return String(s).Reverse()
			default:
				return v
			}
		},
		"$hex": func(v any, _ ...String) any {
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
		"$oct": func(v any, _ ...String) any {
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
		"$bin": func(v any, _ ...String) any {
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
		"$url": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Encode().URL()
			case string:
				return String(s).Encode().URL()
			default:
				return v
			}
		},
		"$html": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Encode().HTML()
			case string:
				return String(s).Encode().HTML()
			default:
				return v
			}
		},
		"$base64": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Encode().Base64()
			case string:
				return String(s).Encode().Base64()
			default:
				return v
			}
		},
		"$rot13": func(v any, _ ...String) any {
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

	for _, custom := range customHandlers {
		handlers.Copy(custom)
	}

	parseModifier := func(mod String) (String, Slice[String]) {
		matches := mod.Regexp().FindSubmatch(modRx)
		if matches.IsSome() {
			name, params := matches.Some()[1], NewSlice[String]()
			if matches.Some().Len().Gt(2) && matches.Some()[2].NotEmpty() {
				params = matches.Some()[2].Split(",").Collect()
			}

			return name, params
		}

		return "", nil
	}

	applyModifiers := func(mods String, v any) String {
		mods.Split(".").Exclude(f.IsZero).ForEach(func(mod String) {
			name, params := parseModifier(mod)
			if handler := handlers.Get(name); handler.IsSome() {
				v = handler.Some()(v, params...)
			}
		})

		return Sprint(v)
	}

	parsePlaceholders := func(pattern *regexp.Regexp, handler func(Slice[String]) String) {
		result = result.Regexp().ReplaceBy(pattern, func(match String) String {
			matches := match.Regexp().FindSubmatch(pattern)
			if matches.IsSome() {
				return handler(matches.Some())
			}

			return match
		})
	}

	parsePlaceholders(fallRx, func(matches Slice[String]) String {
		if matches.Len().Lt(4) {
			return matches[0]
		}

		key, fallbackKey, modifiers := matches[1], matches[2], matches[3]

		value := args.Get(U(key))
		if value.IsNone() {
			value = args.Get(U(fallbackKey))
		}

		if value.IsSome() {
			return applyModifiers(modifiers, value.Some())
		}

		return matches[0]
	})

	parsePlaceholders(placeRx, func(matches Slice[String]) String {
		if matches.Len().Lt(3) {
			return matches[0]
		}

		key, modifiers := matches[1], matches[2]

		value := args.Get(U(key))
		if value.IsSome() {
			return applyModifiers(modifiers, value.Some())
		}

		return matches[0]
	})

	return result
}
