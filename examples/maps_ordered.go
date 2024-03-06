package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/ref"
)

func main() {
	// var mo g.MapOrd[int, string]
	// mo := g.MapOrd[int, string]{}
	// mo := make(g.MapOrd[int, string], 0)

	gos := g.NewMapOrd[int, *g.Slice[int]]()

	for i := range 5 {
		gos.GetOrSet(i, ref.Of(g.NewSlice[int]())).AppendInPlace(i)
	}

	for i := range 10 {
		gos.GetOrSet(i, ref.Of(g.NewSlice[int]())).AppendInPlace(i)
	}

	gos.Print()

	//////////////////////////////////////////////////////////////////////////

	for _, m := range gos {
		fmt.Println(m.Key, m.Value)
	}

	//////////////////////////////////////////////////////////////////////////

	god := g.NewMapOrd[int, g.Slice[int]]()

	for i := range 5 {
		god.Set(i, god.GetOrDefault(i, g.NewSlice[int]()).Append(i))
	}

	for i := range 10 {
		god.Set(i, god.GetOrDefault(i, g.NewSlice[int]()).Append(i))
	}

	god.Print()

	//////////////////////////////////////////////////////////////////////////

	ms := g.NewMapOrd[g.Int, g.Int]()
	ms.
		Set(11, 99).
		Set(12, 2).
		Set(1, 22).
		Set(2, 32).
		Set(222, 2)

	ms1 := ms.Clone()
	ms1.
		Set(888, 000).
		Set(888, 300)

	if v, ok := ms1.Get(888); ok {
		v.Print()
	}

	if v, ok := ms1.Get(11); ok {
		v.Print()
	}

	ms1.Set(1, 223)

	ms.Print()
	ms1.Print()

	fmt.Println(ms.Eq(ms1))
	fmt.Println(ms.Contains(12))

	ms.Iter().ForEach(func(k, v g.Int) { fmt.Println(k, v) })

	ms = ms.Iter().Map(func(k, v g.Int) (g.Int, g.Int) { return k.Mul(2), v.Mul(2) }).Collect()
	ms.Print()

	ms.Delete(22)
	fmt.Println(ms.Contains(22))

	msstr := g.NewMapOrd[g.String, g.String]()
	msstr.
		Set("aaa", "CCC").
		Set("ccc", "AAA").
		Set("bbb", "DDD").
		Set("ddd", "BBB")

	fmt.Println("before sort:", msstr)

	msstr.SortBy(func(a, b g.Pair[g.String, g.String]) bool { return a.Key < b.Key })
	fmt.Println("after sort:", msstr)

	msstr.SortBy(func(a, b g.Pair[g.String, g.String]) bool { return a.Value < b.Value })
	fmt.Println("after sort by value:", msstr)

	mss := g.NewMapOrd[g.Int, g.Slice[int]]()
	mss.Set(22, g.Slice[int]{4, 0, 9, 6, 7})
	mss.Set(11, g.Slice[int]{1, 2, 3, 4})

	fmt.Println("before sort: ", mss)

	mss.SortBy(func(a, b g.Pair[g.Int, g.Slice[int]]) bool { return a.Key < b.Key })
	fmt.Println("after sort by key: ", mss)

	mss.SortBy(func(a, b g.Pair[g.Int, g.Slice[int]]) bool { return a.Value[1] < b.Value[1] })
	fmt.Println("after sort by second value: ", mss)

	// mss.Iter().Keys().Collect().Print()
	// mss.Iter().Values().Collect().Print()

	// g.MapOrdFromStd(mss.ToMap().Std()).Print()
}
