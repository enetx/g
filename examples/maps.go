package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/pkg/ref"
)

func main() {
	gos := NewMap[string, *int]()

	gos.GetOrSet("root", ref.Of(3))
	fmt.Println(*gos.Get("root").Some() == 3)

	*gos.GetOrSet("root", ref.Of(10)) *= 2
	fmt.Println(*gos.Get("root").Some() == 6)

	//////////////////////////////////////////////////////////////////////////

	gos2 := NewMap[int, Slice[int]]()

	for i := range 5 {
		gos2.Set(i, gos2.Get(i).Some().Append(i))
	}

	for i := range 10 {
		gos2.Set(i, gos2.Get(i).Some().Append(i))
	}

	gos2.Print()

	//////////////////////////////////////////////////////////////////////////

	god := NewMap[int, Slice[int]]()

	for i := range 10 {
		god[i] = god.Get(i).Some().Append(i)
	}

	for i := range 10 {
		god[i] = god.Get(i).Some().Append(i)
	}

	god.Print()
}
