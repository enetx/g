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
		value := gos2.Get(i).Some() // get the value at index i, or the default if it does not exist
		value.Push(i)
		gos2.Set(i, value)
	}

	for i := range 10 {
		value := gos2.Get(i).Some()
		value.Push(i)
		gos2.Set(i, value)
	}

	gos2.Println()

	//////////////////////////////////////////////////////////////////////////

	god := NewMap[int, Slice[int]]()

	for i := range 5 {
		value := god.Get(i).Some()
		value.Push(i)
		god[i] = value
	}

	for i := range 10 {
		value := god.Get(i).Some()
		value.Push(i)
		god[i] = value
	}

	god.Println()
}
