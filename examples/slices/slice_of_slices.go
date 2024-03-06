package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/dbg"
)

func main() {
	ns1 := g.NewSlice[g.String]().Append("aaa")
	ns2 := g.NewSlice[g.String]().Append("bbb", "ccc")
	ns3 := g.NewSlice[g.String]().Append("ccc", "dddd", "wwwww")

	nx := g.SliceOf(ns3, ns2, ns1, ns2)

	nx.SortBy(func(a, b g.Slice[g.String]) bool {
		if a.Eq(b) {
			return false
		}

		return true
	})

	nx.Print()

	nx = nx.Iter().Dedup().Collect().Print()

	nx.SortBy(func(a, b g.Slice[g.String]) bool { return a.Get(0).Lt(b.Get(0)) })
	nx.Print()

	nx.SortBy(func(a, b g.Slice[g.String]) bool { return a.Len() < b.Len() })
	nx.Print()

	nx.Reverse()
	nx.Print()

	nx.Random().Print()
	nx.RandomSample(2).Print()

	// ch := nx.Iter().Chunks(2).Collect() // return []Slice[T]
	// chunks := g.SliceOf(ch...)          // make slice chunks
	// chunks.Print()

	// pr := nx.Iter().Permutations().Collect() // return []Slice[T]
	// permutations := g.SliceOf(pr...)         // make slice permutations
	// permutations.Print()

	m := g.NewMap[string, g.Slice[g.Slice[g.String]]]()
	m.Set("one", nx)

	fmt.Println(m.Get("one").Last().Contains("aaa"))

	nestedSlice := g.Slice[any]{
		1,
		g.SliceOf[any](2, 3),
		"abc",
		g.SliceOf[any]("def", "ghi"),
		g.SliceOf[any](4.5, 6.7),
	}

	nestedSlice.Print()                            // Output: Slice[1, Slice[2, 3], abc, Slice[def, ghi], Slice[4.5, 6.7]]
	nestedSlice.Iter().Flatten().Collect().Print() // Output: Slice[1, 2, 3, abc, def, ghi, 4.5, 6.7]

	nestedSlice2 := g.Slice[any]{
		1,
		g.Slice[any]{2, 3},
		"abc",
		g.SliceOf[any]("awe", "som", "e"),
		g.Slice[any]{"co", "ol"},
		g.SliceOf[any](4.5, 6.7),
		g.Slice[any]{4.5, 6.7},
		map[string]string{"a": "ss"},
		g.SliceOf[any](g.MapOrd[int, int]{{1, 1}}, g.MapOrd[int, int]{{2, 2}}),
	}

	dbg.Dbg(nestedSlice2.Iter().Flatten().Collect())
}