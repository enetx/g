package g_test

import (
	"testing"

	"github.com/enetx/g/f"
)

func TestIsZero(t *testing.T) {
	// Testing IsZero function with integers
	if !f.IsZero(0) {
		t.Errorf("IsZero(0) returned false, expected true")
	}

	if f.IsZero(5) {
		t.Errorf("IsZero(5) returned true, expected false")
	}

	// Testing IsZero function with floats
	if !f.IsZero(0.0) {
		t.Errorf("IsZero(0.0) returned false, expected true")
	}

	if f.IsZero(3.14) {
		t.Errorf("IsZero(3.14) returned true, expected false")
	}

	// Testing IsZero function with strings
	if !f.IsZero("") {
		t.Errorf("IsZero(\"\") returned false, expected true")
	}

	if f.IsZero("hello") {
		t.Errorf("IsZero(\"hello\") returned true, expected false")
	}
}

func TestIsEven(t *testing.T) {
	// Testing IsEven function with positive even integers
	if !f.IsEven(2) {
		t.Errorf("IsEven(2) returned false, expected true")
	}

	if !f.IsEven(100) {
		t.Errorf("IsEven(100) returned false, expected true")
	}

	// Testing IsEven function with positive odd integers
	if f.IsEven(3) {
		t.Errorf("IsEven(3) returned true, expected false")
	}

	if f.IsEven(101) {
		t.Errorf("IsEven(101) returned true, expected false")
	}

	// Testing IsEven function with negative even integers
	if !f.IsEven(-2) {
		t.Errorf("IsEven(-2) returned false, expected true")
	}

	if !f.IsEven(-100) {
		t.Errorf("IsEven(-100) returned false, expected true")
	}

	// Testing IsEven function with negative odd integers
	if f.IsEven(-3) {
		t.Errorf("IsEven(-3) returned true, expected false")
	}

	if f.IsEven(-101) {
		t.Errorf("IsEven(-101) returned true, expected false")
	}
}

func TestIsOdd(t *testing.T) {
	// Testing IsOdd function with positive odd integers
	if !f.IsOdd(3) {
		t.Errorf("IsOdd(3) returned false, expected true")
	}

	if !f.IsOdd(101) {
		t.Errorf("IsOdd(101) returned false, expected true")
	}

	// Testing IsOdd function with positive even integers
	if f.IsOdd(2) {
		t.Errorf("IsOdd(2) returned true, expected false")
	}

	if f.IsOdd(100) {
		t.Errorf("IsOdd(100) returned true, expected false")
	}

	// Testing IsOdd function with negative odd integers
	if !f.IsOdd(-3) {
		t.Errorf("IsOdd(-3) returned false, expected true")
	}

	if !f.IsOdd(-101) {
		t.Errorf("IsOdd(-101) returned false, expected true")
	}

	// Testing IsOdd function with negative even integers
	if f.IsOdd(-2) {
		t.Errorf("IsOdd(-2) returned true, expected false")
	}

	if f.IsOdd(-100) {
		t.Errorf("IsOdd(-100) returned true, expected false")
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

func TestEqDeep(t *testing.T) {
	// Test case 1: Test deep equality of slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	if !f.EqDeep(slice1)(slice2) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test deep inequality of maps
	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "c": 3}
	if f.EqDeep(map1)(map2) {
		t.Errorf("Test case 2: Expected false, got true")
	}
}

func TestNeDeep(t *testing.T) {
	// Test case 1: Test deep inequality of slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 4}
	if !f.NeDeep(slice1)(slice2) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test deep equality of structs
	type Person struct {
		Name string
		Age  int
	}

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Bob", Age: 25}
	if !f.NeDeep(person1)(person2) {
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

func TestGtEq(t *testing.T) {
	// Test case 1: Test greater than or equal to comparison of integers
	if !f.GtEq(10)(10) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test greater than or equal to comparison of floats
	if !f.GtEq(3.14)(3.14) {
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

func TestLtEq(t *testing.T) {
	// Test case 1: Test less than or equal to comparison of integers
	if !f.LtEq(5)(5) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test less than or equal to comparison of floats
	if !f.LtEq(3.14)(3.14) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}
