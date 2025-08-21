package main

import (
	"math"
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	Println("=== Parallel FilterMap Examples ===\n")

	// Example 1: Parallel FilterMap with Slices - number processing
	Println("1. Parallel FilterMap with Slices - Safe division:")
	start := time.Now()

	numbers := SliceOf(10, 0, 15, -5, 20, 3, 0, 8)
	result1 := numbers.Iter().
		Parallel(4).
		FilterMap(func(n int) Option[int] {
			// Simulate processing time
			time.Sleep(30 * time.Millisecond)
			// Only keep positive numbers, double them
			if n > 0 {
				return Some(n * 2)
			}
			return None[int]()
		}).
		Collect()

	duration1 := time.Since(start)
	Println("Result: {}", result1)
	Println("Duration: {} (should be ~60ms with parallelism)\n", duration1)

	// Example 2: Parallel FilterMap with Deques - email validation
	Println("2. Parallel FilterMap with Deques - Email processing:")
	deque := NewDeque[String]()
	emails := SliceOf[String](
		"user@example.com",
		"invalid-email",
		"admin@test.org",
		"",
		"support@company.net",
		"no-at-sign",
		"another@valid.com",
	)
	emails.Iter().ForEach(func(email String) {
		deque.PushBack(email)
	})

	start = time.Now()
	result2 := deque.Iter().
		Parallel(3).
		FilterMap(func(email String) Option[String] {
			// Simulate validation time
			time.Sleep(25 * time.Millisecond)
			// Extract domain from valid emails
			if email.Contains("@") && !email.Empty() {
				parts := email.Split("@").Collect()
				if parts.Len() == 2 && !parts[1].Empty() {
					return Some(parts[1].Upper())
				}
			}
			return None[String]()
		}).
		Collect()

	duration2 := time.Since(start)
	Println("Valid domains: {}", result2.String())
	Println("Duration: {} (should be ~50-75ms with parallelism)\n", duration2)

	// Example 3: Parallel FilterMap with Heaps - score processing
	Println("3. Parallel FilterMap with Heaps - Score grading:")
	heap := NewHeap(cmp.Cmp[int])
	scores := SliceOf(85, 42, 95, 67, 23, 89, 91, 55, 78)
	scores.Iter().ForEach(func(score int) {
		heap.Push(score)
	})

	start = time.Now()
	result3 := heap.Iter().
		Parallel(3).
		FilterMap(func(score int) Option[int] {
			// Simulate grading time
			time.Sleep(20 * time.Millisecond)
			// Only return passing scores (≥70), scaled 0-100
			if score >= 90 {
				return Some(100) // A grade
			} else if score >= 80 {
				return Some(85) // B grade
			} else if score >= 70 {
				return Some(75) // C grade
			}
			// Failing grades filtered out
			return None[int]()
		}).
		Collect()

	duration3 := time.Since(start)
	Println("Passing grades: {} elements", result3.Len())
	passingGrades := make([]int, 0)
	for !result3.Empty() {
		passingGrades = append(passingGrades, result3.Pop().Some())
	}
	sorted := SliceOf(passingGrades...)
	sorted.SortBy(cmp.Cmp)
	sorted.Println()
	Println("Duration: {} (should be ~60ms with parallelism)\n", duration3)

	// Example 4: Complex data transformation - User data
	Println("4. Complex data transformation - User data:")
	type User struct {
		ID   int
		Name String
		Age  int
	}

	users := SliceOf(
		User{1, "Alice", 25},
		User{2, "Bob", 17}, // Minor - filtered out
		User{3, "Charlie", 30},
		User{4, "Diana", 16}, // Minor - filtered out
		User{5, "Eve", 28},
		User{6, "Frank", 19},
	)

	start = time.Now()
	result4 := users.Iter().
		Parallel(3).
		FilterMap(func(user User) Option[User] {
			// Simulate user processing
			time.Sleep(15 * time.Millisecond)
			// Only process adult users, increment age
			if user.Age >= 18 {
				user.Age += 1 // Age them by 1 year
				return Some(user)
			}
			return None[User]()
		}).
		Collect()

	duration4 := time.Since(start)
	Println("Adult users processed:")
	result4.Iter().ForEach(func(user User) {
		Println("  {} (age {})", user.Name, user.Age)
	})
	Println("Duration: {} (should be ~30ms with parallelism)\n", duration4)

	// Example 5: Numerical computation - Square roots of perfect squares
	Println("5. Mathematical filtering - Perfect squares:")
	numbers2 := Range(1, 51).Collect() // 1 to 50

	start = time.Now()
	result5 := numbers2.Iter().
		Parallel(5).
		FilterMap(func(n int) Option[int] {
			// Simulate computation
			time.Sleep(5 * time.Millisecond)
			// Check if n is a perfect square
			sqrt := int(math.Sqrt(float64(n)))
			if sqrt*sqrt == n {
				return Some(sqrt)
			}
			return None[int]()
		}).
		Collect()
	result5.SortBy(cmp.Cmp)

	duration5 := time.Since(start)
	Println("Square roots of perfect squares (1-50): {}", result5)
	Println("Duration: {} (should be ~50ms with parallelism)\n", duration5)

	// Example 6: Performance comparison - Sequential vs Parallel
	Println("6. Performance Comparison:")
	largeDataset := Range(1, 101).Collect() // 100 items

	// Sequential processing
	start = time.Now()
	seqResult := largeDataset.Iter().
		FilterMap(func(n int) Option[int] {
			time.Sleep(2 * time.Millisecond)
			if n%3 == 0 && n > 10 {
				return Some(n * 2)
			}
			return None[int]()
		}).
		Collect()
	seqDuration := time.Since(start)

	// Parallel processing
	start = time.Now()
	parResult := largeDataset.Iter().
		Parallel(10).
		FilterMap(func(n int) Option[int] {
			time.Sleep(2 * time.Millisecond)
			if n%3 == 0 && n > 10 {
				return Some(n * 2)
			}
			return None[int]()
		}).
		Collect()
	parDuration := time.Since(start)

	Println("Sequential duration: {} (should be ~200ms)", seqDuration)
	Println("Parallel duration: {} (should be ~20-30ms)", parDuration)
	Println("Speedup: {}x", Float(seqDuration.Milliseconds())/Float(parDuration.Milliseconds()))
	seqResult.SortBy(cmp.Cmp)
	parResult.SortBy(cmp.Cmp)
	Println("Results equal: {}", seqResult.Eq(parResult))
	Println("Filtered count: {}", seqResult.Len())

	/* Expected Output:
	=== Parallel FilterMap Examples ===

	1. Parallel FilterMap with Slices - Safe division:
	Result: Slice[20, 30, 40, 6, 16]
	Duration: 62ms (should be ~60ms with parallelism)

	2. Parallel FilterMap with Deques - Email processing:
	Valid domains: Deque[EXAMPLE.COM, TEST.ORG, COMPANY.NET, VALID.COM]
	Duration: 58ms (should be ~50-75ms with parallelism)

	3. Parallel FilterMap with Heaps - Score grading:
	Passing grades: 5 elements
	Slice[75, 85, 85, 100, 100]
	Duration: 61ms (should be ~60ms with parallelism)

	4. Complex data transformation - User data:
	Adult users processed:
	  Alice (age 26)
	  Charlie (age 31)
	  Eve (age 29)
	  Frank (age 20)
	Duration: 31ms (should be ~30ms with parallelism)

	5. Mathematical filtering - Perfect squares:
	Square roots of perfect squares (1-50): Slice[1, 2, 3, 4, 5, 6, 7]
	Duration: 51ms (should be ~50ms with parallelism)

	6. Performance Comparison:
	Sequential duration: 201ms (should be ~200ms)
	Parallel duration: 22ms (should be ~20-30ms)
	Speedup: 9.1x
	Results equal: true
	Filtered count: 30
	*/
}
