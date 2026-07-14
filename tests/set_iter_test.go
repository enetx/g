package g_test

import (
	"context"
	"testing"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestSet_Iter_Collect(t *testing.T) {
	original := g.NewSet[string]()
	original.Insert("a")
	original.Insert("b")
	original.Insert("c")

	collected := original.Iter().Collect()

	if collected.Len() != 3 {
		t.Errorf("Expected collected set to have 3 elements, got %d", collected.Len())
	}

	if !collected.Contains("a") || !collected.Contains("b") || !collected.Contains("c") {
		t.Error("Collected set should contain all original set elements")
	}
}

func TestSet_Iter_Filter(t *testing.T) {
	set := g.NewSet[int]()
	set.Insert(1)
	set.Insert(2)
	set.Insert(3)
	set.Insert(4)
	set.Insert(5)

	filtered := set.Iter().
		Filter(func(v int) bool { return v%2 == 0 }).
		Collect()

	if filtered.Len() != 2 {
		t.Errorf("Expected 2 even numbers, got %d", filtered.Len())
	}

	// Check that only even numbers are present
	slice := filtered.Slice()
	for _, num := range slice {
		if num%2 != 0 {
			t.Errorf("Filtered result should only contain even numbers, found %d", num)
		}
	}
}

func TestSet_Iter_Count(t *testing.T) {
	set := g.NewSet[string]()
	set.Insert("one")
	set.Insert("two")
	set.Insert("three")

	count := set.Iter().Count()

	if count != 3 {
		t.Errorf("Expected count of 3, got %d", count)
	}
}

func TestSet_Iter_Find(t *testing.T) {
	set := g.NewSet[int]()
	set.Insert(1)
	set.Insert(3)
	set.Insert(5)
	set.Insert(7)

	hasEven := set.Iter().Find(func(v int) bool { return v%2 == 0 })
	if hasEven.IsSome() {
		t.Error("Set of odd numbers should not have any even numbers")
	}

	hasOdd := set.Iter().Find(func(v int) bool { return v%2 == 1 })
	if hasOdd.IsNone() {
		t.Error("Set of odd numbers should have at least one odd number")
	}
}

func TestSet_Iter_Filter_AllEven(t *testing.T) {
	set := g.NewSet[int]()
	set.Insert(2)
	set.Insert(4)
	set.Insert(6)
	set.Insert(8)

	// Filter to get only even numbers - should be all of them
	evenSet := set.Iter().Filter(func(v int) bool { return v%2 == 0 }).Collect()
	if evenSet.Len() != 4 {
		t.Error("Set of even numbers should have all even numbers")
	}

	// Filter to get only odd numbers - should be empty
	oddSet := set.Iter().Filter(func(v int) bool { return v%2 == 1 }).Collect()
	if oddSet.Len() != 0 {
		t.Error("Set of even numbers should not have any odd numbers")
	}
}

func TestSet_Iter_ToSlice(t *testing.T) {
	set := g.NewSet[string]()
	set.Insert("a")
	set.Insert("b")
	set.Insert("c")
	set.Insert("d")
	set.Insert("e")

	slice := set.Slice()

	if len(slice) != 5 {
		t.Errorf("Expected 5 elements in slice, got %d", len(slice))
	}
}

func TestSet_Iter_EmptySet(t *testing.T) {
	set := g.NewSet[string]()

	count := set.Iter().Count()
	if count != 0 {
		t.Errorf("Empty set iterator count should be 0, got %d", count)
	}

	collected := set.Iter().Collect()
	if collected.Len() != 0 {
		t.Error("Empty set iterator should collect empty set")
	}

	hasAny := set.Iter().Find(func(v string) bool { return true })
	if hasAny.IsSome() {
		t.Error("Empty set should not have any elements")
	}
}

func TestSet_Iter_Map(t *testing.T) {
	set := g.NewSet[int]()
	set.Insert(1)
	set.Insert(2)
	set.Insert(3)

	doubled := set.Iter().
		Map(func(v int) int { return v * 2 }).
		Collect()

	if doubled.Len() != 3 {
		t.Errorf("Expected 3 doubled elements, got %d", doubled.Len())
	}

	if !doubled.Contains(2) || !doubled.Contains(4) || !doubled.Contains(6) {
		t.Error("Map should double all values")
	}
}

func TestSet_Iter_Chain(t *testing.T) {
	set1 := g.NewSet[string]()
	set1.Insert("a")
	set1.Insert("b")

	set2 := g.NewSet[string]()
	set2.Insert("c")
	set2.Insert("d")

	chained := set1.Iter().Chain(set2.Iter()).Collect()

	if chained.Len() != 4 {
		t.Errorf("Expected 4 elements after chaining, got %d", chained.Len())
	}

	expected := []string{"a", "b", "c", "d"}
	for _, exp := range expected {
		if !chained.Contains(exp) {
			t.Errorf("Chained iterator should contain %s", exp)
		}
	}
}

func TestSet_Iter_Pull(t *testing.T) {
	set := g.NewSet[string]()
	set.Insert("a")
	set.Insert("b")
	set.Insert("c")

	iter := set.Iter()
	next, stop := iter.Pull()
	defer stop()

	count := 0
	seen := make(map[string]bool)

	for {
		value, ok := next()
		if !ok {
			break
		}
		count++
		seen[value] = true
	}

	if count != 3 {
		t.Errorf("Expected to pull 3 values, got %d", count)
	}

	// Verify all original values were seen
	if !seen["a"] || !seen["b"] || !seen["c"] {
		t.Errorf("Not all set values were pulled: %v", seen)
	}
}

func TestSetIterContext(t *testing.T) {
	t.Run("context cancellation stops iteration", func(t *testing.T) {
		set := g.NewSet[int]()
		set.Insert(1)
		set.Insert(2)
		set.Insert(3)
		set.Insert(4)
		set.Insert(5)

		ctx, cancel := context.WithCancel(context.Background())

		var collected []int
		iter := set.Iter().Context(ctx)

		// Cancel context after processing 3 elements
		count := 0
		iter(func(v int) bool {
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

		// Verify all collected elements are from the original set
		for _, val := range collected {
			if !set.Contains(val) {
				t.Errorf("Collected value %d is not in original set", val)
			}
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		set := g.NewSet[int]()
		set.Insert(1)
		set.Insert(2)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var collected []int
		set.Iter().Context(ctx)(func(v int) bool {
			collected = append(collected, v)
			return true
		})

		// Should collect nothing due to immediate cancellation
		if len(collected) != 0 {
			t.Errorf("Expected 0 elements due to cancelled context, got %d: %v", len(collected), collected)
		}
	})
}

func TestSetIterTake(t *testing.T) {
	t.Run("take first n elements", func(t *testing.T) {
		set := g.NewSet[int]()
		set.Insert(1)
		set.Insert(2)
		set.Insert(3)
		set.Insert(4)
		set.Insert(5)

		taken := set.Iter().Take(3).Collect()

		if taken.Len() != 3 {
			t.Errorf("Expected 3 elements, got %d", taken.Len())
		}

		// Verify all taken elements are from the original set
		taken.Iter().ForEach(func(v int) {
			if !set.Contains(v) {
				t.Errorf("Taken element %d is not in original set", v)
			}
		})
	})

	t.Run("take more than available", func(t *testing.T) {
		set := g.NewSet[int]()
		set.Insert(1)
		set.Insert(2)

		taken := set.Iter().Take(5).Collect()

		if taken.Len() != 2 {
			t.Errorf("Expected 2 elements (all available), got %d", taken.Len())
		}
	})

	t.Run("take zero elements", func(t *testing.T) {
		set := g.NewSet[int]()
		set.Insert(1)
		set.Insert(2)

		taken := set.Iter().Take(0).Collect()

		if taken.Len() != 0 {
			t.Errorf("Expected 0 elements, got %d", taken.Len())
		}
	})
}

func TestSetIterNth(t *testing.T) {
	t.Run("nth element exists", func(t *testing.T) {
		set := g.NewSet[int]()
		set.Insert(10)
		set.Insert(20)
		set.Insert(30)
		set.Insert(40)
		set.Insert(50)

		// Get the 2nd element (0-indexed)
		nth := set.Iter().Nth(2)

		if nth.IsNone() {
			t.Error("Expected Some value, got None")
		} else {
			val := nth.Some()
			if !set.Contains(val) {
				t.Errorf("Nth element %d is not in original set", val)
			}
		}
	})

	t.Run("nth element out of bounds", func(t *testing.T) {
		set := g.NewSet[int]()
		set.Insert(1)
		set.Insert(2)

		nth := set.Iter().Nth(5)

		if nth.IsSome() {
			t.Errorf("Expected None for out of bounds index, got Some(%v)", nth.Some())
		}
	})

	t.Run("negative index", func(t *testing.T) {
		set := g.NewSet[int]()
		set.Insert(1)
		set.Insert(2)

		nth := set.Iter().Nth(-1)

		if nth.IsSome() {
			t.Errorf("Expected None for negative index, got Some(%v)", nth.Some())
		}
	})

	t.Run("empty set", func(t *testing.T) {
		set := g.NewSet[int]()

		nth := set.Iter().Nth(0)

		if nth.IsSome() {
			t.Errorf("Expected None for empty set, got Some(%v)", nth.Some())
		}
	})
}

func TestSeqSetNext(t *testing.T) {
	t.Run("Next with non-empty iterator", func(t *testing.T) {
		set := g.SetOf(1, 2, 3, 4, 5)
		iter := set.Iter()

		// Extract first element (set order is not guaranteed)
		first := iter.Next()
		if !first.IsSome() {
			t.Errorf("Expected Some(value), got None")
		}

		firstValue := first.Some()
		if !set.Contains(firstValue) {
			t.Errorf("First value %v should be in original set", firstValue)
		}

		// Extract second element
		second := iter.Next()
		if !second.IsSome() {
			t.Errorf("Expected Some(value), got None")
		}

		secondValue := second.Some()
		if !set.Contains(secondValue) || secondValue == firstValue {
			t.Errorf("Second value %v should be different from first and in original set", secondValue)
		}

		// Remaining elements
		remaining := iter.Collect()
		if remaining.Len() != 3 {
			t.Errorf("Expected 3 remaining elements, got %d", remaining.Len())
		}
	})

	t.Run("Next with empty iterator", func(t *testing.T) {
		set := g.NewSet[int]()
		iter := set.Iter()

		result := iter.Next()
		if result.IsSome() {
			t.Errorf("Expected None, got Some(%v)", result.Some())
		}
	})

	t.Run("Next until exhausted", func(t *testing.T) {
		set := g.SetOf(1, 2)
		iter := set.Iter()

		// Extract all elements
		first := iter.Next()
		second := iter.Next()
		third := iter.Next()

		if !first.IsSome() {
			t.Errorf("Expected first to be Some(value), got None")
		}
		if !second.IsSome() {
			t.Errorf("Expected second to be Some(value), got None")
		}
		if third.IsSome() {
			t.Errorf("Expected third to be None, got Some(%v)", third.Some())
		}

		// Iterator should be empty now
		remaining := iter.Collect()
		if remaining.Len() != 0 {
			t.Errorf("Expected empty set, got length %d", remaining.Len())
		}
	})
}

func TestSet_Iter_All(t *testing.T) {
	tests := []struct {
		name string
		set  g.Set[int]
		fn   func(int) bool
		want bool
	}{
		{"all_positive", g.SetOf(1, 2, 3, 4), func(v int) bool { return v > 0 }, true},
		{"not_all_even", g.SetOf(2, 4, 5, 6), func(v int) bool { return v%2 == 0 }, false},
		{"empty_is_true", g.NewSet[int](), func(int) bool { return false }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Iter().All(tt.fn); got != tt.want {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Iter_Any(t *testing.T) {
	tests := []struct {
		name string
		set  g.Set[int]
		fn   func(int) bool
		want bool
	}{
		{"has_even", g.SetOf(1, 3, 4, 5), func(v int) bool { return v%2 == 0 }, true},
		{"none_negative", g.SetOf(1, 2, 3), func(v int) bool { return v < 0 }, false},
		{"empty_is_false", g.NewSet[int](), func(int) bool { return true }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Iter().Any(tt.fn); got != tt.want {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Iter_Fold(t *testing.T) {
	set := g.SetOf(1, 2, 3, 4, 5)

	sum := set.Iter().Fold(0, func(acc, val int) int { return acc + val })
	if sum != 15 {
		t.Errorf("Fold sum = %d, want 15", sum)
	}

	empty := g.NewSet[int]()
	if got := empty.Iter().Fold(42, func(acc, val int) int { return acc + val }); got != 42 {
		t.Errorf("Fold on empty = %d, want 42", got)
	}
}

func TestSet_Iter_Reduce(t *testing.T) {
	set := g.SetOf(1, 2, 3, 4, 5)

	product := set.Iter().Reduce(func(a, b int) int { return a * b })
	if !product.IsSome() || product.Some() != 120 {
		t.Errorf("Reduce product = %v, want Some(120)", product)
	}

	empty := g.NewSet[int]()
	if empty.Iter().Reduce(func(a, b int) int { return a + b }).IsSome() {
		t.Error("Reduce on empty set should return None")
	}
}

func TestSet_Iter_Skip(t *testing.T) {
	set := g.SetOf(1, 2, 3, 4, 5, 6)

	skipped := set.Iter().Skip(4).Collect()
	if skipped.Len() != 2 {
		t.Errorf("Expected 2 elements after Skip(4), got %d", skipped.Len())
	}

	// Every remaining element must be a member of the original set.
	skipped.Iter().ForEach(func(v int) {
		if !set.Contains(v) {
			t.Errorf("Skip produced element %d not in original set", v)
		}
	})

	none := set.Iter().Skip(100).Collect()
	if none.Len() != 0 {
		t.Errorf("Expected empty set after Skip(100), got %d", none.Len())
	}
}

func TestSet_Iter_FilterMap(t *testing.T) {
	set := g.SetOf(1, 2, 3, 4, 5)

	result := set.Iter().FilterMap(func(n int) g.Option[int] {
		if n%2 == 0 {
			return g.Some(n * 10)
		}
		return g.None[int]()
	}).Collect()

	if result.Len() != 2 {
		t.Errorf("Expected 2 mapped elements, got %d", result.Len())
	}

	if !result.Contains(20) || !result.Contains(40) {
		t.Errorf("Expected {20, 40}, got %v", result)
	}
}

func TestSeqSetConsistency127(t *testing.T) {
	s := g.SetOf(1, 2, 3, 4, 5)

	// Fold with a foreign accumulator type
	joined := s.Iter().Fold("", func(acc string, v int) string { return acc + "x" })
	if len(joined) != 5 {
		t.Errorf("Fold cross-type = %q, want 5 x's", joined)
	}

	// Difference / Intersection on the iterator
	if got := s.Iter().Difference(g.SetOf(1, 2, 3, 4)).Collect(); !got.Eq(g.SetOf(5)) {
		t.Errorf("Difference = %v, want Set{5}", got)
	}
	if got := s.Iter().Intersection(g.SetOf(4, 5, 6)).Collect(); !got.Eq(g.SetOf(4, 5)) {
		t.Errorf("Intersection = %v, want Set{4, 5}", got)
	}

	// Partition
	even, odd := s.Iter().Partition(func(v int) bool { return v%2 == 0 })
	if !even.Eq(g.SetOf(2, 4)) || !odd.Eq(g.SetOf(1, 3, 5)) {
		t.Errorf("Partition = %v / %v", even, odd)
	}
}

func TestSet_Iter_CounterBy(t *testing.T) {
	set := g.SetOf("a", "bb", "cc", "d", "eee")

	counts := set.Iter().CounterBy(func(s string) int { return len(s) }).Collect()

	if counts.Len() != 3 {
		t.Fatalf("Expected 3 buckets, got %d: %v", counts.Len(), counts)
	}

	if got := counts.Get(1); got.IsNone() || got.Some() != 2 {
		t.Errorf("Expected 2 one-char elements, got %v", got)
	}

	if got := counts.Get(2); got.IsNone() || got.Some() != 2 {
		t.Errorf("Expected 2 two-char elements, got %v", got)
	}

	if got := counts.Get(3); got.IsNone() || got.Some() != 1 {
		t.Errorf("Expected 1 three-char element, got %v", got)
	}

	// Identity counting: a set has no duplicates, so every count is 1.
	identity := g.SetOf(1, 2, 3).Iter().CounterBy(func(v int) int { return v }).Collect()
	identity.Iter().ForEach(func(k int, count g.Int) {
		if count != 1 {
			t.Errorf("Expected count 1 for key %d, got %d", k, count)
		}
	})
}

func TestSet_Iter_Chan(t *testing.T) {
	set := g.SetOf(1, 2, 3)

	collected := g.NewSet[int]()
	for v := range set.Iter().Chan() {
		collected.Insert(v)
	}

	if !collected.Eq(set) {
		t.Errorf("Expected %v from channel, got %v", set, collected)
	}

	// A cancelled context stops the stream early.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	count := 0
	for range set.Iter().Chan(ctx) {
		count++
	}

	if count == 3 {
		t.Error("Expected cancelled context to stop the channel before all elements")
	}
}

func TestSet_Iter_MaxBy(t *testing.T) {
	max := g.SetOf(3, 1, 4, 1, 5).Iter().MaxBy(cmp.Cmp[int])
	if max.IsNone() || max.Some() != 5 {
		t.Errorf("Expected Some(5), got %v", max)
	}

	empty := g.NewSet[int]().Iter().MaxBy(cmp.Cmp[int])
	if empty.IsSome() {
		t.Errorf("Expected None for empty set, got %v", empty)
	}
}

func TestSet_Iter_MinBy(t *testing.T) {
	min := g.SetOf(3, 1, 4, 5).Iter().MinBy(cmp.Cmp[int])
	if min.IsNone() || min.Some() != 1 {
		t.Errorf("Expected Some(1), got %v", min)
	}

	empty := g.NewSet[int]().Iter().MinBy(cmp.Cmp[int])
	if empty.IsSome() {
		t.Errorf("Expected None for empty set, got %v", empty)
	}
}

func TestSet_Iter_First(t *testing.T) {
	set := g.SetOf(1, 2, 3, 4, 5)

	// Iteration order is nondeterministic, so First may yield any element of the set.
	first := set.Iter().First()
	if first.IsNone() {
		t.Fatal("Expected Some for non-empty set, got None")
	}

	if !set.Contains(first.Some()) {
		t.Errorf("Expected First to yield an element of the set, got %v", first.Some())
	}

	single := g.SetOf(42).Iter().First()
	if single.IsNone() || single.Some() != 42 {
		t.Errorf("Expected Some(42), got %v", single)
	}

	empty := g.NewSet[int]().Iter().First()
	if empty.IsSome() {
		t.Errorf("Expected None for empty set, got %v", empty)
	}
}

func TestSet_Iter_Last(t *testing.T) {
	set := g.SetOf(1, 2, 3, 4, 5)

	// Iteration order is nondeterministic, so Last may yield any element of the set.
	last := set.Iter().Last()
	if last.IsNone() {
		t.Fatal("Expected Some for non-empty set, got None")
	}

	if !set.Contains(last.Some()) {
		t.Errorf("Expected Last to yield an element of the set, got %v", last.Some())
	}

	single := g.SetOf(42).Iter().Last()
	if single.IsNone() || single.Some() != 42 {
		t.Errorf("Expected Some(42), got %v", single)
	}

	empty := g.NewSet[int]().Iter().Last()
	if empty.IsSome() {
		t.Errorf("Expected None for empty set, got %v", empty)
	}
}

func TestSet_Iter_StepBy(t *testing.T) {
	set := g.SetOf(1, 2, 3, 4, 5)

	// Step 1 yields every element.
	all := set.Iter().StepBy(1).Collect()
	if all.Ne(set) {
		t.Errorf("Expected StepBy(1) to yield all elements, got %v", all)
	}

	// Step 2 over 5 elements yields 3 of them; which ones depends on iteration order.
	stepped := set.Iter().StepBy(2).Collect()
	if stepped.Len() != 3 {
		t.Errorf("Expected 3 elements from StepBy(2), got %d", stepped.Len())
	}

	if !stepped.Iter().All(set.Contains) {
		t.Errorf("Expected StepBy(2) to yield elements of the set, got %v", stepped)
	}

	// Step larger than the set yields only the first element.
	one := set.Iter().StepBy(10).Collect()
	if one.Len() != 1 {
		t.Errorf("Expected 1 element from StepBy(10), got %d", one.Len())
	}

	empty := g.NewSet[int]().Iter().StepBy(2).Collect()
	if !empty.IsEmpty() {
		t.Errorf("Expected empty result for empty set, got %v", empty)
	}
}

func TestSeqSetSumBy(t *testing.T) {
	set := g.SetOf[g.Int](1, 2, 3, 4)
	if got := set.Iter().SumBy(func(v g.Int) g.Int { return v }); got != 10 {
		t.Errorf("SeqSet.SumBy = %d, want 10", got)
	}
}

func TestSeqSetTryMap(t *testing.T) {
	got := g.SetOf[g.String]("1", "2", "3").Iter().
		TryMap(g.String.TryInt).
		SumBy(func(v g.Int) g.Int { return v })
	if got.IsErr() || got.Ok() != 6 {
		t.Fatalf("Set.TryMap = %v", got)
	}
}
