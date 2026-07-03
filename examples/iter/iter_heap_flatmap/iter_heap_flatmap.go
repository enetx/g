package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	// Example 1: FlatMap to expand priority levels
	priorities := NewHeap(cmp.Cmp[int])
	priorities.Push(1, 3, 5)

	expanded := priorities.Iter().FlatMap(func(priority int) SeqHeap[int] {
		// Create sub-tasks for each priority level
		subHeap := NewHeap(cmp.Cmp[int])
		subHeap.Push(priority, priority*10, priority*100)
		return subHeap.Iter()
	}).Collect(cmp.Cmp)

	Print("Expanded priorities: ")
	expanded.Println() // Output will be sorted by heap order

	// Example 2: FlatMap to generate task hierarchies
	taskLevels := NewHeap(cmp.Cmp[string])
	taskLevels.Push("urgent", "normal", "low")

	taskHierarchy := taskLevels.Iter().FlatMap(func(level string) SeqHeap[string] {
		tasks := NewHeap(cmp.Cmp[string])
		tasks.Push(level+"-task1", level+"-task2")
		return tasks.Iter()
	}).Collect(cmp.Cmp)

	Print("Task hierarchy: ")
	taskHierarchy.Println()

	// Example 3: FlatMap to expand nested scoring systems
	baseScores := NewHeap(cmp.Cmp[float64])
	baseScores.Push(1.5, 2.0, 3.5)

	weightedScores := baseScores.Iter().FlatMap(func(base float64) SeqHeap[float64] {
		weights := NewHeap(cmp.Cmp[float64])
		// Apply different multipliers
		weights.Push(base*1.0, base*1.5, base*2.0)
		return weights.Iter()
	}).Collect(cmp.Cmp)

	Print("Weighted scores: ")
	weightedScores.Println()

	// Example 4: FlatMap with conditional expansion (priorities only)
	numbers := NewHeap(cmp.Cmp[int])
	numbers.Push(-1, 0, 1, 2, 3)

	positiveExpanded := numbers.Iter().FlatMap(func(n int) SeqHeap[int] {
		if n > 0 {
			// Expand positive numbers into ranges
			expanded := NewHeap(cmp.Cmp[int])
			for i := n; i <= n+2; i++ {
				expanded.Push(i)
			}
			return expanded.Iter()
		}
		// Return empty heap for non-positive numbers
		return NewHeap(cmp.Cmp[int]).Iter()
	}).Collect(cmp.Cmp)

	Print("Positive expanded: ")
	positiveExpanded.Println()

	// Example 5: FlatMap for multi-level categorization
	categories := NewHeap(cmp.Cmp[string])
	categories.Push("tech", "science")

	subcategories := categories.Iter().FlatMap(func(category string) SeqHeap[string] {
		subCats := NewHeap(cmp.Cmp[string])
		switch category {
		case "tech":
			subCats.Push("ai", "web", "mobile")
		case "science":
			subCats.Push("physics", "chemistry", "biology")
		}
		return subCats.Iter()
	}).Collect(cmp.Cmp)

	Print("All subcategories: ")
	subcategories.Println()

	// Example 6: FlatMap with number expansion
	bases := NewHeap(cmp.Cmp[int])
	bases.Push(2, 5, 10)

	multiplied := bases.Iter().FlatMap(func(base int) SeqHeap[int] {
		multiples := NewHeap(cmp.Cmp[int])
		multiples.Push(base*1, base*2, base*3)
		return multiples.Iter()
	}).Collect(cmp.Cmp)

	Print("Multiplied values: ")
	multiplied.Println()
}
