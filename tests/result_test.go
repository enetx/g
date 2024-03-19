package g_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/enetx/g"
)

func TestResultOf(t *testing.T) {
	// Test 1: Result with an error
	err1 := errors.New("something went wrong")
	result1 := g.ResultOf(0, err1)

	if result1.IsOk() {
		t.Errorf("Test 1: Expected error, got Ok")
	}

	if !reflect.DeepEqual(result1.Err(), err1) {
		t.Errorf("Test 1: Expected %v, got %v", err1, result1.Err())
	}

	// Test 2: Result with a value
	value2 := 42
	result2 := g.ResultOf(value2, nil)

	if result2.IsErr() {
		t.Errorf("Test 2: Expected Ok, got error")
	}

	if result2.Unwrap() != value2 {
		t.Errorf("Test 2: Expected %d, got %d", value2, result2.Unwrap())
	}
}

func TestResultMap(t *testing.T) {
	// Test 1: Mapping over a Result with a value
	result1 := g.Ok[int](5)

	fn1 := func(x int) g.Result[int] {
		return g.Ok(x * 2)
	}

	mappedResult1 := g.ResultMap(result1, fn1)

	if mappedResult1.IsErr() {
		t.Errorf("Test 1: Expected Ok, got error")
	}

	expectedValue1 := 10
	if mappedResult1.Unwrap() != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, mappedResult1.Unwrap())
	}

	// Test 2: Mapping over a Result with an error
	err2 := errors.New("some error")
	result2 := g.Err[int](err2)

	fn2 := func(x int) g.Result[int] {
		return g.Ok(x * 2)
	}

	mappedResult2 := g.ResultMap(result2, fn2)

	if mappedResult2.IsOk() {
		t.Errorf("Test 2: Expected error, got Ok")
	}

	if !reflect.DeepEqual(mappedResult2.Err(), err2) {
		t.Errorf("Test 2: Expected %v, got %v", err2, mappedResult2.Err())
	}
}

func TestResultOfMap(t *testing.T) {
	// Test 1: Mapping over a Result with a value
	result1 := g.Ok[int](5)

	fn1 := func(x int) (int, error) {
		return x * 2, nil
	}

	mappedResult1 := g.ResultOfMap(result1, fn1)

	if mappedResult1.IsErr() {
		t.Errorf("Test 1: Expected Ok, got error")
	}

	expectedValue1 := 10
	if mappedResult1.Unwrap() != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, mappedResult1.Unwrap())
	}

	// Test 2: Mapping over a Result with an error
	err2 := errors.New("some error")
	result2 := g.Err[int](err2)

	fn2 := func(x int) (int, error) {
		return x * 2, nil
	}

	mappedResult2 := g.ResultOfMap(result2, fn2)

	if mappedResult2.IsOk() {
		t.Errorf("Test 2: Expected error, got Ok")
	}

	if !reflect.DeepEqual(mappedResult2.Err(), err2) {
		t.Errorf("Test 2: Expected %v, got %v", err2, mappedResult2.Err())
	}
}

func TestResult(t *testing.T) {
	// Test 1: Result with a value
	result1 := g.Ok[int](42)
	value1, err1 := result1.Result()

	if err1 != nil {
		t.Errorf("Test 1: Expected nil error, got %v", err1)
	}

	expectedValue1 := 42
	if value1 != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, value1)
	}

	// Test 2: Result with an error
	err2 := errors.New("some error")
	result2 := g.Err[int](err2)

	_, err := result2.Result()
	if err == nil {
		t.Errorf("Test 2: Expected non-nil error, got nil")
	}

	if !reflect.DeepEqual(err2, err) {
		t.Errorf("Test 2: Expected %v, got %v", err2, err)
	}
}

// func TestResultUnwrap(t *testing.T) {
// 	// Test 1: Unwrapping Result with a value
// 	result1 := g.Ok[int](42)
// 	value1 := result1.Unwrap()

// 	expectedValue1 := 42
// 	if value1 != expectedValue1 {
// 		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, value1)
// 	}

// 	// Test 2: Unwrapping Result with an error, should panic
// 	err2 := errors.New("some error")
// 	result2 := g.Err[int](err2)
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("Test 2: The code did not panic")
// 		}
// 	}()

// 	result2.Unwrap()
// }

func TestResultUnwrapOr(t *testing.T) {
	// Test 1: Unwrapping Result with a value
	result1 := g.Ok[int](42)
	value1 := result1.UnwrapOr(10)

	expectedValue1 := 42
	if value1 != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, value1)
	}

	// Test 2: Unwrapping Result with an error, should return default value
	err2 := errors.New("some error")
	result2 := g.Err[int](err2)
	defaultValue2 := 10
	value2 := result2.UnwrapOr(defaultValue2)

	if value2 != defaultValue2 {
		t.Errorf("Test 2: Expected %d, got %d", defaultValue2, value2)
	}
}

// func TestResultExpect(t *testing.T) {
// 	// Test 1: Expecting value from Result with a value
// 	result1 := g.Ok[int](42)
// 	value1 := result1.Expect("Expected Some, got None")

// 	expectedValue1 := 42
// 	if value1 != expectedValue1 {
// 		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, value1)
// 	}

// 	// Test 2: Expecting panic from Result with an error
// 	err2 := errors.New("some error")
// 	result2 := g.Err[int](err2)
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("Test 2: The code did not panic")
// 		}
// 	}()

// 	result2.Expect("Expected Some, got None")
// }

func TestResultThen(t *testing.T) {
	// Test 1: Applying fn to Result with a value
	result1 := g.Ok[int](42)
	fn1 := func(x int) g.Result[int] {
		return g.Ok(x * 2)
	}
	newResult1 := result1.Then(fn1)

	if newResult1.IsErr() {
		t.Errorf("Test 1: Expected Ok, got error")
	}

	expectedValue1 := 84
	if newResult1.Unwrap() != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, newResult1.Unwrap())
	}

	// Test 2: Returning same Result for Result with an error
	err2 := errors.New("some error")
	result2 := g.Err[int](err2)
	fn2 := func(x int) g.Result[int] {
		return g.Ok(x * 2)
	}
	newResult2 := result2.Then(fn2)

	if newResult2.IsOk() {
		t.Errorf("Test 2: Expected error, got Ok")
	}

	if !reflect.DeepEqual(newResult2.Err(), err2) {
		t.Errorf("Test 2: Expected %v, got %v", err2, newResult2.Err())
	}
}

func TestResultThenOf(t *testing.T) {
	// Test 1: Applying fn to Result with a value
	result1 := g.Ok[int](42)
	fn1 := func(x int) (int, error) {
		return x * 2, nil
	}

	newResult1 := result1.ThenOf(fn1)

	if newResult1.IsErr() {
		t.Errorf("Test 1: Expected Ok, got error")
	}

	expectedValue1 := 84
	if newResult1.Unwrap() != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, newResult1.Unwrap())
	}

	// Test 2: Returning same Result for Result with an error
	err2 := errors.New("some error")
	result2 := g.Err[int](err2)
	fn2 := func(x int) (int, error) {
		return x * 2, nil
	}
	newResult2 := result2.ThenOf(fn2)

	if newResult2.IsOk() {
		t.Errorf("Test 2: Expected error, got Ok")
	}

	if !reflect.DeepEqual(newResult2.Err(), err2) {
		t.Errorf("Test 2: Expected %v, got %v", err2, newResult2.Err())
	}
}

func TestResultOption(t *testing.T) {
	// Test 1: Converting Result with a value to Option
	result1 := g.Ok[int](42)
	option1 := result1.Option()

	if option1.IsNone() {
		t.Errorf("Test 1: Expected Some, got None")
	}

	expectedValue1 := 42
	if option1.Some() != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, option1.Some())
	}

	// Test 2: Converting Result with an error to None
	err2 := errors.New("some error")
	result2 := g.Err[int](err2)
	option2 := result2.Option()

	if option2.IsSome() {
		t.Errorf("Test 2: Expected None, got Some")
	}
}
