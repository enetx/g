package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	ns1 := Slice[String]{"aaa"}
	ns2 := Slice[String]{"bbb", "ccc"}
	ns3 := Slice[String]{"ccc", "dddd", "wwwww"}

	nx := SliceOf(ns3, ns2, ns1, ns2)

	nx.SortBy(func(a, b Slice[String]) cmp.Ordering {
		if a.Eq(b) {
			return cmp.Equal
		}

		return cmp.Less
	})

	nx.Println()

	nx = nx.Iter().Dedup().Collect().Println()

	nx.SortBy(func(a, b Slice[String]) cmp.Ordering { return b.Get(0).Cmp(a.Get(0)) })
	nx.Println()

	nx.SortBy(func(a, b Slice[String]) cmp.Ordering { return a.Len().Cmp(b.Len()) })
	nx.Println()

	nx.Reverse()
	nx.Println()

	nx.Random().Println()
	nx.RandomSample(2).Println()

	ch := nx.Iter().Chunks(2).Collect() // return []Slice[T]
	chunks := SliceOf(ch...)            // make slice chunks
	chunks.Println()

	// pr := nx.Iter().Permutations().Collect() // return []Slice[T]
	// permutations := g.SliceOf(pr...)         // make slice permutations
	// permutations.Println()

	m := NewMap[string, Slice[Slice[String]]]()
	m.Set("one", nx)

	fmt.Println(m.Get("one").Some().Last().Contains("aaa"))

	nested := Slice[any]{1, 2, Slice[int]{3, 4, 5}, []any{6, 7, []int{8, 9}}}
	flattened := nested.Iter().Flatten().Collect()
	fmt.Println(flattened)

	nestedSlice := Slice[any]{
		1,
		SliceOf(2, 3),
		"abc",
		SliceOf("def", "ghi"),
		SliceOf(4.5, 6.7),
	}

	nestedSlice.Println()                            // Output: Slice[1, Slice[2, 3], abc, Slice[def, ghi], Slice[4.5, 6.7]]
	nestedSlice.Iter().Flatten().Collect().Println() // Output: Slice[1, 2, 3, abc, def, ghi, 4.5, 6.7]

	nestedSlice2 := Slice[any]{
		1,
		SliceOf(2, 3),
		"abc",
		SliceOf("awe", "som", "e"),
		SliceOf("co", "ol"),
		SliceOf(4.5, 6.7),
		map[string]string{"a": "ss"},
		SliceOf(MapOrd[int, int]{{1, 1}}, MapOrd[int, int]{{2, 2}}),
	}

	// Slice[1, 2, 3, abc, awe, som, e, co, ol, 4.5, 6.7, map[a:ss], {1 1}, {2 2}]
	nestedSlice2.Iter().Flatten().Collect().Println()
}
