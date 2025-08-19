package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g/iter"
)

func TestOnce(t *testing.T) {
	// Test once operation
	result := ToSlice(Once(42))
	expected := []int{42}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Once() = %v, want %v", result, expected)
	}
}

func TestOnceWith(t *testing.T) {
	// Test onceWith operation
	called := 0
	result := ToSlice(OnceWith(func() int {
		called++
		return 42
	}))
	expected := []int{42}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OnceWith() = %v, want %v", result, expected)
	}
	if called != 1 {
		t.Errorf("OnceWith() called %d times, want 1", called)
	}
}

func TestEmpty(t *testing.T) {
	// Test empty operation
	result := ToSlice(Empty[int]())
	if len(result) != 0 {
		t.Errorf("Empty() = %v, want empty slice", result)
	}
}

func TestRepeat(t *testing.T) {
	// Test repeat operation
	result := ToSlice(Take(Repeat(42), 5))
	expected := []int{42, 42, 42, 42, 42}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Repeat() = %v, want %v", result, expected)
	}
}

func TestRepeatWith(t *testing.T) {
	// Test repeatWith operation
	count := 0
	result := ToSlice(Take(RepeatWith(func() int {
		count++
		return count
	}), 3))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RepeatWith() = %v, want %v", result, expected)
	}
}

func TestIota(t *testing.T) {
	// Test basic iota
	result := ToSlice(Iota(1, 5))
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Iota(1, 5) = %v, want %v", result, expected)
	}

	// Test iota with step
	result2 := ToSlice(Iota(1, 10, 2))
	expected2 := []int{1, 3, 5, 7, 9}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Iota(1, 10, 2) = %v, want %v", result2, expected2)
	}

	// Test iota empty range
	result3 := ToSlice(Iota(5, 5))
	if len(result3) != 0 {
		t.Errorf("Iota(5, 5) = %v, want empty slice", result3)
	}
}

func TestIotaInclusive(t *testing.T) {
	// Test basic iota inclusive
	result := ToSlice(IotaInclusive(1, 5))
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("IotaInclusive(1, 5) = %v, want %v", result, expected)
	}

	// Test iota inclusive with step
	result2 := ToSlice(IotaInclusive(1, 9, 2))
	expected2 := []int{1, 3, 5, 7, 9}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("IotaInclusive(1, 9, 2) = %v, want %v", result2, expected2)
	}
}

func TestIotaNegativeStep(t *testing.T) {
	// Test iota with negative step
	result := ToSlice(Iota(5, 1, -1))
	expected := []int{5, 4, 3, 2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Iota(5, 1, -1) = %v, want %v", result, expected)
	}
}

func TestIotaInclusiveNegativeStep(t *testing.T) {
	// Test iota inclusive with negative step
	result := ToSlice(IotaInclusive(5, 1, -1))
	expected := []int{5, 4, 3, 2, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("IotaInclusive(5, 1, -1) = %v, want %v", result, expected)
	}
}

func TestIotaEarlyTermination(t *testing.T) {
	// Test iota with early termination
	count := 0
	Iota(1, 10)(func(x int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Errorf("IotaEarlyTermination count = %v, want 3", count)
	}
}

func TestIotaNegativeStepEarlyTermination(t *testing.T) {
	// Test iota negative step with early termination
	count := 0
	Iota(10, 1, -1)(func(x int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Errorf("IotaNegativeStepEarlyTermination count = %v, want 3", count)
	}
}

func TestIotaInclusiveEarlyTermination(t *testing.T) {
	// Test iota inclusive with early termination
	count := 0
	IotaInclusive(1, 10)(func(x int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Errorf("IotaInclusiveEarlyTermination count = %v, want 3", count)
	}
}

func TestIotaInclusiveNegativeStepEarlyTermination(t *testing.T) {
	// Test iota inclusive negative step with early termination
	count := 0
	IotaInclusive(10, 1, -1)(func(x int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Errorf("IotaInclusiveNegativeStepEarlyTermination count = %v, want 3", count)
	}
}

func TestIotaZeroStep(t *testing.T) {
	// Test iota with zero step
	result := ToSlice(Iota(1, 5, 0))
	if len(result) != 0 {
		t.Errorf("Iota(zero step) = %v, want empty slice", result)
	}
}

func TestIotaInclusiveZeroStep(t *testing.T) {
	// Test iota inclusive with zero step
	result := ToSlice(IotaInclusive(1, 5, 0))
	if len(result) != 0 {
		t.Errorf("IotaInclusive(zero step) = %v, want empty slice", result)
	}
}

func TestCounter(t *testing.T) {
	// Test counter operation
	result := Counter(FromSlice([]int{1, 2, 1, 3, 2, 1}))
	if result[1] != 3 || result[2] != 2 || result[3] != 1 {
		t.Errorf("Counter() = %v, want map with correct counts", result)
	}
}
