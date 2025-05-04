package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/pkg/ref"
)

func main() {
	// var mo g.MapOrd[int, string]
	// mo := g.MapOrd[int, string]{}
	// mo := make(g.MapOrd[int, string], 0)

	gos := NewMapOrd[int, *Slice[int]]()

	for i := range 5 {
		ns := NewSlice[int]()
		gos.GetOrSet(i, ref.Of(ns)).Push(i)
	}

	for i := range 10 {
		ns := NewSlice[int]()
		gos.GetOrSet(i, ref.Of(ns)).Push(i)
	}

	gos.Println()
	gos.Shuffle()
	gos.Println()

	//////////////////////////////////////////////////////////////////////////

	for _, m := range gos {
		fmt.Println(m.Key, m.Value)
	}

	//////////////////////////////////////////////////////////////////////////

	god := NewMapOrd[int, Slice[int]]()

	for i := range 5 {
		god.Set(i, god.Get(i).Some().Append(i))
	}

	for i := range 10 {
		god.Set(i, god.Get(i).Some().Append(i))
	}

	god.Println()

	//////////////////////////////////////////////////////////////////////////

	ms := NewMapOrd[Int, Int]()
	ms.Set(11, 99)
	ms.Set(12, 2)
	ms.Set(1, 22)
	ms.Set(2, 32)
	ms.Set(222, 2)

	ms1 := ms.Clone()
	ms1.Set(888, 0o00)
	ms1.Set(888, 300)

	if r := ms1.Get(888); r.IsSome() {
		r.Some().Println()
	}

	if r := ms1.Get(11); r.IsSome() {
		r.Some().Println()
	}

	ms1.Set(1, 223)

	ms.Println()
	ms1.Println()

	fmt.Println(ms.Eq(ms1))
	fmt.Println(ms.Contains(12))

	ms.Iter().ForEach(func(k, v Int) { fmt.Println(k, v) })

	ms = ms.Iter().Map(func(k, v Int) (Int, Int) { return k.Mul(2), v.Mul(2) }).Collect()
	ms.Println()

	ms.Delete(22)
	fmt.Println(ms.Contains(22))

	msstr := NewMapOrd[String, String]()
	msstr.Set("aaa", "CCC")
	msstr.Set("ccc", "AAA")
	msstr.Set("bbb", "DDD")
	msstr.Set("ddd", "BBB")

	fmt.Println("before sort:", msstr)

	// msstr.SortBy(func(a, b Pair[String, String]) cmp.Ordering { return a.Key.Cmp(b.Key) })
	msstr.SortByKey(func(a, b String) cmp.Ordering { return a.Cmp(b) })
	fmt.Println("after sort:", msstr)

	// msstr.SortBy(func(a, b Pair[String, String]) cmp.Ordering { return a.Value.Cmp(b.Value) })
	msstr.SortByValue(String.Cmp)
	fmt.Println("after sort by value:", msstr)

	mss := NewMapOrd[Int, Slice[int]]()
	mss.Set(22, Slice[int]{4, 0, 9, 6, 7})
	mss.Set(11, Slice[int]{1, 2, 3, 4})

	fmt.Println("before sort: ", mss)
	// mss.SortBy(func(a, b Pair[Int, Slice[int]]) cmp.Ordering { return a.Key.Cmp(b.Key) })
	mss.SortByKey(Int.Cmp)
	fmt.Println("after sort by key: ", mss)

	// mss.SortBy(func(a, b g.Pair[g.Int, g.Slice[int]]) cmp.Ordering { return cmp.Cmp(a.Value[1], b.Value[1]) })
	mss.SortByValue(func(a, b Slice[int]) cmp.Ordering { return cmp.Cmp(a[1], b[1]) })
	fmt.Println("after sort by second value: ", mss)

	// MapOrdFromStd(mss.ToMap().Std()).Println()
}
