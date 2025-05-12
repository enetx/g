package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Example 1: basic OrSet, Transform, Get, Println
	m := NewMap[string, Int]()
	m.Entry("root").
		OrSet(1).
		Transform(func(i *Int) { *i++ })
	m.Entry("root").Get().Some().Println() // prints: 2

	// Example 2: accumulating slices per key
	m2 := NewMap[int, Slice[int]]()

	for i := range 5 {
		m2.Entry(i).
			OrDefault().
			Transform(func(sl *Slice[int]) { sl.Push(i) })
	}

	for i := range 10 {
		m2.Entry(i).
			OrDefault().
			Transform(func(sl *Slice[int]) { sl.Push(i) })
	}

	m2.Println() // prints: Map{0:[0 0] 1:[1 1] ... 4:[4 4] 5:[5] ... 9:[9]}

	// Example 3: lazy initialization with OrSetBy
	m3 := NewMap[string, Slice[string]]()
	m3.Entry("users").
		OrSetBy(func() Slice[string] {
			fmt.Println("initializing users slice")
			return Slice[string]{"alice", "bob"}
		}).
		Transform(func(sl *Slice[string]) { sl.Push("charlie") })
	fmt.Println("m3:", m3)
	// Output:
	// initializing users slice
	// m3: Map{users:Slice[alice, bob, charlie]}

	// Example 4: override with Set and then Delete
	m4 := NewMap[string, Int]()
	m4.Entry("count").OrSet(10)
	m4.Entry("count").Set(100)
	fmt.Println("before delete:", m4) // Map{count:100}
	m4.Entry("count").Delete()
	fmt.Println("after delete:", m4) // Map{}

	// Example 5: chaining OrDefault, Transform, Set
	m5 := NewMap[string, Int]()
	m5.Entry("a").
		OrDefault().
		Transform(func(i *Int) { *i += 5 }).
		Set(42)
	fmt.Println("m5:", m5) // Map{a:42}
}
