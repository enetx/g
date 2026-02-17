package main

import (
	"fmt"
	"regexp"

	. "github.com/enetx/g"
)

func main() {
	// Regexp
	reCache := NewMapSafe[String, *regexp.Regexp]()

	p := String("content=.+")

	re := reCache.Entry(p).OrInsertWith(p.Regexp().Compile().Unwrap)
	re = reCache.Entry(p).OrInsertWith((p + "\\w").Regexp().Compile().Unwrap)

	fmt.Println("re:", re) // re: content=.+

	// Counter
	m := Map[string, Int]{}

	m.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	m.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	m.Entry("counter").AndModify(func(v *Int) { *v++ }).OrInsert(1)
	fmt.Println("counter:", m.Get("counter").Some()) // 3

	// Word frequency
	words := SliceOf("apple", "banana", "apple", "cherry", "banana", "apple")
	freq := Map[string, Int]{}

	words.Iter().ForEach(func(word string) {
		// freq[word]++
		freq.Entry(word).AndModify(func(v *Int) { *v++ }).OrInsert(1)
	})

	fmt.Println("freq:", freq) // Map{apple:3, banana:2, cherry:1}

	// Get value after insert
	val := m.Entry("new_key").OrInsert(100)
	fmt.Println("val:", val) // 100

	// Lazy init
	config := Map[string, String]{}
	host := config.Entry("host").OrInsertWith(func() String {
		fmt.Println("computing default host...")
		return "localhost"
	})

	fmt.Println("host:", host)

	// Pattern matching
	switch e := m.Entry("counter").(type) {
	case OccupiedEntry[string, Int]:
		fmt.Println("exists:", e.Get())
		// e.Remove()
	case VacantEntry[string, Int]:
		e.Insert(42)
	}

	// Grouping
	groups := Map[int, Slice[int]]{}
	for i := range 10 {
		groups.Entry(i % 3).
			AndModify(func(s *Slice[int]) { *s = s.Append(i) }).
			OrInsertWith(func() Slice[int] { return Slice[int]{i} })
	}

	fmt.Println("groups:", groups)
}
