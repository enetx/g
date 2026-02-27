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

func TestOptionFilter(t *testing.T) {
	opt := Some(10)

	res := opt.Filter(func(x int) bool {
		return x > 5
	})

	if res.IsNone() {
		t.Error("Expected Some(10)")
	}

	res2 := opt.Filter(func(x int) bool {
		return x > 20
	})

	if res2.IsSome() {
		t.Error("Expected None")
	}
}

func TestOptionOr(t *testing.T) {
	a := Some(1)
	b := Some(2)
	n := None[int]()

	if a.Or(b).Unwrap() != 1 {
		t.Error("Or should return first Some")
	}

	if n.Or(b).Unwrap() != 2 {
		t.Error("Or should return alternative")
	}
}

func TestOptionOrElse(t *testing.T) {
	a := Some(1)
	n := None[int]()

	called := false

	res := a.OrElse(func() Option[int] {
		called = true
		return Some(2)
	})

	if called {
		t.Error("OrElse should not call fn when Some")
	}
	if res.Unwrap() != 1 {
		t.Error("Unexpected value")
	}

	res2 := n.OrElse(func() Option[int] {
		called = true
		return Some(3)
	})

	if !called {
		t.Error("OrElse should call fn when None")
	}
	if res2.Unwrap() != 3 {
		t.Error("Unexpected value")
	}
}

func TestOptionIsSomeAnd(t *testing.T) {
	opt := Some(10)

	if !opt.IsSomeAnd(func(x int) bool { return x == 10 }) {
		t.Error("Expected true")
	}

	if opt.IsSomeAnd(func(x int) bool { return x > 20 }) {
		t.Error("Expected false")
	}

	if None[int]().IsSomeAnd(func(int) bool { return true }) {
		t.Error("None should return false")
	}
}

func TestOptionInsert(t *testing.T) {
	opt := None[int]()

	ptr := opt.Insert(5)

	if !opt.IsSome() || *ptr != 5 {
		t.Error("Insert failed")
	}
}

func TestOptionGetOrInsert(t *testing.T) {
	opt := None[int]()

	ptr := opt.GetOrInsert(10)

	if *ptr != 10 {
		t.Error("Expected inserted value 10")
	}

	ptr2 := opt.GetOrInsert(20)

	if *ptr2 != 10 {
		t.Error("Should not overwrite existing value")
	}
}

func TestOptionGetOrInsertWith(t *testing.T) {
	opt := None[int]()

	called := false

	ptr := opt.GetOrInsertWith(func() int {
		called = true
		return 30
	})

	if !called || *ptr != 30 {
		t.Error("Lazy insert failed")
	}

	called = false

	opt.GetOrInsertWith(func() int {
		called = true
		return 40
	})

	if called {
		t.Error("Should not call fn when Some")
	}
}

func TestOptionTake(t *testing.T) {
	opt := Some(50)

	old := opt.Take()

	if old.IsNone() || old.Unwrap() != 50 {
		t.Error("Take should return old value")
	}

	if opt.IsSome() {
		t.Error("Option should become None after Take")
	}
}

func TestOptionReplace(t *testing.T) {
	opt := Some(1)

	old := opt.Replace(2)

	if old.Unwrap() != 1 {
		t.Error("Replace should return old value")
	}

	if opt.Unwrap() != 2 {
		t.Error("Replace should set new value")
	}
}

func TestOptionOkOr(t *testing.T) {
	opt := Some(5)

	res := opt.OkOr(errors.New("err"))

	if res.IsErr() || res.Unwrap() != 5 {
		t.Error("Expected Ok(5)")
	}

	none := None[int]()

	res2 := none.OkOr(errors.New("err"))

	if res2.IsOk() {
		t.Error("Expected Err")
	}
}

func TestOptionOkOrElse(t *testing.T) {
	opt := None[int]()

	called := false

	res := opt.OkOrElse(func() error {
		called = true
		return errors.New("generated")
	})

	if !called {
		t.Error("Expected fn to be called")
	}

	if res.IsOk() {
		t.Error("Expected Err")
	}
}

func TestOptionPtr(t *testing.T) {
	opt := Some(5)

	ptr := opt.Ptr()

	if ptr == nil || *ptr != 5 {
		t.Error("Ptr failed")
	}

	if None[int]().Ptr() != nil {
		t.Error("Ptr of None should be nil")
	}
}

func TestOptionFromPtr(t *testing.T) {
	val := 10
	opt := OptionFromPtr(&val)

	if opt.IsNone() || opt.Unwrap() != 10 {
		t.Error("FromPtr failed")
	}

	if OptionFromPtr[int](nil).IsSome() {
		t.Error("FromPtr(nil) should be None")
	}
}
