package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	p1 := SliceOf("bbb", "ddd")

	contains := p1.Iter().Find(f.Eq("bbb")).IsSome()
	fmt.Println(contains)

	p2 := SliceOf("bbb", "yyy")

	containsAll := p2.Iter().All(func(v string) bool { return p1.Contains(v) })
	fmt.Println(containsAll)

	containsAny := p2.Iter().Any(func(v string) bool { return p1.Contains(v) })
	fmt.Println(containsAny)
}
