package g_test

import (
	"reflect"
	"testing"

	"github.com/enetx/g"
)

func TestOptionUnwrapOr(t *testing.T) {
	fn := func(x int) g.Option[int] {
		if x > 10 {
			return g.Some[int](x)
		}
		return g.None[int]()
	}

	result := fn(5).UnwrapOr(10)
	expected := 10

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}

	result = fn(11).UnwrapOr(10)
	expected = 11

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// func TestOptionExpect(t *testing.T) {
// 	// Test 1: Expecting value from Some
// 	option1 := g.Some(42)
// 	result1 := option1.Expect("Expected Some, got None")
// 	expected1 := 42

// 	if result1 != expected1 {
// 		t.Errorf("Test 1: Expected %d, got %d", expected1, result1)
// 	}

// 	// Test 2: Expecting panic from None
// 	option2 := g.None[int]()
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("Test 2: The code did not panic")
// 		}
// 	}()
// 	option2.Expect("Expected Some, got None")
// }

func TestOptionThen(t *testing.T) {
	// Test 1: Applying fn to Some
	option1 := g.Some(5)

	fn1 := func(x int) g.Option[int] {
		return g.Some(x * 2)
	}

	result1 := option1.Then(fn1)
	expected1 := g.Some(10)

	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test 1: Expected %v, got %v", expected1, result1)
	}

	// Test 2: Returning same Option for None
	option2 := g.None[int]()

	fn2 := func(x int) g.Option[int] {
		return g.Some(x * 2)
	}

	result2 := option2.Then(fn2)

	if result2.IsSome() {
		t.Errorf("Test 2: Expected None, got Some")
	}
}

// func TestOptionUnwrap(t *testing.T) {
// 	// Test 1: Unwrapping Some
// 	option1 := g.Some(42)
// 	result1 := option1.Unwrap()
// 	expected1 := 42

// 	if result1 != expected1 {
// 		t.Errorf("Test 1: Expected %d, got %d", expected1, result1)
// 	}

// 	// Test 2: Unwrapping None, should panic
// 	option2 := g.None[int]()
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("Test 2: The code did not panic")
// 		}
// 	}()

// 	option2.Unwrap()
// }

func TestMapOption(t *testing.T) {
	// Test 1: Mapping over Some value
	option1 := g.Some(5)

	fn1 := func(x int) g.Option[int] {
		return g.Some(x * 2)
	}

	result1 := g.MapOption(option1, fn1)
	expected1 := g.Some(10)

	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test 1: Expected %v, got %v", expected1, result1)
	}

	// Test 2: Mapping over None value
	option2 := g.None[int]()

	fn2 := func(x int) g.Option[int] {
		return g.Some(x * 2)
	}

	result2 := g.MapOption(option2, fn2)

	if result2.IsSome() {
		t.Errorf("Test 2: Expected None, got Some")
	}
}
