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
//   - Simple Placeholder: "{key}"
//     - Replaced with the value from the map corresponding to "key".
//   - Placeholder with Fallback: "{key?fallback}"
//     - Uses the "fallback" key if "key" is not found in the map.
//   - Placeholder with Modifiers: "{key.$modifiers}"
//     - Applies transformations to the value using the specified modifiers.
//   - Placeholder with Fallback and Modifiers: "{key?fallback.$modifiers}"
//     - Applies modifiers to the value and uses fallback logic if "key" is not found.
//   - Placeholder with Parameters: "{key.$modifier(param1,param2,...)}"
//     - Passes parameters to the specified modifier for additional control or customization.
//   - Placeholder with Fallback, Modifiers, and Parameters: "{key?fallback.$modifier(param1,param2,...)}"
//     - Combines fallback logic, modifiers, and parameters for maximum flexibility.
//
// Supported Modifiers:
//   - "$abs": Returns the absolute value for numeric inputs.
//   - "$base64d": Decodes a Base64-encoded string.
//   - "$base64e": Encodes the string in Base64 format.
//   - "$bin": Converts a numeric value to binary representation.
//   - "$bool": Converts a boolean value to "true" or "false".
//   - "$date": Formats a time.Time value using the provided format string (e.g., {time.$date(2006-01-02)}).
//   - "$hex": Converts a numeric value to hexadecimal representation.
//   - "$html": Encodes the string with HTML entities.
//   - "$len": Returns the length of the value as a string.
//   - "$lower": Converts the value to lowercase.
//   - "$oct": Converts a numeric value to octal representation.
//   - "$repeat": Repeats the value a specified number of times (e.g., {value.$repeat(3)}).
//   - "$replace": Replaces all occurrences of a substring in the value with another substring (e.g., {text.$replace(old,new)}).
//   - "$reverse": Reverses the string value.
//   - "$rot13": Applies ROT13 encoding to the string.
//   - "$round": Rounds a floating-point value to the nearest integer.
//   - "$substring": Extracts a substring from a string starting at a specified index and ending at another index.
//   - "$title": Converts the value to title case.
//   - "$trim": Trims leading and trailing whitespace from the value by default. If an optional parameter is provided, it trims characters in the parameter instead of whitespace.
//   - "$truncate": Truncates the value to a specified maximum length and appends "..." if truncation occurs (e.g., {text.$truncate(10)}).
//   - "$upper": Converts the value to uppercase.
//   - "$url": Encodes the string as a URL-safe string.
//   - "$xor": Performs a bitwise XOR operation on an integer value with the provided operand (e.g., {value.$xor(42)}).

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

	modrx := String(`(\$\w+)(?:\((.*?)\))?`).Regexp().Compile().Ok()
	fallrx := String(`\{(\w+)\?([\w]+)((?:\.\$[\w]+(?:\([^\)]*\))?)*)\}`).Regexp().Compile().Ok()
	placerx := String(`\{(\w+)((?:\.\$[\w]+(?:\([^\)]*\))?)*)\}`).Regexp().Compile().Ok()

	handlers := Map[String, func(any, ...String) any]{
		"$date": func(v any, params ...String) any {
			if len(params) == 0 {
				return v
			}

			if date, ok := v.(time.Time); ok {
				return date.Format(params[0].Std())
			}

			return v
		},
		"$replace": func(v any, params ...String) any {
			if len(params) < 2 {
				return v
			}

			oldS, newS := params[0], params[1]

			switch s := v.(type) {
			case String:
				return s.ReplaceAll(oldS, newS)
			case string:
				return String(s).ReplaceAll(oldS, newS)
			default:
				return v
			}
		},
		"$repeat": func(v any, params ...String) any {
			if len(params) == 0 {
				return v
			}

			counter := params[0].Trim().ToInt()

			if counter.IsErr() {
				return v
			}

			switch t := v.(type) {
			case String:
				return t.Repeat(counter.Ok())
			case string:
				return String(t).Repeat(counter.Ok())
			case Int:
				return t.String().Repeat(counter.Ok())
			case int:
				return Int(t).String().Repeat(counter.Ok())
			case Float:
				return t.String().Repeat(counter.Ok())
			case float64:
				return Float(t).String().Repeat(counter.Ok())
			default:
				return v
			}
		},
		"$truncate": func(v any, params ...String) any {
			if len(params) == 0 {
				return v
			}

			max := params[0].Trim().ToInt()
			if max.IsErr() {
				return v
			}

			switch s := v.(type) {
			case String:
				return s.Truncate(max.Ok())
			case string:
				return String(s).Truncate(max.Ok())
			default:
				return v
			}
		},
		"$substring": func(v any, params ...String) any {
			if len(params) == 0 || len(params) < 2 {
				return v
			}

			start := params[0].Trim().ToInt()
			end := params[1].Trim().ToInt()
			step := Ok[Int](1)

			if len(params) > 2 {
				step = params[2].Trim().ToInt()
			}

			if start.IsErr() || end.IsErr() || step.IsErr() {
				return v
			}

			switch s := v.(type) {
			case String:
				return s.SubString(start.Ok(), end.Ok(), step.Ok())
			case string:
				return String(s).SubString(start.Ok(), end.Ok(), step.Ok())
			default:
				return v
			}
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
		"$trim": func(v any, params ...String) any {
			if len(params) == 0 {
				switch s := v.(type) {
				case String:
					return s.Trim()
				case string:
					return String(s).Trim()
				default:
					return v
				}
			}

			switch s := v.(type) {
			case String:
				return s.TrimSet(params[0])
			case string:
				return String(s).TrimSet(params[0])
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
		"$base64e": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Encode().Base64()
			case string:
				return String(s).Encode().Base64()
			default:
				return v
			}
		},
		"$base64d": func(v any, _ ...String) any {
			switch s := v.(type) {
			case String:
				return s.Decode().Base64().Ok()
			case string:
				return String(s).Decode().Base64().Ok()
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
		"$xor": func(v any, params ...String) any {
			if len(params) == 0 {
				return v
			}

			switch s := v.(type) {
			case String:
				return s.Encode().XOR(params[0])
			case string:
				return String(s).Encode().XOR(params[0])
			default:
				return v
			}
		},
	}

	for _, custom := range customHandlers {
		handlers.Copy(custom)
	}

	parseModifier := func(mod String) (String, Slice[String]) {
		matches := mod.Regexp().FindSubmatch(modrx)
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

	parsePlaceholders(fallrx, func(matches Slice[String]) String {
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

	parsePlaceholders(placerx, func(matches Slice[String]) String {
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
