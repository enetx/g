package g_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/enetx/g"
)

func TestResultOf(t *testing.T) {
	// Test 1: Result with an error
	err1 := errors.New("something went wrong")
	result1 := ResultOf(0, err1)

	if result1.IsOk() {
		t.Errorf("Test 1: Expected error, got Ok")
	}

	if !reflect.DeepEqual(result1.Err(), err1) {
		t.Errorf("Test 1: Expected %v, got %v", err1, result1.Err())
	}

	// Test 2: Result with a value
	value2 := 42
	result2 := ResultOf(value2, nil)

	if result2.IsErr() {
		t.Errorf("Test 2: Expected Ok, got error")
	}

	if result2.Unwrap() != value2 {
		t.Errorf("Test 2: Expected %d, got %d", value2, result2.Unwrap())
	}
}

func TestTransformResult(t *testing.T) {
	// Test 1: Mapping over a Result with a value
	result1 := Ok(5)

	fn1 := func(x int) Result[int] {
		return Ok(x * 2)
	}

	mappedResult1 := TransformResult(result1, fn1)

	if mappedResult1.IsErr() {
		t.Errorf("Test 1: Expected Ok, got error")
	}

	expectedValue1 := 10
	if mappedResult1.Unwrap() != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, mappedResult1.Unwrap())
	}

	// Test 2: Mapping over a Result with an error
	err2 := errors.New("some error")
	result2 := Err[int](err2)

	fn2 := func(x int) Result[int] {
		return Ok(x * 2)
	}

	mappedResult2 := TransformResult(result2, fn2)

	if mappedResult2.IsOk() {
		t.Errorf("Test 2: Expected error, got Ok")
	}

	if !reflect.DeepEqual(mappedResult2.Err(), err2) {
		t.Errorf("Test 2: Expected %v, got %v", err2, mappedResult2.Err())
	}
}

func TestResultOfMap(t *testing.T) {
	// Test 1: Mapping over a Result with a value
	result1 := Ok(5)

	fn1 := func(x int) (int, error) {
		return x * 2, nil
	}

	mappedResult1 := TransformResultOf(result1, fn1)

	if mappedResult1.IsErr() {
		t.Errorf("Test 1: Expected Ok, got error")
	}

	expectedValue1 := 10
	if mappedResult1.Unwrap() != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, mappedResult1.Unwrap())
	}

	// Test 2: Mapping over a Result with an error
	err2 := errors.New("some error")
	result2 := Err[int](err2)

	fn2 := func(x int) (int, error) {
		return x * 2, nil
	}

	mappedResult2 := TransformResultOf(result2, fn2)

	if mappedResult2.IsOk() {
		t.Errorf("Test 2: Expected error, got Ok")
	}

	if !reflect.DeepEqual(mappedResult2.Err(), err2) {
		t.Errorf("Test 2: Expected %v, got %v", err2, mappedResult2.Err())
	}
}

func TestResult(t *testing.T) {
	// Test 1: Result with a value
	result1 := Ok(42)
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
	result2 := Err[int](err2)

	_, err := result2.Result()
	if err == nil {
		t.Errorf("Test 2: Expected non-nil error, got nil")
	}

	if !reflect.DeepEqual(err2, err) {
		t.Errorf("Test 2: Expected %v, got %v", err2, err)
	}
}

func TestResultUnwrap(t *testing.T) {
	// Test 1: Unwrapping Result with a value
	result1 := Ok[int](42)
	value1 := result1.Unwrap()

	expectedValue1 := 42
	if value1 != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, value1)
	}

	// Test 2: Unwrapping Result with an error, should panic
	err2 := errors.New("some error")
	result2 := Err[int](err2)

	// Test the panic case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Test 2: The code did not panic")
			}
		}()
		result2.Unwrap()
	}()
}

func TestResultUnwrapOr(t *testing.T) {
	// Test 1: Unwrapping Result with a value
	result1 := Ok[int](42)
	value1 := result1.UnwrapOr(10)

	expectedValue1 := 42
	if value1 != expectedValue1 {
		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, value1)
	}

	// Test 2: Unwrapping Result with an error, should return default value
	err2 := errors.New("some error")
	result2 := Err[int](err2)
	defaultValue2 := 10
	value2 := result2.UnwrapOr(defaultValue2)

	if value2 != defaultValue2 {
		t.Errorf("Test 2: Expected %d, got %d", defaultValue2, value2)
	}
}

// func TestResultExpect(t *testing.T) {
// 	// Test 1: Expecting value from Result with a value
// 	result1 := Ok[int](42)
// 	value1 := result1.Expect("Expected Some, got None")

// 	expectedValue1 := 42
// 	if value1 != expectedValue1 {
// 		t.Errorf("Test 1: Expected %d, got %d", expectedValue1, value1)
// 	}

// 	// Test 2: Expecting panic from Result with an error
// 	err2 := errors.New("some error")
// 	result2 := Err[int](err2)
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("Test 2: The code did not panic")
// 		}
// 	}()

// 	result2.Expect("Expected Some, got None")
// }

func TestResultThen(t *testing.T) {
	// Test 1: Applying fn to Result with a value
	result1 := Ok(42)
	fn1 := func(x int) Result[int] {
		return Ok(x * 2)
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
	result2 := Err[int](err2)
	fn2 := func(x int) Result[int] {
		return Ok(x * 2)
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
	result1 := Ok[int](42)
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
	result2 := Err[int](err2)
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
	result1 := Ok[int](42)
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
	result2 := Err[int](err2)
	option2 := result2.Option()

	if option2.IsSome() {
		t.Errorf("Test 2: Expected None, got Some")
	}
}

func TestResultExpect(t *testing.T) {
	// Test Ok case
	okResult := Ok(42)
	value := okResult.Expect("should not panic")
	if value != 42 {
		t.Errorf("Ok.Expect() = %d, want %d", value, 42)
	}

	// Test Err case - should panic
	errResult := Err[int](errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Error("Err.Expect() should panic, but it didn't")
		} else {
			expectedMsg := "Expect() failed: test message"
			if !strings.Contains(fmt.Sprintf("%v", r), expectedMsg) {
				t.Errorf("Expected panic message to contain '%s', got '%v'", expectedMsg, r)
			}
		}
	}()
	errResult.Expect("test message")
}

func TestResultMapErr(t *testing.T) {
	// Test Ok case - should not change
	okResult := Ok(42)
	mappedOk := okResult.MapErr(func(err error) error {
		return errors.New("new error")
	})

	if !mappedOk.IsOk() || mappedOk.Unwrap() != 42 {
		t.Error("MapErr on Ok result should not change the result")
	}

	// Test Err case - should transform error
	errResult := Err[int](errors.New("original error"))
	mappedErr := errResult.MapErr(func(err error) error {
		return errors.New("mapped error")
	})

	if !mappedErr.IsErr() {
		t.Error("MapErr on Err result should remain Err")
	}

	if mappedErr.Err().Error() != "mapped error" {
		t.Errorf("MapErr should transform error. Expected 'mapped error', got '%s'", mappedErr.Err().Error())
	}
}

func TestResultString(t *testing.T) {
	// Test Ok case
	okResult := Ok(42)
	okStr := okResult.String()
	expected := "Ok(42)"
	if okStr != expected {
		t.Errorf("Ok.String() = '%s', want '%s'", okStr, expected)
	}

	// Test Err case
	errResult := Err[int](errors.New("test error"))
	errStr := errResult.String()
	expectedErr := "Err(test error)"
	if errStr != expectedErr {
		t.Errorf("Err.String() = '%s', want '%s'", errStr, expectedErr)
	}
}

var errSentinel = errors.New("sentinel error")

type customError struct {
	code int
	msg  string
}

func (e *customError) Error() string { return e.msg }

func TestResultErrIs(t *testing.T) {
	// Test Ok case - should return false
	okResult := Ok(42)
	if okResult.ErrIs(errSentinel) {
		t.Error("ErrIs on Ok result should return false")
	}

	// Test Err case with matching error
	errResult := Err[int](errSentinel)
	if !errResult.ErrIs(errSentinel) {
		t.Error("ErrIs should return true for matching sentinel error")
	}

	// Test Err case with wrapped error
	wrappedErr := fmt.Errorf("wrapped: %w", errSentinel)
	wrappedResult := Err[int](wrappedErr)
	if !wrappedResult.ErrIs(errSentinel) {
		t.Error("ErrIs should return true for wrapped sentinel error")
	}

	// Test Err case with non-matching error
	otherErr := errors.New("other error")
	otherResult := Err[int](otherErr)
	if otherResult.ErrIs(errSentinel) {
		t.Error("ErrIs should return false for non-matching error")
	}
}

func TestResultErrAs(t *testing.T) {
	// Test Ok case - should return false
	okResult := Ok(42)
	var target *customError
	if okResult.ErrAs(&target) {
		t.Error("ErrAs on Ok result should return false")
	}

	// Test Err case with matching type
	customErr := &customError{code: 404, msg: "not found"}
	errResult := Err[int](customErr)
	var matched *customError
	if !errResult.ErrAs(&matched) {
		t.Error("ErrAs should return true for matching error type")
	}
	if matched.code != 404 {
		t.Errorf("ErrAs should set target, expected code 404, got %d", matched.code)
	}

	// Test Err case with wrapped custom error
	wrappedErr := fmt.Errorf("wrapped: %w", customErr)
	wrappedResult := Err[int](wrappedErr)
	var wrappedMatched *customError
	if !wrappedResult.ErrAs(&wrappedMatched) {
		t.Error("ErrAs should return true for wrapped custom error")
	}
	if wrappedMatched.code != 404 {
		t.Errorf("ErrAs should set target for wrapped error, expected code 404, got %d", wrappedMatched.code)
	}

	// Test Err case with non-matching type
	plainErr := errors.New("plain error")
	plainResult := Err[int](plainErr)
	var notMatched *customError
	if plainResult.ErrAs(&notMatched) {
		t.Error("ErrAs should return false for non-matching error type")
	}
}

func TestResultErrSource(t *testing.T) {
	// Test Ok case - should return None
	okResult := Ok(42)
	if okResult.ErrSource().IsSome() {
		t.Error("ErrSource on Ok result should return None")
	}

	// Test Err case with no wrapped error
	plainErr := errors.New("plain error")
	plainResult := Err[int](plainErr)
	if plainResult.ErrSource().IsSome() {
		t.Error("ErrSource should return None for error without wrapped error")
	}

	// Test Err case with wrapped error
	innerErr := errors.New("inner error")
	wrappedErr := fmt.Errorf("outer: %w", innerErr)
	wrappedResult := Err[int](wrappedErr)
	source := wrappedResult.ErrSource()
	if source.IsNone() {
		t.Error("ErrSource should return Some for wrapped error")
	}
	if source.Some().Error() != "inner error" {
		t.Errorf("ErrSource should return inner error, got '%s'", source.Some().Error())
	}
}

var (
	errContext         = errors.New("context")
	errFailedToProcess = errors.New("failed to process")
	errLevel1          = errors.New("level 1")
	errLevel2          = errors.New("level 2")
)

func TestResultWrap(t *testing.T) {
	// Test Ok case - should return unchanged
	okResult := Ok(42)
	wrapped := okResult.Wrap(errContext)
	if !wrapped.IsOk() || wrapped.Unwrap() != 42 {
		t.Error("Wrap on Ok result should return unchanged")
	}

	// Test Err case - should wrap error
	originalErr := errors.New("original error")
	errResult := Err[int](originalErr)
	wrappedResult := errResult.Wrap(errFailedToProcess)
	if !wrappedResult.IsErr() {
		t.Error("Wrap on Err result should remain Err")
	}
	expectedMsg := "failed to process: original error"
	if wrappedResult.Err().Error() != expectedMsg {
		t.Errorf("Wrap should prepend error, expected '%s', got '%s'", expectedMsg, wrappedResult.Err().Error())
	}

	// Test that both errors are preserved in chain (errors.Is works for both)
	if !wrappedResult.ErrIs(originalErr) {
		t.Error("Wrapped error should preserve original error for errors.Is")
	}
	if !wrappedResult.ErrIs(errFailedToProcess) {
		t.Error("Wrapped error should preserve wrapper error for errors.Is")
	}

	// Test chaining multiple Wraps
	chainedResult := Err[int](originalErr).Wrap(errLevel1).Wrap(errLevel2)
	expectedChained := "level 2: level 1: original error"
	if chainedResult.Err().Error() != expectedChained {
		t.Errorf("Chained Wrap expected '%s', got '%s'", expectedChained, chainedResult.Err().Error())
	}

	// All errors in chain are accessible
	if !chainedResult.ErrIs(originalErr) || !chainedResult.ErrIs(errLevel1) || !chainedResult.ErrIs(errLevel2) {
		t.Error("All errors in chain should be accessible via ErrIs")
	}
}
