package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	// Example 1: FilterMap for priority task filtering
	taskPriorities := NewHeap(cmp.Cmp[int])
	taskPriorities.Push(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	highPriorityTasks := taskPriorities.Iter().FilterMap(func(priority int) Option[int] {
		if priority >= 8 {
			return Some(priority * 100) // Convert to high priority codes
		}
		return None[int]()
	}).Collect(cmp.Cmp)

	Print("High priority codes: ")
	highPriorityTasks.Println()

	// Example 2: FilterMap for score grading system
	scores := NewHeap(func(a, b int) cmp.Ordering { return cmp.Cmp(b, a) }) // Max heap
	scores.Push(95, 87, 76, 92, 45, 88, 71, 34, 99)

	excellentGrades := scores.Iter().FilterMap(func(score int) Option[int] {
		if score >= 90 {
			return Some(score) // Keep excellent scores
		}
		return None[int]()
	}).Collect(cmp.Reverse)

	Print("Excellent grades: ")
	excellentGrades.Println()

	// Example 3: FilterMap for string validation
	words := NewHeap(cmp.Cmp[string])
	words.Push("hello", "", "world", "a", "golang", "short")

	longWords := words.Iter().FilterMap(func(word string) Option[string] {
		if len(word) > 3 { // Only keep longer words
			return Some(word)
		}
		return None[string]()
	}).Collect(cmp.Cmp)

	Print("Long words: ")
	longWords.Println()

	// Example 4: FilterMap for status code processing
	httpCodes := NewHeap(cmp.Cmp[int])
	httpCodes.Push(200, 404, 500, 201, 403, 301, 502, 204)

	errorCodes := httpCodes.Iter().FilterMap(func(code int) Option[int] {
		if code >= 400 { // Only error codes
			return Some(code)
		}
		return None[int]()
	}).Collect(cmp.Cmp)

	Print("Error codes: ")
	errorCodes.Println()

	// Example 5: FilterMap for temperature filtering
	temperatures := NewHeap(cmp.Cmp[int])
	temperatures.Push(-10, 0, 15, 25, 35, 42, -5, 30)

	hotTemps := temperatures.Iter().FilterMap(func(celsius int) Option[int] {
		if celsius > 30 { // Only hot temperatures
			return Some(celsius)
		}
		return None[int]()
	}).Collect(cmp.Cmp)

	Print("Hot temperatures: ")
	hotTemps.Println()

	// Example 6: FilterMap for admin users
	userIDs := NewHeap(cmp.Cmp[int])
	userIDs.Push(1001, 1002, 1003, 2001, 2002, 3001, 1004)

	adminUserIDs := userIDs.Iter().FilterMap(func(userID int) Option[int] {
		// Admin users have IDs starting with 1000-1999
		if userID >= 1000 && userID < 2000 {
			return Some(userID)
		}
		return None[int]()
	}).Collect(cmp.Cmp)

	Print("Admin user IDs: ")
	adminUserIDs.Println()

	// Example 7: FilterMap for large file sizes
	fileSizes := NewHeap(func(a, b int) cmp.Ordering { return cmp.Cmp(b, a) }) // Max heap
	fileSizes.Push(1024, 2048, 512, 4096, 256, 8192, 128)

	largeSizes := fileSizes.Iter().FilterMap(func(sizeInKB int) Option[int] {
		if sizeInKB >= 2048 { // Files >= 2MB
			return Some(sizeInKB)
		}
		return None[int]()
	}).Collect(cmp.Reverse)

	Print("Large file sizes: ")
	largeSizes.Println()
}
