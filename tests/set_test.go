package g_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/enetx/g"
)

func TestSetOf(t *testing.T) {
	// Test empty values
	emptySet := SetOf[int]()
	if emptySet.Len() != 0 {
		t.Errorf("Expected empty set size to be 0, got %d", emptySet.Len())
	}

	// Test single value
	singleSet := SetOf(42)
	if singleSet.Len() != 1 {
		t.Errorf("Expected single set size to be 1, got %d", singleSet.Len())
	}
	if !singleSet.Contains(42) {
		t.Errorf("Expected single set to contain value 42")
	}

	// Test multiple values
	multiSet := SetOf(1, 2, 3, 4, 5)
	expectedValues := []int{1, 2, 3, 4, 5}
	for _, v := range expectedValues {
		if !multiSet.Contains(v) {
			t.Errorf("Expected multi set to contain value %d", v)
		}
	}

	// Test duplicate values
	duplicateSet := SetOf(1, 1, 2, 2, 3, 3)
	if duplicateSet.Len() != 3 {
		t.Errorf("Expected duplicate set size to be 3, got %d", duplicateSet.Len())
	}
}

func TestSetDifference(t *testing.T) {
	set1 := SetOf(1, 2, 3, 4)
	set2 := SetOf(3, 4, 5, 6)
	set5 := SetOf(1, 2)
	set6 := SetOf(2, 3, 4)

	set3 := set1.Difference(set2).Collect()
	set4 := set2.Difference(set1).Collect()
	set7 := set5.Difference(set6).Collect()
	set8 := set6.Difference(set5).Collect()

	if set3.Len() != 2 || set3.Ne(SetOf(1, 2)) {
		t.Errorf("Unexpected result: %v", set3)
	}

	if set4.Len() != 2 || set4.Ne(SetOf(5, 6)) {
		t.Errorf("Unexpected result: %v", set4)
	}

	if set7.Len() != 1 || set7.Ne(SetOf(1)) {
		t.Errorf("Unexpected result: %v", set7)
	}

	if set8.Len() != 2 || set8.Ne(SetOf(3, 4)) {
		t.Errorf("Unexpected result: %v", set8)
	}
}

func TestSetSymmetricDifference(t *testing.T) {
	set1 := NewSet[int](10)
	set2 := set1.Clone()
	result := set1.SymmetricDifference(set2).Collect()

	if !result.Empty() {
		t.Errorf("SymmetricDifference between equal sets should be empty, got %v", result)
	}

	set1 = SetOf(0, 1, 2, 3, 4)
	set2 = SetOf(5, 6, 7, 8, 9)
	result = set1.SymmetricDifference(set2).Collect()
	expected := set1.Union(set2).Collect()

	if !result.Eq(expected) {
		t.Errorf(
			"SymmetricDifference between disjoint sets should be their union, expected %v but got %v",
			expected,
			result,
		)
	}

	set1 = SetOf(0, 1, 2, 3, 4, 5)
	set2 = SetOf(4, 5, 6, 7, 8)
	result = set1.SymmetricDifference(set2).Collect()
	expected2 := SetOf(0, 1, 2, 3, 6, 7, 8)

	if !result.Eq(expected2) {
		t.Errorf(
			"SymmetricDifference between sets with common elements should be correct, expected %v but got %v",
			expected,
			result,
		)
	}
}

func TestSetSubset(t *testing.T) {
	tests := []struct {
		name  string
		s     Set[int]
		other Set[int]
		want  bool
	}{
		{
			name:  "test_subset_1",
			s:     SetOf(1, 2, 3),
			other: SetOf(1, 2, 3, 4, 5),
			want:  true,
		},
		{
			name:  "test_subset_2",
			s:     SetOf(1, 2, 3, 4),
			other: SetOf(1, 2, 3),
			want:  false,
		},
		{
			name:  "test_subset_3",
			s:     SetOf(5, 4, 3, 2, 1),
			other: SetOf(1, 2, 3, 4, 5),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Subset(tt.other); got != tt.want {
				t.Errorf("Subset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetSuperset(t *testing.T) {
	tests := []struct {
		name  string
		s     Set[int]
		other Set[int]
		want  bool
	}{
		{
			name:  "test_superset_1",
			s:     SetOf(1, 2, 3, 4, 5),
			other: SetOf(1, 2, 3),
			want:  true,
		},
		{
			name:  "test_superset_2",
			s:     SetOf(1, 2, 3),
			other: SetOf(1, 2, 3, 4),
			want:  false,
		},
		{
			name:  "test_superset_3",
			s:     SetOf(1, 2, 3, 4, 5),
			other: SetOf(5, 4, 3, 2, 1),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Superset(tt.other); got != tt.want {
				t.Errorf("Superset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetEq(t *testing.T) {
	tests := []struct {
		name  string
		s     Set[int]
		other Set[int]
		want  bool
	}{
		{
			name:  "test_eq_1",
			s:     SetOf(1, 2, 3),
			other: SetOf(1, 2, 3),
			want:  true,
		},
		{
			name:  "test_eq_2",
			s:     SetOf(1, 2, 3),
			other: SetOf(1, 2, 4),
			want:  false,
		},
		{
			name:  "test_eq_3",
			s:     SetOf(1, 2, 3),
			other: SetOf(3, 2, 1),
			want:  true,
		},
		{
			name:  "test_eq_4",
			s:     SetOf(1, 2, 3, 4),
			other: SetOf(1, 2, 3),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Eq(tt.other); got != tt.want {
				t.Errorf("Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRemove(t *testing.T) {
	// Test case 1: Remove a single value
	set := SetOf(1, 2, 3)
	set.Remove(2)
	if _, ok := set[2]; ok {
		t.Errorf("Set should not contain value 2 after removal")
	}

	// Test case 2: Remove multiple values
	set2 := SetOf("a", "b", "c")
	set2.Remove("a", "c")
	if _, ok := set2["a"]; ok {
		t.Errorf("Set should not contain value 'a' after removal")
	}
	if _, ok := set2["c"]; ok {
		t.Errorf("Set should not contain value 'c' after removal")
	}

	// Test case 3: Remove non-existent value
	set3 := SetOf(1.1, 2.2)
	set3.Remove(3.3)
	if len(set3) != 2 {
		t.Errorf("Set should not be modified when removing non-existent value")
	}
}

func TestSetContainsAny(t *testing.T) {
	// Test case 1: Set contains some elements from another set
	set1 := SetOf(1, 2, 3)
	set2 := SetOf(2, 4, 6)
	if !set1.ContainsAny(set2) {
		t.Errorf("Expected Set to contain at least one element from the other Set")
	}

	// Test case 2: Set doesn't contain any elements from another set
	set3 := SetOf("a", "b")
	set4 := SetOf("c", "d", "e")
	if set3.ContainsAny(set4) {
		t.Errorf("Expected Set not to contain any elements from the other Set")
	}

	// Test case 3: Empty sets
	set5 := Set[float64]{}
	set6 := Set[float64]{}
	if set5.ContainsAny(set6) {
		t.Errorf("Expected empty sets not to contain any elements from each other")
	}
}

func TestSetContainsAll(t *testing.T) {
	// Test case 1: Set contains all elements from another set
	set1 := SetOf(1, 2, 3)
	set2 := SetOf(2, 3)
	if !set1.ContainsAll(set2) {
		t.Errorf("Expected Set to contain all elements from the other Set")
	}

	// Test case 2: Set doesn't contain all elements from another set
	set3 := SetOf("a", "b", "c")
	set4 := SetOf("b", "d")
	if set3.ContainsAll(set4) {
		t.Errorf("Expected Set not to contain all elements from the other Set")
	}

	// Test case 3: Empty sets
	set5 := SetOf[float64]()
	set6 := SetOf(1.1, 2.2)
	if set5.ContainsAll(set6) {
		t.Errorf("Expected empty Set not to contain all elements from another non-empty Set")
	}
}

func TestSetToSlice(t *testing.T) {
	// Test case 1: Set with elements
	set1 := SetOf(1, 2, 3)
	expected1 := Slice[int]{1, 2, 3}
	slice1 := set1.ToSlice()
	if len(slice1) != len(expected1) {
		t.Errorf("Expected length of slice to be %d, got %d", len(expected1), len(slice1))
	}

	// Test case 2: Empty Set
	set2 := NewSet[string]()
	expected2 := Slice[string]{}
	slice2 := set2.ToSlice()
	if len(slice2) != len(expected2) {
		t.Errorf("Expected length of slice to be %d, got %d", len(expected2), len(slice2))
	}
}

func TestSetString(t *testing.T) {
	// Test case 1: Set with elements
	set1 := NewSet[int]()
	set1.Insert(1)
	expected1 := "Set{1}"
	if str := set1.String(); str != expected1 {
		t.Errorf("Expected string representation to be %s, got %s", expected1, str)
	}

	// Test case 2: Empty Set
	set2 := NewSet[int]()
	expected2 := "Set{}"
	if str := set2.String(); str != expected2 {
		t.Errorf("Expected string representation to be %s, got %s", expected2, str)
	}
}

func TestSetClear(t *testing.T) {
	// Test case 1: Set with elements
	set1 := SetOf(1, 2, 3)
	set1.Clear()
	if len(set1) != 0 {
		t.Errorf("Expected Set to be empty after calling Clear()")
	}

	// Test case 2: Empty Set
	set2 := NewSet[int]()
	set2.Clear()
	if len(set2) != 0 {
		t.Errorf("Expected Set to remain empty after calling Clear() on an empty Set")
	}
}

func TestSetIntersection(t *testing.T) {
	// Test case 1: Set with elements
	set1 := SetOf(1, 2, 3, 4, 5)
	set2 := SetOf(4, 5, 6, 7, 8)
	expected := SetOf(4, 5)
	intersection := set1.Intersection(set2).Collect()
	if len(intersection) != len(expected) {
		t.Errorf("Expected intersection to have length %d, got %d", len(expected), len(intersection))
	}
	for k := range intersection {
		if _, exists := expected[k]; !exists {
			t.Errorf("Unexpected element in intersection: %d", k)
		}
	}

	// Test case 2: Empty Set
	set3 := SetOf("a", "b")
	set4 := NewSet[string]()
	intersection2 := set3.Intersection(set4).Collect()
	if len(intersection2) != 0 {
		t.Errorf("Expected intersection of an empty set to be empty")
	}

	// Test case 3: Both sets empty
	set5 := NewSet[int]()
	set6 := NewSet[int]()
	intersection = set5.Intersection(set6).Collect()
	if len(intersection) != 0 {
		t.Errorf("Expected intersection of two empty sets to be empty")
	}
}

func TestSetUnion(t *testing.T) {
	// Test case 1: Set with elements
	set1 := SetOf(1, 2, 3)
	set2 := SetOf(3, 4, 5)
	expected := SetOf(1, 2, 3, 4, 5)
	union := set1.Union(set2).Collect()
	if len(union) != len(expected) {
		t.Errorf("Expected union to have length %d, got %d", len(expected), len(union))
	}
	for k := range union {
		if _, exists := expected[k]; !exists {
			t.Errorf("Unexpected element in union: %d", k)
		}
	}

	// Test case 2: Empty Set
	set3 := SetOf("a", "b")
	set4 := NewSet[string]()
	union2 := set3.Union(set4).Collect()
	if len(union2) != len(set3) {
		t.Errorf("Expected union to be equal to the non-empty set")
	}

	// Test case 3: Both sets empty
	set5 := NewSet[int]()
	set6 := NewSet[int]()
	union = set5.Union(set6).Collect()
	if len(union) != 0 {
		t.Errorf("Expected union of two empty sets to be empty")
	}
}

func TestTransformSet(t *testing.T) {
	// Test case 1: Set with elements
	set1 := SetOf(1, 2, 3)
	expected := SetOf("1", "2", "3")
	setMap := TransformSet(set1, func(val int) string { return fmt.Sprintf("%d", val) })
	if len(setMap) != len(expected) {
		t.Errorf("Expected SetMap to have length %d, got %d", len(expected), len(setMap))
	}
	for k := range setMap {
		if _, exists := expected[k]; !exists {
			t.Errorf("Unexpected element in SetMap: %s", k)
		}
	}

	// Test case 2: Empty Set
	set2 := NewSet[int]()
	setMap = TransformSet(set2, func(val int) string { return fmt.Sprintf("%d", val) })
	if len(setMap) != 0 {
		t.Errorf("Expected SetMap of an empty set to be empty")
	}
}

func TestSetIterCount(t *testing.T) {
	// Test case 1: Counting elements in the set
	seq := SetOf(1, 2, 3)
	count := seq.Iter().Count()
	if count != 3 {
		t.Errorf("Expected count to be 3, got %d", count)
	}

	// Test case 2: Counting elements in an empty set
	emptySeq := NewSet[int]()
	emptyCount := emptySeq.Iter().Count()
	if emptyCount != 0 {
		t.Errorf("Expected count to be 0 for an empty set, got %d", emptyCount)
	}
}

func TestSetIterRange(t *testing.T) {
	// Test case 1: Stop iteration when function returns false
	seq := SetOf(1, 2, 3, 4)
	var result []int
	seq.Iter().Range(func(v int) bool {
		if v == 3 {
			result = append(result, v)
			return false
		}
		return true
	})

	expected := []int{3}
	if len(result) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(result))
	}

	// Test case 2: Iterate over the entire set
	emptySeq := NewSet[string]()
	emptyResult := make([]string, 0)

	emptySeq.Iter().Range(func(v string) bool {
		emptyResult = append(emptyResult, v)
		return true
	})

	if len(emptyResult) != 0 {
		t.Errorf("Expected no elements in an empty set, got %d elements", len(emptyResult))
	}
}

func TestSetIterFilter(t *testing.T) {
	// Test case 1: Filter even numbers
	seq := SetOf(1, 2, 3, 4, 5)

	even := seq.Iter().Filter(func(v int) bool {
		return v%2 == 0
	}).Collect()

	expected := SetOf(2, 4)
	if len(even) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(even))
	}
	for k := range even {
		if _, ok := expected[k]; !ok {
			t.Errorf("Unexpected element %v in the result", k)
		}
	}

	// Test case 2: Filter odd numbers
	odd := seq.Iter().Filter(func(v int) bool {
		return v%2 != 0
	}).Collect()

	oddExpected := SetOf(1, 3, 5)
	if len(odd) != len(oddExpected) {
		t.Errorf("Expected %d elements, got %d", len(oddExpected), len(odd))
	}
	for k := range odd {
		if _, ok := oddExpected[k]; !ok {
			t.Errorf("Unexpected element %v in the result", k)
		}
	}

	// Test case 3: Filter all elements
	all := seq.Iter().Filter(func(v int) bool {
		return true
	}).Collect()

	if len(all) != len(seq) {
		t.Errorf("Expected %d elements, got %d", len(seq), len(all))
	}
	for k := range all {
		if _, ok := seq[k]; !ok {
			t.Errorf("Unexpected element %v in the result", k)
		}
	}
}

func TestSetIterExclude(t *testing.T) {
	// Test case 1: Exclude even numbers
	seq := SetOf(1, 2, 3, 4, 5)
	notEven := seq.Iter().Exclude(func(v int) bool {
		return v%2 == 0
	}).Collect()

	expected := SetOf(1, 3, 5)
	if len(notEven) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(notEven))
	}
	for k := range notEven {
		if _, ok := expected[k]; !ok {
			t.Errorf("Unexpected element %v in the result", k)
		}
	}

	// Test case 2: Exclude odd numbers
	notOdd := seq.Iter().Exclude(func(v int) bool {
		return v%2 != 0
	}).Collect()

	notOddExpected := SetOf(2, 4)
	if len(notOdd) != len(notOddExpected) {
		t.Errorf("Expected %d elements, got %d", len(notOddExpected), len(notOdd))
	}
	for k := range notOdd {
		if _, ok := notOddExpected[k]; !ok {
			t.Errorf("Unexpected element %v in the result", k)
		}
	}

	// Test case 3: Exclude all elements
	empty := seq.Iter().Exclude(func(v int) bool {
		return true
	}).Collect()

	if len(empty) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(empty))
	}
}

func TestSetIterMap(t *testing.T) {
	// Test case 1: Double each element
	set := SetOf(1, 2, 3)
	doubled := set.Iter().Map(func(val int) int {
		return val * 2
	}).Collect()

	expected := SetOf(2, 4, 6)
	if !reflect.DeepEqual(doubled, expected) {
		t.Errorf("Expected set after doubling elements to be %v, got %v", expected, doubled)
	}

	// Test case 2: Square each element
	set2 := SetOf(1, 2, 3)
	squared := set2.Iter().Map(func(val int) int {
		return val * val
	}).Collect()

	expected2 := SetOf(1, 4, 9)
	if !reflect.DeepEqual(squared, expected2) {
		t.Errorf("Expected set after squaring elements to be %v, got %v", expected2, squared)
	}
}

func TestSetIterInspect(t *testing.T) {
	// Define a set to iterate over
	s := SetOf(1, 2, 3)

	// Define a slice to store the inspected elements
	inspectedElements := NewSet[int]()

	// Create a new iterator with Inspect and collect the elements
	s.Iter().Inspect(func(v int) {
		inspectedElements.Insert(v)
	}).Collect()

	if !inspectedElements.Eq(s) {
		t.Errorf("Expected inspected elements to be equal to the set, got %v", inspectedElements)
	}
}

func TestSetTransform(t *testing.T) {
	original := Set[int]{1: {}, 2: {}, 3: {}}

	addElement := func(s Set[int]) Set[int] {
		s[4] = struct{}{}
		return s
	}

	expected := Set[int]{1: {}, 2: {}, 3: {}, 4: {}}

	result := original.Transform(addElement)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Transform failed: expected %v, got %v", expected, result)
	}

	removeElement := func(s Set[int]) Set[int] {
		delete(s, 2)
		return s
	}

	expectedAfterRemoval := Set[int]{1: {}, 3: {}, 4: {}}
	resultAfterRemoval := result.Transform(removeElement)

	if !reflect.DeepEqual(resultAfterRemoval, expectedAfterRemoval) {
		t.Errorf("Transform with removal failed: expected %v, got %v", expectedAfterRemoval, resultAfterRemoval)
	}
}

func TestSetIterFind(t *testing.T) {
	// Test case 1: Element found
	seq := Set[int]{1: {}, 2: {}, 3: {}, 4: {}, 5: {}}
	found := seq.Iter().Find(func(i int) bool {
		return i == 2
	})
	if !found.IsSome() {
		t.Error("Expected found option to be Some")
	}
	if found.Some() != 2 {
		t.Errorf("Expected found element to be 2, got %d", found.Some())
	}

	// Test case 2: Element not found
	notFound := seq.Iter().Find(func(i int) bool {
		return i == 6
	})
	if notFound.IsSome() {
		t.Error("Expected not found option to be None")
	}
}
