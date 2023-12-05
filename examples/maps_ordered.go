package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/iter"
	"gitlab.com/x0xO/g/pkg/ref"
)

func main() {
	// var mo *g.MapOrd[int, string]
	// mo := &g.MapOrd[int, string]{}
	// mo := new(g.MapOrd[int, string])
	// mo := ref.Of(make(g.MapOrd[int, string], 0))

	gos := g.NewMapOrd[int, *g.Slice[int]]()

	for i := range iter.N(5) {
		gos.GetOrSet(i, ref.Of(g.NewSlice[int]())).AppendInPlace(i)
	}

	for i := range iter.N(10) {
		gos.GetOrSet(i, ref.Of(g.NewSlice[int]())).AppendInPlace(i)
	}

	gos.Print()

	//////////////////////////////////////////////////////////////////////////

	for _, m := range *gos {
		fmt.Println(m.Key, m.Value)
	}

	//////////////////////////////////////////////////////////////////////////

	god := g.NewMapOrd[int, g.Slice[int]]()

	for i := range iter.N(5) {
		god.Set(i, god.GetOrDefault(i, g.NewSlice[int]()).Append(i))
	}

	for i := range iter.N(10) {
		god.Set(i, god.GetOrDefault(i, g.NewSlice[int]()).Append(i))
	}

	god.Print()

	//////////////////////////////////////////////////////////////////////////

	ms := g.NewMapOrd[g.Int, g.Int]()
	ms.Set(11, 99).Set(12, 2).Set(1, 22).Set(2, 32).Set(222, 2)

	ms1 := ms.Clone()

	ms1.Set(888, 000)
	ms1.Set(888, 300)

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

	ms.ForEach(func(k, v g.Int) { fmt.Println(k, v) })

	ms = ms.Map(func(k, v g.Int) (g.Int, g.Int) { return k.Mul(2), v.Mul(2) })

	ms.Print()

	ms.Delete(12, 1, 222)
	fmt.Println(ms.Contains(12))

	msstr := g.NewMapOrd[g.String, g.String]()
	msstr.Set("aaa", "CCC").Set("ccc", "AAA").Set("bbb", "DDD").Set("ddd", "BBB")
	msstr.Print() // before sort

	msstr.SortBy(func(i, j int) bool { return (*msstr)[i].Key < (*msstr)[j].Key })
	msstr.Print() // after sort by key

	msstr.SortBy(func(i, j int) bool { return (*msstr)[i].Value < (*msstr)[j].Value })

	msstr.Print() // after sort by value

	mss := g.NewMapOrd[g.Int, g.Slice[int]]()
	mss.Set(22, g.Slice[int]{4, 0, 9, 6, 7})
	mss.Set(11, g.Slice[int]{1, 2, 3, 4})
	mss.Print() // before sort

	mss.SortBy(func(i, j int) bool { return (*mss)[i].Key < (*mss)[j].Key })
	mss.Print() // after sort by key

	mss.SortBy(func(i, j int) bool { return (*mss)[i].Value.Get(1) < (*mss)[j].Value.Get(1) })
	mss.Print() // after sort by value

	g.MapOrdFromStd(mss.ToMap().Std()).Print()
}
