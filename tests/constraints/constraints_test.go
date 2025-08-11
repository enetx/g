package g_test

import (
	"testing"

	"github.com/enetx/g/constraints"
)

// Test that Float interface works with generic functions
func TestFloatInterface(t *testing.T) {
	// Test that different float types work with the constraint
	if got := testFloatGeneric(float32(1.5)); got != 3.0 {
		t.Errorf("testFloatGeneric(float32(1.5)) = %v, want 3.0", got)
	}

	if got := testFloatGeneric(float64(2.5)); got != 5.0 {
		t.Errorf("testFloatGeneric(float64(2.5)) = %v, want 5.0", got)
	}
}

// Test that Integer interface works with generic functions
func TestIntegerInterface(t *testing.T) {
	// Test that different integer types work with the constraint
	if got := testIntegerGeneric(int(5)); got != 6 {
		t.Errorf("testIntegerGeneric(int(5)) = %v, want 6", got)
	}

	if got := testIntegerGeneric(int64(10)); got != 11 {
		t.Errorf("testIntegerGeneric(int64(10)) = %v, want 11", got)
	}

	if got := testIntegerGeneric(uint(20)); got != 21 {
		t.Errorf("testIntegerGeneric(uint(20)) = %v, want 21", got)
	}
}

// Generic function tests to ensure constraints work in practice
func testFloatGeneric[T constraints.Float](val T) T {
	return val * 2
}

func testIntegerGeneric[T constraints.Integer](val T) T {
	return val + 1
}

func TestFloatGeneric(t *testing.T) {
	if got := testFloatGeneric(2.5); got != 5.0 {
		t.Errorf("testFloatGeneric(2.5) = %v, want 5.0", got)
	}

	if got := testFloatGeneric(float32(1.5)); got != 3.0 {
		t.Errorf("testFloatGeneric(1.5) = %v, want 3.0", got)
	}
}

func TestIntegerGeneric(t *testing.T) {
	if got := testIntegerGeneric(5); got != 6 {
		t.Errorf("testIntegerGeneric(5) = %v, want 6", got)
	}

	if got := testIntegerGeneric(int64(10)); got != 11 {
		t.Errorf("testIntegerGeneric(10) = %v, want 11", got)
	}

	if got := testIntegerGeneric(uint(20)); got != 21 {
		t.Errorf("testIntegerGeneric(20) = %v, want 21", got)
	}
}
