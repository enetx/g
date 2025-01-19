package g

import (
	"fmt"
	"io"
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

// Format replaces placeholders in the given template string `template` with corresponding values
// from the provided arguments. Placeholders are enclosed in curly braces `{...}` and can be either
// numeric (1-based) or named. Named placeholders can also have an optional fallback and chainable
// modifiers. The function returns the final replaced string.
//
// The arguments passed to Format can include:
//   - Any number of map-like values (map[string]any, map[String]any, Map[string,any], Map[String,any])
//     whose entries are merged into a single combined map for named placeholders.
//   - Any other values (strings, numbers, etc.) that become positional arguments, accessible via
//     numeric placeholders, e.g. {1}, {2}, etc.
//
// Supported placeholder forms:
//
//  1. Numeric:         {1}, {2}, {3}, ...
//     - Uses 1-based indexing into the non-map arguments.
//  2. Named:           {key}
//     - Looks up 'key' in the merged map.
//  3. Fallback:        {key?fallback}
//     - If 'key' is not found in the map, tries 'fallback'. If also missing, leaves the placeholder unchanged.
//  4. Modifiers:       {key.$upper}, {1.$date(2006-01-02)}, etc.
//     - Applies transformations to the placeholder value (see below).
//  5. Combined:        {key?fallback.$trim.$replace(a,b)}
//     - Fallback and multiple modifiers can be used together.
//  6. Parameters:      {key.$modifier(param1,param2,...)}
//     - Passes parameters to the specified modifier function.
//
// Once a placeholder's value is determined (either by numeric position or via the merged map), Format
// checks for a chain of modifiers. Each modifier name (e.g. "$upper") maps to a function with the signature:
//
//	func(value any, params ...String) any
//
// This function transforms 'value' (e.g., changing case, formatting dates, or performing calculations)
// and returns the new value. If a modifier is unrecognized, it is skipped.
//
// The built-in modifiers currently supported include:
//
//   - $date(2006-01-02): formats a time.Time using the specified layout;
//   - $replace(old,new): replaces 'old' with 'new' in the string;
//   - $repeat(n): repeats a string or numeric text n times;
//   - $truncate(n): truncates a string to length n, optionally appending "...";
//   - $substring(start,end[,step]): returns a substring with optional step;
//   - $upper, $lower, $title, $trim, $len: upper/lower/title case, trim, length;
//   - $round, $abs, $bool, $reverse: round floats, absolute values, booleans, reverse strings;
//   - $hex, $oct, $bin: convert numeric or string data into hex, octal, or binary forms;
//   - $url, $html: URL-encode or HTML-encode string data;
//   - $base64e, $base64d: encode/decode using Base64;
//   - $rot13: applies a ROT13 transform to a string;
//   - $xor: performs XOR transformation on a string with the given key.
//
// Example usage:
//
//	format := "Hello, {1}! Fallback: {city?unknown} => {city.$upper}"
//	result := Format(format,
//	                 "Alice",            // {1}
//	                 map[string]any{"city": "New York"})
//	// {1} -> "Alice"
//	// {city?unknown} -> "New York" (no fallback needed)
//	// {city.$upper} -> "NEW YORK"
//
// If a placeholder cannot be resolved (key is missing and no fallback, or numeric index is out of range),
// it remains intact with braces, e.g., `"{unknown}"`. If modifiers are present but the initial value
// is invalid for them, the modifiers are skipped. In all other cases, a transformed value is substituted
// in the final result.
func Format[T ~string](template T, args ...any) String {
	named := NewMap[String, any]()
	var positional Slice[any]

	for _, arg := range args {
		switch value := arg.(type) {
		case map[string]any:
			for k, v := range value {
				named[String(k)] = v
			}
		case map[String]any:
			for k, v := range value {
				named[k] = v
			}
		case Map[String, any]:
			for k, v := range value {
				named[k] = v
			}
		case Map[string, any]:
			for k, v := range value {
				named[String(k)] = v
			}
		default:
			positional = positional.Append(value)
		}
	}

	return parseTmpl(String(template), named, positional)
}

func parseTmpl(tmpl String, named Map[String, any], positional Slice[any]) String {
	builder := NewBuilder()
	length := tmpl.Len()

	for i := Int(0); i < length; {
		if tmpl[i] == '{' {
			cidx := tmpl[i+1:].Index("}")
			if cidx.IsNegative() {
				builder.WriteByte(tmpl[i])
				i++
				continue
			}

			eidx := i + 1 + cidx
			placeholder := tmpl[i+1 : eidx]

			replaced := processPlaceholder(placeholder, named, positional)
			builder.Write(replaced)

			i = eidx + 1
		} else {
			builder.WriteByte(tmpl[i])
			i++
		}
	}

	return builder.String()
}

func processPlaceholder(placeholder String, named Map[String, any], positional Slice[any]) String {
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

	vo := resolveValue(key, fall, named, positional)
	if vo.IsNone() {
		return "{" + placeholder + "}"
	}

	value := vo.Some()

	if mods.NotEmpty() {
		mods.Split(".").Exclude(f.IsZero).ForEach(func(segment String) {
			name, params := parseMod(segment)
			value = applyMod(value, name, params)
		})
	}

	return Sprint(value)
}

func resolveValue(key, fall String, named Map[String, any], positional Slice[any]) Option[any] {
	if num := key.ToInt(); num.IsOk() {
		idx := num.Ok() - 1
		if idx.IsNegative() || idx.Gte(positional.Len()) {
			return None[any]()
		}

		return Some(positional[idx])
	}

	value := named.Get(key)
	if value.IsNone() && fall.NotEmpty() {
		value = named.Get(fall)
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
		switch fl := value.(type) {
		case Float:
			return fl.Round().String()
		case float64:
			return Float(fl).Round().String()
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
