package main

import (
	"fmt"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

// ==============================================================================
// HEAP EXAMPLES - COMPLETE GUIDE
// ==============================================================================
//
// A Heap is a specialized tree-based data structure that satisfies the heap property:
// - Min Heap: parent ≤ children (smallest element at top)
// - Max Heap: parent ≥ children (largest element at top)
//
// Key Operations:
// - Push(item): Add element - O(log n)
// - Pop(): Remove top element - O(log n)
// - Peek(): View top element - O(1)
// - Len(): Get size - O(1)
//
// Common Use Cases:
// - Priority queues
// - Task scheduling
// - Finding top-K elements
// - Heap sort algorithm
// - Event scheduling
// ==============================================================================

// Run all examples
func main() {
	BasicHeapOperations()
	MaxHeap()
	HeapFromSlice()
	TaskScheduler()
	TopKElements()
	EventScheduler()
	HeapSort()
	HeapIterators()
	CustomTypes()
	PerformanceComparison()
	ErrorHandling()
}

// Example 1: Basic Heap Operations
func BasicHeapOperations() {
	fmt.Println("=== Basic Heap Operations ===")

	// Create a min heap (smallest elements first)
	minHeap := NewHeap(cmp.Cmp[int])

	// Add elements
	minHeap.Push(10, 5, 15, 1, 8, 12)
	fmt.Printf("Heap size: %d\n", minHeap.Len()) // 6

	// Peek at top element (smallest)
	if top := minHeap.Peek(); top.IsSome() {
		fmt.Printf("Top element: %d\n", top.Some()) // 1
	}

	// Extract elements in sorted order
	fmt.Print("Elements in order: ")
	for !minHeap.IsEmpty() {
		if elem := minHeap.Pop(); elem.IsSome() {
			fmt.Printf("%d ", elem.Some())
		}
	}
	fmt.Println() // Output: 1 5 8 10 12 15
}

// Example 2: Max Heap (Priority Queue)
func MaxHeap() {
	fmt.Println("\n=== Max Heap (Priority Queue) ===")

	// Create a max heap by reversing comparison
	maxHeap := NewHeap[Int](func(a, b Int) cmp.Ordering {
		return b.Cmp(a) // Reverse order for max heap
	})

	// Add elements
	maxHeap.Push(10, 5, 15, 1, 8, 12)

	// Extract largest elements first
	fmt.Print("Largest to smallest: ")
	for !maxHeap.IsEmpty() {
		if elem := maxHeap.Pop(); elem.IsSome() {
			fmt.Printf("%d ", elem.Some())
		}
	}
	fmt.Println() // Output: 15 12 10 8 5 1
}

// Example 3: Creating Heap from Slice - O(n) Construction
func HeapFromSlice() {
	fmt.Println("\n=== Heap from Slice (Fast Construction) ===")

	// Create slice of numbers
	numbers := SliceOf(64, 34, 25, 12, 22, 11, 90, 5, 77, 30)
	fmt.Printf("Original slice: %v\n", numbers)

	// Convert to min heap in O(n) time
	heap := numbers.Heap(cmp.Cmp[int])

	// Get all elements in sorted order
	sorted := heap.Iter().Collect(cmp.Cmp)
	fmt.Printf("Sorted: %v\n", sorted)
	// Output: [5, 11, 12, 22, 25, 30, 34, 64, 77, 90]
}

// Example 4: Priority Task Scheduler
func TaskScheduler() {
	fmt.Println("\n=== Priority Task Scheduler ===")

	type Task struct {
		Name     string
		Priority int
		Deadline time.Time
	}

	now := time.Now()
	tasks := SliceOf(
		Task{"Email", 2, now.Add(time.Hour)},
		Task{"Meeting", 5, now.Add(30 * time.Minute)},
		Task{"Report", 1, now.Add(2 * time.Hour)},
		Task{"Call", 4, now.Add(15 * time.Minute)},
		Task{"Review", 3, now.Add(45 * time.Minute)},
	)

	// Create priority queue: higher priority first, then by deadline
	scheduler := tasks.Heap(func(a, b Task) cmp.Ordering {
		// First by priority (higher = better)
		return cmp.Cmp(b.Priority, a.Priority).
			// Then by deadline (earlier = better)
			Then(cmp.Cmp(a.Deadline.Unix(), b.Deadline.Unix()))
	})

	fmt.Println("Tasks in execution order:")
	for i := 1; !scheduler.IsEmpty(); i++ {
		if task := scheduler.Pop(); task.IsSome() {
			t := task.Some()
			fmt.Printf("%d. %s (Priority: %d)\n", i, t.Name, t.Priority)
		}
	}
	// Output: Meeting(5) → Call(4) → Review(3) → Email(2) → Report(1)
}

// Example 5: Finding Top-K Elements
func TopKElements() {
	fmt.Println("\n=== Finding Top-K Elements ===")

	// Large dataset
	data := SliceOf(85, 23, 67, 91, 45, 78, 34, 56, 92, 12, 89, 76, 43, 98, 21)
	k := Int(5)

	fmt.Printf("Original data: %v\n", data)

	// Method 1: Using min heap of size k (memory efficient for large data)
	minHeap := NewHeap(cmp.Cmp[int])

	for _, num := range data {
		if minHeap.Len() < k {
			minHeap.Push(num)
		} else if topK := minHeap.Peek(); topK.IsSome() && num > topK.Some() {
			minHeap.Pop()
			minHeap.Push(num)
		}
	}

	topK := minHeap.Iter().Collect(cmp.Cmp)
	fmt.Printf("Top %d elements: %v\n", k, topK)

	// Method 2: Using full heap (simpler but uses more memory)
	allSorted := data.Heap(func(a, b int) cmp.Ordering {
		return cmp.Cmp(b, a) // Max heap
	}).Iter().Take(uint(k)).Collect(cmp.Cmp)

	fmt.Printf("Top %d (method 2): %v\n", k, allSorted)
}

// Example 6: Event Scheduling System
func EventScheduler() {
	fmt.Println("\n=== Event Scheduling System ===")

	type Event struct {
		Name string
		Time time.Time
		Type string
	}

	now := time.Now()
	events := SliceOf(
		Event{"Meeting", now.Add(2 * time.Hour), "Work"},
		Event{"Lunch", now.Add(30 * time.Minute), "Personal"},
		Event{"Call", now.Add(15 * time.Minute), "Work"},
		Event{"Gym", now.Add(4 * time.Hour), "Personal"},
		Event{"Review", now.Add(1 * time.Hour), "Work"},
	)

	// Schedule events by time (earliest first)
	scheduler := events.Heap(func(a, b Event) cmp.Ordering {
		return cmp.Cmp(a.Time.Unix(), b.Time.Unix())
	})

	fmt.Println("Events in chronological order:")
	for !scheduler.IsEmpty() {
		if event := scheduler.Pop(); event.IsSome() {
			e := event.Some()
			duration := e.Time.Sub(now).Round(time.Minute)
			fmt.Printf("- %s (%s) in %v\n", e.Name, e.Type, duration)
		}
	}
}

// Example 7: Heap Sort Algorithm
func HeapSort() {
	fmt.Println("\n=== Heap Sort Algorithm ===")

	unsorted := SliceOf(64, 34, 25, 12, 22, 11, 90, 88, 76, 50, 42)
	fmt.Printf("Unsorted: %v\n", unsorted)

	// Heap sort in ascending order
	ascending := unsorted.Heap(cmp.Cmp).Iter().Collect(cmp.Cmp)
	fmt.Printf("Ascending: %v\n", ascending)

	// Heap sort in descending order
	descending := unsorted.Heap(func(a, b int) cmp.Ordering {
		return cmp.Cmp(b, a)
	}).Iter().Collect(cmp.Reverse)
	fmt.Printf("Descending: %v\n", descending)
}

// Example 8: Iterators and Functional Programming
func HeapIterators() {
	fmt.Println("\n=== Heap Iterators ===")

	numbers := SliceOf(1, 5, 3, 8, 2, 9, 4, 7, 6)
	heap := numbers.Heap(cmp.Cmp)

	// Non-consuming iteration (heap remains intact)
	fmt.Print("All elements (heap unchanged): ")
	heap.Iter().ForEach(func(x int) {
		fmt.Printf("%d ", x)
	})

	fmt.Printf("\nHeap size after iteration: %d\n", heap.Len()) // Still 9

	// Functional operations on heap iterator
	evenNumbers := heap.Iter().
		Filter(func(x int) bool { return x%2 == 0 }).
		Collect(cmp.Cmp)
	fmt.Printf("Even numbers: %v\n", evenNumbers)

	// Take first N elements
	firstThree := heap.Iter().Take(3).Collect(cmp.Cmp)
	fmt.Printf("First 3 smallest: %v\n", firstThree)

	// Consuming iteration (empties the heap)
	fmt.Print("Consuming iteration: ")
	heap.IntoIter().ForEach(func(x int) {
		fmt.Printf("%d ", x)
	})

	fmt.Printf("\nHeap size after consuming: %d\n", heap.Len()) // Now 0
}

// Example 9: Custom Complex Types
func CustomTypes() {
	fmt.Println("\n=== Custom Complex Types ===")

	type Student struct {
		Name  String
		Grade Float
		Age   Int
	}

	students := SliceOf(
		Student{"Alice", 92.5, 20},
		Student{"Bob", 87.3, 19},
		Student{"Charlie", 95.1, 21},
		Student{"Diana", 89.7, 20},
		Student{"Eve", 91.2, 18},
	)

	// Multi-level comparison: Grade (desc), then Age (asc), then Name (asc)
	topStudents := students.Heap(func(a, b Student) cmp.Ordering {
		// Primary: Grade (higher is better)
		return b.Grade.Cmp(a.Grade).
			// Secondary: Age (younger is better)
			Then(a.Age.Cmp(b.Age)).
			// Tertiary: Name (alphabetical)
			Then(a.Name.Cmp(b.Name))
	})

	fmt.Println("Students by rank:")
	for i := 1; !topStudents.IsEmpty(); i++ {
		if student := topStudents.Pop(); student.IsSome() {
			s := student.Some()
			fmt.Printf("%d. %s (Grade: %.1f, Age: %d)\n", i, s.Name, s.Grade, s.Age)
		}
	}
}

// Example 10: Performance Comparison
func PerformanceComparison() {
	fmt.Println("\n=== Performance Comparison ===")

	// Large dataset
	data := RangeInclusive(1000, 0, -1).Collect()

	fmt.Printf("Dataset size: %d elements\n", len(data))

	// Method 1: Build heap from slice - O(n)
	start := time.Now()
	heap1 := data.Heap(cmp.Cmp)
	duration1 := time.Since(start)
	fmt.Printf("Heap() construction: %v\n", duration1)

	// Method 2: Build heap with individual pushes - O(n log n)
	start = time.Now()
	heap2 := NewHeap(cmp.Cmp[int])
	for _, item := range data {
		heap2.Push(item)
	}

	duration2 := time.Since(start)
	fmt.Printf("Individual Push() construction: %v\n", duration2)

	// Verify both heaps produce same result
	result1 := heap1.Iter().Take(10)
	result2 := heap2.Iter().Take(10)
	fmt.Printf("Results identical: %v\n", result1.Eq(result2))
	fmt.Printf("First 10 elements: %v\n", result1.Collect(cmp.Cmp))
}

// Example 11: Error Handling and Edge Cases
func ErrorHandling() {
	fmt.Println("\n=== Error Handling and Edge Cases ===")

	// Empty heap
	emptyHeap := NewHeap(cmp.Cmp[int])

	if top := emptyHeap.Pop(); top.IsNone() {
		fmt.Println("✓ Pop from empty heap returns None")
	}

	if peek := emptyHeap.Peek(); peek.IsNone() {
		fmt.Println("✓ Peek at empty heap returns None")
	}

	fmt.Printf("Empty heap size: %d\n", emptyHeap.Len())

	// Single element
	singleHeap := NewHeap(cmp.Cmp[int])
	singleHeap.Push(42)

	if elem := singleHeap.Pop(); elem.IsSome() {
		fmt.Printf("✓ Single element heap works: %d\n", elem.Some())
	}

	// Duplicate elements
	duplicates := SliceOf[Int](5, 3, 5, 1, 3, 5, 1, 1)
	dupHeap := duplicates.Heap(Int.Cmp)
	sorted := dupHeap.Iter().Collect(Int.Cmp)
	fmt.Printf("✓ Handles duplicates: %v\n", sorted)
}
