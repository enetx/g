package main

import (
	"fmt"
	"sync/atomic"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/pkg/iter"
)

func main() {
	var cnt uint32

	slice := genSlice()

	slice.ForEachParallel(func(i int) {
		atomic.AddUint32(&cnt, uint32(i))
	})

	fmt.Println(cnt) // 704982704

	cnt = 0

	mp := genMap()

	mp.ForEachParallel(func(i int, s g.String) {
		atomic.AddUint32(&cnt, uint32(i))
	})

	fmt.Println(cnt) // 704982704
}

func genSlice() g.Slice[int] {
	slice := g.NewSlice[int](0, 100000)
	for i := range iter.N(100000) {
		slice = slice.Append(i)
	}

	return slice
}

func genMap() g.Map[int, g.String] {
	mp := g.NewMap[int, g.String](100000)
	for i := range iter.N(100000) {
		mp.Set(i, g.NewInt(i).ToString())
	}

	return mp
}
