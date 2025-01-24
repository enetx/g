package g

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

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

// Sprintf processes a template string and replaces placeholders with corresponding values from the provided arguments.
// It supports numeric, named, and auto-indexed placeholders, as well as dynamic invocation of methods on values.
//
// If a placeholder cannot resolve a value or an invoked method fails, the placeholder remains unchanged in the output.
//
// Parameters:
//   - template (T ~string): A string containing placeholders enclosed in `{}`.
//   - args (...any): A variadic list of arguments, which may include:
//   - Positional arguments (numbers, strings, slices, structs, maps, etc.).
//   - A `Named` map for named placeholders.
//
// Placeholder Forms:
//   - Numeric: `{1}`, `{2}` - References positional arguments by their 1-based index.
//   - Named: `{key}`, `{key.MethodName(param1, param2)}` - References keys from a `Named` map and allows method invocation.
//   - Fallback: `{key?fallback}` - Uses `fallback` if the key is not found in the named map.
//   - Auto-index: `{}` - Automatically uses the next positional argument if the placeholder is empty.
//   - Escaping: `\{` and `\}` - Escapes literal braces in the template string.
//
// Returns:
//   - String: A formatted string with all resolved placeholders replaced by their corresponding values.
//
// Notes:
//   - If a placeholder cannot resolve a value (e.g., missing key or out-of-range index), it remains unchanged in the output.
//   - Method invocation supports any type with accessible methods. If the method or its parameters are invalid, the value remains unmodified.
//
// Usage:
//
//	// Example 1: Numeric placeholders
//	result := g.Sprintf("{1} + {2} = {3}", 1, 2, 3)
//
//	// Example 2: Named placeholders
//	named := g.Named{
//		"name": "Alice",
//		"age":  30,
//	}
//	result := g.Sprintf("My name is {name} and I am {age} years old.", named)
//
//	// Example 3: Method invocation on values
//	result := g.Sprintf("Hex: {1.Hex}, Binary: {1.Binary}", g.Int(255))
//
//	// Example 4: Fallbacks and chaining
//	named := g.Named{
//		"name": g.String("   john  "),
//		"city": g.String("New York"),
//	}
//	result := g.Sprintf("Hello, {name.Trim.Title}. Welcome to {city?Unknown}!", named)
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
			Split(".").
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

func applyMod(value any, name String, params Slice[String]) any {
	val := reflect.ValueOf(value)

	method := val.MethodByName(string(name))
	if !method.IsValid() || method.Kind() != reflect.Func {
		return value
	}

	methodType := method.Type()
	numIn := methodType.NumIn()
	isVariadic := methodType.IsVariadic()

	if isVariadic {
		numIn--
	}

	var args []reflect.Value

	for i := range numIn {
		arg := toType(params[i], methodType.In(i))
		if arg.IsErr() {
			return value
		}

		args = append(args, arg.Ok())
	}

	if isVariadic {
		elemType := methodType.In(numIn).Elem()

		for param := range params[numIn:].Iter() {
			arg := toType(param, elemType)
			if arg.IsErr() {
				return value
			}

			args = append(args, arg.Ok())
		}
	}

	results := method.Call(args)

	if len(results) > 0 {
		return results[0].Interface()
	}

	return value
}

func toType(param String, targetType reflect.Type) Result[reflect.Value] {
	switch targetType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64, err := strconv.ParseInt(param.Std(), 10, targetType.Bits())
		if err != nil {
			return Err[reflect.Value](err)
		}
		return Ok(reflect.ValueOf(i64).Convert(targetType))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64, err := strconv.ParseUint(param.Std(), 10, targetType.Bits())
		if err != nil {
			return Err[reflect.Value](err)
		}
		return Ok(reflect.ValueOf(u64).Convert(targetType))
	case reflect.Bool:
		b, err := strconv.ParseBool(param.Std())
		if err != nil {
			return Err[reflect.Value](err)
		}
		return Ok(reflect.ValueOf(b))
	case reflect.Float32, reflect.Float64:
		param = param.ReplaceAll("_", ".")
		fl, err := strconv.ParseFloat(param.Std(), targetType.Bits())
		if err != nil {
			return Err[reflect.Value](err)
		}
		return Ok(reflect.ValueOf(fl).Convert(targetType))
	default:
		switch targetType {
		case reflect.TypeOf(""):
			return Ok(reflect.ValueOf(param.Std()))
		case reflect.TypeOf(String("")):
			return Ok(reflect.ValueOf(param))
		default:
			return Err[reflect.Value](fmt.Errorf("unsupported type: %s", targetType))
		}
	}
}
