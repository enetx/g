package main

import (
	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/filters"
)

func main() {
	// Example 1: Filter integers in a slice and print the result
	g.SliceOf[g.Int](1, 2).
		Iter().
		Filter(func(i g.Int) bool { return i != 1 }).
		Collect().
		Print()

	// Example 2: Chained filtering on a slice of strings and print the result
	fi := g.SliceOf[g.String]("bbb", "ddd", "xxx", "aaa", "ccc").Iter().Filter(func(_ g.String) bool { return true })

	fi = fi.Filter(func(s g.String) bool { return s.Ne("aaa") })
	fi = fi.Filter(func(s g.String) bool { return s.Ne("xxx") })
	fi = fi.Filter(func(s g.String) bool { return s.Ne("ddd") })
	fi = fi.Filter(func(s g.String) bool { return s.Ne("bbb") })

	fi.Collect().Print()

	// Example 3: Exclude a key from a map and print the result
	g.NewMap[int, string]().Set(88, "aa").Set(99, "bb").Set(199, "ii").
		Iter().
		Exclude(func(k int, v string) bool { return k == 99 }).
		Collect().
		Print()

	// Example 4: Exclude empty strings from a slice and print the result
	g.SliceOf[g.String]("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Exclude(filters.IsZero).
		Collect().
		Print()
}
