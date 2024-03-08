package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/ref"
)

func main() {
	gos := g.NewMap[string, *int]()

	gos.GetOrSet("root", ref.Of(3))
	fmt.Println(*gos.Get("root").Some() == 3)

	*gos.GetOrSet("root", ref.Of(10)) *= 2
	fmt.Println(*gos.Get("root").Some() == 6)

	//////////////////////////////////////////////////////////////////////////

	gos2 := g.NewMap[int, g.Slice[int]]()

	for i := range 5 {
		gos2.Set(i, gos2.Get(i).UnwrapOrDefault().Append(i))
	}

	for i := range 10 {
		gos2.Set(i, gos2.Get(i).UnwrapOrDefault().Append(i))
	}

	gos2.Print()

	//////////////////////////////////////////////////////////////////////////

	god := g.NewMap[int, g.Slice[int]]()

	for i := range 10 {
		god[i] = god.Get(i).UnwrapOrDefault().Append(i)
	}

	for i := range 10 {
		god[i] = god.Get(i).UnwrapOrDefault().Append(i)
	}

	god.Print()
}
