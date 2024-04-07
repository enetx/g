package main

import (
	"time"

	"github.com/enetx/g"
	"github.com/enetx/g/filters"
)

func main() {
	// Example 1: Sort an ordered map by key and print the result
	g.MapOrd[int, string]{
		{6, "bb"},
		{0, "dd"},
		{1, "aa"},
		{5, "xx"},
		{2, "cc"},
		{3, "ff"},
		{4, "zz"},
	}.
		Iter().
		SortBy(
			func(a, b g.Pair[int, string]) bool {
				return a.Key < b.Key
				// By value
				// return a.Value < b.Value
			}).
		Collect().
		Print()

	// Example 2: Sort a slice of custom structs and print the result
	type status struct {
		date   time.Time
		name   g.String
		status g.String
	}

	s1 := status{time.Now(), "s1", "good"}
	s2 := status{time.Now(), "s2", "bad"}
	s3 := status{time.Now().Add(time.Second * 10), "s3", "bad"}

	g.SliceOf(s3, s1, s2).Iter().
		SortBy(
			func(a, b status) bool {
				astatus := 5
				switch a.status {
				case "good":
					astatus = 0
				case "bad":
					astatus = 1
				}

				bstatus := 5
				switch b.status {
				case "good":
					bstatus = 0
				case "bad":
					bstatus = 1
				}

				return astatus < bstatus || astatus == bstatus && a.date.Unix() < b.date.Unix()
			}).
		Collect().
		Print()

	// Example 3: Sort a slice of time.Time, deduplicate, and print the result
	g.SliceOf(time.Now().Add(time.Second*20), time.Now()).
		Iter().
		SortBy(func(a, b time.Time) bool { return a.Second() < b.Second() }).
		Dedup().
		Collect().
		Print()

	// Example 4: Sort and deduplicate a slice of integers and print the result
	g.SliceOf(9, 8, 9, 8, 0, 1, 1, 1, 2, 7, 2, 2, 2, 3, 4, 5).
		Iter().
		Sort().
		Dedup().
		Filter(filters.IsOdd).
		Collect().
		Print()

	// Example 5: Sort a slice of strings in descending order and print the result
	g.SliceOf("a", "c", "b").
		Iter().
		SortBy(func(a, b string) bool { return a > b }).
		Collect().
		Print()
}
