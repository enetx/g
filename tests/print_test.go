package g_test

import (
	"errors"
	"fmt"
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
			name:     "Modifier type",
			format:   "Type: {1.type}",
			args:     []any{Int(42)},
			expected: "Type: g.Int",
		},
		{
			name:     "Modifier debug on string",
			format:   "Debug: {1.debug}",
			args:     []any{String("hi")},
			expected: `Debug: "hi"`,
		},
		{
			name:     "Modifier debug on int",
			format:   "Debug: {1.debug}",
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
			name:     "Modifier type on named",
			format:   "Type: {val.type}",
			args:     Named{"val": Int(10)},
			expected: "Type: g.Int",
		},
		{
			name:     "Modifier debug on named",
			format:   "Debug: {val.debug}",
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

	t.Run("positional wrap {1.wrap}", func(t *testing.T) {
		err := Errorf("failed: {1.wrap}", sentinel)
		if err.Error() != "failed: sentinel" {
			t.Fatalf("unexpected message: %s", err)
		}
		if !errors.Is(err, sentinel) {
			t.Fatal("expected errors.Is to find sentinel")
		}
	})

	t.Run("auto-index wrap {.wrap}", func(t *testing.T) {
		err := Errorf("open {}: {.wrap}", "foo.txt", sentinel)
		if err.Error() != "open foo.txt: sentinel" {
			t.Fatalf("unexpected message: %s", err)
		}
		if !errors.Is(err, sentinel) {
			t.Fatal("expected errors.Is to find sentinel")
		}
	})

	t.Run("named wrap {cause.wrap}", func(t *testing.T) {
		err := Errorf("read {file}: {cause.wrap}", Named{"file": "bar.txt", "cause": sentinel})
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
		err := Errorf("{1.wrap} and {2.wrap}", e1, e2)
		if !errors.Is(err, e1) {
			t.Fatal("expected errors.Is to find e1")
		}
		if !errors.Is(err, e2) {
			t.Fatal("expected errors.Is to find e2")
		}
	})

	t.Run("wrap with chained modifier", func(t *testing.T) {
		err := Errorf("msg: {1.wrap}", sentinel)
		if !errors.Is(err, sentinel) {
			t.Fatal("expected wrapping")
		}
	})

	t.Run("wrap on non-error value — no panic", func(t *testing.T) {
		err := Errorf("value: {1.wrap}", "not an error")
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
		err := Errorf("err: {1.wrap}", nil)
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
