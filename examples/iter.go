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
		// Filter(g.String.NotEmpty).
		Exclude(filters.IsZero).
		Map(g.String.Upper).
		Collect().
		Print()

	// ========================================================================

	s := g.SliceOf[g.String]("bbb", "ddd", "xxx", "aaa", "ccc").Iter().Enumerate()
	for next := s.Next(); next.IsSome(); next = s.Next() {
		fmt.Println(next.Some())
	}

	// ========================================================================

	g.NewMapOrd[int, string]().Set(88, "aa").Set(99, "bb").Set(199, "ii").
		Iter().
		Exclude(func(k int, v string) bool { return k == 99 }).
		Map(func(k int, v string) (int, string) { return k, v + "aaa" }).
		Collect().Print()

	// ========================================================================
}
