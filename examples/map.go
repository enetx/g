package main

import (
	. "github.com/enetx/g"
)

func main() {
	gos2 := NewMap[int, Slice[int]]()

	for i := range 5 {
		gos2.Set(i, gos2.Get(i).UnwrapOrDefault().Append(i))
	}

	for i := range 10 {
		gos2.Set(i, gos2.Get(i).UnwrapOrDefault().Append(i))
	}

	gos2.Println()

	//////////////////////////////////////////////////////////////////////////

	god := NewMap[int, Slice[int]]()

	for i := range 10 {
		god[i] = god.Get(i).UnwrapOrDefault().Append(i)
	}

	for i := range 10 {
		god[i] = god.Get(i).UnwrapOrDefault().Append(i)
	}

	god.Println()
}
