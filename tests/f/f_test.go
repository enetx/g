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
		{"int is comparable", 42, true},
		{"string is comparable", "test", true},
		{"slice is not comparable", []int{1, 2, 3}, false},
		{"map is not comparable", map[string]int{"key": 1}, false},
		{"struct is comparable", struct{ Name string }{"test"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.IsComparable(tt.value)
			if got != tt.want {
				t.Errorf("IsComparable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"zero int", 0, true},
		{"non-zero int", 42, false},
		{"zero string", "", true},
		{"non-zero string", "test", false},
		{"zero float", 0.0, true},
		{"non-zero float", 3.14, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case int:
				got := f.IsZero(v)
				if got != tt.want {
					t.Errorf("IsZero(%v) = %v, want %v", v, got, tt.want)
				}
			case string:
				got := f.IsZero(v)
				if got != tt.want {
					t.Errorf("IsZero(%v) = %v, want %v", v, got, tt.want)
				}
			case float64:
				got := f.IsZero(v)
				if got != tt.want {
					t.Errorf("IsZero(%v) = %v, want %v", v, got, tt.want)
				}
			}
		})
	}
}

func TestIsEven(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"even number", 4, true},
		{"odd number", 3, false},
		{"zero is even", 0, true},
		{"negative even", -4, true},
		{"negative odd", -3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.IsEven(tt.value)
			if got != tt.want {
				t.Errorf("IsEven(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsOdd(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"odd number", 3, true},
		{"even number", 4, false},
		{"zero is not odd", 0, false},
		{"negative odd", -3, true},
		{"negative even", -4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.IsOdd(tt.value)
			if got != tt.want {
				t.Errorf("IsOdd(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	pattern, err := regexp.Compile(`^\d+$`)
	if err != nil {
		t.Fatal("Failed to compile regex pattern")
	}

	matchFunc := f.Match[string](pattern)

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"numbers only", "12345", true},
		{"contains letters", "123abc", false},
		{"empty string", "", false},
		{"single digit", "7", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchFunc(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	containsFunc := f.Contains("test")

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"contains substring", "this is a test", true},
		{"exact match", "test", true},
		{"doesn't contain", "hello world", false},
		{"empty string", "", false},
		{"case sensitive", "TEST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsFunc(tt.input)
			if got != tt.want {
				t.Errorf("Contains(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestContainsAnyChars(t *testing.T) {
	containsAnyFunc := f.ContainsAnyChars[string]("abc")

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"contains 'a'", "apple", true},
		{"contains 'b'", "banana", true},
		{"contains 'c'", "cat", true},
		{"contains multiple", "abc", true},
		{"doesn't contain any", "xyz", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsAnyFunc(tt.input)
			if got != tt.want {
				t.Errorf("ContainsAnyChars(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	startsWithFunc := f.StartsWith[string]("hello")

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"starts with prefix", "hello world", true},
		{"exact match", "hello", true},
		{"doesn't start with", "world hello", false},
		{"empty string", "", false},
		{"case sensitive", "Hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := startsWithFunc(tt.input)
			if got != tt.want {
				t.Errorf("StartsWith(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
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

func TestFilterRxMatch(t *testing.T) {
	regex := regexp.MustCompile(`^\d+$`)
	matchDigits := f.Match[string](regex)

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
