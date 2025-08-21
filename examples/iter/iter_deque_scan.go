package main

import (
	. "github.com/enetx/g"
)

func main() {
	// Example 1: Scan to calculate running sum
	numbers := DequeOf(1, 2, 3, 4, 5)
	runningSums := numbers.Iter().Scan(0, func(acc, val int) int {
		return acc + val
	}).Collect()

	runningSums.Println() // Deque[0, 1, 3, 6, 10, 15]

	// Example 2: Scan to calculate running product
	factors := DequeOf(2, 3, 4)
	runningProducts := factors.Iter().Scan(1, func(acc, val int) int {
		return acc * val
	}).Collect()

	runningProducts.Println() // Deque[1, 2, 6, 24]

	// Example 3: Scan to build running maximum
	values := DequeOf(3, 1, 4, 1, 5, 9, 2)
	runningMax := values.Iter().Scan(0, func(acc, val int) int {
		if val > acc {
			return val
		}
		return acc
	}).Collect()

	runningMax.Println() // Deque[0, 3, 3, 4, 4, 5, 9, 9]

	// Example 4: Scan to concatenate strings
	words := DequeOf("Hello", " ", "World", "!")
	buildingSentence := words.Iter().Scan("", func(acc, val string) string {
		return acc + val
	}).Collect()

	buildingSentence.Println() // Deque[, Hello,  Hello , Hello World, Hello World!]

	// Example 5: Scan to track balance changes
	transactions := DequeOf(100, -30, -20, 50, -10)
	balanceHistory := transactions.Iter().Scan(0, func(balance, transaction int) int {
		return balance + transaction
	}).Collect()

	Print("Balance history: ")
	balanceHistory.Println() // Deque[0, 100, 70, 50, 100, 90]

	// Example 6: Scan for simple counting
	items := DequeOf(1, 2, 3, 4, 5)
	indexCounts := items.Iter().Scan(0, func(count, item int) int {
		return count + 1 // Just count items
	}).Collect()

	Print("Item counts: ")
	indexCounts.Println() // Deque[0, 1, 2, 3, 4, 5]

	// Example 7: Scan for vowel counting
	vowelCounts := DequeOf(1, 2, 1, 1, 3) // vowel count per word
	totalVowels := vowelCounts.Iter().Scan(0, func(total, vowels int) int {
		return total + vowels
	}).Collect()

	totalVowels.Println() // Deque[0, 1, 3, 4, 5, 8]

	// Example 8: Scan for character counting
	wordLengths := DequeOf(3, 5, 5, 3) // lengths of "The", "quick", "brown", "fox"
	charCounts := wordLengths.Iter().Scan(0, func(total, length int) int {
		return total + length
	}).Collect()

	Print("Character count progression: ")
	charCounts.Println() // Shows cumulative character counts
}
