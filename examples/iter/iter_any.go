package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	values := Slice[any]{
		123,
		"hello",
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		struct{ X int }{X: 10},
		func() {},
		true,
		3.14,
		Int(11),
		String("hello"),
	}

	fmt.Println("Comparable values:")

	values.Iter().Filter(f.IsComparable).
		ForEach(func(v any) {
			Println("{1.type} - {1.debug}", v)
		})
}
