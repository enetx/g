package g

import (
	"fmt"
	"io"
	"time"

	"github.com/enetx/g/f"
)

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf[T ~string](w io.Writer, format T, args ...any) (int, error) {
	return w.Write(Sprintf(format, args...).Bytes())
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf[T ~string](format T, args ...any) { Sprintf(format, args...).Print() }

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
func Sprintln(a ...any) String { return String(fmt.Sprintln(a...)) }

type Named Map[String, any]

// Sprintf processes a template string and replaces placeholders with values
// provided in the arguments. Placeholders are enclosed in curly braces `{...}` and
// can be either numeric (1-based) or named. Named placeholders can also have an
// optional fallback (`?fallback`) and chainable modifiers (`.$upper`, etc.).
// The function returns the final replaced string.
//
// The arguments passed to Sprintf may include:
//   - Named: a custom type (map-like) that contains key-value pairs for named
//     placeholders, e.g. {key}, {key?fallback}.
//   - Any other values (strings, numbers, etc.) become positional arguments,
//     accessible via numeric placeholders, e.g. {1}, {2}, etc.
//
// Placeholder Forms:
//
//  1. Numeric:    {1}, {2}, {3}, ...
//     - Uses 1-based indexing into the positional arguments.
//
//  2. Named:      {key}
//     - Looks up 'key' in the map of Named arguments.
//
//  3. Fallback:   {key?fallback}
//     - If "key" is not found in the named map, it tries "fallback". If
//     neither is found, the placeholder remains unchanged, e.g. {key?fallback}.
//
//  4. Modifiers:  {key.$upper}, {1.$date(2006-01-02)}, etc.
//     - Applies transformations to the placeholder value (see below).
//
//  5. Combined:   {key?fallback.$trim.$replace(a,b)}
//     - Fallback and multiple modifiers can be used together.
//
//  6. Parameters: {key.$modifier(param1,param2,...)}
//     - Passes parameters to the specified modifier function.
//
//  7. Auto-index: {}
//     - If the placeholder is empty ({}) or starts with a dot ({.something}),
//     the system automatically takes the next positional argument. If there
//     are not enough positional arguments, it leaves the braces as "{}".
//
//  8. Escaping braces with a backslash
//     - If the parser sees `\{`, it interprets it as a literal '{' (and not the
//     start of a placeholder). Likewise, `\}` is interpreted as a literal '}'.
//
// Once a placeholder’s value is determined (either from numeric position or via
// the named map), Sprintf checks for a chain of modifiers. Each modifier name
// (e.g. "$upper") maps to a function with the signature:
//
//	func(value any, params ...String) any
//
// This function transforms the value (e.g., changing case, formatting dates, or
// performing calculations) and returns the new value. If a modifier is
// unrecognized, it is skipped.
//
// Built-in Modifiers:
//
//   - `$date(2006-01-02)`: Formats a time.Time using the specified Go layout.
//   - `$replace(old,new)`: Replaces all occurrences of `old` with `new` in a string.
//   - `$repeat(n)`: Repeats a string or numeric text n times.
//   - `$truncate(n)`: Truncates a string to length n, optionally appending "..."
//   - `$substring(start,end[,step])`: Extracts a substring with an optional step.
//   - `$upper`, `$lower`, `$title`, `$trim`, `$len`: Basic string transformations
//     (uppercase, lowercase, title case, trim) and length calculation.
//   - `$round`, `$abs`, `$bool`, `$reverse`: Rounds floats, calculates absolute
//     values, converts booleans to "true"/"false", or reverses strings.
//   - `$hex`, `$oct`, `$bin`: Converts numeric/string data to hex, octal, or binary forms.
//   - `$url`, `$html`: Encodes the string as URL-safe or HTML entities.
//   - `$base64e`, `$base64d`: Base64-encode or decode string data.
//   - `$rot13`: Applies a ROT13 transform to a string.
//   - `$xor(key)`: Performs XOR transformation on a string with the given key.
//
// Example Usage:
//
//	format := "Hello, {1}! Fallback: {city?unknown} => {city.$upper}"
//	result := g.Sprintf(format,
//	                  "Alice", // => {1}
//	                  g.Named{"city": "New York"})
//
// Explanation:
//
//	// {1} -> "Alice"
//	// {city?unknown} -> "New York" (no fallback needed)
//	// {city.$upper} -> "NEW YORK"
//
//	// If a placeholder cannot be resolved (missing key, out-of-range index),
//	// it remains e.g. "{unknown}". If modifiers don’t apply, they’re skipped.
func Sprintf[T ~string](template T, args ...any) String {
	named := make(Named)
	var positional Slice[any]

	for _, arg := range args {
		switch x := arg.(type) {
		case Named:
			for k, v := range x {
				named[k] = v
			}
		default:
			positional = positional.Append(x)
		}
	}

	return parseTmpl(String(template), named, positional)
}

func parseTmpl(tmpl String, named Named, positional Slice[any]) String {
	builder := NewBuilder()
	length := tmpl.Len()

	var autoidx, idx Int

	for idx < length {
		if tmpl[idx] == '\\' && idx+1 < length {
			next := tmpl[idx+1]
			if next == '{' || next == '}' {
				builder.WriteByte(next)
				idx += 2

				continue
			}
		}

		if tmpl[idx] == '{' {
			cidx := tmpl[idx+1:].Index("}")
			if cidx.IsNegative() {
				builder.WriteByte(tmpl[idx])
				idx++

				continue
			}

			eidx := idx + 1 + cidx
			placeholder := tmpl[idx+1 : eidx]

			trimmed := placeholder.Trim()
			if trimmed.Empty() || trimmed[0] == '.' {
				autoidx++
				if autoidx <= positional.Len() {
					placeholder = autoidx.String() + placeholder
				}
			}

			replaced := processPlaceholder(placeholder, named, positional)
			builder.Write(replaced)

			idx = eidx + 1
		} else {
			builder.WriteByte(tmpl[idx])
			idx++
		}
	}

	return builder.String()
}

func processPlaceholder(placeholder String, named Named, positional Slice[any]) String {
	var (
		keyfall String
		mods    String
		key     String
		fall    String
	)

	if idx := placeholder.Index("."); idx.IsPositive() {
		keyfall = placeholder[:idx]
		mods = placeholder[idx+1:]
	} else {
		keyfall = placeholder
	}

	if idx := keyfall.Index("?"); idx.IsPositive() {
		key = keyfall[:idx]
		fall = keyfall[idx+1:]
	} else {
		key = keyfall
	}

	value := resolveValue(key, fall, named, positional)
	if value == nil {
		return "{" + placeholder + "}"
	}

	if mods.NotEmpty() {
		mods.Split(".").Exclude(f.IsZero).ForEach(func(segment String) {
			name, params := parseMod(segment)
			value = applyMod(value, name, params)
		})
	}

	return Sprint(value)
}

func resolveValue(key, fall String, named Named, positional Slice[any]) any {
	if num := key.ToInt(); num.IsOk() {
		idx := num.Ok() - 1
		if idx.IsNegative() || idx.Gte(positional.Len()) {
			return nil
		}

		return positional[idx]
	}

	value, ok := named[key]
	if !ok && fall.NotEmpty() {
		value, ok = named[fall]
		if !ok {
			return nil
		}
	}

	return value
}

func parseMod(segment String) (String, Slice[String]) {
	oidx := segment.Index("(")
	if oidx.IsNegative() {
		return segment, nil
	}

	cidx := segment.LastIndex(")")
	if cidx.Lt(oidx) {
		return segment, nil
	}

	params := segment[oidx+1 : cidx].Split(",").Collect()
	name := segment[:oidx]

	return name, params
}

func applyMod(value any, name String, params Slice[String]) any {
	switch name {
	case "$fmt":
		if len(params) == 0 {
			return value
		}

		return fmt.Sprintf(params[0].Std(), value)
	case "$date":
		if len(params) == 0 {
			return value
		}

		if date, ok := value.(time.Time); ok {
			return date.Format(params[0].Std())
		}

		return value
	case "$replace":
		if len(params) < 2 {
			return value
		}

		oldS, newS := params[0], params[1]

		switch s := value.(type) {
		case String:
			return s.ReplaceAll(oldS, newS)
		case string:
			return String(s).ReplaceAll(oldS, newS)
		default:
			return value
		}
	case "$repeat":
		if len(params) == 0 {
			return value
		}

		counter := params[0].Trim().ToInt()
		if counter.IsErr() {
			return value
		}

		switch t := value.(type) {
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
			return value
		}
	case "$truncate":
		if len(params) == 0 {
			return value
		}

		max := params[0].Trim().ToInt()
		if max.IsErr() {
			return value
		}

		switch s := value.(type) {
		case String:
			return s.Truncate(max.Ok())
		case string:
			return String(s).Truncate(max.Ok())
		default:
			return value
		}
	case "$substring":
		if len(params) == 0 || len(params) < 2 {
			return value
		}

		start := params[0].Trim().ToInt()
		end := params[1].Trim().ToInt()
		step := Ok[Int](1)

		if len(params) > 2 {
			step = params[2].Trim().ToInt()
		}

		if start.IsErr() || end.IsErr() || step.IsErr() {
			return value
		}

		switch s := value.(type) {
		case String:
			return s.SubString(start.Ok(), end.Ok(), step.Ok())
		case string:
			return String(s).SubString(start.Ok(), end.Ok(), step.Ok())
		default:
			return value
		}
	case "$upper":
		switch s := value.(type) {
		case String:
			return s.Upper()
		case string:
			return String(s).Upper()
		default:
			return value
		}
	case "$lower":
		switch s := value.(type) {
		case String:
			return s.Lower()
		case string:
			return String(s).Lower()
		default:
			return value
		}
	case "$title":
		switch s := value.(type) {
		case String:
			return s.Title()
		case string:
			return String(s).Title()
		default:
			return value
		}
	case "$trim":
		if len(params) == 0 {
			switch s := value.(type) {
			case String:
				return s.Trim()
			case string:
				return String(s).Trim()
			default:
				return value
			}
		}

		switch s := value.(type) {
		case String:
			return s.TrimSet(params[0])
		case string:
			return String(s).TrimSet(params[0])
		default:
			return value
		}
	case "$len":
		switch s := value.(type) {
		case String:
			return s.Len().String()
		case string:
			return String(s).Len().String()
		default:
			return value
		}
	case "$round":
		if len(params) == 0 {
			switch fl := value.(type) {
			case Float:
				return fl.Round().String()
			case float64:
				return Float(fl).Round().String()
			default:
				return value
			}
		}

		precision := params[0].Trim().ToInt()
		if precision.IsErr() {
			return value
		}

		switch fl := value.(type) {
		case Float:
			return fl.RoundDecimal(precision.Ok()).String()
		case float64:
			return Float(fl).RoundDecimal(precision.Ok()).String()
		default:
			return value
		}
	case "$abs":
		switch n := value.(type) {
		case Int:
			return n.Abs().String()
		case int:
			return Int(n).Abs().String()
		case Float:
			return n.Abs().String()
		case float64:
			return Float(n).Abs().String()
		default:
			return value
		}
	case "$bool":
		if b, ok := value.(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
		return value
	case "$reverse":
		switch s := value.(type) {
		case String:
			return s.Reverse()
		case string:
			return String(s).Reverse()
		default:
			return value
		}
	case "$hex":
		switch n := value.(type) {
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
			return value
		}
	case "$oct":
		switch n := value.(type) {
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
			return value
		}
	case "$bin":
		switch n := value.(type) {
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
			return value
		}
	case "$url":
		switch s := value.(type) {
		case String:
			return s.Encode().URL()
		case string:
			return String(s).Encode().URL()
		default:
			return value
		}
	case "$html":
		switch s := value.(type) {
		case String:
			return s.Encode().HTML()
		case string:
			return String(s).Encode().HTML()
		default:
			return value
		}
	case "$base64e":
		switch s := value.(type) {
		case String:
			return s.Encode().Base64()
		case string:
			return String(s).Encode().Base64()
		default:
			return value
		}
	case "$base64d":
		switch s := value.(type) {
		case String:
			return s.Decode().Base64().Ok()
		case string:
			return String(s).Decode().Base64().Ok()
		default:
			return value
		}
	case "$rot13":
		switch s := value.(type) {
		case String:
			return s.Encode().Rot13()
		case string:
			return String(s).Encode().Rot13()
		default:
			return value
		}
	case "$xor":
		if len(params) == 0 {
			return value
		}

		switch s := value.(type) {
		case String:
			return s.Encode().XOR(params[0])
		case string:
			return String(s).Encode().XOR(params[0])
		default:
			return value
		}
	default:
		return value
	}
}
