package main

import (
	. "github.com/enetx/g"
)

func main() {
	// Example 1: FlatMap to expand numbers into ranges
	numbers := DequeOf(1, 3, 5)
	expanded := numbers.Iter().FlatMap(func(n int) SeqDeque[int] {
		// Create a range from n to n+2
		return DequeOf(n, n+1, n+2).Iter()
	}).Collect()

	expanded.Println() // Deque[1, 2, 3, 3, 4, 5, 5, 6, 7]

	// Example 2: FlatMap to split strings into characters
	words := DequeOf("hi", "go", "rust")
	chars := words.Iter().FlatMap(func(word string) SeqDeque[string] {
		result := NewDeque[string]()
		for _, char := range word {
			result.PushBack(string(char))
		}
		return result.Iter()
	}).Collect()

	chars.Println() // Deque[h, i, g, o, r, u, s, t]

	// Example 3: FlatMap to generate multiplication tables
	bases := DequeOf(2, 3)
	multTable := bases.Iter().FlatMap(func(base int) SeqDeque[int] {
		table := NewDeque[int]()
		for i := 1; i <= 3; i++ {
			table.PushBack(base * i)
		}
		return table.Iter()
	}).Collect()

	Print("Multiplication results: ")
	multTable.Println()

	// Example 4: FlatMap to duplicate elements
	duplicateBase := DequeOf(10, 20, 30)
	duplicated := duplicateBase.Iter().FlatMap(func(n int) SeqDeque[int] {
		return DequeOf(n, n).Iter() // Duplicate each number
	}).Collect()

	duplicated.Println() // Deque[10, 10, 20, 20, 30, 30]

	// Example 5: FlatMap with conditional expansion
	values := DequeOf(-1, 0, 1, 2)
	conditionalExpanded := values.Iter().FlatMap(func(n int) SeqDeque[int] {
		if n <= 0 {
			return NewDeque[int]().Iter() // Empty for non-positive
		}
		// Duplicate positive numbers
		return DequeOf(n, n).Iter()
	}).Collect()

	conditionalExpanded.Println() // Deque[1, 1, 2, 2]
}
