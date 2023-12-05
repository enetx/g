package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/iter"
	"gitlab.com/x0xO/g/pkg/ref"
)

func main() {
	gos := g.NewMap[string, *int]()

	gos.GetOrSet("root", ref.Of(3))
	fmt.Println(*gos.Get("root") == 3)

	*gos.GetOrSet("root", ref.Of(10)) *= 2
	fmt.Println(*gos.Get("root") == 6)

	// //////////////////////////////////////////////////////////////////////////

	gos2 := g.NewMap[int, *g.Slice[int]]()

	for i := range iter.N(5) {
		gos2.GetOrSet(i, ref.Of(g.NewSlice[int]())).AppendInPlace(i)
	}

	for i := range iter.N(10) {
		gos2.GetOrSet(i, ref.Of(g.NewSlice[int]())).AppendInPlace(i)
	}

	gos2.Print()

	// //////////////////////////////////////////////////////////////////////////

	god := g.NewMap[int, g.Slice[int]]()

	for i := range iter.N(5) {
		// god[i] = god.GetOrDefault(i, g.NewSlice[int]()).Append(i)
		god.Set(i, god.GetOrDefault(i, g.NewSlice[int]()).Append(i))
		// god.Set(i, god.Get(i).Append(i))
	}

	for i := range iter.N(10) {
		// god[i] = god.GetOrDefault(i, g.NewSlice[int]()).Append(i)
		god.Set(i, god.GetOrDefault(i, g.NewSlice[int]()).Append(i))
		// god.Set(i, god.Get(i).Append(i))
	}

	for i := range iter.N(10) {
		// god[i] = god.GetOrDefault(i, g.NewSlice[int]()).Append(i)
		god.Set(i, god.GetOrDefault(i, g.NewSlice[int]()).Append(i))
	}

	god.Print()
}
