package main

import (
	"fmt"
	"time"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/filters"
)

func main() {
	// ========================================================================

	g.NewMapOrd[int, string]().
		Set(6, "bb").
		Set(0, "dd").
		Set(1, "aa").
		Set(5, "xx").
		Set(2, "cc").
		Set(3, "ff").
		Set(4, "zz").
		Iter().
		SortBy(
			func(a, b g.Pair[int, string]) bool {
				return a.Key < b.Key
				// return a.Value < b.Value
			}).
		Collect().
		Print()

	// ========================================================================

	type status struct {
		name   g.String
		status g.String
		date   time.Time
	}

	s1 := status{"s1", "good", time.Now()}
	s2 := status{"s2", "bad", time.Now()}
	s3 := status{"s3", "bad", time.Now().Add(time.Second * 10)}

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

	g.SliceOf(time.Now().Add(time.Second*20), time.Now()).
		Iter().
		SortBy(func(a, b time.Time) bool { return a.Second() < b.Second() }).
		Dedup().
		Collect().
		Print()

	g.SliceOf(9, 8, 9, 8, 0, 1, 1, 1, 2, 7, 2, 2, 2, 3, 4, 5).
		Iter().
		Sort().
		Dedup().
		Collect().
		Print()

	g.SliceOf("a", "c", "b").
		Iter().
		SortBy(func(a, b string) bool { return a > b }).
		Collect().
		Print()

	// ========================================================================

	is := g.SliceOf[g.Int](1, 2, 3, 4, 5)
	itos := g.TransformSlice(is, g.Int.ToString)

	itos.Iter().
		Fold("0", func(acc, val g.String) g.String { return g.Sprintf("(%s + %s)", acc, val) }).
		Print()

	// ========================================================================

	g.SliceOf(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
		Iter().
		StepBy(3).
		Map(
			func(i int) int {
				return i * i
			}).
		Inspect(
			func(i int) {
				fmt.Println(i)
			}).
		Collect().
		Print()

	// ========================================================================

	windows := g.SliceOf(1, 2, 3, 4).
		Iter().
		Windows(2).
		Collect()

	fmt.Println(windows)

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

	g.NewMapOrd[int, string]().
		Set(0, "aa").
		Set(1, "bb").
		Set(2, "cc").
		Set(3, "dd").
		Set(4, "ee").
		Set(5, "ff").
		Set(6, "gg").
		Iter().
		StepBy(2).
		Exclude(func(k int, v string) bool { return k == 4 }).
		Inspect(func(k int, v string) { fmt.Println(k, v) }).
		Map(func(k int, v string) (int, string) { return k, v + v }).
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
