package main

import (
	. "github.com/enetx/g"
)

func main() {
	Range('a', 'z').ForEach(func(v rune) {
		Println("{}", v)
	})

	RangeInclusive(0, 5).ForEach(func(v int) {
		Println("{}", v)
	})
}
