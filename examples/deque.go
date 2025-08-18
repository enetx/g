package main

import (
	"fmt"
	"time"

	. "github.com/enetx/g"
)

// ==============================================================================
// DEQUE EXAMPLES - COMPLETE GUIDE
// ==============================================================================
//
// A Deque (Double-Ended Queue) is a data structure that allows efficient
// insertion and deletion at both ends. It's implemented as a growable ring buffer.
//
// Key Operations:
// - PushFront/PushBack: Add element to front/back - O(1) amortized
// - PopFront/PopBack: Remove element from front/back - O(1)
// - Front/Back: View front/back element - O(1)
// - Get/Set: Access by index - O(1)
//
// Common Use Cases:
// - Sliding window algorithms
// - Undo/Redo systems
// - Breadth-first search
// - Task queues
// - Browser history
// - Palindrome checking
// ==============================================================================

// Run all examples
func main() {
	BasicDequeOperations()
	SlidingWindowMaximum()
	UndoRedoSystem()
	BreadthFirstSearch()
	TaskQueue()
	BrowserHistory()
	PalindromeChecker()
	DequeIterators()
	RingBufferBehavior()
	PerformanceComparison2()
	ErrorHandling2()
}

// Example 1: Basic Deque Operations
func BasicDequeOperations() {
	fmt.Println("=== Basic Deque Operations ===")

	// Create empty deque
	dq := NewDeque[int]()

	// Add elements to both ends
	dq.PushBack(1)   // [1]
	dq.PushFront(0)  // [0, 1]
	dq.PushBack(2)   // [0, 1, 2]
	dq.PushFront(-1) // [-1, 0, 1, 2]

	fmt.Printf("Deque: %s\n", dq) // Deque[-1, 0, 1, 2]
	fmt.Printf("Length: %d\n", dq.Len())

	// Access front and back
	if front := dq.Front(); front.IsSome() {
		fmt.Printf("Front: %d\n", front.Some()) // -1
	}
	if back := dq.Back(); back.IsSome() {
		fmt.Printf("Back: %d\n", back.Some()) // 2
	}

	// Remove from both ends
	if elem := dq.PopFront(); elem.IsSome() {
		fmt.Printf("Popped front: %d\n", elem.Some()) // -1
	}
	if elem := dq.PopBack(); elem.IsSome() {
		fmt.Printf("Popped back: %d\n", elem.Some()) // 2
	}

	fmt.Printf("After pops: %v\n", dq) // [0, 1]
}

// Example 2: Sliding Window Maximum
func SlidingWindowMaximum() {
	fmt.Println("\n=== Sliding Window Maximum ===")

	// Find maximum in each sliding window of size k
	arr := SliceOf(1, 3, -1, -3, 5, 3, 6, 7)
	k := 3

	fmt.Printf("Array: %v, Window size: %d\n", arr, k)

	// Deque stores indices of array elements in decreasing order of their values
	dq := NewDeque[int]()
	result := NewSlice[int]()

	for i, num := range arr {
		// Remove indices that are out of current window
		for !dq.IsEmpty() {
			if front := dq.Front(); front.IsSome() && front.Some() <= i-k {
				dq.PopFront()
			} else {
				break
			}
		}

		// Remove indices whose corresponding values are smaller than current
		for !dq.IsEmpty() {
			if back := dq.Back(); back.IsSome() && arr[back.Some()] <= num {
				dq.PopBack()
			} else {
				break
			}
		}

		dq.PushBack(i)

		// The front of deque contains index of maximum element of current window
		if i >= k-1 {
			if front := dq.Front(); front.IsSome() {
				result.Push(arr[front.Some()])
			}
		}
	}

	fmt.Printf("Window maximums: %v\n", result) // [3, 3, 5, 5, 6, 7]
}

// Example 3: Undo/Redo System
func UndoRedoSystem() {
	fmt.Println("\n=== Undo/Redo System ===")

	type Command struct {
		Action string
		Data   string
	}

	type Document struct {
		content   String
		undoStack *Deque[Command]
		redoStack *Deque[Command]
	}

	doc := &Document{
		content:   "",
		undoStack: NewDeque[Command](),
		redoStack: NewDeque[Command](),
	}

	execute := func(cmd Command) {
		fmt.Printf("Executing: %s '%s'\n", cmd.Action, cmd.Data)

		// Save current state for undo
		undoCmd := Command{
			Action: "set",
			Data:   doc.content.Std(),
		}
		doc.undoStack.PushBack(undoCmd)

		// Clear redo stack on new action
		doc.redoStack.Clear()

		// Apply command
		switch cmd.Action {
		case "append":
			doc.content += String(cmd.Data)
		case "set":
			doc.content = String(cmd.Data)
		}

		fmt.Printf("Document: '%s'\n", doc.content)
	}

	undo := func() {
		if undoCmd := doc.undoStack.PopBack(); undoCmd.IsSome() {
			// Save current state for redo
			redoCmd := Command{
				Action: "set",
				Data:   doc.content.Std(),
			}
			doc.redoStack.PushBack(redoCmd)

			// Apply undo
			cmd := undoCmd.Some()
			doc.content = String(cmd.Data)
			fmt.Printf("Undo: Document: '%s'\n", doc.content)
		} else {
			fmt.Println("Nothing to undo")
		}
	}

	redo := func() {
		if redoCmd := doc.redoStack.PopBack(); redoCmd.IsSome() {
			// Save current state for undo
			undoCmd := Command{
				Action: "set",
				Data:   doc.content.Std(),
			}
			doc.undoStack.PushBack(undoCmd)

			// Apply redo
			cmd := redoCmd.Some()
			doc.content = String(cmd.Data)
			fmt.Printf("Redo: Document: '%s'\n", doc.content)
		} else {
			fmt.Println("Nothing to redo")
		}
	}

	// Simulate editing
	execute(Command{"append", "Hello"})
	execute(Command{"append", " World"})
	execute(Command{"append", "!"})

	undo() // Remove "!"
	undo() // Remove " World"
	redo() // Add back " World"

	execute(Command{"append", "."}) // This clears redo stack
	undo()                          // Remove "."
}

// Example 4: Breadth-First Search
func BreadthFirstSearch() {
	fmt.Println("\n=== Breadth-First Search ===")

	type Node struct {
		Value    string
		Children []*Node
	}

	// Build tree
	root := &Node{Value: "A"}
	b := &Node{Value: "B"}
	c := &Node{Value: "C"}
	d := &Node{Value: "D"}
	e := &Node{Value: "E"}
	f := &Node{Value: "F"}

	root.Children = []*Node{b, c}
	b.Children = []*Node{d, e}
	c.Children = []*Node{f}

	// BFS traversal using deque as queue
	queue := NewDeque[*Node]()
	visited := NewSlice[string]()

	queue.PushBack(root)

	fmt.Print("BFS traversal: ")
	for !queue.IsEmpty() {
		if nodeOpt := queue.PopFront(); nodeOpt.IsSome() {
			node := nodeOpt.Some()
			visited.Push(node.Value)
			fmt.Printf("%s ", node.Value)

			// Add children to queue
			for _, child := range node.Children {
				queue.PushBack(child)
			}
		}
	}
	fmt.Printf("\nVisited order: %v\n", visited) // [A, B, C, D, E, F]
}

// Example 5: Task Queue with Priorities
func TaskQueue() {
	fmt.Println("\n=== Task Queue ===")

	type Task struct {
		ID       int
		Name     string
		Priority string // "high", "normal", "low"
	}

	// Separate queues for different priorities
	highPriority := NewDeque[Task]()
	normalPriority := NewDeque[Task]()
	lowPriority := NewDeque[Task]()

	addTask := func(task Task) {
		switch task.Priority {
		case "high":
			highPriority.PushBack(task)
		case "normal":
			normalPriority.PushBack(task)
		case "low":
			lowPriority.PushBack(task)
		}
		fmt.Printf("Added task: %s (%s priority)\n", task.Name, task.Priority)
	}

	processNext := func() bool {
		// Process high priority first, then normal, then low
		queues := []*Deque[Task]{highPriority, normalPriority, lowPriority}
		priorities := []string{"high", "normal", "low"}

		for i, queue := range queues {
			if task := queue.PopFront(); task.IsSome() {
				t := task.Some()
				fmt.Printf("Processing: %s (%s priority)\n", t.Name, priorities[i])
				return true
			}
		}
		return false
	}

	// Add some tasks
	addTask(Task{1, "Send email", "normal"})
	addTask(Task{2, "Fix critical bug", "high"})
	addTask(Task{3, "Update docs", "low"})
	addTask(Task{4, "Security patch", "high"})
	addTask(Task{5, "Code review", "normal"})

	fmt.Println("\nProcessing tasks:")
	for processNext() {
		// Continue until all queues are empty
	}
}

// Example 6: Browser History
func BrowserHistory() {
	fmt.Println("\n=== Browser History ===")

	type BrowserSession struct {
		history *Deque[string]
		current int
	}

	browser := &BrowserSession{
		history: NewDeque[string](),
		current: -1,
	}

	visit := func(url string) {
		// Remove all history after current position
		for browser.current < int(browser.history.Len())-1 {
			browser.history.PopBack()
		}

		browser.history.PushBack(url)
		browser.current = int(browser.history.Len()) - 1
		fmt.Printf("Visited: %s\n", url)
	}

	back := func() string {
		if browser.current > 0 {
			browser.current--
			if page := browser.history.Get(Int(browser.current)); page.IsSome() {
				url := page.Some()
				fmt.Printf("Back to: %s\n", url)
				return url
			}
		}
		fmt.Println("Can't go back")
		return ""
	}

	forward := func() string {
		if browser.current < int(browser.history.Len())-1 {
			browser.current++
			if page := browser.history.Get(Int(browser.current)); page.IsSome() {
				url := page.Some()
				fmt.Printf("Forward to: %s\n", url)
				return url
			}
		}
		fmt.Println("Can't go forward")
		return ""
	}

	// Simulate browsing
	visit("google.com")
	visit("github.com")
	visit("stackoverflow.com")

	back()    // github.com
	back()    // google.com
	forward() // github.com

	visit("reddit.com") // This removes stackoverflow.com from history
	back()              // github.com
	forward()           // reddit.com
}

// Example 7: Palindrome Checker
func PalindromeChecker() {
	fmt.Println("\n=== Palindrome Checker ===")

	isPalindrome := func(s string) bool {
		dq := NewDeque[rune]()

		// Add all characters to deque
		for _, char := range s {
			if char != ' ' { // Ignore spaces
				dq.PushBack(char)
			}
		}

		// Compare from both ends
		for dq.Len() > 1 {
			front := dq.PopFront()
			back := dq.PopBack()

			if front.IsNone() || back.IsNone() || front.Some() != back.Some() {
				return false
			}
		}

		return true
	}

	testCases := SliceOf[String]("racecar", "hello", "A man a plan a canal Panama", "race a car", "madam")

	for _, test := range testCases {
		result := isPalindrome(test.Lower().Std())
		fmt.Printf("'%s': %t\n", test, result)
	}
}

// Example 8: Iterators and Functional Programming
func DequeIterators() {
	fmt.Println("\n=== Deque Iterators ===")

	dq := DequeOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	fmt.Printf("Original: %s\n", dq)

	// Forward iteration
	fmt.Print("Forward: ")
	dq.Iter().ForEach(func(x int) {
		fmt.Printf("%d ", x)
	})
	fmt.Println()

	// Reverse iteration
	fmt.Print("Reverse: ")
	dq.IterReverse().ForEach(func(x int) {
		fmt.Printf("%d ", x)
	})
	fmt.Println()

	// Functional operations
	evenNumbers := dq.Iter().
		Filter(func(x int) bool { return x%2 == 0 }).
		Collect()
	fmt.Printf("Even numbers: %s\n", evenNumbers)

	// Chain with transformations
	result := dq.Iter().
		Filter(func(x int) bool { return x > 5 }).
		Map(func(x int) int { return x * 2 }).
		Collect()
	fmt.Printf("Filtered and doubled: %s\n", result)
}

// Example 9: Ring Buffer Behavior
func RingBufferBehavior() {
	fmt.Println("\n=== Ring Buffer Behavior ===")

	// Create deque with initial capacity
	dq := NewDeque[int](5)
	fmt.Printf("Created deque with capacity 5\n")

	// Fill beyond initial capacity to show growth
	for i := 1; i <= 8; i++ {
		dq.PushBack(i)
		fmt.Printf("Added %d, length: %d, capacity: %d\n", i, dq.Len(), dq.Capacity())
	}

	// Show internal representation
	fmt.Printf("Deque: %s\n", dq)

	// Remove from front and add to back (typical ring buffer usage)
	for i := 0; i < 3; i++ {
		if elem := dq.PopFront(); elem.IsSome() {
			dq.PushBack(elem.Some() + 10)
			fmt.Printf("Moved front element to back: %s\n", dq)
		}
	}
}

// Example 10: Performance Comparison
func PerformanceComparison2() {
	fmt.Println("\n=== Performance Comparison ===")

	n := 100000
	fmt.Printf("Testing with %d operations\n", n)

	// Test 1: Deque vs Slice for front operations
	start := time.Now()
	dq := NewDeque[int]()
	for i := 0; i < n; i++ {
		dq.PushFront(i)
	}
	for i := 0; i < n; i++ {
		dq.PopFront()
	}
	dequeTime := time.Since(start)

	start = time.Now()
	sl := NewSlice[int]()
	for i := 0; i < n; i++ {
		sl = append(Slice[int]{i}, sl...) // Insert at front
	}
	for i := 0; i < n; i++ {
		if len(sl) > 0 {
			sl = sl[1:] // Remove from front
		}
	}
	sliceTime := time.Since(start)

	fmt.Printf("Front operations:\n")
	fmt.Printf("  Deque: %v\n", dequeTime)
	fmt.Printf("  Slice: %v\n", sliceTime)
	fmt.Printf("  Deque is %.1fx faster\n", float64(sliceTime)/float64(dequeTime))

	// Test 2: Random access
	dq2 := NewDeque[int]()
	sl2 := NewSlice[int]()
	for i := 0; i < 1000; i++ {
		dq2.PushBack(i)
		sl2.Push(i)
	}

	start = time.Now()
	for i := 0; i < 10000; i++ {
		dq2.Get(Int(i % 1000))
	}
	dequeAccessTime := time.Since(start)

	start = time.Now()
	for i := 0; i < 10000; i++ {
		sl2.Get(Int(i % 1000))
	}
	sliceAccessTime := time.Since(start)

	fmt.Printf("Random access:\n")
	fmt.Printf("  Deque: %v\n", dequeAccessTime)
	fmt.Printf("  Slice: %v\n", sliceAccessTime)
}

// Example 11: Error Handling and Edge Cases
func ErrorHandling2() {
	fmt.Println("\n=== Error Handling and Edge Cases ===")

	// Empty deque
	empty := NewDeque[int]()

	if front := empty.PopFront(); front.IsNone() {
		fmt.Println("✓ PopFront from empty deque returns None")
	}

	if back := empty.PopBack(); back.IsNone() {
		fmt.Println("✓ PopBack from empty deque returns None")
	}

	if front := empty.Front(); front.IsNone() {
		fmt.Println("✓ Front of empty deque returns None")
	}

	if back := empty.Back(); back.IsNone() {
		fmt.Println("✓ Back of empty deque returns None")
	}

	// Index out of bounds
	if elem := empty.Get(0); elem.IsNone() {
		fmt.Println("✓ Get with invalid index returns None")
	}

	// Single element operations
	single := NewDeque[string]()
	single.PushBack("only")

	if elem := single.PopFront(); elem.IsSome() {
		fmt.Printf("✓ Single element deque works: %s\n", elem.Some())
	}

	// Large deque operations
	large := NewDeque[int]()
	for i := 0; i < 1000; i++ {
		large.PushBack(i)
	}

	fmt.Printf("✓ Large deque created with %d elements\n", large.Len())

	// Clear operation
	large.Clear()
	fmt.Printf("✓ Large deque cleared, size: %d\n", large.Len())
}
