package main

import (
	. "github.com/enetx/g"
)

func main() {
	fruits := SliceOf[String]("banana", "avocado", "apple", "kiwifruit").Iter()

	fruits.
		Filter(func(s String) bool { return s.StartsWith("a") }).
		SortBy(String.Cmp).
		Map(String.Upper).
		ForEach(func(v String) { v.Print() })
}
