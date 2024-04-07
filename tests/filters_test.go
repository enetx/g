package g_test

import (
	"testing"

	"github.com/enetx/g/filters"
)

func TestIsZero(t *testing.T) {
	// Testing IsZero function with integers
	if !filters.IsZero(0) {
		t.Errorf("IsZero(0) returned false, expected true")
	}

	if filters.IsZero(5) {
		t.Errorf("IsZero(5) returned true, expected false")
	}

	// Testing IsZero function with floats
	if !filters.IsZero(0.0) {
		t.Errorf("IsZero(0.0) returned false, expected true")
	}

	if filters.IsZero(3.14) {
		t.Errorf("IsZero(3.14) returned true, expected false")
	}

	// Testing IsZero function with strings
	if !filters.IsZero("") {
		t.Errorf("IsZero(\"\") returned false, expected true")
	}

	if filters.IsZero("hello") {
		t.Errorf("IsZero(\"hello\") returned true, expected false")
	}
}

func TestIsEven(t *testing.T) {
	// Testing IsEven function with positive even integers
	if !filters.IsEven(2) {
		t.Errorf("IsEven(2) returned false, expected true")
	}

	if !filters.IsEven(100) {
		t.Errorf("IsEven(100) returned false, expected true")
	}

	// Testing IsEven function with positive odd integers
	if filters.IsEven(3) {
		t.Errorf("IsEven(3) returned true, expected false")
	}

	if filters.IsEven(101) {
		t.Errorf("IsEven(101) returned true, expected false")
	}

	// Testing IsEven function with negative even integers
	if !filters.IsEven(-2) {
		t.Errorf("IsEven(-2) returned false, expected true")
	}

	if !filters.IsEven(-100) {
		t.Errorf("IsEven(-100) returned false, expected true")
	}

	// Testing IsEven function with negative odd integers
	if filters.IsEven(-3) {
		t.Errorf("IsEven(-3) returned true, expected false")
	}

	if filters.IsEven(-101) {
		t.Errorf("IsEven(-101) returned true, expected false")
	}
}

func TestIsOdd(t *testing.T) {
	// Testing IsOdd function with positive odd integers
	if !filters.IsOdd(3) {
		t.Errorf("IsOdd(3) returned false, expected true")
	}

	if !filters.IsOdd(101) {
		t.Errorf("IsOdd(101) returned false, expected true")
	}

	// Testing IsOdd function with positive even integers
	if filters.IsOdd(2) {
		t.Errorf("IsOdd(2) returned true, expected false")
	}

	if filters.IsOdd(100) {
		t.Errorf("IsOdd(100) returned true, expected false")
	}

	// Testing IsOdd function with negative odd integers
	if !filters.IsOdd(-3) {
		t.Errorf("IsOdd(-3) returned false, expected true")
	}

	if !filters.IsOdd(-101) {
		t.Errorf("IsOdd(-101) returned false, expected true")
	}

	// Testing IsOdd function with negative even integers
	if filters.IsOdd(-2) {
		t.Errorf("IsOdd(-2) returned true, expected false")
	}

	if filters.IsOdd(-100) {
		t.Errorf("IsOdd(-100) returned true, expected false")
	}
}
