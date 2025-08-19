package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g/iter"
)

func TestFromSlice(t *testing.T) {
	// Test basic fromSlice operation
	result := ToSlice(FromSlice([]int{1, 2, 3, 4, 5}))
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromSlice() = %v, want %v", result, expected)
	}

	// Test fromSlice with empty slice
	result2 := ToSlice(FromSlice([]int{}))
	if len(result2) != 0 {
		t.Errorf("FromSlice(empty) = %v, want empty slice", result2)
	}
}

func TestFromChan(t *testing.T) {
	// Test fromChan operation
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	result := ToSlice(FromChan(ch))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromChan() = %v, want %v", result, expected)
	}
}

func TestFromMap(t *testing.T) {
	// Test fromMap operation
	m := map[int]string{1: "a", 2: "b"}
	result := ToPairs(FromMap(m))

	// Since map iteration order is not guaranteed, check length and contents
	if len(result) != 2 {
		t.Errorf("FromMap() length = %d, want 2", len(result))
	}

	found := make(map[int]string)
	for _, pair := range result {
		found[pair.Key] = pair.Value
	}

	if found[1] != "a" || found[2] != "b" {
		t.Errorf("FromMap() contents = %v, want map with correct key-value pairs", found)
	}
}

func TestFromChanEmpty(t *testing.T) {
	// Test fromChan with empty channel
	ch := make(chan int)
	close(ch)

	result := ToSlice(FromChan(ch))
	if len(result) != 0 {
		t.Errorf("FromChan(empty) = %v, want empty slice", result)
	}
}

func TestFromChanEarlyTermination(t *testing.T) {
	// Test fromChan with early termination
	ch := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch)

	count := 0
	FromChan(ch)(func(x int) bool {
		count++
		return count < 3
	})

	if count != 3 {
		t.Errorf("FromChanEarlyTermination count = %v, want 3", count)
	}
}

func TestFromMapEarlyTermination(t *testing.T) {
	// Test fromMap with early termination
	m := map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}

	count := 0
	FromMap(m)(func(k int, v string) bool {
		count++
		return count < 3
	})

	if count != 3 {
		t.Errorf("FromMapEarlyTermination count = %v, want 3", count)
	}
}

func TestFromSliceReverse(t *testing.T) {
	// Test reverse operation
	result := ToSlice(FromSliceReverse([]int{1, 2, 3, 4, 5}))
	expected := []int{5, 4, 3, 2, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromSliceReverse() = %v, want %v", result, expected)
	}
}

func TestFromSliceReverseEmpty(t *testing.T) {
	// Test reverse with empty sequence
	result := ToSlice(FromSliceReverse([]int{}))
	if len(result) != 0 {
		t.Errorf("FromSliceReverse(empty) = %v, want empty slice", result)
	}
}

func TestFromSliceReverseEarlyTermination(t *testing.T) {
	// Test reverse with early termination
	count := 0
	FromSliceReverse([]int{1, 2, 3, 4, 5})(func(x int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Errorf("FromSliceReverseEarlyTermination count = %v, want 3", count)
	}
}
