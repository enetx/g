package main

import (
	. "github.com/enetx/g"
)

func main() {
	builder := NewString("builder\n").Builder()

	for range 10 {
		builder.
			Write("a").
			WriteRune('b').
			WriteByte('c').
			Write("\n")
	}

	builder.String().Print()
}
