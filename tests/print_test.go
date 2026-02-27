package g_test

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	. "github.com/enetx/g"
)

func TestFormatAutoIndexAndNumeric(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		{
			name:     "Autoindex multiple",
			format:   "Hello, {} and {}!",
			args:     []any{"Alice", "Bob"},
			expected: "Hello, Alice and Bob!",
		},
		{
			name:     "Autoindex out-of-range",
			format:   "Hello, {} and {}!",
			args:     []any{"Alice"},
			expected: "Hello, Alice and {}!",
		},
		{
			name:     "Numeric placeholders",
			format:   "Values: {1}, {2}, {1.Lower}",
			args:     []any{String("X"), "Y"},
			expected: "Values: X, Y, x",
		},
		{
			name:     "Escaped braces",
			format:   "Show literal \\{{.Upper}\\} here",
			args:     []any{String("upper")},
			expected: "Show literal {UPPER} here",
		},
		{
			name:     "Autoindex with modifier {.Upper}",
			format:   "{.Upper} and {.Lower}",
			args:     []any{String("hello"), String("WORLD")},
			expected: "HELLO and world",
		},
		{
			name:     "Zero positional index {0}",
			format:   "Value: {0}",
			args:     []any{"first"},
			expected: "Value: {0}",
		},
		{
			name:     "Negative positional index {-1}",
			format:   "Value: {-1}",
			args:     []any{"first"},
			expected: "Value: {-1}",
		},
		{
			name:     "Positional index out of range",
			format:   "Value: {5}",
			args:     []any{"only one"},
			expected: "Value: {5}",
		},
		{
			name:     "Nil positional argument",
			format:   "Got: {}",
			args:     []any{nil},
			expected: "Got: <nil>",
		},
		{
			name:     "Nil positional argument with index",
			format:   "Got: {1}",
			args:     []any{nil},
			expected: "Got: <nil>",
		},
		{
			name:     "Mixed positional and named",
			format:   "{1} is {name}",
			args:     []any{"Alice", Named{"name": "great"}},
			expected: "Alice is great",
		},
		{
			name:     "Empty template",
			format:   "",
			args:     []any{"ignored"},
			expected: "",
		},
		{
			name:     "Escaped closing brace",
			format:   "literal \\}",
			args:     nil,
			expected: "literal }",
		},
		{
			name:     "Lone opening brace without closing",
			format:   "broken { here",
			args:     nil,
			expected: "broken { here",
		},
		{
			name:     "Multiple lone opening braces",
			format:   "a { b { c",
			args:     nil,
			expected: "a { b { c",
		},
		{
			name:     "Spec type",
			format:   "Type: {1:T}",
			args:     []any{Int(42)},
			expected: "Type: g.Int",
		},
		{
			name:     "Spec debug on string",
			format:   "Debug: {1:?}",
			args:     []any{String("hi")},
			expected: `Debug: "hi"`,
		},
		{
			name:     "Spec debug on int",
			format:   "Debug: {1:?}",
			args:     []any{Int(7)},
			expected: "Debug: 7",
		},
		{
			name:     "Method with missing params (no panic)",
			format:   "Result: {1.Replace}",
			args:     []any{String("hello")},
			expected: "Result: hello",
		},
		{
			name:     "Method with insufficient params (no panic)",
			format:   "Result: {1.Replace(a)}",
			args:     []any{String("hello")},
			expected: "Result: hello",
		},
		{
			name:     "Nil pointer argument with modifier",
			format:   "Value: {1.Something}",
			args:     []any{(*string)(nil)},
			expected: "Value: <nil>",
		},
		{
			name:     "Adjacent placeholders",
			format:   "{1}{2}{3}",
			args:     []any{"a", "b", "c"},
			expected: "abc",
		},
		{
			name:     "Placeholder with only spaces",
			format:   "Hello, {   }!",
			args:     []any{"world"},
			expected: "Hello, world!",
		},
		{
			name:     "Repeated same index",
			format:   "{1} {1} {1}",
			args:     []any{"echo"},
			expected: "echo echo echo",
		},
		{
			name:     "Backslash not before brace",
			format:   "path\\nhere",
			args:     nil,
			expected: "path\\nhere",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.format, tt.args...)
			if string(result) != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     Named
		expected string
	}{
		// Basic placeholder replacement
		{
			name:     "Basic replacement",
			format:   "Hello, {name}!",
			args:     Named{"name": "John"},
			expected: "Hello, John!",
		},
		// Placeholder with fallback
		{
			name:     "Fallback replacement",
			format:   "Hello, {name?fallback}!",
			args:     Named{"fallback": "Guest"},
			expected: "Hello, Guest!",
		},
		// Placeholder with modifier: upper
		{
			name:     "Modifier: upper",
			format:   "Name: {name.Upper}",
			args:     Named{"name": String("john")},
			expected: "Name: JOHN",
		},
		// Placeholder with modifier: trim and title
		{
			name:     "Modifier: trim and title",
			format:   "Title: {work.Trim.Title}",
			args:     Named{"work": String(" developer ")},
			expected: "Title: Developer",
		},
		// Nested modifiers: trim and len
		{
			name:     "Nested modifiers",
			format:   "Length: {input.Trim.Len}",
			args:     Named{"input": String("  data  ")},
			expected: "Length: 4",
		},
		// Placeholder with fallback and modifier
		{
			name:     "Fallback with modifier",
			format:   "Name: {name?fallback.Upper}",
			args:     Named{"fallback": String("guest")},
			expected: "Name: GUEST",
		},
		// Multiple placeholders
		{
			name:     "Multiple placeholders",
			format:   "{greeting}, {name}! You are {age} years old.",
			args:     Named{"greeting": "Hello", "name": "John", "age": 30},
			expected: "Hello, John! You are 30 years old.",
		},
		// Placeholder with unknown key
		{
			name:     "Unknown placeholder",
			format:   "Hello, {unknown}!",
			args:     Named{"name": "John"},
			expected: "Hello, {unknown}!",
		},
		// Modifier: round for float values
		{
			name:     "Modifier: round",
			format:   "Value: {number.Round}",
			args:     Named{"number": Float(12.7)},
			expected: "Value: 13",
		},
		// Modifier: abs for negative numbers
		{
			name:     "Modifier: abs",
			format:   "Absolute: {value.Abs}",
			args:     Named{"value": Int(-42)},
			expected: "Absolute: 42",
		},
		// Modifier: reverse for strings
		{
			name:     "Modifier: reverse",
			format:   "Reversed: {word.Reverse}",
			args:     Named{"word": String("hello")},
			expected: "Reversed: olleh",
		},
		// Modifier: hex for integers
		{
			name:     "Modifier: hex",
			format:   "Hex: {number.Hex}",
			args:     Named{"number": Int(255)},
			expected: "Hex: ff",
		},
		// Modifier: bin for integers
		{
			name:     "Modifier: bin",
			format:   "Binary: {number.Binary}",
			args:     Named{"number": Int(5)},
			expected: "Binary: 00000101",
		},
		// Modifier: url encoding
		{
			name:     "Modifier: url",
			format:   "URL: {input.Encode.URL}",
			args:     Named{"input": String("hello world")},
			expected: "URL: hello+world",
		},
		// Modifier: base64 encoding
		{
			name:     "Modifier: base64",
			format:   "Base64: {input.Encode.Base64}",
			args:     Named{"input": String("hello")},
			expected: "Base64: aGVsbG8=",
		},
		// Modifier: format for dates
		{
			name:     "Modifier: format date",
			format:   "Date: {today.Format(2006-01-02)}",
			args:     Named{"today": time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC)},
			expected: "Date: 2025-01-17",
		},
		// Test for $replace
		{
			name:     "Modifier: replace",
			format:   "Result: {input.Replace(a,b,-1)}",
			args:     Named{"input": String("banana")},
			expected: "Result: bbnbnb",
		},
		{
			name:     "Modifier: replace with empty string",
			format:   "Result: {input.ReplaceAll(a,)}",
			args:     Named{"input": String("banana")},
			expected: "Result: bnn",
		},
		{
			name:     "Modifier: replace no matches",
			format:   "Result: {input.ReplaceAll(x,y)}",
			args:     Named{"input": String("banana")},
			expected: "Result: banana",
		},
		// Test for $repeat
		{
			name:     "Modifier: repeat string",
			format:   "Repeated: {input.Repeat(3)}",
			args:     Named{"input": String("ha")},
			expected: "Repeated: hahaha",
		},
		{
			name:     "Modifier: repeat with invalid count",
			format:   "Repeated: {input.Repeat(abc)}",
			args:     Named{"input": String("ha")},
			expected: "Repeated: ha",
		},
		// Test for $substring
		{
			name:     "Modifier: substring",
			format:   "Result: {input.SubString(0,-1,2)}",
			args:     Named{"input": String("Hello, World!")},
			expected: "Result: Hlo ol",
		},
		// Test for $truncate
		{
			name:     "Modifier: truncate string",
			format:   "Truncated: {input.Truncate(5)}",
			args:     Named{"input": String("Hello, World!")},
			expected: "Truncated: Hello...",
		},
		{
			name:     "Modifier: truncate with exact length",
			format:   "Truncated: {input.Truncate(5)}",
			args:     Named{"input": String("Hello")},
			expected: "Truncated: Hello",
		},
		{
			name:     "Modifier: truncate with no truncation",
			format:   "Truncated: {input.Truncate(15)}",
			args:     Named{"input": String("Hello, World!")},
			expected: "Truncated: Hello, World!",
		},
		{
			name:     "Modifier: truncate with invalid max",
			format:   "Truncated: {input.Truncate(abc)}",
			args:     Named{"input": String("Hello, World!")},
			expected: "Truncated: Hello, World!",
		},
		// A format string with no placeholders at all.
		{
			name:     "No placeholders",
			format:   "Just a normal text",
			args:     make(Named),
			expected: "Just a normal text",
		},
		// An empty placeholder (e.g., "Hello, {}!")
		{
			name:     "Empty placeholder",
			format:   "Hello, {}!",
			args:     make(Named),
			expected: "Hello, {}!",
		},
		//  Multiple chained modifiers (e.g., trim, lower, replace, reverse).
		{
			name:   "Multiple chain modifiers",
			format: "{word.Trim.Lower.ReplaceAll(e,a).Reverse}",
			args:   Named{"word": String("  EXAMPLE ")},
			// Explanation:
			//   "  EXAMPLE " -> $trim => "EXAMPLE"
			//   -> $lower => "example"
			//   -> $replace(e,a) => "axampla"
			//   -> $reverse => "alpmaxa"
			expected: "alpmaxa",
		},
		{
			name:     "Named nil value",
			format:   "Value: {key}",
			args:     Named{"key": nil},
			expected: "Value: {key}",
		},
		{
			name:     "Both key and fallback missing",
			format:   "Hello, {x?y}!",
			args:     Named{"z": "nope"},
			expected: "Hello, {x?y}!",
		},
		{
			name:     "Spec type on named",
			format:   "Type: {val:T}",
			args:     Named{"val": Int(10)},
			expected: "Type: g.Int",
		},
		{
			name:     "Spec debug on named",
			format:   "Debug: {val:?}",
			args:     Named{"val": String("test")},
			expected: `Debug: "test"`,
		},
		{
			name:     "Empty named map",
			format:   "Hi {name}",
			args:     Named{},
			expected: "Hi {name}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.format, tt.args)
			if result != String(tt.expected) {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestFormatFormatWithErrors(t *testing.T) {
	errorTests := []struct {
		name     string
		format   string
		args     Named
		expected string
	}{
		// Placeholder with invalid syntax
		{
			name:     "Invalid placeholder syntax",
			format:   "Hello, {name?",
			args:     Named{"name": "John"},
			expected: "Hello, {name?",
		},
		// Modifier with invalid syntax
		{
			name:     "Invalid modifier syntax",
			format:   "Value: {number.Unknown(",
			args:     Named{"number": 42},
			expected: "Value: {number.Unknown(",
		},
		// Unsupported modifier
		{
			name:     "Unsupported modifier",
			format:   "Value: {number.Unsupported}",
			args:     Named{"number": 42},
			expected: "Value: 42",
		},
		// Fallback key missing
		{
			name:     "Missing fallback key",
			format:   "Hello, {name?fallback}!",
			args:     Named{},
			expected: "Hello, {name?fallback}!",
		},
		// Placeholder with unsupported type
		{
			name:     "Unsupported type for modifier",
			format:   "Value: {obj.Upper}",
			args:     Named{"obj": Unit{}},
			expected: "Value: {}",
		},
		{
			name:     "Nested unclosed brace",
			format:   "Hello {name.Upper",
			args:     Named{"name": String("john")},
			expected: "Hello {name.Upper",
		},
		{
			name:     "Empty brace pair in named context",
			format:   "A {} B",
			args:     Named{"x": "y"},
			expected: "A {} B",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.format, tt.args)
			if result != String(tt.expected) {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestFormatTrimSetModifier(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     Named
		expected string
	}{
		// Basic trimming
		{
			name:   "Trim specific characters",
			format: "Result: {value.TrimSet(#)}",
			args:   Named{"value": String("###Hello###")}, expected: "Result: Hello",
		},
		// Trim multiple characters
		{
			name:     "Trim multiple characters",
			format:   "Result: {value.TrimSet(#$)}",
			args:     Named{"value": String("$$#Hello#$")},
			expected: "Result: Hello",
		},
		// No trimming (no matching characters)
		{
			name:     "No trimming needed",
			format:   "Result: {value.TrimSet(%)}",
			args:     Named{"value": String("Hello")},
			expected: "Result: Hello",
		},
		// Empty value
		{
			name:     "Empty value",
			format:   "Result: {value.TrimSet(#)}",
			args:     Named{"value": String("")},
			expected: "Result: ",
		},
		// Empty set
		{
			name:     "Empty trim set",
			format:   "Result: {value.Trim}",
			args:     Named{"value": String("###Hello###")},
			expected: "Result: ###Hello###",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.format, tt.args)
			if result != String(tt.expected) {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// go test -bench=. -benchmem -count=4

func BenchmarkSprintf(b *testing.B) {
	name := "World"

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = fmt.Sprintf("Hello, %s!", name)
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = Format("Hello, {}!", name)
		}
	})
}

func BenchmarkSprintfPositional(b *testing.B) {
	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = fmt.Sprintf("%[2]s comes before %[1]s", "World", "Hello")
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = Format("{2} comes before {1}", "World", "Hello")
		}
	})
}

func BenchmarkSprintfNamedAccess(b *testing.B) {
	data := Named{"email": "alice@example.com"}

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = fmt.Sprintf("Email: %s", data["email"])
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = Format("Email: {email}", data)
		}
	})
}

func BenchmarkSprintfFormatSpecifiers(b *testing.B) {
	num := Int(255)

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = fmt.Sprintf("Hex: %x, Binary: %b", num, num)
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = Format("Hex: {1.Hex}, Binary: {1.Binary}", num)
		}
	})
}

func TestFormatGet(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []any
		expected string
	}{
		{
			name:     "Simple Map Access",
			template: "Value: {1.Get(key).Some}",
			args: []any{
				Map[String, String]{"key": "value"},
			},
			expected: "Value: value",
		},
		{
			name:     "Simple Map Any Access",
			template: "Value: {1.Get(key).Some}",
			args: []any{
				Map[String, any]{"key": "value"},
			},
			expected: "Value: value",
		},
		{
			name:     "Nested Map Access",
			template: "Deep Value: {1.Get(key).Some.Get(subkey).Some}",
			args: []any{
				Map[string, Map[String, string]]{
					"key": {"subkey": "deepvalue"},
				},
			},
			expected: "Deep Value: deepvalue",
		},
		{
			name:     "Map with Float Keys",
			template: "Float Key: {1.Get(3_14).Some}",
			args: []any{
				Map[Float, string]{3.14: "pi"},
			},
			expected: "Float Key: pi",
		},
		{
			name:     "Slice Index Access",
			template: "Index 1: {1.Get(1).Some}",
			args: []any{
				Slice[string]{"first", "second", "third"},
			},
			expected: "Index 1: second",
		},
		{
			name:     "Nested Slice Access",
			template: "Nested Index: {1.Get(1).Some.Get(0).Some}",
			args: []any{
				Slice[Slice[Int]]{{100, 200}, {300, 400}},
			},
			expected: "Nested Index: 300",
		},
		{
			name:     "Map with Int Keys",
			template: "Int Key: {1.Get(42).Some}",
			args: []any{
				Map[int, string]{42: "intvalue"},
			},
			expected: "Int Key: intvalue",
		},
		{
			name:     "Boolean Key Map",
			template: "Bool Key: {1.Get(true).Some}",
			args: []any{
				Map[bool, string]{true: "boolvalue"},
			},
			expected: "Bool Key: boolvalue",
		},
		{
			name:     "Full Complexity",
			template: "Access: {1.Get(map).Some.Get(slice).Some.Get(1).Some.Get(0).Some.Get(field).Some}",
			args: []any{
				Map[String, Map[string, Map[String, Slice[Map[string, string]]]]]{
					"map": {
						"slice": {
							"1": Slice[Map[string, string]]{
								{"field": "subfieldvalue"},
							},
						},
					},
				},
			},
			expected: "Access: subfieldvalue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.template, tt.args...)
			if result != String(tt.expected) {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestFormatGetNamed(t *testing.T) {
	tests := []struct {
		name     string
		template string
		named    Named
		expected string
	}{
		{
			name:     "Simple Map Access",
			template: "Value: {map.Get(key).Some}",
			named:    Named{"map": Map[String, String]{"key": "value"}},
			expected: "Value: value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.template, tt.named)
			if result != String(tt.expected) {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestFormatMapSliceStruct(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []any
		expected string
	}{
		{
			name:     "Simple Map Access",
			template: "Value: {1.key}",
			args: []any{
				map[String]String{"key": "value"},
			},
			expected: "Value: value",
		},
		{
			name:     "Simple Map Any Access",
			template: "Value: {1.key}",
			args: []any{
				map[String]any{"key": "value"},
			},
			expected: "Value: value",
		},
		{
			name:     "Nested Map Access",
			template: "Deep Value: {1.key.subkey}",
			args: []any{
				map[string]map[String]string{
					"key": {"subkey": "deepvalue"},
				},
			},
			expected: "Deep Value: deepvalue",
		},
		{
			name:     "Map with Float Keys",
			template: "Float Key: {1.3_14}",
			args: []any{
				map[Float]string{3.14: "pi"},
			},
			expected: "Float Key: pi",
		},
		{
			name:     "Slice Index Access",
			template: "Index 1: {1.1}",
			args: []any{
				[]string{"first", "second", "third"},
			},
			expected: "Index 1: second",
		},
		{
			name:     "Nested Slice Access",
			template: "Nested Index: {1.1.0}",
			args: []any{
				[][]Int{{100, 200}, {300, 400}},
			},
			expected: "Nested Index: 300",
		},
		{
			name:     "Struct Field Access",
			template: "Struct Field: {1.Field}",
			args: []any{
				struct {
					Field string
				}{Field: "fieldvalue"},
			},
			expected: "Struct Field: fieldvalue",
		},
		{
			name:     "Complex Struct Field Access",
			template: "Complex Struct: {1.SubStruct.InnerField}",
			args: []any{
				struct {
					SubStruct struct {
						InnerField String
					}
				}{SubStruct: struct {
					InnerField String
				}{InnerField: "inner"}},
			},
			expected: "Complex Struct: inner",
		},
		{
			name:     "Map with Int Keys",
			template: "Int Key: {1.42}",
			args: []any{
				map[int]string{42: "intvalue"},
			},
			expected: "Int Key: intvalue",
		},
		{
			name:     "Boolean Key Map",
			template: "Bool Key: {1.true}",
			args: []any{
				map[bool]string{true: "boolvalue"},
			},
			expected: "Bool Key: boolvalue",
		},
		{
			name:     "Full Complexity",
			template: "Access: {1.map.slice.1.0.field}",
			args: []any{
				map[String]map[string]map[String][]map[string]string{
					"map": {
						"slice": {
							"1": []map[string]string{
								{"field": "subfieldvalue"},
							},
						},
					},
				},
			},
			expected: "Access: subfieldvalue",
		},
		{
			name:     "Slice index out of range",
			template: "Value: {1.10}",
			args: []any{
				[]string{"only", "two"},
			},
			expected: "Value: [only two]",
		},
		{
			name:     "Struct missing field",
			template: "Value: {1.NonExistent}",
			args: []any{
				struct{ Name string }{Name: "test"},
			},
			expected: "Value: {test}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.template, tt.args...)
			if result != String(tt.expected) {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestFormatComplex(t *testing.T) {
	tests := []struct {
		name     string
		template string
		named    Named
		expected string
	}{
		{
			name:     "Simple Map Access",
			template: "Value: {map.key}",
			named:    Named{"map": map[String]String{"key": "value"}},
			expected: "Value: value",
		},
		{
			name:     "Full Complexity",
			template: "Access: {complex.map.slice.1.0.field} {struct.SubStruct.InnerField}",
			named: Named{
				"complex": map[String]map[string]map[String][]map[string]string{
					"map": {
						"slice": {
							"1": []map[string]string{
								{"field": "subfieldvalue"},
							},
						},
					},
				},
				"struct": struct {
					SubStruct struct {
						InnerField String
					}
				}{SubStruct: struct {
					InnerField String
				}{InnerField: "inner"}},
			},
			expected: "Access: subfieldvalue inner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.template, tt.named)
			if result != String(tt.expected) {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	var sb strings.Builder

	res := Write(&sb, "test message") // pass *strings.Builder
	if res.IsErr() {
		t.Fatalf("Write error: %v", res.Err())
	}

	if n := res.UnwrapOr(0); n != len("test message") {
		t.Fatalf("bytes written mismatch: got %d", n)
	}

	if got := sb.String(); got != "test message" {
		t.Fatalf("content mismatch: %q", got)
	}
}

func TestWriteWithPlaceholders(t *testing.T) {
	var sb strings.Builder

	res := Write(&sb, "Hello, {}!", "world")
	if res.IsErr() {
		t.Fatalf("Write error: %v", res.Err())
	}

	if got := sb.String(); got != "Hello, world!" {
		t.Fatalf("content mismatch: %q", got)
	}
}

func TestWriteln(t *testing.T) {
	var sb strings.Builder

	res := Writeln(&sb, "test message")
	if res.IsErr() {
		t.Fatalf("Writeln error: %v", res.Err())
	}

	if n := res.UnwrapOr(0); n != len("test message\n") {
		t.Fatalf("bytes written mismatch: got %d", n)
	}

	if got := sb.String(); got != "test message\n" {
		t.Fatalf("content mismatch: %q", got)
	}
}

func TestWritelnWithPlaceholders(t *testing.T) {
	var sb strings.Builder

	res := Writeln(&sb, "Hello, {}", "world")
	if res.IsErr() {
		t.Fatalf("Writeln error: %v", res.Err())
	}

	if got := sb.String(); got != "Hello, world\n" {
		t.Fatalf("content mismatch: %q", got)
	}
}

func TestPrint(t *testing.T) {
	// Test that Print doesn't crash - we can't easily capture stdout in unit tests
	Print("test print")
}

func TestPrintln(t *testing.T) {
	// Test that Println doesn't crash - we can't easily capture stdout in unit tests
	Println("test println")
}

func TestEprint(t *testing.T) {
	// Test that Eprint doesn't crash - we can't easily capture stderr in unit tests
	Eprint("test eprint")
}

func TestEprintln(t *testing.T) {
	// Test that Eprintln doesn't crash - we can't easily capture stderr in unit tests
	Eprintln("test eprintln")
}

func TestErrorf(t *testing.T) {
	sentinel := errors.New("sentinel")

	t.Run("no wrap — plain error", func(t *testing.T) {
		err := Errorf("something went wrong: {}", sentinel)
		if err.Error() != "something went wrong: sentinel" {
			t.Fatalf("unexpected message: %s", err)
		}
		if errors.Is(err, sentinel) {
			t.Fatal("expected no wrapping")
		}
	})

	t.Run("positional wrap {1:w}", func(t *testing.T) {
		err := Errorf("failed: {1:w}", sentinel)
		if err.Error() != "failed: sentinel" {
			t.Fatalf("unexpected message: %s", err)
		}
		if !errors.Is(err, sentinel) {
			t.Fatal("expected errors.Is to find sentinel")
		}
	})

	t.Run("auto-index wrap {:w}", func(t *testing.T) {
		err := Errorf("open {}: {:w}", "foo.txt", sentinel)
		if err.Error() != "open foo.txt: sentinel" {
			t.Fatalf("unexpected message: %s", err)
		}
		if !errors.Is(err, sentinel) {
			t.Fatal("expected errors.Is to find sentinel")
		}
	})

	t.Run("named wrap {cause:w}", func(t *testing.T) {
		err := Errorf("read {file}: {cause:w}", Named{"file": "bar.txt", "cause": sentinel})
		if err.Error() != "read bar.txt: sentinel" {
			t.Fatalf("unexpected message: %s", err)
		}
		if !errors.Is(err, sentinel) {
			t.Fatal("expected errors.Is to find sentinel")
		}
	})

	t.Run("multiple wraps", func(t *testing.T) {
		e1 := errors.New("e1")
		e2 := errors.New("e2")
		err := Errorf("{1:w} and {2:w}", e1, e2)
		if !errors.Is(err, e1) {
			t.Fatal("expected errors.Is to find e1")
		}
		if !errors.Is(err, e2) {
			t.Fatal("expected errors.Is to find e2")
		}
	})

	t.Run("wrap with positional spec", func(t *testing.T) {
		err := Errorf("msg: {1:w}", sentinel)
		if !errors.Is(err, sentinel) {
			t.Fatal("expected wrapping")
		}
	})

	t.Run("wrap on non-error value — no panic", func(t *testing.T) {
		err := Errorf("value: {1:w}", "not an error")
		if err.Error() != "value: not an error" {
			t.Fatalf("unexpected message: %s", err)
		}
		// no wrapping since value is not an error
	})

	t.Run("nil argument in Errorf", func(t *testing.T) {
		err := Errorf("value is {}", nil)
		if err.Error() != "value is <nil>" {
			t.Fatalf("unexpected message: %s", err)
		}
	})

	t.Run("no placeholders", func(t *testing.T) {
		err := Errorf("static error message")
		if err.Error() != "static error message" {
			t.Fatalf("unexpected message: %s", err)
		}
	})

	t.Run("wrap nil error — no panic", func(t *testing.T) {
		// nil is converted to "<nil>" string, so wrap should not match error interface
		err := Errorf("err: {1:w}", nil)
		if err.Error() != "err: <nil>" {
			t.Fatalf("unexpected message: %s", err)
		}
	})
}

func TestFormatMapOrdAccess(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []any
		expected string
	}{
		{
			name:     "MapOrd Any String Key Access",
			template: "Value: {1.stringkey}",
			args: func() []any {
				m := NewMapOrd[string, string]()
				m.Insert("stringkey", "stringvalue")
				return []any{m}
			}(),
			expected: "Value: stringvalue",
		},
		{
			name:     "MapOrd Any String Type Key Access",
			template: "Value: {1.stringkey}",
			args: func() []any {
				m := NewMapOrd[String, string]()
				m.Insert(String("stringkey"), "stringvalue")
				return []any{m}
			}(),
			expected: "Value: stringvalue",
		},
		{
			name:     "MapOrd Any Int Key Access",
			template: "Value: {1.42}",
			args: func() []any {
				m := NewMapOrd[int, string]()
				m.Insert(42, "intvalue")
				return []any{m}
			}(),
			expected: "Value: intvalue",
		},
		{
			name:     "MapOrd Any Int Type Key Access",
			template: "Value: {1.42}",
			args: func() []any {
				m := NewMapOrd[Int, string]()
				m.Insert(Int(42), "intvalue")
				return []any{m}
			}(),
			expected: "Value: intvalue",
		},
		{
			name:     "MapOrd Any Float Key Access",
			template: "Value: {1.3}",
			args: func() []any {
				m := NewMapOrd[float64, string]()
				m.Insert(3.0, "floatvalue")
				return []any{m}
			}(),
			expected: "Value: floatvalue",
		},
		{
			name:     "MapOrd Any Float Type Key Access",
			template: "Value: {1.3}",
			args: func() []any {
				m := NewMapOrd[Float, string]()
				m.Insert(Float(3.0), "floatvalue")
				return []any{m}
			}(),
			expected: "Value: floatvalue",
		},
		{
			name:     "MapOrd Any Missing Key (returns whole map)",
			template: "Value: {1.missing}",
			args: func() []any {
				m := NewMapOrd[string, string]()
				m.Insert("existing", "value")
				return []any{m}
			}(),
			expected: "Value: MapOrd{existing:value}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.template, tt.args...)
			if result != String(tt.expected) {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestFormatSpec(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		// ── Verb: x (hex lowercase) ──
		{name: "x int", format: "{:x}", args: []any{255}, expected: "ff"},
		{name: "x zero", format: "{:x}", args: []any{0}, expected: "0"},
		{name: "x negative", format: "{:x}", args: []any{-255}, expected: "-ff"},
		{name: "x uint", format: "{:x}", args: []any{uint(255)}, expected: "ff"},
		{name: "x uint8", format: "{:x}", args: []any{uint8(255)}, expected: "ff"},
		{name: "x uint16", format: "{:x}", args: []any{uint16(1024)}, expected: "400"},
		{name: "x uint32", format: "{:x}", args: []any{uint32(65535)}, expected: "ffff"},
		{name: "x uint64", format: "{:x}", args: []any{uint64(0xdeadbeef)}, expected: "deadbeef"},
		{name: "x int8", format: "{:x}", args: []any{int8(127)}, expected: "7f"},
		{name: "x int16", format: "{:x}", args: []any{int16(256)}, expected: "100"},
		{name: "x int32", format: "{:x}", args: []any{int32(255)}, expected: "ff"},
		{name: "x int64", format: "{:x}", args: []any{int64(4096)}, expected: "1000"},
		{name: "x g.Int", format: "{:x}", args: []any{Int(255)}, expected: "ff"},
		{name: "x g.Int negative", format: "{:x}", args: []any{Int(-1)}, expected: "-1"},
		{name: "x non-numeric", format: "{:x}", args: []any{"hello"}, expected: "hello"},

		// ── Verb: X (hex uppercase) ──
		{name: "X int", format: "{:X}", args: []any{255}, expected: "FF"},
		{name: "X uint64", format: "{:X}", args: []any{uint64(0xabcdef)}, expected: "ABCDEF"},

		// ── Verb: o (octal) ──
		{name: "o int", format: "{:o}", args: []any{255}, expected: "377"},
		{name: "o zero", format: "{:o}", args: []any{0}, expected: "0"},
		{name: "o negative", format: "{:o}", args: []any{-8}, expected: "-10"},
		{name: "o uint", format: "{:o}", args: []any{uint(8)}, expected: "10"},
		{name: "o g.Int", format: "{:o}", args: []any{Int(255)}, expected: "377"},
		{name: "o non-numeric", format: "{:o}", args: []any{"hello"}, expected: "hello"},

		// ── Verb: b (binary) ──
		{name: "b int", format: "{:b}", args: []any{42}, expected: "101010"},
		{name: "b zero", format: "{:b}", args: []any{0}, expected: "0"},
		{name: "b negative", format: "{:b}", args: []any{-5}, expected: "-101"},
		{name: "b uint", format: "{:b}", args: []any{uint(255)}, expected: "11111111"},
		{name: "b g.Int", format: "{:b}", args: []any{Int(42)}, expected: "101010"},
		{name: "b non-numeric", format: "{:b}", args: []any{"hello"}, expected: "hello"},

		// ── Verb: e/E (exponential) ──
		{name: "e float64", format: "{:e}", args: []any{3.14}, expected: "3.140000e+00"},
		{name: "E float64", format: "{:E}", args: []any{3.14}, expected: "3.140000E+00"},
		{name: "e float32", format: "{:e}", args: []any{float32(3.14)}, expected: "3.140000e+00"},
		{name: "e g.Float", format: "{:e}", args: []any{Float(3.14)}, expected: "3.140000e+00"},
		{name: "e int", format: "{:e}", args: []any{100}, expected: "1.000000e+02"},
		{name: "e precision", format: "{:.2e}", args: []any{3.14159}, expected: "3.14e+00"},
		{name: "E precision", format: "{:.2E}", args: []any{3.14159}, expected: "3.14E+00"},
		{name: "e precision 0", format: "{:.0e}", args: []any{3.14}, expected: "3e+00"},
		{name: "e negative", format: "{:e}", args: []any{-3.14}, expected: "-3.140000e+00"},
		{name: "e with sign", format: "{:+e}", args: []any{3.14}, expected: "+3.140000e+00"},
		{name: "e non-numeric", format: "{:e}", args: []any{"hello"}, expected: "hello"},

		// ── Verb: ? (debug) ──
		{name: "debug string", format: "{:?}", args: []any{"hello"}, expected: `"hello"`},
		{name: "debug int", format: "{:?}", args: []any{42}, expected: "42"},
		{name: "debug nil", format: "{:?}", args: []any{(*int)(nil)}, expected: "(*int)(nil)"},
		{name: "debug slice", format: "{:?}", args: []any{[]int{1, 2, 3}}, expected: "[]int{1, 2, 3}"},
		{name: "debug map", format: "{:?}", args: []any{map[string]int{"a": 1}}, expected: `map[string]int{"a":1}`},

		// ── Verb: #? (pretty debug) ──
		{
			name:     "pretty debug",
			format:   "{:#?}",
			args:     []any{struct{ X int }{42}},
			expected: "struct {\n   X int \n}{\n  X:42\n}",
		},
		{
			name:     "pretty debug strings",
			format:   "{:#?}",
			args:     []any{struct{ A, B string }{"hi", "bye"}},
			expected: "struct {\n   A string; B string \n}{\n  A:\"hi\",\n   B:\"bye\"\n}",
		},

		// ── Verb: T (type) ──
		{name: "T int", format: "{:T}", args: []any{42}, expected: "int"},
		{name: "T string", format: "{:T}", args: []any{"hello"}, expected: "string"},
		{name: "T g.Int", format: "{:T}", args: []any{Int(1)}, expected: "g.Int"},
		{name: "T g.String", format: "{:T}", args: []any{String("s")}, expected: "g.String"},
		{name: "T g.Float", format: "{:T}", args: []any{Float(1)}, expected: "g.Float"},
		{name: "T slice", format: "{:T}", args: []any{[]int{1}}, expected: "[]int"},

		// ── Verb: p (pointer) ──
		{name: "pointer", format: "{:p}", args: []any{&struct{}{}}},

		// ── Alternate form (#) ──
		{name: "#x", format: "{:#x}", args: []any{255}, expected: "0xff"},
		{name: "#X", format: "{:#X}", args: []any{255}, expected: "0xFF"},
		{name: "#o", format: "{:#o}", args: []any{255}, expected: "0o377"},
		{name: "#b", format: "{:#b}", args: []any{42}, expected: "0b101010"},
		{name: "#x zero", format: "{:#x}", args: []any{0}, expected: "0x0"},
		{name: "#x negative", format: "{:#x}", args: []any{-255}, expected: "-0xff"},
		{name: "#x uint", format: "{:#x}", args: []any{uint(255)}, expected: "0xff"},
		{name: "#x g.Int", format: "{:#x}", args: []any{Int(255)}, expected: "0xff"},

		// ── Alignment: > (right) ──
		{name: "> string", format: "{:>10}", args: []any{"hello"}, expected: "     hello"},
		{name: "> int", format: "{:>10}", args: []any{42}, expected: "        42"},
		{name: "> exact width", format: "{:>5}", args: []any{"hello"}, expected: "hello"},
		{name: "> content wider", format: "{:>3}", args: []any{"hello"}, expected: "hello"},

		// ── Alignment: < (left) ──
		{name: "< string", format: "{:<10}", args: []any{"hello"}, expected: "hello     "},
		{name: "< int", format: "{:<10}", args: []any{42}, expected: "42        "},

		// ── Alignment: ^ (center) ──
		{name: "^ string", format: "{:^10}", args: []any{"hello"}, expected: "  hello   "},
		{name: "^ even", format: "{:^6}", args: []any{"ab"}, expected: "  ab  "},
		{name: "^ odd", format: "{:^7}", args: []any{"ab"}, expected: "  ab   "},

		// ── Fill + alignment ──
		{name: "fill - center", format: "{:-^20}", args: []any{"hello"}, expected: "-------hello--------"},
		{name: "fill * right", format: "{:*>10}", args: []any{"hi"}, expected: "********hi"},
		{name: "fill . left", format: "{:.<10}", args: []any{"hi"}, expected: "hi........"},
		{name: "fill = center", format: "{:=^10}", args: []any{"hi"}, expected: "====hi===="},
		{name: "fill 0 left", format: "{:0<10}", args: []any{"hi"}, expected: "hi00000000"},

		// ── Default alignment (no explicit align) ──
		{name: "default string left", format: "{:10}", args: []any{"hi"}, expected: "hi        "},
		{name: "default int right", format: "{:10}", args: []any{42}, expected: "        42"},
		{name: "default float right", format: "{:10}", args: []any{3.14}, expected: "      3.14"},
		{name: "default g.Int right", format: "{:10}", args: []any{Int(42)}, expected: "        42"},

		// ── Precision ──
		{name: "prec float .2", format: "{:.2}", args: []any{3.14159}, expected: "3.14"},
		{name: "prec float .0", format: "{:.0}", args: []any{3.14159}, expected: "3"},
		{name: "prec float .4", format: "{:.4}", args: []any{3.14159}, expected: "3.1416"},
		{name: "prec float .10", format: "{:.10}", args: []any{3.14}, expected: "3.1400000000"},
		{name: "prec g.Float .2", format: "{:.2}", args: []any{Float(3.14159)}, expected: "3.14"},
		{name: "prec float32 .2", format: "{:.2}", args: []any{float32(3.14)}, expected: "3.14"},
		{name: "prec negative float", format: "{:.2}", args: []any{-3.14159}, expected: "-3.14"},
		{name: "prec string truncate", format: "{:.3}", args: []any{"hello"}, expected: "hel"},
		{name: "prec string no trunc", format: "{:.10}", args: []any{"hello"}, expected: "hello"},
		{name: "prec string exact", format: "{:.5}", args: []any{"hello"}, expected: "hello"},
		{name: "prec string zero", format: "{:.0}", args: []any{"hello"}, expected: ""},
		{name: "prec unicode trunc", format: "{:.3}", args: []any{"привет"}, expected: "при"},

		// ── Sign ──
		{name: "sign + positive int", format: "{:+}", args: []any{42}, expected: "+42"},
		{name: "sign + negative int", format: "{:+}", args: []any{-42}, expected: "-42"},
		{name: "sign + zero", format: "{:+}", args: []any{0}, expected: "+0"},
		{name: "sign + positive float", format: "{:+.2}", args: []any{3.14}, expected: "+3.14"},
		{name: "sign + negative float", format: "{:+.2}", args: []any{-3.14}, expected: "-3.14"},
		{name: "sign + g.Int", format: "{:+}", args: []any{Int(42)}, expected: "+42"},
		{name: "sign + g.Float", format: "{:+.1}", args: []any{Float(3.14)}, expected: "+3.1"},
		{name: "sign space positive", format: "{: }", args: []any{42}, expected: " 42"},
		{name: "sign space negative", format: "{: }", args: []any{-42}, expected: "-42"},
		{name: "sign space float", format: "{: .2}", args: []any{3.14}, expected: " 3.14"},
		{name: "sign + hex", format: "{:+x}", args: []any{255}, expected: "+ff"},
		{name: "sign + octal", format: "{:+o}", args: []any{255}, expected: "+377"},
		{name: "sign + binary", format: "{:+b}", args: []any{42}, expected: "+101010"},

		// ── Zero-padding ──
		{name: "0pad int", format: "{:05}", args: []any{42}, expected: "00042"},
		{name: "0pad negative", format: "{:05}", args: []any{-42}, expected: "-0042"},
		{name: "0pad +sign", format: "{:+05}", args: []any{42}, expected: "+0042"},
		{name: "0pad #x", format: "{:#010x}", args: []any{255}, expected: "0x000000ff"},
		{name: "0pad #X", format: "{:#010X}", args: []any{255}, expected: "0x000000FF"},
		{name: "0pad #b", format: "{:#010b}", args: []any{42}, expected: "0b00101010"},
		{name: "0pad #o", format: "{:#010o}", args: []any{255}, expected: "0o00000377"},
		{name: "0pad binary", format: "{:08b}", args: []any{42}, expected: "00101010"},
		{name: "0pad hex", format: "{:04x}", args: []any{15}, expected: "000f"},
		{name: "0pad wider than needed", format: "{:02}", args: []any{123}, expected: "123"},
		{name: "0pad exact width", format: "{:03}", args: []any{123}, expected: "123"},
		{name: "0pad g.Int", format: "{:05}", args: []any{Int(42)}, expected: "00042"},
		{name: "0pad negative hex", format: "{:08x}", args: []any{-255}, expected: "-00000ff"},

		// ── Width + verb combinations ──
		{name: "width + hex right", format: "{:>10x}", args: []any{255}, expected: "        ff"},
		{name: "width + hex left", format: "{:<10x}", args: []any{255}, expected: "ff        "},
		{name: "width + binary center", format: "{:^10b}", args: []any{42}, expected: "  101010  "},
		{name: "width + octal fill", format: "{:_>10o}", args: []any{255}, expected: "_______377"},
		{name: "width + #x", format: "{:>#10x}", args: []any{255}, expected: "      0xff"},

		// ── Width + precision ──
		{name: "w+p float right", format: "{:>10.2}", args: []any{3.14159}, expected: "      3.14"},
		{name: "w+p float left", format: "{:<10.2}", args: []any{3.14159}, expected: "3.14      "},
		{name: "w+p float center", format: "{:^10.2}", args: []any{3.14159}, expected: "   3.14   "},
		{name: "w+p float fill", format: "{:*>10.2}", args: []any{3.14159}, expected: "******3.14"},
		{name: "w+p string trunc + pad", format: "{:>10.3}", args: []any{"hello"}, expected: "       hel"},
		{name: "w+p sign + prec + width", format: "{:+10.2}", args: []any{3.14}, expected: "     +3.14"},

		// ── g.Int / g.Float all verbs ──
		{name: "g.Int octal", format: "{:o}", args: []any{Int(255)}, expected: "377"},
		{name: "g.Int exponential", format: "{:e}", args: []any{Int(100)}, expected: "1.000000e+02"},
		{name: "g.Int sign", format: "{:+}", args: []any{Int(42)}, expected: "+42"},
		{name: "g.Int #x", format: "{:#x}", args: []any{Int(255)}, expected: "0xff"},
		{name: "g.Float exponential", format: "{:e}", args: []any{Float(3.14)}, expected: "3.140000e+00"},
		{name: "g.Float sign", format: "{:+}", args: []any{Float(3.14)}, expected: "+3.14"},

		// ── Modifier + spec combinations ──
		{name: "mod Abs + 0pad", format: "{1.Abs:05}", args: []any{Int(-42)}, expected: "00042"},
		{name: "mod Upper + right", format: "{.Upper:>10}", args: []any{String("hello")}, expected: "     HELLO"},
		{name: "mod chain + spec", format: "{.Trim.Upper:>10}", args: []any{String("  go  ")}, expected: "        GO"},
		{name: "named + spec", format: "{num:x}", args: []any{Named{"num": 255}}, expected: "ff"},
		{
			name:     "named mod + spec",
			format:   "{name.Upper:>10}",
			args:     []any{Named{"name": String("go")}},
			expected: "        GO",
		},
		{name: "fallback + spec", format: "{x?y:x}", args: []any{Named{"y": 255}}, expected: "ff"},

		// ── Auto-index + spec ──
		{name: "auto x and b", format: "{:x} and {:b}", args: []any{255, 42}, expected: "ff and 101010"},
		{name: "auto multiple", format: "{:x} {:o} {:b}", args: []any{255, 255, 255}, expected: "ff 377 11111111"},
		{name: "auto with plain", format: "{} {:x} {}", args: []any{"a", 255, "b"}, expected: "a ff b"},

		// ── Same arg different specs ──
		{
			name:     "same arg multi spec",
			format:   "{1:x} {1:o} {1:b} {1}",
			args:     []any{255},
			expected: "ff 377 11111111 255",
		},

		// ── Edge cases ──
		{name: "empty spec", format: "{:}", args: []any{"hello"}, expected: "hello"},
		{name: "empty spec int", format: "{:}", args: []any{42}, expected: "42"},
		{
			name:     "colon in parens",
			format:   "{1.Format(15:04:05)}",
			args:     []any{time.Date(2025, 1, 17, 14, 30, 45, 0, time.UTC)},
			expected: "14:30:45",
		},
		{name: "+#x combined", format: "{:+#x}", args: []any{255}, expected: "+0xff"},
		{name: " #x combined", format: "{: #x}", args: []any{255}, expected: " 0xff"},
		{
			name:     "spec on non-existent key",
			format:   "{missing:x}",
			args:     []any{Named{"other": 1}},
			expected: "{missing}",
		},
		{name: "width 1 string", format: "{:1}", args: []any{"hello"}, expected: "hello"},
		{name: "width 0", format: "{:0}", args: []any{"hello"}, expected: "hello"},

		// ── parseFmtSpec: signMinus branch ──
		{name: "sign - positive", format: "{:-}", args: []any{42}, expected: "42"},
		{name: "sign - negative", format: "{:-}", args: []any{-42}, expected: "-42"},

		// ── parseFmtSpec: precision dot with no digits ──
		{name: "prec dot only float", format: "{:.}", args: []any{3.14}, expected: "3"},
		{name: "prec dot only string", format: "{:.}", args: []any{"hello"}, expected: ""},

		// ── fmtToInt64: uint types via sign path ──
		{name: "sign + uint", format: "{:+}", args: []any{uint(42)}, expected: "+42"},
		{name: "sign + uint16", format: "{:+}", args: []any{uint16(42)}, expected: "+42"},
		{name: "sign + uint32", format: "{:+}", args: []any{uint32(42)}, expected: "+42"},
		{name: "sign + uint64", format: "{:+}", args: []any{uint64(42)}, expected: "+42"},
		{
			name:     "sign + uint64 overflow",
			format:   "{:+}",
			args:     []any{uint64(math.MaxUint64)},
			expected: "18446744073709551615",
		},

		// ── fmtToInt64: uint types via exponential path ──
		{name: "e uint8", format: "{:e}", args: []any{uint8(1)}, expected: "1.000000e+00"},
		{name: "e uint16", format: "{:e}", args: []any{uint16(100)}, expected: "1.000000e+02"},
		{name: "e uint32", format: "{:e}", args: []any{uint32(100)}, expected: "1.000000e+02"},

		// ── fmtToInt64: float types via intbase path ──
		{name: "x Float", format: "{:x}", args: []any{Float(255)}, expected: "ff"},
		{name: "x float32", format: "{:x}", args: []any{float32(255)}, expected: "ff"},
		{name: "x float64", format: "{:x}", args: []any{float64(255)}, expected: "ff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.format, tt.args...)
			if tt.name == "pointer" {
				if !strings.HasPrefix(string(result), "0x") {
					t.Errorf("expected pointer format starting with 0x, got '%s'", result)
				}
				return
			}
			if string(result) != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
