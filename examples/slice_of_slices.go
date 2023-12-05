package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/dbg"
)

func main() {
	ns1 := g.NewSlice[g.String]().Append("aaa")
	ns2 := g.NewSlice[g.String]().Append("bbb").Append("ccc")
	ns3 := g.NewSlice[g.String]().Append("ccc").Append("dddd").Append("wwwww")

	nx := g.SliceOf(ns3, ns2, ns1, ns2)

	nx.Flatten().Print()
	nx.Flatten().Last().(g.String).Upper().Print()

	nx = nx.Unique()

	nx.SortBy(func(i, j int) bool { return nx[i].Get(0).Lt(nx[j].Get(0)) }).Print()
	nx.SortBy(func(i, j int) bool { return nx[i].Len() < nx[j].Len() }).Print()

	nx.Reverse().Print()

	nx.Random().Print()
	nx.RandomSample(2).Print()

	ch := nx.Chunks(2)         // return []Slice[T]
	chunks := g.SliceOf(ch...) // make slice batches
	chunks.Print()

	pr := nx.Permutations()          // return []Slice[T]
	permutations := g.SliceOf(pr...) // make slice permutations
	permutations.Print()

	m := g.NewMap[string, g.Slice[g.Slice[g.String]]]()
	m.Set("one", nx)

	fmt.Println(m.Get("one").Last().Contains("aaa"))

	nestedSlice := g.Slice[any]{
		1,
		g.SliceOf(2, 3),
		"abc",
		g.SliceOf("def", "ghi"),
		g.SliceOf(4.5, 6.7),
	}

	nestedSlice.Print()           // Output: Slice[1, Slice[2, 3], abc, Slice[def, ghi], Slice[4.5, 6.7]]
	nestedSlice.Flatten().Print() // Output: Slice[1, 2, 3, abc, def, ghi, 4.5, 6.7]

	nestedSlice2 := g.Slice[any]{
		1,
		[]int{2, 3},
		"abc",
		g.SliceOf("awe", "som", "e"),
		[]string{"co", "ol"},
		g.SliceOf(4.5, 6.7),
		[]float64{4.5, 6.7},
		map[string]string{"a": "ss"},
		g.SliceOf(g.NewMapOrd[int, int]().Set(1, 1), g.NewMapOrd[int, int]().Set(2, 2)),
	}

	// Output: Slice[1, [2 3], abc, Slice[awe, som, e], [co ol], Slice[4.5, 6.7], [4.5 6.7], map[a:ss], Slice[MapOrd{1:1}, MapOrd{2:2}]]
	nestedSlice2.Print()

	// Output: Slice[1, 2, 3, abc, awe, som, e, lol, ov, 4.5, 6.7, 4.5, 6.7, map[a:ss], MapOrd{1:1}, MapOrd{2:2}]
	dbg.Dbg(nestedSlice2.Flatten())
}
