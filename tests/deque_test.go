package g_test

import (
	"testing"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestNewDeque(t *testing.T) {
	// Test creating empty deque
	dq := g.NewDeque[int]()
	if !dq.IsEmpty() {
		t.Errorf("Expected empty deque, got length %d", dq.Len())
	}
	if dq.Len() != 0 {
		t.Errorf("Expected length 0, got %d", dq.Len())
	}

	// Test creating deque with capacity
	dq2 := g.NewDeque[int](10)
	if !dq2.IsEmpty() {
		t.Errorf("Expected empty deque, got length %d", dq2.Len())
	}
	if dq2.Capacity() != 10 {
		t.Errorf("Expected capacity 10, got %d", dq2.Capacity())
	}
}

func TestDequeOf(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)
	if dq.Len() != 5 {
		t.Errorf("Expected length 5, got %d", dq.Len())
	}

	// Check order is maintained
	for i := 0; i < 5; i++ {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != i+1 {
			t.Errorf("Expected element %d at index %d, got %v", i+1, i, val)
		}
	}
}

func TestDequePushPopFront(t *testing.T) {
	dq := g.NewDeque[int]()

	// Test pushing to front
	dq.PushFront(1)
	dq.PushFront(2)
	dq.PushFront(3)

	if dq.Len() != 3 {
		t.Errorf("Expected length 3, got %d", dq.Len())
	}

	// Test front element
	front := dq.Front()
	if !front.IsSome() || front.Unwrap() != 3 {
		t.Errorf("Expected front element 3, got %v", front)
	}

	// Test popping from front
	val1 := dq.PopFront()
	if !val1.IsSome() || val1.Unwrap() != 3 {
		t.Errorf("Expected popped value 3, got %v", val1)
	}

	val2 := dq.PopFront()
	if !val2.IsSome() || val2.Unwrap() != 2 {
		t.Errorf("Expected popped value 2, got %v", val2)
	}

	val3 := dq.PopFront()
	if !val3.IsSome() || val3.Unwrap() != 1 {
		t.Errorf("Expected popped value 1, got %v", val3)
	}

	// Test popping from empty deque
	val4 := dq.PopFront()
	if val4.IsSome() {
		t.Errorf("Expected None from empty deque, got %v", val4)
	}
}

func TestDequePushPopBack(t *testing.T) {
	dq := g.NewDeque[int]()

	// Test pushing to back
	dq.PushBack(1)
	dq.PushBack(2)
	dq.PushBack(3)

	if dq.Len() != 3 {
		t.Errorf("Expected length 3, got %d", dq.Len())
	}

	// Test back element
	back := dq.Back()
	if !back.IsSome() || back.Unwrap() != 3 {
		t.Errorf("Expected back element 3, got %v", back)
	}

	// Test popping from back
	val1 := dq.PopBack()
	if !val1.IsSome() || val1.Unwrap() != 3 {
		t.Errorf("Expected popped value 3, got %v", val1)
	}

	val2 := dq.PopBack()
	if !val2.IsSome() || val2.Unwrap() != 2 {
		t.Errorf("Expected popped value 2, got %v", val2)
	}

	val3 := dq.PopBack()
	if !val3.IsSome() || val3.Unwrap() != 1 {
		t.Errorf("Expected popped value 1, got %v", val3)
	}

	// Test popping from empty deque
	val4 := dq.PopBack()
	if val4.IsSome() {
		t.Errorf("Expected None from empty deque, got %v", val4)
	}
}

func TestDequeMixedOperations(t *testing.T) {
	dq := g.NewDeque[int]()

	// Mix front and back operations
	dq.PushBack(1)  // [1]
	dq.PushFront(2) // [2, 1]
	dq.PushBack(3)  // [2, 1, 3]
	dq.PushFront(4) // [4, 2, 1, 3]

	if dq.Len() != 4 {
		t.Errorf("Expected length 4, got %d", dq.Len())
	}

	// Check elements
	expected := []int{4, 2, 1, 3}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Pop mixed
	front := dq.PopFront() // Should get 4, leaving [2, 1, 3]
	if !front.IsSome() || front.Unwrap() != 4 {
		t.Errorf("Expected front 4, got %v", front)
	}

	back := dq.PopBack() // Should get 3, leaving [2, 1]
	if !back.IsSome() || back.Unwrap() != 3 {
		t.Errorf("Expected back 3, got %v", back)
	}

	if dq.Len() != 2 {
		t.Errorf("Expected length 2, got %d", dq.Len())
	}
}

func TestDequeGetSet(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)

	// Test valid gets
	for i := 0; i < 5; i++ {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != i+1 {
			t.Errorf("Expected element %d at index %d, got %v", i+1, i, val)
		}
	}

	// Test invalid gets
	val := dq.Get(-1)
	if val.IsSome() {
		t.Errorf("Expected None for negative index, got %v", val)
	}

	val = dq.Get(5)
	if val.IsSome() {
		t.Errorf("Expected None for out of bounds index, got %v", val)
	}

	// Test set
	success := dq.Set(2, 10)
	if !success {
		t.Errorf("Expected successful set operation")
	}

	val = dq.Get(2)
	if !val.IsSome() || val.Unwrap() != 10 {
		t.Errorf("Expected element 10 at index 2 after set, got %v", val)
	}

	// Test invalid set
	success = dq.Set(-1, 20)
	if success {
		t.Errorf("Expected failed set operation for negative index")
	}

	success = dq.Set(5, 20)
	if success {
		t.Errorf("Expected failed set operation for out of bounds index")
	}
}

func TestDequeInsertRemove(t *testing.T) {
	dq := g.DequeOf(1, 3, 5)

	// Insert at beginning
	dq.Insert(0, 0) // [0, 1, 3, 5]
	expected := []int{0, 1, 3, 5}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After insert at 0: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Insert in middle
	dq.Insert(2, 2) // [0, 1, 2, 3, 5]
	expected = []int{0, 1, 2, 3, 5}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After insert at 2: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Insert at end
	dq.Insert(5, 6) // [0, 1, 2, 3, 5, 6]
	expected = []int{0, 1, 2, 3, 5, 6}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After insert at end: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Remove from middle
	removed := dq.Remove(3) // Remove value 3, leaving [0, 1, 2, 5, 6]
	if !removed.IsSome() || removed.Unwrap() != 3 {
		t.Errorf("Expected removed value 3, got %v", removed)
	}

	expected = []int{0, 1, 2, 5, 6}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After remove: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Test remove out of bounds
	removed = dq.Remove(-1)
	if removed.IsSome() {
		t.Errorf("Expected None for remove at negative index, got %v", removed)
	}

	removed = dq.Remove(10)
	if removed.IsSome() {
		t.Errorf("Expected None for remove at out of bounds index, got %v", removed)
	}
}

func TestDequeClear(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)
	if dq.Len() != 5 {
		t.Errorf("Expected length 5 before clear, got %d", dq.Len())
	}

	dq.Clear()
	if !dq.IsEmpty() {
		t.Errorf("Expected empty deque after clear, got length %d", dq.Len())
	}

	// Test that we can still use it after clear
	dq.PushBack(10)
	if dq.Len() != 1 {
		t.Errorf("Expected length 1 after push to cleared deque, got %d", dq.Len())
	}
}

func TestDequeSwap(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)

	dq.Swap(0, 4) // Swap first and last elements

	expected := []int{5, 2, 3, 4, 1}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After swap: Expected element %d at index %d, got %v", exp, i, val)
		}
	}
}

func TestDequeRotate(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)

	// Test rotate left by 2
	dq.RotateLeft(2) // [3, 4, 5, 1, 2]

	expected := []int{3, 4, 5, 1, 2}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After rotate left: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Test rotate right by 1
	dq2 := g.DequeOf(1, 2, 3, 4, 5)
	dq2.RotateRight(1) // [5, 1, 2, 3, 4]

	expected = []int{5, 1, 2, 3, 4}
	for i, exp := range expected {
		val := dq2.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After rotate right: Expected element %d at index %d, got %v", exp, i, val)
		}
	}
}

func TestDequeContains(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)

	// Test existing elements
	for i := 1; i <= 5; i++ {
		if !dq.Contains(i) {
			t.Errorf("Expected deque to contain %d", i)
		}
	}

	// Test non-existing element
	if dq.Contains(10) {
		t.Errorf("Expected deque not to contain 10")
	}
}

func TestDequeIndex(t *testing.T) {
	dq := g.DequeOf(10, 20, 30, 20, 40)

	// Test finding existing elements
	idx := dq.Index(10)
	if idx != 0 {
		t.Errorf("Expected index 0 for value 10, got %d", idx)
	}

	idx = dq.Index(20) // Should find first occurrence
	if idx != 1 {
		t.Errorf("Expected index 1 for first occurrence of 20, got %d", idx)
	}

	idx = dq.Index(40)
	if idx != 4 {
		t.Errorf("Expected index 4 for value 40, got %d", idx)
	}

	// Test non-existing element
	idx = dq.Index(100)
	if idx != -1 {
		t.Errorf("Expected index -1 for non-existing value, got %d", idx)
	}
}

func TestDequeClone(t *testing.T) {
	original := g.DequeOf(1, 2, 3, 4, 5)
	clone := original.Clone()

	// Check they have same contents
	if clone.Len() != original.Len() {
		t.Errorf("Expected clone length %d, got %d", original.Len(), clone.Len())
	}

	for i := g.Int(0); i < original.Len(); i++ {
		origVal := original.Get(i)
		cloneVal := clone.Get(i)
		if origVal != cloneVal {
			t.Errorf("Expected same value at index %d, got original: %v, clone: %v", i, origVal, cloneVal)
		}
	}

	// Modify clone and ensure original is unchanged
	clone.PushBack(10)
	if original.Len() == clone.Len() {
		t.Errorf("Expected clone modification not to affect original")
	}
}

func TestDequeEq(t *testing.T) {
	dq1 := g.DequeOf(1, 2, 3, 4, 5)
	dq2 := g.DequeOf(1, 2, 3, 4, 5)
	dq3 := g.DequeOf(1, 2, 3, 4, 6)
	dq4 := g.DequeOf(1, 2, 3, 4)

	// Test equal deques
	if !dq1.Eq(dq2) {
		t.Errorf("Expected dq1 to equal dq2")
	}

	// Test different values
	if dq1.Eq(dq3) {
		t.Errorf("Expected dq1 not to equal dq3")
	}

	// Test different lengths
	if dq1.Eq(dq4) {
		t.Errorf("Expected dq1 not to equal dq4")
	}
}

func TestDequeRetain(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// Retain only even numbers
	dq.Retain(func(x int) bool {
		return x%2 == 0
	})

	expected := []int{2, 4, 6, 8, 10}
	if dq.Len() != g.Int(len(expected)) {
		t.Errorf("Expected length %d after retain, got %d", len(expected), dq.Len())
	}

	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After retain: Expected element %d at index %d, got %v", exp, i, val)
		}
	}
}

func TestDequeToSlice(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)
	slice := dq.ToSlice()

	if len(slice) != int(dq.Len()) {
		t.Errorf("Expected slice length %d, got %d", dq.Len(), len(slice))
	}

	for i := 0; i < len(slice); i++ {
		if slice[i] != i+1 {
			t.Errorf("Expected slice[%d] = %d, got %d", i, i+1, slice[i])
		}
	}
}

func TestDequeGrowth(t *testing.T) {
	dq := g.NewDeque[int]()

	// Test that deque grows automatically
	for i := 0; i < 100; i++ {
		dq.PushBack(i)
	}

	if dq.Len() != 100 {
		t.Errorf("Expected length 100 after adding 100 elements, got %d", dq.Len())
	}

	// Verify all elements are still there
	for i := 0; i < 100; i++ {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != i {
			t.Errorf("Expected element %d at index %d, got %v", i, i, val)
		}
	}
}

func TestDequeReserve(t *testing.T) {
	dq := g.NewDeque[int]()
	initialCap := dq.Capacity()

	dq.Reserve(100)
	newCap := dq.Capacity()

	if newCap < initialCap+100 {
		t.Errorf("Expected capacity to increase by at least 100, got increase of %d", newCap-initialCap)
	}

	// Test that it's still empty
	if !dq.IsEmpty() {
		t.Errorf("Expected deque to remain empty after reserve")
	}
}

func TestDequeShrinkToFit(t *testing.T) {
	dq := g.NewDeque[int](100)

	// Add only a few elements
	for i := 0; i < 5; i++ {
		dq.PushBack(i)
	}

	initialCap := dq.Capacity()
	dq.ShrinkToFit()
	newCap := dq.Capacity()

	if newCap >= initialCap {
		t.Errorf("Expected capacity to shrink from %d, but got %d", initialCap, newCap)
	}

	if newCap != dq.Len() {
		t.Errorf("Expected capacity %d to equal length %d after shrink to fit", newCap, dq.Len())
	}

	// Verify elements are still there
	for i := 0; i < 5; i++ {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != i {
			t.Errorf("Expected element %d at index %d after shrink, got %v", i, i, val)
		}
	}
}

func TestDequeIter(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)

	var collected []int
	dq.Iter()(func(val int) bool {
		collected = append(collected, val)
		return true
	})

	expected := []int{1, 2, 3, 4, 5}
	if len(collected) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(collected))
	}

	for i, exp := range expected {
		if collected[i] != exp {
			t.Errorf("Expected element %d at position %d, got %d", exp, i, collected[i])
		}
	}
}

func TestDequeIterReverse(t *testing.T) {
	dq := g.DequeOf(1, 2, 3, 4, 5)

	var collected []int
	dq.IterReverse()(func(val int) bool {
		collected = append(collected, val)
		return true
	})

	expected := []int{5, 4, 3, 2, 1}
	if len(collected) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(collected))
	}

	for i, exp := range expected {
		if collected[i] != exp {
			t.Errorf("Expected element %d at position %d, got %d", exp, i, collected[i])
		}
	}
}

func TestDequeString(t *testing.T) {
	// Test empty deque
	dq := g.NewDeque[int]()
	str := dq.String()
	if str != "Deque[]" {
		t.Errorf("Expected 'Deque[]' for empty deque, got '%s'", str)
	}

	// Test non-empty deque
	dq = g.DequeOf(1, 2, 3)
	str = dq.String()
	expected := "Deque[1, 2, 3]"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestDequeRingBufferBehavior(t *testing.T) {
	// Test that the ring buffer works correctly when wrapping around
	dq := g.NewDeque[int](4) // Small capacity to force wrapping

	// Fill the deque
	for i := 0; i < 4; i++ {
		dq.PushBack(i)
	}

	// Remove from front and add to back (should wrap around)
	dq.PopFront()  // Remove 0
	dq.PushBack(4) // Add 4

	// Check the order is maintained: [1, 2, 3, 4]
	expected := []int{1, 2, 3, 4}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("Ring buffer test: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Test front operations with wrapping
	dq.PopBack()    // Remove 4: [1, 2, 3]
	dq.PushFront(0) // Add 0 to front: [0, 1, 2, 3]

	expected = []int{0, 1, 2, 3}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("Ring buffer front test: Expected element %d at index %d, got %v", exp, i, val)
		}
	}
}

// Additional tests for missing coverage

func TestDequeFrontBackEmpty(t *testing.T) {
	dq := g.NewDeque[int]()

	// Test Front() on empty deque
	front := dq.Front()
	if front.IsSome() {
		t.Errorf("Expected None from empty deque Front(), got %v", front)
	}

	// Test Back() on empty deque
	back := dq.Back()
	if back.IsSome() {
		t.Errorf("Expected None from empty deque Back(), got %v", back)
	}
}

func TestDequeInsertEdgeCases(t *testing.T) {
	dq := g.NewDeque[int]()

	// Insert into empty deque
	dq.Insert(0, 42)
	if dq.Len() != 1 {
		t.Errorf("Expected length 1 after insert into empty deque, got %d", dq.Len())
	}

	val := dq.Get(0)
	if !val.IsSome() || val.Unwrap() != 42 {
		t.Errorf("Expected 42 at index 0, got %v", val)
	}

	// Test insert at valid end position
	dq2 := g.DequeOf(1, 2, 3)

	// Insert at end (valid position)
	dq2.Insert(dq2.Len(), 999) // Insert at length (valid)
	if dq2.Get(dq2.Len()-1).Some() != 999 {
		t.Errorf("Expected insert at end to append")
	}
}

func TestDequeInsertPanicsOnNegativeIndex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Insert with negative index to panic")
		}
	}()

	dq := g.DequeOf(1, 2, 3)
	dq.Insert(-1, 999) // Should panic
}

func TestDequeInsertPanicsOnHighIndex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Insert with out-of-bounds high index to panic")
		}
	}()

	dq := g.DequeOf(1, 2, 3)
	dq.Insert(100, 999) // Should panic
}

func TestDequeRemoveEdgeCases(t *testing.T) {
	// Test remove from single element deque
	dq := g.DequeOf(42)
	removed := dq.Remove(0)
	if !removed.IsSome() || removed.Unwrap() != 42 {
		t.Errorf("Expected to remove 42 from single element deque, got %v", removed)
	}
	if !dq.IsEmpty() {
		t.Errorf("Expected deque to be empty after removing only element")
	}

	// Test remove from empty deque
	empty := g.NewDeque[int]()
	removed = empty.Remove(0)
	if removed.IsSome() {
		t.Errorf("Expected None when removing from empty deque, got %v", removed)
	}
}

func TestDequeSwapEdgeCases(t *testing.T) {
	// Test swap with same index
	dq := g.DequeOf(1, 2, 3)
	dq.Swap(1, 1)

	// Should remain unchanged
	expected := []int{1, 2, 3}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After self-swap: Expected element %d at index %d, got %v", exp, i, val)
		}
	}
}

func TestDequeSwapPanicsOnOutOfBounds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Swap with out-of-bounds index to panic")
		}
	}()

	dq := g.DequeOf(1, 2, 3)
	dq.Swap(-1, 0) // Should panic
}

func TestDequeRotateEdgeCases(t *testing.T) {
	// Test rotate empty deque
	empty := g.NewDeque[int]()
	empty.RotateLeft(5)  // Should not crash
	empty.RotateRight(3) // Should not crash

	// Test rotate single element
	single := g.DequeOf(42)
	single.RotateLeft(10)
	if single.Get(0).Some() != 42 {
		t.Errorf("Single element should remain after rotation")
	}
	single.RotateRight(10)
	if single.Get(0).Some() != 42 {
		t.Errorf("Single element should remain after rotation")
	}

	// Test rotate by 0
	dq := g.DequeOf(1, 2, 3)
	expected := []int{1, 2, 3}
	dq.RotateLeft(0)
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After rotate left 0: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	dq.RotateRight(0)
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After rotate right 0: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Test rotate by full length (should return to original)
	dq = g.DequeOf(1, 2, 3, 4, 5)
	dq.RotateLeft(5)
	expected = []int{1, 2, 3, 4, 5}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After full rotation: Expected element %d at index %d, got %v", exp, i, val)
		}
	}
}

func TestDequeMakeContiguous(t *testing.T) {
	// Create a deque that will have non-contiguous storage
	dq := g.NewDeque[int](4)

	// Fill it
	for i := 0; i < 4; i++ {
		dq.PushBack(i)
	}

	// Remove from front and add to back to create wrap-around
	dq.PopFront()
	dq.PopFront()
	dq.PushBack(4)
	dq.PushBack(5)

	// Now storage should be non-contiguous: [4, 5, 2, 3] in ring buffer
	// Call MakeContiguous
	dq.MakeContiguous()

	// Verify elements are still in correct order
	expected := []int{2, 3, 4, 5}
	for i, exp := range expected {
		val := dq.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After MakeContiguous: Expected element %d at index %d, got %v", exp, i, val)
		}
	}

	// Test MakeContiguous on empty deque
	empty := g.NewDeque[int]()
	empty.MakeContiguous() // Should not crash

	// Test MakeContiguous on already contiguous deque
	contiguous := g.DequeOf(1, 2, 3)
	contiguous.MakeContiguous()
	expected = []int{1, 2, 3}
	for i, exp := range expected {
		val := contiguous.Get(g.Int(i))
		if !val.IsSome() || val.Unwrap() != exp {
			t.Errorf("After MakeContiguous on contiguous: Expected element %d at index %d, got %v", exp, i, val)
		}
	}
}

func TestDequeIterReverseEdgeCases(t *testing.T) {
	// Test IterReverse on empty deque
	empty := g.NewDeque[int]()
	var collected []int
	empty.IterReverse()(func(val int) bool {
		collected = append(collected, val)
		return true
	})

	if len(collected) != 0 {
		t.Errorf("Expected no elements from empty IterReverse, got %d", len(collected))
	}

	// Test IterReverse early termination
	dq := g.DequeOf(1, 2, 3, 4, 5)
	collected = nil
	dq.IterReverse()(func(val int) bool {
		collected = append(collected, val)
		return val > 3 // Stop when we reach 3 or lower
	})

	expected := []int{5, 4, 3} // Should include the 3 that triggers the stop
	if len(collected) != len(expected) {
		t.Errorf("Expected %d elements with early termination, got %d: %v", len(expected), len(collected), collected)
	}
}

func TestDequeReserveEdgeCases(t *testing.T) {
	// Test reserve 0
	dq := g.NewDeque[int]()
	initialCap := dq.Capacity()
	dq.Reserve(0)
	if dq.Capacity() != initialCap {
		t.Errorf("Reserve(0) should not change capacity")
	}

	// Test reserve less than current capacity
	dq = g.NewDeque[int](10)
	dq.Reserve(5) // Should not shrink
	if dq.Capacity() < 10 {
		t.Errorf("Reserve should not shrink capacity")
	}
}

func TestDequeShrinkToFitEdgeCases(t *testing.T) {
	// Test ShrinkToFit on empty deque
	empty := g.NewDeque[int](10)
	empty.ShrinkToFit()
	// Should not crash and should handle empty case

	// Test ShrinkToFit when already optimal
	optimal := g.DequeOf(1, 2, 3)
	optimal.ShrinkToFit()
	// Should not crash
}

func TestDequeContainsComprehensive(t *testing.T) {
	t.Run("comparable_types", func(t *testing.T) {
		// Test with integers (comparable type)
		dq := g.DequeOf(1, 2, 3, 4, 5)

		// Test existing elements
		for i := 1; i <= 5; i++ {
			if !dq.Contains(i) {
				t.Errorf("Expected deque to contain %d", i)
			}
		}

		// Test non-existing element
		if dq.Contains(10) {
			t.Errorf("Expected deque not to contain 10")
		}

		// Test with zero value
		dqWithZero := g.DequeOf(0, 1, 2)
		if !dqWithZero.Contains(0) {
			t.Errorf("Expected deque to contain zero value")
		}

		// Test with strings (comparable type)
		stringDq := g.DequeOf("hello", "world", "test")
		if !stringDq.Contains("hello") {
			t.Errorf("Expected deque to contain 'hello'")
		}
		if !stringDq.Contains("world") {
			t.Errorf("Expected deque to contain 'world'")
		}
		if stringDq.Contains("missing") {
			t.Errorf("Expected deque not to contain 'missing'")
		}
	})

	t.Run("non_comparable_types", func(t *testing.T) {
		// Test with slices (non-comparable type)
		type TestSlice []int
		dq := g.NewDeque[TestSlice]()
		dq.PushBack(TestSlice{1, 2, 3})
		dq.PushBack(TestSlice{4, 5, 6})
		dq.PushBack(TestSlice{7, 8, 9})

		// Test existing slice
		if !dq.Contains(TestSlice{1, 2, 3}) {
			t.Errorf("Expected deque to contain {1, 2, 3}")
		}
		if !dq.Contains(TestSlice{4, 5, 6}) {
			t.Errorf("Expected deque to contain {4, 5, 6}")
		}
		if !dq.Contains(TestSlice{7, 8, 9}) {
			t.Errorf("Expected deque to contain {7, 8, 9}")
		}

		// Test non-existing slice
		if dq.Contains(TestSlice{1, 2, 4}) {
			t.Errorf("Expected deque not to contain {1, 2, 4}")
		}
		if dq.Contains(TestSlice{}) {
			t.Errorf("Expected deque not to contain empty slice")
		}

		// Test with maps (non-comparable type)
		type TestMap map[string]int
		mapDq := g.NewDeque[TestMap]()
		map1 := TestMap{"a": 1, "b": 2}
		map2 := TestMap{"c": 3, "d": 4}
		mapDq.PushBack(map1)
		mapDq.PushBack(map2)

		if !mapDq.Contains(map1) {
			t.Errorf("Expected deque to contain map1")
		}
		if !mapDq.Contains(map2) {
			t.Errorf("Expected deque to contain map2")
		}
		if mapDq.Contains(TestMap{"e": 5}) {
			t.Errorf("Expected deque not to contain different map")
		}
	})

	t.Run("edge_cases", func(t *testing.T) {
		// Test Contains on empty deque
		empty := g.NewDeque[int]()
		if empty.Contains(42) {
			t.Errorf("Empty deque should not contain any element")
		}

		// Test Contains with single element
		single := g.DequeOf(42)
		if !single.Contains(42) {
			t.Errorf("Single element deque should contain its element")
		}
		if single.Contains(43) {
			t.Errorf("Single element deque should not contain other elements")
		}

		// Test Contains with duplicate elements
		duplicates := g.DequeOf(1, 2, 2, 3, 2, 4)
		if !duplicates.Contains(2) {
			t.Errorf("Expected deque with duplicates to contain 2")
		}

		// Test Contains after modifications
		dq := g.DequeOf(1, 2, 3)
		if !dq.Contains(2) {
			t.Errorf("Expected deque to contain 2 before modification")
		}

		// Remove middle element and test
		dq.Remove(1) // Remove element at index 1 (value 2)
		if dq.Contains(2) {
			t.Errorf("Expected deque not to contain 2 after removal")
		}
		if !dq.Contains(1) || !dq.Contains(3) {
			t.Errorf("Expected deque to still contain 1 and 3")
		}

		// Test with negative numbers
		negatives := g.DequeOf(-5, -10, 0, 5, 10)
		if !negatives.Contains(-5) {
			t.Errorf("Expected deque to contain -5")
		}
		if !negatives.Contains(-10) {
			t.Errorf("Expected deque to contain -10")
		}
		if negatives.Contains(-1) {
			t.Errorf("Expected deque not to contain -1")
		}
	})

	t.Run("with_circular_buffer", func(t *testing.T) {
		// Test Contains when deque wraps around (circular buffer behavior)
		dq := g.NewDeque[int]()

		// Fill deque and then pop/push to cause wrapping
		for i := 0; i < 10; i++ {
			dq.PushBack(i)
		}

		// Remove front elements
		for i := 0; i < 5; i++ {
			dq.PopFront()
		}

		// Add more elements to cause wrapping
		for i := 10; i < 15; i++ {
			dq.PushBack(i)
		}

		// Now deque should contain: [5, 6, 7, 8, 9, 10, 11, 12, 13, 14]
		// Test contains on wrapped buffer
		for i := 5; i < 15; i++ {
			if !dq.Contains(i) {
				t.Errorf("Expected wrapped deque to contain %d", i)
			}
		}

		// Test elements that should not be there
		for i := 0; i < 5; i++ {
			if dq.Contains(i) {
				t.Errorf("Expected wrapped deque not to contain %d", i)
			}
		}
		if dq.Contains(15) {
			t.Errorf("Expected wrapped deque not to contain 15")
		}
	})
}

func TestDequeContainsEdgeCases(t *testing.T) {
	// Test Contains on empty deque
	empty := g.NewDeque[int]()
	if empty.Contains(42) {
		t.Errorf("Empty deque should not contain any element")
	}

	// Test Contains with single element
	single := g.DequeOf(42)
	if !single.Contains(42) {
		t.Errorf("Single element deque should contain its element")
	}
	if single.Contains(43) {
		t.Errorf("Single element deque should not contain other elements")
	}
}

func TestDequeIndexEdgeCases(t *testing.T) {
	// Test Index on empty deque
	empty := g.NewDeque[int]()
	if empty.Index(42) != -1 {
		t.Errorf("Empty deque should return -1 for any element")
	}

	// Test Index with duplicates at edges
	dq := g.DequeOf(5, 1, 2, 3, 5)
	idx := dq.Index(5)
	if idx != 0 {
		t.Errorf("Index should return first occurrence, expected 0, got %d", idx)
	}
}

func TestDequeBinarySearch(t *testing.T) {
	// Test BinarySearch on sorted deque
	dq := g.DequeOf(1, 3, 5, 7, 9)

	// Test finding existing elements
	idx, found := dq.BinarySearch(5, cmp.Cmp[int])
	if !found || idx != 2 {
		t.Errorf("Expected to find 5 at index 2, got index %d, found %t", idx, found)
	}

	idx, found = dq.BinarySearch(1, cmp.Cmp[int])
	if !found || idx != 0 {
		t.Errorf("Expected to find 1 at index 0, got index %d, found %t", idx, found)
	}

	idx, found = dq.BinarySearch(9, cmp.Cmp[int])
	if !found || idx != 4 {
		t.Errorf("Expected to find 9 at index 4, got index %d, found %t", idx, found)
	}

	// Test finding non-existing elements (should return insertion point)
	idx, found = dq.BinarySearch(4, cmp.Cmp[int])
	if found || idx != 2 {
		t.Errorf("Expected not to find 4, insertion point should be 2, got index %d, found %t", idx, found)
	}

	idx, found = dq.BinarySearch(0, cmp.Cmp[int])
	if found || idx != 0 {
		t.Errorf("Expected not to find 0, insertion point should be 0, got index %d, found %t", idx, found)
	}

	idx, found = dq.BinarySearch(10, cmp.Cmp[int])
	if found || idx != 5 {
		t.Errorf("Expected not to find 10, insertion point should be 5, got index %d, found %t", idx, found)
	}

	// Test BinarySearch on empty deque
	empty := g.NewDeque[int]()
	idx, found = empty.BinarySearch(5, cmp.Cmp[int])
	if found || idx != 0 {
		t.Errorf("Expected not to find 5 in empty deque, insertion point should be 0, got index %d, found %t", idx, found)
	}

	// Test BinarySearch on single element
	single := g.DequeOf(5)
	idx, found = single.BinarySearch(5, cmp.Cmp[int])
	if !found || idx != 0 {
		t.Errorf("Expected to find 5 at index 0 in single element deque, got index %d, found %t", idx, found)
	}

	idx, found = single.BinarySearch(3, cmp.Cmp[int])
	if found || idx != 0 {
		t.Errorf("Expected not to find 3, insertion point should be 0, got index %d, found %t", idx, found)
	}

	idx, found = single.BinarySearch(7, cmp.Cmp[int])
	if found || idx != 1 {
		t.Errorf("Expected not to find 7, insertion point should be 1, got index %d, found %t", idx, found)
	}
}

func TestDequeEqEdgeCases(t *testing.T) {
	// Test Eq with empty deques
	empty1 := g.NewDeque[int]()
	empty2 := g.NewDeque[int]()
	if !empty1.Eq(empty2) {
		t.Errorf("Two empty deques should be equal")
	}

	// Test Eq with one empty, one non-empty
	nonEmpty := g.DequeOf(1)
	if empty1.Eq(nonEmpty) {
		t.Errorf("Empty and non-empty deques should not be equal")
	}

	// Test Eq with nil pointer (should not crash)
	// This test depends on implementation - some may panic, others may handle gracefully
	// For safety, let's skip this test or make it implementation-specific

	// Test Eq with same reference
	dq := g.DequeOf(1, 2, 3)
	if !dq.Eq(dq) {
		t.Errorf("Deque should be equal to itself")
	}
}

func TestDequePrint(t *testing.T) {
	// Test Print method - should return the deque unchanged
	dq := g.DequeOf(1, 2, 3)
	result := dq.Print()

	// Should return the same deque instance
	if result != dq {
		t.Errorf("Print() should return the same deque instance")
	}

	// Deque should be unchanged
	if dq.Len() != 3 {
		t.Errorf("Print() should not modify deque, expected length 3, got %d", dq.Len())
	}
}

func TestDequePrintln(t *testing.T) {
	// Test Println method - should return the deque unchanged
	dq := g.DequeOf(1, 2, 3)
	result := dq.Println()

	// Should return the same deque instance
	if result != dq {
		t.Errorf("Println() should return the same deque instance")
	}

	// Deque should be unchanged
	if dq.Len() != 3 {
		t.Errorf("Println() should not modify deque, expected length 3, got %d", dq.Len())
	}
}
