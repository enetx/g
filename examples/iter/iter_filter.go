package main

import (
	"github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	// Example 1: Filter integers in a slice and print the result
	g.SliceOf(1, 2).
		Iter().
		// Filter(func(i int) bool { return i != 1 }).
		// Filter(f.Ne(1)).
		Filter(f.Ne(1)).
		Collect().
		Print() // Slice[2]

	// Example 2: Chained filtering on a slice of strings and print the result
	fi := g.SliceOf("bbb", "ddd", "xxx", "aaa", "ccc").Iter()

	fi = fi.Filter(f.Ne("aaa"))
	fi = fi.Filter(f.Ne("xxx"))
	fi = fi.Filter(f.Ne("ddd"))
	fi = fi.Filter(f.Ne("bbb"))

	// fi = fi.Filter(func(s string) bool { return s != "aaa" })
	// fi = fi.Filter(func(s string) bool { return s != "xxx" })
	// fi = fi.Filter(func(s string) bool { return s != "ddd" })
	// fi = fi.Filter(func(s string) bool { return s != "bbb" })

	fi.Collect().Print() // Slice[ccc]

	// Example 3: Exclude a key from a map and print the result
	g.NewMap[int, string]().Set(88, "aa").Set(99, "bb").Set(199, "ii").
		Iter().
		Exclude(func(k int, _ string) bool { return k == 99 }).
		Collect().
		Print() // Map{88:aa, 199:ii}

	// Example 4: Exclude empty strings from a slice and print the result
	g.SliceOf[g.String]("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Exclude(f.Zero).
		Collect().
		Print() // Slice[bbb, ddd, aaa, ccc]
}
