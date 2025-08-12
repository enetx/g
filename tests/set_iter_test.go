package g_test

import (
	"testing"

	"github.com/enetx/g"
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
	slice := filtered.ToSlice()
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

	slice := set.ToSlice()

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
