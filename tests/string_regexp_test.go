package g_test

import (
	"reflect"
	"regexp"
	"testing"

	. "github.com/enetx/g"
)

func TestToRegexp(t *testing.T) {
	tests := []struct {
		name         string
		input        String
		expectsOk    bool
		expectsError bool
		pattern      string
	}{
		{
			name:         "ValidRegex",
			input:        String(`^\d+$`),
			expectsOk:    true,
			expectsError: false,
			pattern:      `^\d+$`,
		},
		{
			name:         "InvalidRegex",
			input:        String(`^[`),
			expectsOk:    false,
			expectsError: true,
			pattern:      "",
		},
		{
			name:         "EmptyRegex",
			input:        String(""),
			expectsOk:    true,
			expectsError: false,
			pattern:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Regexp().Compile()
			if tt.expectsOk {
				if !result.IsOk() {
					t.Errorf("Expected result to be OK, got error: %v", result.Err())
				}
				compiledRegex := result.Unwrap()
				if tt.pattern != "" && compiledRegex.String() != tt.pattern {
					t.Errorf("Compiled pattern mismatch: got %v, want %v", compiledRegex.String(), tt.pattern)
				}
			} else {
				if result.IsOk() {
					t.Errorf("Expected result to be an error, but got OK")
				}
				if result.Err() == nil {
					t.Errorf("Expected an error, but got nil")
				}
			}
		})
	}
}

func TestRxReplace(t *testing.T) {
	testCases := []struct {
		input     String
		pattern   *regexp.Regexp
		newString String
		expected  String
	}{
		// Test case 1: Regular replacement
		{
			input:     "Hello, world!",
			pattern:   regexp.MustCompile(`\bworld\b`),
			newString: "universe",
			expected:  "Hello, universe!",
		},
		// Test case 2: Replacement with empty string
		{
			input:     "apple, orange, apple, banana",
			pattern:   regexp.MustCompile(`\bapple\b`),
			newString: "",
			expected:  ", orange, , banana",
		},
		// Test case 3: Replacement with special characters
		{
			input:     "1 + 2 = 3",
			pattern:   regexp.MustCompile(`\d`),
			newString: "x",
			expected:  "x + x = x",
		},
		// Test case 4: No match
		{
			input:     "Hello, world!",
			pattern:   regexp.MustCompile(`\buniverse\b`),
			newString: "galaxy",
			expected:  "Hello, world!",
		},
		// Test case 5: Empty input
		{
			input:     "",
			pattern:   regexp.MustCompile(`\d`),
			newString: "x",
			expected:  "",
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().Replace(tc.pattern, tc.newString)
		if result != tc.expected {
			t.Errorf("Expected %s, but got %s for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestRxFind(t *testing.T) {
	testCases := []struct {
		pattern  *regexp.Regexp
		expected Option[String]
		input    String
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, world!",
			pattern:  regexp.MustCompile(`\bworld\b`),
			expected: Some[String]("world"),
		},
		// Test case 2: Match with special characters
		{
			input:    "Hello, 12345!",
			pattern:  regexp.MustCompile(`\d+`),
			expected: Some[String]("12345"),
		},
		// Test case 3: No match
		{
			input:    "Hello, world!",
			pattern:  regexp.MustCompile(`\buniverse\b`),
			expected: None[String](),
		},
		// Test case 4: Empty input
		{
			input:    "",
			pattern:  regexp.MustCompile(`\d`),
			expected: None[String](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().Find(tc.pattern)
		if result.IsSome() {
			if result.Some() != tc.expected.Some() {
				t.Errorf("Expected %s, but got %s for input %s", tc.expected.Some(), result.Some(), tc.input)
			}
		} else {
			if result.IsSome() != tc.expected.IsSome() {
				t.Errorf("Expected None, but got Some for input %s", tc.input)
			}
		}
	}
}

func TestRxMatch(t *testing.T) {
	testCases := []struct {
		pattern  *regexp.Regexp
		input    String
		expected bool
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, world!",
			pattern:  regexp.MustCompile(`\bworld\b`),
			expected: true,
		},
		// Test case 2: Match with special characters
		{
			input:    "Hello, 12345!",
			pattern:  regexp.MustCompile(`\d+`),
			expected: true,
		},
		// Test case 3: No match
		{
			input:    "Hello, world!",
			pattern:  regexp.MustCompile(`\buniverse\b`),
			expected: false,
		},
		// Test case 4: Empty input
		{
			input:    "",
			pattern:  regexp.MustCompile(`\d`),
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().Match(tc.pattern)
		if result != tc.expected {
			t.Errorf("Expected %v, but got %v for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestRxMatchAny(t *testing.T) {
	testCases := []struct {
		input    String
		patterns Slice[*regexp.Regexp]
		expected bool
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, world!",
			patterns: Slice[*regexp.Regexp]{regexp.MustCompile(`\bworld\b`)},
			expected: true,
		},
		// Test case 2: Multiple patterns, one matches
		{
			input:    "Hello, world!",
			patterns: Slice[*regexp.Regexp]{regexp.MustCompile(`\bworld\b`), regexp.MustCompile(`\d+`)},
			expected: true,
		},
		// Test case 3: Multiple patterns, none matches
		{
			input:    "Hello, world!",
			patterns: Slice[*regexp.Regexp]{regexp.MustCompile(`\buniverse\b`), regexp.MustCompile(`\d`)},
			expected: false,
		},
		// Test case 4: Empty input
		{
			input:    "",
			patterns: Slice[*regexp.Regexp]{regexp.MustCompile(`\d`)},
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().MatchAny(tc.patterns...)
		if result != tc.expected {
			t.Errorf("Expected %v, but got %v for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestRxMatchAll(t *testing.T) {
	testCases := []struct {
		input    String
		patterns Slice[*regexp.Regexp]
		expected bool
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, world!",
			patterns: Slice[*regexp.Regexp]{regexp.MustCompile(`\bworld\b`)},
			expected: true,
		},
		// Test case 2: Multiple patterns, all match
		{
			input:    "Hello, 12345!",
			patterns: Slice[*regexp.Regexp]{regexp.MustCompile(`\bHello\b`), regexp.MustCompile(`\d+`)},
			expected: true,
		},
		// Test case 3: Multiple patterns, some match
		{
			input:    "Hello, world!",
			patterns: Slice[*regexp.Regexp]{regexp.MustCompile(`\bworld\b`), regexp.MustCompile(`\d`)},
			expected: false,
		},
		// Test case 4: Empty input
		{
			input:    "",
			patterns: Slice[*regexp.Regexp]{regexp.MustCompile(`\d`)},
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().MatchAll(tc.patterns...)
		if result != tc.expected {
			t.Errorf("Expected %v, but got %v for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestRxSplit(t *testing.T) {
	testCases := []struct {
		input    String
		expected Slice[String]
		pattern  *regexp.Regexp
	}{
		// Test case 1: Regular split
		{
			input:    "one,two,three",
			pattern:  regexp.MustCompile(`,`),
			expected: Slice[String]{"one", "two", "three"},
		},
		// Test case 2: Split with multiple patterns
		{
			input:    "1, 2, 3, 4",
			pattern:  regexp.MustCompile(`\s*,\s*`),
			expected: Slice[String]{"1", "2", "3", "4"},
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  regexp.MustCompile(`,`),
			expected: Slice[String]{""},
		},
		// Test case 4: No match
		{
			input:    "abcdefgh",
			pattern:  regexp.MustCompile(`,`),
			expected: Slice[String]{"abcdefgh"},
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().Split(tc.pattern)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Expected %v, but got %v for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestRxSplitN(t *testing.T) {
	testCases := []struct {
		expected Option[Slice[String]]
		input    String
		pattern  *regexp.Regexp
		n        Int
	}{
		// Test case 1: Regular split with n = 2
		{
			input:    "one,two,three",
			pattern:  regexp.MustCompile(`,`),
			n:        2,
			expected: Some(Slice[String]{"one", "two,three"}),
		},
		// Test case 2: Split with multiple patterns with n = 0
		{
			input:    "1, 2, 3, 4",
			pattern:  regexp.MustCompile(`\s*,\s*`),
			n:        0,
			expected: None[Slice[String]](),
		},
		// Test case 3: Empty input with n = 1
		{
			input:    "",
			pattern:  regexp.MustCompile(`,`),
			n:        1,
			expected: Some(Slice[String]{""}),
		},
		// Test case 4: No match with n = -1
		{
			input:    "abcdefgh",
			pattern:  regexp.MustCompile(`,`),
			n:        -1,
			expected: Some(Slice[String]{"abcdefgh"}),
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().SplitN(tc.pattern, tc.n)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Expected %v, but got %v for input %s with n = %d", tc.expected, result, tc.input, tc.n)
		}
	}
}

func TestRxIndex(t *testing.T) {
	testCases := []struct {
		expected Option[Slice[Int]]
		input    String
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`World`),
			expected: Some(Slice[Int]{7, 12}),
		},
		// Test case 2: No match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			expected: None[Slice[Int]](),
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`World`),
			expected: None[Slice[Int]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().Index(&tc.pattern)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf(
				"Expected %v, but got %v for input %s with pattern %s",
				tc.expected,
				result,
				tc.input,
				tc.pattern.String(),
			)
		}
	}
}

func TestRxFindAll(t *testing.T) {
	testCases := []struct {
		expected Option[Slice[String]]
		input    String
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular matches
		{
			input:    "Hello, World! Hello, Universe!",
			pattern:  *regexp.MustCompile(`Hello`),
			expected: Some(Slice[String]{"Hello", "Hello"}),
		},
		// Test case 2: No match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			expected: None[Slice[String]](),
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			expected: None[Slice[String]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().FindAll(&tc.pattern)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf(
				"Expected %v, but got %v for input %s with pattern %s",
				tc.expected,
				result,
				tc.input,
				tc.pattern.String(),
			)
		}
	}
}

func TestRxFindAllN(t *testing.T) {
	testCases := []struct {
		expected Option[Slice[String]]
		input    String
		pattern  regexp.Regexp
		n        Int
	}{
		// Test case 1: Regular matches with n = 2
		{
			input:    "Hello, World! Hello, Universe!",
			pattern:  *regexp.MustCompile(`Hello`),
			n:        2,
			expected: Some(Slice[String]{"Hello", "Hello"}),
		},
		// Test case 2: No match with n = -1
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			n:        -1,
			expected: None[Slice[String]](),
		},
		// Test case 3: Empty input with n = 1
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			n:        1,
			expected: None[Slice[String]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().FindAllN(&tc.pattern, tc.n)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf(
				"Expected %v, but got %v for input %s with pattern %s and n = %d",
				tc.expected,
				result,
				tc.input,
				tc.pattern.String(),
				tc.n,
			)
		}
	}
}

func TestRxFindSubmatch(t *testing.T) {
	testCases := []struct {
		expected Option[Slice[String]]
		input    String
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Hello, (\w+)!`),
			expected: Some(Slice[String]{"Hello, World!", "World"}),
		},
		// Test case 2: No match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			expected: None[Slice[String]](),
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			expected: None[Slice[String]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().FindSubmatch(&tc.pattern)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf(
				"Expected %v, but got %v for input %s with pattern %s",
				tc.expected,
				result,
				tc.input,
				tc.pattern.String(),
			)
		}
	}
}

func TestRxFindAllSubmatch(t *testing.T) {
	testCases := []struct {
		expected Option[Slice[Slice[String]]]
		input    String
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular matches
		{
			input:    "Hello, World! Hello, Universe!",
			pattern:  *regexp.MustCompile(`Hello, (\w+)!`),
			expected: Some(Slice[Slice[String]]{{"Hello, World!", "World"}, {"Hello, Universe!", "Universe"}}),
		},
		// Test case 2: No match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			expected: None[Slice[Slice[String]]](),
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			expected: None[Slice[Slice[String]]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().FindAllSubmatch(&tc.pattern)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf(
				"Expected %v, but got %v for input %s with pattern %s",
				tc.expected,
				result,
				tc.input,
				tc.pattern.String(),
			)
		}
	}
}

func TestRxFindAllSubmatchN(t *testing.T) {
	testCases := []struct {
		expected Option[Slice[Slice[String]]]
		input    String
		pattern  regexp.Regexp
		n        Int
	}{
		// Test case 1: Regular matches with n = 2
		{
			input:    "Hello, World! Hello, Universe!",
			pattern:  *regexp.MustCompile(`Hello, (\w+)!`),
			n:        2,
			expected: Some(Slice[Slice[String]]{{"Hello, World!", "World"}, {"Hello, Universe!", "Universe"}}),
		},
		// Test case 2: No match with n = -1
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			n:        -1,
			expected: None[Slice[Slice[String]]](),
		},
		// Test case 3: Empty input with n = 1
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			n:        1,
			expected: None[Slice[Slice[String]]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.Regexp().FindAllSubmatchN(&tc.pattern, tc.n)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf(
				"Expected %v, but got %v for input %s with pattern %s and n = %d",
				tc.expected,
				result,
				tc.input,
				tc.pattern.String(),
				tc.n,
			)
		}
	}
}
