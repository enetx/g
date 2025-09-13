package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	// Example 1: Scan for running total in a priority heap
	priorities := NewHeap(cmp.Cmp[int])
	priorities.Push(3, 1, 4, 1, 5, 9, 2)

	runningTotal := priorities.Iter().Scan(0, func(acc, val int) int {
		return acc + val
	}).Collect(cmp.Cmp)

	Print("Priority heap running totals: ")
	runningTotal.Println()

	// Example 2: Scan for maximum tracking in score heap
	scores := NewHeap(func(a, b int) cmp.Ordering { return cmp.Cmp(b, a) }) // Max heap
	scores.Push(85, 92, 78, 95, 88, 91)

	maxSoFar := scores.Iter().Scan(0, func(maxVal, score int) int {
		if score > maxVal {
			return score
		}
		return maxVal
	}).Collect(func(a, b int) cmp.Ordering { return cmp.Cmp(b, a) })

	Print("Running maximum: ")
	maxSoFar.Println()

	// Example 3: Scan for simple counting
	items := NewHeap(cmp.Cmp[int])
	items.Push(10, 20, 30, 40)

	itemCounts := items.Iter().Scan(0, func(count, _ int) int {
		return count + 1 // Just count items
	}).Collect(cmp.Cmp)

	Print("Item counts: ")
	itemCounts.Println()

	// Example 4: Scan for cumulative product
	factors := NewHeap(cmp.Cmp[int])
	factors.Push(2, 3, 4)

	products := factors.Iter().Scan(1, func(acc, factor int) int {
		return acc * factor
	}).Collect(cmp.Cmp)

	Print("Running products: ")
	products.Println()

	// Example 5: Scan for balance tracking
	transactions := NewHeap(cmp.Cmp[int])
	transactions.Push(-25, 100, -15, 50, -10)

	balances := transactions.Iter().Scan(0, func(balance, transaction int) int {
		return balance + transaction
	}).Collect(cmp.Cmp)

	Print("Balance progression: ")
	balances.Println()

	// Example 6: Scan with string concatenation (same type required)
	numbers := NewHeap(cmp.Cmp[int])
	numbers.Push(1, 2, 3, 4)

	runningSum := numbers.Iter().Scan(0, func(sum, num int) int {
		return sum + num
	}).Collect(cmp.Cmp)

	Print("Running sum: ")
	runningSum.Println()

	// Example 7: Scan for minimum tracking
	values := NewHeap(cmp.Cmp[int])
	values.Push(5, 2, 8, 1, 9)

	runningMin := values.Iter().Scan(1000, func(minVal, val int) int {
		if val < minVal {
			return val
		}
		return minVal
	}).Collect(cmp.Cmp)

	Print("Running minimum: ")
	runningMin.Println()
}
