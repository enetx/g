package g_test

import (
	"testing"
	"time"

	. "github.com/enetx/g"
)

func TestSprinfAutoIndexAndNumeric(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sprintf(tt.format, tt.args...)
			if string(result) != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSprintf(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sprintf(tt.format, tt.args)
			if result != String(tt.expected) {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSprintfFormatWithErrors(t *testing.T) {
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
			args:     Named{"obj": struct{}{}},
			expected: "Value: {}",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sprintf(tt.format, tt.args)
			if result != String(tt.expected) {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSprintfTrimSetModifier(t *testing.T) {
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
			result := Sprintf(tt.format, tt.args)
			if result != String(tt.expected) {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
