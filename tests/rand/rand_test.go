package g_test

import (
	"testing"

	"github.com/enetx/g/rand"
)

func TestN(t *testing.T) {
	tests := []struct {
		name string
		max  int
		want func(int) bool // validation function
	}{
		{
			"positive max",
			10,
			func(result int) bool { return result >= 0 && result < 10 },
		},
		{
			"max = 1",
			1,
			func(result int) bool { return result == 0 },
		},
		{
			"zero max",
			0,
			func(result int) bool { return result == 0 }, // should be treated as max = 1
		},
		{
			"negative max",
			-5,
			func(result int) bool { return result == 0 }, // should be treated as max = 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test multiple times to ensure randomness works correctly
			for range 100 {
				got := rand.N(tt.max)
				if !tt.want(got) {
					t.Errorf("N(%v) = %v, validation failed", tt.max, got)
					break
				}
			}
		})
	}
}

func TestN_Distribution(t *testing.T) {
	// Test that the function produces different values over multiple calls
	max := 100
	results := make(map[int]bool)

	// Generate 1000 random numbers
	for range 1000 {
		result := rand.N(max)
		results[result] = true

		// Verify range
		if result < 0 || result >= max {
			t.Errorf("N(%v) = %v, out of expected range [0, %v)", max, result, max)
		}
	}

	// We should have gotten multiple different values
	// This is probabilistically almost certain with 1000 calls and max=100
	if len(results) < 10 {
		t.Errorf("Expected multiple different random values, got only %d unique values", len(results))
	}
}

func TestN_WithDifferentTypes(t *testing.T) {
	// Test with different integer types
	t.Run("int8", func(t *testing.T) {
		result := rand.N(int8(10))
		if result < 0 || result >= 10 {
			t.Errorf("N(int8(10)) = %v, out of range", result)
		}
	})

	t.Run("int16", func(t *testing.T) {
		result := rand.N(int16(100))
		if result < 0 || result >= 100 {
			t.Errorf("N(int16(100)) = %v, out of range", result)
		}
	})

	t.Run("int32", func(t *testing.T) {
		result := rand.N(int32(1000))
		if result < 0 || result >= 1000 {
			t.Errorf("N(int32(1000)) = %v, out of range", result)
		}
	})

	t.Run("int64", func(t *testing.T) {
		result := rand.N(int64(10000))
		if result < 0 || result >= 10000 {
			t.Errorf("N(int64(10000)) = %v, out of range", result)
		}
	})

	t.Run("uint", func(t *testing.T) {
		result := rand.N(uint(50))
		if result >= 50 {
			t.Errorf("N(uint(50)) = %v, out of range", result)
		}
	})
}

func TestN_EdgeCases(t *testing.T) {
	t.Run("max = 1", func(t *testing.T) {
		// Should always return 0
		for range 10 {
			result := rand.N(1)
			if result != 0 {
				t.Errorf("N(1) = %v, expected 0", result)
			}
		}
	})

	t.Run("large max", func(t *testing.T) {
		max := 1000000
		result := rand.N(max)
		if result < 0 || result >= max {
			t.Errorf("N(%v) = %v, out of range", max, result)
		}
	})
}
