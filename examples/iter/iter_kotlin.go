package main

import (
	"strings"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

// https://kotlinlang.org/docs/basic-syntax.html#collections

func main() {
	fruits := SliceOf("banana", "avocado", "apple", "kiwifruit").Iter()
	fruits.
		// Filter(f.RxMatch[string](regexp.MustCompile("a"))).
		Filter(f.StartsWith("a")).
		SortBy(cmp.Cmp).
		Map(strings.ToUpper).
		ForEach(func(v string) { println(v) })
}
