package main

import (
	. "github.com/enetx/g"
)

func main() {
	// var b Builder
	// b := new(Builder)

	// b.WriteString("builder\n")

	b := String("builder\n").Builder()

	for range 10 {
		b.WriteString("a")
		b.WriteRune('b')
		b.WriteByte('c')
		b.WriteString("\n")
	}

	b.String().Println()
}
