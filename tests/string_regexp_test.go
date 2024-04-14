package g_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/enetx/g"
)

func TestReplaceRegexp(t *testing.T) {
	testCases := []struct {
		input     g.String
		pattern   *regexp.Regexp
		newString g.String
		expected  g.String
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
		result := tc.input.ReplaceRegexp(tc.pattern, tc.newString)
		if result != tc.expected {
			t.Errorf("Expected %s, but got %s for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestFindRegexp(t *testing.T) {
	testCases := []struct {
		pattern  *regexp.Regexp
		expected g.Option[g.String]
		input    g.String
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, world!",
			pattern:  regexp.MustCompile(`\bworld\b`),
			expected: g.Some[g.String]("world"),
		},
		// Test case 2: Match with special characters
		{
			input:    "Hello, 12345!",
			pattern:  regexp.MustCompile(`\d+`),
			expected: g.Some[g.String]("12345"),
		},
		// Test case 3: No match
		{
			input:    "Hello, world!",
			pattern:  regexp.MustCompile(`\buniverse\b`),
			expected: g.None[g.String](),
		},
		// Test case 4: Empty input
		{
			input:    "",
			pattern:  regexp.MustCompile(`\d`),
			expected: g.None[g.String](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.FindRegexp(tc.pattern)
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

func TestContainsRegexp(t *testing.T) {
	testCases := []struct {
		pattern  g.String
		input    g.String
		expected bool
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, world!",
			pattern:  `\bworld\b`,
			expected: true,
		},
		// Test case 2: Match with special characters
		{
			input:    "Hello, 12345!",
			pattern:  `\d+`,
			expected: true,
		},
		// Test case 3: No match
		{
			input:    "Hello, world!",
			pattern:  `\buniverse\b`,
			expected: false,
		},
		// Test case 4: Empty input
		{
			input:    "",
			pattern:  `\d`,
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := tc.input.ContainsRegexp(tc.pattern).Ok()
		if result != tc.expected {
			t.Errorf("Expected %v, but got %v for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestContainsRegexpAny(t *testing.T) {
	testCases := []struct {
		input    g.String
		patterns g.Slice[g.String]
		expected bool
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, world!",
			patterns: g.Slice[g.String]{`\bworld\b`},
			expected: true,
		},
		// Test case 2: Multiple patterns, one matches
		{
			input:    "Hello, world!",
			patterns: g.Slice[g.String]{`\bworld\b`, `\d+`},
			expected: true,
		},
		// Test case 3: Multiple patterns, none matches
		{
			input:    "Hello, world!",
			patterns: g.Slice[g.String]{`\buniverse\b`, `\d`},
			expected: false,
		},
		// Test case 4: Empty input
		{
			input:    "",
			patterns: g.Slice[g.String]{`\d`},
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := tc.input.ContainsRegexpAny(tc.patterns...).Ok()
		if result != tc.expected {
			t.Errorf("Expected %v, but got %v for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestContainsRegexpAll(t *testing.T) {
	testCases := []struct {
		input    g.String
		patterns g.Slice[g.String]
		expected bool
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, world!",
			patterns: g.Slice[g.String]{`\bworld\b`},
			expected: true,
		},
		// Test case 2: Multiple patterns, all match
		{
			input:    "Hello, 12345!",
			patterns: g.Slice[g.String]{`\bHello\b`, `\d+`},
			expected: true,
		},
		// Test case 3: Multiple patterns, some match
		{
			input:    "Hello, world!",
			patterns: g.Slice[g.String]{`\bworld\b`, `\d`},
			expected: false,
		},
		// Test case 4: Empty input
		{
			input:    "",
			patterns: g.Slice[g.String]{`\d`},
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := tc.input.ContainsRegexpAll(tc.patterns...).Ok()
		if result != tc.expected {
			t.Errorf("Expected %v, but got %v for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestSplitRegexp(t *testing.T) {
	testCases := []struct {
		input    g.String
		expected g.Slice[g.String]
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular split
		{
			input:    "one,two,three",
			pattern:  *regexp.MustCompile(`,`),
			expected: g.Slice[g.String]{"one", "two", "three"},
		},
		// Test case 2: Split with multiple patterns
		{
			input:    "1, 2, 3, 4",
			pattern:  *regexp.MustCompile(`\s*,\s*`),
			expected: g.Slice[g.String]{"1", "2", "3", "4"},
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`,`),
			expected: g.Slice[g.String]{""},
		},
		// Test case 4: No match
		{
			input:    "abcdefgh",
			pattern:  *regexp.MustCompile(`,`),
			expected: g.Slice[g.String]{"abcdefgh"},
		},
	}

	for _, tc := range testCases {
		result := tc.input.SplitRegexp(tc.pattern)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Expected %v, but got %v for input %s", tc.expected, result, tc.input)
		}
	}
}

func TestSplitRegexpN(t *testing.T) {
	testCases := []struct {
		expected g.Option[g.Slice[g.String]]
		input    g.String
		pattern  regexp.Regexp
		n        g.Int
	}{
		// Test case 1: Regular split with n = 2
		{
			input:    "one,two,three",
			pattern:  *regexp.MustCompile(`,`),
			n:        2,
			expected: g.Some(g.Slice[g.String]{"one", "two,three"}),
		},
		// Test case 2: Split with multiple patterns with n = 0
		{
			input:    "1, 2, 3, 4",
			pattern:  *regexp.MustCompile(`\s*,\s*`),
			n:        0,
			expected: g.None[g.Slice[g.String]](),
		},
		// Test case 3: Empty input with n = 1
		{
			input:    "",
			pattern:  *regexp.MustCompile(`,`),
			n:        1,
			expected: g.Some(g.Slice[g.String]{""}),
		},
		// Test case 4: No match with n = -1
		{
			input:    "abcdefgh",
			pattern:  *regexp.MustCompile(`,`),
			n:        -1,
			expected: g.Some(g.Slice[g.String]{"abcdefgh"}),
		},
	}

	for _, tc := range testCases {
		result := tc.input.SplitRegexpN(tc.pattern, tc.n)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Expected %v, but got %v for input %s with n = %d", tc.expected, result, tc.input, tc.n)
		}
	}
}

func TestIndexRegexp(t *testing.T) {
	testCases := []struct {
		expected g.Option[g.Slice[g.Int]]
		input    g.String
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`World`),
			expected: g.Some(g.Slice[g.Int]{7, 12}),
		},
		// Test case 2: No match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			expected: g.None[g.Slice[g.Int]](),
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`World`),
			expected: g.None[g.Slice[g.Int]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.IndexRegexp(&tc.pattern)
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

func TestFindAllRegexp(t *testing.T) {
	testCases := []struct {
		expected g.Option[g.Slice[g.String]]
		input    g.String
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular matches
		{
			input:    "Hello, World! Hello, Universe!",
			pattern:  *regexp.MustCompile(`Hello`),
			expected: g.Some(g.Slice[g.String]{"Hello", "Hello"}),
		},
		// Test case 2: No match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			expected: g.None[g.Slice[g.String]](),
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			expected: g.None[g.Slice[g.String]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.FindAllRegexp(&tc.pattern)
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

func TestFindAllRegexpN(t *testing.T) {
	testCases := []struct {
		expected g.Option[g.Slice[g.String]]
		input    g.String
		pattern  regexp.Regexp
		n        g.Int
	}{
		// Test case 1: Regular matches with n = 2
		{
			input:    "Hello, World! Hello, Universe!",
			pattern:  *regexp.MustCompile(`Hello`),
			n:        2,
			expected: g.Some(g.Slice[g.String]{"Hello", "Hello"}),
		},
		// Test case 2: No match with n = -1
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			n:        -1,
			expected: g.None[g.Slice[g.String]](),
		},
		// Test case 3: Empty input with n = 1
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			n:        1,
			expected: g.None[g.Slice[g.String]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.FindAllRegexpN(&tc.pattern, tc.n)
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

func TestFindSubmatchRegexp(t *testing.T) {
	testCases := []struct {
		expected g.Option[g.Slice[g.String]]
		input    g.String
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Hello, (\w+)!`),
			expected: g.Some(g.Slice[g.String]{"Hello, World!", "World"}),
		},
		// Test case 2: No match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			expected: g.None[g.Slice[g.String]](),
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			expected: g.None[g.Slice[g.String]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.FindSubmatchRegexp(&tc.pattern)
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

func TestFindAllSubmatchRegexp(t *testing.T) {
	testCases := []struct {
		expected g.Option[g.Slice[g.Slice[g.String]]]
		input    g.String
		pattern  regexp.Regexp
	}{
		// Test case 1: Regular matches
		{
			input:    "Hello, World! Hello, Universe!",
			pattern:  *regexp.MustCompile(`Hello, (\w+)!`),
			expected: g.Some(g.Slice[g.Slice[g.String]]{{"Hello, World!", "World"}, {"Hello, Universe!", "Universe"}}),
		},
		// Test case 2: No match
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			expected: g.None[g.Slice[g.Slice[g.String]]](),
		},
		// Test case 3: Empty input
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			expected: g.None[g.Slice[g.Slice[g.String]]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.FindAllSubmatchRegexp(&tc.pattern)
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

func TestFindAllSubmatchRegexpN(t *testing.T) {
	testCases := []struct {
		expected g.Option[g.Slice[g.Slice[g.String]]]
		input    g.String
		pattern  regexp.Regexp
		n        g.Int
	}{
		// Test case 1: Regular matches with n = 2
		{
			input:    "Hello, World! Hello, Universe!",
			pattern:  *regexp.MustCompile(`Hello, (\w+)!`),
			n:        2,
			expected: g.Some(g.Slice[g.Slice[g.String]]{{"Hello, World!", "World"}, {"Hello, Universe!", "Universe"}}),
		},
		// Test case 2: No match with n = -1
		{
			input:    "Hello, World!",
			pattern:  *regexp.MustCompile(`Earth`),
			n:        -1,
			expected: g.None[g.Slice[g.Slice[g.String]]](),
		},
		// Test case 3: Empty input with n = 1
		{
			input:    "",
			pattern:  *regexp.MustCompile(`Hello`),
			n:        1,
			expected: g.None[g.Slice[g.Slice[g.String]]](),
		},
	}

	for _, tc := range testCases {
		result := tc.input.FindAllSubmatchRegexpN(&tc.pattern, tc.n)
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
