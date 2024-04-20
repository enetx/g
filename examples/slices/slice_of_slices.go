package main

import (
	"fmt"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	ns1 := g.NewSlice[g.String]().Append("aaa")
	ns2 := g.NewSlice[g.String]().Append("bbb", "ccc")
	ns3 := g.NewSlice[g.String]().Append("ccc", "dddd", "wwwww")

	nx := g.SliceOf(ns3, ns2, ns1, ns2)

	nx.SortBy(func(a, b g.Slice[g.String]) cmp.Ordering {
		if a.Eq(b) {
			return cmp.Equal
		}

		return cmp.Less
	})

	nx.Print()

	nx = nx.Iter().Dedup().Collect().Print()

	nx.SortBy(func(a, b g.Slice[g.String]) cmp.Ordering { return b.Get(0).Cmp(a.Get(0)) })
	nx.Print()

	nx.SortBy(func(a, b g.Slice[g.String]) cmp.Ordering { return a.Len().Cmp(b.Len()) })
	nx.Print()

	nx.Reverse()
	nx.Print()

	nx.Random().Print()
	nx.RandomSample(2).Print()

	ch := nx.Iter().Chunks(2).Collect() // return []Slice[T]
	chunks := g.SliceOf(ch...)          // make slice chunks
	chunks.Print()

	// pr := nx.Iter().Permutations().Collect() // return []Slice[T]
	// permutations := g.SliceOf(pr...)         // make slice permutations
	// permutations.Print()

	m := g.NewMap[string, g.Slice[g.Slice[g.String]]]()
	m.Set("one", nx)

	fmt.Println(m.Get("one").Some().Last().Contains("aaa"))

	nested := g.Slice[any]{1, 2, g.Slice[int]{3, 4, 5}, []any{6, 7, []int{8, 9}}}
	flattened := nested.Iter().Flatten().Collect()
	fmt.Println(flattened)

	nestedSlice := g.Slice[any]{
		1,
		g.SliceOf(2, 3),
		"abc",
		g.SliceOf("def", "ghi"),
		g.SliceOf(4.5, 6.7),
	}

	nestedSlice.Print()                            // Output: Slice[1, Slice[2, 3], abc, Slice[def, ghi], Slice[4.5, 6.7]]
	nestedSlice.Iter().Flatten().Collect().Print() // Output: Slice[1, 2, 3, abc, def, ghi, 4.5, 6.7]

	nestedSlice2 := g.Slice[any]{
		1,
		g.SliceOf(2, 3),
		"abc",
		g.SliceOf("awe", "som", "e"),
		g.SliceOf("co", "ol"),
		g.SliceOf(4.5, 6.7),
		map[string]string{"a": "ss"},
		g.SliceOf(g.MapOrd[int, int]{{1, 1}}, g.MapOrd[int, int]{{2, 2}}),
	}

	nestedSlice2.Iter().Flatten().Collect().Print()
}
