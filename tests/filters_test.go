package g_test

import (
	"regexp"
	"testing"

	"github.com/enetx/g/f"
)

func TestIsComparable(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"Int", 10, true},
		{"String", "Hello", true},
		{"Slice", []int{1, 2, 3}, false},
		{"Map", make(map[string]int), false},
		{"Struct", struct{ X int }{}, true},
		{"Func", (func())(nil), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := f.IsComparable(tt.value); got != tt.want {
				t.Errorf("Comparable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	// Testing Zero function with integers
	if !f.IsZero(0) {
		t.Errorf("Zero(0) returned false, expected true")
	}

	if f.IsZero(5) {
		t.Errorf("Zero(5) returned true, expected false")
	}

	// Testing Zero function with floats
	if !f.IsZero(0.0) {
		t.Errorf("Zero(0.0) returned false, expected true")
	}

	if f.IsZero(3.14) {
		t.Errorf("Zero(3.14) returned true, expected false")
	}

	// Testing Zero function with strings
	if !f.IsZero("") {
		t.Errorf("Zero(\"\") returned false, expected true")
	}

	if f.IsZero("hello") {
		t.Errorf("Zero(\"hello\") returned true, expected false")
	}
}

func TestIsEven(t *testing.T) {
	// Testing Even function with positive even integers
	if !f.IsEven(2) {
		t.Errorf("Even(2) returned false, expected true")
	}

	if !f.IsEven(100) {
		t.Errorf("Even(100) returned false, expected true")
	}

	// Testing Even function with positive odd integers
	if f.IsEven(3) {
		t.Errorf("Even(3) returned true, expected false")
	}

	if f.IsEven(101) {
		t.Errorf("Even(101) returned true, expected false")
	}

	// Testing Even function with negative even integers
	if !f.IsEven(-2) {
		t.Errorf("Even(-2) returned false, expected true")
	}

	if !f.IsEven(-100) {
		t.Errorf("Even(-100) returned false, expected true")
	}

	// Testing Even function with negative odd integers
	if f.IsEven(-3) {
		t.Errorf("Even(-3) returned true, expected false")
	}

	if f.IsEven(-101) {
		t.Errorf("Even(-101) returned true, expected false")
	}
}

func TestIsOdd(t *testing.T) {
	// Testing Odd function with positive odd integers
	if !f.IsOdd(3) {
		t.Errorf("Odd(3) returned false, expected true")
	}

	if !f.IsOdd(101) {
		t.Errorf("Odd(101) returned false, expected true")
	}

	// Testing Odd function with positive even integers
	if f.IsOdd(2) {
		t.Errorf("Odd(2) returned true, expected false")
	}

	if f.IsOdd(100) {
		t.Errorf("Odd(100) returned true, expected false")
	}

	// Testing Odd function with negative odd integers
	if !f.IsOdd(-3) {
		t.Errorf("Odd(-3) returned false, expected true")
	}

	if !f.IsOdd(-101) {
		t.Errorf("Odd(-101) returned false, expected true")
	}

	// Testing Odd function with negative even integers
	if f.IsOdd(-2) {
		t.Errorf("Odd(-2) returned true, expected false")
	}

	if f.IsOdd(-100) {
		t.Errorf("Odd(-100) returned true, expected false")
	}
}

func TestEq(t *testing.T) {
	// Test case 1: Test equality of integers
	if !f.Eq(5)(5) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test inequality of strings
	if f.Eq("hello")("world") {
		t.Errorf("Test case 2: Expected false, got true")
	}
}

func TestNe(t *testing.T) {
	// Test case 1: Test inequality of integers
	if !f.Ne(5)(6) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test equality of strings
	if !f.Ne("hello")("world") {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestEqd(t *testing.T) {
	// Test case 1: Test deep equality of slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	if !f.Eqd(slice1)(slice2) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test deep inequality of maps
	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "c": 3}
	if f.Eqd(map1)(map2) {
		t.Errorf("Test case 2: Expected false, got true")
	}
}

func TestNed(t *testing.T) {
	// Test case 1: Test deep inequality of slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 4}
	if !f.Ned(slice1)(slice2) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test deep equality of structs
	type Person struct {
		Name string
		Age  int
	}

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Bob", Age: 25}
	if !f.Ned(person1)(person2) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestGt(t *testing.T) {
	// Test case 1: Test greater than comparison of integers
	if !f.Gt(10)(15) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test greater than comparison of floats
	if !f.Gt(3.14)(3.1416) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestGte(t *testing.T) {
	// Test case 1: Test greater than or equal to comparison of integers
	if !f.Gte(10)(10) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test greater than or equal to comparison of floats
	if !f.Gte(3.14)(3.14) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestLt(t *testing.T) {
	// Test case 1: Test less than comparison of integers
	if !f.Lt(5)(3) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test less than comparison of floats
	if !f.Lt(3.14)(2.71) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestLte(t *testing.T) {
	// Test case 1: Test less than or equal to comparison of integers
	if !f.Lte(5)(5) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test less than or equal to comparison of floats
	if !f.Lte(3.14)(3.14) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestContains(t *testing.T) {
	substr := "world"
	containsFunc := f.Contains(substr)

	input := "hello world"

	if !containsFunc(input) {
		t.Errorf("Expected %q to contain %q, but it did not.", input, substr)
	}

	missingInput := "hello universe"
	if containsFunc(missingInput) {
		t.Errorf("Expected %q not to contain %q, but it did.", missingInput, substr)
	}

	emptySubstring := f.Contains("")
	if !emptySubstring(input) {
		t.Errorf("Expected any string to contain an empty substring, but %q did not.", input)
	}

	if !emptySubstring("") {
		t.Errorf("Expected an empty string to contain an empty substring, but it did not.")
	}
}

func TestContainsAnyChars(t *testing.T) {
	charSet := "abc"
	containsAny := f.ContainsAnyChars(charSet)

	inputWithChars := "hello cat"
	if !containsAny(inputWithChars) {
		t.Errorf("Expected %q to contain at least one of %q, but it did not.", inputWithChars, charSet)
	}

	inputWithoutChars := "hello dog"
	if containsAny(inputWithoutChars) {
		t.Errorf("Expected %q not to contain any of %q, but it did.", inputWithoutChars, charSet)
	}

	exactMatch := "abc"
	if !containsAny(exactMatch) {
		t.Errorf("Expected %q to contain characters from %q, but it did not.", exactMatch, charSet)
	}
}

func TestStartsWith(t *testing.T) {
	prefix := "hello"
	startsWith := f.StartsWith(prefix)

	inputWithPrefix := "hello world"
	if !startsWith(inputWithPrefix) {
		t.Errorf("Expected %q to start with %q, but it did not.", inputWithPrefix, prefix)
	}

	inputWithoutPrefix := "world hello"
	if startsWith(inputWithoutPrefix) {
		t.Errorf("Expected %q not to start with %q, but it did.", inputWithoutPrefix, prefix)
	}

	emptyPrefix := f.StartsWith("")
	if !emptyPrefix(inputWithPrefix) {
		t.Errorf("Expected any string to start with an empty prefix, but %q did not.", inputWithPrefix)
	}

	if !emptyPrefix("") {
		t.Errorf("Expected an empty string to start with an empty prefix, but it did not.")
	}

	exactMatch := "hello"
	if !startsWith(exactMatch) {
		t.Errorf("Expected %q to start with %q, but it did not.", exactMatch, prefix)
	}
}

func TestEndsWith(t *testing.T) {
	suffix := "world"
	endsWith := f.EndsWith(suffix)

	inputWithSuffix := "hello world"
	if !endsWith(inputWithSuffix) {
		t.Errorf("Expected %q to end with %q, but it did not.", inputWithSuffix, suffix)
	}

	inputWithoutSuffix := "world hello"
	if endsWith(inputWithoutSuffix) {
		t.Errorf("Expected %q not to end with %q, but it did.", inputWithoutSuffix, suffix)
	}

	emptySuffix := f.EndsWith("")
	if !emptySuffix(inputWithSuffix) {
		t.Errorf("Expected any string to end with an empty suffix, but %q did not.", inputWithSuffix)
	}

	if !emptySuffix("") {
		t.Errorf("Expected an empty string to end with an empty suffix, but it did not.")
	}

	exactMatch := "world"
	if !endsWith(exactMatch) {
		t.Errorf("Expected %q to end with %q, but it did not.", exactMatch, suffix)
	}
}

func TestFilterRxMatch(t *testing.T) {
	regex := regexp.MustCompile(`^\d+$`)
	matchDigits := f.RxMatch[string](regex)

	tests := []struct {
		input    string
		expected bool
	}{
		{"12345", true},     // Matches: Only digits
		{"abc123", false},   // Does not match: Contains letters
		{"", false},         // Does not match: Empty string
		{" 12345 ", false},  // Does not match: Leading/trailing spaces
		{"0000", true},      // Matches: Only digits
		{"123\n456", false}, // Does not match: Contains newline
	}

	for _, test := range tests {
		result := matchDigits(test.input)
		if result != test.expected {
			t.Errorf("Filter RxMatch failed for input %q: expected %v, got %v", test.input, test.expected, result)
		}
	}
}
