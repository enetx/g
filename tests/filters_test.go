package g_test

import (
	"testing"

	"github.com/enetx/g/f"
)

func TestComparable(t *testing.T) {
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
			if got := f.Comparable(tt.value); got != tt.want {
				t.Errorf("Comparable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZero(t *testing.T) {
	// Testing Zero function with integers
	if !f.Zero(0) {
		t.Errorf("Zero(0) returned false, expected true")
	}

	if f.Zero(5) {
		t.Errorf("Zero(5) returned true, expected false")
	}

	// Testing Zero function with floats
	if !f.Zero(0.0) {
		t.Errorf("Zero(0.0) returned false, expected true")
	}

	if f.Zero(3.14) {
		t.Errorf("Zero(3.14) returned true, expected false")
	}

	// Testing Zero function with strings
	if !f.Zero("") {
		t.Errorf("Zero(\"\") returned false, expected true")
	}

	if f.Zero("hello") {
		t.Errorf("Zero(\"hello\") returned true, expected false")
	}
}

func TestEven(t *testing.T) {
	// Testing Even function with positive even integers
	if !f.Even(2) {
		t.Errorf("Even(2) returned false, expected true")
	}

	if !f.Even(100) {
		t.Errorf("Even(100) returned false, expected true")
	}

	// Testing Even function with positive odd integers
	if f.Even(3) {
		t.Errorf("Even(3) returned true, expected false")
	}

	if f.Even(101) {
		t.Errorf("Even(101) returned true, expected false")
	}

	// Testing Even function with negative even integers
	if !f.Even(-2) {
		t.Errorf("Even(-2) returned false, expected true")
	}

	if !f.Even(-100) {
		t.Errorf("Even(-100) returned false, expected true")
	}

	// Testing Even function with negative odd integers
	if f.Even(-3) {
		t.Errorf("Even(-3) returned true, expected false")
	}

	if f.Even(-101) {
		t.Errorf("Even(-101) returned true, expected false")
	}
}

func TestOdd(t *testing.T) {
	// Testing Odd function with positive odd integers
	if !f.Odd(3) {
		t.Errorf("Odd(3) returned false, expected true")
	}

	if !f.Odd(101) {
		t.Errorf("Odd(101) returned false, expected true")
	}

	// Testing Odd function with positive even integers
	if f.Odd(2) {
		t.Errorf("Odd(2) returned true, expected false")
	}

	if f.Odd(100) {
		t.Errorf("Odd(100) returned true, expected false")
	}

	// Testing Odd function with negative odd integers
	if !f.Odd(-3) {
		t.Errorf("Odd(-3) returned false, expected true")
	}

	if !f.Odd(-101) {
		t.Errorf("Odd(-101) returned false, expected true")
	}

	// Testing Odd function with negative even integers
	if f.Odd(-2) {
		t.Errorf("Odd(-2) returned true, expected false")
	}

	if f.Odd(-100) {
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
