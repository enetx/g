package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Example 1: Basic number manipulation
	m1 := NewMapSafe[string, Int]()

	em1 := m1.Entry("root") // Get entry for key "root"
	em1.OrDefault()         // Insert zero value (0) if the key is missing

	em1.Transform(
		func(i Int) Int {
			return i + 10 // Atomically add 10 to the current value
		})

	em1.Get().Some().Println() // Print the final value: 10

	// Example 2: Appending to slices in a loop
	m2 := NewMapSafe[int, Slice[int]]()

	for i := range 5 {
		em2 := m2.Entry(i)
		em2.OrDefault() // Insert an empty slice if missing
		em2.Transform(
			func(s Slice[int]) Slice[int] {
				return s.Append(i) // Append the current index to the slice
			})
	}

	for i := range 10 {
		em2 := m2.Entry(i)
		em2.OrDefault() // Ensure the key exists
		em2.Transform(
			func(s Slice[int]) Slice[int] {
				return s.Append(i) // Append again for overlapping keys (0–4)
			})
	}

	// Final map: keys 0–4 contain two values, 5–9 contain one
	m2.Println() // Output: Map{0:[0 0] 1:[1 1] ... 4:[4 4] 5:[5] ... 9:[9]}

	// Example 3: Lazy initialization using OrSetBy
	m3 := NewMapSafe[string, Slice[string]]()

	em3 := m3.Entry("users")

	em3.OrSetBy(
		func() Slice[string] {
			fmt.Println("initializing users slice")
			return Slice[string]{"alice", "bob"} // Only called if the key is missing
		})

	em3.Transform(
		func(s Slice[string]) Slice[string] {
			return s.Append("charlie") // Append "charlie" to the existing slice
		})

	fmt.Println("m3:", m3) // Output: Map{users:[alice bob charlie]}

	// Example 4: Manual delete after update
	m4 := NewMapSafe[string, Int]()

	em4 := m4.Entry("count")
	em4.OrSet(10) // Insert 10 if missing

	em4.Set(100)                      // Update to 100
	fmt.Println("before delete:", m4) // Map{count:100}

	em4.Delete()                     // Remove the key
	fmt.Println("after delete:", m4) // Map{}

	// Example 5: Combined operations
	m5 := NewMapSafe[string, Int]()

	em5 := m5.Entry("a")
	em5.OrDefault() // Insert zero (0)

	em5.Transform(
		func(i Int) Int {
			return i + 5 // Increase to 5
		})

	em5.Set(42) // Overwrite with 42

	fmt.Println("m5:", m5) // Map{a:42}
}
