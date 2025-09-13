package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	Println("=== Parallel MaxBy/MinBy Examples ===\n")

	// Example 1: Parallel MaxBy/MinBy with Slices - basic comparison
	Println("1. Parallel MaxBy/MinBy with Slices - Numbers:")
	start := time.Now()

	numbers := SliceOf(42, 7, 93, 15, 88, 3, 67, 21, 99, 56)

	maxResult := numbers.Iter().
		Parallel(3).
		Inspect(func(int) {
			// Simulate processing time
			time.Sleep(10 * time.Millisecond)
		}).
		MaxBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})

	minResult := numbers.Iter().
		Parallel(3).
		Inspect(func(int) {
			time.Sleep(10 * time.Millisecond)
		}).
		MinBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})

	duration1 := time.Since(start)
	Println("Numbers: {}", numbers)
	Println("Maximum: {}", maxResult)
	Println("Minimum: {}", minResult)
	Println("Duration: {} (should be ~35ms with parallelism)\n", duration1)

	// Example 2: Parallel MaxBy/MinBy with Deques - string comparison
	Println("2. Parallel MaxBy/MinBy with Deques - Strings by length:")
	deque := NewDeque[String]()
	words := SliceOf[String]("cat", "elephant", "dog", "butterfly", "ant", "hippopotamus", "bee", "tiger")
	words.Iter().ForEach(func(word String) {
		deque.PushBack(word)
	})

	start = time.Now()
	longestWord := deque.Iter().
		Parallel(3).
		Inspect(func(String) {
			// Simulate string analysis
			time.Sleep(15 * time.Millisecond)
		}).
		MaxBy(func(a, b String) cmp.Ordering {
			return cmp.Cmp(a.Len(), b.Len())
		})

	shortestWord := deque.Iter().
		Parallel(3).
		Inspect(func(String) {
			time.Sleep(15 * time.Millisecond)
		}).
		MinBy(func(a, b String) cmp.Ordering {
			return cmp.Cmp(a.Len(), b.Len())
		})

	duration2 := time.Since(start)
	Println("Words: {}", words)
	Println("Longest word: {} (length: {})", longestWord, longestWord.Some().Len())
	Println("Shortest word: {} (length: {})", shortestWord, shortestWord.Some().Len())
	Println("Duration: {} (should be ~45ms with parallelism)\n", duration2)

	// Example 3: Parallel MaxBy/MinBy with Heaps - custom comparison
	Println("3. Parallel MaxBy/MinBy with Heaps - Points by distance:")
	type Point struct {
		X, Y float64
		Name String
	}

	distanceFromOrigin := func(p Point) float64 {
		return (p.X*p.X + p.Y*p.Y) // squared distance (avoiding sqrt for performance)
	}

	heap := NewHeap(func(a, b Point) cmp.Ordering {
		return cmp.Cmp(a.Name, b.Name) // Heap ordered by name
	})

	points := SliceOf(
		Point{3, 4, "A"},  // distance² = 25
		Point{1, 1, "B"},  // distance² = 2
		Point{5, 12, "C"}, // distance² = 169
		Point{0, 1, "D"},  // distance² = 1
		Point{8, 6, "E"},  // distance² = 100
		Point{2, 3, "F"},  // distance² = 13
	)
	points.Iter().ForEach(func(p Point) {
		heap.Push(p)
	})

	start = time.Now()
	farthestPoint := heap.Iter().
		Parallel(3).
		Inspect(func(Point) {
			// Simulate complex calculation
			time.Sleep(20 * time.Millisecond)
		}).
		MaxBy(func(a, b Point) cmp.Ordering {
			distA := distanceFromOrigin(a)
			distB := distanceFromOrigin(b)
			return cmp.Cmp(distA, distB)
		})

	nearestPoint := heap.Iter().
		Parallel(3).
		Inspect(func(Point) {
			time.Sleep(20 * time.Millisecond)
		}).
		MinBy(func(a, b Point) cmp.Ordering {
			distA := distanceFromOrigin(a)
			distB := distanceFromOrigin(b)
			return cmp.Cmp(distA, distB)
		})

	duration3 := time.Since(start)
	Println("Points: {}", points)
	if farthestPoint.IsSome() {
		fp := farthestPoint.Some()
		Println("Farthest point: {} at ({}, {}) - distance²: {}",
			fp.Name, fp.X, fp.Y, distanceFromOrigin(fp))
	}
	if nearestPoint.IsSome() {
		np := nearestPoint.Some()
		Println("Nearest point: {} at ({}, {}) - distance²: {}",
			np.Name, np.X, np.Y, distanceFromOrigin(np))
	}
	Println("Duration: {} (should be ~40ms with parallelism)\n", duration3)

	// Example 4: Complex data structures - finding best/worst performance
	Println("4. Complex data structures - Employee performance:")
	type Employee struct {
		Name       String
		Department String
		Score      float64
		Years      int
	}

	employees := SliceOf(
		Employee{"Alice", "Engineering", 95.5, 5},
		Employee{"Bob", "Sales", 87.2, 3},
		Employee{"Charlie", "Engineering", 92.1, 7},
		Employee{"Diana", "Marketing", 89.8, 4},
		Employee{"Eve", "Sales", 94.3, 6},
		Employee{"Frank", "Engineering", 88.7, 2},
	)

	// Custom scoring: performance score weighted by experience
	complexScore := func(e Employee) float64 {
		return e.Score + float64(e.Years)*2.5 // Bonus points for experience
	}

	start = time.Now()
	topPerformer := employees.Iter().
		Parallel(3).
		Inspect(func(Employee) {
			// Simulate performance analysis
			time.Sleep(25 * time.Millisecond)
		}).
		MaxBy(func(a, b Employee) cmp.Ordering {
			return cmp.Cmp(complexScore(a), complexScore(b))
		})

	lowestPerformer := employees.Iter().
		Parallel(3).
		Inspect(func(Employee) {
			time.Sleep(25 * time.Millisecond)
		}).
		MinBy(func(a, b Employee) cmp.Ordering {
			return cmp.Cmp(complexScore(a), complexScore(b))
		})

	duration4 := time.Since(start)
	if topPerformer.IsSome() {
		tp := topPerformer.Some()
		Println("Top performer: {} ({}) - Score: {}, Experience: {} years, Weighted: {}",
			tp.Name, tp.Department, tp.Score, tp.Years, complexScore(tp))
	}
	if lowestPerformer.IsSome() {
		lp := lowestPerformer.Some()
		Println("Needs improvement: {} ({}) - Score: {}, Experience: {} years, Weighted: {}",
			lp.Name, lp.Department, lp.Score, lp.Years, complexScore(lp))
	}
	Println("Duration: {} (should be ~50ms with parallelism)\n", duration4)

	// Example 5: Edge cases
	Println("5. Edge cases:")

	// Empty collection
	emptySlice := NewSlice[int]()
	emptyMax := emptySlice.Iter().
		Parallel(2).
		MaxBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})
	Println("Empty collection MaxBy: {}", emptyMax)

	// Single element
	singleMax := SliceOf(42).Iter().
		Parallel(2).
		MaxBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})
	Println("Single element MaxBy: {}", singleMax)

	// All equal elements
	equalElements := SliceOf(5, 5, 5, 5)
	equalMax := equalElements.Iter().
		Parallel(2).
		MaxBy(func(a, b int) cmp.Ordering {
			return cmp.Cmp(a, b)
		})
	Println("All equal elements MaxBy: {}", equalMax)
	Println("")

	// Example 6: Performance comparison - Sequential vs Parallel
	Println("6. Performance Comparison:")
	largeDataset := Range(1, 1001).Collect() // 1000 items

	// Sequential processing
	start = time.Now()
	seqMax := largeDataset.Iter().
		Inspect(func(int) {
			time.Sleep(1 * time.Millisecond)
		}).
		MaxBy(cmp.Cmp)
	seqDuration := time.Since(start)

	// Parallel processing
	start = time.Now()
	parMax := largeDataset.Iter().
		Parallel(20).
		Inspect(func(int) {
			time.Sleep(1 * time.Millisecond)
		}).
		MaxBy(cmp.Cmp)
	parDuration := time.Since(start)

	Println("Sequential duration: {} (should be ~1000ms)", seqDuration)
	Println("Parallel duration: {} (should be ~50ms)", parDuration)
	Println("Speedup: {}x", Float(seqDuration.Milliseconds())/Float(parDuration.Milliseconds()))
	seqSome := seqMax.IsSome() && parMax.IsSome() && seqMax.Some() == parMax.Some()
	Println("Results equal: {}", seqSome)
	Println("Maximum found: {}", seqMax)
}
