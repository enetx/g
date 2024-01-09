package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/filters"
)

func main() {
	// ========================================================================
	// retrieve the filter iterator and reuse it
	fi := g.SliceOf[g.String]("bbb", "ddd", "xxx", "aaa", "ccc").Iter().Filter(func(_ g.String) bool { return true })

	fi = fi.Filter(func(s g.String) bool { return s.Ne("aaa") })
	fi = fi.Filter(func(s g.String) bool { return s.Ne("xxx") })
	fi = fi.Filter(func(s g.String) bool { return s.Ne("ddd") })
	fi = fi.Filter(func(s g.String) bool { return s.Ne("bbb") })

	fi.Collect().Print()

	// ========================================================================

	g.NewMap[int, string]().Set(88, "aa").Set(99, "bb").Set(199, "ii").
		Iter().
		Exclude(func(k int, v string) bool { return k == 99 }).
		Map(func(k int, v string) (int, string) { return k, v + "aaa" }).
		Collect().
		Print()

	// ========================================================================

	g.SliceOf[g.String]("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Cycle().
		Take(20).
		Exclude(filters.IsZero).
		Map(g.String.Upper).
		Collect().
		Print()

	// ========================================================================

	pairs := g.SliceOf[g.String]("bbb", "ddd", "xxx", "aaa", "ccc").
		Iter().
		Enumerate().
		Collect()

	ps := g.MapOrd[uint, g.String](pairs)
	ps.Print()

	// ========================================================================

	g.NewMapOrd[int, string]().Set(88, "aa").Set(99, "bb").Set(199, "ii").
		Iter().
		Exclude(func(k int, v string) bool { return k == 99 }).
		Map(func(k int, v string) (int, string) { return k, v + "aaa" }).
		Collect().Print()

	// ========================================================================

	slice1 := g.Slice[int]{1, 2, 3}.Iter()
	slice2 := g.Slice[int]{4, 5, 6}.Iter()
	slice3 := g.Slice[int]{7, 8, 9}.Iter()

	zipped := slice1.Zip(slice2, slice3).Collect()

	for _, v := range zipped {
		v.Print()
	}

	// ========================================================================

	p1 := g.SliceOf[g.String]("bbb", "ddd")
	p2 := g.SliceOf[g.String]("xxx", "aaa")

	pp := p1.
		Iter().
		Chain(p2.Iter()).
		Map(g.String.Upper).
		Permutations().
		Collect()

	for _, v := range pp {
		v.Print()
	}

	// ========================================================================

	g.SliceOf[g.String]("bbb", "ddd", "bbb", "aaa", "bbb").
		Iter().
		Unique().
		Map(g.String.Upper).
		Collect().
		Print()

	// ========================================================================

	chunks := g.SliceOf[g.String]("bbb", "ddd", "bbb", "ccc", "aaa", "bbb", "ccc").
		Iter().
		Unique().
		Chunks(2).
		Collect()

	fmt.Println(chunks)
}
