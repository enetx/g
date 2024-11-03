package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

func main() {
	// Example 1: Sort an ordered map by key and print the result
	MapOrd[int, string]{
		{6, "bb"},
		{0, "dd"},
		{1, "aa"},
		{5, "xx"},
		{2, "cc"},
		{3, "ff"},
		{4, "zz"},
	}.
		Iter().
		// SortBy(
		// 	func(a, b Pair[int, string]) cmp.Ordering {
		// 		return cmp.Cmp(a.Key, b.Key)
		// 		// By value
		// 		// return cmp.Cmp(a.Value, b.Value)
		// 	}).
		// or
		SortByKey(cmp.Cmp). // By key
		// SortByValue(cmp.Cmp). // By value
		Collect().
		Print() // MapOrd{0:dd, 1:aa, 2:cc, 3:ff, 4:zz, 5:xx, 6:bb}

	// Example 2: Sort a slice of custom structs and print the result
	type status struct {
		date   time.Time
		name   String
		status String
	}

	s1 := status{time.Now(), "s1", "good"}
	s2 := status{time.Now(), "s2", "bad"}
	s3 := status{time.Now().Add(time.Second * 10), "s3", "bad"}

	SliceOf(s3, s1, s2).Iter().
		SortBy(
			func(a, b status) cmp.Ordering {
				var astatus Int = 5
				switch a.status {
				case "good":
					astatus = 0
				case "bad":
					astatus = 1
				}

				var bstatus Int = 5
				switch b.status {
				case "good":
					bstatus = 0
				case "bad":
					bstatus = 1
				}

				return astatus.Cmp(bstatus).Then(cmp.Cmp(a.date.Unix(), b.date.Unix()))
			}).
		Collect().
		Print()

	// Example 3: Sort a slice of time.Time, deduplicate, and print the result
	SliceOf(time.Now().Add(time.Second*20), time.Now()).
		Iter().
		SortBy(func(a, b time.Time) cmp.Ordering { return cmp.Cmp(a.Second(), b.Second()) }).
		Collect().
		Print()

	// Example 4: Sort and deduplicate a slice of integers and print the result
	SliceOf(9, 8, 9, 8, 0, 1, 1, 1, 2, 7, 2, 2, 2, 3, 4, 5).
		Iter().
		// SortBy(func(a, b int) cmp.Ordering { return cmp.Cmp(a, b) }). // or
		SortBy(cmp.Cmp).
		Dedup().
		Filter(f.Odd).
		Collect().
		Print() // Slice[1, 3, 5, 7, 9]

	// Example 5: Sort a slice of strings in descending order and print the result
	SliceOf("a", "c", "b").
		Iter().
		SortBy(cmp.Cmp).
		Collect().
		Print() // Slice[c, b, a]
}
