package g_test

import (
	"testing"
	"time"

	. "github.com/enetx/g"
)

func TestStringFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     map[string]any
		expected string
	}{
		// Basic placeholder replacement
		{
			name:     "Basic replacement",
			format:   "Hello, {name}!",
			args:     map[string]any{"name": "John"},
			expected: "Hello, John!",
		},
		// Placeholder with fallback
		{
			name:     "Fallback replacement",
			format:   "Hello, {name?fallback}!",
			args:     map[string]any{"fallback": "Guest"},
			expected: "Hello, Guest!",
		},
		// Placeholder with modifier: upper
		{
			name:     "Modifier: upper",
			format:   "Name: {name.$upper}",
			args:     map[string]any{"name": "john"},
			expected: "Name: JOHN",
		},
		// Placeholder with modifier: trim and title
		{
			name:     "Modifier: trim and title",
			format:   "Title: {work.$trim.$title}",
			args:     map[string]any{"work": " developer "},
			expected: "Title: Developer",
		},
		// Nested modifiers: trim and len
		{
			name:     "Nested modifiers",
			format:   "Length: {input.$trim.$len}",
			args:     map[string]any{"input": "  data  "},
			expected: "Length: 4",
		},
		// Placeholder with fallback and modifier
		{
			name:     "Fallback with modifier",
			format:   "Name: {name?fallback.$upper}",
			args:     map[string]any{"fallback": "guest"},
			expected: "Name: GUEST",
		},
		// Multiple placeholders
		{
			name:     "Multiple placeholders",
			format:   "{greeting}, {name}! You are {age} years old.",
			args:     map[string]any{"greeting": "Hello", "name": "John", "age": 30},
			expected: "Hello, John! You are 30 years old.",
		},
		// Placeholder with unknown key
		{
			name:     "Unknown placeholder",
			format:   "Hello, {unknown}!",
			args:     map[string]any{"name": "John"},
			expected: "Hello, {unknown}!",
		},
		// Modifier: round for float values
		{
			name:     "Modifier: round",
			format:   "Value: {number.$round}",
			args:     map[string]any{"number": 12.7},
			expected: "Value: 13",
		},
		// Modifier: abs for negative numbers
		{
			name:     "Modifier: abs",
			format:   "Absolute: {value.$abs}",
			args:     map[string]any{"value": -42},
			expected: "Absolute: 42",
		},
		// Modifier: reverse for strings
		{
			name:     "Modifier: reverse",
			format:   "Reversed: {word.$reverse}",
			args:     map[string]any{"word": "hello"},
			expected: "Reversed: olleh",
		},
		// Modifier: hex for integers
		{
			name:     "Modifier: hex",
			format:   "Hex: {number.$hex}",
			args:     map[string]any{"number": 255},
			expected: "Hex: ff",
		},
		// Modifier: bin for integers
		{
			name:     "Modifier: bin",
			format:   "Binary: {number.$bin}",
			args:     map[string]any{"number": 5},
			expected: "Binary: 00000101",
		},
		// Modifier: url encoding
		{
			name:     "Modifier: url",
			format:   "URL: {input.$url}",
			args:     map[string]any{"input": "hello world"},
			expected: "URL: hello+world",
		},
		// Modifier: base64 encoding
		{
			name:     "Modifier: base64",
			format:   "Base64: {input.$base64e}",
			args:     map[string]any{"input": "hello"},
			expected: "Base64: aGVsbG8=",
		},
		// Modifier: format for dates
		{
			name:     "Modifier: format date",
			format:   "Date: {today.$date(2006-01-02)}",
			args:     map[string]any{"today": time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC)},
			expected: "Date: 2025-01-17",
		},
		// Test for $replace
		{
			name:     "Modifier: replace",
			format:   "Result: {input.$replace(a,b)}",
			args:     map[string]any{"input": "banana"},
			expected: "Result: bbnbnb",
		},
		{
			name:     "Modifier: replace with empty string",
			format:   "Result: {input.$replace(a,)}",
			args:     map[string]any{"input": "banana"},
			expected: "Result: bnn",
		},
		{
			name:     "Modifier: replace no matches",
			format:   "Result: {input.$replace(x,y)}",
			args:     map[string]any{"input": "banana"},
			expected: "Result: banana",
		},

		// Test for $repeat
		{
			name:     "Modifier: repeat string",
			format:   "Repeated: {input.$repeat(3)}",
			args:     map[string]any{"input": "ha"},
			expected: "Repeated: hahaha",
		},
		{
			name:     "Modifier: repeat int",
			format:   "Repeated: {input.$repeat(4)}",
			args:     map[string]any{"input": 5},
			expected: "Repeated: 5555",
		},
		{
			name:     "Modifier: repeat with invalid count",
			format:   "Repeated: {input.$repeat(abc)}",
			args:     map[string]any{"input": "ha"},
			expected: "Repeated: ha",
		},

		// Test for $substring
		{
			name:     "Modifier: substring",
			format:   "Result: {input.$substring(0,-1,2)}",
			args:     map[string]any{"input": "Hello, World!"},
			expected: "Result: Hlo ol",
		},

		// Test for $truncate
		{
			name:     "Modifier: truncate string",
			format:   "Truncated: {input.$truncate(5)}",
			args:     map[string]any{"input": "Hello, World!"},
			expected: "Truncated: Hello...",
		},
		{
			name:     "Modifier: truncate with exact length",
			format:   "Truncated: {input.$truncate(5)}",
			args:     map[string]any{"input": "Hello"},
			expected: "Truncated: Hello",
		},
		{
			name:     "Modifier: truncate with no truncation",
			format:   "Truncated: {input.$truncate(15)}",
			args:     map[string]any{"input": "Hello, World!"},
			expected: "Truncated: Hello, World!",
		},
		{
			name:     "Modifier: truncate with invalid max",
			format:   "Truncated: {input.$truncate(abc)}",
			args:     map[string]any{"input": "Hello, World!"},
			expected: "Truncated: Hello, World!",
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

func TestStringFormatWithErrors(t *testing.T) {
	errorTests := []struct {
		name     string
		format   string
		args     map[string]any
		expected string
	}{
		// Placeholder with invalid syntax
		{
			name:     "Invalid placeholder syntax",
			format:   "Hello, {name?",
			args:     map[string]any{"name": "John"},
			expected: "Hello, {name?",
		},
		// Modifier with invalid syntax
		{
			name:     "Invalid modifier syntax",
			format:   "Value: {number.$unknown(",
			args:     map[string]any{"number": 42},
			expected: "Value: {number.$unknown(",
		},
		// Unsupported modifier
		{
			name:     "Unsupported modifier",
			format:   "Value: {number.$unsupported}",
			args:     map[string]any{"number": 42},
			expected: "Value: 42",
		},
		// Fallback key missing
		{
			name:     "Missing fallback key",
			format:   "Hello, {name?fallback}!",
			args:     make(map[string]any),
			expected: "Hello, {name?fallback}!",
		},
		// Placeholder with unsupported type
		{
			name:     "Unsupported type for modifier",
			format:   "Value: {obj.$upper}",
			args:     map[string]any{"obj": struct{}{}},
			expected: "Value: {}",
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
		args     map[string]any
		expected string
	}{
		// Basic trimming
		{
			name:     "Trim specific characters",
			format:   "Result: {value.$trim(#)}",
			args:     map[string]any{"value": "###Hello###"},
			expected: "Result: Hello",
		},
		// Trim multiple characters
		{
			name:     "Trim multiple characters",
			format:   "Result: {value.$trim(#$)}",
			args:     map[string]any{"value": "$$#Hello#$"},
			expected: "Result: Hello",
		},
		// No trimming (no matching characters)
		{
			name:     "No trimming needed",
			format:   "Result: {value.$trim(%)}",
			args:     map[string]any{"value": "Hello"},
			expected: "Result: Hello",
		},
		// Empty value
		{
			name:     "Empty value",
			format:   "Result: {value.$trim(#)}",
			args:     map[string]any{"value": ""},
			expected: "Result: ",
		},
		// Empty set
		{
			name:     "Empty trim set",
			format:   "Result: {value.$trim()}",
			args:     map[string]any{"value": "###Hello###"},
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
