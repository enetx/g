package g_test

import (
	"testing"

	"github.com/enetx/g/cmp"
)

func TestMinBy(t *testing.T) {
	// Test case 1: Minimum integer
	minInt := cmp.MinBy(cmp.Cmp, 3, 1, 4, 2, 5)
	expectedMinInt := 1
	if minInt != expectedMinInt {
		t.Errorf("MinBy(IntCompare, 3, 1, 4, 2, 5) = %d; want %d", minInt, expectedMinInt)
	}

	// Test case 2: Minimum string
	minString := cmp.MinBy(cmp.Cmp, "banana", "apple", "orange")
	expectedMinString := "apple"
	if minString != expectedMinString {
		t.Errorf("MinBy(StringCompare, \"banana\", \"apple\", \"orange\") = %s; want %s", minString, expectedMinString)
	}
}

func TestMaxBy(t *testing.T) {
	// Test case 1: Maximum integer
	maxInt := cmp.MaxBy(cmp.Cmp, 3, 1, 4, 2, 5)
	expectedMaxInt := 5
	if maxInt != expectedMaxInt {
		t.Errorf("MaxBy(IntCompare, 3, 1, 4, 2, 5) = %d; want %d", maxInt, expectedMaxInt)
	}

	// Test case 2: Maximum string
	maxString := cmp.MaxBy(cmp.Cmp, "banana", "apple", "orange")
	expectedMaxString := "orange"
	if maxString != expectedMaxString {
		t.Errorf("MaxBy(StringCompare, \"banana\", \"apple\", \"orange\") = %s; want %s", maxString, expectedMaxString)
	}
}

func TestOrderingIsLt(t *testing.T) {
	if !cmp.Less.IsLt() {
		t.Errorf("Expected Less to be less than other values")
	}
	if cmp.Equal.IsLt() {
		t.Errorf("Expected Equal not to be less than other values")
	}
	if cmp.Greater.IsLt() {
		t.Errorf("Expected Greater not to be less than other values")
	}
}

func TestOrderingIsEq(t *testing.T) {
	if !cmp.Equal.IsEq() {
		t.Errorf("Expected Equal to be equal to itself")
	}
	if cmp.Less.IsEq() {
		t.Errorf("Expected Less not to be equal to other values")
	}
	if cmp.Greater.IsEq() {
		t.Errorf("Expected Greater not to be equal to other values")
	}
}

func TestOrderingIsGt(t *testing.T) {
	if !cmp.Greater.IsGt() {
		t.Errorf("Expected Greater to be greater than other values")
	}
	if cmp.Equal.IsGt() {
		t.Errorf("Expected Equal not to be greater than other values")
	}
	if cmp.Less.IsGt() {
		t.Errorf("Expected Less not to be greater than other values")
	}
}

func TestMaxAndMin(t *testing.T) {
	// Test cases for Max function
	maxInt := cmp.Max(1, 2, 3, 4, 5)
	if maxInt != 5 {
		t.Errorf("cmp.Max(1, 2, 3, 4, 5) = %d; want 5", maxInt)
	}

	maxFloat := cmp.Max(1.5, 2.5, 3.5, 4.5, 5.5)
	if maxFloat != 5.5 {
		t.Errorf("cmp.Max(1.5, 2.5, 3.5, 4.5, 5.5) = %f; want 5.5", maxFloat)
	}

	maxempt := cmp.Max[int]()
	if maxempt != 0 {
		t.Errorf("cmp.Max() = %d; want 0", maxempt)
	}

	// Test cases for Min function
	minInt := cmp.Min(5, 4, 3, 2, 1)
	if minInt != 1 {
		t.Errorf("cmp.Min(5, 4, 3, 2, 1) = %d; want 1", minInt)
	}

	minFloat := cmp.Min(5.5, 4.5, 3.5, 2.5, 1.5)
	if minFloat != 1.5 {
		t.Errorf("cmp.Min(5.5, 4.5, 3.5, 2.5, 1.5) = %f; want 1.5", minFloat)
	}

	minempt := cmp.Min[int]()
	if minempt != 0 {
		t.Errorf("cmp.Min() = %d; want 0", minempt)
	}
}

func TestOrderingThen(t *testing.T) {
	tests := []struct {
		name     string
		o        cmp.Ordering
		other    cmp.Ordering
		expected cmp.Ordering
	}{
		{"Non-zero receiver", cmp.Ordering(2), cmp.Ordering(3), cmp.Ordering(2)},
		{"Zero receiver", cmp.Ordering(0), cmp.Ordering(3), cmp.Ordering(3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.o.Then(tt.other)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestOrderingCmp(t *testing.T) {
	tests := []struct {
		name     string
		x        cmp.Ordering
		y        cmp.Ordering
		expected cmp.Ordering
	}{
		{"x < y", cmp.Ordering(2), cmp.Ordering(3), cmp.Ordering(-1)},
		{"x = y", cmp.Ordering(2), cmp.Ordering(2), cmp.Ordering(0)},
		{"x > y", cmp.Ordering(3), cmp.Ordering(2), cmp.Ordering(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmp.Cmp(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestOrderingReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    cmp.Ordering
		expected cmp.Ordering
	}{
		{"Reverse_Less", cmp.Less, cmp.Greater},
		{"Reverse_Equal", cmp.Equal, cmp.Equal},
		{"Reverse_Greater", cmp.Greater, cmp.Less},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Reverse()
			if result != tt.expected {
				t.Errorf("Expected reverse of %v to be %v, but got %v", tt.input, tt.expected, result)
			}
		})
	}
}

func TestOrderingString(t *testing.T) {
	tests := []struct {
		name     string
		input    cmp.Ordering
		expected string
	}{
		{"Less", cmp.Less, "Less"},
		{"Equal", cmp.Equal, "Equal"},
		{"Greater", cmp.Greater, "Greater"},
		{"Unknown", 10, "Unknown Ordering value: 10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("Expected %s for Ordering value %d, but got %s", tt.expected, tt.input, result)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	t.Run("ReverseInt", func(t *testing.T) {
		// Test that Reverse inverts the comparison
		tests := []struct {
			name     string
			x, y     int
			expected cmp.Ordering
		}{
			{"x < y should be Greater", 1, 2, cmp.Greater},
			{"x = y should be Equal", 2, 2, cmp.Equal},
			{"x > y should be Less", 3, 2, cmp.Less},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := cmp.Reverse(tt.x, tt.y)
				if result != tt.expected {
					t.Errorf("cmp.Reverse(%d, %d) = %v; want %v", tt.x, tt.y, result, tt.expected)
				}
			})
		}
	})

	t.Run("ReverseString", func(t *testing.T) {
		tests := []struct {
			name     string
			x, y     string
			expected cmp.Ordering
		}{
			{"'apple' vs 'banana' reversed", "apple", "banana", cmp.Greater},
			{"'banana' vs 'apple' reversed", "banana", "apple", cmp.Less},
			{"'apple' vs 'apple' reversed", "apple", "apple", cmp.Equal},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := cmp.Reverse(tt.x, tt.y)
				if result != tt.expected {
					t.Errorf("cmp.Reverse(%q, %q) = %v; want %v", tt.x, tt.y, result, tt.expected)
				}
			})
		}
	})

	t.Run("ReverseFloat", func(t *testing.T) {
		tests := []struct {
			name     string
			x, y     float64
			expected cmp.Ordering
		}{
			{"1.5 vs 2.5 reversed", 1.5, 2.5, cmp.Greater},
			{"2.5 vs 1.5 reversed", 2.5, 1.5, cmp.Less},
			{"2.5 vs 2.5 reversed", 2.5, 2.5, cmp.Equal},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := cmp.Reverse(tt.x, tt.y)
				if result != tt.expected {
					t.Errorf("cmp.Reverse(%f, %f) = %v; want %v", tt.x, tt.y, result, tt.expected)
				}
			})
		}
	})
}

func TestReverseWithSorting(t *testing.T) {
	// Test Reverse function with actual sorting scenarios
	t.Run("IntegerSorting", func(t *testing.T) {
		// Simulate what would happen in sorting
		slice := []int{1, 3, 2, 5, 4}

		// Test that Reverse gives opposite results to Cmp
		for i := 0; i < len(slice)-1; i++ {
			for j := i + 1; j < len(slice); j++ {
				normalCmp := cmp.Cmp(slice[i], slice[j])
				reverseCmp := cmp.Reverse(slice[i], slice[j])

				// Reverse should be the opposite of normal comparison
				if normalCmp.IsLt() && !reverseCmp.IsGt() {
					t.Errorf("Expected Reverse to be Greater when Cmp is Less for %d vs %d", slice[i], slice[j])
				}
				if normalCmp.IsGt() && !reverseCmp.IsLt() {
					t.Errorf("Expected Reverse to be Less when Cmp is Greater for %d vs %d", slice[i], slice[j])
				}
				if normalCmp.IsEq() && !reverseCmp.IsEq() {
					t.Errorf("Expected Reverse to be Equal when Cmp is Equal for %d vs %d", slice[i], slice[j])
				}
			}
		}
	})
}

func TestReverseConsistency(t *testing.T) {
	// Test that Reverse is consistent with Ordering.Reverse()
	t.Run("ConsistentWithOrderingReverse", func(t *testing.T) {
		tests := []struct {
			x, y int
		}{
			{1, 2},
			{2, 1},
			{2, 2},
			{-1, 5},
			{100, -50},
		}

		for _, tt := range tests {
			cmpResult := cmp.Cmp(tt.x, tt.y)
			reverseFunc := cmp.Reverse(tt.x, tt.y)
			reverseMethod := cmpResult.Reverse()

			if reverseFunc != reverseMethod {
				t.Errorf("cmp.Reverse(%d, %d) = %v, but cmp.Cmp(%d, %d).Reverse() = %v",
					tt.x, tt.y, reverseFunc, tt.x, tt.y, reverseMethod)
			}
		}
	})
}
