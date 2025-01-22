package g

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/enetx/g/f"
)

var (

	// Print writes the output to standard output using the default formats for its operands.
	// It returns the number of bytes written and any write error encountered.
	// func Print(a ...any) (int, error) { return fmt.Print(a...) }
	Print = fmt.Print

	// Println writes the output to standard output followed by a newline using the default formats for its operands.
	// It returns the number of bytes written and any write error encountered.
	Println = fmt.Println

	// Fprint writes the output to w using the default formats for its operands.
	// It returns the number of bytes written and any write error encountered.
	Fprint = fmt.Fprint

	// Fprintln writes the output to w followed by a newline using the default formats for its operands.
	// It returns the number of bytes written and any write error encountered.
	Fprintln = fmt.Fprintln
)

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf[T ~string](w io.Writer, format T, args ...any) (int, error) {
	return w.Write(Sprintf(format, args...).Bytes())
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf[T ~string](format T, args ...any) { Sprintf(format, args...).Print() }

// Sprint formats using the default formats for its operands and returns the resulting String.
// Spaces are added between operands when neither is a string.
func Sprint(a ...any) String { return NewString(fmt.Sprint(a...)) }

// Sprintln formats using the default formats for its operands and returns the resulting String.
// Spaces are added between operands when neither is a string. A newline is appended.
func Sprintln(a ...any) String { return String(fmt.Sprintln(a...)) }

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
//     neither is found, the placeholder remains unchanged (e.g. `{key?fallback}`).
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
//     - If the placeholder is empty (`{}`) or starts with a dot (`{.something}`),
//     the system automatically takes the next positional argument. If there
//     are not enough positional arguments, it leaves the braces as "{}".
//
//  8. Escaping braces with a backslash
//     - If the parser sees `\{`, it interprets it as a literal '{' (and not the
//     start of a placeholder). Likewise, `\}` is interpreted as a literal '}'.
//
// Once a placeholder’s value is determined (either from numeric position or via
// the named map), Sprintf checks for a chain of modifiers. Each modifier name
// (e.g. `$upper`) maps to a function with the signature:
//
//	func(value any, params ...String) any
//
// This function transforms the value (e.g., changing case, formatting dates, or
// performing calculations) and returns the new value. If a modifier is
// unrecognized, it is skipped.
//
// Built-in Modifiers:
//
//   - `$get(path)`
//     Retrieves a nested value from a map, slice, array, or struct by following
//     the dot-separated path. For maps, underscores in `path` are interpreted as
//     literal dots in the key. Returns `nil` if the path cannot be resolved.
//
//   - `$json`
//     Marshals the current value as JSON. If marshalling fails, returns the original value.
//
//   - `$fmt(formatString)`
//     Passes the current value to `fmt.Sprintf(formatString, value)`. If no parameters are given,
//     returns the original value.
//
//   - `$date(2006-01-02)`
//     If the current value is a `time.Time`, formats it using the Go time layout string.
//
//   - `$replace(old,new)`
//     Replaces all occurrences of `old` with `new` in the current string.
//
//   - `$repeat(n)`
//     Repeats the current string representation `n` times.
//
//   - `$truncate(n)`
//     Truncates the current string to length `n`. (Implementation can optionally append `...`.)
//
//   - `$substring(start,end[,step])`
//     Extracts a substring from the current string (start-inclusive, end-exclusive). An optional
//     `step` can also be provided.
//
//   - `$upper`, `$lower`, `$title`
//     Converts the string to uppercase, lowercase, or title case.
//
//   - `$trim`
//     Trims whitespace from the current string. If a parameter is given (e.g. `$trim(abc)`),
//     it trims only those runes (`a`, `b`, `c`) instead.
//
//   - `$len`
//     Computes the length of the current string and returns it as a string.
//
//   - `$round`, `$abs`
//     For numeric values: `$round` rounds a float to an integer, or `$round(n)` rounds to `n` decimal
//     places. `$abs` returns the absolute value.
//
//   - `$bool`
//     Converts a `bool` to the string `"true"` or `"false"`. Other types are returned unchanged.
//
//   - `$reverse`
//     Reverses the current string.
//
//   - `$hex`, `$oct`, `$bin`
//     Converts numeric values to hexadecimal/octal/binary string representations. If the current
//     value is a string, it converts the *byte-encoded* version of that string instead.
//
//   - `$url`, `$html`
//     Escapes the current string for safe inclusion in a URL or HTML. (Typically applies
//     percent-encoding for `$url`, and HTML entities for `$html`.)
//
//   - `$base64e`, `$base64d`
//     Base64-encodes (`$base64e`) or decodes (`$base64d`) the current string.
//
//   - `$rot13`
//     Applies a simple ROT13 transformation to the current string.
//
//   - `$xor(key)`
//     XOR-encrypts the current string using the provided key.
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
	var (
		named      Named
		positional Slice[any]
	)

	for _, arg := range args {
		switch x := arg.(type) {
		case Named:
			named = x
		default:
			positional = positional.Append(x)
		}
	}

	return parseTmpl(String(template), named, positional)
}

func parseTmpl(tmpl String, named Named, positional Slice[any]) String {
	builder := NewBuilder()
	length := tmpl.Len()
	builder.Grow(length)

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
		mods.
			Prepend(".").
			Split(".$").
			Exclude(f.IsZero).
			ForEach(func(segment String) {
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

	value := Map[String, any](named).Get(key)
	if value.IsNone() && fall.NotEmpty() {
		value = Map[String, any](named).Get(fall)
	}

	return value.Some()
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

func constructKey(keyType reflect.Type, key string) Result[reflect.Value] {
	switch keyType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64, err := strconv.ParseInt(key, 10, keyType.Bits())
		if err != nil {
			return Err[reflect.Value](err)
		}
		return Ok(reflect.ValueOf(i64).Convert(keyType))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64, err := strconv.ParseUint(key, 10, keyType.Bits())
		if err != nil {
			return Err[reflect.Value](err)
		}
		return Ok(reflect.ValueOf(u64).Convert(keyType))
	case reflect.Bool:
		b, err := strconv.ParseBool(key)
		if err != nil {
			return Err[reflect.Value](err)
		}
		return Ok(reflect.ValueOf(b))
	case reflect.Float32, reflect.Float64:
		fl, err := strconv.ParseFloat(key, keyType.Bits())
		if err != nil {
			return Err[reflect.Value](err)
		}
		return Ok(reflect.ValueOf(fl).Convert(keyType))
	default:
		switch keyType {
		case reflect.TypeOf(""):
			return Ok(reflect.ValueOf(key))
		case reflect.TypeOf(String("")):
			return Ok(reflect.ValueOf(String(key)))
		default:
			return Err[reflect.Value](fmt.Errorf("unsupported key type: %s", keyType))
		}
	}
}

func extractFromMapOrd(v reflect.Value, key String) Option[any] {
	slice := v.Interface().(MapOrd[any, any])
	if slice.Empty() {
		return None[any]()
	}

	mapKeyType := reflect.ValueOf(slice[0].Key).Type()

	switch mapKeyType {
	case reflect.TypeOf(""):
		return slice.Get(key.Std())
	case reflect.TypeOf(String("")):
		return slice.Get(key)
	case reflect.TypeOf(0):
		return slice.Get(key.ToInt().Ok().Std())
	case reflect.TypeOf(Int(0)):
		return slice.Get(key.ToInt().Ok())
	case reflect.TypeOf(0.0):
		return slice.Get(key.ToFloat().Ok().Std())
	case reflect.TypeOf(Float(0.0)):
		return slice.Get(key.ToFloat().Ok())
	default:
		return None[any]()
	}
}

func resolveIndirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}
		}

		v = v.Elem()
	}

	return v
}

func applyMod(value any, name String, params Slice[String]) any {
	switch name {
	case "get":
		if len(params) == 0 {
			return value
		}

		path := params[0].Split(".").Map(func(s String) String { return s.ReplaceAll("_", ".") })
		current := reflect.ValueOf(value)

		for part := range path {
			if !current.IsValid() {
				return nil
			}

			switch current.Kind() {
			case reflect.Map:
				keyType := current.Type().Key()

				key := constructKey(keyType, part.Std())
				if key.IsErr() {
					return nil
				}

				current = current.MapIndex(key.Ok())
				current = resolveIndirect(current)
			case reflect.Slice, reflect.Array:
				if current.Type() == reflect.TypeOf(MapOrd[any, any]{}) {
					pair := extractFromMapOrd(current, part)
					if pair.IsNone() {
						return nil
					}
					current = reflect.ValueOf(pair.Some())
				} else {
					index := part.ToInt()
					if index.IsErr() || index.Ok().Gte(Int(current.Len())) {
						return nil
					}
					current = current.Index(index.Ok().Std())
				}
			case reflect.Struct:
				current = current.FieldByName(part.Std())
			default:
				return nil
			}
		}

		if current.IsValid() && current.CanInterface() {
			return current.Interface()
		}

		return nil
	case "json":
		jsonData, err := json.Marshal(value)
		if err != nil {
			return value
		}

		return String(jsonData)
	case "fmt":
		if len(params) == 0 {
			return value
		}

		return fmt.Sprintf(params[0].Std(), value)
	case "date":
		if len(params) == 0 {
			return value
		}

		if date, ok := value.(time.Time); ok {
			return date.Format(params[0].Std())
		}

		return value
	case "replace":
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
	case "repeat":
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
	case "truncate":
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
	case "substring":
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
	case "upper":
		switch s := value.(type) {
		case String:
			return s.Upper()
		case string:
			return String(s).Upper()
		default:
			return value
		}
	case "lower":
		switch s := value.(type) {
		case String:
			return s.Lower()
		case string:
			return String(s).Lower()
		default:
			return value
		}
	case "title":
		switch s := value.(type) {
		case String:
			return s.Title()
		case string:
			return String(s).Title()
		default:
			return value
		}
	case "trim":
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
	case "len":
		switch s := value.(type) {
		case String:
			return s.Len().String()
		case string:
			return String(s).Len().String()
		default:
			return value
		}
	case "round":
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
	case "abs":
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
	case "bool":
		if b, ok := value.(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
		return value
	case "reverse":
		switch s := value.(type) {
		case String:
			return s.Reverse()
		case string:
			return String(s).Reverse()
		default:
			return value
		}
	case "hex":
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
	case "oct":
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
	case "bin":
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
	case "url":
		switch s := value.(type) {
		case String:
			return s.Encode().URL()
		case string:
			return String(s).Encode().URL()
		default:
			return value
		}
	case "html":
		switch s := value.(type) {
		case String:
			return s.Encode().HTML()
		case string:
			return String(s).Encode().HTML()
		default:
			return value
		}
	case "base64e":
		switch s := value.(type) {
		case String:
			return s.Encode().Base64()
		case string:
			return String(s).Encode().Base64()
		default:
			return value
		}
	case "base64d":
		switch s := value.(type) {
		case String:
			return s.Decode().Base64().Ok()
		case string:
			return String(s).Decode().Base64().Ok()
		default:
			return value
		}
	case "rot13":
		switch s := value.(type) {
		case String:
			return s.Encode().Rot13()
		case string:
			return String(s).Encode().Rot13()
		default:
			return value
		}
	case "xor":
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
