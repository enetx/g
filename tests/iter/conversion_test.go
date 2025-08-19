package g_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/enetx/g/iter"
)

func TestToSlice(t *testing.T) {
	// Test toSlice operation
	result := ToSlice(FromSlice([]int{1, 2, 3}))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToSlice() = %v, want %v", result, expected)
	}
}

func TestPull2(t *testing.T) {
	// Test pull2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	next, stop := Pull2(FromPairs(pairs))
	defer stop()

	k1, v1, ok1 := next()
	if !ok1 || k1 != 1 || v1 != "a" {
		t.Errorf("Pull2() first = %v, %v, %v, want 1, a, true", k1, v1, ok1)
	}

	k2, v2, ok2 := next()
	if !ok2 || k2 != 2 || v2 != "b" {
		t.Errorf("Pull2() second = %v, %v, %v, want 2, b, true", k2, v2, ok2)
	}

	_, _, ok3 := next()
	if ok3 {
		t.Errorf("Pull2() third ok = %v, want false", ok3)
	}
}

func TestContext(t *testing.T) {
	// Test context operation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := ToSlice(Context(FromSlice([]int{1, 2, 3}), ctx))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Context() = %v, want %v", result, expected)
	}
}

func TestContext2(t *testing.T) {
	// Test context2 operation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := ToPairs(Context2(FromPairs(pairs), ctx))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Context2() = %v, want %v", result, expected)
	}
}

func TestToChan(t *testing.T) {
	// Test toChan operation
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := ToChan(FromSlice([]int{1, 2, 3}), ctx)

	var result []int
	for val := range ch {
		result = append(result, val)
	}

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToChan() = %v, want %v", result, expected)
	}
}

func TestToChan2(t *testing.T) {
	// Test toChan2 operation
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	ch := ToChan2(FromPairs(pairs), ctx)

	var result []Pair[int, string]
	for pair := range ch {
		result = append(result, pair)
	}

	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToChan2() = %v, want %v", result, expected)
	}
}

func TestContextCancellation(t *testing.T) {
	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	count := 0
	Context(FromSlice([]int{1, 2, 3, 4, 5}), ctx)(func(x int) bool {
		count++
		return true
	})

	// Should not process any items due to cancelled context
	if count != 0 {
		t.Errorf("ContextCancellation count = %v, want 0", count)
	}
}

func TestContextAlreadyCanceled(t *testing.T) {
	// Test context that's already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := ToSlice(Context(FromSlice([]int{1, 2, 3}), ctx))
	if len(result) != 0 {
		t.Errorf("ContextAlreadyCanceled = %v, want empty slice", result)
	}
}

func TestContext2Cancellation(t *testing.T) {
	// Test context2 cancellation
	ctx, cancel := context.WithCancel(context.Background())

	count := 0
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	Context2(FromPairs(pairs), ctx)(func(k int, v string) bool {
		count++
		if count == 2 {
			cancel() // Cancel after processing 2 items
		}
		return true
	})

	// Should process 2 items then stop due to cancellation
	if count != 2 {
		t.Errorf("Context2Cancellation count = %v, want 2", count)
	}
}

func TestToChanContextCancellation(t *testing.T) {
	// Test toChan with context cancellation
	ctx, cancel := context.WithCancel(context.Background())

	ch := ToChan(FromSlice([]int{1, 2, 3, 4, 5}), ctx)

	var result []int
	for val := range ch {
		result = append(result, val)
		if len(result) == 2 {
			cancel() // Cancel after receiving 2 values
		}
	}

	// Should receive at least 2 values before cancellation
	if len(result) < 2 {
		t.Errorf("ToChanContextCancellation length = %d, want at least 2", len(result))
	}
}

func TestToChan2ContextCancellation(t *testing.T) {
	// Test toChan2 with context cancellation
	ctx, cancel := context.WithCancel(context.Background())

	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	ch := ToChan2(FromPairs(pairs), ctx)

	var result []Pair[int, string]
	for pair := range ch {
		result = append(result, pair)
		if len(result) == 2 {
			cancel() // Cancel after receiving 2 pairs
		}
	}

	// Should receive at least 2 pairs before cancellation
	if len(result) < 2 {
		t.Errorf("ToChan2ContextCancellation length = %d, want at least 2", len(result))
	}
}

func TestToChanContextCancellationDuringIteration(t *testing.T) {
	// Test toChan context cancellation during iteration
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Create a slow sequence
	slowSeq := func(yield func(int) bool) {
		for i := 1; i <= 10; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(5 * time.Millisecond)
				if !yield(i) {
					return
				}
			}
		}
	}

	ch := ToChan(slowSeq, ctx)

	var result []int
	for val := range ch {
		result = append(result, val)
	}

	// Should receive fewer than 10 values due to timeout
	if len(result) >= 10 {
		t.Errorf("ToChanContextCancellationDuringIteration length = %d, want < 10", len(result))
	}
}

func TestToChan2ContextCancellationDuringIteration(t *testing.T) {
	// Test toChan2 context cancellation during iteration
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Create a slow sequence
	slowSeq := func(yield func(int, string) bool) {
		for i := 1; i <= 10; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(5 * time.Millisecond)
				if !yield(i, fmt.Sprintf("val%d", i)) {
					return
				}
			}
		}
	}

	ch := ToChan2(slowSeq, ctx)

	var result []Pair[int, string]
	for pair := range ch {
		result = append(result, pair)
	}

	// Should receive fewer than 10 pairs due to timeout
	if len(result) >= 10 {
		t.Errorf("ToChan2ContextCancellationDuringIteration length = %d, want < 10", len(result))
	}
}

func TestToMap(t *testing.T) {
	// Test basic toMap operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToMap(FromPairs(pairs))

	expected := map[int]string{1: "a", 2: "b", 3: "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToMap() = %v, want %v", result, expected)
	}

	// Test empty sequence
	result2 := ToMap(FromPairs([]Pair[int, string]{}))
	if len(result2) != 0 {
		t.Errorf("ToMap(empty) = %v, want empty map", result2)
	}

	// Test duplicate keys (later values should overwrite earlier ones)
	pairs3 := []Pair[int, string]{{1, "a"}, {2, "b"}, {1, "c"}}
	result3 := ToMap(FromPairs(pairs3))
	if result3[1] != "c" {
		t.Errorf("ToMap(duplicate keys) result[1] = %v, want 'c'", result3[1])
	}
}

func TestContextEarlyReturn(t *testing.T) {
	// Test Context with early return from yield function
	ctx := context.Background()

	count := 0
	Context(FromSlice([]int{1, 2, 3, 4, 5}), ctx)(func(x int) bool {
		count++
		return x != 3 // Stop when we see 3
	})

	if count != 3 {
		t.Errorf("Context early return count = %v, want 3", count)
	}
}

func TestContext2EarlyReturn(t *testing.T) {
	// Test Context2 with early return from yield function
	ctx := context.Background()

	count := 0
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	Context2(FromPairs(pairs), ctx)(func(k int, v string) bool {
		count++
		return k != 3 // Stop when we see key 3
	})

	if count != 3 {
		t.Errorf("Context2 early return count = %v, want 3", count)
	}
}

func TestToChanBuffering(t *testing.T) {
	// Test ToChan with buffering (should not block)
	ctx := context.Background()

	ch := ToChan(FromSlice([]int{1, 2, 3}), ctx)

	// Read first value immediately
	val1, ok1 := <-ch
	if !ok1 || val1 != 1 {
		t.Errorf("ToChan buffering first value = %v, %v, want 1, true", val1, ok1)
	}

	// Read remaining values
	var remaining []int
	for val := range ch {
		remaining = append(remaining, val)
	}

	expected := []int{2, 3}
	if !reflect.DeepEqual(remaining, expected) {
		t.Errorf("ToChan buffering remaining = %v, want %v", remaining, expected)
	}
}

func TestToChan2Buffering(t *testing.T) {
	// Test ToChan2 with buffering (should not block)
	ctx := context.Background()

	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	ch := ToChan2(FromPairs(pairs), ctx)

	// Read first pair immediately
	pair1, ok1 := <-ch
	if !ok1 || pair1.Key != 1 || pair1.Value != "a" {
		t.Errorf("ToChan2 buffering first pair = %v, %v, want {1, a}, true", pair1, ok1)
	}

	// Read remaining pairs
	var remaining []Pair[int, string]
	for pair := range ch {
		remaining = append(remaining, pair)
	}

	expected := []Pair[int, string]{{2, "b"}, {3, "c"}}
	if !reflect.DeepEqual(remaining, expected) {
		t.Errorf("ToChan2 buffering remaining = %v, want %v", remaining, expected)
	}
}
