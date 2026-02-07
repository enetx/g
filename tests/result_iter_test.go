package g_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/pool"
)

func TestSeqResultAll(t *testing.T) {
	t.Run("all pass, no error", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		res := seq.All(func(v int) bool { return v > 0 })
		if res.IsErr() {
			t.Errorf("unexpected error: %v", res.Err())
		} else if !res.Ok() {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("one fails predicate, no error", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(-2)) // fails the predicate
			yield(Ok(3))
		})

		res := seq.All(func(v int) bool { return v > 0 })
		if res.IsErr() {
			t.Errorf("unexpected error: %v", res.Err())
		} else if res.Ok() {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("encounter error early", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("some error")))
			yield(Ok(3)) // never reached
		})

		res := seq.All(func(v int) bool { return v > 0 })
		if !res.IsErr() {
			t.Errorf("expected an error, got Ok(%v)", res.Ok())
		} else {
			wantErrMsg := "some error"
			if res.Err().Error() != wantErrMsg {
				t.Errorf("expected error message %q, got %q", wantErrMsg, res.Err().Error())
			}
		}
	})
}

func TestSeqResultAny(t *testing.T) {
	t.Run("no match, no error", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(3))
			yield(Ok(5))
		})

		res := seq.Any(func(x int) bool { return x%2 == 0 }) // Searching for an even number

		if res.IsErr() {
			t.Fatalf("expected Ok(false), got Err: %v", res.Err())
		}
		if res.Ok() {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("match on first item", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(2)) // even number
			yield(Ok(3))
			yield(Ok(4))
		})

		res := seq.Any(func(x int) bool { return x%2 == 0 })

		if res.IsErr() {
			t.Fatalf("expected Ok(true), got Err: %v", res.Err())
		}
		if !res.Ok() {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("match on second item", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(5)) // not even
			yield(Ok(6)) // even
			yield(Ok(7))
		})

		res := seq.Any(func(x int) bool { return x%2 == 0 })

		if res.IsErr() {
			t.Fatalf("expected Ok(true), got Err: %v", res.Err())
		}
		if !res.Ok() {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("encounter error", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("some error")))
			yield(Ok(3)) // won't be yielded because of early stop
		})

		res := seq.Any(func(x int) bool { return x == 2 })

		if !res.IsErr() {
			t.Errorf("expected an error, got Ok(%v)", res.Ok())
		} else {
			wantErrMsg := "some error"
			if res.Err().Error() != wantErrMsg {
				t.Errorf("expected error message %q, got %q", wantErrMsg, res.Err().Error())
			}
		}
	})

	t.Run("no elements at all", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {})

		res := seq.Any(func(x int) bool { return x == 100 })

		if res.IsErr() {
			t.Fatalf("expected Ok(false), got Err: %v", res.Err())
		}
		if res.Ok() {
			t.Errorf("expected false, got true")
		}
	})
}

func TestSeqResultCollect(t *testing.T) {
	t.Run("all Ok", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		res := seq
		if res.FirstErr().IsSome() {
			t.Fatalf("expected Ok([10, 20, 30]) but got Err: %v", res.FirstErr().Some())
		}

		collected := res.Ok().Collect()
		if len(collected) != 3 {
			t.Errorf("expected 3 elements, got %d", len(collected))
		}
		if collected[0] != 10 || collected[1] != 20 || collected[2] != 30 {
			t.Errorf("expected [10, 20, 30], got %v", collected)
		}
	})

	t.Run("error on first element", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](errors.New("boom")))
			yield(Ok(20)) // not reached
		})

		res := seq
		if res.FirstErr().IsNone() {
			t.Fatalf("expected Err, got %v", res.Ok().Collect())
		}

		// Optional: check error message
		wantErr := "boom"
		if res.FirstErr().Some().Error() != wantErr {
			t.Errorf("expected error %q, got %q", wantErr, res.FirstErr().Some().Error())
		}
	})

	t.Run("error in the middle", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("middle error")))
			yield(Ok(3)) // won't be collected
		})

		firstErr := seq.FirstErr()
		if firstErr.IsNone() {
			collected := seq.Ok().Collect()
			t.Fatalf("expected Err but got Ok(%v)", collected)
		}

		wantErr := "middle error"
		if firstErr.Some().Error() != wantErr {
			t.Errorf("expected error %q, got %q", wantErr, firstErr.Some().Error())
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {})

		firstErr := seq.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}

		collected := seq.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected 0 elements, got %d", len(collected))
		}
	})
}

func TestSeqResultCount(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		got := seq.Count()
		want := Int(0)
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	})

	t.Run("all Ok", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		got := seq.Count()
		want := Int(3)
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	})

	t.Run("some Ok, some Err", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("boom")))
			yield(Ok(2))
			yield(Ok(3))
		})

		got := seq.Count()
		// The Count method ignores whether values are Ok or Err; it counts them all
		want := Int(4)
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	})
}

func TestSeqResultMap(t *testing.T) {
	t.Run("all Ok", func(t *testing.T) {
		// Create a sequence of Ok values
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		// Transform each Ok value: multiply by 10
		mapped := seq.Map(func(v int) int {
			return v * 10
		})

		// Collect transformed results
		firstErr := mapped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected all Ok, got Err: %v", firstErr.Some())
		}

		values := mapped.Ok().Collect()
		if len(values) != 3 {
			t.Errorf("expected 3 items, got %d", len(values))
		}
		expected := []int{10, 20, 30}
		for i, val := range values {
			if val != expected[i] {
				t.Errorf("expected %v at index %d, got %v", expected[i], i, val)
			}
		}
	})

	t.Run("err in the middle", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Err[int](errors.New("boom")))
			yield(Ok(30)) // won't be reached, because iteration stops on first Err
		})

		// Transform each Ok value by adding 1
		mapped := seq.Map(func(v int) int {
			return v + 1
		})

		// Check if there's an error
		firstErr := mapped.FirstErr()
		if firstErr.IsNone() {
			collected := mapped.Ok().Collect()
			t.Fatalf("expected an error, got Ok(%v)", collected)
		}

		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error message \"boom\", got %q", errMsg)
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		// Any transform function won't matter because the sequence is empty
		mapped := seq.Map(func(v int) int {
			return v * 2
		})

		firstErr := mapped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}
		collected := mapped.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected 0 items, got %d", len(collected))
		}
	})

	t.Run("err is first element", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](errors.New("immediate error")))
			yield(Ok(99)) // never reached
		})

		mapped := seq.Map(func(v int) int {
			return v + 1
		})

		firstErr := mapped.FirstErr()
		if firstErr.IsNone() {
			collected := mapped.Ok().Collect()
			t.Fatalf("expected an error, got Ok(%v)", collected)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "immediate error" {
			t.Errorf("expected \"immediate error\", got %q", errMsg)
		}
	})
}

func TestSeqResultFilter(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		filtered := seq.Filter(func(v int) bool {
			return v > 0
		})

		firstErr := filtered.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}

		collected := filtered.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected an empty slice, got %d elements", len(collected))
		}
	})

	t.Run("all Ok, some fail predicate", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))  // pass predicate (v > 0)
			yield(Ok(-1)) // fail predicate
			yield(Ok(2))  // pass
			yield(Ok(0))  // fail
			yield(Ok(3))  // pass
		})

		filtered := seq.Filter(func(v int) bool {
			return v > 0
		})

		firstErr := filtered.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := filtered.Ok().Collect()
		if len(collected) != 3 {
			t.Errorf("expected 3 elements, got %d", len(collected))
		}
		want := []int{1, 2, 3}
		for i, gotVal := range collected {
			if gotVal != want[i] {
				t.Errorf("expected %d, got %d at index %d", want[i], gotVal, i)
			}
		}
	})

	t.Run("encounter error", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Err[int](errors.New("boom")))
			yield(Ok(20)) // never reached, iteration stops on error
		})

		filtered := seq.Filter(func(v int) bool {
			return v > 0
		})

		firstErr := filtered.FirstErr()
		if firstErr.IsNone() {
			collected := filtered.Ok().Collect()
			t.Fatalf("expected an error, got Ok(%v)", collected)
		}

		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error \"boom\", got %q", errMsg)
		}
	})

	t.Run("all fail predicate", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		filtered := seq.Filter(func(v int) bool {
			return v < 0 // none are < 0
		})

		firstErr := filtered.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}

		collected := filtered.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected 0 elements, got %d", len(collected))
		}
	})
}

func TestSeqResultExclude(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		excluded := seq.Exclude(func(v int) bool {
			return v < 0 // doesn't matter, there are no items
		})

		firstErr := excluded.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}
		collected := excluded.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected empty slice, got %d items", len(collected))
		}
	})

	t.Run("exclude some Ok", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1)) // keep
			yield(Ok(2)) // exclude
			yield(Ok(3)) // keep
			yield(Ok(4)) // exclude
		})

		// Exclude even numbers
		excluded := seq.Exclude(func(v int) bool {
			return v%2 == 0
		})

		firstErr := excluded.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := excluded.Ok().Collect()
		if len(collected) != 2 {
			t.Errorf("expected 2 items, got %d", len(collected))
		}
		want := []int{1, 3}
		for i, v := range collected {
			if v != want[i] {
				t.Errorf("at index %d: want %d, got %d", i, want[i], v)
			}
		}
	})

	t.Run("encounter error", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Err[int](errors.New("boom")))
			yield(Ok(20)) // not reached, iteration stops on error
		})

		excluded := seq.Exclude(func(v int) bool {
			return v < 0
		})

		firstErr := excluded.FirstErr()
		if firstErr.IsNone() {
			collected := excluded.Ok().Collect()
			t.Fatalf("expected Err, got Ok(%v)", collected)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected \"boom\", got %q", errMsg)
		}
	})

	t.Run("all satisfy fn (exclude all)", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(2))
			yield(Ok(4))
			yield(Ok(6))
		})

		// Exclude all evens
		excluded := seq.Exclude(func(v int) bool {
			return v%2 == 0
		})

		firstErr := excluded.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}
		collected := excluded.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected no items, got %d", len(collected))
		}
	})

	t.Run("none satisfy fn (exclude none)", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		// Exclude if < 0, but all are >= 1
		excluded := seq.Exclude(func(v int) bool {
			return v < 0
		})

		firstErr := excluded.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := excluded.Ok().Collect()
		if len(collected) != 3 {
			t.Errorf("expected all 3 items, got %d", len(collected))
		}
		want := []int{1, 2, 3}
		for i, v := range collected {
			if v != want[i] {
				t.Errorf("at index %d: want %d, got %d", i, want[i], v)
			}
		}
	})
}

func TestSeqResultDedup(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		deduped := seq.Dedup()
		firstErr := deduped.FirstErr()

		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}

		collected := deduped.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected empty slice, got %d items", len(collected))
		}
	})

	t.Run("all duplicates", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(10))
			yield(Ok(10))
			yield(Ok(10))
		})

		deduped := seq.Dedup()
		firstErr := deduped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := deduped.Ok().Collect()
		if len(collected) != 1 {
			t.Errorf("expected only 1 item after dedup, got %d", len(collected))
		}
		if len(collected) > 0 && collected[0] != 10 {
			t.Errorf("expected 10, got %d", collected[0])
		}
	})

	t.Run("no duplicates (already unique)", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		deduped := seq.Dedup()
		firstErr := deduped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := deduped.Ok().Collect()
		if len(collected) != 3 {
			t.Errorf("expected 3 items, got %d", len(collected))
		}
		expected := []int{1, 2, 3}
		for i, val := range collected {
			if val != expected[i] {
				t.Errorf("index %d: expected %d, got %d", i, expected[i], val)
			}
		}
	})

	t.Run("some duplicates", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(5))
			yield(Ok(5)) // duplicate
			yield(Ok(10))
			yield(Ok(10)) // duplicate
			yield(Ok(10)) // duplicate
			yield(Ok(15))
			yield(Ok(15)) // duplicate
			yield(Ok(20))
		})

		deduped := seq.Dedup()
		firstErr := deduped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := deduped.Ok().Collect()
		expected := []int{5, 10, 15, 20}
		if len(collected) != len(expected) {
			t.Errorf("expected %d items, got %d", len(expected), len(collected))
		}
		for i, val := range collected {
			if val != expected[i] {
				t.Errorf("index %d: expected %d, got %d", i, expected[i], val)
			}
		}
	})

	t.Run("error in the middle", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(1)) // duplicate
			yield(Ok(2))
			yield(Err[int](errors.New("boom")))
			yield(Ok(2)) // never reached
			yield(Ok(3)) // never reached
		})

		deduped := seq.Dedup()
		firstErr := deduped.FirstErr()
		if firstErr.IsNone() {
			collected := deduped.Ok().Collect()
			t.Fatalf("expected Err, got Ok(%v)", collected)
		}

		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error message \"boom\", got %q", errMsg)
		}
	})
}

func TestSeqResultUnique(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		unique := seq.Unique()
		firstErr := unique.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}
		collected := unique.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected 0 elements, got %d", len(collected))
		}
	})

	t.Run("all duplicates", func(t *testing.T) {
		// For instance, all values are 10
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(10))
			yield(Ok(10))
			yield(Ok(10))
		})

		unique := seq.Unique()
		firstErr := unique.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := unique.Ok().Collect()
		// Only the first 10 should appear
		if len(collected) != 1 {
			t.Errorf("expected 1 element, got %d", len(collected))
		} else if collected[0] != 10 {
			t.Errorf("expected 10, got %d", collected[0])
		}
	})

	t.Run("some duplicates", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(1)) // duplicate
			yield(Ok(3))
			yield(Ok(2)) // duplicate
			yield(Ok(4))
		})

		unique := seq.Unique()
		firstErr := unique.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := unique.Ok().Collect()
		// Expect [1, 2, 3, 4], skipping repeated 1 and 2
		want := []int{1, 2, 3, 4}
		if len(collected) != len(want) {
			t.Errorf("expected %d elements, got %d", len(want), len(collected))
		}
		for i, val := range collected {
			if val != want[i] {
				t.Errorf("index %d: expected %d, got %d", i, want[i], val)
			}
		}
	})

	t.Run("all unique", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		unique := seq.Unique()
		firstErr := unique.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := unique.Ok().Collect()
		if len(collected) != 3 {
			t.Errorf("expected 3 elements, got %d", len(collected))
		}
		want := []int{10, 20, 30}
		for i, val := range collected {
			if val != want[i] {
				t.Errorf("index %d: want %d, got %d", i, want[i], val)
			}
		}
	})

	t.Run("error in the middle", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Err[int](errors.New("boom")))
			yield(Ok(2)) // never reached
			yield(Ok(3)) // never reached
		})

		unique := seq.Unique()
		firstErr := unique.FirstErr()
		if firstErr.IsNone() {
			collected := unique.Ok().Collect()
			t.Fatalf("expected Err, got Ok(%v)", collected)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error message \"boom\", got %q", errMsg)
		}
	})
}

func TestSeqResultForEach(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		var callCount int

		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		// ForEach should not call fn at all in this case
		seq.ForEach(func(v Result[int]) {
			callCount++
		})

		if callCount != 0 {
			t.Errorf("expected 0 calls, got %d", callCount)
		}
	})

	t.Run("multiple items including Err", func(t *testing.T) {
		// We'll yield an Ok(1), then an Err, then Ok(2).
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("boom")))
			yield(Ok(2))
		})

		var visited []Result[int]
		seq.ForEach(func(v Result[int]) {
			visited = append(visited, v)
		})

		// ForEach does not stop on Err; it continues for all items.
		if len(visited) != 3 {
			t.Errorf("expected 3 items, got %d", len(visited))
		} else {
			if !visited[0].IsOk() || visited[0].Ok() != 1 {
				t.Errorf("expected Ok(1), got %v", visited[0])
			}
			if !visited[1].IsErr() {
				t.Errorf("expected Err, got %v", visited[1])
			} else {
				errMsg := visited[1].Err().Error()
				if errMsg != "boom" {
					t.Errorf("expected error 'boom', got %q", errMsg)
				}
			}
			if !visited[2].IsOk() || visited[2].Ok() != 2 {
				t.Errorf("expected Ok(2), got %v", visited[2])
			}
		}
	})

	t.Run("all Ok items", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(100))
			yield(Ok(200))
			yield(Ok(300))
		})

		var collected []int
		seq.ForEach(func(v Result[int]) {
			if v.IsOk() {
				collected = append(collected, v.Ok())
			}
		})

		if len(collected) != 3 {
			t.Errorf("expected 3 items, got %d", len(collected))
		}
		want := []int{100, 200, 300}
		for i, got := range collected {
			if got != want[i] {
				t.Errorf("index %d: expected %d, got %d", i, want[i], got)
			}
		}
	})
}

func TestSeqResultRange(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		callCount := 0

		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		seq.Range(func(v Result[int]) bool {
			callCount++
			return true
		})

		if callCount != 0 {
			t.Errorf("expected 0 calls, got %d", callCount)
		}
	})

	t.Run("all Ok, never returns false", func(t *testing.T) {
		callCount := 0
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		seq.Range(func(v Result[int]) bool {
			callCount++
			if v.IsErr() {
				t.Errorf("unexpected error: %v", v.Err())
			} else {
				// Check the Ok value if desired
				_ = v.Ok()
			}
			return true // never stops
		})

		if callCount != 3 {
			t.Errorf("expected 3 calls, got %d", callCount)
		}
	})

	t.Run("stop in the middle", func(t *testing.T) {
		var visited []Result[int]

		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			if !yield(Ok(1)) {
				return
			}
			if !yield(Ok(2)) {
				return
			}
			if !yield(Ok(3)) {
				return
			}
			if !yield(Ok(4)) {
				return
			}
			if !yield(Ok(5)) {
				return
			}
		})

		seq.Range(func(v Result[int]) bool {
			visited = append(visited, v)
			// stop when we see the third item or above
			return v.Ok() < 3
		})

		// Now we expect [Ok(1), Ok(2), Ok(3)]
		if len(visited) != 3 {
			t.Errorf("expected 3 items visited, got %d", len(visited))
		} else {
			if visited[0].Ok() != 1 || visited[1].Ok() != 2 {
				t.Errorf("unexpected values: %v", visited)
			}
		}
	})

	t.Run("encounter error, decide to stop", func(t *testing.T) {
		var visited []Result[int]

		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			if !yield(Ok(10)) {
				return
			}
			if !yield(Err[int](errors.New("boom"))) {
				return
			}
			if !yield(Ok(20)) {
				return
			}
		})

		seq.Range(func(v Result[int]) bool {
			visited = append(visited, v)
			if v.IsErr() {
				// stop iteration if we see an error
				return false
			}
			return true
		})

		// Now we expect [Ok(10), Err("boom")]
		if len(visited) != 2 {
			t.Errorf("expected 2 items visited, got %d", len(visited))
		} else {
			if !visited[0].IsOk() || visited[0].Ok() != 10 {
				t.Errorf("expected Ok(10), got %v", visited[0])
			}
			if !visited[1].IsErr() || visited[1].Err().Error() != "boom" {
				t.Errorf("expected Err(\"boom\"), got %v", visited[1])
			}
		}
	})
}

func TestSeqResultSkip(t *testing.T) {
	t.Run("skip = 0", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		skipped := seq.Skip(0) // should skip none
		firstErr := skipped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := skipped.Ok().Collect()
		if len(collected) != 3 {
			t.Errorf("expected 3 items, got %d", len(collected))
		} else {
			want := []int{1, 2, 3}
			for i, val := range collected {
				if val != want[i] {
					t.Errorf("index %d: want %d, got %d", i, want[i], val)
				}
			}
		}
	})

	t.Run("skip first 2 ok", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
			yield(Ok(40))
			yield(Ok(50))
		})

		skipped := seq.Skip(2) // skip the first 2 Ok
		allResults := skipped.Collect()

		// Check for errors
		var firstErr error
		for _, r := range allResults {
			if r.IsErr() {
				firstErr = r.Err()
				break
			}
		}
		if firstErr != nil {
			t.Fatalf("expected Ok, got Err: %v", firstErr)
		}

		// Collect Ok values
		var collected []int
		for _, r := range allResults {
			if r.IsOk() {
				collected = append(collected, r.Ok())
			}
		}

		// After skipping 10, 20 => we should have [30, 40, 50]
		if len(collected) != 3 {
			t.Errorf("expected 3 items, got %d", len(collected))
		} else {
			want := []int{30, 40, 50}
			for i, val := range collected {
				if val != want[i] {
					t.Errorf("index %d: want %d, got %d", i, want[i], val)
				}
			}
		}
	})

	t.Run("skip more than total ok", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		skipped := seq.Skip(5) // skip 5, but we only have 3 Ok items
		allResults := skipped.Collect()

		// Check for errors
		var firstErr error
		for _, r := range allResults {
			if r.IsErr() {
				firstErr = r.Err()
				break
			}
		}
		if firstErr != nil {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr)
		}

		// Collect Ok values
		var collected []int
		for _, r := range allResults {
			if r.IsOk() {
				collected = append(collected, r.Ok())
			}
		}
		// We should have no remaining items
		if len(collected) != 0 {
			t.Errorf("expected 0 items, got %d", len(collected))
		}
	})

	t.Run("encounter error before skip is done", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Err[int](errors.New("boom")))
			yield(Ok(20)) // never reached, because iteration stops on error
			yield(Ok(30)) // never reached
		})

		skipped := seq.Skip(2) // wants to skip 2 Ok items, but there's an error after the first Ok
		firstErr := skipped.FirstErr()
		if firstErr.IsNone() {
			collected := skipped.Ok().Collect()
			t.Fatalf("expected an Err, got Ok(%v)", collected)
		}

		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error message \"boom\", got %q", errMsg)
		}
	})
}

func TestSeqResultStepBy(t *testing.T) {
	t.Run("step = 1 (yield all Ok)", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		stepped := seq.StepBy(1)
		firstErr := stepped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := stepped.Ok().Collect()
		if len(collected) != 3 {
			t.Errorf("expected 3 items, got %d", len(collected))
		}
		want := []int{10, 20, 30}
		for i, v := range collected {
			if v != want[i] {
				t.Errorf("index %d: want %d, got %d", i, want[i], v)
			}
		}
	})

	t.Run("step = 2 (every second Ok)", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
			yield(Ok(4))
			yield(Ok(5))
		})

		stepped := seq.StepBy(2)
		firstErr := stepped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := stepped.Ok().Collect()
		// We yield i=1 => keep, i=2 => skip, i=3 => keep, i=4 => skip, i=5 => keep
		want := []int{1, 3, 5}
		if len(collected) != len(want) {
			t.Errorf("expected %d items, got %d", len(want), len(collected))
		} else {
			for i, v := range collected {
				if v != want[i] {
					t.Errorf("index %d: want %d, got %d", i, want[i], v)
				}
			}
		}
	})

	t.Run("step > total Ok items", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
		})

		stepped := seq.StepBy(5) // 5 > total Ok items = 2
		firstErr := stepped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := stepped.Ok().Collect()
		// We yield i=1 => keep, i=2 => skip, no more items
		// => [10]
		if len(collected) != 1 {
			t.Errorf("expected 1 item, got %d", len(collected))
		} else {
			if collected[0] != 10 {
				t.Errorf("expected [10], got %v", collected)
			}
		}
	})

	t.Run("encounter error", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("boom")))
			yield(Ok(2))
			yield(Ok(3))
		})

		stepped := seq.StepBy(2)
		firstErr := stepped.FirstErr()
		if firstErr.IsNone() {
			collected := stepped.Ok().Collect()
			t.Fatalf("expected an Err, got Ok(%v)", collected)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error message \"boom\", got %q", errMsg)
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		stepped := seq.StepBy(3)
		firstErr := stepped.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}
		collected := stepped.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected empty, got %d items", len(collected))
		}
	})
}

func TestSeqResultTake(t *testing.T) {
	t.Run("take = 0", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
		})

		taken := seq.Take(0)
		firstErr := taken.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}
		collected := taken.Ok().Collect()
		// No elements should be yielded
		if len(collected) != 0 {
			t.Errorf("expected 0, got %d", len(collected))
		}
	})

	t.Run("take < total Ok", func(t *testing.T) {
		// Yields 3 Ok values but we only want 2
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		taken := seq.Take(2)
		allResults := taken.Collect()

		// Check for errors
		var firstErr error
		for _, r := range allResults {
			if r.IsErr() {
				firstErr = r.Err()
				break
			}
		}
		if firstErr != nil {
			t.Fatalf("expected Ok, got Err: %v", firstErr)
		}

		// Collect Ok values
		var collected []int
		for _, r := range allResults {
			if r.IsOk() {
				collected = append(collected, r.Ok())
			}
		}
		// We should have only the first 2 Ok values
		if len(collected) != 2 {
			t.Errorf("expected 2 items, got %d", len(collected))
		} else {
			want := []int{10, 20}
			for i, val := range collected {
				if val != want[i] {
					t.Errorf("index %d: want %d, got %d", i, want[i], val)
				}
			}
		}
	})

	t.Run("take = total Ok", func(t *testing.T) {
		// Yields exactly 3 Ok values, and we want 3
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		taken := seq.Take(3)
		allResults := taken.Collect()

		// Check for errors
		var firstErr error
		for _, r := range allResults {
			if r.IsErr() {
				firstErr = r.Err()
				break
			}
		}
		if firstErr != nil {
			t.Fatalf("expected Ok, got Err: %v", firstErr)
		}

		// Collect Ok values
		var collected []int
		for _, r := range allResults {
			if r.IsOk() {
				collected = append(collected, r.Ok())
			}
		}
		if len(collected) != 3 {
			t.Errorf("expected 3, got %d", len(collected))
		} else {
			want := []int{1, 2, 3}
			for i, val := range collected {
				if val != want[i] {
					t.Errorf("index %d: want %d, got %d", i, want[i], val)
				}
			}
		}
	})

	t.Run("take > total Ok", func(t *testing.T) {
		// Yields 2 Ok, but we want 5
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(7))
			yield(Ok(14))
		})

		taken := seq.Take(5)
		firstErr := taken.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := taken.Ok().Collect()
		// We can only get the 2 Ok items
		if len(collected) != 2 {
			t.Errorf("expected 2, got %d", len(collected))
		} else {
			want := []int{7, 14}
			for i, val := range collected {
				if val != want[i] {
					t.Errorf("index %d: want %d, got %d", i, want[i], val)
				}
			}
		}
	})

	t.Run("error in the middle", func(t *testing.T) {
		// Yields Ok(1), then Err, then Ok(3) which won't be reached
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("boom")))
			yield(Ok(3)) // not reached
		})

		taken := seq.Take(2)
		firstErr := taken.FirstErr()
		if firstErr.IsNone() {
			collected := taken.Ok().Collect()
			t.Fatalf("expected Err, got Ok(%v)", collected)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error message \"boom\", got %q", errMsg)
		}
	})
}

func TestSeqResultChain(t *testing.T) {
	t.Run("no additional sequences", func(t *testing.T) {
		// Single sequence with a few Ok items
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
		})

		// Chain with no additional sequences
		chained := seq.Chain()
		firstErr := chained.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := chained.Ok().Collect()
		if len(collected) != 2 {
			t.Errorf("expected 2 items, got %d", len(collected))
		} else {
			want := []int{1, 2}
			for i, val := range collected {
				if val != want[i] {
					t.Errorf("index %d: want %d, got %d", i, want[i], val)
				}
			}
		}
	})

	t.Run("chain multiple, all Ok", func(t *testing.T) {
		// First sequence
		seq1 := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
		})
		// Second sequence
		seq2 := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(30))
			yield(Ok(40))
		})
		// Third sequence
		seq3 := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(50))
		})

		chained := seq1.Chain(seq2, seq3)
		firstErr := chained.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}
		collected := chained.Ok().Collect()
		// Expect them in the order: seq1 -> seq2 -> seq3
		want := []int{10, 20, 30, 40, 50}
		if len(collected) != len(want) {
			t.Errorf("expected %d items, got %d", len(want), len(collected))
		} else {
			for i, val := range collected {
				if val != want[i] {
					t.Errorf("index %d: want %d, got %d", i, want[i], val)
				}
			}
		}
	})

	t.Run("error in the first sequence", func(t *testing.T) {
		// First sequence has an error quickly
		seq1 := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("error in seq1")))
			yield(Ok(2))
		})
		// Second sequence, won't be reached if error stops iteration
		seq2 := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(3))
		})

		chained := seq1.Chain(seq2)
		firstErr := chained.FirstErr()

		if firstErr.IsNone() {
			collected := chained.Ok().Collect()
			t.Fatalf("expected an error, got Ok(%v)", collected)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "error in seq1" {
			t.Errorf("expected error message %q, got %q", "error in seq1", errMsg)
		}
	})

	t.Run("error in the second sequence", func(t *testing.T) {
		seq1 := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(100))
			yield(Ok(200))
		})
		seq2 := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(300))
			yield(Err[int](errors.New("error in seq2")))
			// Would have more items, but we won't get here
			yield(Ok(999))
		})
		seq3 := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(400)) // won't be reached if iteration stops on error
		})

		chained := seq1.Chain(seq2, seq3)
		firstErr := chained.FirstErr()
		if firstErr.IsNone() {
			collected := chained.Ok().Collect()
			t.Fatalf("expected an error, got Ok(%v)", collected)
		}

		// The items from seq1 and the Ok(300) from seq2 should be seen,
		// then we see the error from seq2 and stop. seq3 is never visited.
		// If you want to check the partial results before the error,
		// you'd have to test with a short-circuit approach, or re-implement differently.
		errMsg := firstErr.Some().Error()
		if errMsg != "error in seq2" {
			t.Errorf("expected error message %q, got %q", "error in seq2", errMsg)
		}
	})

	t.Run("chaining empty sequences", func(t *testing.T) {
		emptySeq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})
		seqWithItems := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(999))
		})

		chained := emptySeq.Chain(emptySeq, seqWithItems, emptySeq)
		firstErr := chained.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := chained.Ok().Collect()
		// We only expect the one Ok(999) from seqWithItems
		if len(collected) != 1 || collected[0] != 999 {
			t.Errorf("expected [999], got %v", collected)
		}
	})
}

func TestSeqResultIntersperse(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[string](func(yield func(Result[string]) bool) {
			// no items
		})

		inters := seq.Intersperse(",")
		firstErr := inters.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}

		collected := inters.Ok().Collect()
		if len(collected) != 0 {
			t.Errorf("expected 0 items, got %d", len(collected))
		}
	})

	t.Run("single Ok item", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(42))
		})

		inters := seq.Intersperse(999)
		firstErr := inters.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := inters.Ok().Collect()
		// Only one item => no separators
		if len(collected) != 1 || collected[0] != 42 {
			t.Errorf("expected [42], got %v", collected)
		}
	})

	t.Run("multiple Ok items", func(t *testing.T) {
		seq := SeqResult[string](func(yield func(Result[string]) bool) {
			yield(Ok("a"))
			yield(Ok("b"))
			yield(Ok("c"))
		})

		inters := seq.Intersperse(",")
		firstErr := inters.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok, got Err: %v", firstErr.Some())
		}

		collected := inters.Ok().Collect()
		// We expect: ["a", ",", "b", ",", "c"]
		want := []string{"a", ",", "b", ",", "c"}
		if len(collected) != len(want) {
			t.Errorf("expected %d items, got %d", len(want), len(collected))
		} else {
			for i, v := range collected {
				if v != want[i] {
					t.Errorf("index %d: want %q, got %q", i, want[i], v)
				}
			}
		}
	})

	t.Run("error in the middle", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("boom")))
			yield(Ok(2)) // never reached, iteration stops
		})

		inters := seq.Intersperse(999)
		firstErr := inters.FirstErr()
		if firstErr.IsNone() {
			collected := inters.Ok().Collect()
			t.Fatalf("expected Err, got Ok(%v)", collected)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error message \"boom\", got %q", errMsg)
		}
	})

	t.Run("error first", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](errors.New("immediate failure")))
			yield(Ok(10)) // never reached
		})

		inters := seq.Intersperse(999)
		firstErr := inters.FirstErr()
		if firstErr.IsNone() {
			collected := inters.Ok().Collect()
			t.Fatalf("expected Err, got Ok(%v)", collected)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "immediate failure" {
			t.Errorf("expected \"immediate failure\", got %q", errMsg)
		}
	})
}

func TestSeqResultInspect(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		var inspectedCount int
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		inspectedSeq := seq.Inspect(func(v int) {
			inspectedCount++
		})

		firstErr := inspectedSeq.FirstErr()
		if firstErr.IsSome() {
			t.Fatalf("expected Ok([]), got Err: %v", firstErr.Some())
		}
		if inspectedCount != 0 {
			t.Errorf("expected fn to be called 0 times, got %d", inspectedCount)
		}
	})

	t.Run("all Ok", func(t *testing.T) {
		var inspected []int

		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		inspectedSeq := seq.Inspect(func(v int) {
			inspected = append(inspected, v)
		})

		allResults := inspectedSeq.Collect()

		// Check for errors
		var firstErr error
		for _, r := range allResults {
			if r.IsErr() {
				firstErr = r.Err()
				break
			}
		}
		if firstErr != nil {
			t.Fatalf("expected Ok, got Err: %v", firstErr)
		}

		// Collect Ok values
		var collected []int
		for _, r := range allResults {
			if r.IsOk() {
				collected = append(collected, r.Ok())
			}
		}

		// We expect Inspect to be called once per Ok value before yielding them.
		if len(inspected) != 3 {
			t.Errorf("expected 3 items inspected, got %d", len(inspected))
		}
		wantInspected := []int{10, 20, 30}
		for i, v := range inspected {
			if v != wantInspected[i] {
				t.Errorf("index %d: expected %d, got %d", i, wantInspected[i], v)
			}
		}

		// Also check that the final collected items are unchanged.
		if len(collected) != 3 {
			t.Errorf("expected 3 collected items, got %d", len(collected))
		}
		wantCollected := []int{10, 20, 30}
		for i, v := range collected {
			if v != wantCollected[i] {
				t.Errorf("index %d: expected %d, got %d", i, wantCollected[i], v)
			}
		}
	})

	t.Run("error in the middle", func(t *testing.T) {
		var inspected []int

		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			if !yield(Err[int](errors.New("boom"))) {
				return
			}
			yield(Ok(2)) // not reached
		})

		inspectedSeq := seq.Inspect(func(v int) {
			inspected = append(inspected, v)
		})

		firstErr := inspectedSeq.FirstErr()
		if firstErr.IsNone() {
			collected := inspectedSeq.Ok().Collect()
			t.Fatalf("expected Err, got Ok(%v)", collected)
		}
		if len(inspected) != 1 {
			t.Errorf("expected 1 item inspected before the error, got %d", len(inspected))
		}
		if inspected[0] != 1 {
			t.Errorf("expected inspected[0] = 1, got %d", inspected[0])
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "boom" {
			t.Errorf("expected error \"boom\", got %q", errMsg)
		}
	})

	t.Run("error first", func(t *testing.T) {
		var inspectedCount int

		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			if !yield(Err[int](errors.New("immediate error"))) {
				return
			}
			yield(Ok(999)) // never reached
		})

		inspectedSeq := seq.Inspect(func(v int) {
			inspectedCount++
		})

		firstErr := inspectedSeq.FirstErr()
		if firstErr.IsNone() {
			collected := inspectedSeq.Ok().Collect()
			t.Fatalf("expected Err, got Ok(%v)", collected)
		}
		if inspectedCount != 0 {
			t.Errorf("expected 0 items inspected, got %d", inspectedCount)
		}
		errMsg := firstErr.Some().Error()
		if errMsg != "immediate error" {
			t.Errorf("expected \"immediate error\", got %q", errMsg)
		}
	})
}

func TestSeqResultFind(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// no items
		})

		res := seq.Find(func(v int) bool { return v == 100 })
		if res.IsErr() {
			t.Fatalf("expected Ok(None), got Err(%v)", res.Err())
		}

		opt := res.UnwrapOr(None[int]())
		if opt.IsSome() {
			t.Errorf("expected None, got Some(%v)", opt.Unwrap())
		}
	})

	t.Run("all Ok, no match", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		res := seq.Find(func(v int) bool { return v > 100 })
		if res.IsErr() {
			t.Fatalf("expected Ok(None), got Err(%v)", res.Err())
		}

		opt := res.UnwrapOr(None[int]())
		if opt.IsSome() {
			t.Errorf("expected None, got Some(%v)", opt.Unwrap())
		}
	})

	t.Run("match found", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		res := seq.Find(func(v int) bool { return v == 20 })
		if res.IsErr() {
			t.Fatalf("expected Ok(Some(20)), got Err(%v)", res.Err())
		}

		opt := res.UnwrapOr(None[int]())
		if !opt.IsSome() {
			t.Fatalf("expected Some(20), got None")
		}
		if opt.Unwrap() != 20 {
			t.Errorf("expected 20, got %d", opt.Unwrap())
		}
	})

	t.Run("encounter error", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			if !yield(Err[int](errors.New("boom"))) {
				return
			}
			yield(Ok(2))
		})

		res := seq.Find(func(v int) bool { return v == 2 })
		if !res.IsErr() {
			t.Fatalf("expected Err(\"boom\"), got Ok(%v)", res.UnwrapOr(None[int]()))
		}

		errMsg := res.Err().Error()
		if errMsg != "boom" {
			t.Errorf("expected error \"boom\", got %q", errMsg)
		}
	})
}

func TestSeqResultPull(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// No items to yield
		})

		next, stop := seq.Pull()
		defer stop()

		var items []Result[int]
		for {
			res, ok := next()
			if !ok {
				break
			}
			items = append(items, res)
		}

		if len(items) != 0 {
			t.Errorf("expected 0 items, got %d", len(items))
		}
	})

	t.Run("all Ok", func(t *testing.T) {
		seq := SeqResult[string](func(yield func(Result[string]) bool) {
			if !yield(Ok("alpha")) {
				return
			}
			if !yield(Ok("beta")) {
				return
			}
			if !yield(Ok("gamma")) {
				return
			}
		})

		next, stop := seq.Pull()
		defer stop()

		var items []Result[string]
		for {
			res, ok := next()
			if !ok {
				break
			}
			items = append(items, res)
		}

		expected := []string{"alpha", "beta", "gamma"}
		if len(items) != len(expected) {
			t.Fatalf("expected %d items, got %d", len(expected), len(items))
		}

		for i, res := range items {
			if res.IsErr() {
				t.Errorf("unexpected error at index %d: %v", i, res.Err())
				continue
			}
			if res.Ok() != expected[i] {
				t.Errorf("at index %d: expected %q, got %q", i, expected[i], res.Ok())
			}
		}
	})

	t.Run("error in the middle", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			if !yield(Ok(1)) {
				return
			}
			if !yield(Err[int](errors.New("encountered error"))) {
				return
			}
			if !yield(Ok(2)) {
				return
			} // Should not be reached
		})

		next, stop := seq.Pull()
		defer stop()

		var items []Result[int]
		for {
			res, ok := next()
			if !ok {
				break
			}

			items = append(items, res)

			if res.IsErr() {
				break
			}
		}

		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}

		// First item should be Ok(1)
		if !items[0].IsOk() || items[0].Ok() != 1 {
			t.Errorf("expected first item to be Ok(1), got %v", items[0])
		}

		// Second item should be Err("encountered error")
		if !items[1].IsErr() || items[1].Err().Error() != "encountered error" {
			t.Errorf("expected second item to be Err(\"encountered error\"), got %v", items[1])
		}
	})

	t.Run("error first", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			if !yield(Err[int](errors.New("immediate failure"))) {
				return
			}
			if !yield(Ok(1)) {
				return
			} // Should not be reached
		})

		next, stop := seq.Pull()
		defer stop()

		var items []Result[int]
		for {
			res, ok := next()
			if !ok {
				break
			}

			items = append(items, res)

			if res.IsErr() {
				break
			}
		}

		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}

		if !items[0].IsErr() || items[0].Err().Error() != "immediate failure" {
			t.Errorf("expected first item to be Err(\"immediate failure\"), got %v", items[0])
		}
	})

	t.Run("stop early using stop function", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			if !yield(Ok(100)) {
				return
			}
			if !yield(Ok(200)) {
				return
			}
			if !yield(Ok(300)) {
				return
			}
		})

		next, stop := seq.Pull()

		var items []Result[int]
		for {
			res, ok := next()
			if !ok {
				break
			}
			items = append(items, res)
			if res.IsOk() && res.Ok() == 200 {
				stop() // Stop iteration early
			}
		}

		// Should have collected [Ok(100), Ok(200)]
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}

		if !items[0].IsOk() || items[0].Ok() != 100 {
			t.Errorf("expected first item to be Ok(100), got %v", items[0])
		}
		if !items[1].IsOk() || items[1].Ok() != 200 {
			t.Errorf("expected second item to be Ok(200), got %v", items[1])
		}
	})

	t.Run("multiple Pull calls", func(t *testing.T) {
		seq := SeqResult[string](func(yield func(Result[string]) bool) {
			if !yield(Ok("first")) {
				return
			}
			if !yield(Ok("second")) {
				return
			}
			if !yield(Ok("third")) {
				return
			}
		})

		next, stop := seq.Pull()
		defer stop()

		// First Pull call
		res1, ok1 := next()
		if !ok1 {
			t.Fatalf("expected first item, got none")
		}
		if !res1.IsOk() || res1.Ok() != "first" {
			t.Errorf("expected first item to be Ok(\"first\"), got %v", res1)
		}

		// Second Pull call
		res2, ok2 := next()
		if !ok2 {
			t.Fatalf("expected second item, got none")
		}
		if !res2.IsOk() || res2.Ok() != "second" {
			t.Errorf("expected second item to be Ok(\"second\"), got %v", res2)
		}

		// Third Pull call
		res3, ok3 := next()
		if !ok3 {
			t.Fatalf("expected third item, got none")
		}
		if !res3.IsOk() || res3.Ok() != "third" {
			t.Errorf("expected third item to be Ok(\"third\"), got %v", res3)
		}

		// Fourth Pull call (no more items)
		_, ok4 := next()
		if ok4 {
			t.Errorf("expected no more items, but got some")
		}
	})

	t.Run("concurrent Pull and Push", func(t *testing.T) {
		// Note: Go's testing framework runs tests sequentially,
		// and since Pull is designed to convert push to pull,
		// concurrent access might not be a typical use case.
		// This test ensures that Pull works correctly even if the consumer
		// stops early.

		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			for i := 1; i <= 5; i++ {
				if !yield(Ok(i)) {
					return
				}
			}
		})

		next, stop := seq.Pull()
		defer stop()

		var items []Result[int]
		for {
			res, ok := next()
			if !ok {
				break
			}
			items = append(items, res)
			if res.IsOk() && res.Ok() == 3 {
				stop() // Stop after 3
			}
		}

		// Should have [Ok(1), Ok(2), Ok(3)]
		if len(items) != 3 {
			t.Fatalf("expected 3 items, got %d", len(items))
		}

		for i, expected := range []int{1, 2, 3} {
			if !items[i].IsOk() || items[i].Ok() != expected {
				t.Errorf("at index %d: expected Ok(%d), got %v", i, expected, items[i])
			}
		}
	})
}

func TestSeqResultContext(t *testing.T) {
	t.Run("context cancellation stops iteration", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
			yield(Ok(4))
			yield(Ok(5))
		})

		ctx, cancel := context.WithCancel(context.Background())

		var collected []Result[int]
		iter := seq.Context(ctx)

		// Cancel context after processing 3 elements
		count := 0
		iter(func(v Result[int]) bool {
			collected = append(collected, v)
			count++
			if count == 3 {
				cancel()
			}
			return true
		})

		// Should have processed exactly 3 elements before cancellation
		if len(collected) != 3 {
			t.Errorf("Expected 3 elements, got %d: %v", len(collected), collected)
		}

		// Verify all collected elements are Ok and have expected values
		for i, result := range collected {
			if result.IsErr() {
				t.Errorf("Expected Ok result at index %d, got Err: %v", i, result.Err())
			} else if result.Ok() != i+1 {
				t.Errorf("Expected value %d at index %d, got %d", i+1, i, result.Ok())
			}
		}
	})

	t.Run("context cancellation with error in sequence", func(t *testing.T) {
		testErr := errors.New("test error")
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](testErr))
			yield(Ok(3))
		})

		ctx, cancel := context.WithCancel(context.Background())

		var collected []Result[int]
		iter := seq.Context(ctx)

		// Cancel context after processing 2 elements
		count := 0
		iter(func(v Result[int]) bool {
			collected = append(collected, v)
			count++
			if count == 2 {
				cancel()
			}
			return true
		})

		// Should have processed 2 elements before cancellation
		if len(collected) != 2 {
			t.Errorf("Expected 2 elements, got %d: %v", len(collected), collected)
		}

		if collected[0].IsErr() || collected[0].Ok() != 1 {
			t.Errorf("First element should be Ok(1), got %v", collected[0])
		}

		if collected[1].IsOk() || collected[1].Err().Error() != "test error" {
			t.Errorf("Second element should be Err(test error), got %v", collected[1])
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
		})

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var collected []Result[int]
		seq.Context(ctx)(func(v Result[int]) bool {
			collected = append(collected, v)
			return true
		})

		// Should collect nothing due to immediate cancellation
		if len(collected) != 0 {
			t.Errorf("Expected 0 elements due to cancelled context, got %d: %v", len(collected), collected)
		}
	})
}

func TestSeqResultFirst(t *testing.T) {
	t.Run("first Ok element exists", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		first := seq.First()

		if first.IsErr() {
			t.Errorf("Expected Ok, got Err: %v", first.Err())
		} else if first.Ok().IsNone() {
			t.Error("Expected Some value, got None")
		} else if first.Ok().Some() != 10 {
			t.Errorf("Expected 10, got %d", first.Ok().Some())
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {})

		first := seq.First()

		if first.IsErr() {
			t.Errorf("Expected Ok, got Err: %v", first.Err())
		} else if first.Ok().IsSome() {
			t.Errorf("Expected None for empty sequence, got Some(%v)", first.Ok().Some())
		}
	})

	t.Run("sequence with only Err", func(t *testing.T) {
		testErr := errors.New("test error")
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](testErr))
		})

		first := seq.First()

		if first.IsOk() {
			t.Errorf("Expected Err, got Ok: %v", first.Ok())
		}
	})

	t.Run("single Ok element", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(42))
		})

		first := seq.First()

		if first.IsErr() {
			t.Errorf("Expected Ok, got Err: %v", first.Err())
		} else if first.Ok().IsNone() {
			t.Error("Expected Some value, got None")
		} else if first.Ok().Some() != 42 {
			t.Errorf("Expected 42, got %d", first.Ok().Some())
		}
	})
}

func TestSeqResultLast(t *testing.T) {
	t.Run("last Ok element exists", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
		})

		last := seq.Last()

		if last.IsErr() {
			t.Errorf("Expected Ok, got Err: %v", last.Err())
		} else if last.Ok().IsNone() {
			t.Error("Expected Some value, got None")
		} else if last.Ok().Some() != 30 {
			t.Errorf("Expected 30, got %d", last.Ok().Some())
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {})

		last := seq.Last()

		if last.IsErr() {
			t.Errorf("Expected Ok, got Err: %v", last.Err())
		} else if last.Ok().IsSome() {
			t.Errorf("Expected None for empty sequence, got Some(%v)", last.Ok().Some())
		}
	})

	t.Run("sequence with only Err", func(t *testing.T) {
		testErr := errors.New("test error")
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](testErr))
		})

		last := seq.Last()

		if last.IsOk() {
			t.Errorf("Expected Err, got Ok: %v", last.Ok())
		}
	})

	t.Run("single Ok element", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(42))
		})

		last := seq.Last()

		if last.IsErr() {
			t.Errorf("Expected Ok, got Err: %v", last.Err())
		} else if last.Ok().IsNone() {
			t.Error("Expected Some value, got None")
		} else if last.Ok().Some() != 42 {
			t.Errorf("Expected 42, got %d", last.Ok().Some())
		}
	})
}

func TestSeqResultNth(t *testing.T) {
	t.Run("nth Ok element exists", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
			yield(Ok(30))
			yield(Ok(40))
			yield(Ok(50))
		})

		// Get the 2nd Ok element (0-indexed) - should be 30
		nth := seq.Nth(2)

		if nth.IsErr() {
			t.Errorf("Expected Ok result, got Err: %v", nth.Err())
		} else if nth.Ok().IsNone() {
			t.Error("Expected Some value, got None")
		} else if nth.Ok().Some() != 30 {
			t.Errorf("Expected 30, got %d", nth.Ok().Some())
		}
	})

	t.Run("nth element with error before nth", func(t *testing.T) {
		testErr := errors.New("test error")
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Err[int](testErr))
			yield(Ok(30))
		})

		// Should return the error before reaching index 2
		nth := seq.Nth(2)

		if nth.IsOk() {
			t.Errorf("Expected Err result, got Ok: %v", nth.Ok())
		} else if nth.Err().Error() != "test error" {
			t.Errorf("Expected 'test error', got %v", nth.Err())
		}
	})

	t.Run("nth element out of bounds", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
		})

		nth := seq.Nth(5)

		if nth.IsErr() {
			t.Errorf("Expected Ok result, got Err: %v", nth.Err())
		} else if nth.Ok().IsSome() {
			t.Errorf("Expected None for out of bounds index, got Some(%v)", nth.Ok().Some())
		}
	})

	t.Run("negative index", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(10))
			yield(Ok(20))
		})

		nth := seq.Nth(-1)

		if nth.IsErr() {
			t.Errorf("Expected Ok result, got Err: %v", nth.Err())
		} else if nth.Ok().IsSome() {
			t.Errorf("Expected None for negative index, got Some(%v)", nth.Ok().Some())
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// Empty sequence
		})

		nth := seq.Nth(0)

		if nth.IsErr() {
			t.Errorf("Expected Ok result, got Err: %v", nth.Err())
		} else if nth.Ok().IsSome() {
			t.Errorf("Expected None for empty sequence, got Some(%v)", nth.Ok().Some())
		}
	})

	t.Run("sequence with only errors", func(t *testing.T) {
		testErr := errors.New("first error")
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](testErr))
			yield(Err[int](errors.New("second error")))
		})

		nth := seq.Nth(0)

		if nth.IsOk() {
			t.Errorf("Expected Err result, got Ok: %v", nth.Ok())
		} else if nth.Err().Error() != "first error" {
			t.Errorf("Expected 'first error', got %v", nth.Err())
		}
	})
}

func TestSeqResultNext(t *testing.T) {
	t.Run("Next with mixed results", func(t *testing.T) {
		iter := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("error1")))
			yield(Ok(3))
		})

		// First result (Ok)
		first := iter.Next()
		if !first.IsSome() {
			t.Errorf("Expected Some(Result), got None")
		}

		firstResult := first.Some()
		if firstResult.IsErr() || firstResult.Ok() != 1 {
			t.Errorf("Expected Ok(1), got %v", firstResult)
		}

		// Second result (Err)
		second := iter.Next()
		if !second.IsSome() {
			t.Errorf("Expected Some(Result), got None")
		}

		secondResult := second.Some()
		if secondResult.IsOk() || secondResult.Err().Error() != "error1" {
			t.Errorf("Expected Err(error1), got %v", secondResult)
		}

		// Remaining results
		remaining := iter.Collect()
		remainingSlice := make([]int, 0)
		for _, r := range remaining {
			if r.IsOk() {
				remainingSlice = append(remainingSlice, r.Ok())
			}
		}
		if len(remainingSlice) != 1 {
			t.Errorf("Expected 1 remaining result, got %d", len(remainingSlice))
		}
		if len(remainingSlice) > 0 && remainingSlice[0] != 3 {
			t.Errorf("Expected value 3, got %v", remainingSlice[0])
		}
	})

	t.Run("Next with empty iterator", func(t *testing.T) {
		iter := SeqResult[int](func(func(Result[int]) bool) {
			// empty
		})

		result := iter.Next()
		if result.IsSome() {
			t.Errorf("Expected None, got Some(%v)", result.Some())
		}
	})

	t.Run("Next until exhausted", func(t *testing.T) {
		iter := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
		})

		// Extract all elements
		first := iter.Next()
		second := iter.Next()
		third := iter.Next()

		if !first.IsSome() {
			t.Errorf("Expected first to be Some(Result), got None")
		}
		if !second.IsSome() {
			t.Errorf("Expected second to be Some(Result), got None")
		}
		if third.IsSome() {
			t.Errorf("Expected third to be None, got Some(%v)", third.Some())
		}

		// Iterator should be empty now
		remaining := iter.Collect()
		remainingSlice := make([]int, 0)
		for _, r := range remaining {
			if r.IsOk() {
				remainingSlice = append(remainingSlice, r.Ok())
			}
		}
		if len(remainingSlice) != 0 {
			t.Errorf("Expected empty results, got %d", len(remainingSlice))
		}
	})
}

func TestSeqResultPartition(t *testing.T) {
	t.Run("mixed Ok and Err values", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("error1")))
			yield(Ok(2))
			yield(Err[int](errors.New("error2")))
			yield(Ok(3))
		})

		okValues, errValues := seq.Partition()

		expectedOk := Slice[int]{1, 2, 3}
		if len(okValues) != len(expectedOk) {
			t.Errorf("Expected %d Ok values, got %d", len(expectedOk), len(okValues))
		}
		for i, v := range expectedOk {
			if i >= len(okValues) || okValues[i] != v {
				t.Errorf("Expected Ok value %d at index %d, got %v", v, i, okValues[i])
			}
		}

		if len(errValues) != 2 {
			t.Errorf("Expected 2 error values, got %d", len(errValues))
		}
		if len(errValues) > 0 && errValues[0].Error() != "error1" {
			t.Errorf("Expected first error 'error1', got '%s'", errValues[0].Error())
		}
		if len(errValues) > 1 && errValues[1].Error() != "error2" {
			t.Errorf("Expected second error 'error2', got '%s'", errValues[1].Error())
		}
	})

	t.Run("only Ok values", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		okValues, errValues := seq.Partition()

		expectedOk := Slice[int]{1, 2, 3}
		if len(okValues) != len(expectedOk) {
			t.Errorf("Expected %d Ok values, got %d", len(expectedOk), len(okValues))
		}
		for i, v := range expectedOk {
			if i >= len(okValues) || okValues[i] != v {
				t.Errorf("Expected Ok value %d at index %d, got %v", v, i, okValues[i])
			}
		}

		if len(errValues) != 0 {
			t.Errorf("Expected 0 error values, got %d", len(errValues))
		}
	})

	t.Run("only Err values", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](errors.New("error1")))
			yield(Err[int](errors.New("error2")))
		})

		okValues, errValues := seq.Partition()

		if len(okValues) != 0 {
			t.Errorf("Expected 0 Ok values, got %d", len(okValues))
		}

		if len(errValues) != 2 {
			t.Errorf("Expected 2 error values, got %d", len(errValues))
		}
		if len(errValues) > 0 && errValues[0].Error() != "error1" {
			t.Errorf("Expected first error 'error1', got '%s'", errValues[0].Error())
		}
		if len(errValues) > 1 && errValues[1].Error() != "error2" {
			t.Errorf("Expected second error 'error2', got '%s'", errValues[1].Error())
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// empty
		})

		okValues, errValues := seq.Partition()

		if len(okValues) != 0 {
			t.Errorf("Expected 0 Ok values, got %d", len(okValues))
		}
		if len(errValues) != 0 {
			t.Errorf("Expected 0 error values, got %d", len(errValues))
		}
	})
}

func TestSeqResultOk(t *testing.T) {
	t.Run("filter Ok values from mixed sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("error1")))
			yield(Ok(2))
			yield(Err[int](errors.New("error2")))
			yield(Ok(3))
		})

		okSeq := seq.Ok()
		collected := okSeq.Collect()

		expected := Slice[int]{1, 2, 3}
		if len(collected) != len(expected) {
			t.Errorf("Expected %d Ok values, got %d", len(expected), len(collected))
		}
		for i, v := range expected {
			if i >= len(collected) || collected[i] != v {
				t.Errorf("Expected Ok value %d at index %d, got %v", v, i, collected[i])
			}
		}
	})

	t.Run("only Ok values", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		okSeq := seq.Ok()
		collected := okSeq.Collect()

		expected := Slice[int]{1, 2, 3}
		if len(collected) != len(expected) {
			t.Errorf("Expected %d Ok values, got %d", len(expected), len(collected))
		}
		for i, v := range expected {
			if i >= len(collected) || collected[i] != v {
				t.Errorf("Expected Ok value %d at index %d, got %v", v, i, collected[i])
			}
		}
	})

	t.Run("only Err values", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](errors.New("error1")))
			yield(Err[int](errors.New("error2")))
		})

		okSeq := seq.Ok()
		collected := okSeq.Collect()

		if len(collected) != 0 {
			t.Errorf("Expected 0 Ok values, got %d", len(collected))
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// empty
		})

		okSeq := seq.Ok()
		collected := okSeq.Collect()

		if len(collected) != 0 {
			t.Errorf("Expected 0 Ok values, got %d", len(collected))
		}
	})
}

func TestSeqResultErr(t *testing.T) {
	t.Run("filter Err values from mixed sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Err[int](errors.New("error1")))
			yield(Ok(2))
			yield(Err[int](errors.New("error2")))
			yield(Ok(3))
		})

		errSeq := seq.Err()
		collected := errSeq.Collect()

		if len(collected) != 2 {
			t.Errorf("Expected 2 error values, got %d", len(collected))
		}
		if len(collected) > 0 && collected[0].Error() != "error1" {
			t.Errorf("Expected first error 'error1', got '%s'", collected[0].Error())
		}
		if len(collected) > 1 && collected[1].Error() != "error2" {
			t.Errorf("Expected second error 'error2', got '%s'", collected[1].Error())
		}
	})

	t.Run("only Err values", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](errors.New("error1")))
			yield(Err[int](errors.New("error2")))
		})

		errSeq := seq.Err()
		collected := errSeq.Collect()

		if len(collected) != 2 {
			t.Errorf("Expected 2 error values, got %d", len(collected))
		}
		if len(collected) > 0 && collected[0].Error() != "error1" {
			t.Errorf("Expected first error 'error1', got '%s'", collected[0].Error())
		}
		if len(collected) > 1 && collected[1].Error() != "error2" {
			t.Errorf("Expected second error 'error2', got '%s'", collected[1].Error())
		}
	})

	t.Run("only Ok values", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		errSeq := seq.Err()
		collected := errSeq.Collect()

		if len(collected) != 0 {
			t.Errorf("Expected 0 error values, got %d", len(collected))
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// empty
		})

		errSeq := seq.Err()
		collected := errSeq.Collect()

		if len(collected) != 0 {
			t.Errorf("Expected 0 error values, got %d", len(collected))
		}
	})
}

func TestSeqResultFirstErr(t *testing.T) {
	t.Run("first error in mixed sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Err[int](errors.New("first error")))
			yield(Ok(3))
			yield(Err[int](errors.New("second error")))
		})

		firstErr := seq.FirstErr()

		if firstErr.IsNone() {
			t.Errorf("Expected Some(error), got None")
		}
		if firstErr.IsSome() && firstErr.Some().Error() != "first error" {
			t.Errorf("Expected 'first error', got '%s'", firstErr.Some().Error())
		}
	})

	t.Run("error at beginning", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](errors.New("immediate error")))
			yield(Ok(1))
			yield(Ok(2))
		})

		firstErr := seq.FirstErr()

		if firstErr.IsNone() {
			t.Errorf("Expected Some(error), got None")
		}
		if firstErr.IsSome() && firstErr.Some().Error() != "immediate error" {
			t.Errorf("Expected 'immediate error', got '%s'", firstErr.Some().Error())
		}
	})

	t.Run("no errors in sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Ok(1))
			yield(Ok(2))
			yield(Ok(3))
		})

		firstErr := seq.FirstErr()

		if firstErr.IsSome() {
			t.Errorf("Expected None, got Some(%v)", firstErr.Some())
		}
	})

	t.Run("only errors in sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			yield(Err[int](errors.New("error1")))
			yield(Err[int](errors.New("error2")))
			yield(Err[int](errors.New("error3")))
		})

		firstErr := seq.FirstErr()

		if firstErr.IsNone() {
			t.Errorf("Expected Some(error), got None")
		}
		if firstErr.IsSome() && firstErr.Some().Error() != "error1" {
			t.Errorf("Expected 'error1', got '%s'", firstErr.Some().Error())
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := SeqResult[int](func(yield func(Result[int]) bool) {
			// empty
		})

		firstErr := seq.FirstErr()

		if firstErr.IsSome() {
			t.Errorf("Expected None, got Some(%v)", firstErr.Some())
		}
	})
}

func TestFromResultChan(t *testing.T) {
	t.Run("closed empty channel", func(t *testing.T) {
		ch := make(chan Result[int])
		close(ch)

		collected := FromResultChan(ch).Collect()
		if len(collected) != 0 {
			t.Errorf("expected 0 results, got %d", len(collected))
		}
	})

	t.Run("buffered channel with Ok values", func(t *testing.T) {
		ch := make(chan Result[int], 3)
		ch <- Ok(1)
		ch <- Ok(2)
		ch <- Ok(3)
		close(ch)

		seq := FromResultChan(ch)
		collected := seq.Ok().Collect()

		want := []int{1, 2, 3}
		if len(collected) != len(want) {
			t.Fatalf("expected %d items, got %d", len(want), len(collected))
		}
		for i, v := range collected {
			if v != want[i] {
				t.Errorf("index %d: want %d, got %d", i, want[i], v)
			}
		}
	})

	t.Run("mixed Ok and Err", func(t *testing.T) {
		ch := make(chan Result[int], 4)
		ch <- Ok(10)
		ch <- Err[int](errors.New("fail"))
		ch <- Ok(30)
		ch <- Err[int](errors.New("fail2"))
		close(ch)

		okVals, errVals := FromResultChan(ch).Partition()

		if len(okVals) != 2 {
			t.Errorf("expected 2 ok values, got %d", len(okVals))
		}
		if len(errVals) != 2 {
			t.Errorf("expected 2 errors, got %d", len(errVals))
		}
	})

	t.Run("only errors", func(t *testing.T) {
		ch := make(chan Result[int], 2)
		ch <- Err[int](errors.New("e1"))
		ch <- Err[int](errors.New("e2"))
		close(ch)

		firstErr := FromResultChan(ch).FirstErr()
		if firstErr.IsNone() {
			t.Fatal("expected an error, got None")
		}
		if firstErr.Some().Error() != "e1" {
			t.Errorf("expected 'e1', got %q", firstErr.Some().Error())
		}
	})

	t.Run("async producer", func(t *testing.T) {
		ch := make(chan Result[int])
		go func() {
			defer close(ch)
			for i := range 5 {
				ch <- Ok(i * 10)
			}
		}()

		collected := FromResultChan(ch).Ok().Collect()
		want := []int{0, 10, 20, 30, 40}
		if len(collected) != len(want) {
			t.Fatalf("expected %d items, got %d", len(want), len(collected))
		}
		for i, v := range collected {
			if v != want[i] {
				t.Errorf("index %d: want %d, got %d", i, want[i], v)
			}
		}
	})

	t.Run("with pool Stream", func(t *testing.T) {
		p := pool.New[int]().Limit(3)
		ch := p.Stream(func() {
			for i := range 10 {
				p.Go(func() Result[int] {
					if i%3 == 0 {
						return Err[int](errors.New("divisible by 3"))
					}
					return Ok(i)
				})
			}
		})

		okVals, errVals := FromResultChan(ch).Partition()

		// i=0,3,6,9 fail → 4 errors; i=1,2,4,5,7,8 succeed → 6 ok
		if len(okVals) != 6 {
			t.Errorf("expected 6 ok values, got %d", len(okVals))
		}
		if len(errVals) != 4 {
			t.Errorf("expected 4 errors, got %d", len(errVals))
		}
	})

	t.Run("chaining SeqResult methods", func(t *testing.T) {
		ch := make(chan Result[int], 5)
		ch <- Ok(1)
		ch <- Ok(2)
		ch <- Ok(3)
		ch <- Ok(4)
		ch <- Ok(5)
		close(ch)

		collected := FromResultChan(ch).
			Filter(func(v int) bool { return v%2 != 0 }).
			Map(func(v int) int { return v * 100 }).
			Ok().
			Collect()

		want := []int{100, 300, 500}
		if len(collected) != len(want) {
			t.Fatalf("expected %d items, got %d", len(want), len(collected))
		}
		for i, v := range collected {
			if v != want[i] {
				t.Errorf("index %d: want %d, got %d", i, want[i], v)
			}
		}
	})

	t.Run("count", func(t *testing.T) {
		ch := make(chan Result[int], 3)
		ch <- Ok(1)
		ch <- Err[int](errors.New("x"))
		ch <- Ok(3)
		close(ch)

		count := FromResultChan(ch).Count()
		if count != 3 {
			t.Errorf("expected count 3, got %d", count)
		}
	})

	t.Run("Find on channel results", func(t *testing.T) {
		ch := make(chan Result[int], 5)
		ch <- Ok(10)
		ch <- Ok(20)
		ch <- Ok(30)
		ch <- Ok(40)
		ch <- Ok(50)
		close(ch)

		res := FromResultChan(ch).Find(func(v int) bool { return v == 30 })
		if res.IsErr() {
			t.Fatalf("expected Ok, got Err: %v", res.Err())
		}
		opt := res.UnwrapOr(None[int]())
		if !opt.IsSome() || opt.Some() != 30 {
			t.Errorf("expected Some(30), got %v", opt)
		}
	})
}
