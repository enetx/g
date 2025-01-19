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
			format:   "Values: {1}, {2}, {1.$lower}",
			args:     []any{"X", "Y"},
			expected: "Values: X, Y, x",
		},
		{
			name:     "Escaped braces",
			format:   "Show literal \\{{.$upper}\\} here",
			args:     []any{"upper"},
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
			format:   "Name: {name.$upper}",
			args:     Named{"name": "john"},
			expected: "Name: JOHN",
		},
		// Placeholder with modifier: trim and title
		{
			name:     "Modifier: trim and title",
			format:   "Title: {work.$trim.$title}",
			args:     Named{"work": " developer "},
			expected: "Title: Developer",
		},
		// Nested modifiers: trim and len
		{
			name:     "Nested modifiers",
			format:   "Length: {input.$trim.$len}",
			args:     Named{"input": "  data  "},
			expected: "Length: 4",
		},
		// Placeholder with fallback and modifier
		{
			name:     "Fallback with modifier",
			format:   "Name: {name?fallback.$upper}",
			args:     Named{"fallback": "guest"},
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
			format:   "Value: {number.$round}",
			args:     Named{"number": 12.7},
			expected: "Value: 13",
		},
		// Modifier: abs for negative numbers
		{
			name:     "Modifier: abs",
			format:   "Absolute: {value.$abs}",
			args:     Named{"value": -42},
			expected: "Absolute: 42",
		},
		// Modifier: reverse for strings
		{
			name:     "Modifier: reverse",
			format:   "Reversed: {word.$reverse}",
			args:     Named{"word": "hello"},
			expected: "Reversed: olleh",
		},
		// Modifier: hex for integers
		{
			name:     "Modifier: hex",
			format:   "Hex: {number.$hex}",
			args:     Named{"number": 255},
			expected: "Hex: ff",
		},
		// Modifier: bin for integers
		{
			name:     "Modifier: bin",
			format:   "Binary: {number.$bin}",
			args:     Named{"number": 5},
			expected: "Binary: 00000101",
		},
		// Modifier: url encoding
		{
			name:     "Modifier: url",
			format:   "URL: {input.$url}",
			args:     Named{"input": "hello world"},
			expected: "URL: hello+world",
		},
		// Modifier: base64 encoding
		{
			name:     "Modifier: base64",
			format:   "Base64: {input.$base64e}",
			args:     Named{"input": "hello"},
			expected: "Base64: aGVsbG8=",
		},
		// Modifier: format for dates
		{
			name:     "Modifier: format date",
			format:   "Date: {today.$date(2006-01-02)}",
			args:     Named{"today": time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC)},
			expected: "Date: 2025-01-17",
		},
		// Test for $replace
		{
			name:     "Modifier: replace",
			format:   "Result: {input.$replace(a,b)}",
			args:     Named{"input": "banana"},
			expected: "Result: bbnbnb",
		},
		{
			name:     "Modifier: replace with empty string",
			format:   "Result: {input.$replace(a,)}",
			args:     Named{"input": "banana"},
			expected: "Result: bnn",
		},
		{
			name:     "Modifier: replace no matches",
			format:   "Result: {input.$replace(x,y)}",
			args:     Named{"input": "banana"},
			expected: "Result: banana",
		},
		// Test for $repeat
		{
			name:     "Modifier: repeat string",
			format:   "Repeated: {input.$repeat(3)}",
			args:     Named{"input": "ha"},
			expected: "Repeated: hahaha",
		},
		{
			name:     "Modifier: repeat int",
			format:   "Repeated: {input.$repeat(4)}",
			args:     Named{"input": 5},
			expected: "Repeated: 5555",
		},
		{
			name:     "Modifier: repeat with invalid count",
			format:   "Repeated: {input.$repeat(abc)}",
			args:     Named{"input": "ha"},
			expected: "Repeated: ha",
		},
		// Test for $substring
		{
			name:     "Modifier: substring",
			format:   "Result: {input.$substring(0,-1,2)}",
			args:     Named{"input": "Hello, World!"},
			expected: "Result: Hlo ol",
		},
		// Test for $truncate
		{
			name:     "Modifier: truncate string",
			format:   "Truncated: {input.$truncate(5)}",
			args:     Named{"input": "Hello, World!"},
			expected: "Truncated: Hello...",
		},
		{
			name:     "Modifier: truncate with exact length",
			format:   "Truncated: {input.$truncate(5)}",
			args:     Named{"input": "Hello"},
			expected: "Truncated: Hello",
		},
		{
			name:     "Modifier: truncate with no truncation",
			format:   "Truncated: {input.$truncate(15)}",
			args:     Named{"input": "Hello, World!"},
			expected: "Truncated: Hello, World!",
		},
		{
			name:     "Modifier: truncate with invalid max",
			format:   "Truncated: {input.$truncate(abc)}",
			args:     Named{"input": "Hello, World!"},
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
			format: "{word.$trim.$lower.$replace(e,a).$reverse}",
			args:   Named{"word": "  EXAMPLE "},
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
			format:   "Value: {number.$unknown(",
			args:     Named{"number": 42},
			expected: "Value: {number.$unknown(",
		},
		// Unsupported modifier
		{
			name:     "Unsupported modifier",
			format:   "Value: {number.$unsupported}",
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
			format:   "Value: {obj.$upper}",
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
			name:     "Trim specific characters",
			format:   "Result: {value.$trim(#)}",
			args:     Named{"value": "###Hello###"},
			expected: "Result: Hello",
		},
		// Trim multiple characters
		{
			name:     "Trim multiple characters",
			format:   "Result: {value.$trim(#$)}",
			args:     Named{"value": "$$#Hello#$"},
			expected: "Result: Hello",
		},
		// No trimming (no matching characters)
		{
			name:     "No trimming needed",
			format:   "Result: {value.$trim(%)}",
			args:     Named{"value": "Hello"},
			expected: "Result: Hello",
		},
		// Empty value
		{
			name:     "Empty value",
			format:   "Result: {value.$trim(#)}",
			args:     Named{"value": ""},
			expected: "Result: ",
		},
		// Empty set
		{
			name:     "Empty trim set",
			format:   "Result: {value.$trim()}",
			args:     Named{"value": "###Hello###"},
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
