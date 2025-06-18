package main

import (
	"unsafe"

	. "github.com/enetx/g"
)

func main() {
	a := [][]string{{"aaa", "bbb"}, {"ccc"}}

	b := *(*(Slice[Slice[String]]))(unsafe.Pointer(&a))
	c := *(*(Slice[Slice[string]]))(unsafe.Pointer(&a))

	b[len(b)-1].Push("ddd")

	Println("{}", a)
	b.Println()
	c.Println()

	// [[aaa bbb] [ccc ddd]]
	// Slice[Slice[aaa, bbb], Slice[ccc, ddd]]
	// Slice[Slice[aaa, bbb], Slice[ccc, ddd]]
}
