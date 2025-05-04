package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/pkg/ref"
)

func main() {
	gos := NewMap[string, *Int]()
	gos.GetOrSet("root", ref.Of(Int(3)))
	fmt.Println(*gos.Get("root").Some() == 3)

	*gos.GetOrSet("root", ref.Of(Int(10))) *= 2
	fmt.Println(*gos.Get("root").Some() == 6)

	//////////////////////////////////////////////////////////////////////////

	gos2 := NewMap[int, Slice[int]]()

	for i := range 5 {
		gos2.Set(i, gos2.Get(i).Some().Append(i))
	}

	for i := range 10 {
		gos2.Set(i, gos2.Get(i).Some().Append(i))
	}

	gos2.Println()

	//////////////////////////////////////////////////////////////////////////

	god := NewMap[int, Slice[int]]()

	for i := range 10 {
		god[i] = god.Get(i).Some().Append(i)
	}

	for i := range 10 {
		god[i] = god.Get(i).Some().Append(i)
	}

	god.Println()
}
