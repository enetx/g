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

// Generic helpers exercising the remaining constraints.
func testSignedGeneric[T constraints.Signed](v T) T   { return v - 1 }
func testUnsignedGeneric[T constraints.Unsigned](v T) T { return v + 1 }
func testComplexGeneric[T constraints.Complex](v T) T  { return v * 2 }
func testNumberGeneric[T constraints.Number](v T) T    { return v + 1 }

func testOrderedGeneric[T constraints.Ordered](a, b T) bool { return a < b }

func TestSignedConstraint(t *testing.T) {
	if got := testSignedGeneric(int8(5)); got != 4 {
		t.Errorf("testSignedGeneric(int8(5)) = %v, want 4", got)
	}
	if got := testSignedGeneric(int64(-3)); got != -4 {
		t.Errorf("testSignedGeneric(int64(-3)) = %v, want -4", got)
	}
}

func TestUnsignedConstraint(t *testing.T) {
	if got := testUnsignedGeneric(uint8(5)); got != 6 {
		t.Errorf("testUnsignedGeneric(uint8(5)) = %v, want 6", got)
	}
	if got := testUnsignedGeneric(uint(0)); got != 1 {
		t.Errorf("testUnsignedGeneric(uint(0)) = %v, want 1", got)
	}
}

func TestComplexConstraint(t *testing.T) {
	if got := testComplexGeneric(complex64(1 + 2i)); got != complex64(2+4i) {
		t.Errorf("testComplexGeneric(1+2i) = %v, want (2+4i)", got)
	}
	if got := testComplexGeneric(complex128(3 + 1i)); got != complex128(6+2i) {
		t.Errorf("testComplexGeneric(3+1i) = %v, want (6+2i)", got)
	}
}

func TestNumberConstraint(t *testing.T) {
	if got := testNumberGeneric(5); got != 6 {
		t.Errorf("testNumberGeneric(5) = %v, want 6", got)
	}
	if got := testNumberGeneric(2.5); got != 3.5 {
		t.Errorf("testNumberGeneric(2.5) = %v, want 3.5", got)
	}
	if got := testNumberGeneric(uint(7)); got != 8 {
		t.Errorf("testNumberGeneric(uint(7)) = %v, want 8", got)
	}
}

func TestOrderedConstraint(t *testing.T) {
	if !testOrderedGeneric(1, 2) {
		t.Errorf("testOrderedGeneric(1, 2) = false, want true")
	}
	if !testOrderedGeneric("apple", "banana") {
		t.Errorf("testOrderedGeneric(apple, banana) = false, want true")
	}
	if testOrderedGeneric(2.5, 1.5) {
		t.Errorf("testOrderedGeneric(2.5, 1.5) = true, want false")
	}
}
