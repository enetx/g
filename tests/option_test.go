package g_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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

func TestOptionUnwrap(t *testing.T) {
	// Test unwrapping Some value
	some := Some(42)
	value := some.Unwrap()
	if value != 42 {
		t.Errorf("Unwrap() of Some(42) should return 42, got %d", value)
	}

	// Test unwrapping None should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Unwrap() of None should panic")
		}
	}()

	none := None[int]()
	none.Unwrap() // This should panic
}

func TestOptionExpect(t *testing.T) {
	// Test expect with Some value
	some := Some(42)
	value := some.Expect("should have value")
	if value != 42 {
		t.Errorf("Expect() of Some(42) should return 42, got %d", value)
	}

	// Test expect with None should panic with message
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expect() of None should panic")
		} else if !strings.Contains(fmt.Sprintf("%v", r), "test message") {
			t.Error("Expect() should panic with custom message")
		}
	}()

	none := None[int]()
	none.Expect("test message") // This should panic with message
}

func TestOptionResult(t *testing.T) {
	// Test Some to Result
	some := Some(42)
	result := some.Result(errors.New("error message"))
	if result.IsErr() {
		t.Error("Some should convert to Ok Result")
	}
	if result.Unwrap() != 42 {
		t.Errorf("Result value should be 42, got %d", result.Unwrap())
	}

	// Test None to Result
	none := None[int]()
	result2 := none.Result(errors.New("test error"))
	if result2.IsOk() {
		t.Error("None should convert to Err Result")
	}
}

func TestOptionString(t *testing.T) {
	// Test Some string representation
	some := Some(42)
	str := some.String()
	if !strings.Contains(str, "Some") || !strings.Contains(str, "42") {
		t.Errorf("Some string should contain 'Some' and value, got %s", str)
	}

	// Test None string representation
	none := None[int]()
	str2 := none.String()
	if !strings.Contains(str2, "None") {
		t.Errorf("None string should contain 'None', got %s", str2)
	}
}
