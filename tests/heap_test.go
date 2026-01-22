package g

import (
	"strings"
	"testing"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestNewHeap(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])

	if heap.Len() != 0 {
		t.Errorf("Expected new heap to be empty, got length %d", heap.Len())
	}

	if !heap.IsEmpty() {
		t.Error("Expected new heap to be empty")
	}
}

func TestHeap_Push_Pop(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])

	// Test pushing single elements
	heap.Push(5)
	if heap.Len() != 1 {
		t.Errorf("Expected length 1, got %d", heap.Len())
	}

	// Test pushing multiple elements
	heap.Push(3, 7, 1, 9)
	if heap.Len() != 5 {
		t.Errorf("Expected length 5, got %d", heap.Len())
	}

	// Test min heap property - should pop in ascending order
	expected := []int{1, 3, 5, 7, 9}
	for _, exp := range expected {
		val := heap.Pop()
		if !val.IsSome() {
			t.Errorf("Expected to pop value %d, got None", exp)
			continue
		}
		if val.Some() != exp {
			t.Errorf("Expected to pop %d, got %d", exp, val.Some())
		}
	}

	// Test popping from empty heap
	empty := heap.Pop()
	if empty.IsSome() {
		t.Error("Expected None when popping from empty heap")
	}
}

func TestHeap_Peek(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])

	// Test peek on empty heap
	val := heap.Peek()
	if val.IsSome() {
		t.Error("Expected None when peeking empty heap")
	}

	// Test peek with elements
	heap.Push(5, 3, 7, 1)
	val = heap.Peek()
	if !val.IsSome() || val.Some() != 1 {
		t.Errorf("Expected to peek 1, got %v", val)
	}

	// Ensure peek doesn't modify heap
	if heap.Len() != 4 {
		t.Errorf("Expected length to remain 4 after peek, got %d", heap.Len())
	}
}

func TestHeap_MaxHeap(t *testing.T) {
	// Create max heap by reversing comparison
	heap := g.NewHeap(func(a, b int) cmp.Ordering {
		return cmp.Cmp(b, a) // reverse comparison for max heap
	})

	heap.Push(5, 3, 7, 1, 9)

	// Should pop in descending order (max heap)
	expected := []int{9, 7, 5, 3, 1}
	for _, exp := range expected {
		val := heap.Pop()
		if !val.IsSome() {
			t.Errorf("Expected to pop value %d, got None", exp)
			continue
		}
		if val.Some() != exp {
			t.Errorf("Expected to pop %d, got %d", exp, val.Some())
		}
	}
}

func TestHeap_StringComparison(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[string])

	heap.Push("banana", "apple", "cherry", "date")

	// Should pop in alphabetical order
	expected := []string{"apple", "banana", "cherry", "date"}
	for _, exp := range expected {
		val := heap.Pop()
		if !val.IsSome() {
			t.Errorf("Expected to pop value %s, got None", exp)
			continue
		}
		if val.Some() != exp {
			t.Errorf("Expected to pop %s, got %s", exp, val.Some())
		}
	}
}

func TestHeap_Len_Empty(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])

	if heap.Len() != 0 {
		t.Errorf("Expected length 0, got %d", heap.Len())
	}

	if !heap.IsEmpty() {
		t.Error("Expected heap to be empty")
	}

	heap.Push(1, 2, 3)

	if heap.Len() != 3 {
		t.Errorf("Expected length 3, got %d", heap.Len())
	}

	if heap.IsEmpty() {
		t.Error("Expected heap to not be empty")
	}

	heap.Pop()
	heap.Pop()
	heap.Pop()

	if heap.Len() != 0 {
		t.Errorf("Expected length 0 after popping all, got %d", heap.Len())
	}

	if !heap.IsEmpty() {
		t.Error("Expected heap to be empty after popping all")
	}
}

func TestHeap_ToSlice(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(5, 3, 7, 1)

	slice := heap.ToSlice()

	if len(slice) != 4 {
		t.Errorf("Expected slice length 4, got %d", len(slice))
	}

	// ToSlice should contain all elements (order not guaranteed)
	elements := make(map[int]bool)
	for _, v := range slice {
		elements[v] = true
	}

	expected := map[int]bool{1: true, 3: true, 5: true, 7: true}
	for k := range expected {
		if !elements[k] {
			t.Errorf("Expected element %d in slice", k)
		}
	}
}

func TestHeap_Clear(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3, 4, 5)

	if heap.Len() != 5 {
		t.Errorf("Expected length 5 before clear, got %d", heap.Len())
	}

	heap.Clear()

	if heap.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", heap.Len())
	}

	if !heap.IsEmpty() {
		t.Error("Expected heap to be empty after clear")
	}

	val := heap.Pop()
	if val.IsSome() {
		t.Error("Expected None when popping from cleared heap")
	}
}

func TestHeap_Clone(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(3, 1, 4, 1, 5)

	cloned := heap.Clone()

	if cloned.Len() != heap.Len() {
		t.Errorf("Expected cloned heap length %d, got %d", heap.Len(), cloned.Len())
	}

	// Verify both heaps have same elements by popping
	originalElements := make([]int, 0)
	clonedElements := make([]int, 0)

	for !heap.IsEmpty() {
		originalElements = append(originalElements, heap.Pop().Some())
	}

	for !cloned.IsEmpty() {
		clonedElements = append(clonedElements, cloned.Pop().Some())
	}

	if len(originalElements) != len(clonedElements) {
		t.Errorf("Expected same number of elements, got original=%d, cloned=%d",
			len(originalElements), len(clonedElements))
	}

	for i, v := range originalElements {
		if i >= len(clonedElements) || clonedElements[i] != v {
			t.Errorf("Element at index %d: original=%d, cloned=%d", i, v,
				func() int {
					if i < len(clonedElements) {
						return clonedElements[i]
					}
					return -1
				}())
		}
	}
}

func TestHeap_Transform(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1, 2, 3)

	result := heap.Transform(func(h *g.Heap[int]) *g.Heap[int] {
		h.Push(4, 5)
		return h
	})

	if result.Len() != 5 {
		t.Errorf("Expected transformed heap length 5, got %d", result.Len())
	}

	// Verify the transformation was applied to the original heap
	if heap.Len() != 5 {
		t.Errorf("Expected original heap length 5, got %d", heap.Len())
	}
}

func TestHeap_DuplicateElements(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(2, 1, 2, 3, 1, 3, 2)

	expected := []int{1, 1, 2, 2, 2, 3, 3}
	actual := make([]int, 0)

	for !heap.IsEmpty() {
		actual = append(actual, heap.Pop().Some())
	}

	if len(actual) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(actual))
	}

	for i, v := range expected {
		if i >= len(actual) || actual[i] != v {
			t.Errorf("At index %d: expected %d, got %d", i, v,
				func() int {
					if i < len(actual) {
						return actual[i]
					}
					return -1
				}())
		}
	}
}

func TestHeap_SingleElement(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(42)

	if heap.Len() != 1 {
		t.Errorf("Expected length 1, got %d", heap.Len())
	}

	peek := heap.Peek()
	if !peek.IsSome() || peek.Some() != 42 {
		t.Errorf("Expected peek to return 42, got %v", peek)
	}

	pop := heap.Pop()
	if !pop.IsSome() || pop.Some() != 42 {
		t.Errorf("Expected pop to return 42, got %v", pop)
	}

	if !heap.IsEmpty() {
		t.Error("Expected heap to be empty after popping single element")
	}
}

func TestHeap_MaintainsProperty(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])

	// Push elements in random order
	elements := []int{50, 25, 75, 10, 30, 60, 80, 5, 15, 35}
	for _, elem := range elements {
		heap.Push(elem)
	}

	// Pop all elements and verify they come out sorted
	var previous int = -1
	for !heap.IsEmpty() {
		current := heap.Pop().Some()
		if previous >= 0 && current < previous {
			t.Errorf("Heap property violated: %d came after %d", current, previous)
		}
		previous = current
	}
}

func TestHeap_CustomComparison(t *testing.T) {
	// Heap that sorts by absolute value
	heap := g.NewHeap(func(a, b int) cmp.Ordering {
		absA := a
		if a < 0 {
			absA = -a
		}
		absB := b
		if b < 0 {
			absB = -b
		}
		return cmp.Cmp(absA, absB)
	})

	heap.Push(-5, 3, -1, 4, -2)

	// Should pop in order of increasing absolute value
	expected := []int{-1, -2, 3, 4, -5} // abs values: 1, 2, 3, 4, 5
	for _, exp := range expected {
		val := heap.Pop()
		if !val.IsSome() {
			t.Errorf("Expected to pop value %d, got None", exp)
			continue
		}
		if val.Some() != exp {
			t.Errorf("Expected to pop %d, got %d", exp, val.Some())
		}
	}
}

func TestHeap_IntoIter(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(3, 1, 4, 1, 5)

	// Test that IntoIter drains the heap in sorted order
	var result []int
	heap.IntoIter()(func(v int) bool {
		result = append(result, v)
		return true
	})

	// Should be in ascending order (min heap)
	expected := []int{1, 1, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("IntoIter: expected length %d, got %d", len(expected), len(result))
	}

	for i, exp := range expected {
		if i < len(result) && result[i] != exp {
			t.Errorf("IntoIter at index %d: expected %d, got %d", i, exp, result[i])
		}
	}

	// Heap should now be empty
	if !heap.IsEmpty() {
		t.Errorf("Expected heap to be empty after IntoIter, got length %d", heap.Len())
	}
}

func TestHeap_IntoIterEarlyTermination(t *testing.T) {
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(5, 3, 7, 1, 9)

	// Test early termination of IntoIter
	var result []int
	heap.IntoIter()(func(v int) bool {
		result = append(result, v)
		return v < 5 // Stop when we reach 5 or higher
	})

	// Should stop at 5
	expected := []int{1, 3, 5}
	if len(result) != len(expected) {
		t.Errorf("IntoIter early termination: expected length %d, got %d", len(expected), len(result))
	}

	// Heap should still contain remaining elements
	if heap.IsEmpty() {
		t.Errorf("Expected heap to still have elements after early termination")
	}
}

// Test to trigger the heapify method indirectly
func TestHeap_HeapifyIndirect(t *testing.T) {
	// The heapify method is private, but we can test it indirectly
	// by creating scenarios that would require it to be called

	// Create a heap and manipulate it in ways that would require heapification
	heap := g.NewHeap(cmp.Cmp[int])

	// Push elements in a way that might require heapification
	elements := []int{10, 5, 15, 3, 7, 12, 20, 1, 4, 6}
	for _, elem := range elements {
		heap.Push(elem)
	}

	// Verify heap property is maintained by popping all elements
	var result []int
	for !heap.IsEmpty() {
		val := heap.Pop()
		if val.IsSome() {
			result = append(result, val.Some())
		}
	}

	// Should be in ascending order
	for i := 1; i < len(result); i++ {
		if result[i] < result[i-1] {
			t.Errorf("Heap property violated: %d came before %d", result[i-1], result[i])
		}
	}
}

func TestHeap_String(t *testing.T) {
	// Test empty heap
	heap := g.NewHeap(cmp.Cmp[int])
	str := heap.String()
	if str != "Heap[]" {
		t.Errorf("Expected 'Heap[]' for empty heap, got '%s'", str)
	}

	// Test heap with elements
	heap.Push(3, 1, 4, 2)
	str = heap.String()

	// Should start with "Heap[" and end with "]"
	if !strings.HasPrefix(str, "Heap[") {
		t.Errorf("Expected string to start with 'Heap[', got '%s'", str)
	}
	if !strings.HasSuffix(str, "]") {
		t.Errorf("Expected string to end with ']', got '%s'", str)
	}

	// Should contain all elements (order not guaranteed due to heap property)
	expectedElements := []string{"1", "2", "3", "4"}
	for _, elem := range expectedElements {
		if !strings.Contains(str, elem) {
			t.Errorf("Expected string to contain '%s', got '%s'", elem, str)
		}
	}

	// Test single element heap
	singleHeap := g.NewHeap(cmp.Cmp[string])
	singleHeap.Push("hello")
	singleStr := singleHeap.String()
	expected := "Heap[hello]"
	if singleStr != expected {
		t.Errorf("Expected '%s' for single element heap, got '%s'", expected, singleStr)
	}
}

func TestHeapHeapify(t *testing.T) {
	// Create a heap and add elements to trigger heapify
	heap := g.NewHeap(cmp.Cmp[int])

	// Add elements in non-heap order
	elements := []int{5, 3, 8, 1, 9, 2}
	for _, elem := range elements {
		heap.Push(elem)
	}

	// Verify heap property is maintained
	sorted := make([]int, 0, len(elements))
	for !heap.IsEmpty() {
		val := heap.Pop()
		if val.IsSome() {
			sorted = append(sorted, val.Unwrap())
		}
	}

	// Should be sorted in ascending order (min-heap)
	for i := 1; i < len(sorted); i++ {
		if sorted[i-1] > sorted[i] {
			t.Errorf("Heap property violated: elements not in order: %v", sorted)
			break
		}
	}
}

func TestHeapPrint(t *testing.T) {
	// Test Print method - should return the heap unchanged
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)

	result := heap.Print()

	// Should return the same heap instance
	if result != heap {
		t.Errorf("Print() should return the same heap instance")
	}

	// Heap should be unchanged
	if heap.Len() != 3 {
		t.Errorf("Print() should not modify heap, expected length 3, got %d", heap.Len())
	}
}

func TestHeapPrintln(t *testing.T) {
	// Test Println method - should return the heap unchanged
	heap := g.NewHeap(cmp.Cmp[int])
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)

	result := heap.Println()

	// Should return the same heap instance
	if result != heap {
		t.Errorf("Println() should return the same heap instance")
	}

	// Heap should be unchanged
	if heap.Len() != 3 {
		t.Errorf("Println() should not modify heap, expected length 3, got %d", heap.Len())
	}
}
