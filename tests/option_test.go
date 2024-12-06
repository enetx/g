package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g"
)

func TestOptionOf(t *testing.T) {
	t.Run("Returns Some when ok is true", func(t *testing.T) {
		value := 42
		ok := true

		option := OptionOf(value, ok)

		if option.IsNone() {
			t.Errorf("Expected Some, but got None")
		}
	})

	t.Run("Returns None when ok is false", func(t *testing.T) {
		value := 42
		ok := false

		option := OptionOf(value, ok)

		if option.IsSome() {
			t.Errorf("Expected None, but got Some")
		}
	})

	t.Run("Works with different types", func(t *testing.T) {
		strValue := "test"
		strOption := OptionOf(strValue, true)
		if strOption.IsNone() {
			t.Errorf("Expected Some for string value, but got None")
		}

		floatValue := 3.14
		floatOption := OptionOf(floatValue, false)
		if floatOption.IsSome() {
			t.Errorf("Expected None for float value, but got Some")
		}
	})
}

func TestOptionUnwrapOr(t *testing.T) {
	fn := func(x int) Option[int] {
		if x > 10 {
			return Some(x)
		}
		return None[int]()
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
// 	option1 := Some(42)
// 	result1 := option1.Expect("Expected Some, got None")
// 	expected1 := 42

// 	if result1 != expected1 {
// 		t.Errorf("Test 1: Expected %d, got %d", expected1, result1)
// 	}

// 	// Test 2: Expecting panic from None
// 	option2 := None[int]()
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("Test 2: The code did not panic")
// 		}
// 	}()
// 	option2.Expect("Expected Some, got None")
// }

func TestOptionThen(t *testing.T) {
	// Test 1: Applying fn to Some
	option1 := Some(5)

	fn1 := func(x int) Option[int] {
		return Some(x * 2)
	}

	result1 := option1.Then(fn1)
	expected1 := Some(10)

	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test 1: Expected %v, got %v", expected1, result1)
	}

	// Test 2: Returning same Option for None
	option2 := None[int]()

	fn2 := func(x int) Option[int] {
		return Some(x * 2)
	}

	result2 := option2.Then(fn2)

	if result2.IsSome() {
		t.Errorf("Test 2: Expected None, got Some")
	}
}

// func TestOptionUnwrap(t *testing.T) {
// 	// Test 1: Unwrapping Some
// 	option1 := Some(42)
// 	result1 := option1.Unwrap()
// 	expected1 := 42

// 	if result1 != expected1 {
// 		t.Errorf("Test 1: Expected %d, got %d", expected1, result1)
// 	}

// 	// Test 2: Unwrapping None, should panic
// 	option2 := None[int]()
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("Test 2: The code did not panic")
// 		}
// 	}()

// 	option2.Unwrap()
// }

func TestTransformOption(t *testing.T) {
	// Test 1: Mapping over Some value
	option1 := Some(5)

	fn1 := func(x int) Option[int] {
		return Some(x * 2)
	}

	result1 := TransformOption(option1, fn1)
	expected1 := Some(10)

	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test 1: Expected %v, got %v", expected1, result1)
	}

	// Test 2: Mapping over None value
	option2 := None[int]()

	fn2 := func(x int) Option[int] {
		return Some(x * 2)
	}

	result2 := TransformOption(option2, fn2)

	if result2.IsSome() {
		t.Errorf("Test 2: Expected None, got Some")
	}
}
